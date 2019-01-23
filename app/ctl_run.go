package app

import (
	"encoding/json"

	"git.parallelcoin.io/pod/pkg/clog"
)

func runCtl(args []string) {
	j, _ := json.MarshalIndent(Config, "", "  ")
	log <- cl.Tracef{"running with configuration:\n%s", string(j)}
	ctl.Main(args, Config)
}
