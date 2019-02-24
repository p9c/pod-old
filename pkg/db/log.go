package database

import (
	cl "git.parallelcoin.io/pod/pkg/util/clog"
)

// Log is the logger for the peer package
var Log = cl.NewSubSystem("pkg/db         ", "info")
var log = Log.Ch

// UseLogger uses a specified Logger to output package logging info.
func UseLogger(
	logger *cl.SubSystem) {

	Log = logger
	log = Log.Ch
}
