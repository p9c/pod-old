package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"git.parallelcoin.io/pod/cmd/node"
	"git.parallelcoin.io/pod/cmd/node/mempool"
	walletmain "git.parallelcoin.io/pod/cmd/wallet"
	ww "git.parallelcoin.io/pod/cmd/wallet/wallet"
	cl "git.parallelcoin.io/pod/pkg/clog"
	"github.com/tucnak/climax"
)

// ShellCfg is the combined app and logging configuration data
type ShellCfg struct {
	DataDir      string
	AppDataDir   string
	ConfFileName string
	Node         *node.Config
	Wallet       *walletmain.Config
	Levels       map[string]string
}

var (
	// ShellConfig is the combined app and log levels configuration
	ShellConfig = DefaultShellConfig()
)

// ShellCommand is a command to send RPC queries to bitcoin RPC protocol server for node and wallet queries
var ShellCommand = climax.Command{
	Name:  "shell",
	Brief: "parallelcoin combined full node and wallet",
	Help:  "distrubutes, verifies and mines blocks for the parallelcoin duo cryptocurrency, as well as optionally providing search indexes for transactions in the database, and provides RPC and GUI interfaces for a built-in wallet",
	Flags: []climax.Flag{
		t("version", "V", "show version number and quit"),

		s("datadir", "D", "base path for pod data and configurations"),
		f("appdatadir", "path to store application data in"),
		t("configfile", "C", "path to configuration file"),

		f("init", "resets configuration to defaults"),
		f("save", "saves current configuration"),

		f("gui", "launch GUI"),

		s("username", "-u", "username for wallet RPC"),
		s("password", "-P", "password for wallet RPC"),
		f("walletpass", "public data password for wallet"),

		f("listeners", "sets an address to listen for P2P connections"),
		f("disablelisten", "disables the P2P listener"),

		f("addpeers", "adds a peer to the peers database to try to connect to"),
		f("connectpeers", "adds a peer to a connect-only whitelist"),

		f("network", "connect to (mainnet|testnet|simnet)"),

		f("legacyrpclisteners", "add a listener for the legacy RPC"),
		f("legacyrpcmaxclients", "max connections for legacy RPC"),
		f("legacyrpcmaxwebsockets", "max websockets for legacy RPC"),
		f("experimentalrpclisteners", "address for experimental RPC listener"),

		f("proxy", "SOCKS5 proxy address"),
		f("proxyuser", "username for proxy server"),
		f("proxypass", "password for proxy server"),

		f("onion", "tor proxy address"),
		f("onionuser", "username for tor proxy"),
		f("onionpass", "password for tor proxy"),
		f("noonion", "disable onion proxy (if user/pass is set)"),
		f("torisolation", "use a different tor session for each peer"),

		f("rpccert", "file containing the RPC TLS certificate"),
		f("rpckey", "file containing RPC TLS key"),
		f("cafile", "custom certificate authority for TLS"),
		f("onetimetlskey", "generate TLS certificates but don't save key"),
		f("enableservertls", "enable TLS on wallet RPC"),
		f("skipverify", `do not verify tls certificates`),

		f("noinitialload", "launch without unlocking a wallet"),
		f("disablednsseed", "disable dns seeding"),
		f("upnp", "use uPNP to auto-configure NAT port redirection"),

		f("txindex", "enable transaction search API"),
		f("addrindex", "enable address search API"),

		f("create", "create a new wallet if it does not exist"),
		f("createtemp", "create temporary wallet (pass=walletpass)"),

		t("dropcfindex", "", "deletes the committed filtering (CF) and exits"),
		t("droptxindex", "", "deletes the transaction index and exits"),
		t("dropaddrindex", "", "deletes the address index and exits"),

		f("trickleinterval", "time between sending messages to a peer"),
		f("maxpeers", "sets max number of peers to open connect to at once"),
		f("disablebanning", "disable banning of misbehaving peers"),
		f("banduration", "time to enforce ban on a misbehaving peer"),
		f("banthreshold", "ban score above which a ban is triggered"),
		f("whitelists", "ip or network in which peers are not banned"),
		f("externalips", "extra listeners on different address/interfaces"),

		f("minrelaytxfee", "the minimum fee in DUO/Kb to relay a transaction"),
		f("freetxrelaylimit", "kb/min of sub-minimum tx fees"),
		f("norelaypriority", "relay transactions regardless of fee size"),
		f("maxorphantxs", "max orphan transactions to store in memory"),

		f("addcheckpoints", `add custom checkpoints "height:hash"`),
		f("disablecheckpoints", "disable all checkpoints"),

		f("blockminsize", "minimum block size for miners"),
		f("blockmaxsize", "max block size for miners"),
		f("blockminweight", "mininum block weight for miners"),
		f("blockmaxweight", "max block weight for miners"),
		f("blockprioritysize", "max size low fee transactions get priority"),

		f("generate", "set CPU miner to generate blocks"),
		f("genthreads", "set number of threads to generate blocks with"),
		f("algo", "set algorithm to be used by cpu miner"),
		f("miningaddrs", "addresses to pay to on mined blocks"),
		f("minerlistener", "port for miner controller subscriptions"),
		f("minerpass", "encrypt miner traffic for insecure networks"),

		f("relaynonstd", "relay non-standard transactions"),
		f("nopeerbloomfilters", "disable bloom filters"),
		f("nocfilters", "disable committed filtering (CF) support"),
		f("blocksonly", "do not accept transactions from remote peers"),
		f("rejectnonstd", "reject non-standard transactions"),

		f("uacomment", "comment to add to the P2P network user agent string"),
		f("dbtype", "set database backend type"),
		f("sigcachemaxsize", "maxi number of signature cache entries to cache"),

		f("profile", "start HTTP profiling server on given address"),
		f("cpuprofile", "start CPU profiling server on given address"),

		s("debuglevel", "d", "sets debuglevel base (logging is per-library)"),
		l("lib-addrmgr"), l("lib-blockchain"), l("lib-connmgr"), l("lib-database"), l("lib-database"), l("lib-mining"), l("lib-mining"), l("lib-netsync"), l("lib-peer"), l("lib-rpcclient"), l("lib-txscript"), l("node"), l("node-mempool"), l("spv"), l("wallet"), l("wallet-chain"), l("wallet-legacyrpc"), l("wallet-rpcserver"), l("wallet-tx"), l("wallet-votingpool"), l("wallet-waddrmgr"), l("wallet-wallet"), l("wallet-wtxmgr"),
	},
	Examples: []climax.Example{
		// {
		// 	Usecase:     "--init --rpcuser=user --rpcpass=pa55word --save",
		// 	Description: "resets the configuration file to default, sets rpc username and password and saves the changes to config after parsing",
		// },
	},
	Handle: func(ctx climax.Context) int {
		var dl string
		var ok bool
		if dl, ok = ctx.Get("debuglevel"); ok {
			log <- cl.Tracef{
				"setting debug level %s",
				dl,
			}
			Log.SetLevel(dl)
			ll := GetAllSubSystems()
			for i := range ShellConfig.Levels {
				ll[i].SetLevel(dl)
			}
		}
		log <- cl.Debugf{
			"pod/shell version %s",
			Version(),
		}
		if ctx.Is("version") {
			fmt.Println("shell version", Version())
			fmt.Println("pod version", node.Version())
			fmt.Println("wallet version", walletmain.Version())
			cl.Shutdown()
		}
		var cfgFile string
		if cfgFile, ok = ctx.Get("configfile"); !ok {
			cfgFile = DefaultShellConfFileName
		}
		if ctx.Is("init") {
			log <- cl.Debugf{
				"writing default configuration to %s", cfgFile,
			}
			WriteDefaultShellConfig(cfgFile)
			configShell(&ctx, cfgFile)
		} else {
			log <- cl.Infof{
				"loading configuration from %s",
				cfgFile,
			}
			if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
				log <- cl.Wrn(
					"configuration file does not exist, creating new one",
				)
				WriteDefaultShellConfig(cfgFile)
				configShell(&ctx, cfgFile)
			} else {
				log <- cl.Debug{
					"reading app configuration from", cfgFile,
				}
				cfgData, err := ioutil.ReadFile(cfgFile)
				if err != nil {
					log <- cl.Error{"reading app config file:", err.Error()}
					cl.Shutdown()
				}
				log <- cl.Tracef{"parsing app configuration\n%s", cfgData}
				err = json.Unmarshal(cfgData, &ShellConfig)
				log <- cl.Dbg("finished processing config file")
				if err != nil {
					log <- cl.Error{"parsing app configuration:", err.Error()}
					cl.Shutdown()
				}
				configShell(&ctx, cfgFile)
			}
		}
		runShell(ctx.Args)
		cl.Shutdown()
		return 0
	},
}

