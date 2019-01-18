package pod

import (
	"git.parallelcoin.io/pod/module/ctl"
	"git.parallelcoin.io/pod/module/node"
	"git.parallelcoin.io/pod/module/wallet"
	"git.parallelcoin.io/pod/run/shell"
	"git.parallelcoin.io/pod/run/util"
)

func init() {
	podutil.EnsureDir(ctl.DefaultConfigFile)
	podutil.EnsureDir(node.DefaultConfigFile)
	podutil.EnsureDir(walletmain.DefaultConfigFile)
	podutil.EnsureDir(shell.DefaultConfFileName)
}
