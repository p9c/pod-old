package main

import (
	"encoding/json"
	"fmt"

	"git.parallelcoin.io/pod/walletmain"
)

var wallet walletCfg

func (n *walletCfg) Execute(args []string) (err error) {
	fmt.Println("running wallet")
	joined := walletmain.Config{
		ConfigFile:               walletmain.DefaultConfigFile,
		ShowVersion:              cfg.General.ShowVersion,
		Create:                   n.WalletLaunch.Create,
		CreateTemp:               n.WalletLaunch.CreateTemp,
		AppDataDir:               walletmain.DefaultAppDataDir,
		TestNet3:                 cfg.Network.TestNet3,
		SimNet:                   cfg.Network.SimNet,
		NoInitialLoad:            n.WalletLaunch.NoInitialLoad,
		LogDir:                   walletmain.DefaultLogDir,
		Profile:                  n.WalletLaunch.Profile,
		GUI:                      walletmain.DefaultGUI,
		WalletPass:               "password",
		RPCConnect:               n.WalletNode.RPCConnect,
		CAFile:                   n.WalletNode.CAFile,
		EnableClientTLS:          n.WalletNode.EnableClientTLS,
		PodUsername:              n.WalletNode.PodUsername,
		PodPassword:              n.WalletNode.PodPassword,
		Proxy:                    n.WalletNode.Proxy,
		ProxyUser:                n.WalletNode.ProxyUser,
		ProxyPass:                n.WalletNode.ProxyPass,
		AddPeers:                 n.WalletNode.AddPeers,
		ConnectPeers:             n.WalletNode.ConnectPeers,
		MaxPeers:                 n.WalletNode.MaxPeers,
		BanDuration:              n.WalletNode.BanDuration,
		BanThreshold:             n.WalletNode.BanThreshold,
		RPCCert:                  walletmain.DefaultRPCCertFile,
		RPCKey:                   walletmain.DefaultRPCKeyFile,
		OneTimeTLSKey:            n.WalletRPC.OneTimeTLSKey,
		EnableServerTLS:          n.WalletRPC.EnableServerTLS,
		LegacyRPCListeners:       n.WalletRPC.LegacyRPCListeners,
		LegacyRPCMaxClients:      walletmain.DefaultRPCMaxClients,
		LegacyRPCMaxWebsockets:   walletmain.DefaultRPCMaxWebsockets,
		Username:                 n.WalletRPC.Username,
		Password:                 n.WalletRPC.Password,
		ExperimentalRPCListeners: n.WalletRPC.ExperimentalRPCListeners,
		DataDir:                  walletmain.DefaultDataDir,
	}
	switch {
	case cfg.General.ConfigFile != "":
		joined.ConfigFile = cfg.General.ConfigFile
	case cfg.General.DataDir != "":
		joined.AppDataDir = cfg.General.DataDir
		joined.DataDir = cfg.General.DataDir
	case cfg.General.LogDir != "":
		joined.LogDir = cfg.General.LogDir
	case n.WalletRPC.RPCCert != "":
		joined.RPCCert = n.WalletRPC.RPCCert
	case n.WalletRPC.RPCKey != "":
		joined.RPCKey = n.WalletRPC.RPCKey
	case n.WalletRPC.LegacyRPCMaxClients != 0:
		joined.LegacyRPCMaxClients = n.WalletRPC.LegacyRPCMaxClients
	case n.WalletRPC.LegacyRPCMaxWebsockets != 0:
		joined.LegacyRPCMaxWebsockets = n.WalletRPC.LegacyRPCMaxWebsockets
	}
	j, _ := json.MarshalIndent(joined, "", "  ")
	fmt.Println(string(j))
	return
}
