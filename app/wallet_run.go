package app

import (
	"encoding/json"
	"fmt"

	walletmain "git.parallelcoin.io/pod/cmd/wallet"

	cl "git.parallelcoin.io/pod/pkg/clog"
)

func runWallet(wc *walletmain.Config) int {
	j, _ := json.MarshalIndent(wc, "", "  ")
	log <- cl.Tracef{"running with configuration:\n%s", string(j)}
	err := walletmain.Main(wc)
	if err != nil {
		fmt.Print(err)
		return 1
	}
	return 0
}
