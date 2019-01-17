package pod

import (
	"git.parallelcoin.io/pod/run/ctl"
	"github.com/tucnak/climax"
)

var interrupt <-chan struct{}

// PodApp is the climax main app controller for pod
var PodApp = climax.Application{
	Name:     "pod",
	Brief:    "multi-application launcher for Parallelcoin Pod",
	Version:  version(),
	Commands: []climax.Command{},
	Topics:   []climax.Topic{},
	Groups:   []climax.Group{},
	Default:  nil,
}

// Main is the real pod main
func Main() (err error) {
	PodApp.AddCommand(ctl.Command)
	PodApp.Run()
	// interrupt = interruptListener()
	// defer clog.Shutdown()
	// if interruptRequested(interrupt) {
	// 	return nil
	// }
	// <-interrupt
	return
}
