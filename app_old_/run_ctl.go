package app_old

import (
	"encoding/json"

	cl "git.parallelcoin.io/clog"
	"git.parallelcoin.io/pod/cmd/ctl"
)

func runCtl(
	args []string,
	cc *ctl.Config,
) {

	j, _ := json.MarshalIndent(cc, "", "  ")
	log <- cl.Tracef{"running with configuration:\n%s", string(j)}
	ctl.Main(args, cc)
}
