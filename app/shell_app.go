package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	n "git.parallelcoin.io/pod/cmd/node"
	"git.parallelcoin.io/pod/cmd/node/mempool"
	w "git.parallelcoin.io/pod/cmd/wallet"
	"git.parallelcoin.io/pod/pkg/netparams"
	cl "git.parallelcoin.io/pod/pkg/clog"
	"github.com/tucnak/climax"
)

// ShellCfg is the combined app and logging configuration data
type ShellCfg struct {
	Node   *n.Config
	Wallet *w.Config
	Levels map[string]string
}

// ShellConfig is the combined app and log levels configuration
var ShellConfig = DefaultShellConfig()

// ShellCommand is a command to send RPC queries to bitcoin RPC protocol server for node and wallet queries
var ShellCommand = climax.Command{
	Name:  "shell",
	Brief: "parallelcoin shell",
	Help:  "check balances, make payments, manage contacts, search the chain, it slices, it dices",
	Flags: []climax.Flag{
		t("version", "V", "show version number and quit"),

		s("configfile", "C", DefaultShellConfFileName, "path to configuration file"),
		s("datadir", "D", n.DefaultDataDir, "set the pod base directory"),
		f("appdatadir", "shell", "set app data directory for wallet, configuration and logs"),

		t("init", "i", "resets configuration to defaults"),
		t("save", "S", "saves current flags into configuration"),

		f("noinitialload", "false", "defer wallet load to be triggered by RPC"),
		f("network", "mainnet", "connect to (mainnet|testnet|regtestnet|simnet)"),

		t("create", "", "create a new wallet if it does not exist"),
		f("createtemp", "false", "create temporary wallet (pass=walletpass) requires --datadir"),

		f("gui", "false", "launch GUI"),

		f("walletpass", "", "the public wallet password - only required if the wallet was created with one"),

		f("profile", "false", "enable HTTP profiling on given port (1024-65536)"),

		f("legacyrpclisteners", "127.0.0.1:11046", "add a listener for the legacy RPC"),
		f("legacyrpcmaxclients", fmt.Sprint(n.DefaultMaxRPCClients),
			"max connections for legacy RPC"),
		f("legacyrpcmaxwebsockets", fmt.Sprint(n.DefaultMaxRPCWebsockets),
			"max websockets for legacy RPC"),

		f("username", "user", "username for wallet RPC"),
		f("password", "pa55word", "password for wallet RPC"),

		f("rpccert", n.DefaultRPCCertFile,
			"file containing the RPC tls certificate"),
		f("rpckey", n.DefaultRPCKeyFile,
			"file containing RPC TLS key"),
		f("onetimetlskey", "false",
			"generate a new TLS certpair don't save key"),
		f("cafile", w.DefaultCAFile,
			"certificate authority for custom TLS CA"),
		f("enableservertls", "false", "enable TLS on wallet RPC server"),

		f("proxy", "", "proxy address for outbound connections"),
		f("proxyuser", "", "username for proxy server"),
		f("proxypass", "", "password for proxy server"),

		f("txindex", "true", "enable transaction index"),
		f("addrindex", "true", "enable address index"),
		t("dropcfindex", "", "delete committed filtering (CF) index then exit"),
		t("droptxindex", "", "deletes transaction index then exit"),
		t("dropaddrindex", "", "deletes the address index then exits"),

		s("listeners", "S", n.DefaultListener, "sets an address to listen for P2P connections"),
		f("externalips", "", "additional P2P listeners"),
		f("disablelisten", "false", "disables the P2P listener"),

		f("addpeers", "", "adds a peer to the peers database to try to connect to"),
		f("connectpeers", "", "adds a peer to a connect-only whitelist"),
		f(`maxpeers`, fmt.Sprint(n.DefaultMaxPeers),
			"sets max number of peers to connect to to at once"),
		f(`disablebanning`, "false",
			"disable banning of misbehaving peers"),
		f("banduration", "1d",
			"time to ban misbehaving peers (d/h/m/s)"),
		f("banthreshold", fmt.Sprint(n.DefaultBanThreshold),
			"banscore that triggers a ban"),
		f("whitelists", "", "addresses and networks immune to banning"),

		f("onion", "", "connect via tor proxy relay"),
		f("onionuser", "", "username for onion proxy server"),
		f("onionpass", "", "password for onion proxy server"),
		f("noonion", "false", "disable onion proxy"),
		f("torisolation", "false", "use a different user/pass for each peer"),

		f("trickleinterval", fmt.Sprint(n.DefaultTrickleInterval),
			"time between sending inventory batches to peers"),
		f("minrelaytxfee", "0",
			"min fee in DUO/kb to relay transaction"),
		f("freetxrelaylimit", fmt.Sprint(n.DefaultFreeTxRelayLimit),
			"limit below min fee transactions in kb/bin"),
		f("norelaypriority", "false",
			"do not discriminate transactions for relaying"),

		f("nopeerbloomfilters", "false",
			"disable bloom filtering support"),
		f("nocfilters", "false",
			"disable committed filtering (CF) support"),
		f("blocksonly", "false", "do not accept transactions from peers"),
		f("relaynonstd", "false", "relay nonstandard transactions"),
		f("rejectnonstd", "false", "reject nonstandard transactions"),

		f("maxorphantxs", fmt.Sprint(n.DefaultMaxOrphanTransactions),
			"max number of orphan transactions to store"),
		f("sigcachemaxsize", fmt.Sprint(n.DefaultSigCacheMaxSize),
			"maximum number of signatures to store in memory"),

		f("generate", fmt.Sprint(n.DefaultGenerate),
			"set CPU miner to generate blocks"),
		f("genthreads", fmt.Sprint(n.DefaultGenThreads),
			"set number of threads to generate using CPU, -1 = all"),
		f("algo", n.DefaultAlgo, "set algorithm to be used by cpu miner"),
		f("miningaddrs", "", "add address to pay block rewards to"),
		f("minerlistener", n.DefaultMinerListener,
			"address to listen for mining work subscriptions"),
		f("minerpass", "",
			"PSK to prevent snooping/spoofing of miner traffic"),

		f("addcheckpoints", "false", `add custom checkpoints "height:hash"`),
		f("disablecheckpoints", "false", "disable all checkpoints"),

		f("blockminsize", fmt.Sprint(n.DefaultBlockMinSize),
			"min block size for miners"),
		f("blockmaxsize", fmt.Sprint(n.DefaultBlockMaxSize),
			"max block size for miners"),
		f("blockminweight", fmt.Sprint(n.DefaultBlockMinWeight),
			"min block weight for miners"),
		f("blockmaxweight", fmt.Sprint(n.DefaultBlockMaxWeight),
			"max block weight for miners"),
		f("blockprioritysize", fmt.Sprint(),
			"size in bytes of high priority blocks"),

		f("uacomment", "", "comment to add to the P2P network user agent"),
		f("upnp", "false", "use UPNP to automatically port forward to node"),
		f("dbtype", "ffldb", "set database backend type"),
		f("disablednsseed", "false", "disable dns seeding"),

		f("profile", "false", "start HTTP profiling server on given address"),
		f("cpuprofile", "false", "start cpu profiling server on given address"),

		s("debuglevel", "d", "info", "sets debuglevel, specify per-library below"),

		l("lib-addrmgr"), l("lib-blockchain"), l("lib-connmgr"), l("lib-database-ffldb"), l("lib-database"), l("lib-mining-cpuminer"), l("lib-mining"), l("lib-netsync"), l("lib-peer"), l("lib-rpcclient"), l("lib-txscript"), l("node"), l("node-mempool"), l("spv"), l("wallet"), l("wallet-chain"), l("wallet-legacyrpc"), l("wallet-rpcserver"), l("wallet-tx"), l("wallet-votingpool"), l("wallet-waddrmgr"), l("wallet-wallet"), l("wallet-wtxmgr"),
	},
	Examples: []climax.Example{
		{
			Usecase:     "--init --rpcuser=user --rpcpass=pa55word --save",
			Description: "resets the configuration file to default, sets rpc username and password and saves the changes to config after parsing",
		},
	},
	Handle: func(ctx climax.Context) int {
		var dl string
		var ok bool
		if dl, ok = ctx.Get("debuglevel"); ok {
			log <- cl.Tracef{"setting debug level %s", dl}
			Log.SetLevel(dl)
			ll := GetAllSubSystems()
			for i := range ll {
				ll[i].SetLevel(dl)
			}
		}
		log <- cl.Trc("starting wallet app")
		log <- cl.Debugf{"pod/wallet version %s", w.Version()}
		if ctx.Is("version") {
			fmt.Println("pod/wallet version", w.Version())
			cl.Shutdown()
		}
		var cfgFile string
		if cfgFile, ok = ctx.Get("configfile"); !ok {
			cfgFile = w.DefaultConfigFile
		}
		if ctx.Is("init") {
			log <- cl.Debug{"writing default configuration to", cfgFile}
			WriteDefaultShellConfig(cfgFile)
		}
		log <- cl.Info{"loading configuration from", cfgFile}
		if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
			log <- cl.Wrn("configuration file does not exist, creating new one")
			WriteDefaultShellConfig(cfgFile)
		} else {
			log <- cl.Debug{"reading app configuration from", cfgFile}
			cfgData, err := ioutil.ReadFile(cfgFile)
			if err != nil {
				log <- cl.Error{"reading app config file", err.Error()}
				WriteDefaultShellConfig(cfgFile)
			}
			log <- cl.Tracef{"parsing app configuration\n%s", cfgData}
			err = json.Unmarshal(cfgData, &ShellConfig)
			if err != nil {
				log <- cl.Error{"parsing app config file", err.Error()}
				WriteDefaultShellConfig(cfgFile)
			}
		}
		configShell(&ctx, cfgFile)
		runShell(ctx.Args)
		cl.Shutdown()
		return 0
	},
}

