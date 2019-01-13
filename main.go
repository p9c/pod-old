package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/EXCCoin/exccd/limits"
)

const (
	showHelpMessage = "Specify -h to show available options"
	listCmdMessage  = "Specify -l to list available commands"
)

var (
	cfg *config
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
	fmt.Fprintln(os.Stderr, listCmdMessage)
}

// winServiceMain is only invoked on Windows.  It detects when pod is running as a service and reacts accordingly.
var winServiceMain func() (bool, error)

// Main is the real pod main
func Main() (err error) {
	// Load configuration and parse command line.  This function also initializes logging and configures it accordingly.
	tcfg, args, err := loadConfig()
	if err != nil {
		return err
	}
	cfg = tcfg
	// Get a channel that will be closed when a shutdown signal has been triggered either from an OS signal such as SIGINT (Ctrl+C) or from another subsystem such as the RPC server.
	interrupt := interruptListener()
	defer podLog.Info("Shutdown complete")
	// Show version at startup.
	podLog.Infof("Version %s", version())
	if len(args) < 1 {
		usage("No command specified")
		os.Exit(1)
	}
	// Return now if an interrupt signal was triggered.
	if interruptRequested(interrupt) {
		return nil
	}
	// Wait until the interrupt signal is received from an OS signal or shutdown is requested through one of the subsystems such as the RPC server.
	<-interrupt
	return
}

func main() {
	// Use all processor cores.
	runtime.GOMAXPROCS(runtime.NumCPU())
	// Block and transaction processing can cause bursty allocations.  This limits the garbage collector from excessively overallocating during bursts.  This value was arrived at with the help of profiling live usage.
	debug.SetGCPercent(10)
	// Up some limits.
	if err := limits.SetLimits(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to set limits: %v\n", err)
		os.Exit(1)
	}
	// Call serviceMain on Windows to handle running as a service.  When the return isService flag is true, exit now since we ran as a service.  Otherwise, just fall through to normal operation.
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
	// Work around defer not working after os.Exit()
	if err := Main(); err != nil {
		os.Exit(1)
	}
}
