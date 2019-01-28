package app

import (
	"path/filepath"

	"git.parallelcoin.io/pod/cmd/node"
	"git.parallelcoin.io/pod/pkg/util"
)

var (
	// AppName is the name of this application
	AppName = "pod"
	// DefaultDataDir is the default location for the data
	DefaultDataDir = util.AppDataDir(AppName, false)
	// DefaultShellDataDir is the default data directory for the shell
	DefaultShellDataDir = filepath.Join(
		node.DefaultHomeDir, "shell")
	// DefaultShellConfFileName is
	DefaultShellConfFileName = filepath.Join(
		filepath.Join(node.DefaultHomeDir, "shell"), "conf")
	f = GenFlag
	t = GenTrig
	s = GenShort
	l = GenLog
)
