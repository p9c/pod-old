package cl

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/mitchellh/colorstring"
)

// Og is the root channel that processes logging messages, so, cl.Og <- Fatalf{"format string %s %d", stringy, inty} sends to the root
var Og = make(chan interface{})

// Value is the generic list of things processed by the log chan
type Value []interface{}

// Fatal is a log value that indicates level and how to interpret the interface slice
type Fatal Value

// Fatalf is a log value that indicates level and how to interpret the interface slice
type Fatalf Value

// Ftl is a log type that is just one string
type Ftl string

// Error is a log value that indicates level and how to interpret the interface slice
type Error Value

// Errorf is a log value that indicates level and how to interpret the interface slice
type Errorf Value

// Err is a log type that is just one string
type Err string

// Warn is a log value that indicates level and how to interpret the interface slice
type Warn Value

// Warnf is a log value that indicates level and how to interpret the interface slice
type Warnf Value

// Wrn is a log type that is just one string
type Wrn string

// Info is a log value that indicates level and how to interpret the interface slice
type Info Value

// Infof is a log value that indicates level and how to interpret the interface slice
type Infof Value

// Inf is a log type that is just one string
type Inf string

// Debug is a log value that indicates level and how to interpret the interface slice
type Debug Value

// Debugf is a log value that indicates level and how to interpret the interface slice
type Debugf Value

// Dbg is a log type that is just one string
type Dbg string

// Trace is a log value that indicates level and how to interpret the interface slice
type Trace Value

// Tracef is a log value that indicates level and how to interpret the interface slice
type Tracef Value

// Trc is a log type that is just one string
type Trc string

// A SubSystem is a logger with a specific prefix name prepended  to the entry
type SubSystem struct {
	Ch          chan interface{}
	Level       int
	LevelString string
}

// Close a SubSystem logger
func (s *SubSystem) Close() {
	close(s.Ch)
}

// Writer is the place the
var Writer = io.MultiWriter(os.Stdout)

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

// SetLevel changes the level of a subsystem by level name
func (s *SubSystem) SetLevel(level string) {
	if i, ok := Levels[level]; ok {
		s.Level = i
		s.LevelString = level
	} else {
		s.Level = _off
		s.LevelString = "off"
	}
}

const errFmt = "ERR:FMT\n  "

// Color turns on and off colouring of error type tag
var Color bool

// ColorChan accepts a bool and flips the state accordingly
var ColorChan = make(chan bool)

// Started will block when the logger starts
var Started = make(chan bool)

// ShuttingDown indicates if the shutdown switch has been triggered
var ShuttingDown bool

