package netsync

import (
	"git.parallelcoin.io/pod/lib/clog"
)

// Log is the logger for the netsync package
var Log = clog.NewSubSystem("netsync", clog.Ninf)

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

// // UseLogger uses a specified Logger to output package logging info.
// // This should be used in preference to SetLogWriter if the caller is also
// // using log.
// func UseLogger(logger l.Logger) {
// 	log = logger
// }