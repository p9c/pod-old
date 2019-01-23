package app

import (
	"git.parallelcoin.io/pod/pkg/clog"
)

// Log is the logger for node
var Log = cl.NewSubSystem("module/node", "info")
var log = Log.Ch
