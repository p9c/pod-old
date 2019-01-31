package app

import (
	"encoding/json"
	"time"

	walletmain "git.parallelcoin.io/pod/cmd/wallet"
	"git.parallelcoin.io/pod/pkg/interrupt"

	cl "git.parallelcoin.io/pod/pkg/clog"
)

func runShell() int {
	j, _ := json.MarshalIndent(ShellConfig, "", "  ")
	log <- cl.Tracef{"running with configuration:\n%s", string(j)}
	go runNode(ShellConfig.Node)
	time.Sleep(time.Second * 3)
	go runWallet(ShellConfig.Wallet, walletmain.ActiveNet)
	shutdown := make(chan struct{})
	interrupt.AddHandler(func() {
		close(shutdown)
	})
	<-shutdown
	return 0
}
