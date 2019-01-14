package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type walletnodeCfg struct {
	LaunchGroup walletCfgLaunchGroup `group:"launch options"`

	NodeRPCGroup   nodeCfgRPCGroup    `group:"RPC options"`
	NodeP2PGroup   nodeCfgP2PGroup    `group:"P2P options"`
	NodeChainGroup nodeCfgChainGroup  `group:"Chain options"`
	MiningGroup    nodeCfgMiningGroup `group:"Mining options"`

	WalletNodeCfgGroup walletNodeCfg     `group:"node connection options"`
	WalletRPCCfgGroup  walletRPCCfgGroup `group:"wallet RPC configuration"`
}

var walletnode walletnodeCfg

func (n *walletnodeCfg) Execute(args []string) (err error) {
	fmt.Println("running wallet with full node")
	j, _ := json.MarshalIndent(n, "", "\t")
	fmt.Println(string(j))
	fmt.Println("not implemented - quitting")
	os.Exit(1)
	return
}
