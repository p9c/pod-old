package node

import (
	"encoding/json"
	"fmt"

	cl "git.parallelcoin.io/pod/lib/clog"
	"git.parallelcoin.io/pod/module/node"
)

func runNode() {
	j, _ := json.MarshalIndent(Config, "", "  ")
	log <- cl.Tracef{"running with configuration:\n%s", string(j)}
	fmt.Println(Config.Node.DbType)
	node.Main(Config.Node, nil)
}
