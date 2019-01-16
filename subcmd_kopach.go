package main

import (
	"encoding/json"
	"fmt"

	"git.parallelcoin.io/pod/ctl"
	"git.parallelcoin.io/pod/node"
)

var minercfg minerCfg

// KopachCfg is the configuration structure for the miner
type KopachCfg struct {
	General generalCfg
	Network networkGroup
	Miner   minerCfg
}

func (n *minerCfg) Execute(args []string) (err error) {
	fmt.Println("running miner")
	joined := KopachCfg{
		General: generalCfg{
			ShowVersion: cfg.General.ShowVersion,
			ConfigFile:  ctl.DefaultConfigFile,
		},
		Network: networkGroup{
			TestNet3: cfg.Network.TestNet3,
			SimNet:   cfg.Network.SimNet,
		},
		Miner: minerCfg{
			Algo:       "random",
			Controller: node.DefaultMinerListener,
			Password:   "pa55word",
		},
	}
	switch {
	case n.Algo != "":
		joined.Miner.Algo = n.Algo
	case n.Controller != "":
		joined.Miner.Controller = n.Controller
	case n.Password != "":
		joined.Miner.Password = n.Password
	}
	j, _ := json.MarshalIndent(joined, "", "  ")
	fmt.Println(string(j))
	fmt.Println(args)
	return
}
