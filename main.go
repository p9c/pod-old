package main

import (
	"fmt"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/debug"

	"git.parallelcoin.io/pod/pkg/limits"
	"git.parallelcoin.io/pod/app"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	debug.SetGCPercent(1)
	if err := limits.SetLimits(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to set limits: %v\n", err)
		os.Exit(1)
	}

	os.Exit(pod.Main())
}
