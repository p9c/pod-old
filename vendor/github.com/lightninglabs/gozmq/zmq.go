// Package gozmq provides a ZMQ pubsub client.
//
// It implements the protocol described here:
// http://rfc.zeromq.org/spec:23/ZMTP/
package gozmq

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
	"strings"
	"time"
)

var (
	// MaxBodySize is the maximum frame size we're willing to accept.
	// The maximum size of a Bitcoin message is 32MiB.
	MaxBodySize uint64 = 0x02000000
)

// reconnectError wraps any error with a timeout, which allows clients to keep
// attempting to read normally while the underlying socket attempts to
// reconnect.
type reconnectError struct {
	error
}

// Timeout implements the timeout interface from net and other packages.
func (e *reconnectError) Timeout() bool {
	return true
}

func connFromAddr(addr string) (net.Conn, error) {
	re, err := regexp.Compile(`((tcp|unix|ipc)://)?([^:]*):?(\d*)`)
	if err != nil {
		return nil, err
	}

	submatch := re.FindStringSubmatch(addr)

	if addrs, err := net.LookupIP(submatch[3]); (err == nil &&
		len(addrs) > 0) || submatch[2] == "tcp" ||
		net.ParseIP(submatch[3]) != nil || submatch[4] != "" {

		// We have a TCP address.
		return net.DialTimeout("tcp", strings.Replace(submatch[3], "*",
			"", 1)+":"+submatch[4], time.Minute)

	} else if _, err := os.Stat(submatch[3]); err == nil ||
		submatch[2] == "unix" || submatch[2] == "ipc" {

		// We have a UNIX socket.
		return net.DialTimeout("unix", submatch[3], time.Minute)
	}

	return nil, fmt.Errorf("couldn't resolve address %s", addr)
}

// Conn is a connection to a ZMQ server.
type Conn struct {
	conn    net.Conn
	topics  []string
	timeout time.Duration
}

func (c *Conn) writeAll(buf []byte) error {
	written := 0
	for written < len(buf) {
		n, err := c.conn.Write(buf[written:])
		if err != nil {
			return err
		}
		written += n
	}
	return nil
}

func (c *Conn) writeGreeting() error {
	signature := []byte{0xff, 0, 0, 0, 0, 0, 0, 0, 0, 0x7f}
	version := []byte{3, 0}
	mechanism := make([]byte, 20)
	for i, chr := range "NULL" {
		mechanism[i] = byte(chr)
	}
	server := []byte{0}

	var greeting []byte
	greeting = append(greeting, signature...)
	greeting = append(greeting, version...)
	greeting = append(greeting, mechanism...)
	greeting = append(greeting, server...)
	for len(greeting) < 64 {
		greeting = append(greeting, 0)
	}

	return c.writeAll(greeting)
}

func (c *Conn) readGreeting() error {
	greet := make([]byte, 64)
	if _, err := io.ReadFull(c.conn, greet); err != nil {
		return err
	}

	if greet[0] != 0xff || greet[9] != 0x7f {
		return errors.New("invalid signature")
	}
	if greet[10] < 3 {
		return errors.New("peer version is too old")
	}
	if string(greet[12:17]) != "NULL\x00" {
		return errors.New("unsupported security mechanism")
	}
	if greet[32] != 0 {
		return errors.New("as-server must be zero for NULL")
	}

	return nil
}

func (c *Conn) writeFrame(flag byte, buf []byte) error {
	if flag&0xf8 != 0 {
		return errors.New("invalid flag")
	}

	if flag&2 != 0 {
		return errors.New("caller must not specify long frame flag")
	}

	var header []byte

	if len(buf) > 255 {
		flag = flag | 2
		header = []byte{flag, 0, 0, 0, 0, 0, 0, 0, 0}
		size := len(buf)
		i := len(header) - 1
		for size > 0 {
			header[i] = byte(size & 0xff)
			size = size >> 8
			i--
		}
	} else {
		header = []byte{flag, byte(len(buf))}
	}

	if err := c.writeAll(header); err != nil {
		return err
	}
	return c.writeAll(buf)
}

func (c *Conn) writeCommand(name string, data []byte) error {
	size := len(name)
	if size > 255 {
		return errors.New("command name is too long")
	}
	body := append([]byte{byte(size)}, []byte(name)...)
	buf := append(body, data...)
	var flag byte = 4
	return c.writeFrame(flag, buf)
}

func (c *Conn) writeReady(socketType string) error {
	if len(socketType) > 255 {
		return errors.New("socket type too long")
	}
	const socketTypeName = "Socket-Type"
	metadata := []byte{byte(len(socketTypeName))}
	metadata = append(metadata, []byte(socketTypeName)...)
	metadata = append(metadata, []byte{0, 0, 0, byte(len(socketType))}...)
	metadata = append(metadata, []byte(socketType)...)
	return c.writeCommand(commandReady, metadata)
}

func (c *Conn) writeMessage(parts [][]byte) error {
	if len(parts) == 0 {
		return errors.New("empty message")
	}
	for _, msg := range parts[:len(parts)-1] {
		if err := c.writeFrame(1, msg); err != nil {
			return err
		}
	}
	return c.writeFrame(0, parts[len(parts)-1])
}

func (c *Conn) subscribe(prefix string) error {
	msg := append([]byte{1}, []byte(prefix)...)
	return c.writeMessage([][]byte{msg})
}

