package main

import (
	"encoding/json"
	"fmt"
	"os"

	"git.parallelcoin.io/pod/ctl"
)

var ctlcfg ctlCfg

func (n *ctlCfg) Execute(args []string) (err error) {
	fmt.Println("running ctl")
	joined := ctl.Config{
		ConfigFile: ctl.DefaultConfigFile,
		RPCServer:  ctl.DefaultRPCServer,
		RPCCert:    ctl.DefaultRPCCertFile,
	}
	joined.ShowVersion = cfg.General.ShowVersion
	joined.ListCommands = n.CtlLaunch.ListCommands
	if cfg.General.ConfigFile != "" {
		joined.ConfigFile = cfg.General.ConfigFile
	}
	joined.RPCUser = n.CtlRPC.RPCUser
	joined.RPCPassword = n.CtlRPC.RPCPassword
	if n.CtlRPC.RPCServer != "" {
		joined.RPCServer = n.CtlRPC.RPCServer
	}
	if n.CtlRPC.RPCCert != "" {
		joined.RPCCert = n.CtlRPC.RPCCert
	}
	joined.TLS = n.CtlRPC.TLS
	joined.Proxy = n.CtlRPC.Proxy
	joined.ProxyUser = n.CtlRPC.ProxyUser
	joined.ProxyPass = n.CtlRPC.ProxyUser
	joined.TestNet3 = cfg.Network.TestNet3
	joined.SimNet = cfg.Network.SimNet
	joined.TLSSkipVerify = n.CtlRPC.TLSSkipVerify
	joined.Wallet = n.CtlLaunch.Wallet
	j, _ := json.MarshalIndent(joined, "", "  ")
	fmt.Println(string(j))
	fmt.Println(args)
	os.Exit(1)
	return
}
