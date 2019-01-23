package walletrun

import (
	"encoding/json"

	"git.parallelcoin.io/pod/lib/clog"
	"git.parallelcoin.io/pod/module/wallet"
)

func runNode() {
	j, _ := json.MarshalIndent(Config, "", "  ")
	log <- cl.Trc("running with configuration:\n" + string(j))
	walletmain.Main(Config.Wallet)
}
