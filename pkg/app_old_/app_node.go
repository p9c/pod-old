package app_old

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"git.parallelcoin.io/dev/pod/cmd/node"
	"git.parallelcoin.io/dev/pod/pkg/util"
	cl "git.parallelcoin.io/dev/pod/pkg/util/cl"
	"github.com/tucnak/climax"
)

// NodeCfg is the combined app and logging configuration data
type NodeCfg struct {
	Node      *node.Config
	LogLevels map[string]string
	params    *node.Params
}

// serviceOptions defines the configuration options for the daemon as a service on Windows.
type serviceOptions struct {
	ServiceCommand string `short:"s" long:"service" description:"Service command {install, remove, start, stop}"`
}

// NodeCommand is a command to send RPC queries to bitcoin RPC protocol server for node and wallet queries
var NodeCommand = climax.Command{
	Name:  "node",
	Brief: "parallelcoin full node",
	Help:  "distrubutes, verifies and mines blocks for the parallelcoin duo cryptocurrency, as well as optionally providing search indexes for transactions in the database",
	Flags: []climax.Flag{

		t("version", "V", "show version number and quit"),

		s("configfile", "C", node.DefaultConfigFile, "path to configuration file"),
		s("datadir", "D", node.DefaultDataDir, "path to configuration directory"),

		t("init", "", "resets configuration to defaults"),
		t("save", "", "saves current configuration"),

		f("network", "mainnet", "connect to (mainnet|testnet|simnet)"),

		f("txindex", "true", "enable transaction index"),
		f("addrindex", "true", "enable address index"),
		t("dropcfindex", "", "delete committed filtering (CF) index then exit"),
		t("droptxindex", "", "deletes transaction index then exit"),
		t("dropaddrindex", "", "deletes the address index then exits"),

		s("listeners", "S", node.DefaultListener, "sets an address to listen for P2P connections"),
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
		{
			Usecase:     " -D test -d trace",
			Description: "run using the configuration in the 'test' directory with trace logging",
		},
	},
	Handle: func(ctx climax.Context) int {
		var dl string
		var ok bool
		if dl, ok = ctx.Get("debuglevel"); ok {
			log <- cl.Tracef{"setting debug level %s", dl}
			NodeConfig.Node.DebugLevel = dl
			Log.SetLevel(dl)
			ll := GetAllSubSystems()
			for i := range ll {
				ll[i].SetLevel(dl)
			}
		}
		if ctx.Is("version") {
			fmt.Println("pod/node version", node.Version())
			return 0
		}
		var datadir, cfgFile string
		if datadir, ok = ctx.Get("datadir"); !ok {

			datadir = util.AppDataDir("pod", false)
		}
		cfgFile = filepath.Join(filepath.Join(datadir, "node"), "conf.json")
		log <- cl.Debug{"DataDir", datadir, "cfgFile", cfgFile}
		if r, ok := getIfIs(&ctx, "configfile"); ok {

			cfgFile = r
		}
		if ctx.Is("init") {

			log <- cl.Debugf{"writing default configuration to %s", cfgFile}
			WriteDefaultNodeConfig(datadir)
		} else {

			log <- cl.Infof{"loading configuration from %s", cfgFile}
			if _, err := os.Stat(cfgFile); os.IsNotExist(err) {

				log <- cl.Wrn("configuration file does not exist, creating new one")
				WriteDefaultNodeConfig(datadir)
			} else {

				log <- cl.Debug{"reading app configuration from", cfgFile}
				cfgData, err := ioutil.ReadFile(cfgFile)
				if err != nil {

					log <- cl.Error{"reading app config file:", err.Error()}
					WriteDefaultNodeConfig(datadir)
				} else {

					log <- cl.Trace{"parsing app configuration", string(cfgData)}
					err = json.Unmarshal(cfgData, &NodeConfig)
					if err != nil {

						log <- cl.Error{"parsing app config file:", err.Error()}
						WriteDefaultNodeConfig(datadir)
					}
				}
			}
			switch {
			case NodeConfig.Node.TestNet3:
				log <- cl.Info{"running on testnet"}
				NodeConfig.params = &node.TestNet3Params
			case NodeConfig.Node.SimNet:
				log <- cl.Info{"running on simnet"}
				NodeConfig.params = &node.SimNetParams
			default:
				log <- cl.Info{"running on mainnet"}
				NodeConfig.params = &node.MainNetParams
			}
		}
		configNode(NodeConfig.Node, &ctx, cfgFile)
		runNode(NodeConfig.Node, NodeConfig.params)
		return 0
	},
}

// NodeConfig is the combined app and log levels configuration
var NodeConfig = DefaultNodeConfig(node.DefaultDataDir)

// StateCfg is a reference to the main node state configuration struct
var StateCfg = node.StateCfg

var aN = filepath.Base(os.Args[0])

var appName = strings.TrimSuffix(aN, filepath.Ext(aN))

// runServiceCommand is only set to a real function on Windows.  It is used to parse and execute service commands specified via the -s flag.
var runServiceCommand func(string) error

var usageMessage = fmt.Sprintf(
	"use `%s help node` to show usage", appName)
