package app

import (
	"encoding/json"
	"time"

	"git.parallelcoin.io/pod/cmd/node"
	"git.parallelcoin.io/pod/pkg/interrupt"
	"git.parallelcoin.io/pod/pkg/netparams"

	cl "git.parallelcoin.io/pod/pkg/clog"
)

func runShell(
	nodeActiveNet *node.Params,
	walletActiveNet *netparams.Params,
) int {
	j, _ := json.MarshalIndent(ShellConfig, "", "  ")
	log <- cl.Tracef{"running with configuration:\n%s", string(j)}
	go runNode(ShellConfig.Node, nodeActiveNet)
	time.Sleep(time.Second * 3)
	go runWallet(ShellConfig.Wallet, walletActiveNet)
	shutdown := make(chan struct{})
	interrupt.AddHandler(func() {
		close(shutdown)
	})
	<-shutdown
	return 0
}
