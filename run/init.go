package pod

import (
	"git.parallelcoin.io/pod/module/ctl"
	"git.parallelcoin.io/pod/module/node"
	"git.parallelcoin.io/pod/module/wallet"
	"git.parallelcoin.io/pod/run/shell"
	"git.parallelcoin.io/pod/run/util"
)

func init() {
	pu.EnsureDir(ctl.DefaultConfigFile)
	pu.EnsureDir(node.DefaultConfigFile)
	pu.EnsureDir(walletmain.DefaultConfigFile)
	pu.EnsureDir(shell.DefaultConfFileName)
}
