package app_old

import (
	"encoding/json"
	"fmt"

	walletmain "git.parallelcoin.io/dev/pod/cmd/wallet"
	netparams "git.parallelcoin.io/dev/pod/pkg/chain/config/params"

	cl "git.parallelcoin.io/dev/pod/pkg/util/cl"
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
