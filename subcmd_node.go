package main

import (
	"encoding/json"
	"fmt"

	"git.parallelcoin.io/pod/node"
	"git.parallelcoin.io/pod/node/mempool"
)

var nodecfg nodeCfg

func (n *nodeCfg) Execute(args []string) (err error) {
	fmt.Println("running node")
	joined := node.Config{}
	joined.ShowVersion = cfg.General.ShowVersion
	joined.ConfigFile = node.DefaultConfigFile
	if cfg.General.ConfigFile != "" {
		joined.ConfigFile = cfg.General.ConfigFile
	}
	joined.DataDir = node.DefaultDataDir
	if cfg.General.DataDir != "" {
		joined.DataDir = cfg.General.DataDir
	}
	joined.LogDir = node.DefaultLogDir
	if cfg.General.LogDir != "" {
		joined.LogDir = cfg.General.LogDir
	}
	joined.AddPeers = n.NodeP2P.AddPeers
	joined.ConnectPeers = n.NodeP2P.ConnectPeers
	joined.Listeners = []string{node.DefaultListener}
	if cfg.Node.NodeP2P.Listeners != nil {
		joined.Listeners = cfg.Node.NodeP2P.Listeners
	}
	joined.DisableListen = n.NodeP2P.DisableListen
	joined.MaxPeers = n.NodeP2P.MaxPeers
	joined.DisableBanning = n.NodeP2P.DisableBanning
	joined.BanDuration = node.DefaultBanDuration
	if n.NodeP2P.BanDuration != 0 {
		joined.BanDuration = n.NodeP2P.BanDuration
	}
	joined.BanThreshold = node.DefaultBanThreshold
	if n.NodeP2P.BanThreshold != 0 {
		joined.BanThreshold = n.NodeP2P.BanThreshold
	}
	joined.Whitelists = n.NodeP2P.Whitelists
	joined.RPCUser = n.NodeRPC.RPCUser
	joined.RPCPass = n.NodeRPC.RPCPass
	joined.RPCLimitUser = n.NodeRPC.RPCLimitUser
	joined.RPCLimitPass = n.NodeRPC.RPCLimitPass
	joined.RPCListeners = []string{fmt.Sprintf("127.0.0.1:%d", node.ActiveNetParams.RPCPort)}
	if n.NodeRPC.RPCListeners != nil {
		joined.RPCListeners = n.NodeRPC.RPCListeners
	}
	joined.RPCCert = node.DefaultRPCCertFile
	if n.NodeRPC.RPCCert != "" {
		joined.RPCCert = n.NodeRPC.RPCCert
	}
	joined.RPCKey = node.DefaultRPCKeyFile
	if n.NodeRPC.RPCKey != "" {
		joined.RPCKey = n.NodeRPC.RPCKey
	}
	joined.RPCMaxClients = node.DefaultMaxRPCClients
	if n.NodeRPC.RPCMaxClients != 0 {
		joined.RPCMaxClients = int(n.NodeRPC.RPCMaxClients)
	}
	joined.RPCMaxWebsockets = node.DefaultMaxRPCWebsockets
	if n.NodeRPC.RPCMaxWebsockets != 0 {
		joined.RPCMaxWebsockets = int(n.NodeRPC.RPCMaxWebsockets)
	}
	joined.RPCMaxConcurrentReqs = node.DefaultMaxRPCConcurrentReqs
	if n.NodeRPC.RPCMaxConcurrentReqs != 0 {
		joined.RPCMaxConcurrentReqs = int(n.NodeRPC.RPCMaxConcurrentReqs)
	}
	joined.RPCQuirks = n.NodeRPC.RPCQuirks
	joined.DisableRPC = n.NodeRPC.DisableRPC
	joined.TLS = n.NodeRPC.TLS
	joined.DisableDNSSeed = n.NodeP2P.DisableDNSSeed
	joined.ExternalIPs = n.NodeP2P.ExternalIPs
	joined.Proxy = n.NodeP2P.Proxy
	joined.ProxyUser = n.NodeP2P.ProxyUser
	joined.ProxyPass = n.NodeP2P.ProxyPass
	joined.OnionProxy = n.NodeP2P.OnionProxy
	joined.OnionProxyUser = n.NodeP2P.OnionProxy
	joined.OnionProxyPass = n.NodeP2P.OnionProxyPass
	joined.NoOnion = n.NodeP2P.NoOnion
	joined.TorIsolation = n.NodeP2P.TorIsolation
	joined.TestNet3 = cfg.Network.TestNet3
	joined.RegressionTest = cfg.Network.RegressionTest
	joined.SimNet = cfg.Network.SimNet
	joined.AddCheckpoints = n.NodeChain.AddCheckpoints
	joined.DisableCheckpoints = n.NodeChain.DisableCheckpoints
	joined.DbType = node.DefaultDbType
	if n.NodeChain.DbType != "" {
		joined.DbType = n.NodeChain.DbType
	}
	joined.Upnp = n.NodeP2P.Upnp
	joined.MinRelayTxFee = mempool.DefaultMinRelayTxFee.ToDUO()
	if n.NodeChain.MinRelayTxFee != 0 {
		joined.MinRelayTxFee = n.NodeChain.MinRelayTxFee
	}
	joined.FreeTxRelayLimit = node.DefaultFreeTxRelayLimit
	if n.NodeChain.FreeTxRelayLimit != 0 {
		joined.FreeTxRelayLimit = n.NodeChain.FreeTxRelayLimit
	}
	joined.TrickleInterval = node.DefaultTrickleInterval
	if n.NodeChain.TrickleInterval != 0 {
		joined.TrickleInterval = n.NodeChain.TrickleInterval
	}
	joined.MaxOrphanTxs = node.DefaultMaxOrphanTransactions
	if n.NodeChain.MaxOrphanTxs != 0 {
		joined.MaxOrphanTxs = n.NodeChain.MaxOrphanTxs
	}
	joined.Algo = node.DefaultAlgo
	joined.Generate = node.DefaultGenerate
	joined.GenThreads = node.DefaultGenThreads
	//joined.=MiningAddr
	joined.MinerController = n.NodeMining.MinerController
	joined.MinerPort = node.DefaultMinerPort
	//joined.=MinerPas
	joined.BlockMinSize = node.DefaultBlockMinSize
	joined.BlockMaxSize = node.DefaultBlockMaxSize
	joined.BlockMinWeight = node.DefaultBlockMinWeight
	joined.BlockMaxWeight = node.DefaultBlockMaxWeight
	joined.BlockPrioritySize = mempool.DefaultBlockPrioritySize
	joined.UserAgentComments = n.NodeP2P.UserAgentComments
	joined.NoPeerBloomFilters = n.NodeP2P.NoPeerBloomFilters
	joined.NoCFilters = n.NodeP2P.NoCFilters
	joined.DropCfIndex = n.NodeLaunch.DropCfIndex
	joined.SigCacheMaxSize = node.DefaultSigCacheMaxSize
	joined.BlocksOnly = n.NodeP2P.BlocksOnly
	joined.TxIndex = node.DefaultTxIndex
	joined.DropTxIndex = n.NodeLaunch.DropTxIndex
	joined.AddrIndex = node.DefaultAddrIndex
	joined.DropAddrIndex = n.NodeLaunch.DropAddrIndex
	joined.RelayNonStd = n.NodeP2P.RejectNonStd
	joined.RejectNonStd = n.NodeP2P.RejectNonStd
	joined.ShowVersion = cfg.General.ShowVersion

	j, _ := json.MarshalIndent(joined, "", "  ")
	fmt.Println(string(j))
	fmt.Println(args)

	// node.PreMain()
	return
}
