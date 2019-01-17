package pod

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"git.parallelcoin.io/pod/module/ctl"
)

var ctlcfg ctlCfg

func (n *ctlCfg) Execute(args []string) (err error) {
	fmt.Println("running ctl")
	joined := ctl.Config{
		ShowVersion:   cfg.General.ShowVersion,
		ListCommands:  cfg.Ctl.CtlLaunch.ListCommands,
		ConfigFile:    ctl.DefaultConfigFile,
		RPCServer:     ctl.DefaultRPCServer,
		RPCCert:       ctl.DefaultRPCCertFile,
		RPCUser:       defaultUser,
		RPCPassword:   defaultPass,
		TLS:           n.CtlRPC.TLS,
		Proxy:         n.CtlRPC.Proxy,
		ProxyUser:     n.CtlRPC.ProxyUser,
		ProxyPass:     n.CtlRPC.ProxyUser,
		TestNet3:      cfg.Network.TestNet3,
		SimNet:        cfg.Network.SimNet,
		TLSSkipVerify: n.CtlRPC.TLSSkipVerify,
		Wallet:        n.CtlLaunch.Wallet,
	}
	switch {
	case n.CtlRPC.RPCUser != "":
		joined.RPCUser = n.CtlRPC.RPCUser
	case n.CtlRPC.RPCPassword != "":
		joined.RPCPassword = n.CtlRPC.RPCPassword
	case !n.CtlLaunch.ListCommands:
		joined.ListCommands = n.CtlLaunch.ListCommands
	case cfg.General.ConfigFile != "":
		joined.ConfigFile = cfg.General.ConfigFile
	case n.CtlRPC.RPCServer != "":
		joined.RPCServer = n.CtlRPC.RPCServer
	case n.CtlRPC.RPCCert != "":
		joined.RPCCert = n.CtlRPC.RPCCert
	}
	if cfg.General.SaveConfig {
		j, _ := json.MarshalIndent(joined, "", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Println(joined.ConfigFile)
		ensureDir(joined.ConfigFile)
		err := ioutil.WriteFile(joined.ConfigFile, j, 0600)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(j))
	}
	return
}
