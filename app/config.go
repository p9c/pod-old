package app

import(
 "git.parallelcoin.io/pod/pkg/util"
 "git.parallelcoin.io/pod/cmd/node"
 "path/filepath"
)
var (
	AppName = "pod"
 	DefaultDataDir = util.AppDataDir(AppName, false)
	DefaultAppDataDir   = filepath.Join(
		node.DefaultHomeDir, "shell")
	DefaultConfFileName = filepath.Join(
		filepath.Join(node.DefaultHomeDir, "shell"), "conf"	)
	f = GenFlag
	t = GenTrig
	s = GenShort
	l = GenLog
)
