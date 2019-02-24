package controller

import (
	cl "git.parallelcoin.io/pod/pkg/util/clog"
)

// Log is the logger for the peer package
var Log = cl.NewSubSystem("mining/controller", "info")
var log = Log.Ch

// UseLogger uses a specified Logger to output package logging info. This should be used in preference to SetLogWriter if the caller is also using log.
func UseLogger(
	logger *cl.SubSystem) {

	Log = logger
	log = Log.Ch
}
