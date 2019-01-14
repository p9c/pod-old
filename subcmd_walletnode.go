package main

import (
	"encoding/json"
	"fmt"
	"os"
)

var walletnode walletnodeCfg

func (n *walletnodeCfg) Execute(args []string) (err error) {
	fmt.Println("running wallet with full node")
	j, _ := json.MarshalIndent(n, "", "\t")
	fmt.Println(string(j))
	fmt.Println("not implemented - quitting")
	os.Exit(1)
	return
}
