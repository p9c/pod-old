package app

import (
	"encoding/json"

	cl "git.parallelcoin.io/pod/pkg/clog"
<<<<<<< HEAD
	"git.parallelcoin.io/pod/module/node"
=======
	"git.parallelcoin.io/pod/cmd/node"
>>>>>>> master
)

func runNode() {
	j, _ := json.MarshalIndent(Config, "", "  ")
	log <- cl.Tracef{"running with configuration:\n%s", string(j)}
	node.Main(Config.Node, nil)
}
