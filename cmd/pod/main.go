package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"runtime/trace"

	"git.parallelcoin.io/dev/pod/app"
	"git.parallelcoin.io/dev/pod/pkg/util/interrupt"
	"git.parallelcoin.io/dev/pod/pkg/util/limits"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	debug.SetGCPercent(100)

	if err := limits.SetLimits(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to set limits: %v\n", err)
		os.Exit(1)
	}
	f, err := os.Create("trace.out")

	if err != nil {
		panic(err)
	}
	defer f.Close()
	err = trace.Start(f)

	if err != nil {
		panic(err)
	}

	mf, err := os.Create("mem.prof")

	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	runtime.GC() // get up-to-date statistics

	if err := pprof.WriteHeapProfile(mf); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}

	cf, err := os.Create("cpu.prof")

	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}

	if err := pprof.StartCPUProfile(cf); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}

	go func() {

		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	interrupt.AddHandler(

		func() {

			fmt.Println("stopping trace")
			trace.Stop()

			pprof.StopCPUProfile()
			f.Close()
		},
	)

	os.Exit(app.Main())

}
