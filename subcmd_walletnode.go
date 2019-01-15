package main

import (
	"encoding/json"
	"fmt"
)

type walletnodeCfgJoined struct {
	General    generalCfg
	Network    networkGroup
	LogBase    logTopLevel
	Logging    logSubSystems
	NodeLaunch nodeLaunchGroup
	NodeRPC    nodeCfgRPCGroup
	NodeP2P    nodeCfgP2PGroup
	NodeChain  nodeCfgChainGroup
	NodeMining nodeCfgMiningGroup
}

var walletnode walletnodeCfg

func (n *walletnodeCfg) Execute(args []string) (err error) {
	fmt.Println("running wallet with full node")
	joined := walletnodeCfgJoined{
		General: cfg.General,
	}
	j, _ := json.MarshalIndent(joined, "", "  ")
	fmt.Println(string(j))
	fmt.Println(args)
	return
}
