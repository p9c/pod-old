package addrmgr

import (
	cl "git.parallelcoin.io/clog"
)

// Log is the logger for the addrmgr package
var Log = cl.NewSubSystem("peer/addrmgr", "info")
var log = Log.Ch
