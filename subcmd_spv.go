package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type spvCfg struct{}

var spv spvCfg

func (n *spvCfg) Execute(args []string) (err error) {
	fmt.Println("running spv node")
	j, _ := json.MarshalIndent(n, "", "\t")
	fmt.Println(string(j))
	fmt.Println("not implemented - quitting")
	os.Exit(1)
	return
}
