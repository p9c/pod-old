package netsync

import (
	"git.parallelcoin.io/pod/lib/clog"
)

// Log is the logger for the netsync package
var Log = clog.NewSubSystem("lib/netsync", clog.Ndbg)

// import (
// 	l "git.parallelcoin.io/pod/lib/log"
// )

// // log is a logger that is initialized with no output filters.  This
// // means the package will not perform any logging by default until the caller
// // requests it.
// var log l.Logger

// // DisableLog disables all library log output.  Logging output is disabled
// // by default until either UseLogger or SetLogWriter are called.
// func DisableLog() {
// 	log = l.Disabled
// }

// UseLogger uses a specified Logger to output package logging info.
func UseLogger(logger *clog.SubSystem) {
	Log = logger
}
