package main

import (
	"encoding/json"
	"fmt"

	"git.parallelcoin.io/pod/node"
	"git.parallelcoin.io/pod/node/mempool"
	"git.parallelcoin.io/pod/walletmain"
)

type walletnodeCfgJoined struct {
	WalletCfg walletmain.Config
	NodeCfg   node.Config
}

var walletnode walletnodeCfg

func (n *walletnodeCfg) Execute(args []string) (err error) {
	fmt.Println("running wallet with full node")
	joined := walletnodeCfgJoined{
		WalletCfg: walletmain.Config{
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
			RPCConnect:               node.DefaultRPCListener,
			CAFile:                   n.WalletNode.CAFile,
			EnableClientTLS:          false,
			PodUsername:              autoUser,
			PodPassword:              autoPass,
			Proxy:                    "", // n.WalletNode.Proxy,
			ProxyUser:                "", // n.WalletNode.ProxyUser,
			ProxyPass:                "", // n.WalletNode.ProxyPass,
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
			Username:                 "user",
			Password:                 "pa55word",
			ExperimentalRPCListeners: n.WalletRPC.ExperimentalRPCListeners,
			DataDir:                  walletmain.DefaultDataDir,
		},
		NodeCfg: node.Config{
			ShowVersion:          cfg.General.ShowVersion,
			ConfigFile:           node.DefaultConfigFile,
			DataDir:              node.DefaultDataDir,
			LogDir:               node.DefaultLogDir,
			AddPeers:             n.NodeP2P.AddPeers,
			ConnectPeers:         n.NodeP2P.ConnectPeers,
			DisableListen:        n.NodeP2P.DisableListen,
			Listeners:            []string{node.DefaultListener},
			MaxPeers:             node.DefaultMaxPeers,
			DisableBanning:       n.NodeP2P.DisableBanning,
			BanDuration:          node.DefaultBanDuration,
			BanThreshold:         node.DefaultBanThreshold,
			Whitelists:           n.NodeP2P.Whitelists,
			RPCUser:              autoUser,
			RPCPass:              autoPass,
			RPCLimitUser:         "",
			RPCLimitPass:         "",
			RPCListeners:         []string{node.DefaultRPCListener},
			RPCCert:              node.DefaultRPCCertFile,
			RPCKey:               n.NodeRPC.RPCKey,
			RPCMaxClients:        node.DefaultMaxRPCClients,
			RPCMaxWebsockets:     node.DefaultMaxRPCWebsockets,
			RPCMaxConcurrentReqs: node.DefaultMaxRPCConcurrentReqs,
			RPCQuirks:            false, // n.NodeRPC.RPCQuirks,
			DisableRPC:           false, // n.NodeRPC.DisableRPC,
			TLS:                  false,
			DisableDNSSeed:       n.NodeP2P.DisableDNSSeed,
			ExternalIPs:          n.NodeP2P.ExternalIPs,
			Proxy:                n.NodeP2P.Proxy,
			ProxyUser:            n.NodeP2P.ProxyUser,
			ProxyPass:            n.NodeP2P.ProxyPass,
			OnionProxy:           n.NodeP2P.OnionProxy,
			OnionProxyUser:       n.NodeP2P.OnionProxyUser,
			OnionProxyPass:       n.NodeP2P.OnionProxyPass,
			NoOnion:              n.NodeP2P.NoOnion,
			TorIsolation:         n.NodeP2P.TorIsolation,
			TestNet3:             cfg.Network.TestNet3,
			RegressionTest:       cfg.Network.RegressionTest,
			SimNet:               cfg.Network.SimNet,
			AddCheckpoints:       n.NodeChain.AddCheckpoints,
			DisableCheckpoints:   n.NodeChain.DisableCheckpoints,
			DbType:               node.DefaultDbType,
			Profile:              n.NodeLaunch.Profile,
			CPUProfile:           n.NodeLaunch.CPUProfile,
			Upnp:                 n.NodeP2P.Upnp,
			MinRelayTxFee:        n.NodeChain.MinRelayTxFee,
			FreeTxRelayLimit:     node.DefaultFreeTxRelayLimit,
			NoRelayPriority:      n.NodeChain.NoRelayPriority,
			TrickleInterval:      node.DefaultTrickleInterval,
			MaxOrphanTxs:         node.DefaultMaxOrphanTransactions,
			Algo:                 node.DefaultAlgo,
			Generate:             n.NodeMining.Generate,
			GenThreads:           node.DefaultGenThreads,
			MiningAddrs:          n.NodeMining.MiningAddrs,
			MinerListener:        node.DefaultMinerListener,
			MinerPass:            "pa55word",
			BlockMinSize:         node.DefaultBlockMinSize,
			BlockMaxSize:         node.DefaultBlockMaxSize,
			BlockMinWeight:       node.DefaultBlockMinWeight,
			BlockMaxWeight:       node.DefaultBlockMaxWeight,
			BlockPrioritySize:    mempool.DefaultBlockPrioritySize,
			UserAgentComments:    n.NodeP2P.UserAgentComments,
			NoPeerBloomFilters:   n.NodeP2P.NoPeerBloomFilters,
			NoCFilters:           n.NodeP2P.NoCFilters,
			DropCfIndex:          n.NodeLaunch.DropCfIndex,
			SigCacheMaxSize:      node.DefaultSigCacheMaxSize,
			BlocksOnly:           n.NodeP2P.BlocksOnly,
			TxIndex:              n.NodeChain.TxIndex,
			DropTxIndex:          n.NodeLaunch.DropTxIndex,
			AddrIndex:            n.NodeChain.AddrIndex,
			DropAddrIndex:        n.NodeLaunch.DropAddrIndex,
			RelayNonStd:          n.NodeP2P.RelayNonStd,
			RejectNonStd:         n.NodeP2P.RejectNonStd,
			// lookup: ,
			// oniondial: ,
			// dial: ,
			// addCheckpoints: ,
			// miningAddrs: ,
			// minerKey: ,
			// minRelayTxFee: ,
			// whitelists: ,
		},
	}
	switch {
	// wallet items
	case n.WalletNode.RPCConnect != "":
		joined.WalletCfg.RPCConnect = n.WalletNode.RPCConnect
	case n.WalletRPC.Username != "":
		joined.WalletCfg.Username = n.WalletRPC.Username
	case n.WalletRPC.Password != "":
		joined.WalletCfg.Password = n.WalletRPC.Password
	case n.WalletNode.PodUsername != "":
		joined.WalletCfg.PodUsername = n.WalletNode.PodUsername
	case n.WalletNode.PodPassword != "":
		joined.WalletCfg.PodPassword = n.WalletNode.PodPassword
	case cfg.General.ConfigFile != "":
		joined.WalletCfg.ConfigFile = cfg.General.ConfigFile
	case cfg.General.DataDir != "":
		joined.WalletCfg.AppDataDir = cfg.General.DataDir
		joined.WalletCfg.DataDir = cfg.General.DataDir
	case cfg.General.LogDir != "":
		joined.WalletCfg.LogDir = cfg.General.LogDir
	case n.WalletRPC.RPCCert != "":
		joined.WalletCfg.RPCCert = n.WalletRPC.RPCCert
	case n.WalletRPC.RPCKey != "":
		joined.WalletCfg.RPCKey = n.WalletRPC.RPCKey
	case n.WalletRPC.LegacyRPCMaxClients != 0:
		joined.WalletCfg.LegacyRPCMaxClients = n.WalletRPC.LegacyRPCMaxClients
	case n.WalletRPC.LegacyRPCMaxWebsockets != 0:
		joined.WalletCfg.LegacyRPCMaxWebsockets = n.WalletRPC.LegacyRPCMaxWebsockets
		// node items
	case n.NodeRPC.RPCUser != "":
		joined.NodeCfg.RPCUser = n.NodeRPC.RPCUser
	case n.NodeRPC.RPCPass != "":
		joined.NodeCfg.RPCPass = n.NodeRPC.RPCUser
	case n.NodeRPC.RPCLimitUser != "":
		joined.NodeCfg.RPCLimitUser = n.NodeRPC.RPCLimitUser
	case n.NodeRPC.RPCLimitPass != "":
		joined.NodeCfg.RPCLimitPass = n.NodeRPC.RPCLimitPass
	case cfg.General.ConfigFile != "":
		joined.NodeCfg.ConfigFile = cfg.General.ConfigFile
	case cfg.General.DataDir != "":
		joined.NodeCfg.DataDir = cfg.General.DataDir
	case n.NodeP2P.Listeners != nil:
		joined.NodeCfg.Listeners = n.NodeP2P.Listeners
	case n.NodeP2P.MaxPeers != 0:
		joined.NodeCfg.MaxPeers = n.NodeP2P.MaxPeers
	case n.NodeP2P.BanDuration != 0:
		joined.NodeCfg.BanDuration = n.NodeP2P.BanDuration
	case n.NodeP2P.BanThreshold != 0:
		joined.NodeCfg.BanThreshold = n.NodeP2P.BanThreshold
	case n.NodeRPC.RPCListeners != nil:
		joined.NodeCfg.RPCListeners = n.NodeRPC.RPCListeners
	case n.NodeRPC.RPCCert != "":
		joined.NodeCfg.RPCCert = n.NodeRPC.RPCCert
	case n.NodeRPC.RPCMaxClients != 0:
		joined.NodeCfg.RPCMaxClients = int(n.NodeRPC.RPCMaxClients)
	case n.NodeRPC.RPCMaxWebsockets != 0:
		joined.NodeCfg.RPCMaxWebsockets = int(n.NodeRPC.RPCMaxWebsockets)
	case n.NodeRPC.RPCMaxConcurrentReqs != 0:
		joined.NodeCfg.RPCMaxConcurrentReqs = int(n.NodeRPC.RPCMaxConcurrentReqs)
	case n.NodeChain.DbType != "":
		joined.NodeCfg.DbType = n.NodeChain.DbType
	case n.NodeChain.FreeTxRelayLimit != 0:
		joined.NodeCfg.FreeTxRelayLimit = n.NodeChain.FreeTxRelayLimit
	case n.NodeChain.TrickleInterval != 0:
		joined.NodeCfg.TrickleInterval = n.NodeChain.TrickleInterval
	case n.NodeChain.MaxOrphanTxs != 0:
		joined.NodeCfg.MaxOrphanTxs = n.NodeChain.MaxOrphanTxs
	case n.NodeMining.Algo != "":
		joined.NodeCfg.Algo = n.NodeMining.Algo
	case n.NodeMining.GenThreads != 0:
		joined.NodeCfg.GenThreads = n.NodeMining.GenThreads
	case n.NodeMining.MinerListener != "":
		joined.NodeCfg.MinerListener = n.NodeMining.MinerListener
	case n.NodeChain.BlockMinSize != 0:
		joined.NodeCfg.BlockMinSize = n.NodeChain.BlockMinSize
	case n.NodeChain.BlockMaxSize != 0:
		joined.NodeCfg.BlockMaxSize = n.NodeChain.BlockMaxSize
	case n.NodeChain.BlockMinWeight != 0:
		joined.NodeCfg.BlockMinWeight = n.NodeChain.BlockMinWeight
	case n.NodeChain.BlockMaxWeight != 0:
		joined.NodeCfg.BlockMaxWeight = n.NodeChain.BlockMaxWeight
	case n.NodeChain.SigCacheMaxSize != 0:
		joined.NodeCfg.SigCacheMaxSize = n.NodeChain.SigCacheMaxSize
	}
	j, _ := json.MarshalIndent(joined, "", "  ")
	fmt.Println(string(j))
	return
}
