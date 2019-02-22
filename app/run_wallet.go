package app

import (
	"encoding/json"
	"fmt"

	walletmain "git.parallelcoin.io/pod/cmd/wallet"

	cl "git.parallelcoin.io/pod/pkg/clog"
	"git.parallelcoin.io/pod/pkg/netparams"
)

func runWallet(
	wc *walletmain.Config,
	activeNet *netparams.Params,
) int {

	j, _ := json.MarshalIndent(wc, "", "  ")
	log <- cl.Tracef{"running with configuration:\n%s", string(j)}
	err := walletmain.Main(wc, activeNet)
	if err != nil {
		fmt.Print(err)
		return 1
	}
	return 0
}
