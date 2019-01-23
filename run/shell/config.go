package shell

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"git.parallelcoin.io/pod/lib/clog"
	n "git.parallelcoin.io/pod/module/node"
	"git.parallelcoin.io/pod/module/node/mempool"
	w "git.parallelcoin.io/pod/module/wallet"
	ww "git.parallelcoin.io/pod/module/wallet/wallet"
	"git.parallelcoin.io/pod/run/logger"
	"git.parallelcoin.io/pod/run/util"
	"github.com/tucnak/climax"
)

// Log is the shell main logger
var Log = cl.NewSubSystem("run/shell", "info")
var log = Log.Ch

// UseLogger uses a specified Logger to output package logging info. This should be used in preference to SetLogWriter if the caller is also using log.
func UseLogger(logger *cl.SubSystem) {
	Log = logger
	log = Log.Ch
}

// Cfg is the combined app and logging configuration data
type Cfg struct {
	DataDir      string
	AppDataDir   string
	ConfFileName string
	Node         *n.Config
	Wallet       *w.Config
	Levels       map[string]string
}

var (
	DefaultDataDir      = n.DefaultDataDir
	DefaultAppDataDir   = filepath.Join(n.DefaultHomeDir, "shell")
	DefaultConfFileName = filepath.Join(filepath.Join(n.DefaultHomeDir, "shell"), "conf")
)

// Config is the combined app and log levels configuration
var Config = DefaultConfig()
var f = pu.GenFlag
var t = pu.GenTrig
var s = pu.GenShort
var l = pu.GenLog

// Command is a command to send RPC queries to bitcoin RPC protocol server for node and wallet queries
var Command = climax.Command{
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
			for i := range logger.Levels {
				logger.Levels[i].SetLevel(dl)
			}
		}
		log <- cl.Debugf{
			"pod/shell version %s",
			Version(),
		}
		if ctx.Is("version") {
			fmt.Println("shell version", Version())
			fmt.Println("pod version", n.Version())
			fmt.Println("wallet version", w.Version())
			cl.Shutdown()
		}
		var cfgFile string
		if cfgFile, ok = ctx.Get("configfile"); !ok {
			cfgFile = DefaultConfFileName
		}
		if ctx.Is("init") {
			log <- cl.Debugf{
				"writing default configuration to %s", cfgFile,
			}
			WriteDefaultConfig(cfgFile)
			configNode(&ctx, cfgFile)
		} else {
			log <- cl.Infof{
				"loading configuration from %s",
				cfgFile,
			}
			if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
				log <- cl.Wrn(
					"configuration file does not exist, creating new one",
				)
				WriteDefaultConfig(cfgFile)
				configNode(&ctx, cfgFile)
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
				err = json.Unmarshal(cfgData, &Config)
				if err != nil {
					log <- cl.Error{"parsing app configuration:", err.Error()}
					cl.Shutdown()
				}
				configNode(&ctx, cfgFile)
			}
		}
		runShell()
		cl.Shutdown()
		return 0
	},
}

func getIfIs(ctx *climax.Context, name string, r *string) (ok bool) {
	if ctx.Is(name) {
		var s string
		s, ok = ctx.Get(name)
		r = &s
	}
	return
}

