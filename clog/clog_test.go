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
	Init()
	Color(false)
	done := make(chan bool)
	LogLevel = Trc.Num
	go tests(done)
	<-done
}

func tests(done chan bool) {
	for i := range L {
		txt := "testing log"
		L[i].Chan <- txt
	}
	time.Sleep(time.Millisecond)
	done <- true
}
