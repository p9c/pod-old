package pod

import (
	"encoding/json"
	"fmt"

	"git.parallelcoin.io/pod/module/node"
	"git.parallelcoin.io/pod/module/node/mempool"
)

var nodecfg nodeCfg

func (n *nodeCfg) Execute(args []string) (err error) {
	fmt.Println("running node")
	// Load all defaults and from parser the remainder
	joined := node.Config{
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
		RPCUser:              "user",
		RPCPass:              "pa55word",
		RPCLimitUser:         "",
		RPCLimitPass:         "",
		RPCListeners:         []string{node.DefaultRPCListener},
		RPCCert:              node.DefaultRPCCertFile,
		RPCKey:               n.NodeRPC.RPCKey,
		RPCMaxClients:        node.DefaultMaxRPCClients,
		RPCMaxWebsockets:     node.DefaultMaxRPCWebsockets,
		RPCMaxConcurrentReqs: node.DefaultMaxRPCConcurrentReqs,
		RPCQuirks:            n.NodeRPC.RPCQuirks,
		DisableRPC:           n.NodeRPC.DisableRPC,
		TLS:                  n.NodeRPC.TLS,
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
	}

	switch {
	case n.NodeRPC.RPCUser != "":
		joined.RPCUser = n.NodeRPC.RPCUser
	case n.NodeRPC.RPCPass != "":
		joined.RPCPass = n.NodeRPC.RPCUser
	case n.NodeRPC.RPCLimitUser != "":
		joined.RPCLimitUser = n.NodeRPC.RPCLimitUser
	case n.NodeRPC.RPCLimitPass != "":
		joined.RPCLimitPass = n.NodeRPC.RPCLimitPass
	case cfg.General.ConfigFile != "":
		joined.ConfigFile = cfg.General.ConfigFile
	case cfg.General.DataDir != "":
		joined.DataDir = cfg.General.DataDir
	case n.NodeP2P.Listeners != nil:
		joined.Listeners = n.NodeP2P.Listeners
	case n.NodeP2P.MaxPeers != 0:
		joined.MaxPeers = n.NodeP2P.MaxPeers
	case n.NodeP2P.BanDuration != 0:
		joined.BanDuration = n.NodeP2P.BanDuration
	case n.NodeP2P.BanThreshold != 0:
		joined.BanThreshold = n.NodeP2P.BanThreshold
	case n.NodeRPC.RPCListeners != nil:
		joined.RPCListeners = n.NodeRPC.RPCListeners
	case n.NodeRPC.RPCCert != "":
		joined.RPCCert = n.NodeRPC.RPCCert
	case n.NodeRPC.RPCMaxClients != 0:
		joined.RPCMaxClients = int(n.NodeRPC.RPCMaxClients)
	case n.NodeRPC.RPCMaxWebsockets != 0:
		joined.RPCMaxWebsockets = int(n.NodeRPC.RPCMaxWebsockets)
	case n.NodeRPC.RPCMaxConcurrentReqs != 0:
		joined.RPCMaxConcurrentReqs = int(n.NodeRPC.RPCMaxConcurrentReqs)
	case n.NodeChain.DbType != "":
		joined.DbType = n.NodeChain.DbType
	case n.NodeChain.FreeTxRelayLimit != 0:
		joined.FreeTxRelayLimit = n.NodeChain.FreeTxRelayLimit
	case n.NodeChain.TrickleInterval != 0:
		joined.TrickleInterval = n.NodeChain.TrickleInterval
	case n.NodeChain.MaxOrphanTxs != 0:
		joined.MaxOrphanTxs = n.NodeChain.MaxOrphanTxs
	case n.NodeMining.Algo != "":
		joined.Algo = n.NodeMining.Algo
	case n.NodeMining.GenThreads != 0:
		joined.GenThreads = n.NodeMining.GenThreads
	case n.NodeMining.MinerListener != "":
		joined.MinerListener = n.NodeMining.MinerListener
	case n.NodeChain.BlockMinSize != 0:
		joined.BlockMinSize = n.NodeChain.BlockMinSize
	case n.NodeChain.BlockMaxSize != 0:
		joined.BlockMaxSize = n.NodeChain.BlockMaxSize
	case n.NodeChain.BlockMinWeight != 0:
		joined.BlockMinWeight = n.NodeChain.BlockMinWeight
	case n.NodeChain.BlockMaxWeight != 0:
		joined.BlockMaxWeight = n.NodeChain.BlockMaxWeight
	case n.NodeChain.SigCacheMaxSize != 0:
		joined.SigCacheMaxSize = n.NodeChain.SigCacheMaxSize
	}
	j, _ := json.MarshalIndent(joined, "", "  ")
	fmt.Println(string(j))
	fmt.Println(args)

	// node.PreMain()
	return
}
