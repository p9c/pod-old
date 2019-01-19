package txscript

import (
	"git.parallelcoin.io/pod/lib/clog"
)

// Log is the logger for the peer package
var Log = clog.NewSubSystem("pod/lib/txscript", clog.Ndbg)

// // log is a logger that is initialized with no output filters.  This means the package will not perform any logging by default until the caller requests it.
// var log = l.Disabled

// // The default amount of logging is none.
// func init() {
// 	// DisableLog()
// }

// // DisableLog disables all library log output.  Logging output is disabled by default until UseLogger is called.
// func DisableLog() {
// 	log = l.Disabled
// }

// UseLogger uses a specified Logger to output package logging info.
func UseLogger(logger *clog.SubSystem) {
	Log = logger
}

// LogClosure is a closure that can be printed with %v to be used to generate expensive-to-create data for a detailed log level and avoid doing the work if the data isn't printed.
type logClosure func() string

func (c logClosure) String() string {
	return c()
}
func newLogClosure(c func() string) logClosure {
	return logClosure(c)
}
