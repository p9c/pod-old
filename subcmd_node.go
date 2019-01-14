package main

import (
	"fmt"
	"os"

	"git.parallelcoin.io/pod/node"
)

var nodecfg nodeCfg

func (n *nodeCfg) Execute(args []string) (err error) {
	node.PreMain()
	// j, _ := json.MarshalIndent(n, "", "\t")
	// fmt.Println(string(j))
	fmt.Println("not implemented - quitting")
	os.Exit(1)
	return
}
