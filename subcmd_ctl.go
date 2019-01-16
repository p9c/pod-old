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
		ShowVersion:   cfg.General.ShowVersion,
		ConfigFile:    ctl.DefaultConfigFile,
		RPCServer:     ctl.DefaultRPCServer,
		RPCCert:       ctl.DefaultRPCCertFile,
		RPCUser:       n.CtlRPC.RPCUser,
		RPCPassword:   n.CtlRPC.RPCPassword,
		TLS:           n.CtlRPC.TLS,
		Proxy:         n.CtlRPC.Proxy,
		ProxyUser:     n.CtlRPC.ProxyUser,
		ProxyPass:     n.CtlRPC.ProxyUser,
		TestNet3:      cfg.Network.TestNet3,
		SimNet:        cfg.Network.SimNet,
		TLSSkipVerify: n.CtlRPC.TLSSkipVerify,
		Wallet:        n.CtlLaunch.Wallet,
	}
	joined.ListCommands = n.CtlLaunch.ListCommands
	if cfg.General.ConfigFile != "" {
		joined.ConfigFile = cfg.General.ConfigFile
	}
	if n.CtlRPC.RPCServer != "" {
		joined.RPCServer = n.CtlRPC.RPCServer
	}
	if n.CtlRPC.RPCCert != "" {
		joined.RPCCert = n.CtlRPC.RPCCert
	}
	j, _ := json.MarshalIndent(joined, "", "  ")
	fmt.Println(string(j))
	fmt.Println(args)
	os.Exit(1)
	return
}
