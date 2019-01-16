package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"

	"git.parallelcoin.io/pod/lib/limits"
	"git.parallelcoin.io/pod/run"
)

const (
	showHelpMessage = "Specify -h to show available options"
)

// usage displays the general usage when the help flag is not displayed and and an invalid command was specified.  The commandUsage function is used instead when a valid command was specified.
func usage(errorMessage string) {
	appName := filepath.Base(os.Args[0])
	appName = strings.TrimSuffix(appName, filepath.Ext(appName))
	fmt.Fprintln(os.Stderr, errorMessage)
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintf(os.Stderr, "  %s [OPTIONS] <command> <args...>\n\n",
		appName)
	fmt.Fprintln(os.Stderr, showHelpMessage)
}

// winServiceMain is only invoked on Windows.  It detects when pod is running as a service and reacts accordingly.
var winServiceMain func() (bool, error)

// Main is the real pod main
func Main(args []string) (err error) {
	pod.LoadConfig()
	// interrupt := interruptListener()
	// defer fmt.Println("Shutdown complete")
	// if interruptRequested(interrupt) {
	// 	return nil
	// }
	// <-interrupt
	return
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	debug.SetGCPercent(10)
	if err := limits.SetLimits(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to set limits: %v\n", err)
		os.Exit(1)
	}
	if runtime.GOOS == "windows" {
		isService, err := winServiceMain()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if isService {
			os.Exit(0)
		}
	}
	if err := Main(os.Args); err != nil {
		os.Exit(1)
	}
}
