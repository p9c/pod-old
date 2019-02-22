package walletmain

import (
	cl "git.parallelcoin.io/pod/pkg/clog"
)

// Log is the logger for node
var Log = cl.NewSubSystem("cmd/wallet     ", "info")
var log = Log.Ch

// UseLogger uses a specified Logger to output package logging info. This should be used in preference to SetLogWriter if the caller is also using log.
func UseLogger(
	logger *cl.SubSystem) {
	Log = logger
	log = Log.Ch
}
