package chain

import (
	cl "git.parallelcoin.io/dev/pod/pkg/util/cl"
)

// Log is the logger for the peer package
var Log = cl.NewSubSystem("wallet/chain", "info")
var log = Log.Ch
