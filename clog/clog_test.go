package clog

import (
	"testing"
	"time"
)

var (
	lFtl = &Ftl.Chan
	lErr = &Err.Chan
	lWrn = &Wrn.Chan
	lInf = &Inf.Chan
	lDbg = &Dbg.Chan
	lTrc = &Trc.Chan
)

func TestClog(t *testing.T) {
	Color(true)
	done := make(chan bool)
	go tests(done)
	<-done
	close(done)
	done = make(chan bool)
	ss := NewSubSystem("TEST", Trc.Num)
	go testSubSystem(ss, done)
	<-done
	close(Quit)
}

func tests(done chan bool) {
	for i := range L {
		txt := "testing log"
		L[i].Chan <- txt
	}
	time.Sleep(time.Millisecond)
	done <- true
}

func testSubSystem(ss *SubSystem, done chan bool) {
	ss.Fatal <- "testing"
	ss.Error <- "testing"
	ss.Warn <- "testing"
	ss.Info <- "testing"
	ss.Debug <- "testing"
	ss.Trace <- "testing"
	done <- true
}