// ShellFlags is the list of flags and the default values stored in the Usage field
var ShellFlags = GetFlags(ShellCommand)

func configShell(ctx *climax.Context, cfgFile string) {
	log <- cl.Trace{"configuring from command line flags ", os.Args}
	if ctx.Is("create") {
		log <- cl.Dbg("request to make new wallet")
		ShellConfig.Wallet.Create = true
	}
	if ctx.Is("createtemp") {
		log <- cl.Dbg("request to make temp wallet")
		ShellConfig.Wallet.CreateTemp = true
	}
	if r, ok := getIfIs(ctx, "appdatadir"); ok {
		log <- cl.Debug{"appdatadir set to", r}
		ShellConfig.Wallet.AppDataDir = n.CleanAndExpandPath(r)
	}
	if r, ok := getIfIs(ctx, "noinitialload"); ok {
		log <- cl.Dbg("no initial load requested")
		ShellConfig.Wallet.NoInitialLoad = r == "true"
	}
	if r, ok := getIfIs(ctx, "logdir"); ok {
		log <- cl.Dbg("")
		ShellConfig.Wallet.LogDir = n.CleanAndExpandPath(r)
	}
	if r, ok := getIfIs(ctx, "profile"); ok {
		log <- cl.Dbg("")
		NormalizeAddress(r, "3131", &ShellConfig.Wallet.Profile)
	}
	if r, ok := getIfIs(ctx, "gui"); ok {
		log <- cl.Dbg("")
		ShellConfig.Wallet.GUI = r == "true"
	}
	if r, ok := getIfIs(ctx, "walletpass"); ok {
		log <- cl.Dbg("")
		ShellConfig.Wallet.WalletPass = r
	}
	if r, ok := getIfIs(ctx, "rpcconnect"); ok {
		log <- cl.Dbg("")
		NormalizeAddress(r, "11048", &ShellConfig.Wallet.RPCConnect)
	}
	if r, ok := getIfIs(ctx, "cafile"); ok {
		log <- cl.Dbg("")
		ShellConfig.Wallet.CAFile = n.CleanAndExpandPath(r)
	}
	if r, ok := getIfIs(ctx, "enableclienttls"); ok {
		log <- cl.Dbg("")
		ShellConfig.Wallet.EnableClientTLS = r == "true"
	}
	if r, ok := getIfIs(ctx, "podusername"); ok {
		log <- cl.Dbg("")
		ShellConfig.Wallet.PodUsername = r
	}
	if r, ok := getIfIs(ctx, "podpassword"); ok {
		log <- cl.Dbg("")
		ShellConfig.Wallet.PodPassword = r
	}
	if r, ok := getIfIs(ctx, "proxy"); ok {
		log <- cl.Dbg("")
		NormalizeAddress(r, "11048", &ShellConfig.Wallet.Proxy)
	}
	if r, ok := getIfIs(ctx, "proxyuser"); ok {
		log <- cl.Dbg("")
		ShellConfig.Wallet.ProxyUser = r
	}
	if r, ok := getIfIs(ctx, "proxypass"); ok {
		log <- cl.Dbg("")
		ShellConfig.Wallet.ProxyPass = r
	}
	if r, ok := getIfIs(ctx, "rpccert"); ok {
		log <- cl.Dbg("")
		ShellConfig.Wallet.RPCCert = n.CleanAndExpandPath(r)
	}
	if r, ok := getIfIs(ctx, "rpckey"); ok {
		log <- cl.Dbg("")
		ShellConfig.Wallet.RPCKey = n.CleanAndExpandPath(r)
	}
	if r, ok := getIfIs(ctx, "onetimetlskey"); ok {
		log <- cl.Dbg("")
		ShellConfig.Wallet.OneTimeTLSKey = r == "true"
	}
	if r, ok := getIfIs(ctx, "enableservertls"); ok {
		log <- cl.Dbg("")
		ShellConfig.Wallet.EnableServerTLS = r == "true"
	}
	if r, ok := getIfIs(ctx, "legacyrpclisteners"); ok {
		log <- cl.Dbg("")
		NormalizeAddresses(r, "11046", &ShellConfig.Wallet.LegacyRPCListeners)
	}
	if r, ok := getIfIs(ctx, "legacyrpcmaxclients"); ok {
		log <- cl.Dbg("")
		var bt int
		if err := ParseInteger(r, "legacyrpcmaxclients", &bt); err != nil {
			log <- cl.Wrn(err.Error())
		} else {
			ShellConfig.Wallet.LegacyRPCMaxClients = int64(bt)
		}
	}
	if r, ok := getIfIs(ctx, "legacyrpcmaxwebsockets"); ok {
		log <- cl.Dbg("")
		_, err := fmt.Sscanf(r, "%d", ShellConfig.Wallet.LegacyRPCMaxWebsockets)
		if err != nil {
			log <- cl.Errorf{
				"malformed legacyrpcmaxwebsockets: `%s` leaving set at `%d`",
				r, ShellConfig.Wallet.LegacyRPCMaxWebsockets,
			}
		}
	}
	if r, ok := getIfIs(ctx, "username"); ok {
		log <- cl.Dbg("")
		ShellConfig.Wallet.Username = r
	}
	if r, ok := getIfIs(ctx, "password"); ok {
		log <- cl.Dbg("")
		ShellConfig.Wallet.Password = r
	}
	if r, ok := getIfIs(ctx, "experimentalrpclisteners"); ok {
		log <- cl.Dbg("")
		NormalizeAddresses(r, "11045", &ShellConfig.Wallet.ExperimentalRPCListeners)
	}
	if r, ok := getIfIs(ctx, "datadir"); ok {
		log <- cl.Dbg("")
		ShellConfig.Wallet.DataDir = r
	}
	if r, ok := getIfIs(ctx, "network"); ok {
		log <- cl.Dbg("")
		switch r {
		case "testnet":
			ShellConfig.Wallet.TestNet3, ShellConfig.Wallet.SimNet = true, false
			w.ActiveNet = &netparams.TestNet3Params
		case "simnet":
			ShellConfig.Wallet.TestNet3, ShellConfig.Wallet.SimNet = false, true
			w.ActiveNet = &netparams.SimNetParams
		default:
			ShellConfig.Wallet.TestNet3, ShellConfig.Wallet.SimNet = false, false
			w.ActiveNet = &netparams.MainNetParams
		}
	}

	// finished configuration

	SetLogging(ctx)

	if ctx.Is("save") {
		log <- cl.Info{"saving config file to", cfgFile}
		j, err := json.MarshalIndent(ShellConfig, "", "  ")
		if err != nil {
			log <- cl.Error{"writing app config file", err}
		}
		j = append(j, '\n')
		log <- cl.Trace{"JSON formatted config file\n", string(j)}
		ioutil.WriteFile(cfgFile, j, 0600)
	}
}

