package app

import (
	"git.parallelcoin.io/pod/cmd/ctl"
	"git.parallelcoin.io/pod/cmd/node"
)

func init() {
	pu.EnsureDir(ctl.DefaultConfigFile)
	pu.EnsureDir(node.DefaultConfigFile)
	pu.EnsureDir(walletmain.DefaultConfigFile)
	pu.EnsureDir(shell.DefaultConfFileName)
}
