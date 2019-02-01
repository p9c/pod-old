package app

import (
	"encoding/json"

	cl "git.parallelcoin.io/pod/pkg/clog"
)

func runShell() int {
	j, _ := json.MarshalIndent(ShellConfig, "", "  ")
	log <- cl.Tracef{"running with configuration:\n%s", string(j)}
	shutdown := make(chan struct{})
	go func() {
		go runWallet(ShellConfig.Wallet, ShellConfig.walletActiveNet)
		runNode(ShellConfig.Node, ShellConfig.nodeActiveNet)
		close(shutdown)
	}()
	<-shutdown
	return 0
}
