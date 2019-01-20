package ctl

import (
	"encoding/json"

	"git.parallelcoin.io/pod/lib/clog"
)

func runCtl() {
	j, _ := json.MarshalIndent(Config, "", "  ")
	log <- cl.Tracef{"running with configuration:\n%s", string(j)}
}