// WriteShellConfig creates and writes the config file in the requested location
func WriteShellConfig(cfgFile string, c *ShellCfg) {
	log <- cl.Dbg("writing config")
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

// WriteDefaultShellConfig creates and writes a default config to the requested location
func WriteDefaultShellConfig(cfgFile string) {
	defCfg := DefaultShellConfig()
	defCfg.Wallet.ConfigFile = cfgFile
	j, err := json.MarshalIndent(defCfg, "", "  ")
	if err != nil {
		log <- cl.Error{"marshalling configuration", err}
		panic(err)
	}
	j = append(j, '\n')
	log <- cl.Trace{"JSON formatted config file\n", string(j)}
	err = ioutil.WriteFile(cfgFile, j, 0600)
	if err != nil {
		log <- cl.Error{"writing app config file", err}
		panic(err)
	}
	// if we are writing default config we also want to use it
	ShellConfig = defCfg
}

// DefaultShellConfig returns a default configuration
func DefaultShellConfig() *ShellCfg {
	log <- cl.Dbg("getting default config")
	return &ShellCfg{
		Node: &n.Config{
			RPCUser:              "",
			RPCPass:              "",
			Listeners:            []string{n.DefaultListener},
			RPCListeners:         []string{n.DefaultRPCListener},
			DebugLevel:           "info",
			ConfigFile:           n.DefaultConfigFile,
			MaxPeers:             n.DefaultMaxPeers,
			BanDuration:          n.DefaultBanDuration,
			BanThreshold:         n.DefaultBanThreshold,
			RPCMaxClients:        n.DefaultMaxRPCClients,
			RPCMaxWebsockets:     n.DefaultMaxRPCWebsockets,
			RPCMaxConcurrentReqs: n.DefaultMaxRPCConcurrentReqs,
			DataDir:              n.DefaultDataDir,
			LogDir:               n.DefaultLogDir,
			DbType:               n.DefaultDbType,
			RPCKey:               n.DefaultRPCKeyFile,
			RPCCert:              n.DefaultRPCCertFile,
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
			NoInitialLoad:          false,
			ConfigFile:             w.DefaultConfigFile,
			DataDir:                w.DefaultDataDir,
			AppDataDir:             w.DefaultAppDataDir,
			LogDir:                 w.DefaultLogDir,
			RPCKey:                 w.DefaultRPCKeyFile,
			RPCCert:                w.DefaultRPCCertFile,
			WalletPass:             "",
			CAFile:                 "",
			LegacyRPCMaxClients:    w.DefaultRPCMaxClients,
			LegacyRPCMaxWebsockets: w.DefaultRPCMaxWebsockets,
		},
		Levels: GetDefaultLogLevelsConfig(),
	}
}
