package addrmgr

import (
	"git.parallelcoin.io/pod/lib/clog"
)

// Log is the logger for the addrmgr package
var Log = cl.NewSubSystem("lib/addrmgr", "trace")
var log = Log.Ch

// UseLogger uses a specified Logger to output package logging info. This should be used in preference to SetLogWriter if the caller is also using log.
func UseLogger(logger *cl.SubSystem) {
	Log = logger
	log = Log.Ch
}
