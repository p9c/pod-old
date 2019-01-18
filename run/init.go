package pod

import (
	c "git.parallelcoin.io/pod/module/ctl"
	n "git.parallelcoin.io/pod/module/node"
)

func init() {
	ensureDir(c.DefaultConfigFile)
	ensureDir(n.DefaultConfigFile)
}
