package app

import (
	"encoding/json"

	walletmain "git.parallelcoin.io/pod/cmd/wallet"

	cl "git.parallelcoin.io/pod/pkg/clog"
)

func runWallet(args []string) {
	j, _ := json.MarshalIndent(WalletConfig.Wallet, "", "  ")
	log <- cl.Tracef{"running with configuration:\n%s", string(j)}
	walletmain.Main(WalletConfig.Wallet)
}
