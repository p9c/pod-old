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
	close(done)
	done = make(chan bool)
	ssf := NewSubSystemf("TEST", Trc.Num)
	go testSubSystemf(ssf, done)
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
	ss.Fatal.Print("testing1", "testing2", "'")
	ss.Error.Print("testing1", "testing2", "'")
	ss.Warn.Print("testing1", "testing2", "'")
	ss.Info.Print("testing1", "testing2", "'")
	ss.Debug.Print("testing1", "testing2", "'")
	ss.Trace.Print("testing1", "testing2", "'")
	done <- true
}

func testSubSystemf(ss *SubSystemf, done chan bool) {
	ss.Fatal <- Fmt{"testing %s %d", T{"test", 100}}
	ss.Error <- Fmt{"testing %s %d", T{"test", 100}}
	ss.Warn <- Fmt{"testing %s %d", T{"test", 100}}
	ss.Info <- Fmt{"testing %s %d", T{"test", 100}}
	ss.Debug <- Fmt{"testing %s %d", T{"test", 100}}
	ss.Trace <- Fmt{"testing %s %d", T{"test", 100}}
	ss.Fatal.Print("%s %d", "print()", 11)
	ss.Error.Print("%s %d", "print()", 11)
	ss.Warn.Print("%s %d", "print()", 11)
	ss.Info.Print("%s %d", "print()", 11)
	ss.Debug.Print("%s %d", "print()", 11)
	ss.Trace.Print("%s %d", "print()", 11)
	done <- true
}
