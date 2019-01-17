package pod

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"git.parallelcoin.io/pod/module/node"
	"git.parallelcoin.io/pod/module/node/mempool"
)

var nodecfg nodeCfg

var joined = node.Config{}

func (n *nodeCfg) Execute(args []string) (err error) {
	fmt.Println("running node")
	if _, err := os.Stat(cfg.General.ConfigFile); os.IsNotExist(err) {
		cfg.General.SaveConfig = true
		// Load all defaults and from parser the remainder
		joined = defaultNodeCfg()
	} else {
		cfgfile, err := ioutil.ReadFile(cfg.General.ConfigFile)
		if err != nil {
			fmt.Println(err)
		}
		err = json.Unmarshal(cfgfile, joined)
		if err != nil {
			fmt.Println(err)
		}
	}
	switch {
	case cfg.General.ShowVersion:
		joined.ShowVersion = cfg.General.ShowVersion          
	case n.NodeP2P.AddPeers != nil:
		joined.AddPeers = n.NodeP2P.AddPeers             
	case n.NodeP2P.ConnectPeers != nil:
		joined.ConnectPeers = n.NodeP2P.ConnectPeers         
	case n.NodeP2P.DisableListen:
		joined.DisableListen = n.NodeP2P.DisableListen        
	case n.NodeP2P.DisableBanning:
		joined.DisableBanning = n.NodeP2P.DisableBanning       
	case n.NodeP2P.Whitelists != :
		joined.Whitelists = n.NodeP2P.Whitelists           
	case n.NodeRPC.RPCKey != :
		joined.RPCKey = n.NodeRPC.RPCKey               
	case n.NodeRPC.RPCQuirks != :
		joined.RPCQuirks = n.NodeRPC.RPCQuirks            
	case n.NodeRPC.DisableRPC != :
		joined.DisableRPC = n.NodeRPC.DisableRPC           
	case n.NodeRPC.TLS != :
		joined.TLS = n.NodeRPC.TLS                  
	case n.NodeP2P.DisableDNSSeed != :
		joined.DisableDNSSeed = n.NodeP2P.DisableDNSSeed       
	case n.NodeP2P.ExternalIPs != :
		joined.ExternalIPs = n.NodeP2P.ExternalIPs          
	case n.NodeP2P.Proxy != :
		joined.Proxy = n.NodeP2P.Proxy                
	case n.NodeP2P.ProxyUser != :
		joined.ProxyUser = n.NodeP2P.ProxyUser            
	case n.NodeP2P.ProxyPass != :
		joined.ProxyPass = n.NodeP2P.ProxyPass            
	case n.NodeP2P.OnionProxy != :
		joined.OnionProxy = n.NodeP2P.OnionProxy           
	case n.NodeP2P.OnionProxyUser != :
		joined.OnionProxyUser = n.NodeP2P.OnionProxyUser       
	case n.NodeP2P.OnionProxyPass != :
		joined.OnionProxyPass = n.NodeP2P.OnionProxyPass       
	case n.NodeP2P.NoOnion != :
		joined.NoOnion = n.NodeP2P.NoOnion              
	case n.NodeP2P.TorIsolation != :
		joined.TorIsolation = n.NodeP2P.TorIsolation         
	case cfg.Network.TestNet3 != :
		joined.TestNet3 = cfg.Network.TestNet3             
	case cfg.Network.RegressionTest != :
		joined.RegressionTest = cfg.Network.RegressionTest       
	case cfg.Network.SimNet != :
		joined.SimNet = cfg.Network.SimNet               
	case n.NodeChain.AddCheckpoints != :
		joined.AddCheckpoints = n.NodeChain.AddCheckpoints       
	case n.NodeChain.DisableCheckpoints != :
		joined.DisableCheckpoints = n.NodeChain.DisableCheckpoints   
	case n.NodeLaunch.Profile != :
		joined.Profile = n.NodeLaunch.Profile              
	case n.NodeLaunch.CPUProfile != :
		joined.CPUProfile = n.NodeLaunch.CPUProfile           
	case n.NodeP2P.Upnp != :
		joined.Upnp = n.NodeP2P.Upnp                 
	case n.NodeChain.MinRelayTxFee != :
		joined.MinRelayTxFee = n.NodeChain.MinRelayTxFee        
	case n.NodeChain.NoRelayPriority != :
		joined.NoRelayPriority = n.NodeChain.NoRelayPriority      
	case n.NodeMining.Generate != :
		joined.Generate = n.NodeMining.Generate             
	case n.NodeMining.MiningAddrs != :
		joined.MiningAddrs = n.NodeMining.MiningAddrs          
	case n.NodeP2P.UserAgentComments != :
		joined.UserAgentComments = n.NodeP2P.UserAgentComments    
	case n.NodeP2P.NoPeerBloomFilters != :
		joined.NoPeerBloomFilters = n.NodeP2P.NoPeerBloomFilters   
	case n.NodeP2P.NoCFilters != :
		joined.NoCFilters = n.NodeP2P.NoCFilters           
	case n.NodeLaunch.DropCfIndex != :
		joined.DropCfIndex = n.NodeLaunch.DropCfIndex          
	case n.NodeP2P.BlocksOnly != :
		joined.BlocksOnly = n.NodeP2P.BlocksOnly           
	case n.NodeChain.TxIndex != :
		joined.TxIndex = n.NodeChain.TxIndex              
	case n.NodeLaunch.DropTxIndex != :
		joined.DropTxIndex = n.NodeLaunch.DropTxIndex          
	case n.NodeChain.AddrIndex != :
		joined.AddrIndex = n.NodeChain.AddrIndex            
	case n.NodeLaunch.DropAddrIndex != :
		joined.DropAddrIndex = n.NodeLaunch.DropAddrIndex        
	case n.NodeP2P.RelayNonStd != :
		joined.RelayNonStd = n.NodeP2P.RelayNonStd          
	case n.NodeP2P.RejectNonStd != :
		joined.RejectNonStd = n.NodeP2P.RejectNonStd         
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

func defaultNodeCfg() (out node.Config) {
	out = node.Config{
		ConfigFile:           node.DefaultConfigFile,
		DataDir:              node.DefaultDataDir,
		LogDir:               node.DefaultLogDir,
		Listeners:            []string{node.DefaultListener},
		MaxPeers:             node.DefaultMaxPeers,
		BanDuration:          node.DefaultBanDuration,
		BanThreshold:         node.DefaultBanThreshold,
		RPCUser:              "user",
		RPCPass:              "pa55word",
		RPCLimitUser:         "",
		RPCLimitPass:         "",
		RPCListeners:         []string{node.DefaultRPCListener},
		RPCCert:              node.DefaultRPCCertFile,
		RPCMaxClients:        node.DefaultMaxRPCClients,
		RPCMaxWebsockets:     node.DefaultMaxRPCWebsockets,
		RPCMaxConcurrentReqs: node.DefaultMaxRPCConcurrentReqs,
		DbType:               node.DefaultDbType,
		FreeTxRelayLimit:     node.DefaultFreeTxRelayLimit,
		TrickleInterval:      node.DefaultTrickleInterval,
		MaxOrphanTxs:         node.DefaultMaxOrphanTransactions,
		Algo:                 node.DefaultAlgo,
		GenThreads:           node.DefaultGenThreads,
		MinerListener:        node.DefaultMinerListener,
		MinerPass:            "pa55word",
		BlockMinSize:         node.DefaultBlockMinSize,
		BlockMaxSize:         node.DefaultBlockMaxSize,
		BlockMinWeight:       node.DefaultBlockMinWeight,
		BlockMaxWeight:       node.DefaultBlockMaxWeight,
		BlockPrioritySize:    mempool.DefaultBlockPrioritySize,
		SigCacheMaxSize:      node.DefaultSigCacheMaxSize,
	}
	return
}
