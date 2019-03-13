package cl

import (
	"fmt"
	"time"

	"git.parallelcoin.io/pod/pkg/util/interrupt"

	"github.com/mitchellh/colorstring"
)

// Close a SubSystem logger
func (s *SubSystem) Close() {

	close(s.Ch)
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

// NewSubSystem starts up a new subsystem logger
func NewSubSystem(name, level string) (ss *SubSystem) {

	wg.Add(1)
	ss = new(SubSystem)
	ss.Ch = make(chan interface{})
	ss.Name = name
	ss.SetLevel(level)
	Register.Add(ss)

	// The main subsystem processing loop
	go func() {

		for i := range ss.Ch {
			// fmt.Println("chan:i := <-ss.Ch")
			if ShuttingDown {
				break
			}
			if i == nil {
				fmt.Println("got nil")
				continue
			}
			n := name
			if Color {
				n = colorstring.Color("[bold]" + n + "[reset]")
			} else {
				n += ":"
			}
			switch I := i.(type) {

			case Ftl:
				if ss.Level > _off {
					Og <- Ftl(n+" ") + I
				}
			case Err:
				if ss.Level > _fatal {
					Og <- Err(n+" ") + I
				}
			case Wrn:
				if ss.Level > _error {
					Og <- Wrn(n+" ") + I
				}
			case Inf:
				if ss.Level > _warn {
					Og <- Inf(n+" ") + I
				}
			case Dbg:
				if ss.Level > _info {
					Og <- Dbg(n+" ") + I
				}
			case Trc:
				if ss.Level > _debug {
					Og <- Trc(n+" ") + I
				}
			case Fatalc:
				if ss.Level > _off {
					fn := func() string {
						o := n + " "
						o += i.(Fatalc)()
						return o
					}
					Og <- Fatalc(fn)
				}
			case Errorc:
				if ss.Level > _fatal {
					fn := func() string {
						o := n + " "
						o += i.(Errorc)()
						return o
					}
					Og <- Errorc(fn)
				}
			case Warnc:
				if ss.Level > _error {
					fn := func() string {
						o := n + " "
						o += i.(Warnc)()
						return o
					}
					Og <- Warnc(fn)
				}
			case Infoc:
				if ss.Level > _warn {
					fn := func() string {
						o := n + " "
						o += i.(Infoc)()
						return o
					}
					Og <- Infoc(fn)
				}
			case Debugc:
				if ss.Level > _info {
					fn := func() string {
						o := n + " "
						o += i.(Debugc)()
						return o
					}
					Og <- Debugc(fn)
				}
			case Tracec:
				if ss.Level > _debug {
					fn := func() string {
						o := n + " "
						o += i.(Tracec)()
						return o
					}
					Og <- Tracec(fn)
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
				if ss.Level > _error {
					Og <- append(Warnf{n + " " + i.(Warnf)[0].(string)}, i.(Warnf)[1:]...)
				}
			case Infof:
				if ss.Level > _warn {
					Og <- append(Infof{n + " " + i.(Infof)[0].(string)}, i.(Infof)[1:]...)
				}
			case Debugf:
				if ss.Level > _info {
					Og <- append(Debugf{n + " " + i.(Debugf)[0].(string)}, i.(Debugf)[1:]...)
				}
			case Tracef:
				if ss.Level > _debug {
					Og <- append(Tracef{n + " " + i.(Tracef)[0].(string)}, i.(Tracef)[1:]...)
				}
			}
		}
	}()
	wg.Done()
	return
}

func init() {

	wg.Add(1)
	worker := func() {

		var t, s string
		for {
			// fmt.Println("clog loop")
			select {
			case <-Quit:
				// fmt.Println("chan:<-Quit")
				ShuttingDown = true
				break
			case Color = <-ColorChan:
				// fmt.Println("chan:Color = <-ColorChan")
			case i := <-Og:
				// fmt.Println("chan:i := <-Og")
				if ShuttingDown {
					break
				}
				if i == nil {
					fmt.Println("received nil")
					continue
				}
				color := Color
				if color {
					s = colorstring.Color("[reset]")
				}
				t = time.Now().UTC().Format("06-01-02 15:04:05.000")
				switch ii := i.(type) {

				case Fatalc:
					s += ii() + "\n"
				case Errorc:
					s += ii() + "\n"
				case Warnc:
					s += ii() + "\n"
				case Infoc:
					s += ii() + "\n"
				case Debugc:
					s += ii() + "\n"
				case Tracec:
					s += ii() + "\n"
				case Ftl:
					s += string(ii) + "\n"
				case Err:
					s += string(ii) + "\n"
				case Wrn:
					s += string(ii) + "\n"
				case Inf:
					s += string(ii) + "\n"
				case Dbg:
					s += string(ii) + "\n"
				case Trc:
					s += string(ii) + "\n"
				case Fatal:
					s += fmt.Sprintln(ii...)
				case Error:
					s += fmt.Sprintln(ii...)
				case Warn:
					s += fmt.Sprintln(ii...)
				case Info:
					s += fmt.Sprintln(ii...)
				case Debug:
					s += fmt.Sprintln(ii...)
				case Trace:
					s += fmt.Sprintln(ii...)
				case Fatalf:
					if I, ok := ii[0].(string); ok {
						s += fmt.Sprintf(I, ii[1:]...) + "\n"
					}
				case Errorf:
					if I, ok := ii[0].(string); ok {
						s += fmt.Sprintf(I, ii[1:]...) + "\n"
					}
				case Warnf:
					if I, ok := ii[0].(string); ok {
						s += fmt.Sprintf(I, ii[1:]...) + "\n"
					}
				case Infof:
					if I, ok := ii[0].(string); ok {
						s += fmt.Sprintf(I, ii[1:]...) + "\n"
					}
				case Debugf:
					if I, ok := ii[0].(string); ok {
						s += fmt.Sprintf(I, ii[1:]...) + "\n"
					}
				case Tracef:
					if I, ok := ii[0].(string); ok {
						s += fmt.Sprintf(I, ii[1:]...) + "\n"
					}
				}
				switch i.(type) {

				case Ftl, Fatal, Fatalf, Fatalc:
					s = ftlTag(color) + s
				case Err, Error, Errorf, Errorc:
					s = errTag(color) + s
				case Wrn, Warn, Warnf, Warnc:
					s = wrnTag(color) + s
				case Inf, Info, Infof, Infoc:
					s = infTag(color) + s
				case Dbg, Debug, Debugf, Debugc:
					s = dbgTag(color) + s
				case Trc, Trace, Tracef, Tracec:
					s = trcTag(color) + s
				}
				if color {
					t = colorstring.Color("[light_gray]" + fmt.Sprintf("%-16v", t) + "[dark_gray]")
				}
				fmt.Fprint(Writer, t+s)
			}
		}
	}
	go worker()
	wg.Done()
}

// Shutdown the application, allowing the logger a moment to clear the channels
func Shutdown() {

	close(Quit)
	wg.Wait()
	<-interrupt.HandlersDone
}
