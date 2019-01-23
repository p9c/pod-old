package app

import "git.parallelcoin.io/pod/pkg/util"

var (
	AppName = "pod"
 	DefaultDataDir = util.AppDataDir(AppName, false)
	DefaultAppDataDir   = filepath.Join(
		node.DefaultHomeDir, "shell")
	DefaultConfFileName = filepath.Join(
		filepath.Join(node.DefaultHomeDir, "shell"), "conf"	)
	f = pu.GenFlag
	t = pu.GenTrig
	s = pu.GenShort
	l = pu.GenLog
)