func (c *Conn) readCommand() (string, []byte, error) {
	flag, buf, err := c.readFrame()
	if err != nil {
		return "", nil, err
	}
	if flag&4 != 4 {
		return "", nil, errors.New("expected command frame")
	}
	if len(buf) < 1 {
		return "", nil, errors.New("empty command buffer")
	}
	size := int(buf[0])
	buf = buf[1:]
	if size > len(buf) {
		return "", nil, errors.New("invalid command name size")
	}
	name := string(buf[:size])
	data := buf[size:]
	return name, data, nil
}

const commandReady = "READY"

func (c *Conn) readReady() error {
	name, metadata, err := c.readCommand()
	if err != nil {
		return err
	}
	if name != commandReady {
		return errors.New("expected ready command")
	}

	m := make(map[string]string)
	for len(metadata) > 0 {
		size := int(metadata[0])
		metadata = metadata[1:]
		if size > len(metadata) {
			return errors.New("invalid metadata")
		}
		name := metadata[:size]
		metadata = metadata[size:]

		if len(metadata) < 4 {
			return errors.New("invalid metadata")
		}
		var valueSize uint32
		for i := 0; i < 4; i++ {
			valueSize = valueSize<<8 + uint32(metadata[i])
		}
		metadata = metadata[4:]

		if int(valueSize) > len(metadata) {
			return errors.New("invalid metadata")
		}
		value := metadata[:valueSize]
		metadata = metadata[valueSize:]

		m[string(name)] = string(value)
	}

	return nil
}

// Read a frame from the socket, setting deadline before each read to prevent
// timeouts during or between frames.
func (c *Conn) readFrame() (byte, []byte, error) {
	var flagBuf [1]byte
	c.conn.SetReadDeadline(time.Now().Add(c.timeout))
	if _, err := io.ReadFull(c.conn, flagBuf[:1]); err != nil {
		return 0, nil, err
	}

	flag := flagBuf[0]
	if flag&0xf8 != 0 {
		return 0, nil, errors.New("invalid flag")
	}

	var size uint64

	if flag&2 == 2 {
		// Long form
		var buf [8]byte
		c.conn.SetReadDeadline(time.Now().Add(c.timeout))
		if _, err := io.ReadFull(c.conn, buf[:8]); err != nil {
			return 0, nil, err
		}
		for _, b := range buf {
			size = (size << 8) | uint64(b)
		}
	} else {
		// Short form
		var buf [1]byte
		c.conn.SetReadDeadline(time.Now().Add(c.timeout))
		if _, err := io.ReadFull(c.conn, buf[:1]); err != nil {
			return 0, nil, err
		}
		size = uint64(buf[0])
	}

	if size > MaxBodySize {
		return 0, nil, errors.New("frame too large")
	}

	buf := make([]byte, size)
	// Prevent timeout during large data read in case of slow connection.
	c.conn.SetReadDeadline(time.Time{})
	if _, err := io.ReadFull(c.conn, buf); err != nil {
		return 0, nil, err
	}

	return flag, buf, nil
}

// Read a message from the socket.
func (c *Conn) readMessage() ([][]byte, error) {
	var parts [][]byte
	for {
		flag, buf, err := c.readFrame()
		if err != nil {
			return nil, err
		}
		if flag&4 != 0 {
			return nil, errors.New("expected message frame")
		}

		parts = append(parts, buf)

		if flag&1 == 0 {
			break
		}

		if len(parts) > 16 {
			return nil, errors.New("message has too many parts")
		}
	}
	return parts, nil
}

// Subscribe connects to a publisher server and subscribes to the given topics.
func Subscribe(addr string, topics []string, timeout time.Duration) (*Conn, error) {
	conn, err := connFromAddr(addr)
	if err != nil {
		return nil, err
	}

	conn.SetDeadline(time.Now().Add(10 * time.Second))

	c := &Conn{conn, topics, timeout}

	if err := c.writeGreeting(); err != nil {
		conn.Close()
		return nil, err
	}
	if err := c.readGreeting(); err != nil {
		conn.Close()
		return nil, err
	}

	if err := c.writeReady("SUB"); err != nil {
		conn.Close()
		return nil, err
	}

	if err := c.readReady(); err != nil {
		conn.Close()
		return nil, err
	}

	for _, topic := range topics {
		if err := c.subscribe(topic); err != nil {
			conn.Close()
			return nil, err
		}
	}

	conn.SetDeadline(time.Time{})

	return c, nil
}

// Receive a message from the publisher. It blocks until a new message is
// received.
func (c *Conn) Receive() ([][]byte, error) {
	messages, err := c.readMessage()
	// If the error is either nil or a non-EOF error, we return it as-is.
	if err != io.EOF {
		return messages, err
	}
	// We got an EOF, so our socket is disconnected. We attempt to
	// reconnect. If successful, replace the existing connection with the
	// new one. Either way, return a timeout error.
	errTimeout := &net.OpError{
		Op:     "read",
		Net:    c.conn.LocalAddr().Network(),
		Source: c.conn.LocalAddr(),
		Addr:   c.conn.RemoteAddr(),
		Err:    &reconnectError{err},
	}
	newConn, err := Subscribe(c.conn.RemoteAddr().String(), c.topics,
		c.timeout)
	if err != nil {
		// Prevent CPU overuse by refused reconnection attempts.
		time.Sleep(c.timeout)
	} else {
		c.Close()
		*c = *newConn
	}
	return nil, errTimeout
}

// Close the underlying connection. Any further operations will fail.
func (c *Conn) Close() error {
	return c.conn.Close()
}
