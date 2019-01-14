package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type walletSpvCfg struct{}

var walletspv walletSpvCfg

func (n *walletSpvCfg) Execute(args []string) (err error) {
	fmt.Println("running wallet with spv node")
	j, _ := json.MarshalIndent(n, "", "\t")
	fmt.Println(string(j))
	fmt.Println("not implemented - quitting")
	os.Exit(1)
	return
}
