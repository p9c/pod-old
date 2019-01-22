package pod

import (
	"git.parallelcoin.io/pod/run/ctl"
	"git.parallelcoin.io/pod/run/node"
	"git.parallelcoin.io/pod/run/shell"
	"git.parallelcoin.io/pod/run/wallet"
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
func Main() int {
	PodApp.AddCommand(ctl.Command)
	PodApp.AddCommand(node.Command)
	PodApp.AddCommand(walletrun.Command)
	PodApp.AddCommand(shell.Command)
	return PodApp.Run()
}
