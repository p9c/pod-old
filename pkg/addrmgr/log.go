package addrmgr

import (
	cl "git.parallelcoin.io/pod/pkg/clog"
)

// Log is the logger for the addrmgr package
var Log = cl.NewSubSystem("pkg/addrmgr ", "info")
var log = Log.Ch
