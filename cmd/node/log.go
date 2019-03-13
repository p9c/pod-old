package node

import (
	cl "git.parallelcoin.io/clog"
)

// Log is the logger for node
var Log = cl.NewSubSystem("cmd/node", "info")

var log = Log.Ch

// UseLogger uses a specified Logger to output package logging info. This should be used in preference to SetLogWriter if the caller is also using log.
func UseLogger(
	logger *cl.SubSystem,
) {

	Log = logger
	log = Log.Ch
}

// directionString is a helper function that returns a string that represents the direction of a connection (inbound or outbound).
func directionString(
	inbound bool,
) string {

	if inbound {
		return "inbound"
	}
	return "outbound"
}

// pickNoun returns the singular or plural form of a noun depending on the count n.
func pickNoun(
	n uint64,
	singular,
	plural string,
) string {

	if n == 1 {
		return singular
	}
	return plural
}
