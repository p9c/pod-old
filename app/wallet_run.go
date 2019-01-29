package app

import (
	"encoding/json"
	"fmt"

	walletmain "git.parallelcoin.io/pod/cmd/wallet"

	cl "git.parallelcoin.io/pod/pkg/clog"
)

func runWallet() int {
	j, _ := json.MarshalIndent(WalletConfig, "", "  ")
	log <- cl.Tracef{"running with configuration:\n%s", string(j)}
	go func() {
		for {
			select {
			case <-walletmain.WalletDone:
				break
			}
		}
	}()
	err := walletmain.Main(WalletConfig.Wallet)
	if err != nil {
		fmt.Print(err)
		return 1
	}
	return 0
}
