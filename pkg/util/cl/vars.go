package cl

import (
	"io"
	"os"
	"sync"
)

const (
	_off = iota
	_fatal
	_error
	_warn
	_info
	_debug
	_trace
)

// Levels is the map of name to level
var Levels = map[string]int{
	"off":   _off,
	"fatal": _fatal,
	"error": _error,
	"warn":  _warn,
	"info":  _info,
	"debug": _debug,
	"trace": _trace,
}

// const errFmt = "ERR:FMT\n  "

// Color turns on and off colouring of error type tag
var Color = true

// ColorChan accepts a bool and flips the state accordingly
var ColorChan = make(chan bool)

// ShuttingDown indicates if the shutdown switch has been triggered
var ShuttingDown bool

// Writer is the place thelogs put out
var Writer = io.MultiWriter(os.Stdout)

// Og is the root channel that processes logging messages, so, cl.Og <- Fatalf{"format string %s %d", stringy, inty} sends to the root
var Og = make(chan interface{})

var wg sync.WaitGroup

// Quit signals the logger to stop. Invoke like this:
//     close(clog.Quit)
// You can call init() again to start it up again
var Quit = make(chan struct{})

var Register = make(Registry)
