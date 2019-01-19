package node

import (
	"git.parallelcoin.io/pod/lib/clog"
)

// Log is the logger for node
var Log = clog.NewSubSystem("pod/node", clog.Ninf)

// directionString is a helper function that returns a string that represents the direction of a connection (inbound or outbound).
func directionString(inbound bool) string {
	if inbound {
		return "inbound"
	}
	return "outbound"
}

// pickNoun returns the singular or plural form of a noun depending on the count n.
func pickNoun(n uint64, singular, plural string) string {
	if n == 1 {
		return singular
	}
	return plural
}
