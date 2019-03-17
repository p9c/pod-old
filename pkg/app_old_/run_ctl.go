package app_old

import (
	"encoding/json"

	"git.parallelcoin.io/dev/pod/cmd/ctl"
	cl "git.parallelcoin.io/dev/pod/pkg/util/cl"
)

func runCtl(
	args []string,
	cc *ctl.Config,

) {

	j, _ := json.MarshalIndent(cc, "", "  ")
	log <- cl.Tracef{"running with configuration:\n%s", string(j)}

	ctl.Main(args, cc)
}
