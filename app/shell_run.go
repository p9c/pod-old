package app

import (
	"encoding/json"
	"time"

	"git.parallelcoin.io/pod/cmd/node"
	walletmain "git.parallelcoin.io/pod/cmd/wallet"
	cl "git.parallelcoin.io/pod/pkg/clog"
)

func runShell() int {
	j, _ := json.MarshalIndent(ShellConfig, "", "  ")
	log <- cl.Tracef{"running with configuration:\n%s", string(j)}
	time.Sleep(time.Second)
	go runWallet()
	go runNode()
	var walletdone, nodedone bool
	for {
		select {
		case <-node.NodeDone:
			nodedone = true
			if walletdone == true {
				break
			}
		case <-walletmain.WalletDone:
			walletdone = true
			if nodedone == true {
				break
			}
		}
	}
}
