package app_old

import (
	"encoding/json"
	"fmt"

	"git.parallelcoin.io/dev/pod/cmd/node"
	cl "git.parallelcoin.io/dev/pod/pkg/util/cl"
)

func runNode(
	nc *node.Config,
	activeNet *node.Params,
) int {

	j, _ := json.MarshalIndent(nc, "", "  ")
	log <- cl.Tracef{"running with configuration:\n%s", string(j)}
	err := node.Main(nc, activeNet, nil)
	if err != nil {
		fmt.Print(err)
		return 1
	}
	return 0
}
