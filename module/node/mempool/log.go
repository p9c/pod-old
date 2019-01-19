package mempool

import (
	"git.parallelcoin.io/pod/lib/clog"
)

// Log is the logger for the peer package
var Log = clog.NewSubSystem("pod/node/mempool", clog.Ndbg)

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

// pickNoun returns the singular or plural form of a noun depending on the count n.
func pickNoun(n int, singular, plural string) string {
	if n == 1 {
		return singular
	}
	return plural
}
