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

	"git.parallelcoin.io/pod/pkg/blockchain"
	cl "git.parallelcoin.io/pod/pkg/clog"
	"git.parallelcoin.io/pod/pkg/connmgr"
	"git.parallelcoin.io/pod/pkg/fork"
	"git.parallelcoin.io/pod/pkg/util"
	n "git.parallelcoin.io/pod/cmd/node"
	"git.parallelcoin.io/pod/cmd/node/mempool"
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
var f = pu.GenFlag
var t = pu.GenTrig
var s = pu.GenShort
var l = pu.GenLog

// NodeCommand is a command to send RPC queries to bitcoin RPC protocol server for node and wallet queries
var NodeCommand = climax.Command{
	Name:  "node",
	Brief: "parallelcoin full node",
	Help:  "distrubutes, verifies and mines blocks for the parallelcoin duo cryptocurrency, as well as optionally providing search indexes for transactions in the database",
	Flags: []climax.Flag{

		t("version", "V", "show version number and quit"),

		s("configfile", "C", "path to configuration file"),
		s("datadir", "D", "path to configuration file"),

		t("init", "", "resets configuration to defaults"),
		t("save", "", "saves current configuration"),

		f("network", "connect to (mainnet|testnet|simnet"),

		f("txindex", "enable transaction index"),
		f("addrindex", "enable address index"),
		t("dropcfindex", "", "delete committed filtering (CF) index then exit"),
		t("droptxindex", "", "deletes transaction index then exit"),
		t("dropaddrindex", "", "deletes the address index then exits"),

		s("listeners", "S", "sets an address to listen for P2P connections"),
		f("externalips", "additional P2P listeners"),
		f("disablelisten", "disables the P2P listener"),

		f("addpeers", "adds a peer to the peers database to try to connect to"),
		f("connectpeers", "adds a peer to a connect-only whitelist"),
		f(`maxpeers`, "sets max number of peers to connect to to at once"),
		f(`disablebanning`, "disable banning of misbehaving peers"),
		f("banduration", "time to ban misbehaving peers (d/h/m/s)"),
		f("banthreshold", "banscore that triggers a ban"),
		f("whitelists", "addresses and networks immune to banning"),

		s("rpcuser", "u", "RPC username"),
		s("rpcpass", "P", "RPC password"),

		f("rpclimituser", "limited user RPC username"),
		f("rpclimitpass", "limited user RPC password"),

		s("rpclisteners", "s", "RPC server to connect to"),

		f("rpccert", "RPC server tls certificate chain for validation"),
		f("rpckey", "RPC server tls key for authentication"),
		f("tls", "enable TLS"),
		f("skipverify", "do not verify tls certificates"),

		f("proxy", "connect via SOCKS5 proxy server"),
		f("proxyuser", "username for proxy server"),
		f("proxypass", "password for proxy server"),

		f("onion", "connect via tor proxy relay"),
		f("onionuser", "username for onion proxy server"),
		f("onionpass", "password for onion proxy server"),
		f("noonion", "disable onion proxy"),
		f("torisolation", "use a different user/pass for each peer"),

		f("trickleinterval", "time between sending inventory batches to peers"),
		f("minrelaytxfee", "min fee in DUO/kb to relay transaction"),
		f("freetxrelaylimit", "limit below min fee transactions in kb/bin"),
		f("norelaypriority", "do not discriminate transactions for relaying"),

		f("nopeerbloomfilters", "disable bloom filtering support"),
		f("nocfilters", "disable committed filtering (CF) support"),
		f("blocksonly", "do not accept transactions from peers"),
		f("relaynonstd", "relay nonstandard transactions"),
		f("rejectnonstd", "reject nonstandard transactions"),

		f("maxorphantxs", "max number of orphan transactions to store"),
		f("sigcachemaxsize", "maximum number of signatures to store in memory"),

		f("generate", "set CPU miner to generate blocks"),
		f("genthreads", "set number of threads to generate using CPU, -1 = all"),
		f("algo", "set algorithm to be used by cpu miner"),
		f("miningaddrs", "add address to pay block rewards to"),
		f("minerlistener", "address to listen for mining work subscriptions"),
		f("minerpass", "PSK to prevent snooping/spoofing of miner traffic"),

		f("addcheckpoints", `add custom checkpoints "height:hash"`),
		f("disablecheckpoints", "disable all checkpoints"),

		f("blockminsize", "min block size for miners"),
		f("blockmaxsize", "max block size for miners"),
		f("blockminweight", "min block weight for miners"),
		f("blockmaxweight", "max block weight for miners"),
		f("blockprioritysize", "size in bytes of high priority blocks"),

		f("uacomment", "comment to add to the P2P network user agent"),
		f("upnp", "use UPNP to automatically port forward to node"),
		f("dbtype", "set database backend type (ffldb)"),
		f("disablednsseed", "disable dns seeding"),

		f("profile", "start HTTP profiling server on given address"),
		f("cpuprofile", "start cpu profiling server on given address"),

		s("debuglevel", "d", "sets log level for those unspecified below"),

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
			Config.Node.DebugLevel = dl
			Log.SetLevel(dl)
			for i := range logger.Levels {
				logger.Levels[i].SetLevel(dl)
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
				err = json.Unmarshal(cfgData, &Config)
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

<<<<<<< HEAD
func getIfIs(ctx *climax.Context, name string, r *string) (ok bool) {
	if ctx.Is(name) {
		var s string
		s, ok = ctx.Get(name)
		r = &s
	}
	return
}

=======
>>>>>>> master
func configNode(ctx *climax.Context, cfgFile string) {
	Nodecfg := Config.Node
	var err error
	var r *string
	t := ""
	r = &t
	if getIfIs(ctx, "debuglevel", r) {
		switch *r {
		case "fatal", "error", "warn", "info", "debug", "trace":
			Nodecfg.DebugLevel = *r
		default:
			Nodecfg.DebugLevel = "info"
		}
		Log.SetLevel(Nodecfg.DebugLevel)
	}
	if getIfIs(ctx, "datadir", r) {
		Nodecfg.DataDir = n.CleanAndExpandPath(*r)
	}
	if getIfIs(ctx, "addpeers", r) {
		pu.NormalizeAddresses(*r, n.DefaultPort, &Nodecfg.AddPeers)
	}
	if getIfIs(ctx, "connectpeers", r) {
		pu.NormalizeAddresses(*r, n.DefaultPort, &Nodecfg.ConnectPeers)
	}
	if getIfIs(ctx, "disablelisten", r) {
		Nodecfg.DisableListen = *r == "true"
	}
	if getIfIs(ctx, "listeners", r) {
		pu.NormalizeAddresses(*r, n.DefaultPort, &Nodecfg.Listeners)
	}
	if getIfIs(ctx, "maxpeers", r) {
		if err := pu.ParseInteger(*r, "maxpeers", &Nodecfg.MaxPeers); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if getIfIs(ctx, "disablebanning", r) {
		Nodecfg.DisableBanning = *r == "true"
	}
	if getIfIs(ctx, "banduration", r) {
		if err := pu.ParseDuration(*r, "banduration", &Nodecfg.BanDuration); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if getIfIs(ctx, "banthreshold", r) {
		var bt int
		if err := pu.ParseInteger(*r, "banthtreshold", &bt); err != nil {
			log <- cl.Wrn(err.Error())
		} else {
			Nodecfg.BanThreshold = uint32(bt)
		}
	}
	if getIfIs(ctx, "whitelists", r) {
		pu.NormalizeAddresses(*r, n.DefaultPort, &Nodecfg.Whitelists)
	}
	if getIfIs(ctx, "rpcuser", r) {
		Nodecfg.RPCUser = *r
	}
	if getIfIs(ctx, "rpcpass", r) {
		Nodecfg.RPCPass = *r
	}
	if getIfIs(ctx, "rpclimituser", r) {
		Nodecfg.RPCLimitUser = *r
	}
	if getIfIs(ctx, "rpclimitpass", r) {
		Nodecfg.RPCLimitPass = *r
	}
	if getIfIs(ctx, "rpclisteners", r) {
		pu.NormalizeAddresses(*r, n.DefaultRPCPort, &Nodecfg.RPCListeners)
	}
	if getIfIs(ctx, "rpccert", r) {
		Nodecfg.RPCCert = n.CleanAndExpandPath(*r)
	}
	if getIfIs(ctx, "rpckey", r) {
		Nodecfg.RPCKey = n.CleanAndExpandPath(*r)
	}
	if getIfIs(ctx, "tls", r) {
		Nodecfg.TLS = *r == "true"
	}
	if getIfIs(ctx, "disablednsseed", r) {
		Nodecfg.DisableDNSSeed = *r == "true"
	}
	if getIfIs(ctx, "externalips", r) {
		pu.NormalizeAddresses(*r, n.DefaultPort, &Nodecfg.ExternalIPs)
	}
	if getIfIs(ctx, "proxy", r) {
		pu.NormalizeAddress(*r, "9050", &Nodecfg.Proxy)
	}
	if getIfIs(ctx, "proxyuser", r) {
		Nodecfg.ProxyUser = *r
	}
	if getIfIs(ctx, "proxypass", r) {
		Nodecfg.ProxyPass = *r
	}
	if getIfIs(ctx, "onion", r) {
		pu.NormalizeAddress(*r, "9050", &Nodecfg.OnionProxy)
	}
	if getIfIs(ctx, "onionuser", r) {
		Nodecfg.OnionProxyUser = *r
	}
	if getIfIs(ctx, "onionpass", r) {
		Nodecfg.OnionProxyPass = *r
	}
	if getIfIs(ctx, "noonion", r) {
		Nodecfg.NoOnion = *r == "true"
	}
	if getIfIs(ctx, "torisolation", r) {
		Nodecfg.TorIsolation = *r == "true"
	}
	if getIfIs(ctx, "network", r) {
		switch *r {
		case "testnet":
			Nodecfg.TestNet3, Nodecfg.RegressionTest, Nodecfg.SimNet = true, false, false
		case "regtest":
			Nodecfg.TestNet3, Nodecfg.RegressionTest, Nodecfg.SimNet = false, true, false
		case "simnet":
			Nodecfg.TestNet3, Nodecfg.RegressionTest, Nodecfg.SimNet = false, false, true
		default:
			Nodecfg.TestNet3, Nodecfg.RegressionTest, Nodecfg.SimNet = false, false, false
		}
	}
	if getIfIs(ctx, "addcheckpoints", r) {
		Nodecfg.AddCheckpoints = strings.Split(*r, " ")
	}
	if getIfIs(ctx, "disablecheckpoints", r) {
		Nodecfg.DisableCheckpoints = *r == "true"
	}
	if getIfIs(ctx, "dbtype", r) {
		Nodecfg.DbType = *r
	}
	if getIfIs(ctx, "profile", r) {
		var p int
		if err = pu.ParseInteger(*r, "profile", &p); err == nil {
			Nodecfg.Profile = fmt.Sprint(p)
		}
	}
	if getIfIs(ctx, "cpuprofile", r) {
		Nodecfg.CPUProfile = *r
	}
	if getIfIs(ctx, "upnp", r) {
		Nodecfg.Upnp = *r == "true"
	}
	if getIfIs(ctx, "minrelaytxfee", r) {
		if err := pu.ParseFloat(*r, "minrelaytxfee", &Nodecfg.MinRelayTxFee); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if getIfIs(ctx, "freetxrelaylimit", r) {
		if err := pu.ParseFloat(*r, "freetxrelaylimit", &Nodecfg.FreeTxRelayLimit); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if getIfIs(ctx, "norelaypriority", r) {
		Nodecfg.NoRelayPriority = *r == "true"
	}
	if getIfIs(ctx, "trickleinterval", r) {
		if err := pu.ParseDuration(*r, "trickleinterval", &Nodecfg.TrickleInterval); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if getIfIs(ctx, "maxorphantxs", r) {
		if err := pu.ParseInteger(*r, "maxorphantxs", &Nodecfg.MaxOrphanTxs); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if getIfIs(ctx, "algo", r) {
		Nodecfg.Algo = *r
	}
	if getIfIs(ctx, "generate", r) {
		Nodecfg.Generate = *r == "true"
	}
	if getIfIs(ctx, "genthreads", r) {
		var gt int
		if err := pu.ParseInteger(*r, "genthreads", &gt); err != nil {
			log <- cl.Wrn(err.Error())
		} else {
			Nodecfg.GenThreads = int32(gt)
		}
	}
	if getIfIs(ctx, "miningaddrs", r) {
		Nodecfg.MiningAddrs = strings.Split(*r, " ")
	}
	if getIfIs(ctx, "minerlistener", r) {
		pu.NormalizeAddress(*r, n.DefaultRPCPort, &Nodecfg.MinerListener)
	}
	if getIfIs(ctx, "minerpass", r) {
		Nodecfg.MinerPass = *r
	}
	if getIfIs(ctx, "blockminsize", r) {
		if err := pu.ParseUint32(*r, "blockminsize", &Nodecfg.BlockMinSize); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if getIfIs(ctx, "blockmaxsize", r) {
		if err := pu.ParseUint32(*r, "blockmaxsize", &Nodecfg.BlockMaxSize); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if getIfIs(ctx, "blockminweight", r) {
		if err := pu.ParseUint32(*r, "blockminweight", &Nodecfg.BlockMinWeight); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if getIfIs(ctx, "blockmaxweight", r) {
		if err := pu.ParseUint32(*r, "blockmaxweight", &Nodecfg.BlockMaxWeight); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if getIfIs(ctx, "blockprioritysize", r) {
		if err := pu.ParseUint32(*r, "blockmaxweight", &Nodecfg.BlockPrioritySize); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if getIfIs(ctx, "uacomment", r) {
		Nodecfg.UserAgentComments = strings.Split(*r, " ")
	}
	if getIfIs(ctx, "nopeerbloomfilters", r) {
		Nodecfg.NoPeerBloomFilters = *r == "true"
	}
	if getIfIs(ctx, "nocfilters", r) {
		Nodecfg.NoCFilters = *r == "true"
	}
	if ctx.Is("dropcfindex") {
		Nodecfg.DropCfIndex = true
	}
	if getIfIs(ctx, "sigcachemaxsize", r) {
		var scms int
		if err := pu.ParseInteger(*r, "sigcachemaxsize", &scms); err != nil {
			log <- cl.Wrn(err.Error())
		} else {
			Nodecfg.SigCacheMaxSize = uint(scms)
		}
	}
	if getIfIs(ctx, "blocksonly", r) {
		Nodecfg.BlocksOnly = *r == "true"
	}
	if getIfIs(ctx, "txindex", r) {
		Nodecfg.TxIndex = *r == "true"
	}
	if ctx.Is("droptxindex") {
		Nodecfg.DropTxIndex = true
	}
	if ctx.Is("addrindex") {
		r, _ := ctx.Get("addrindex")
		Nodecfg.AddrIndex = r == "true"
	}
	if ctx.Is("dropaddrindex") {
		Nodecfg.DropAddrIndex = true
	}
	if getIfIs(ctx, "relaynonstd", r) {
		Nodecfg.RelayNonStd = *r == "true"
	}
	if getIfIs(ctx, "rejectnonstd", r) {
		Nodecfg.RejectNonStd = *r == "true"
	}
	logger.SetLogging(ctx)
	if ctx.Is("save") {
		log <- cl.Infof{
			"saving config file to %s",
			cfgFile,
		}
		j, err := json.MarshalIndent(Config, "", "  ")
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
	if Nodecfg.RegressionTest && len(Nodecfg.AddPeers) > 0 {
		Nodecfg.AddPeers = nil
	}
	// Set the mining algorithm correctly, default to random if unrecognised
	switch Nodecfg.Algo {
	case "blake14lr", "cryptonight7v2", "keccak", "lyra2rev2", "scrypt", "skein", "x11", "stribog", "random", "easy":
	default:
		Nodecfg.Algo = "random"
	}
	relayNonStd := ActiveNetParams.RelayNonStdTxs
	funcName := "loadConfig"
	switch {
	case Nodecfg.RelayNonStd && Nodecfg.RejectNonStd:
		str := "%s: rejectnonstd and relaynonstd cannot be used together -- choose only one"
		err := fmt.Errorf(str, funcName)
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	case Nodecfg.RejectNonStd:
		relayNonStd = false
	case Nodecfg.RelayNonStd:
		relayNonStd = true
	}
	Nodecfg.RelayNonStd = relayNonStd
	// Append the network type to the data directory so it is "namespaced" per network.  In addition to the block database, there are other pieces of data that are saved to disk such as address manager state. All data is specific to a network, so namespacing the data directory means each individual piece of serialized data does not have to worry about changing names per network and such.
	Nodecfg.DataDir = n.CleanAndExpandPath(Nodecfg.DataDir)
	Nodecfg.DataDir = filepath.Join(Nodecfg.DataDir, netName(ActiveNetParams))
	// Append the network type to the log directory so it is "namespaced" per network in the same fashion as the data directory.
	Nodecfg.LogDir = n.CleanAndExpandPath(Nodecfg.LogDir)
	Nodecfg.LogDir = filepath.Join(Nodecfg.LogDir, netName(ActiveNetParams))

	// Initialize log rotation.  After log rotation has been initialized, the logger variables may be used.
	// initLogRotator(filepath.Join(Nodecfg.LogDir, DefaultLogFilename))
	// Validate database type.
	if !n.ValidDbType(Nodecfg.DbType) {
		str := "%s: The specified database type [%v] is invalid -- " +
			"supported types %v"
		err := fmt.Errorf(str, funcName, Nodecfg.DbType, n.KnownDbTypes)
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// Validate profile port number
	if Nodecfg.Profile != "" {
		profilePort, err := strconv.Atoi(Nodecfg.Profile)
		if err != nil || profilePort < 1024 || profilePort > 65535 {
			str := "%s: The profile port must be between 1024 and 65535"
			err := fmt.Errorf(str, funcName)
			fmt.Fprintln(os.Stderr, err)
			fmt.Fprintln(os.Stderr, usageMessage)
			cl.Shutdown()
		}
	}
	// Don't allow ban durations that are too short.
	if Nodecfg.BanDuration < time.Second {
		str := "%s: The banduration option may not be less than 1s -- parsed [%v]"
		err := fmt.Errorf(str, funcName, Nodecfg.BanDuration)
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// Validate any given whitelisted IP addresses and networks.
	if len(Nodecfg.Whitelists) > 0 {
		var ip net.IP
		StateCfg.ActiveWhitelists = make([]*net.IPNet, 0, len(Nodecfg.Whitelists))
		for _, addr := range Nodecfg.Whitelists {
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
	if len(Nodecfg.AddPeers) > 0 && len(Nodecfg.ConnectPeers) > 0 {
		str := "%s: the --addpeer and --connect options can not be " +
			"mixed"
		err := fmt.Errorf(str, funcName)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
	}
	// --proxy or --connect without --listen disables listening.
	if (Nodecfg.Proxy != "" || len(Nodecfg.ConnectPeers) > 0) &&
		len(Nodecfg.Listeners) == 0 {
		Nodecfg.DisableListen = true
	}
	// Connect means no DNS seeding.
	if len(Nodecfg.ConnectPeers) > 0 {
		Nodecfg.DisableDNSSeed = true
	}
	// Add the default listener if none were specified. The default listener is all addresses on the listen port for the network we are to connect to.
	if len(Nodecfg.Listeners) == 0 {
		Nodecfg.Listeners = []string{
			net.JoinHostPort("", ActiveNetParams.DefaultPort),
		}
	}
	// Check to make sure limited and admin users don't have the same username
	if Nodecfg.RPCUser == Nodecfg.RPCLimitUser && Nodecfg.RPCUser != "" {
		str := "%s: --rpcuser and --rpclimituser must not specify the same username"
		err := fmt.Errorf(str, funcName)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// Check to make sure limited and admin users don't have the same password
	if Nodecfg.RPCPass == Nodecfg.RPCLimitPass && Nodecfg.RPCPass != "" {
		str := "%s: --rpcpass and --rpclimitpass must not specify the " +
			"same password"
		err := fmt.Errorf(str, funcName)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// The RPC server is disabled if no username or password is provided.
	if (Nodecfg.RPCUser == "" || Nodecfg.RPCPass == "") &&
		(Nodecfg.RPCLimitUser == "" || Nodecfg.RPCLimitPass == "") {
		Nodecfg.DisableRPC = true
	}
	if Nodecfg.DisableRPC {
		log <- cl.Inf("RPC service is disabled")
	}
	// Default RPC to listen on localhost only.
	if !Nodecfg.DisableRPC && len(Nodecfg.RPCListeners) == 0 {
		addrs, err := net.LookupHost(n.DefaultRPCListener)
		if err != nil {
			log <- cl.Err(err.Error())
			cl.Shutdown()
		}
		Nodecfg.RPCListeners = make([]string, 0, len(addrs))
		for _, addr := range addrs {
			addr = net.JoinHostPort(addr, ActiveNetParams.RPCPort)
			Nodecfg.RPCListeners = append(Nodecfg.RPCListeners, addr)
		}
	}
	if Nodecfg.RPCMaxConcurrentReqs < 0 {
		str := "%s: The rpcmaxwebsocketconcurrentrequests option may not be less than 0 -- parsed [%d]"
		err := fmt.Errorf(str, funcName, Nodecfg.RPCMaxConcurrentReqs)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// Validate the the minrelaytxfee.
	StateCfg.ActiveMinRelayTxFee, err = util.NewAmount(Nodecfg.MinRelayTxFee)
	if err != nil {
		str := "%s: invalid minrelaytxfee: %v"
		err := fmt.Errorf(str, funcName, err)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// Limit the max block size to a sane value.
	if Nodecfg.BlockMaxSize < n.BlockMaxSizeMin || Nodecfg.BlockMaxSize >
		n.BlockMaxSizeMax {
		str := "%s: The blockmaxsize option must be in between %d and %d -- parsed [%d]"
		err := fmt.Errorf(str, funcName, n.BlockMaxSizeMin,
			n.BlockMaxSizeMax, Nodecfg.BlockMaxSize)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// Limit the max block weight to a sane value.
	if Nodecfg.BlockMaxWeight < n.BlockMaxWeightMin ||
		Nodecfg.BlockMaxWeight > n.BlockMaxWeightMax {
		str := "%s: The blockmaxweight option must be in between %d and %d -- parsed [%d]"
		err := fmt.Errorf(str, funcName, n.BlockMaxWeightMin,
			n.BlockMaxWeightMax, Nodecfg.BlockMaxWeight)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// Limit the max orphan count to a sane vlue.
	if Nodecfg.MaxOrphanTxs < 0 {
		str := "%s: The maxorphantx option may not be less than 0 -- parsed [%d]"
		err := fmt.Errorf(str, funcName, Nodecfg.MaxOrphanTxs)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// Limit the block priority and minimum block sizes to max block size.
	Nodecfg.BlockPrioritySize = minUint32(Nodecfg.BlockPrioritySize, Nodecfg.BlockMaxSize)
	Nodecfg.BlockMinSize = minUint32(Nodecfg.BlockMinSize, Nodecfg.BlockMaxSize)
	Nodecfg.BlockMinWeight = minUint32(Nodecfg.BlockMinWeight, Nodecfg.BlockMaxWeight)
	switch {
	// If the max block size isn't set, but the max weight is, then we'll set the limit for the max block size to a safe limit so weight takes precedence.
	case Nodecfg.BlockMaxSize == n.DefaultBlockMaxSize &&
		Nodecfg.BlockMaxWeight != n.DefaultBlockMaxWeight:
		Nodecfg.BlockMaxSize = blockchain.MaxBlockBaseSize - 1000
	// If the max block weight isn't set, but the block size is, then we'll scale the set weight accordingly based on the max block size value.
	case Nodecfg.BlockMaxSize != n.DefaultBlockMaxSize &&
		Nodecfg.BlockMaxWeight == n.DefaultBlockMaxWeight:
		Nodecfg.BlockMaxWeight = Nodecfg.BlockMaxSize * blockchain.WitnessScaleFactor
	}
	// Look for illegal characters in the user agent comments.
	for _, uaComment := range Nodecfg.UserAgentComments {
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
	if Nodecfg.TxIndex && Nodecfg.DropTxIndex {
		err := fmt.Errorf("%s: the --txindex and --droptxindex options may  not be activated at the same time",
			funcName)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()

	}
	// --addrindex and --dropaddrindex do not mix.
	if Nodecfg.AddrIndex && Nodecfg.DropAddrIndex {
		err := fmt.Errorf("%s: the --addrindex and --dropaddrindex "+
			"options may not be activated at the same time",
			funcName)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// --addrindex and --droptxindex do not mix.
	if Nodecfg.AddrIndex && Nodecfg.DropTxIndex {
		err := fmt.Errorf("%s: the --addrindex and --droptxindex options may not be activated at the same time "+
			"because the address index relies on the transaction index",
			funcName)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// Check mining addresses are valid and saved parsed versions.
	StateCfg.ActiveMiningAddrs = make([]util.Address, 0, len(Nodecfg.MiningAddrs))
	for _, strAddr := range Nodecfg.MiningAddrs {
		addr, err := util.DecodeAddress(strAddr, ActiveNetParams.Params)
		if err != nil {
			str := "%s: mining address '%s' failed to decode: %v"
			err := fmt.Errorf(str, funcName, strAddr, err)
			log <- cl.Err(err.Error())
			fmt.Fprintln(os.Stderr, usageMessage)
			cl.Shutdown()
		}
		if !addr.IsForNet(ActiveNetParams.Params) {
			str := "%s: mining address '%s' is on the wrong network"
			err := fmt.Errorf(str, funcName, strAddr)
			log <- cl.Err(err.Error())
			fmt.Fprintln(os.Stderr, usageMessage)
			cl.Shutdown()
		}
		StateCfg.ActiveMiningAddrs = append(StateCfg.ActiveMiningAddrs, addr)
	}
	// Ensure there is at least one mining address when the generate flag is set.
	if (Nodecfg.Generate || Nodecfg.MinerListener != "") && len(Nodecfg.MiningAddrs) == 0 {
		str := "%s: the generate flag is set, but there are no mining addresses specified "
		err := fmt.Errorf(str, funcName)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()

	}
	if Nodecfg.MinerPass != "" {
		StateCfg.ActiveMinerKey = fork.Argon2i([]byte(Nodecfg.MinerPass))
	}
	// Add default port to all listener addresses if needed and remove duplicate addresses.
	Nodecfg.Listeners = n.NormalizeAddresses(Nodecfg.Listeners,
		ActiveNetParams.DefaultPort)
	// Add default port to all rpc listener addresses if needed and remove duplicate addresses.
	Nodecfg.RPCListeners = n.NormalizeAddresses(Nodecfg.RPCListeners,
		ActiveNetParams.RPCPort)
	if !Nodecfg.DisableRPC && !Nodecfg.TLS {
		for _, addr := range Nodecfg.RPCListeners {
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
	Nodecfg.AddPeers = n.NormalizeAddresses(Nodecfg.AddPeers,
		ActiveNetParams.DefaultPort)
	Nodecfg.ConnectPeers = n.NormalizeAddresses(Nodecfg.ConnectPeers,
		ActiveNetParams.DefaultPort)
	// --noonion and --onion do not mix.
	if Nodecfg.NoOnion && Nodecfg.OnionProxy != "" {
		err := fmt.Errorf("%s: the --noonion and --onion options may not be activated at the same time", funcName)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// Check the checkpoints for syntax errors.
	StateCfg.AddedCheckpoints, err = n.ParseCheckpoints(Nodecfg.AddCheckpoints)
	if err != nil {
		str := "%s: Error parsing checkpoints: %v"
		err := fmt.Errorf(str, funcName, err)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// Tor stream isolation requires either proxy or onion proxy to be set.
	if Nodecfg.TorIsolation && Nodecfg.Proxy == "" && Nodecfg.OnionProxy == "" {
		str := "%s: Tor stream isolation requires either proxy or onionproxy to be set"
		err := fmt.Errorf(str, funcName)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// Setup dial and DNS resolution (lookup) functions depending on the specified options.  The default is to use the standard net.DialTimeout function as well as the system DNS resolver.  When a proxy is specified, the dial function is set to the proxy specific dial function and the lookup is set to use tor (unless --noonion is specified in which case the system DNS resolver is used).
	StateCfg.Dial = net.DialTimeout
	StateCfg.Lookup = net.LookupIP
	if Nodecfg.Proxy != "" {
		_, _, err := net.SplitHostPort(Nodecfg.Proxy)
		if err != nil {
			str := "%s: Proxy address '%s' is invalid: %v"
			err := fmt.Errorf(str, funcName, Nodecfg.Proxy, err)
			log <- cl.Err(err.Error())
			fmt.Fprintln(os.Stderr, usageMessage)
			cl.Shutdown()
		}
		// Tor isolation flag means proxy credentials will be overridden unless there is also an onion proxy configured in which case that one will be overridden.
		torIsolation := false
		if Nodecfg.TorIsolation && Nodecfg.OnionProxy == "" &&
			(Nodecfg.ProxyUser != "" || Nodecfg.ProxyPass != "") {
			torIsolation = true
			fmt.Fprintln(os.Stderr, "Tor isolation set -- "+
				"overriding specified proxy user credentials")
		}
		proxy := &socks.Proxy{
			Addr:         Nodecfg.Proxy,
			Username:     Nodecfg.ProxyUser,
			Password:     Nodecfg.ProxyPass,
			TorIsolation: torIsolation,
		}
		StateCfg.Dial = proxy.DialTimeout
		// Treat the proxy as tor and perform DNS resolution through it unless the --noonion flag is set or there is an onion-specific proxy configured.
		if !Nodecfg.NoOnion && Nodecfg.OnionProxy == "" {
			StateCfg.Lookup = func(host string) ([]net.IP, error) {
				return connmgr.TorLookupIP(host, Nodecfg.Proxy)
			}
		}
	}
	// Setup onion address dial function depending on the specified options. The default is to use the same dial function selected above.  However, when an onion-specific proxy is specified, the onion address dial function is set to use the onion-specific proxy while leaving the normal dial function as selected above.  This allows .onion address traffic to be routed through a different proxy than normal traffic.
	if Nodecfg.OnionProxy != "" {
		_, _, err := net.SplitHostPort(Nodecfg.OnionProxy)
		if err != nil {
			str := "%s: Onion proxy address '%s' is invalid: %v"
			err := fmt.Errorf(str, funcName, Nodecfg.OnionProxy, err)
			log <- cl.Err(err.Error())
			fmt.Fprintln(os.Stderr, usageMessage)
			cl.Shutdown()
		}
		// Tor isolation flag means onion proxy credentials will be overridden.
		if Nodecfg.TorIsolation &&
			(Nodecfg.OnionProxyUser != "" || Nodecfg.OnionProxyPass != "") {
			fmt.Fprintln(os.Stderr, "Tor isolation set -- "+
				"overriding specified onionproxy user "+
				"credentials ")
		}
		StateCfg.Oniondial = func(network, addr string, timeout time.Duration) (net.Conn, error) {
			proxy := &socks.Proxy{
				Addr:         Nodecfg.OnionProxy,
				Username:     Nodecfg.OnionProxyUser,
				Password:     Nodecfg.OnionProxyPass,
				TorIsolation: Nodecfg.TorIsolation,
			}
			return proxy.DialTimeout(network, addr, timeout)
		}
		// When configured in bridge mode (both --onion and --proxy are configured), it means that the proxy configured by --proxy is not a tor proxy, so override the DNS resolution to use the onion-specific proxy.
		if Nodecfg.Proxy != "" {
			StateCfg.Lookup = func(host string) ([]net.IP, error) {
				return connmgr.TorLookupIP(host, Nodecfg.OnionProxy)
			}
		}
	} else {
		StateCfg.Oniondial = StateCfg.Dial
	}
	// Specifying --noonion means the onion address dial function results in an error.
	if Nodecfg.NoOnion {
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
	Config = defCfg
}

// DefaultNodeConfig is the default configuration for node
func DefaultNodeConfig() *NodeCfg {
	user := pu.GenKey()
	pass := pu.GenKey()
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
		LogLevels: logger.GetDefaultConfig(),
	}
}
