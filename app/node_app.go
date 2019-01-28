package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"git.parallelcoin.io/pod/cmd/node"
	n "git.parallelcoin.io/pod/cmd/node"
	"git.parallelcoin.io/pod/cmd/node/mempool"
	"git.parallelcoin.io/pod/pkg/blockchain"
	cl "git.parallelcoin.io/pod/pkg/clog"
	"git.parallelcoin.io/pod/pkg/connmgr"
	"git.parallelcoin.io/pod/pkg/fork"
	"git.parallelcoin.io/pod/pkg/util"
	"github.com/btcsuite/go-socks/socks"
	"github.com/davecgh/go-spew/spew"
	"github.com/tucnak/climax"
)

// serviceOptions defines the configuration options for the daemon as a service on Windows.
type serviceOptions struct {
	ServiceCommand string `short:"s" long:"service" description:"Service command {install, remove, start, stop}"`
}

// StateCfg is a reference to the main node state configuration struct
var StateCfg = n.StateCfg

// runServiceCommand is only set to a real function on Windows.  It is used to parse and execute service commands specified via the -s flag.
var runServiceCommand func(string) error

var aN = filepath.Base(os.Args[0])
var appName = strings.TrimSuffix(aN, filepath.Ext(aN))

var usageMessage = fmt.Sprintf("use `%s help node` to show usage", appName)

// NodeCfg is the combined app and logging configuration data
type NodeCfg struct {
	Node      *n.Config
	LogLevels map[string]string
}

// NodeConfig is the combined app and log levels configuration
var NodeConfig = DefaultNodeConfig()

// NodeCommand is a command to send RPC queries to bitcoin RPC protocol server for node and wallet queries
var NodeCommand = climax.Command{
	Name:  "node",
	Brief: "parallelcoin full node",
	Help:  "distrubutes, verifies and mines blocks for the parallelcoin duo cryptocurrency, as well as optionally providing search indexes for transactions in the database",
	Flags: []climax.Flag{

		t("version", "V", "show version number and quit"),

		s("configfile", "C", n.DefaultConfigFile, "path to configuration file"),
		s("datadir", "D", n.DefaultDataDir, "path to configuration directory"),

		t("init", "", "resets configuration to defaults"),
		t("save", "", "saves current configuration"),

		f("network", "mainnet", "connect to (mainnet|testnet|simnet)"),

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
		f(`maxpeers`, fmt.Sprint(node.DefaultMaxPeers),
			"sets max number of peers to connect to to at once"),
		f(`disablebanning`, "false",
			"disable banning of misbehaving peers"),
		f("banduration", "1d",
			"time to ban misbehaving peers (d/h/m/s)"),
		f("banthreshold", fmt.Sprint(node.DefaultBanThreshold),
			"banscore that triggers a ban"),
		f("whitelists", "", "addresses and networks immune to banning"),

		s("rpcuser", "u", "user", "RPC username"),
		s("rpcpass", "P", "pa55word", "RPC password"),

		f("rpclimituser", "user", "limited user RPC username"),
		f("rpclimitpass", "pa55word", "limited user RPC password"),

		s("rpclisteners", "s", node.DefaultRPCListener, "RPC server to connect to"),

		f("rpccert", node.DefaultRPCCertFile,
			"RPC server tls certificate chain for validation"),
		f("rpckey", node.DefaultRPCKeyFile,
			"RPC server tls key for authentication"),
		f("tls", "false", "enable TLS"),
		f("skipverify", "false", "do not verify tls certificates"),

		f("proxy", "", "connect via SOCKS5 proxy server"),
		f("proxyuser", "", "username for proxy server"),
		f("proxypass", "", "password for proxy server"),

		f("onion", "", "connect via tor proxy relay"),
		f("onionuser", "", "username for onion proxy server"),
		f("onionpass", "", "password for onion proxy server"),
		f("noonion", "false", "disable onion proxy"),
		f("torisolation", "false", "use a different user/pass for each peer"),

		f("trickleinterval", fmt.Sprint(node.DefaultTrickleInterval),
			"time between sending inventory batches to peers"),
		f("minrelaytxfee", "0", "min fee in DUO/kb to relay transaction"),
		f("freetxrelaylimit", fmt.Sprint(node.DefaultFreeTxRelayLimit),
			"limit below min fee transactions in kb/bin"),
		f("norelaypriority", "false",
			"do not discriminate transactions for relaying"),

		f("nopeerbloomfilters", "false",
			"disable bloom filtering support"),
		f("nocfilters", "false",
			"disable committed filtering (CF) support"),
		f("blocksonly", "false", "do not accept transactions from peers"),
		f("relaynonstd", "false", "relay nonstandard transactions"),
		f("rejectnonstd", "true", "reject nonstandard transactions"),

		f("maxorphantxs", fmt.Sprint(node.DefaultMaxOrphanTransactions),
			"max number of orphan transactions to store"),
		f("sigcachemaxsize", fmt.Sprint(node.DefaultSigCacheMaxSize),
			"maximum number of signatures to store in memory"),

		f("generate", "false", "set CPU miner to generate blocks"),
		f("genthreads", "-1", "set number of threads to generate using CPU, -1 = all"),
		f("algo", "random", "set algorithm to be used by cpu miner"),
		f("miningaddrs", "", "add address to pay block rewards to"),
		f("minerlistener", node.DefaultMinerListener,
			"address to listen for mining work subscriptions"),
		f("minerpass", "", "Preshared Key to prevent snooping/spoofing of miner traffic"),

		f("addcheckpoints", "", `add custom checkpoints "height:hash"`),
		f("disablecheckpoints", "", "disable all checkpoints"),

		f("blockminsize", fmt.Sprint(node.DefaultBlockMinSize),
			"min block size for miners"),
		f("blockmaxsize", fmt.Sprint(node.DefaultBlockMaxSize),
			"max block size for miners"),
		f("blockminweight", fmt.Sprint(node.DefaultBlockMinWeight),
			"min block weight for miners"),
		f("blockmaxweight", fmt.Sprint(node.DefaultBlockMaxWeight),
			"max block weight for miners"),
		f("blockprioritysize", "0", "size in bytes of high priority blocks"),

		f("uacomment", "", "comment to add to the P2P network user agent"),
		f("upnp", "false", "use UPNP to automatically port forward to node"),
		f("dbtype", "ffldb", "set database backend type"),
		f("disablednsseed", "false", "disable dns seeding"),

		f("profile", "false", "start HTTP profiling server on given address"),
		f("cpuprofile", "false", "start cpu profiling server on given address"),

		s("debuglevel", "d", "info", "sets log level for those unspecified below"),

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
			log <- cl.Tracef{
				"setting debug level %s",
				dl,
			}
			CtlCfg.DebugLevel = dl
			Log.SetLevel(dl)
			ll := GetAllSubSystems()
			for i := range ll {
				ll[i].SetLevel(dl)
			}
		}
		log <- cl.Debugf{"pod/node version %s", n.Version()}
		if ctx.Is("version") {
			fmt.Println("pod/node version", n.Version())
			cl.Shutdown()
		}
		var cfgFile string
		if cfgFile, ok = ctx.Get("configfile"); !ok {
			cfgFile = n.DefaultConfigFile
		}
		if ctx.Is("init") {
			log <- cl.Debugf{
				"writing default configuration to %s",
				cfgFile,
			}
			WriteDefaultNodeConfig(cfgFile)
		} else {
			log <- cl.Infof{
				"loading configuration from %s",
				cfgFile,
			}
			if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
				log <- cl.Warn{"configuration file does not exist, creating new one"}
				WriteDefaultNodeConfig(cfgFile)
			} else {
				log <- cl.Debug{
					"reading app configuration from",
					cfgFile,
				}
				cfgData, err := ioutil.ReadFile(cfgFile)
				if err != nil {
					log <- cl.Error{
						"reading app config file:",
						err.Error(),
					}
					cl.Shutdown()
				}
				log <- cl.Trace{
					"parsing app configuration",
					cfgData,
				}
				err = json.Unmarshal(cfgData, &NodeConfig)
				if err != nil {
					log <- cl.Error{
						"parsing app config file:",
						err.Error(),
					}
					WriteDefaultNodeConfig(cfgFile)
				}
			}
		}
		configNode(&ctx, cfgFile)
		runNode()
		cl.Shutdown()
		return 0
	},
}

