package main

import (
	"fmt"
	"os"
)

type walletGUICfgLaunchGroup struct {
	ShowVersion   bool   `short:"V" long:"version" description:"display version information and exit"`
	ConfigFile    string `short:"C" long:"configfile" description:"path to configuration file"`
	DataDir       string `short:"b" long:"datadir" description:"directory to store data"`
	LogDir        string `long:"logdir" description:"directory to log output"`
	Profile       string `long:"profile" description:"Enable HTTP profiling on given port -- NOTE port must be between 1024 and 65536"`
	CPUProfile    string `long:"cpuprofile" description:"Write CPU profile to the specified file"`
	DropCfIndex   bool   `long:"dropcfindex" description:"Deletes the index used for committed filtering (CF) support from the database on start up and then exits."`
	DropTxIndex   bool   `long:"droptxindex" description:"Deletes the hash-based transaction index from the database on start up and then exits."`
	DropAddrIndex bool   `long:"dropaddrindex" description:"Deletes the address-based transaction index from the database on start up and then exits."`
}

type walletGUICfg struct {
	LaunchGroup    walletGUICfgLaunchGroup `group:"launch options"`
	NodeP2PGroup   nodeCfgP2PGroup         `group:"P2P options"`
	NodeChainGroup nodeCfgChainGroup       `group:"Chain options"`
}

var walletGUI walletGUICfg

func (n *walletGUICfg) Execute(args []string) (err error) {
	fmt.Println("running wallet gui")
	fmt.Println("not implemented - quitting")
	os.Exit(1)
	return
}
