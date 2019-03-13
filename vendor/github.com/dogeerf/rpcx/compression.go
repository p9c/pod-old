package rpcx

import (
	"compress/flate"
	"fmt"
	"io"
	"net"

	"github.com/DataDog/zstd"
	"github.com/golang/snappy"
	"github.com/pierrec/lz4"
)

// CompressType is compression type. Currently only support zip and snappy
type CompressType byte

const (
	// CompressNone represents no compression
	CompressNone CompressType = iota
	// CompressFlate represents zip
	CompressFlate
	// CompressSnappy represents snappy
	CompressSnappy
	// CompressZstd represents Facebook/Zstandard
	CompressZstd
	// CompressLZ4 represents LZ4 (http://www.lz4.org)
	CompressLZ4
)

type writeFlusher struct {
	w *flate.Writer
}

func (wf *writeFlusher) Write(p []byte) (int, error) {
	n, err := wf.w.Write(p)
	if err != nil {
		return n, err
	}
	if err := wf.w.Flush(); err != nil {
		return 0, err
	}
	return n, nil
}

// CompressConn wraps a net.Conn and supports compression
type CompressConn struct {
	net.Conn
	r            io.Reader
	w            io.Writer
	compressType CompressType
}

// NewCompressConn creates a wrapped net.Conn with CompressType
func NewCompressConn(conn net.Conn, compressType CompressType) net.Conn {
	cc := &CompressConn{Conn: conn, compressType: compressType}
	r := io.Reader(cc.Conn)

	switch compressType {
	case CompressNone:
	case CompressFlate:
		r = flate.NewReader(r)
	case CompressSnappy:
		r = snappy.NewReader(r)
	case CompressZstd:
		r = zstd.NewReader(r)
	case CompressLZ4:
		r = lz4.NewReader(r)
	}
	cc.r = r

	w := io.Writer(cc.Conn)
	switch compressType {
	case CompressNone:
	case CompressFlate:
		zw, err := flate.NewWriter(w, flate.DefaultCompression)
		if err != nil {
			panic(fmt.Sprintf("BUG: flate.NewWriter(%d) returned non-nil err: %s", flate.DefaultCompression, err))
		}
		w = &writeFlusher{w: zw}
	case CompressSnappy:
		w = snappy.NewBufferedWriter(w)
	case CompressZstd:
		w = zstd.NewWriter(w)
	case CompressLZ4:
		w = lz4.NewWriter(w)
	}
	cc.w = w
	return cc
}

func (c *CompressConn) Read(b []byte) (n int, err error) {
	return c.r.Read(b)
}

func (c *CompressConn) Write(b []byte) (n int, err error) {
	return c.w.Write(b)
}
