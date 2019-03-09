package app

import (
	"path/filepath"

	"git.parallelcoin.io/pod/cmd/node"
	"git.parallelcoin.io/pod/pkg/util"
)

// AppName is the name of this application
var AppName = "pod"

// DefaultDataDir is the default location for the data
var DefaultDataDir = util.AppDataDir(AppName, false)

// DefaultShellConfFileName is
var DefaultShellConfFileName = filepath.Join(
	filepath.Join(node.DefaultHomeDir, "shell"), "conf")

// DefaultShellDataDir is the default data directory for the shell
var DefaultShellDataDir = filepath.Join(
	node.DefaultHomeDir, "shell")

var f = GenFlag
var l = GenLog
var s = GenShort
var t = GenTrig
