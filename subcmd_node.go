package main

import (
	"encoding/json"
	"fmt"
	"os"
)

var node nodeCfg

func (n *nodeCfg) Execute(args []string) (err error) {
	fmt.Println("running full node")
	j, _ := json.MarshalIndent(n, "", "\t")
	fmt.Println(string(j))
	fmt.Println("not implemented - quitting")
	os.Exit(1)
	return
}
