package app

import (
	cl "git.parallelcoin.io/pod/pkg/clog"
)

// Log is the logger for node
var Log = cl.NewSubSystem("pod", "info")
var log = Log.Ch
