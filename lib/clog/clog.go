package clog

// Miscellaneous functions

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/logrusorgru/aurora"
)

// Joined is a channel that can be used to redirect the entire log output to another channel
var Joined chan string

// Chan is a simple string channel
type Chan chan string

// FmtChan is a chan for Fmt
type FmtChan chan Fmt

// A SubSystem is a logger that intercepts a signal, adds a 'name' prefix and passes it to the main logger channel
type SubSystem struct {
	Fatal  Chan
	Fatalf FmtChan
	Error  Chan
	Errorf FmtChan
	Warn   Chan
	Warnf  FmtChan
	Info   Chan
	Infof  FmtChan
	Debug  Chan
	Debugf FmtChan
	Trace  Chan
	Tracef FmtChan
	Level  int
	Quit   chan struct{}
}

// SetLevel sets the debug level according to a string
func (s *SubSystem) SetLevel(level string) {
	switch strings.ToLower(level) {
	case "fatal":
		s.Level = Nftl
	case "error":
		s.Level = Nerr
	case "warn":
		s.Level = Nwrn
	case "info":
		s.Level = Ninf
	case "debug":
		s.Level = Ndbg
	case "trace":
		s.Level = Ntrc
	default:
		s.Level = Ninf
	}
}

// NewSubSystem creates a new clog logger that adds a prefix to the log entry for subsystem control
func NewSubSystem(name string, level int) *SubSystem {
	ss := SubSystem{
		Fatal:  make(Chan),
		Fatalf: make(FmtChan),
		Error:  make(Chan),
		Errorf: make(FmtChan),
		Warn:   make(Chan),
		Warnf:  make(FmtChan),
		Info:   make(Chan),
		Infof:  make(FmtChan),
		Debug:  make(Chan),
		Debugf: make(FmtChan),
		Trace:  make(Chan),
		Tracef: make(FmtChan),
		Level:  level,
		Quit:   make(chan struct{}),
	}
	go func() {
		for {
			select {
			case s := <-ss.Fatal:
				if ss.Level >= Nftl {
					Ftl.Chan <- name + ": " + s
				}
			case s := <-ss.Error:
				if ss.Level >= Nerr {
					Err.Chan <- name + ": " + s
				}
			case s := <-ss.Warn:
				if ss.Level >= Nwrn {
					Wrn.Chan <- name + ": " + s
				}
			case s := <-ss.Info:
				if ss.Level >= Ninf {
					Inf.Chan <- name + ": " + s
				}
			case s := <-ss.Debug:
				if ss.Level >= Ndbg {
					Dbg.Chan <- name + ": " + s
				}
			case s := <-ss.Trace:
				if ss.Level >= Ntrc {
					Trc.Chan <- name + ": " + s
				}
			case s := <-ss.Fatalf:
				if ss.Level >= Nftl {
					Ftl.Chan <- name + ": " + fmt.Sprintf(s.Fmt, s.Items...)
				}
			case s := <-ss.Errorf:
				if ss.Level >= Nerr {
					Err.Chan <- name + ": " + fmt.Sprintf(s.Fmt, s.Items...)
				}
			case s := <-ss.Warnf:
				if ss.Level >= Nwrn {
					Wrn.Chan <- name + ": " + fmt.Sprintf(s.Fmt, s.Items...)
				}
			case s := <-ss.Infof:
				if ss.Level >= Ninf {
					Inf.Chan <- name + ": " + fmt.Sprintf(s.Fmt, s.Items...)
				}
			case s := <-ss.Debugf:
				if ss.Level >= Ndbg {
					Dbg.Chan <- name + ": " + fmt.Sprintf(s.Fmt, s.Items...)
				}
			case s := <-ss.Tracef:
				if ss.Level >= Ntrc {
					Trc.Chan <- name + ": " + fmt.Sprintf(s.Fmt, s.Items...)
				}
			case <-ss.Quit:
				break
			case <-Quit:
				break
			}
		}
	}()
	return &ss
}

// Print is a function that shortens the invocation for pushing to a logging channel so you can call like this:
//     log.Info.Print("string")
// though you can always just directly do what is done here in this function, where log.Info is the channel. This is just here for completeness with the FmtChan Print, which makes a less weildy invocation by using the builtin variant arguments processing
// This format also allows you to use comma separated list of strings, instead of using concatenation operators, it also interposes spaces between the strings
func (c *Chan) Print(s ...string) {
	tmp := ""
	for i := range s {
		if i > 0 {
			tmp += " "
		}
		tmp += s[i]
	}
	*c <- tmp
}

// Print is a shortcut to assemble an Fmt struct literal. It should be inlined by the compiler
func (s *FmtChan) Print(fmt string, items ...interface{}) {
	*s <- Fmt{fmt, items}
}

// T is a slice of interface{} for feeding to a *Printf function
type T []interface{}

// Fmt is a printf type formatter struct, it is used like this:
//     logf.Fmt("format string %s %d", "test", 100)
// When all parts are strings it is faster to use a SubSystem. Many types provide stringers.
type Fmt struct {
	Fmt   string
	Items []interface{}
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

	// Quit signals the logger to stop. Invoke like this:
	//     close(clog.Quit)
	// You can call Init again to start it up again
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

func emptyJoiner(string) {
	// This just swallows the input and discards it
}

// Joiner is the function that will pipe the output
var Joiner = emptyJoiner

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
	go func() {
		for {
			select {
			case s := <-Joined:
				Joiner(s)
			default:
			}
		}
	}()
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

// Shutdown the application, allowing the logger a moment to clear the channels
func Shutdown() {
	// wait a moment to let log channel clear
	time.Sleep(time.Millisecond * 50)
	os.Exit(0)
}