func configShell(ctx *climax.Context, cfgFile string) {
	fmt.Println("configuring from command line flags")
	// Node and general stuff
	if r, ok := getIfIs(ctx, "debuglevel"); ok {
		switch r {
		case "fatal", "error", "warn", "info", "debug", "trace":
			ShellConfig.Node.DebugLevel = r
		default:
			ShellConfig.Node.DebugLevel = "info"
		}
		Log.SetLevel(ShellConfig.Node.DebugLevel)
	}
	if r, ok := getIfIs(ctx, "datadir"); ok {
		ShellConfig.Node.DataDir = node.CleanAndExpandPath(r)
	}
	if r, ok := getIfIs(ctx, "addpeers"); ok {
		NormalizeAddresses(r, node.DefaultPort, &ShellConfig.Node.AddPeers)
	}
	if r, ok := getIfIs(ctx, "connectpeers"); ok {
		NormalizeAddresses(r, node.DefaultPort, &ShellConfig.Node.ConnectPeers)
	}
	if r, ok := getIfIs(ctx, "disablelisten"); ok {
		ShellConfig.Node.DisableListen = r == "true"
	}
	if r, ok := getIfIs(ctx, "listeners"); ok {
		NormalizeAddresses(r, node.DefaultPort, &ShellConfig.Node.Listeners)
	}
	if r, ok := getIfIs(ctx, "maxpeers"); ok {
		if err := ParseInteger(r, "maxpeers", &ShellConfig.Node.MaxPeers); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if r, ok := getIfIs(ctx, "disablebanning"); ok {
		ShellConfig.Node.DisableBanning = r == "true"
	}
	if r, ok := getIfIs(ctx, "banduration"); ok {
		if err := ParseDuration(r, "banduration", &ShellConfig.Node.BanDuration); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if r, ok := getIfIs(ctx, "banthreshold"); ok {
		var bt int
		if err := ParseInteger(r, "banthtreshold", &bt); err != nil {
			log <- cl.Wrn(err.Error())
		} else {
			ShellConfig.Node.BanThreshold = uint32(bt)
		}
	}
	if r, ok := getIfIs(ctx, "whitelists"); ok {
		NormalizeAddresses(r, node.DefaultPort, &ShellConfig.Node.Whitelists)
	}
	// if getIfIs(ctx, "rpcuser");ok {
	// 	ShellConfig.Node.RPCUser = r
	// }
	// if getIfIs(ctx, "rpcpass");ok {
	// 	ShellConfig.Node.RPCPass = r
	// }
	// if getIfIs(ctx, "rpclimituser");ok {
	// 	ShellConfig.Node.RPCLimitUser = r
	// }
	// if getIfIs(ctx, "rpclimitpass");ok {
	// 	ShellConfig.Node.RPCLimitPass = r
	// }
	// if getIfIs(ctx, "rpclisteners");ok {
	NormalizeAddresses(node.DefaultRPCListener, node.DefaultRPCPort, &ShellConfig.Node.RPCListeners)
	// }
	// if getIfIs(ctx, "rpccert");ok {
	// 	ShellConfig.Node.RPCCert = node.CleanAndExpandPath(r)
	// }
	// if getIfIs(ctx, "rpckey");ok {
	// 	ShellConfig.Node.RPCKey = node.CleanAndExpandPath(r)
	// }
	// if getIfIs(ctx, "tls");ok {
	// ShellConfig.Node.TLS = r == "true"
	// }
	ShellConfig.Node.TLS = false
	if r, ok := getIfIs(ctx, "disablednsseed"); ok {
		ShellConfig.Node.DisableDNSSeed = r == "true"
	}
	if r, ok := getIfIs(ctx, "externalips"); ok {
		NormalizeAddresses(r, node.DefaultPort, &ShellConfig.Node.ExternalIPs)
	}
	if r, ok := getIfIs(ctx, "proxy"); ok {
		NormalizeAddress(r, "9050", &ShellConfig.Node.Proxy)
	}
	if r, ok := getIfIs(ctx, "proxyuser"); ok {
		ShellConfig.Node.ProxyUser = r
	}
	if r, ok := getIfIs(ctx, "proxypass"); ok {
		ShellConfig.Node.ProxyPass = r
	}
	if r, ok := getIfIs(ctx, "onion"); ok {
		NormalizeAddress(r, "9050", &ShellConfig.Node.OnionProxy)
	}
	if r, ok := getIfIs(ctx, "onionuser"); ok {
		ShellConfig.Node.OnionProxyUser = r
	}
	if r, ok := getIfIs(ctx, "onionpass"); ok {
		ShellConfig.Node.OnionProxyPass = r
	}
	if r, ok := getIfIs(ctx, "noonion"); ok {
		ShellConfig.Node.NoOnion = r == "true"
	}
	if r, ok := getIfIs(ctx, "torisolation"); ok {
		ShellConfig.Node.TorIsolation = r == "true"
	}
	if r, ok := getIfIs(ctx, "network"); ok {
		switch r {
		case "testnet":
			ShellConfig.Node.TestNet3, ShellConfig.Node.RegressionTest, ShellConfig.Node.SimNet = true, false, false
		case "regtest":
			ShellConfig.Node.TestNet3, ShellConfig.Node.RegressionTest, ShellConfig.Node.SimNet = false, true, false
		case "simnet":
			ShellConfig.Node.TestNet3, ShellConfig.Node.RegressionTest, ShellConfig.Node.SimNet = false, false, true
		default:
			ShellConfig.Node.TestNet3, ShellConfig.Node.RegressionTest, ShellConfig.Node.SimNet = false, false, false
		}
	}
	if r, ok := getIfIs(ctx, "addcheckpoints"); ok {
		ShellConfig.Node.AddCheckpoints = strings.Split(r, " ")
	}
	if r, ok := getIfIs(ctx, "disablecheckpoints"); ok {
		ShellConfig.Node.DisableCheckpoints = r == "true"
	}
	if r, ok := getIfIs(ctx, "dbtype"); ok {
		ShellConfig.Node.DbType = r
	}
	if r, ok := getIfIs(ctx, "profile"); ok {
		ShellConfig.Node.Profile = node.NormalizeAddress(r, "11034")
	}
	if r, ok := getIfIs(ctx, "cpuprofile"); ok {
		ShellConfig.Node.CPUProfile = node.NormalizeAddress(r, "11033")
	}
	if r, ok := getIfIs(ctx, "upnp"); ok {
		ShellConfig.Node.Upnp = r == "true"
	}
	if r, ok := getIfIs(ctx, "minrelaytxfee"); ok {
		if err := ParseFloat(r, "minrelaytxfee", &ShellConfig.Node.MinRelayTxFee); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if r, ok := getIfIs(ctx, "freetxrelaylimit"); ok {
		if err := ParseFloat(r, "freetxrelaylimit", &ShellConfig.Node.FreeTxRelayLimit); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if r, ok := getIfIs(ctx, "norelaypriority"); ok {
		ShellConfig.Node.NoRelayPriority = r == "true"
	}
	if r, ok := getIfIs(ctx, "trickleinterval"); ok {
		if err := ParseDuration(r, "trickleinterval", &ShellConfig.Node.TrickleInterval); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if r, ok := getIfIs(ctx, "maxorphantxs"); ok {
		if err := ParseInteger(r, "maxorphantxs", &ShellConfig.Node.MaxOrphanTxs); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if r, ok := getIfIs(ctx, "algo"); ok {
		ShellConfig.Node.Algo = r
	}
	if r, ok := getIfIs(ctx, "generate"); ok {
		ShellConfig.Node.Generate = r == "true"
	}
	if r, ok := getIfIs(ctx, "genthreads"); ok {
		var gt int
		if err := ParseInteger(r, "genthreads", &gt); err != nil {
			log <- cl.Wrn(err.Error())
		} else {
			ShellConfig.Node.GenThreads = int32(gt)
		}
	}
	if r, ok := getIfIs(ctx, "miningaddrs"); ok {
		ShellConfig.Node.MiningAddrs = strings.Split(r, " ")
	}
	if r, ok := getIfIs(ctx, "minerlistener"); ok {
		NormalizeAddress(r, node.DefaultRPCPort, &ShellConfig.Node.MinerListener)
	}
	if r, ok := getIfIs(ctx, "minerpass"); ok {
		ShellConfig.Node.MinerPass = r
	}
	if r, ok := getIfIs(ctx, "blockminsize"); ok {
		if err := ParseUint32(r, "blockminsize", &ShellConfig.Node.BlockMinSize); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if r, ok := getIfIs(ctx, "blockmaxsize"); ok {
		if err := ParseUint32(r, "blockmaxsize", &ShellConfig.Node.BlockMaxSize); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if r, ok := getIfIs(ctx, "blockminweight"); ok {
		if err := ParseUint32(r, "blockminweight", &ShellConfig.Node.BlockMinWeight); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if r, ok := getIfIs(ctx, "blockmaxweight"); ok {
		if err := ParseUint32(r, "blockmaxweight", &ShellConfig.Node.BlockMaxWeight); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if r, ok := getIfIs(ctx, "blockprioritysize"); ok {
		if err := ParseUint32(r, "blockmaxweight", &ShellConfig.Node.BlockPrioritySize); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if r, ok := getIfIs(ctx, "uacomment"); ok {
		ShellConfig.Node.UserAgentComments = strings.Split(r, " ")
	}
	if r, ok := getIfIs(ctx, "nopeerbloomfilters"); ok {
		ShellConfig.Node.NoPeerBloomFilters = r == "true"
	}
	if r, ok := getIfIs(ctx, "nocfilters"); ok {
		ShellConfig.Node.NoCFilters = r == "true"
	}
	if ctx.Is("dropcfindex") {
		ShellConfig.Node.DropCfIndex = true
	}
	if r, ok := getIfIs(ctx, "sigcachemaxsize"); ok {
		var scms int
		if err := ParseInteger(r, "sigcachemaxsize", &scms); err != nil {
			log <- cl.Wrn(err.Error())
		} else {
			ShellConfig.Node.SigCacheMaxSize = uint(scms)
		}
	}
	if r, ok := getIfIs(ctx, "blocksonly"); ok {
		ShellConfig.Node.BlocksOnly = r == "true"
	}
	if r, ok := getIfIs(ctx, "txindex"); ok {
		ShellConfig.Node.TxIndex = r == "true"
	}
	if ctx.Is("droptxindex") {
		ShellConfig.Node.DropTxIndex = true
	}
	if r, ok := getIfIs(ctx, "addrindex"); ok {
		ShellConfig.Node.AddrIndex = r == "true"
	}
	if ctx.Is("dropaddrindex") {
		ShellConfig.Node.DropAddrIndex = true
	}
	if r, ok := getIfIs(ctx, "relaynonstd"); ok {
		ShellConfig.Node.RelayNonStd = r == "true"
	}
	if r, ok := getIfIs(ctx, "rejectnonstd"); ok {
		ShellConfig.Node.RejectNonStd = r == "true"
	}

	// Wallet stuff

	if ctx.Is("create") {
		ShellConfig.Wallet.Create = true
	}
	if ctx.Is("createtemp") {
		ShellConfig.Wallet.CreateTemp = true
	}
	if r, ok := getIfIs(ctx, "appdatadir"); ok {
		ShellConfig.Wallet.AppDataDir = node.CleanAndExpandPath(r)
	}
	if r, ok := getIfIs(ctx, "noinitialload"); ok {
		ShellConfig.Wallet.NoInitialLoad = r == "true"
	}
	if r, ok := getIfIs(ctx, "logdir"); ok {
		ShellConfig.Wallet.LogDir = node.CleanAndExpandPath(r)
	}
	if r, ok := getIfIs(ctx, "profile"); ok {
		NormalizeAddress(r, "3131", &ShellConfig.Wallet.Profile)
	}
	if r, ok := getIfIs(ctx, "gui"); ok {
		ShellConfig.Wallet.GUI = r == "true"
	}
	if r, ok := getIfIs(ctx, "walletpass"); ok {
		ShellConfig.Wallet.WalletPass = r
	}
	// if getIfIs(ctx, "rpcconnect");ok {
	NormalizeAddress(node.DefaultRPCListener, "11048", &ShellConfig.Wallet.RPCConnect)
	// }
	if r, ok := getIfIs(ctx, "cafile"); ok {
		ShellConfig.Wallet.CAFile = node.CleanAndExpandPath(r)
	}
	// if getIfIs(ctx, "enableclienttls");ok {
	// 	ShellConfig.Wallet.EnableClientTLS = r == "true"
	// }
	// if getIfIs(ctx, "podusername");ok {
	ShellConfig.Wallet.PodUsername = ShellConfig.Node.RPCUser
	// }
	// if getIfIs(ctx, "podpassword");ok {
	ShellConfig.Wallet.PodPassword = ShellConfig.Node.RPCPass
	// }
	if r, ok := getIfIs(ctx, "onetimetlskey"); ok {
		ShellConfig.Wallet.OneTimeTLSKey = r == "true"
	}
	if r, ok := getIfIs(ctx, "enableservertls"); ok {
		ShellConfig.Wallet.EnableServerTLS = r == "true"
	}
	if r, ok := getIfIs(ctx, "legacyrpclisteners"); ok {
		NormalizeAddresses(r, "11046", &ShellConfig.Wallet.LegacyRPCListeners)
	}
	if r, ok := getIfIs(ctx, "legacyrpcmaxclients"); ok {
		var bt int
		if err := ParseInteger(r, "legacyrpcmaxclients", &bt); err != nil {
			log <- cl.Wrn(err.Error())
		} else {
			ShellConfig.Wallet.LegacyRPCMaxClients = int64(bt)
		}
	}
	if r, ok := getIfIs(ctx, "legacyrpcmaxwebsockets"); ok {
		_, err := fmt.Sscanf(r, "%d", ShellConfig.Wallet.LegacyRPCMaxWebsockets)
		if err != nil {
			log <- cl.Errorf{
				`malformed legacyrpcmaxwebsockets: "%s" leaving set at "%d"`,
				r,
				ShellConfig.Wallet.LegacyRPCMaxWebsockets,
			}
		}
	}
	if r, ok := getIfIs(ctx, "username"); ok {
		ShellConfig.Wallet.Username = r
	}
	if r, ok := getIfIs(ctx, "password"); ok {
		ShellConfig.Wallet.Password = r
	}
	if r, ok := getIfIs(ctx, "experimentalrpclisteners"); ok {
		NormalizeAddresses(r, "11045", &ShellConfig.Wallet.ExperimentalRPCListeners)
	}
	if r, ok := getIfIs(ctx, "network"); ok {
		switch r {
		case "testnet":
			ShellConfig.Wallet.TestNet3, ShellConfig.Wallet.SimNet = true, false
		case "simnet":
			ShellConfig.Wallet.TestNet3, ShellConfig.Wallet.SimNet = false, true
		default:
			ShellConfig.Wallet.TestNet3, ShellConfig.Wallet.SimNet = false, false
		}
	}

	SetLogging(ctx)
	if ctx.Is("save") {
		log <- cl.Infof{
			"saving config file to %s",
			cfgFile,
		}
		j, err := json.MarshalIndent(ShellConfig, "", "  ")
		if err != nil {
			log <- cl.Error{"saving config file:", err.Error()}
		}
		j = append(j, '\n')
		log <- cl.Tracef{"JSON formatted config file\n%s", j}
		err = ioutil.WriteFile(cfgFile, j, 0600)
		if err != nil {
			log <- cl.Error{"writing app config file:", err.Error()}
		}
	}
}

// WriteShellConfig creates and writes the config file in the requested location
func WriteShellConfig(cfgFile string, c *ShellCfg) {
	j, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		panic(err.Error())
	}
	j = append(j, '\n')
	err = ioutil.WriteFile(cfgFile, j, 0600)
	if err != nil {
		panic(err.Error())
	}
}

// WriteDefaultShellConfig creates and writes a default config to the specified path
func WriteDefaultShellConfig(cfgFile string) {
	defCfg := DefaultShellConfig()
	defCfg.ConfFileName = cfgFile
	j, err := json.MarshalIndent(defCfg, "", "  ")
	if err != nil {
		log <- cl.Error{"marshalling default config:" + err.Error()}
	}
	j = append(j, '\n')
	log <- cl.Tracef{"JSON formatted config file\n%s", j}
	err = ioutil.WriteFile(cfgFile, j, 0600)
	if err != nil {
		log <- cl.Error{"writing default config:", err.Error()}
	}
	// if we are writing default config we also want to use it
	ShellConfig = defCfg
}

// DefaultShellConfig returns a default configuration
func DefaultShellConfig() *ShellCfg {
	rpcusername := GenKey()
	rpcpassword := GenKey()
	return &ShellCfg{
		DataDir:      DefaultDataDir,
		AppDataDir:   DefaultShellDataDir,
		ConfFileName: DefaultShellConfFileName,
		Node: &node.Config{
			RPCUser:              rpcusername,
			RPCPass:              rpcpassword,
			RPCMaxClients:        node.DefaultMaxRPCClients,
			RPCMaxWebsockets:     node.DefaultMaxRPCWebsockets,
			RPCMaxConcurrentReqs: node.DefaultMaxRPCConcurrentReqs,
			DbType:               node.DefaultDbType,
			RPCListeners:         []string{"127.0.0.1:11048"},
			TLS:                  false,
			MinRelayTxFee:        mempool.DefaultMinRelayTxFee.ToDUO(),
			FreeTxRelayLimit:     node.DefaultFreeTxRelayLimit,
			TrickleInterval:      node.DefaultTrickleInterval,
			BlockMinSize:         node.DefaultBlockMinSize,
			BlockMaxSize:         node.DefaultBlockMaxSize,
			BlockMinWeight:       node.DefaultBlockMinWeight,
			BlockMaxWeight:       node.DefaultBlockMaxWeight,
			BlockPrioritySize:    mempool.DefaultBlockPrioritySize,
			MaxOrphanTxs:         node.DefaultMaxOrphanTransactions,
			SigCacheMaxSize:      node.DefaultSigCacheMaxSize,
			Generate:             node.DefaultGenerate,
			GenThreads:           1,
			TxIndex:              node.DefaultTxIndex,
			AddrIndex:            node.DefaultAddrIndex,
			Algo:                 node.DefaultAlgo,
		},
		Wallet: &walletmain.Config{
			NoInitialLoad:          true,
			RPCConnect:             "127.0.0.1:11048",
			PodUsername:            rpcusername,
			PodPassword:            rpcpassword,
			RPCKey:                 walletmain.DefaultRPCKeyFile,
			RPCCert:                walletmain.DefaultRPCCertFile,
			WalletPass:             ww.InsecurePubPassphrase,
			EnableClientTLS:        false,
			LegacyRPCMaxClients:    walletmain.DefaultRPCMaxClients,
			LegacyRPCMaxWebsockets: walletmain.DefaultRPCMaxWebsockets,
		},
		Levels: GetDefaultLogLevelsConfig(),
	}
}
