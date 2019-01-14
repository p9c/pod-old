package main

import (
	"fmt"
	"os"
)

type walletSpvCfg struct{}

var walletspv walletSpvCfg

func (n *walletSpvCfg) Execute(args []string) (err error) {
	fmt.Println("running wallet with spv node")
	fmt.Println("not implemented - quitting")
	os.Exit(1)
	return
}
