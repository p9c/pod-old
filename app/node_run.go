package app

import (
	"encoding/json"

	cl "git.parallelcoin.io/pod/pkg/clog"
	"git.parallelcoin.io/pod/cmd/node"
)

func runNode() {
	j, _ := json.MarshalIndent(NodeConfig, "", "  ")
	log <- cl.Tracef{"running with configuration:\n%s", string(j)}
	node.Main(NodeConfig.Node, nil)
}
