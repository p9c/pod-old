package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"

	"git.parallelcoin.io/pod/lib/limits"
	"git.parallelcoin.io/pod/run"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	debug.SetGCPercent(10)
	if err := limits.SetLimits(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to set limits: %v\n", err)
		os.Exit(1)
	}
	if err := pod.Main(); err != nil {
		os.Exit(1)
	}
}
