package node

import (
	"encoding/json"
	"fmt"

	"git.parallelcoin.io/pod/module/node"
)

func runNode() {
	j, _ := json.MarshalIndent(Config, "", "  ")
	Log.Tracef.Print("running with configuration:\n%s", string(j))
	fmt.Println(Config.Node.DbType)
	node.Main(Config.Node, nil)
}
