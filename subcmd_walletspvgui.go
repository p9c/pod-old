package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type walletSpvGUICfg struct{}

var walletspvgui walletSpvGUICfg

func (n *walletSpvGUICfg) Execute(args []string) (err error) {
	fmt.Println("running wallet with spv and gui")
	j, _ := json.MarshalIndent(n, "", "\t")
	fmt.Println(string(j))
	fmt.Println("not implemented - quitting")
	os.Exit(1)
	return
}
