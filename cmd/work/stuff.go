package main

import "github.com/parallelcointeam/sub/clog"

var (
	lf = &clog.Ftl.Chan
	le = &clog.Err.Chan
	lw = &clog.Wrn.Chan
	li = &clog.Inf.Chan
	ld = &clog.Dbg.Chan
	lt = &clog.Trc.Chan
)
