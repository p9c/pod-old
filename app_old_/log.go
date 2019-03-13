package app_old

import (
	cl "git.parallelcoin.io/pod/pkg/util/cl"
)

// Log is the logger for node
var Log = cl.NewSubSystem("pod/app", "info")

var log = Log.Ch
