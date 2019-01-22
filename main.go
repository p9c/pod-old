package main

import (
	"fmt"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/debug"

	"git.parallelcoin.io/pod/lib/limits"
	"git.parallelcoin.io/pod/run"
)

func main() {
	// f, err := os.Create("trace.out")
	// if err != nil {
	// 	panic(err)
	// }
	// defer f.Close()
	// err = trace.Start(f)
	// if err != nil {
	// 	panic(err)
	// }
	// defer trace.Stop()

	runtime.GOMAXPROCS(runtime.NumCPU())
	debug.SetGCPercent(1)
	if err := limits.SetLimits(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to set limits: %v\n", err)
		os.Exit(1)
	}

	os.Exit(pod.Main())

	// cf, err := os.Create("cpu.prof")
	// if err != nil {
	// 	log.Fatal("could not create CPU profile: ", err)
	// }
	// if err := pprof.StartCPUProfile(cf); err != nil {
	// 	log.Fatal("could not start CPU profile: ", err)
	// }
	// defer pprof.StopCPUProfile()

	// go func() {
	// 	log.Println(http.ListenAndServe("localhost:6060", nil))
	// }()

	// if err := pod.Main(); err != nil {
	// }
	// mf, err := os.Create("mem.prof")
	// if err != nil {
	// 	log.Fatal("could not create memory profile: ", err)
	// }
	// runtime.GC() // get up-to-date statistics
	// if err := pprof.WriteHeapProfile(mf); err != nil {
	// 	log.Fatal("could not write memory profile: ", err)
	// }
	// f.Close()
}