// NewSubSystem starts up a new subsystem logger
func NewSubSystem(name, level string) (ss *SubSystem) {
	ss = new(SubSystem)
	ss.Ch = make(chan interface{})
	ss.SetLevel(level)
	Og <- Infof{"started subsystem '%s'", name}

	go func() {
		for {
			if ShuttingDown {
				break
			}
			select {
			case i := <-ss.Ch:
				if ShuttingDown {
					break
				}
				n := name
				if Color {
					n = colorstring.Color("[bold]" + n + "[reset]")
				} else {
					n += ":"
				}
				switch i.(type) {
				case Ftl:
					if ss.Level > _off {
						Og <- Ftl(n+" ") + i.(Ftl)
					}
				case Err:
					if ss.Level > _fatal {
						Og <- Err(n+" ") + i.(Err)
					}
				case Wrn:
					if ss.Level > _error {
						Og <- Wrn(n+" ") + i.(Wrn)
					}
				case Inf:
					if ss.Level > _warn {
						Og <- Inf(n+" ") + i.(Inf)
					}
				case Dbg:
					if ss.Level > _info {
						Og <- Dbg(n+" ") + i.(Dbg)
					}
				case Trc:
					if ss.Level > _debug {
						Og <- Trc(n+" ") + i.(Trc)
					}
				case Fatal:
					if ss.Level > _off {
						Og <- append(Fatal{n}, i.(Fatal)...)
					}
				case Error:
					if ss.Level > _fatal {
						Og <- append(Error{n}, i.(Error)...)
					}
				case Warn:
					if ss.Level > _error {
						Og <- append(Warn{n}, i.(Warn)...)
					}
				case Info:
					if ss.Level > _warn {
						Og <- append(Info{n}, i.(Info)...)
					}
				case Debug:
					if ss.Level > _info {
						Og <- append(Debug{n}, i.(Debug)...)
					}
				case Trace:
					if ss.Level > _debug {
						Og <- append(Trace{n}, i.(Trace)...)
					}
				case Fatalf:
					if ss.Level > _off {
						Og <- append(Fatalf{n + " " + i.(Fatalf)[0].(string)}, i.(Fatalf)[1:]...)
					}
				case Errorf:
					if ss.Level > _fatal {
						Og <- append(Errorf{n + " " + i.(Errorf)[0].(string)}, i.(Errorf)[1:]...)
					}
				case Warnf:
					if ss.Level > _warn {
						Og <- append(Warnf{n + " " + i.(Warnf)[0].(string)}, i.(Warnf)[1:]...)
					}
				case Infof:
					if ss.Level > _info {
						Og <- append(Infof{n + " " + i.(Infof)[0].(string)}, i.(Infof)[1:]...)
					}
				case Debugf:
					if ss.Level > _debug {
						Og <- append(Debugf{n + " " + i.(Debugf)[0].(string)}, i.(Debugf)[1:]...)
					}
				case Tracef:
					if ss.Level > _trace {
						Og <- append(Tracef{n + " " + i.(Tracef)[0].(string)}, i.(Tracef)[1:]...)
					}
				}
			}

		}
	}()
	return
}

func init() {
	fmt.Fprintln(os.Stderr, "starting clog")
	go func() {
		Started <- true
		for {
			var t, s string
			select {
			case <-Quit:
				ShuttingDown = true
				break
			case Color = <-ColorChan:
			case i := <-Og:
				color := Color
				if i == "" {
					continue
				}
				t = time.Now().Format("2006-01-02 15:04:05.000000 MST")
				if color {
					s = colorstring.Color("[reset]")
				}
				switch i.(type) {
				case Ftl:
					s += string(i.(Ftl)) + "\n"
				case Err:
					s += string(i.(Err)) + "\n"
				case Wrn:
					s += string(i.(Wrn)) + "\n"
				case Inf:
					s += string(i.(Inf)) + "\n"
				case Dbg:
					s += string(i.(Dbg)) + "\n"
				case Trc:
					s += string(i.(Trc)) + "\n"
				case Fatal:
					s += fmt.Sprintln(i.(Fatal)...)
				case Error:
					s += fmt.Sprintln(i.(Error)...)
				case Warn:
					s += fmt.Sprintln(i.(Warn)...)
				case Info:
					s += fmt.Sprintln(i.(Info)...)
				case Debug:
					s += fmt.Sprintln(i.(Debug)...)
				case Trace:
					s += fmt.Sprintln(i.(Trace)...)
				case Fatalf:
					I := i.(Fatalf)
					switch I[0].(type) {
					case string:
						s += fmt.Sprintf(I[0].(string), I[1:]...) + "\n"
					}
				case Errorf:
					I := i.(Errorf)
					switch I[0].(type) {
					case string:
						s += fmt.Sprintf(I[0].(string), I[1:]...) + "\n"
					}
				case Warnf:
					I := i.(Warnf)
					switch I[0].(type) {
					case string:
						s += fmt.Sprintf(I[0].(string), I[1:]...) + "\n"
					}
				case Infof:
					I := i.(Infof)
					switch I[0].(type) {
					case string:
						s += fmt.Sprintf(I[0].(string), I[1:]...) + "\n"
					}
				case Debugf:
					I := i.(Debugf)
					switch I[0].(type) {
					case string:
						s += fmt.Sprintf(I[0].(string), I[1:]...) + "\n"
					}
				case Tracef:
					I := i.(Tracef)
					switch I[0].(type) {
					case string:
						s += fmt.Sprintf(I[0].(string), I[1:]...) + "\n"
					}
				}
				switch i.(type) {
				case Ftl, Fatal, Fatalf:
					s = ftlTag(color) + s
				case Err, Error, Errorf:
					s = errTag(color) + s
				case Wrn, Warn, Warnf:
					s = wrnTag(color) + s
				case Inf, Info, Infof:
					s = infTag(color) + s
				case Dbg, Debug, Debugf:
					s = dbgTag(color) + s
				case Trc, Trace, Tracef:
					s = trcTag(color) + s
				}
				if color {
					t = colorstring.Color("[light_gray]" + t + "[dark_gray]")
				}
			}
			fmt.Fprint(Writer, t+s)
		}
	}()
}

