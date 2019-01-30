package app

import (
	"encoding/json"
	"fmt"

	"git.parallelcoin.io/pod/cmd/node"
	cl "git.parallelcoin.io/pod/pkg/clog"
)

func runNode(nc *node.Config) int {
	j, _ := json.MarshalIndent(nc, "", "  ")
	log <- cl.Tracef{"running with configuration:\n%s", string(j)}
	go func() {
		for {
			select {
			case <-node.NodeDone:
				break
			}
		}
	}()
	err := node.Main(nc, nil)
	if err != nil {
		fmt.Print(err)
		return 1
	}
	return 0
}
