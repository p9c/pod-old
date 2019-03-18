package cl

import (
	"fmt"
	"time"

	"git.parallelcoin.io/dev/pod/pkg/util/interrupt"

	"github.com/mitchellh/colorstring"
)

// Close a SubSystem logger
func (s *SubSystem) Close() {

	close(s.Ch)
}

// SetLevel changes the level of a subsystem by level name
func (s *SubSystem) SetLevel(level string) {

	if i, ok := Levels[level]; ok {
		s.mutex.Lock()
		defer s.mutex.Unlock()
		s.Level = i
		s.LevelString = level

	} else {
		s.mutex.Lock()
		defer s.mutex.Unlock()
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

	if len(name) > maxLen {

		maxLen = len(name)
	}

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
			n := fmt.Sprintf("%-"+fmt.Sprint(maxLen)+"v", name)

			if Color {

				n = colorstring.Color("[bold]" + n + "[reset]")

			} else {

				n += ":"
			}

			ss.mutex.Lock()
			sslevel := ss.Level
			ss.mutex.Unlock()

			switch I := i.(type) {

			case Ftl:

				if sslevel > _off {

					Og <- Ftl(n+" ") + I
				}
			case Err:

				if sslevel > _fatal {

					Og <- Err(n+" ") + I
				}
			case Wrn:

				if sslevel > _error {

					Og <- Wrn(n+" ") + I
				}
			case Inf:

				if sslevel > _warn {

					Og <- Inf(n+" ") + I
				}
			case Dbg:

				if sslevel > _info {

					Og <- Dbg(n+" ") + I
				}
			case Trc:

				if sslevel > _debug {

					Og <- Trc(n+" ") + I
				}
			case Fatalc:

				if sslevel > _off {

					fn := func() string {

						o := n + " "
						o += I()
						return o
					}
					Og <- Fatalc(fn)
				}
			case Errorc:

				if sslevel > _fatal {

					fn := func() string {

						o := n + " "
						o += I()
						return o
					}
					Og <- Errorc(fn)
				}
			case Warnc:

				if sslevel > _error {

					fn := func() string {

						o := n + " "
						o += I()
						return o
					}
					Og <- Warnc(fn)
				}
			case Infoc:

				if sslevel > _warn {

					fn := func() string {

						o := n + " "
						o += I()
						return o
					}
					Og <- Infoc(fn)
				}
			case Debugc:

				if sslevel > _info {

					fn := func() string {

						o := n + " "
						o += I()
						return o
					}
					Og <- Debugc(fn)
				}
			case Tracec:

				if sslevel > _debug {

					fn := func() string {

						o := n + " "
						o += I()
						return o
					}
					Og <- Tracec(fn)
				}
			case Fatal:

				if sslevel > _off {

					Og <- append(Fatal{n}, i.(Fatal)...)
				}
			case Error:

				if sslevel > _fatal {

					Og <- append(Error{n}, i.(Error)...)
				}
			case Warn:

				if sslevel > _error {

					Og <- append(Warn{n}, i.(Warn)...)
				}
			case Info:

				if sslevel > _warn {

					Og <- append(Info{n}, i.(Info)...)
				}
			case Debug:

				if sslevel > _info {

					Og <- append(Debug{n}, i.(Debug)...)
				}
			case Trace:

				if sslevel > _debug {

					Og <- append(Trace{n}, i.(Trace)...)
				}
			case Fatalf:

				if sslevel > _off {

					Og <- append(Fatalf{n + " " + i.(Fatalf)[0].(string)}, i.(Fatalf)[1:]...)
				}
			case Errorf:

				if sslevel > _fatal {

					Og <- append(Errorf{n + " " + i.(Errorf)[0].(string)}, i.(Errorf)[1:]...)
				}
			case Warnf:

				if sslevel > _error {

					Og <- append(Warnf{n + " " + i.(Warnf)[0].(string)}, i.(Warnf)[1:]...)
				}
			case Infof:

				if sslevel > _warn {

					Og <- append(Infof{n + " " + i.(Infof)[0].(string)}, i.(Infof)[1:]...)
				}
			case Debugf:

				if sslevel > _info {

					Og <- append(Debugf{n + " " + i.(Debugf)[0].(string)}, i.(Debugf)[1:]...)
				}
			case Tracef:

				if sslevel > _debug {

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

					t = colorstring.Color("[light_gray]" + t + "[dark_gray]")
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
