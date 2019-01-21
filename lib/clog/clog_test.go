package cl

import (
	"io"
	"os"
	"testing"
	"time"
)

func TestClog(t *testing.T) {
	logfile, err := os.OpenFile("/tmp/clog", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		t.Fatal(err.Error())
	}
	Writer = io.MultiWriter(os.Stdout, logfile)
	defer logfile.Close()

	ColorChan <- false

	done := make(chan bool)
	go tests(&Og, done)
	<-done
	close(done)

	ss := *NewSubSystem("testsubsystem", "trace")

	done = make(chan bool)
	go tests(&ss.Ch, done)
	<-done
	close(done)

	ColorChan <- true

	done = make(chan bool)
	go tests(&Og, done)
	<-done
	close(done)

	done = make(chan bool)
	go tests(&ss.Ch, done)
	<-done
	close(done)

	testString := "testing closure"
	var testClosure StringClosure = func() string {
		time.Sleep(time.Millisecond * 100)
		return testString
	}
	ss.Ch <- ss.Ftlc(testClosure)
	ss.Ch <- ss.Errc(testClosure)
	ss.Ch <- ss.Wrnc(testClosure)
	ss.Ch <- ss.Infc(testClosure)
	ss.Ch <- ss.Dbgc(testClosure)
	ss.Ch <- ss.Trcc(testClosure)

	time.Sleep(time.Second)

	ss.Close()

	close(Og)
}

func tests(ch *chan interface{}, done chan bool) {

	*ch <- Ftl("fatal")
	*ch <- Fatal{1, "test", Og}
	*ch <- Fatalf{"%d %s %v", 3, "test3", []byte{1, 2, 3}}
	*ch <- Err("err")
	*ch <- Error{1, "test", Og}
	*ch <- Errorf{"%d %s %v", 3, "test3", []byte{1, 2, 3}}
	*ch <- Wrn("warn")
	*ch <- Warn{1, "test", Og}
	*ch <- Warnf{"%d %s %v", 3, "test3", []byte{1, 2, 3}}
	*ch <- Inf("info")
	*ch <- Info{1, "test", Og}
	*ch <- Infof{"%d %s %v", 3, "test3", []byte{1, 2, 3}}
	*ch <- Dbg("debug")
	*ch <- Debug{1, "test", Og}
	*ch <- Debugf{"%d %s %v", 3, "test3", []byte{1, 2, 3}}
	*ch <- Trc("trace")
	*ch <- Trace{1, "test", Og}
	*ch <- Tracef{"%d %s %v", 3, "test3", []byte{1, 2, 3}}
	done <- true
}
