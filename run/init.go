package pod

import (
	"git.parallelcoin.io/pod/module/ctl"
	"git.parallelcoin.io/pod/module/node"
	"git.parallelcoin.io/pod/module/wallet"
	"git.parallelcoin.io/pod/run/shell"
)

func init() {
	ensureDir(ctl.DefaultConfigFile)
	ensureDir(node.DefaultConfigFile)
	ensureDir(walletmain.DefaultConfigFile)
	ensureDir(shell.DefaultConfFileName)
}