func configNode(ctx *climax.Context, cfgFile string) {
	var r *string
	t := ""
	r = &t

	// Node and general stuff
	if getIfIs(ctx, "debuglevel", r) {
		switch *r {
		case "fatal", "error", "warn", "info", "debug", "trace":
			Config.Node.DebugLevel = *r
		default:
			Config.Node.DebugLevel = "info"
		}
		Log.SetLevel(Config.Node.DebugLevel)
	}
	if getIfIs(ctx, "datadir", r) {
		Config.Node.DataDir = n.CleanAndExpandPath(*r)
	}
	if getIfIs(ctx, "addpeers", r) {
		pu.NormalizeAddresses(*r, n.DefaultPort, &Config.Node.AddPeers)
	}
	if getIfIs(ctx, "connectpeers", r) {
		pu.NormalizeAddresses(*r, n.DefaultPort, &Config.Node.ConnectPeers)
	}
	if getIfIs(ctx, "disablelisten", r) {
		Config.Node.DisableListen = *r == "true"
	}
	if getIfIs(ctx, "listeners", r) {
		pu.NormalizeAddresses(*r, n.DefaultPort, &Config.Node.Listeners)
	}
	if getIfIs(ctx, "maxpeers", r) {
		if err := pu.ParseInteger(*r, "maxpeers", &Config.Node.MaxPeers); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if getIfIs(ctx, "disablebanning", r) {
		Config.Node.DisableBanning = *r == "true"
	}
	if getIfIs(ctx, "banduration", r) {
		if err := pu.ParseDuration(*r, "banduration", &Config.Node.BanDuration); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if getIfIs(ctx, "banthreshold", r) {
		var bt int
		if err := pu.ParseInteger(*r, "banthtreshold", &bt); err != nil {
			log <- cl.Wrn(err.Error())
		} else {
			Config.Node.BanThreshold = uint32(bt)
		}
	}
	if getIfIs(ctx, "whitelists", r) {
		pu.NormalizeAddresses(*r, n.DefaultPort, &Config.Node.Whitelists)
	}
	// if getIfIs(ctx, "rpcuser", r) {
	// 	Config.Node.RPCUser = *r
	// }
	// if getIfIs(ctx, "rpcpass", r) {
	// 	Config.Node.RPCPass = *r
	// }
	// if getIfIs(ctx, "rpclimituser", r) {
	// 	Config.Node.RPCLimitUser = *r
	// }
	// if getIfIs(ctx, "rpclimitpass", r) {
	// 	Config.Node.RPCLimitPass = *r
	// }
	// if getIfIs(ctx, "rpclisteners", r) {
	pu.NormalizeAddresses(n.DefaultRPCListener, n.DefaultRPCPort, &Config.Node.RPCListeners)
	// }
	// if getIfIs(ctx, "rpccert", r) {
	// 	Config.Node.RPCCert = n.CleanAndExpandPath(*r)
	// }
	// if getIfIs(ctx, "rpckey", r) {
	// 	Config.Node.RPCKey = n.CleanAndExpandPath(*r)
	// }
	// if getIfIs(ctx, "tls", r) {
	// Config.Node.TLS = *r == "true"
	// }
	Config.Node.TLS = false
	if getIfIs(ctx, "disablednsseed", r) {
		Config.Node.DisableDNSSeed = *r == "true"
	}
	if getIfIs(ctx, "externalips", r) {
		pu.NormalizeAddresses(*r, n.DefaultPort, &Config.Node.ExternalIPs)
	}
	if getIfIs(ctx, "proxy", r) {
		pu.NormalizeAddress(*r, "9050", &Config.Node.Proxy)
	}
	if getIfIs(ctx, "proxyuser", r) {
		Config.Node.ProxyUser = *r
	}
	if getIfIs(ctx, "proxypass", r) {
		Config.Node.ProxyPass = *r
	}
	if getIfIs(ctx, "onion", r) {
		pu.NormalizeAddress(*r, "9050", &Config.Node.OnionProxy)
	}
	if getIfIs(ctx, "onionuser", r) {
		Config.Node.OnionProxyUser = *r
	}
	if getIfIs(ctx, "onionpass", r) {
		Config.Node.OnionProxyPass = *r
	}
	if getIfIs(ctx, "noonion", r) {
		Config.Node.NoOnion = *r == "true"
	}
	if getIfIs(ctx, "torisolation", r) {
		Config.Node.TorIsolation = *r == "true"
	}
	if getIfIs(ctx, "network", r) {
		switch *r {
		case "testnet":
			Config.Node.TestNet3, Config.Node.RegressionTest, Config.Node.SimNet = true, false, false
		case "regtest":
			Config.Node.TestNet3, Config.Node.RegressionTest, Config.Node.SimNet = false, true, false
		case "simnet":
			Config.Node.TestNet3, Config.Node.RegressionTest, Config.Node.SimNet = false, false, true
		default:
			Config.Node.TestNet3, Config.Node.RegressionTest, Config.Node.SimNet = false, false, false
		}
	}
	if getIfIs(ctx, "addcheckpoints", r) {
		Config.Node.AddCheckpoints = strings.Split(*r, " ")
	}
	if getIfIs(ctx, "disablecheckpoints", r) {
		Config.Node.DisableCheckpoints = *r == "true"
	}
	if getIfIs(ctx, "dbtype", r) {
		Config.Node.DbType = *r
	}
	if getIfIs(ctx, "profile", r) {
		Config.Node.Profile = n.NormalizeAddress(*r, "11034")
	}
	if getIfIs(ctx, "cpuprofile", r) {
		Config.Node.CPUProfile = n.NormalizeAddress(*r, "11033")
	}
	if getIfIs(ctx, "upnp", r) {
		Config.Node.Upnp = *r == "true"
	}
	if getIfIs(ctx, "minrelaytxfee", r) {
		if err := pu.ParseFloat(*r, "minrelaytxfee", &Config.Node.MinRelayTxFee); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if getIfIs(ctx, "freetxrelaylimit", r) {
		if err := pu.ParseFloat(*r, "freetxrelaylimit", &Config.Node.FreeTxRelayLimit); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if getIfIs(ctx, "norelaypriority", r) {
		Config.Node.NoRelayPriority = *r == "true"
	}
	if getIfIs(ctx, "trickleinterval", r) {
		if err := pu.ParseDuration(*r, "trickleinterval", &Config.Node.TrickleInterval); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if getIfIs(ctx, "maxorphantxs", r) {
		if err := pu.ParseInteger(*r, "maxorphantxs", &Config.Node.MaxOrphanTxs); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if getIfIs(ctx, "algo", r) {
		Config.Node.Algo = *r
	}
	if getIfIs(ctx, "generate", r) {
		Config.Node.Generate = *r == "true"
	}
	if getIfIs(ctx, "genthreads", r) {
		var gt int
		if err := pu.ParseInteger(*r, "genthreads", &gt); err != nil {
			log <- cl.Wrn(err.Error())
		} else {
			Config.Node.GenThreads = int32(gt)
		}
	}
	if getIfIs(ctx, "miningaddrs", r) {
		Config.Node.MiningAddrs = strings.Split(*r, " ")
	}
	if getIfIs(ctx, "minerlistener", r) {
		pu.NormalizeAddress(*r, n.DefaultRPCPort, &Config.Node.MinerListener)
	}
	if getIfIs(ctx, "minerpass", r) {
		Config.Node.MinerPass = *r
	}
	if getIfIs(ctx, "blockminsize", r) {
		if err := pu.ParseUint32(*r, "blockminsize", &Config.Node.BlockMinSize); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if getIfIs(ctx, "blockmaxsize", r) {
		if err := pu.ParseUint32(*r, "blockmaxsize", &Config.Node.BlockMaxSize); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if getIfIs(ctx, "blockminweight", r) {
		if err := pu.ParseUint32(*r, "blockminweight", &Config.Node.BlockMinWeight); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if getIfIs(ctx, "blockmaxweight", r) {
		if err := pu.ParseUint32(*r, "blockmaxweight", &Config.Node.BlockMaxWeight); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if getIfIs(ctx, "blockprioritysize", r) {
		if err := pu.ParseUint32(*r, "blockmaxweight", &Config.Node.BlockPrioritySize); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if getIfIs(ctx, "uacomment", r) {
		Config.Node.UserAgentComments = strings.Split(*r, " ")
	}
	if getIfIs(ctx, "nopeerbloomfilters", r) {
		Config.Node.NoPeerBloomFilters = *r == "true"
	}
	if getIfIs(ctx, "nocfilters", r) {
		Config.Node.NoCFilters = *r == "true"
	}
	if ctx.Is("dropcfindex") {
		Config.Node.DropCfIndex = true
	}
	if getIfIs(ctx, "sigcachemaxsize", r) {
		var scms int
		if err := pu.ParseInteger(*r, "sigcachemaxsize", &scms); err != nil {
			log <- cl.Wrn(err.Error())
		} else {
			Config.Node.SigCacheMaxSize = uint(scms)
		}
	}
	if getIfIs(ctx, "blocksonly", r) {
		Config.Node.BlocksOnly = *r == "true"
	}
	if getIfIs(ctx, "txindex", r) {
		Config.Node.TxIndex = *r == "true"
	}
	if ctx.Is("droptxindex") {
		Config.Node.DropTxIndex = true
	}
	if ctx.Is("addrindex") {
		r, _ := ctx.Get("addrindex")
		Config.Node.AddrIndex = r == "true"
	}
	if ctx.Is("dropaddrindex") {
		Config.Node.DropAddrIndex = true
	}
	if getIfIs(ctx, "relaynonstd", r) {
		Config.Node.RelayNonStd = *r == "true"
	}
	if getIfIs(ctx, "rejectnonstd", r) {
		Config.Node.RejectNonStd = *r == "true"
	}

	// Wallet stuff

	if ctx.Is("create") {
		Config.Wallet.Create = true
	}
	if ctx.Is("createtemp") {
		Config.Wallet.CreateTemp = true
	}
	if getIfIs(ctx, "appdatadir", r) {
		Config.Wallet.AppDataDir = n.CleanAndExpandPath(*r)
	}
	if getIfIs(ctx, "noinitialload", r) {
		Config.Wallet.NoInitialLoad = *r == "true"
	}
	if getIfIs(ctx, "logdir", r) {
		Config.Wallet.LogDir = n.CleanAndExpandPath(*r)
	}
	if getIfIs(ctx, "profile", r) {
		pu.NormalizeAddress(*r, "3131", &Config.Wallet.Profile)
	}
	if getIfIs(ctx, "gui", r) {
		Config.Wallet.GUI = *r == "true"
	}
	if getIfIs(ctx, "walletpass", r) {
		Config.Wallet.WalletPass = *r
	}
	// if getIfIs(ctx, "rpcconnect", r) {
	pu.NormalizeAddress(n.DefaultRPCListener, "11048", &Config.Wallet.RPCConnect)
	// }
	if getIfIs(ctx, "cafile", r) {
		Config.Wallet.CAFile = n.CleanAndExpandPath(*r)
	}
	// if getIfIs(ctx, "enableclienttls", r) {
	// 	Config.Wallet.EnableClientTLS = *r == "true"
	// }
	// if getIfIs(ctx, "podusername", r) {
	Config.Wallet.PodUsername = Config.Node.RPCUser
	// }
	// if getIfIs(ctx, "podpassword", r) {
	Config.Wallet.PodPassword = Config.Node.RPCPass
	// }
	if getIfIs(ctx, "onetimetlskey", r) {
		Config.Wallet.OneTimeTLSKey = *r == "true"
	}
	if getIfIs(ctx, "enableservertls", r) {
		Config.Wallet.EnableServerTLS = *r == "true"
	}
	if getIfIs(ctx, "legacyrpclisteners", r) {
		pu.NormalizeAddresses(*r, "11046", &Config.Wallet.LegacyRPCListeners)
	}
	if getIfIs(ctx, "legacyrpcmaxclients", r) {
		var bt int
		if err := pu.ParseInteger(*r, "legacyrpcmaxclients", &bt); err != nil {
			log <- cl.Wrn(err.Error())
		} else {
			Config.Wallet.LegacyRPCMaxClients = int64(bt)
		}
	}
	if getIfIs(ctx, "legacyrpcmaxwebsockets", r) {
		_, err := fmt.Sscanf(*r, "%d", Config.Wallet.LegacyRPCMaxWebsockets)
		if err != nil {
			log <- cl.Errorf{
				`malformed legacyrpcmaxwebsockets: "%s" leaving set at "%d"`,
				r,
				Config.Wallet.LegacyRPCMaxWebsockets,
			}
		}
	}
	if getIfIs(ctx, "username", r) {
		Config.Wallet.Username = *r
	}
	if getIfIs(ctx, "password", r) {
		Config.Wallet.Password = *r
	}
	if getIfIs(ctx, "experimentalrpclisteners", r) {
		pu.NormalizeAddresses(*r, "11045", &Config.Wallet.ExperimentalRPCListeners)
	}
	if getIfIs(ctx, "network", r) {
		switch *r {
		case "testnet":
			Config.Wallet.TestNet3, Config.Wallet.SimNet = true, false
		case "simnet":
			Config.Wallet.TestNet3, Config.Wallet.SimNet = false, true
		default:
			Config.Wallet.TestNet3, Config.Wallet.SimNet = false, false
		}
	}

	logger.SetLogging(ctx)
	if ctx.Is("save") {
		log <- cl.Infof{
			"saving config file to %s",
			cfgFile,
		}
		j, err := json.MarshalIndent(Config, "", "  ")
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

// WriteConfig creates and writes the config file in the requested location
func WriteConfig(cfgFile string, c *Cfg) {
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

// WriteDefaultConfig creates and writes a default config to the specified path
func WriteDefaultConfig(cfgFile string) {
	defCfg := DefaultConfig()
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
	Config = defCfg
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Cfg {
	rpcusername := pu.GenKey()
	rpcpassword := pu.GenKey()
	return &Cfg{
		DataDir:      DefaultDataDir,
		AppDataDir:   DefaultAppDataDir,
		ConfFileName: DefaultConfFileName,
		Node: &n.Config{
			RPCUser:              rpcusername,
			RPCPass:              rpcpassword,
			RPCMaxClients:        n.DefaultMaxRPCClients,
			RPCMaxWebsockets:     n.DefaultMaxRPCWebsockets,
			RPCMaxConcurrentReqs: n.DefaultMaxRPCConcurrentReqs,
			DbType:               n.DefaultDbType,
			RPCListeners:         []string{"127.0.0.1:11048"},
			TLS:                  false,
			MinRelayTxFee:        mempool.DefaultMinRelayTxFee.ToDUO(),
			FreeTxRelayLimit:     n.DefaultFreeTxRelayLimit,
			TrickleInterval:      n.DefaultTrickleInterval,
			BlockMinSize:         n.DefaultBlockMinSize,
			BlockMaxSize:         n.DefaultBlockMaxSize,
			BlockMinWeight:       n.DefaultBlockMinWeight,
			BlockMaxWeight:       n.DefaultBlockMaxWeight,
			BlockPrioritySize:    mempool.DefaultBlockPrioritySize,
			MaxOrphanTxs:         n.DefaultMaxOrphanTransactions,
			SigCacheMaxSize:      n.DefaultSigCacheMaxSize,
			Generate:             n.DefaultGenerate,
			GenThreads:           1,
			TxIndex:              n.DefaultTxIndex,
			AddrIndex:            n.DefaultAddrIndex,
			Algo:                 n.DefaultAlgo,
		},
		Wallet: &w.Config{
			NoInitialLoad:          true,
			RPCConnect:             "127.0.0.1:11048",
			PodUsername:            rpcusername,
			PodPassword:            rpcpassword,
			RPCKey:                 w.DefaultRPCKeyFile,
			RPCCert:                w.DefaultRPCCertFile,
			WalletPass:             ww.InsecurePubPassphrase,
			EnableClientTLS:        false,
			LegacyRPCMaxClients:    w.DefaultRPCMaxClients,
			LegacyRPCMaxWebsockets: w.DefaultRPCMaxWebsockets,
		},
		Levels: logger.GetDefaultConfig(),
	}
}
