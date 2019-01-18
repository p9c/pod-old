package pod

import (
	"git.parallelcoin.io/pod/module/ctl"
	"git.parallelcoin.io/pod/module/node"
	shell "git.parallelcoin.io/pod/module/wallet"
)

func init() {
	ensureDir(ctl.DefaultConfigFile)
	ensureDir(node.DefaultConfigFile)
	ensureDir(shell.DefaultConfigFile)
}
