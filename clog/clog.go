package clog

// Miscellaneous functions

import (
	"fmt"
	"time"

	"github.com/logrusorgru/aurora"
)

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

	// LogLevel is a dynamically settable log level filter that excludes higher values from output
	LogLevel = Trc.Num

	// Quit signals the logger to stop
	Quit = make(chan struct{})

	// LogIt is the function that performs the output, can be loaded by the caller
	LogIt = Print

	color = true
)

// Color sets whether tags are coloured or not, 0 color
func Color(on bool) {
	color = on
}

// GetColor returns if color is turned on
func GetColor() bool {
	return color
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
	fmt.Printf("%s [%s] %s\n",
		time.Now().UTC().Format("2006-01-02 15:04:05.000000 MST"),
		name,
		txt,
	)
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
			if ch <= LogLevel {
				if color {
					LogIt(L[ch].Name(1), txt)
				} else {
					LogIt(L[ch].Name(), txt)
				}
				if ch == Nftl {
					panic(txt)
				}
			}
			continue
		default:
		}
	}
}