func configNode(ctx *climax.Context, cfgFile string) {
	var err error
	if r, ok := getIfIs(ctx, "debuglevel"); ok {
		switch r {
		case "fatal", "error", "warn", "info", "debug", "trace":
			NodeConfig.Node.DebugLevel = r
		default:
			NodeConfig.Node.DebugLevel = "info"
		}
		Log.SetLevel(NodeConfig.Node.DebugLevel)
	}
	if r, ok := getIfIs(ctx, "datadir"); ok {
		NodeConfig.Node.DataDir = n.CleanAndExpandPath(r)
	}
	if r, ok := getIfIs(ctx, "addpeers"); ok {
		NormalizeAddresses(r, n.DefaultPort, &NodeConfig.Node.AddPeers)
	}
	if r, ok := getIfIs(ctx, "connectpeers"); ok {
		NormalizeAddresses(r, n.DefaultPort, &NodeConfig.Node.ConnectPeers)
	}
	if r, ok := getIfIs(ctx, "disablelisten"); ok {
		NodeConfig.Node.DisableListen = r == "true"
	}
	if r, ok := getIfIs(ctx, "listeners"); ok {
		NormalizeAddresses(r, n.DefaultPort, &NodeConfig.Node.Listeners)
	}
	if r, ok := getIfIs(ctx, "maxpeers"); ok {
		if err := ParseInteger(r, "maxpeers", &NodeConfig.Node.MaxPeers); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if r, ok := getIfIs(ctx, "disablebanning"); ok {
		NodeConfig.Node.DisableBanning = r == "true"
	}
	if r, ok := getIfIs(ctx, "banduration"); ok {
		if err := ParseDuration(r, "banduration", &NodeConfig.Node.BanDuration); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if r, ok := getIfIs(ctx, "banthreshold"); ok {
		var bt int
		if err := ParseInteger(r, "banthtreshold", &bt); err != nil {
			log <- cl.Wrn(err.Error())
		} else {
			NodeConfig.Node.BanThreshold = uint32(bt)
		}
	}
	if r, ok := getIfIs(ctx, "whitelists"); ok {
		NormalizeAddresses(r, n.DefaultPort, &NodeConfig.Node.Whitelists)
	}
	if r, ok := getIfIs(ctx, "rpcuser"); ok {
		NodeConfig.Node.RPCUser = r
	}
	if r, ok := getIfIs(ctx, "rpcpass"); ok {
		NodeConfig.Node.RPCPass = r
	}
	if r, ok := getIfIs(ctx, "rpclimituser"); ok {
		NodeConfig.Node.RPCLimitUser = r
	}
	if r, ok := getIfIs(ctx, "rpclimitpass"); ok {
		NodeConfig.Node.RPCLimitPass = r
	}
	if r, ok := getIfIs(ctx, "rpclisteners"); ok {
		NormalizeAddresses(r, n.DefaultRPCPort, &NodeConfig.Node.RPCListeners)
	}
	if r, ok := getIfIs(ctx, "rpccert"); ok {
		NodeConfig.Node.RPCCert = n.CleanAndExpandPath(r)
	}
	if r, ok := getIfIs(ctx, "rpckey"); ok {
		NodeConfig.Node.RPCKey = n.CleanAndExpandPath(r)
	}
	if r, ok := getIfIs(ctx, "tls"); ok {
		NodeConfig.Node.TLS = r == "true"
	}
	if r, ok := getIfIs(ctx, "disablednsseed"); ok {
		NodeConfig.Node.DisableDNSSeed = r == "true"
	}
	if r, ok := getIfIs(ctx, "externalips"); ok {
		NormalizeAddresses(r, n.DefaultPort, &NodeConfig.Node.ExternalIPs)
	}
	if r, ok := getIfIs(ctx, "proxy"); ok {
		NormalizeAddress(r, "9050", &NodeConfig.Node.Proxy)
	}
	if r, ok := getIfIs(ctx, "proxyuser"); ok {
		NodeConfig.Node.ProxyUser = r
	}
	if r, ok := getIfIs(ctx, "proxypass"); ok {
		NodeConfig.Node.ProxyPass = r
	}
	if r, ok := getIfIs(ctx, "onion"); ok {
		NormalizeAddress(r, "9050", &NodeConfig.Node.OnionProxy)
	}
	if r, ok := getIfIs(ctx, "onionuser"); ok {
		NodeConfig.Node.OnionProxyUser = r
	}
	if r, ok := getIfIs(ctx, "onionpass"); ok {
		NodeConfig.Node.OnionProxyPass = r
	}
	if r, ok := getIfIs(ctx, "noonion"); ok {
		NodeConfig.Node.NoOnion = r == "true"
	}
	if r, ok := getIfIs(ctx, "torisolation"); ok {
		NodeConfig.Node.TorIsolation = r == "true"
	}
	if r, ok := getIfIs(ctx, "network"); ok {
		switch r {
		case "testnet":
			NodeConfig.Node.TestNet3, NodeConfig.Node.RegressionTest, NodeConfig.Node.SimNet = true, false, false
		case "regtest":
			NodeConfig.Node.TestNet3, NodeConfig.Node.RegressionTest, NodeConfig.Node.SimNet = false, true, false
		case "simnet":
			NodeConfig.Node.TestNet3, NodeConfig.Node.RegressionTest, NodeConfig.Node.SimNet = false, false, true
		default:
			NodeConfig.Node.TestNet3, NodeConfig.Node.RegressionTest, NodeConfig.Node.SimNet = false, false, false
		}
	}
	if r, ok := getIfIs(ctx, "addcheckpoints"); ok {
		NodeConfig.Node.AddCheckpoints = strings.Split(r, " ")
	}
	if r, ok := getIfIs(ctx, "disablecheckpoints"); ok {
		NodeConfig.Node.DisableCheckpoints = r == "true"
	}
	if r, ok := getIfIs(ctx, "dbtype"); ok {
		NodeConfig.Node.DbType = r
	}
	if r, ok := getIfIs(ctx, "profile"); ok {
		var p int
		if err = ParseInteger(r, "profile", &p); err == nil {
			NodeConfig.Node.Profile = fmt.Sprint(p)
		}
	}
	if r, ok := getIfIs(ctx, "cpuprofile"); ok {
		NodeConfig.Node.CPUProfile = r
	}
	if r, ok := getIfIs(ctx, "upnp"); ok {
		NodeConfig.Node.Upnp = r == "true"
	}
	if r, ok := getIfIs(ctx, "minrelaytxfee"); ok {
		if err := ParseFloat(r, "minrelaytxfee", &NodeConfig.Node.MinRelayTxFee); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if r, ok := getIfIs(ctx, "freetxrelaylimit"); ok {
		if err := ParseFloat(r, "freetxrelaylimit", &NodeConfig.Node.FreeTxRelayLimit); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if r, ok := getIfIs(ctx, "norelaypriority"); ok {
		NodeConfig.Node.NoRelayPriority = r == "true"
	}
	if r, ok := getIfIs(ctx, "trickleinterval"); ok {
		if err := ParseDuration(r, "trickleinterval", &NodeConfig.Node.TrickleInterval); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if r, ok := getIfIs(ctx, "maxorphantxs"); ok {
		if err := ParseInteger(r, "maxorphantxs", &NodeConfig.Node.MaxOrphanTxs); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if r, ok := getIfIs(ctx, "algo"); ok {
		NodeConfig.Node.Algo = r
	}
	if r, ok := getIfIs(ctx, "generate"); ok {
		NodeConfig.Node.Generate = r == "true"
	}
	if r, ok := getIfIs(ctx, "genthreads"); ok {
		var gt int
		if err := ParseInteger(r, "genthreads", &gt); err != nil {
			log <- cl.Wrn(err.Error())
		} else {
			NodeConfig.Node.GenThreads = int32(gt)
		}
	}
	if r, ok := getIfIs(ctx, "miningaddrs"); ok {
		NodeConfig.Node.MiningAddrs = strings.Split(r, " ")
	}
	if r, ok := getIfIs(ctx, "minerlistener"); ok {
		NormalizeAddress(r, n.DefaultRPCPort, &NodeConfig.Node.MinerListener)
	}
	if r, ok := getIfIs(ctx, "minerpass"); ok {
		NodeConfig.Node.MinerPass = r
	}
	if r, ok := getIfIs(ctx, "blockminsize"); ok {
		if err := ParseUint32(r, "blockminsize", &NodeConfig.Node.BlockMinSize); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if r, ok := getIfIs(ctx, "blockmaxsize"); ok {
		if err := ParseUint32(r, "blockmaxsize", &NodeConfig.Node.BlockMaxSize); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if r, ok := getIfIs(ctx, "blockminweight"); ok {
		if err := ParseUint32(r, "blockminweight", &NodeConfig.Node.BlockMinWeight); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if r, ok := getIfIs(ctx, "blockmaxweight"); ok {
		if err := ParseUint32(r, "blockmaxweight", &NodeConfig.Node.BlockMaxWeight); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if r, ok := getIfIs(ctx, "blockprioritysize"); ok {
		if err := ParseUint32(r, "blockmaxweight", &NodeConfig.Node.BlockPrioritySize); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if r, ok := getIfIs(ctx, "uacomment"); ok {
		NodeConfig.Node.UserAgentComments = strings.Split(r, " ")
	}
	if r, ok := getIfIs(ctx, "nopeerbloomfilters"); ok {
		NodeConfig.Node.NoPeerBloomFilters = r == "true"
	}
	if r, ok := getIfIs(ctx, "nocfilters"); ok {
		NodeConfig.Node.NoCFilters = r == "true"
	}
	if ctx.Is("dropcfindex") {
		NodeConfig.Node.DropCfIndex = true
	}
	if r, ok := getIfIs(ctx, "sigcachemaxsize"); ok {
		var scms int
		if err := ParseInteger(r, "sigcachemaxsize", &scms); err != nil {
			log <- cl.Wrn(err.Error())
		} else {
			NodeConfig.Node.SigCacheMaxSize = uint(scms)
		}
	}
	if r, ok := getIfIs(ctx, "blocksonly"); ok {
		NodeConfig.Node.BlocksOnly = r == "true"
	}
	if r, ok := getIfIs(ctx, "txindex"); ok {
		NodeConfig.Node.TxIndex = r == "true"
	}
	if ctx.Is("droptxindex") {
		NodeConfig.Node.DropTxIndex = true
	}
	if ctx.Is("addrindex") {
		r, _ := ctx.Get("addrindex")
		NodeConfig.Node.AddrIndex = r == "true"
	}
	if ctx.Is("dropaddrindex") {
		NodeConfig.Node.DropAddrIndex = true
	}
	if r, ok := getIfIs(ctx, "relaynonstd"); ok {
		NodeConfig.Node.RelayNonStd = r == "true"
	}
	if r, ok := getIfIs(ctx, "rejectnonstd"); ok {
		NodeConfig.Node.RejectNonStd = r == "true"
	}
	SetLogging(ctx)
	if ctx.Is("save") {
		log <- cl.Infof{
			"saving config file to %s",
			cfgFile,
		}
		j, err := json.MarshalIndent(NodeConfig, "", "  ")
		if err != nil {
			log <- cl.Error{
				"saving config file:",
				err.Error(),
			}
		}
		j = append(j, '\n')
		log <- cl.Tracef{
			"JSON formatted config file\n%s",
			j,
		}
		err = ioutil.WriteFile(cfgFile, j, 0600)
		if err != nil {
			log <- cl.Error{"writing app config file:", err.Error()}
		}
	}
	// Service options which are only added on Windows.
	serviceOpts := serviceOptions{}
	// Perform service command and exit if specified.  Invalid service commands show an appropriate error.  Only runs on Windows since the runServiceCommand function will be nil when not on Windows.
	if serviceOpts.ServiceCommand != "" && runServiceCommand != nil {
		err := runServiceCommand(serviceOpts.ServiceCommand)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(0)
	}
	// Don't add peers from the config file when in regression test mode.
	if NodeConfig.Node.RegressionTest && len(NodeConfig.Node.AddPeers) > 0 {
		NodeConfig.Node.AddPeers = nil
	}
	// Set the mining algorithm correctly, default to random if unrecognised
	switch NodeConfig.Node.Algo {
	case "blake14lr", "cryptonight7v2", "keccak", "lyra2rev2", "scrypt", "skein", "x11", "stribog", "random", "easy":
	default:
		NodeConfig.Node.Algo = "random"
	}
	relayNonStd := n.ActiveNetParams.RelayNonStdTxs
	funcName := "loadConfig"
	switch {
	case NodeConfig.Node.RelayNonStd && NodeConfig.Node.RejectNonStd:
		str := "%s: rejectnonstd and relaynonstd cannot be used together -- choose only one"
		err := fmt.Errorf(str, funcName)
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	case NodeConfig.Node.RejectNonStd:
		relayNonStd = false
	case NodeConfig.Node.RelayNonStd:
		relayNonStd = true
	}
	NodeConfig.Node.RelayNonStd = relayNonStd
	// Append the network type to the data directory so it is "namespaced" per network.  In addition to the block database, there are other pieces of data that are saved to disk such as address manager state. All data is specific to a network, so namespacing the data directory means each individual piece of serialized data does not have to worry about changing names per network and such.
	NodeConfig.Node.DataDir = n.CleanAndExpandPath(NodeConfig.Node.DataDir)
	NodeConfig.Node.DataDir = filepath.Join(NodeConfig.Node.DataDir, n.NetName(n.ActiveNetParams))
	// Append the network type to the log directory so it is "namespaced" per network in the same fashion as the data directory.
	NodeConfig.Node.LogDir = n.CleanAndExpandPath(NodeConfig.Node.LogDir)
	NodeConfig.Node.LogDir = filepath.Join(NodeConfig.Node.LogDir, n.NetName(n.ActiveNetParams))

	// Initialize log rotation.  After log rotation has been initialized, the logger variables may be used.
	// initLogRotator(filepath.Join(NodeConfig.Node.LogDir, DefaultLogFilename))
	// Validate database type.
	if !n.ValidDbType(NodeConfig.Node.DbType) {
		str := "%s: The specified database type [%v] is invalid -- " +
			"supported types %v"
		err := fmt.Errorf(str, funcName, NodeConfig.Node.DbType, n.KnownDbTypes)
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// Validate profile port number
	if NodeConfig.Node.Profile != "" {
		profilePort, err := strconv.Atoi(NodeConfig.Node.Profile)
		if err != nil || profilePort < 1024 || profilePort > 65535 {
			str := "%s: The profile port must be between 1024 and 65535"
			err := fmt.Errorf(str, funcName)
			fmt.Fprintln(os.Stderr, err)
			fmt.Fprintln(os.Stderr, usageMessage)
			cl.Shutdown()
		}
	}
	// Don't allow ban durations that are too short.
	if NodeConfig.Node.BanDuration < time.Second {
		str := "%s: The banduration option may not be less than 1s -- parsed [%v]"
		err := fmt.Errorf(str, funcName, NodeConfig.Node.BanDuration)
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// Validate any given whitelisted IP addresses and networks.
	if len(NodeConfig.Node.Whitelists) > 0 {
		var ip net.IP
		StateCfg.ActiveWhitelists = make([]*net.IPNet, 0, len(NodeConfig.Node.Whitelists))
		for _, addr := range NodeConfig.Node.Whitelists {
			_, ipnet, err := net.ParseCIDR(addr)
			if err != nil {
				ip = net.ParseIP(addr)
				if ip == nil {
					str := "%s: The whitelist value of '%s' is invalid"
					err = fmt.Errorf(str, funcName, addr)
					log <- cl.Err(err.Error())
					fmt.Fprintln(os.Stderr, usageMessage)
					cl.Shutdown()
				}
				var bits int
				if ip.To4() == nil {
					// IPv6
					bits = 128
				} else {
					bits = 32
				}
				ipnet = &net.IPNet{
					IP:   ip,
					Mask: net.CIDRMask(bits, bits),
				}
			}
			StateCfg.ActiveWhitelists = append(StateCfg.ActiveWhitelists, ipnet)
		}
	}
	// --addPeer and --connect do not mix.
	if len(NodeConfig.Node.AddPeers) > 0 && len(NodeConfig.Node.ConnectPeers) > 0 {
		str := "%s: the --addpeer and --connect options can not be " +
			"mixed"
		err := fmt.Errorf(str, funcName)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
	}
	// --proxy or --connect without --listen disables listening.
	if (NodeConfig.Node.Proxy != "" || len(NodeConfig.Node.ConnectPeers) > 0) &&
		len(NodeConfig.Node.Listeners) == 0 {
		NodeConfig.Node.DisableListen = true
	}
	// Connect means no DNS seeding.
	if len(NodeConfig.Node.ConnectPeers) > 0 {
		NodeConfig.Node.DisableDNSSeed = true
	}
	// Add the default listener if none were specified. The default listener is all addresses on the listen port for the network we are to connect to.
	if len(NodeConfig.Node.Listeners) == 0 {
		NodeConfig.Node.Listeners = []string{
			net.JoinHostPort("", n.ActiveNetParams.DefaultPort),
		}
	}
	// Check to make sure limited and admin users don't have the same username
	if NodeConfig.Node.RPCUser == NodeConfig.Node.RPCLimitUser && NodeConfig.Node.RPCUser != "" {
		str := "%s: --rpcuser and --rpclimituser must not specify the same username"
		err := fmt.Errorf(str, funcName)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// Check to make sure limited and admin users don't have the same password
	if NodeConfig.Node.RPCPass == NodeConfig.Node.RPCLimitPass && NodeConfig.Node.RPCPass != "" {
		str := "%s: --rpcpass and --rpclimitpass must not specify the " +
			"same password"
		err := fmt.Errorf(str, funcName)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// The RPC server is disabled if no username or password is provided.
	if (NodeConfig.Node.RPCUser == "" || NodeConfig.Node.RPCPass == "") &&
		(NodeConfig.Node.RPCLimitUser == "" || NodeConfig.Node.RPCLimitPass == "") {
		NodeConfig.Node.DisableRPC = true
	}
	if NodeConfig.Node.DisableRPC {
		log <- cl.Inf("RPC service is disabled")
	}
	// Default RPC to listen on localhost only.
	if !NodeConfig.Node.DisableRPC && len(NodeConfig.Node.RPCListeners) == 0 {
		addrs, err := net.LookupHost(n.DefaultRPCListener)
		if err != nil {
			log <- cl.Err(err.Error())
			cl.Shutdown()
		}
		NodeConfig.Node.RPCListeners = make([]string, 0, len(addrs))
		for _, addr := range addrs {
			addr = net.JoinHostPort(addr, n.ActiveNetParams.RPCPort)
			NodeConfig.Node.RPCListeners = append(NodeConfig.Node.RPCListeners, addr)
		}
	}
	if NodeConfig.Node.RPCMaxConcurrentReqs < 0 {
		str := "%s: The rpcmaxwebsocketconcurrentrequests option may not be less than 0 -- parsed [%d]"
		err := fmt.Errorf(str, funcName, NodeConfig.Node.RPCMaxConcurrentReqs)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// Validate the the minrelaytxfee.
	StateCfg.ActiveMinRelayTxFee, err = util.NewAmount(NodeConfig.Node.MinRelayTxFee)
	if err != nil {
		str := "%s: invalid minrelaytxfee: %v"
		err := fmt.Errorf(str, funcName, err)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// Limit the max block size to a sane value.
	if NodeConfig.Node.BlockMaxSize < n.BlockMaxSizeMin || NodeConfig.Node.BlockMaxSize >
		n.BlockMaxSizeMax {
		str := "%s: The blockmaxsize option must be in between %d and %d -- parsed [%d]"
		err := fmt.Errorf(str, funcName, n.BlockMaxSizeMin,
			n.BlockMaxSizeMax, NodeConfig.Node.BlockMaxSize)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// Limit the max block weight to a sane value.
	if NodeConfig.Node.BlockMaxWeight < n.BlockMaxWeightMin ||
		NodeConfig.Node.BlockMaxWeight > n.BlockMaxWeightMax {
		str := "%s: The blockmaxweight option must be in between %d and %d -- parsed [%d]"
		err := fmt.Errorf(str, funcName, n.BlockMaxWeightMin,
			n.BlockMaxWeightMax, NodeConfig.Node.BlockMaxWeight)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// Limit the max orphan count to a sane vlue.
	if NodeConfig.Node.MaxOrphanTxs < 0 {
		str := "%s: The maxorphantx option may not be less than 0 -- parsed [%d]"
		err := fmt.Errorf(str, funcName, NodeConfig.Node.MaxOrphanTxs)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// Limit the block priority and minimum block sizes to max block size.
	NodeConfig.Node.BlockPrioritySize = minUint32(NodeConfig.Node.BlockPrioritySize, NodeConfig.Node.BlockMaxSize)
	NodeConfig.Node.BlockMinSize = minUint32(NodeConfig.Node.BlockMinSize, NodeConfig.Node.BlockMaxSize)
	NodeConfig.Node.BlockMinWeight = minUint32(NodeConfig.Node.BlockMinWeight, NodeConfig.Node.BlockMaxWeight)
	switch {
	// If the max block size isn't set, but the max weight is, then we'll set the limit for the max block size to a safe limit so weight takes precedence.
	case NodeConfig.Node.BlockMaxSize == n.DefaultBlockMaxSize &&
		NodeConfig.Node.BlockMaxWeight != n.DefaultBlockMaxWeight:
		NodeConfig.Node.BlockMaxSize = blockchain.MaxBlockBaseSize - 1000
	// If the max block weight isn't set, but the block size is, then we'll scale the set weight accordingly based on the max block size value.
	case NodeConfig.Node.BlockMaxSize != n.DefaultBlockMaxSize &&
		NodeConfig.Node.BlockMaxWeight == n.DefaultBlockMaxWeight:
		NodeConfig.Node.BlockMaxWeight = NodeConfig.Node.BlockMaxSize * blockchain.WitnessScaleFactor
	}
	// Look for illegal characters in the user agent comments.
	for _, uaComment := range NodeConfig.Node.UserAgentComments {
		if strings.ContainsAny(uaComment, "/:()") {
			err := fmt.Errorf("%s: The following characters must not "+
				"appear in user agent comments: '/', ':', '(', ')'",
				funcName)
			log <- cl.Err(err.Error())
			fmt.Fprintln(os.Stderr, usageMessage)
			cl.Shutdown()

		}
	}
	// --txindex and --droptxindex do not mix.
	if NodeConfig.Node.TxIndex && NodeConfig.Node.DropTxIndex {
		err := fmt.Errorf("%s: the --txindex and --droptxindex options may  not be activated at the same time",
			funcName)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()

	}
	// --addrindex and --dropaddrindex do not mix.
	if NodeConfig.Node.AddrIndex && NodeConfig.Node.DropAddrIndex {
		err := fmt.Errorf("%s: the --addrindex and --dropaddrindex "+
			"options may not be activated at the same time",
			funcName)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// --addrindex and --droptxindex do not mix.
	if NodeConfig.Node.AddrIndex && NodeConfig.Node.DropTxIndex {
		err := fmt.Errorf("%s: the --addrindex and --droptxindex options may not be activated at the same time "+
			"because the address index relies on the transaction index",
			funcName)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// Check mining addresses are valid and saved parsed versions.
	StateCfg.ActiveMiningAddrs = make([]util.Address, 0, len(NodeConfig.Node.MiningAddrs))
	for _, strAddr := range NodeConfig.Node.MiningAddrs {
		addr, err := util.DecodeAddress(strAddr, n.ActiveNetParams.Params)
		if err != nil {
			str := "%s: mining address '%s' failed to decode: %v"
			err := fmt.Errorf(str, funcName, strAddr, err)
			log <- cl.Err(err.Error())
			fmt.Fprintln(os.Stderr, usageMessage)
			cl.Shutdown()
		}
		if !addr.IsForNet(n.ActiveNetParams.Params) {
			str := "%s: mining address '%s' is on the wrong network"
			err := fmt.Errorf(str, funcName, strAddr)
			log <- cl.Err(err.Error())
			fmt.Fprintln(os.Stderr, usageMessage)
			cl.Shutdown()
		}
		StateCfg.ActiveMiningAddrs = append(StateCfg.ActiveMiningAddrs, addr)
	}
	// Ensure there is at least one mining address when the generate flag is set.
	if (NodeConfig.Node.Generate || NodeConfig.Node.MinerListener != "") && len(NodeConfig.Node.MiningAddrs) == 0 {
		str := "%s: the generate flag is set, but there are no mining addresses specified "
		err := fmt.Errorf(str, funcName)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()

	}
	if NodeConfig.Node.MinerPass != "" {
		StateCfg.ActiveMinerKey = fork.Argon2i([]byte(NodeConfig.Node.MinerPass))
	}
	// Add default port to all listener addresses if needed and remove duplicate addresses.
	NodeConfig.Node.Listeners = n.NormalizeAddresses(NodeConfig.Node.Listeners,
		n.ActiveNetParams.DefaultPort)
	// Add default port to all rpc listener addresses if needed and remove duplicate addresses.
	NodeConfig.Node.RPCListeners = n.NormalizeAddresses(NodeConfig.Node.RPCListeners,
		n.ActiveNetParams.RPCPort)
	if !NodeConfig.Node.DisableRPC && !NodeConfig.Node.TLS {
		for _, addr := range NodeConfig.Node.RPCListeners {
			if err != nil {
				str := "%s: RPC listen interface '%s' is invalid: %v"
				err := fmt.Errorf(str, funcName, addr, err)
				log <- cl.Err(err.Error())
				fmt.Fprintln(os.Stderr, usageMessage)
				cl.Shutdown()
			}
		}
	}
	// Add default port to all added peer addresses if needed and remove duplicate addresses.
	NodeConfig.Node.AddPeers = n.NormalizeAddresses(NodeConfig.Node.AddPeers,
		n.ActiveNetParams.DefaultPort)
	NodeConfig.Node.ConnectPeers = n.NormalizeAddresses(NodeConfig.Node.ConnectPeers,
		n.ActiveNetParams.DefaultPort)
	// --noonion and --onion do not mix.
	if NodeConfig.Node.NoOnion && NodeConfig.Node.OnionProxy != "" {
		err := fmt.Errorf("%s: the --noonion and --onion options may not be activated at the same time", funcName)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// Check the checkpoints for syntax errors.
	StateCfg.AddedCheckpoints, err = n.ParseCheckpoints(NodeConfig.Node.AddCheckpoints)
	if err != nil {
		str := "%s: Error parsing checkpoints: %v"
		err := fmt.Errorf(str, funcName, err)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// Tor stream isolation requires either proxy or onion proxy to be set.
	if NodeConfig.Node.TorIsolation && NodeConfig.Node.Proxy == "" && NodeConfig.Node.OnionProxy == "" {
		str := "%s: Tor stream isolation requires either proxy or onionproxy to be set"
		err := fmt.Errorf(str, funcName)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// Setup dial and DNS resolution (lookup) functions depending on the specified options.  The default is to use the standard net.DialTimeout function as well as the system DNS resolver.  When a proxy is specified, the dial function is set to the proxy specific dial function and the lookup is set to use tor (unless --noonion is specified in which case the system DNS resolver is used).
	StateCfg.Dial = net.DialTimeout
	StateCfg.Lookup = net.LookupIP
	if NodeConfig.Node.Proxy != "" {
		_, _, err := net.SplitHostPort(NodeConfig.Node.Proxy)
		if err != nil {
			str := "%s: Proxy address '%s' is invalid: %v"
			err := fmt.Errorf(str, funcName, NodeConfig.Node.Proxy, err)
			log <- cl.Err(err.Error())
			fmt.Fprintln(os.Stderr, usageMessage)
			cl.Shutdown()
		}
		// Tor isolation flag means proxy credentials will be overridden unless there is also an onion proxy configured in which case that one will be overridden.
		torIsolation := false
		if NodeConfig.Node.TorIsolation && NodeConfig.Node.OnionProxy == "" &&
			(NodeConfig.Node.ProxyUser != "" || NodeConfig.Node.ProxyPass != "") {
			torIsolation = true
			fmt.Fprintln(os.Stderr, "Tor isolation set -- "+
				"overriding specified proxy user credentials")
		}
		proxy := &socks.Proxy{
			Addr:         NodeConfig.Node.Proxy,
			Username:     NodeConfig.Node.ProxyUser,
			Password:     NodeConfig.Node.ProxyPass,
			TorIsolation: torIsolation,
		}
		StateCfg.Dial = proxy.DialTimeout
		// Treat the proxy as tor and perform DNS resolution through it unless the --noonion flag is set or there is an onion-specific proxy configured.
		if !NodeConfig.Node.NoOnion && NodeConfig.Node.OnionProxy == "" {
			StateCfg.Lookup = func(host string) ([]net.IP, error) {
				return connmgr.TorLookupIP(host, NodeConfig.Node.Proxy)
			}
		}
	}
	// Setup onion address dial function depending on the specified options. The default is to use the same dial function selected above.  However, when an onion-specific proxy is specified, the onion address dial function is set to use the onion-specific proxy while leaving the normal dial function as selected above.  This allows .onion address traffic to be routed through a different proxy than normal traffic.
	if NodeConfig.Node.OnionProxy != "" {
		_, _, err := net.SplitHostPort(NodeConfig.Node.OnionProxy)
		if err != nil {
			str := "%s: Onion proxy address '%s' is invalid: %v"
			err := fmt.Errorf(str, funcName, NodeConfig.Node.OnionProxy, err)
			log <- cl.Err(err.Error())
			fmt.Fprintln(os.Stderr, usageMessage)
			cl.Shutdown()
		}
		// Tor isolation flag means onion proxy credentials will be overridden.
		if NodeConfig.Node.TorIsolation &&
			(NodeConfig.Node.OnionProxyUser != "" || NodeConfig.Node.OnionProxyPass != "") {
			fmt.Fprintln(os.Stderr, "Tor isolation set -- "+
				"overriding specified onionproxy user "+
				"credentials ")
		}
		StateCfg.Oniondial = func(network, addr string, timeout time.Duration) (net.Conn, error) {
			proxy := &socks.Proxy{
				Addr:         NodeConfig.Node.OnionProxy,
				Username:     NodeConfig.Node.OnionProxyUser,
				Password:     NodeConfig.Node.OnionProxyPass,
				TorIsolation: NodeConfig.Node.TorIsolation,
			}
			return proxy.DialTimeout(network, addr, timeout)
		}
		// When configured in bridge mode (both --onion and --proxy are configured), it means that the proxy configured by --proxy is not a tor proxy, so override the DNS resolution to use the onion-specific proxy.
		if NodeConfig.Node.Proxy != "" {
			StateCfg.Lookup = func(host string) ([]net.IP, error) {
				return connmgr.TorLookupIP(host, NodeConfig.Node.OnionProxy)
			}
		}
	} else {
		StateCfg.Oniondial = StateCfg.Dial
	}
	// Specifying --noonion means the onion address dial function results in an error.
	if NodeConfig.Node.NoOnion {
		StateCfg.Oniondial = func(a, b string, t time.Duration) (net.Conn, error) {
			return nil, errors.New("tor has been disabled")
		}
	}
}

// WriteNodeConfig writes the current config to the requested location
func WriteNodeConfig(cfgFile string, c *NodeCfg) {
	log <- cl.Dbg("writing config")
	c.Node.ConfigFile = cfgFile
	j, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		log <- cl.Error{`marshalling default app config file: "`, err, `"`}
		log <- cl.Err(spew.Sdump(c))
		return
	}
	j = append(j, '\n')
	log <- cl.Tracef{
		"JSON formatted config file\n%s",
		j,
	}
	err = ioutil.WriteFile(cfgFile, j, 0600)
	if err != nil {
		log <- cl.Error{"writing default app config file:", err.Error()}
		return
	}
}

// WriteDefaultNodeConfig creates a default config and writes it to the requested location
func WriteDefaultNodeConfig(cfgFile string) {
	log <- cl.Dbg("writing default config")
	defCfg := DefaultNodeConfig()
	defCfg.Node.ConfigFile = cfgFile
	j, err := json.MarshalIndent(defCfg, "", "  ")
	if err != nil {
		log <- cl.Error{`marshalling default app config file: "`, err, `"`}
		log <- cl.Err(spew.Sdump(defCfg))
		return
	}
	j = append(j, '\n')
	log <- cl.Tracef{
		"JSON formatted config file\n%s",
		j,
	}
	err = ioutil.WriteFile(cfgFile, j, 0600)
	if err != nil {
		log <- cl.Error{"writing default app config file:", err.Error()}
		return
	}
	// if we are writing default config we also want to use it
	NodeConfig = defCfg
}

// DefaultNodeConfig is the default configuration for node
func DefaultNodeConfig() *NodeCfg {
	user := GenKey()
	pass := GenKey()
	return &NodeCfg{
		Node: &n.Config{
			RPCUser:              user,
			RPCPass:              pass,
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
		LogLevels: GetDefaultLogLevelsConfig(),
	}
}
