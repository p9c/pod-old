package app

import (
	"git.parallelcoin.io/pod/pkg/clog"
)

var Log = cl.NewSubSystem("pod", "info")
var log = Log.Ch
