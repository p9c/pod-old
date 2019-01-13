package main

import (
	"fmt"
	"net"
	"time"

	l "github.com/parallelcointeam/sub/clog"
)

type nodeCfg struct {
	Listener   string
	Worker     string
	BufferSize int
}

type node struct {
	cfg        nodeCfg
	connection *net.UDPConn
	kill       chan bool
}

const (
	uNet = "udp4"
)

var (
	_n = nodeCfg{
		Listener:   "127.0.0.1:11011",
		Worker:     "127.0.0.1:11012",
		BufferSize: 10240,
	}
	_w = nodeCfg{
		Listener:   "127.0.0.1:11012",
		Worker:     "127.0.0.1:11011",
		BufferSize: 10240,
	}
)

func main() {
	l.Init()
	*ld <- "starting up"
	n := newNode(_n)

	n.setupListener()
	time.Sleep(time.Second * 1)
	w := newNode(_w)
	go n.readFromSocket()
	for {
		time.Sleep(time.Second)
		go w.send([]byte("hello world"))
	}
}

func newNode(nc nodeCfg) (n *node) {
	n = &node{
		cfg:  nc,
		kill: make(chan bool),
	}
	return
}

func (n *node) setupListener() {
	addr, err := net.ResolveUDPAddr(uNet, n.cfg.Listener)
	check(err)
	n.connection, err = net.ListenUDP(uNet, addr)
	check(err)
}

func (n *node) readFromSocket() {
	for {
		var b = make([]byte, n.cfg.BufferSize)
		count, addr, err := n.connection.ReadFromUDP(b[0:])
		check(err)
		b = b[:count]
		if count > 0 {
			*li <- fmt.Sprint("'", string(b), "' <- ", addr)
			select {
			case <-n.kill:
				*li <- "closing socket"
				break
			default:
			}
		}
	}
}

func (n *node) send(b []byte) {
	addr, err := net.ResolveUDPAddr("udp4", n.cfg.Worker)
	check(err)
	conn, err := net.DialUDP(uNet, nil, addr)
	check(err)
	_, err = conn.Write(b)
	check(err)
	*li <- "'" + string(b) + "' -> " + n.cfg.Worker
}

func check(err error) {
	if err != nil {
		*le <- err.Error()
	}
}
