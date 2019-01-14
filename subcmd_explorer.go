package main

import (
	"encoding/json"
	"fmt"
	"os"
)

var explorer explorerCfg

func (n *explorerCfg) Execute(args []string) (err error) {
	fmt.Println("running explorer")
	j, _ := json.MarshalIndent(n, "", "\t")
	fmt.Println(string(j))
	fmt.Println("not implemented - quitting")
	os.Exit(1)
	return
}
