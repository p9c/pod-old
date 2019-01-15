package main

import (
	"git.parallelcoin.io/pod/node"
)

var nodecfg nodeCfg

func (n *nodeCfg) Execute(args []string) (err error) {
	node.PreMain()
	return
}
