package addrmgr

import (
	"git.parallelcoin.io/pod/pkg/clog"
)

// Log is the logger for the addrmgr package
var Log = cl.NewSubSystem("lib/addrmgr", "info")
var log = Log.Ch
