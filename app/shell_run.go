package app

import (
	"encoding/json"

	"git.parallelcoin.io/pod/cmd/ctl"
	cl "git.parallelcoin.io/pod/pkg/clog"
)

func runShell(args []string) {
	j, _ := json.MarshalIndent(CtlCfg, "", "  ")
	log <- cl.Tracef{"running with configuration:\n%s", string(j)}
	ctl.Main(args, CtlCfg)
}
