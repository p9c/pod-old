package clog

// Miscellaneous functions

import (
	"fmt"
	"os"
	"time"

	"github.com/logrusorgru/aurora"
)

// Joined is a channel that can be used to redirect the entire log output to another channel
var Joined chan string

// A SubSystem is a logger that intercepts a signal, adds a 'name' prefix and passes it to the main logger channel
type SubSystem struct {
	Fatal chan string
	Error chan string
	Warn  chan string
	Info  chan string
	Debug chan string
	Trace chan string
}

// NewSubSystem creates a new clog logger that adds a prefix to the log entry for subsystem control
func NewSubSystem(name string, level int) *SubSystem {
	ss := SubSystem{
		Fatal: make(chan string),
		Error: make(chan string),
		Warn:  make(chan string),
		Info:  make(chan string),
		Debug: make(chan string),
		Trace: make(chan string),
	}
	go func() {
		for {
			select {
			case s := <-ss.Fatal:
				if level >= Nftl {
					Ftl.Chan <- name + ": " + s
				}
			case s := <-ss.Error:
				if level >= Nerr {
					Err.Chan <- name + ": " + s
				}
			case s := <-ss.Warn:
				if level >= Nwrn {
					Wrn.Chan <- name + ": " + s
				}
			case s := <-ss.Info:
				if level >= Ninf {
					Inf.Chan <- name + ": " + s
				}
			case s := <-ss.Debug:
				if level >= Ndbg {
					Dbg.Chan <- name + ": " + s
				}
			case s := <-ss.Trace:
				if level >= Ntrc {
					Trc.Chan <- name + ": " + s
				}
			case <-Quit:
				break
			}
		}
	}()
	return &ss
}

// Check checks if an error exists, if so, prints a log to the specified log level with a string and returns if err was nil
func Check(err error, tag int, where string) (wasError bool) {
	if err != nil {
		wasError = true
		L[tag].Chan <- L[tag].Name() + " " + err.Error()
		if tag == Ftl.Num {
			panic("died")
		}
	}
	return
}

// Lvl is a log level data structure
type Lvl struct {
	Num  int
	Name func(c ...int) string
	Chan chan string
}

const (
	// Nftl is the number for fatal errors
	Nftl = iota
	// Nerr is the number for errors
	Nerr
	// Nwrn is the number for warnings
	Nwrn
	// Ninf is the number for information
	Ninf
	// Ndbg is the number for debugging
	Ndbg
	// Ntrc is the number for trace
	Ntrc
)

var (
	ftlFn = func(c ...int) string {
		out := "FTL"
		if len(c) > 0 {
			return aurora.BgRed(out).String()
		}
		return out
	}
	errFn = func(c ...int) string {
		out := "ERR"
		if len(c) > 0 {
			return aurora.Red(out).String()
		}
		return out
	}
	wrnFn = func(c ...int) string {
		out := "ERR"
		if len(c) > 0 {
			return aurora.Brown(out).String()
		}
		return out
	}
	infFn = func(c ...int) string {
		out := "INF"
		if len(c) > 0 {
			return aurora.Green(out).String()
		}
		return out
	}
	dbgFn = func(c ...int) string {
		out := "DBG"
		if len(c) > 0 {
			return aurora.Blue(out).String()
		}
		return out
	}
	trcFn = func(c ...int) string {
		out := "TRC"
		if len(c) > 0 {
			return aurora.BgBlue(out).String()
		}
		return out
	}
	// Ftl is for critical/fatal errors
	Ftl = &Lvl{0, ftlFn, nil}
	// Err is an error that does block continuation
	Err = &Lvl{1, errFn, nil}
	// Wrn is is a warning of a correctable condition
	Wrn = &Lvl{2, wrnFn, nil}
	// Inf is is general information
	Inf = &Lvl{3, infFn, nil}
	// Dbg is debug level information
	Dbg = &Lvl{4, dbgFn, nil}
	// Trc is detailed outputs of contents of variables
	Trc = &Lvl{5, trcFn, nil}

	// L is an array of log levels that can be selected given the level number
	L = []*Lvl{
		Ftl,
		Err,
		Wrn,
		Inf,
		Dbg,
		Trc,
	}

	// Quit signals the logger to stop
	Quit = make(chan struct{})

	// OutFile sets the output file for the logger
	OutFile *os.File

	// LogIt is the function that performs the output, can be loaded by the caller
	LogIt = Print

	color = true
)

// Disabled is a no-op print function
func Disabled(name, txt string) {
}

// SetPrinter loads a different print function
func SetPrinter(fn func(name, txt string)) {
	LogIt = fn
}

// Color sets whether tags are coloured or not, 0 color
func Color(on bool) {
	color = on
}

// GetColor returns if color is turned on
func GetColor() bool {
	return color
}

func init() {
	Init()
}

// Init manually starts a clog
func Init(fn ...func(name, txt string)) bool {
	var ready []chan bool
	Ftl.Chan = make(chan string)
	Err.Chan = make(chan string)
	Wrn.Chan = make(chan string)
	Inf.Chan = make(chan string)
	Dbg.Chan = make(chan string)
	Trc.Chan = make(chan string)
	// override the output function if one is given
	if len(fn) > 0 {
		LogIt = fn[0]
	}
	for range L {
		ready = append(ready, make(chan bool))
	}
	for i := range L {
		go startChan(i, ready[i])
	}
	for i := range ready {
		<-ready[i]
	}
	Dbg.Chan <- "logger started"
	return true
}

// Print out a formatted log message
func Print(name, txt string) {
	out := fmt.Sprintf("%s [%s] %s\n",
		time.Now().UTC().Format("2006-01-02 15:04:05.000000 MST"),
		name,
		txt,
	)
	if OutFile != nil {
		fmt.Fprint(OutFile, out)
	}
	fmt.Print(out)
	if Joined != nil {
		Joined <- out
	}
}

func startChan(ch int, ready chan bool) {
	L[ch].Chan = make(chan string)
	ready <- true
	done := true
	for done {
		select {
		case <-Quit:
			done = false
			continue
		case txt := <-L[ch].Chan:
			if color {
				LogIt(L[ch].Name(1), txt)
			} else {
				LogIt(L[ch].Name(), txt)
			}
			continue
		default:
		}
	}
}
