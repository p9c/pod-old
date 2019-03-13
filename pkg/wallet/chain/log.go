package chain

import (
	cl "git.parallelcoin.io/pod/pkg/util/cl"
)

// Log is the logger for the peer package
var Log = cl.NewSubSystem("wallet/chain", "info")
var log = Log.Ch

/*
// LogClosure is a closure that can be printed with %v to be used to generate expensive-to-create data for a detailed log level and avoid doing the work if the data isn't printed.
type logClosure func() string

// String invokes the log closure and returns the results string.
func (
	c logClosure,
	) String() string {
		return c()
	}


	// UseLogger uses a specified Logger to output package logging info. This should be used in preference to SetLogWriter if the caller is also using log.
	func UseLogger(
		logger *cl.SubSystem) {

			Log = logger
			log = Log.Ch
		}

		// newLogClosure returns a new closure over the passed function which allows
		// it to be used as a parameter in a logging function that is only invoked when
		// the logging level is such that the message will actually be logged.
		func newLogClosure(
			c func() string) logClosure {
				return logClosure(c)
			}


			// pickNoun returns the singular or plural form of a noun depending
			// on the count n.
			func pickNoun(
				n int, singular, plural string) string {
					if n == 1 {
						return singular
					}
					return plural
				}

*/