func ftlTag(color bool) string {
	tag := "FTL"
	if color {
		pre := ""  // colorstring.Color("[light_gray][[dark_gray]")
		post := "" // colorstring.Color("[light_gray]][dark_gray]")
		tag = pre + colorstring.Color("[red]"+colorstring.Color("[bold]"+tag)) + post
	} else {
		tag = "[" + tag + "]"
	}
	return " " + tag + " "
}

func errTag(color bool) string {
	tag := "ERR"
	if color {
		pre := ""  // colorstring.Color("[light_gray][[dark_gray]")
		post := "" // colorstring.Color("[light_gray]][dark_gray]")
		tag = pre + colorstring.Color("[yellow]"+colorstring.Color("[bold]"+tag)) + post
	} else {
		tag = "[" + tag + "]"
	}
	return " " + tag + " "

}

func wrnTag(color bool) string {
	tag := "WRN"
	if color {
		pre := ""  // colorstring.Color("[light_gray][[dark_gray]")
		post := "" // colorstring.Color("[light_gray]][dark_gray]")
		tag = pre + colorstring.Color("[green]"+colorstring.Color("[bold]"+tag)) + post
	} else {
		tag = "[" + tag + "]"
	}
	return " " + tag + " "

}

func infTag(color bool) string {
	tag := "INF"
	if color {
		pre := ""  // colorstring.Color("[light_gray][[dark_gray]")
		post := "" // colorstring.Color("[light_gray]][dark_gray]")
		tag = pre + colorstring.Color("[cyan]"+colorstring.Color("[bold]"+tag)) + post
	} else {
		tag = "[" + tag + "]"
	}
	return " " + tag + " "

}
func dbgTag(color bool) string {
	tag := "DBG"
	if color {
		pre := ""  // colorstring.Color("[light_gray][[dark_gray]")
		post := "" // colorstring.Color("[light_gray]][dark_gray]")
		tag = pre + colorstring.Color("[blue]"+colorstring.Color("[bold]"+tag)) + post
	} else {
		tag = "[" + tag + "]"
	}
	return " " + tag + " "

}
func trcTag(color bool) string {
	tag := "TRC"
	if color {
		pre := ""  // colorstring.Color("[light_gray][[dark_gray]")
		post := "" // colorstring.Color("[light_gray]][dark_gray]")
		tag = pre + colorstring.Color("[magenta]"+colorstring.Color("[bold]"+tag)) + post
	} else {
		tag = "[" + tag + "]"
	}
	return " " + tag + " "

}

// Quit signals the logger to stop. Invoke like this:
//     close(clog.Quit)
// You can call Init again to start it up again
var Quit = make(chan struct{})

// Shutdown the application, allowing the logger a moment to clear the channels
func Shutdown() {
	close(Quit)
	// wait a moment to let log channel clear
	time.Sleep(time.Millisecond * 50)
	os.Exit(0)
}
