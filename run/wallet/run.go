package walletrun

import (
	"encoding/json"

	cl "git.parallelcoin.io/pod/lib/clog"
)

func runNode() {
	j, _ := json.MarshalIndent(Config, "", "  ")
	log <- cl.Trc("running with configuration:\n" + string(j))
}
