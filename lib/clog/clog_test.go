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
	ss := NewSubSystem("TEST", Ntrc)
	go testSubSystem(ss, done)
	<-done
	close(done)
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
	ss.Fatal.Print("testing1", "eleventy", "'")
	ss.Error.Print("testing1", "eleventy", "'")
	ss.Warn.Print("testing1", "eleventy", "'")
	ss.Info.Print("testing1", "eleventy", "'")
	ss.Debug.Print("testing1", "eleventy", "'")
	ss.Trace.Print("testing1", "eleventy", "'")
	ss.Fatalf.Print("testing1 string[%s] number[%d]", "eleventy", 110)
	ss.Errorf.Print("testing1 string[%s] number[%d]", "eleventy", 110)
	ss.Warnf.Print("testing1 string[%s] number[%d]", "eleventy", 110)
	ss.Infof.Print("testing1 string[%s] number[%d]", "eleventy", 110)
	ss.Debugf.Print("testing1 string[%s] number[%d]", "eleventy", 110)
	ss.Tracef.Print("testing1 string[%s] number[%d]", "eleventy", 110)
	done <- true
}
