package main

import (
	"fmt"
	"os"
)

type walletSpvGUICfg struct{}

var walletspvgui walletSpvGUICfg

func (n *walletSpvGUICfg) Execute(args []string) (err error) {
	fmt.Println("running wallet with spv and gui")
	fmt.Println("not implemented - quitting")
	os.Exit(1)
	return
}
