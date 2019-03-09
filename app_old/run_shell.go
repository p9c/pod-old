package app

import (
	"encoding/json"
	"sync"
	"time"

	cl "git.parallelcoin.io/clog"
	"git.parallelcoin.io/pod/pkg/util/interrupt"
)

func runShell() (
	out int,
) {

	j, _ := json.MarshalIndent(ShellConfig, "", "  ")
	log <- cl.Tracef{"running with configuration:\n%s", string(j)}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {

		out = runNode(ShellConfig.Node, ShellConfig.GetNodeActiveNet())
		wg.Done()
	}()
	time.Sleep(time.Second * 2)
	wg.Add(1)
	go func() {

		out = runWallet(ShellConfig.Wallet, ShellConfig.GetWalletActiveNet())
		wg.Done()
	}()
	wg.Wait()
	<-interrupt.HandlersDone
	return 0
}
