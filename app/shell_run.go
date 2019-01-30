package app

import (
	"encoding/json"
	"time"

	cl "git.parallelcoin.io/pod/pkg/clog"
)

func runShell() int {
	j, _ := json.MarshalIndent(ShellConfig, "", "  ")
	log <- cl.Infof{"running with configuration:\n%s", string(j)}
	time.Sleep(time.Second)
	go runNode(ShellConfig.Node)
	time.Sleep(time.Second * 3)
	go runWallet(ShellConfig.Wallet)
	return 0
}
