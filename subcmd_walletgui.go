package main

import (
	"encoding/json"
	"fmt"
	"os"
)

var walletGUI walletGUICfg

func (n *walletGUICfg) Execute(args []string) (err error) {
	fmt.Println("running wallet gui")
	j, _ := json.MarshalIndent(cfg.General, "", "  ")
	fmt.Println(string(j))
	j, _ = json.MarshalIndent(cfg.Network, "", "  ")
	fmt.Println(string(j))
	j, _ = json.MarshalIndent(n, "", "\t")
	fmt.Println(string(j))
	fmt.Println("not implemented - quitting")
	os.Exit(1)
	return
}
