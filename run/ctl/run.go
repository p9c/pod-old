package ctl

import (
	"encoding/json"

	"git.parallelcoin.io/pod/lib/clog"
	"git.parallelcoin.io/pod/module/ctl"
)

func runCtl(args []string) {
	j, _ := json.MarshalIndent(Config, "", "  ")
	log <- cl.Tracef{"running with configuration:\n%s", string(j)}
	ctl.Main(args, Config)
}
