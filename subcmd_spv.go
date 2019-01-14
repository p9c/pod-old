package main

import (
	"fmt"
	"os"
)

type spvCfg struct{}

var spv spvCfg

func (n *spvCfg) Execute(args []string) (err error) {
	fmt.Println("running spv node")
	fmt.Println("not implemented - quitting")
	os.Exit(1)
	return
}
