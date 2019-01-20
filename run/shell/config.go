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
var Log = cl.NewSubSystem("run/shell", "trace")
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
	Levels       map[string]*cl.SubSystem
}

var (
	DefaultDataDir      = n.DefaultDataDir
	DefaultAppDataDir   = filepath.Join(n.DefaultHomeDir, "shell")
	DefaultConfFileName = filepath.Join(filepath.Join(n.DefaultHomeDir, "shell"), "conf")
)

// Config is the combined app and log levels configuration
var Config = DefaultConfig()

// Command is a command to send RPC queries to bitcoin RPC protocol server for node and wallet queries
var Command = climax.Command{
	Name:  "shell",
	Brief: "parallelcoin combined full node and wallet",
	Help:  "distrubutes, verifies and mines blocks for the parallelcoin duo cryptocurrency, as well as optionally providing search indexes for transactions in the database, and provides RPC and GUI interfaces for a built-in wallet",
	Flags: []climax.Flag{
		podutil.GenerateFlag("version", "V", `--version`, `show version number and quit`, false),
		podutil.GenerateFlag("configfile", "C", "--configfile=/path/to/conf", "path to configuration file", true),
		podutil.GenerateFlag("datadir", "D", "--datadir=/home/user/.pod", "set the base directory for elements shared between modules", true),

		podutil.GenerateFlag("init", "", "--init", "resets configuration to defaults", false),
		podutil.GenerateFlag("save", "", "--save", "saves current configuration", false),

		podutil.GenerateFlag("create", "", "--create", "create a new wallet if it does not exist", false),
		podutil.GenerateFlag("createtemp", "", "--createtemp", "create temporary wallet (pass=password), must call with --datadir", false),

		podutil.GenerateFlag(`dropcfindex`, ``, `--dropcfindex`, `deletes the index used for committed filtering (CF) support from the database on start up and then exits`, false),
		podutil.GenerateFlag(`droptxindex`, ``, `--droptxindex`, `deletes the hash-based transaction index from the database on start up and then exits.`, false),
		podutil.GenerateFlag(`dropaddrindex`, ``, `--dropaddrindex`, `deletes the address-based transaction index from the database on start up and then exits`, false),

		podutil.GenerateFlag(`addpeers`, ``, `--addpeers=some.peer.com:11047`, `adds a peer to the peers database to try to connect to`, true),
		podutil.GenerateFlag(`connectpeers`, ``, `--connectpeers=some.peer.com:11047`, `adds a peer to a connect-only whitelist`, true),
		podutil.GenerateFlag(`disablelisten`, ``, `--disablelisten=true`, `disables the P2P listener`, true),
		podutil.GenerateFlag(`listeners`, `S`, `--listeners=127.0.0.1:11047`, `sets an address to listen for P2P connections`, true),
		podutil.GenerateFlag(`maxpeers`, ``, `--maxpeers=100`, `sets max number of peers to open connections to at once`, true),
		podutil.GenerateFlag(`disablebanning`, ``, `--disablebanning`, `disable banning of misbehaving peers`, false),
		podutil.GenerateFlag(`banduration`, ``, `--banduration=1h`, `how long to ban misbehaving peers - valid time units are {s, m, h},  minimum 1s`, true),
		podutil.GenerateFlag(`banthreshold`, ``, `--banthreshold=100`, `maximum allowed ban score before disconnecting and banning misbehaving peers`, true),
		podutil.GenerateFlag(`whitelists`, ``, `--whitelists=127.0.0.1:11047`, `add an IP network or IP that will not be banned - eg. 192.168.1.0/24 or ::1`, true),
		// podutil.GenerateFlag(`rpcuser`, `u`, `--rpcuser=username`, `RPC username`, true),
		// podutil.GenerateFlag(`rpcpass`, `P`, `--rpcpass=password`, `RPC password`, true),
		// podutil.GenerateFlag(`rpclimituser`, `u`, `--rpclimituser=username`, `limited user RPC username`, true),
		// podutil.GenerateFlag(`rpclimitpass`, `P`, `--rpclimitpass=password`, `limited user RPC password`, true),
		// podutil.GenerateFlag(`rpclisteners`, `s`, `--rpclisteners=127.0.0.1:11048`, `RPC server to connect to`, true),
		// podutil.GenerateFlag(`rpccert`, `c`, `--rpccert=/path/to/rpn.cert`, `RPC server tls certificate chain for validation`, true),
		// podutil.GenerateFlag(`rpckey`, `c`, `--rpccert=/path/to/rpn.key`, `RPC server tls key for validation`, true),
		// podutil.GenerateFlag(`tls`, ``, `--tls=false`, `enable TLS`, true),
		podutil.GenerateFlag(`disablednsseed`, ``, `--disablednsseed=false`, `disable dns seeding`, true),
		podutil.GenerateFlag(`externalips`, ``, `--externalips=192.168.0.1:11048`, `set additional listeners on different address/interfaces`, true),
		podutil.GenerateFlag(`proxy`, ``, `--proxy 127.0.0.1:9050`, `connect via SOCKS5 proxy (eg. 127.0.0.1:9050)`, true),
		podutil.GenerateFlag(`proxyuser`, ``, `--proxyuser username`, `username for proxy server`, true),
		podutil.GenerateFlag(`proxypass`, ``, `--proxypass password`, `password for proxy server`, true),
		podutil.GenerateFlag(`onion`, ``, `--onion 127.0.0.1:9050`, `connect via onion proxy (eg. 127.0.0.1:9050)`, true),
		podutil.GenerateFlag(`onionuser`, ``, `--onionuser username`, `username for onion proxy server`, true),
		podutil.GenerateFlag(`onionpass`, ``, `--onionpass password`, `password for onion proxy server`, true),
		podutil.GenerateFlag(`noonion`, ``, `--noonion=true`, `disable onion proxy`, true),
		podutil.GenerateFlag(`torisolation`, ``, `--torisolation=true`, `enable tor stream isolation by randomising user credentials for each connection`, true),
		podutil.GenerateFlag(`network`, ``, `--network=mainnet`, `connect to specified network: mainnet, testnet, regtestnet or simnet`, true),
		podutil.GenerateFlag(`skipverify`, ``, `--skipverify=false`, `do not verify tls certificates (not recommended!)`, true),
		podutil.GenerateFlag(`addcheckpoints`, ``, `--addcheckpoints <height>:<hash>`, `add custom checkpoints`, true),
		podutil.GenerateFlag(`disablecheckpoints`, ``, `--disablecheckpoints=true`, `disable all checkpoints`, true),
		podutil.GenerateFlag(`dbtype`, ``, `--dbtype=ffldb`, `set database backend type`, true),
		podutil.GenerateFlag(`profile`, ``, `--profile=127.0.0.1:3131`, `start HTTP profiling server on given address`, true),
		podutil.GenerateFlag(`cpuprofile`, ``, `--cpuprofile=127.0.0.1:3232`, `start cpu profiling server on given address`, true),
		podutil.GenerateFlag(`upnp`, ``, `--upnp=true`, `enables the use of UPNP to establish inbound port redirections`, true),
		podutil.GenerateFlag(`minrelaytxfee`, ``, `--minrelaytxfee=1`, `the minimum transaction fee in DUO/Kb to be considered a nonzero fee`, true),
		podutil.GenerateFlag(`freetxrelaylimit`, ``, `--freetxrelaylimit=100`, `limit amount of free transactions relayed in thousand bytes per minute`, true),
		podutil.GenerateFlag(`norelaypriority`, ``, `--norelaypriority=true`, `do not require free or low-fee transactions to have high priority for relaying`, true),
		podutil.GenerateFlag(`trickleinterval`, ``, `--trickleinterval=1`, `time in seconds between attempts to send new inventory to a connected peer`, true),
		podutil.GenerateFlag(`maxorphantxs`, ``, `--maxorphantxs=100`, `set maximum number of orphans transactions to keep in memory`, true),
		podutil.GenerateFlag(`algo`, ``, `--algo=random`, `set algorithm to be used by cpu miner`, true),
		podutil.GenerateFlag(`generate`, ``, `--generate=true`, `set CPU miner to generate blocks`, true),
		podutil.GenerateFlag(`genthreads`, ``, `--genthreads=-1`, `set number of threads to generate using CPU, -1 = all available`, true),
		podutil.GenerateFlag(`miningaddrs`, ``, `--miningaddrs=aoeuaoe0760oeu0`, `add an address to the list of addresses to make block payments to from miners`, true),
		podutil.GenerateFlag(`minerlistener`, ``, `--minerlistener=127.0.0.1:11011`, `set the port for a miner work dispatch server to listen on`, true),
		podutil.GenerateFlag(`minerpass`, ``, `--minerpass=pa55word`, `set the encryption password to prevent leaking or MiTM attacks on miners`, true),
		podutil.GenerateFlag(`blockminsize`, ``, `--blockminsize=80`, `mininum block size in bytes to be used when creating a block`, true),
		podutil.GenerateFlag(`blockmaxsize`, ``, `--blockmaxsize=1024000`, `maximum block size in bytes to be used when creating a block`, true),
		podutil.GenerateFlag(`blockminweight`, ``, `--blockminweight=500`, `mininum block weight to be used when creating a block`, true),
		podutil.GenerateFlag(`blockmaxweight`, ``, `--blockmaxweight=10000`, `maximum block weight to be used when creating a block`, true),
		podutil.GenerateFlag(`blockprioritysize`, ``, `--blockprioritysize=256`, `size in bytes for high-priority/low-fee transactions when creating a block`, true),
		podutil.GenerateFlag(`uacomment`, ``, `--uacomment=joeblogsminers`, `comment to add to the user agent - see BIP 14 for more information.`, true),
		podutil.GenerateFlag(`nopeerbloomfilters`, ``, `--nopeerbloomfilters=false`, `disable bloom filtering support`, true),
		podutil.GenerateFlag(`nocfilters`, ``, `--nocfilters=false`, `disable committed filtering (CF) support`, true),
		podutil.GenerateFlag(`sigcachemaxsize`, ``, `--sigcachemaxsize=1000`, `the maximum number of entries in the signature verification cache`, true),
		podutil.GenerateFlag(`blocksonly`, ``, `--blocksonly=true`, `do not accept transactions from remote peers`, true),
		podutil.GenerateFlag(`txindex`, ``, `--txindex=true`, `maintain a full hash-based transaction index which makes all transactions available via the getrawtransaction`, true),
		podutil.GenerateFlag(`addrindex`, ``, `--addrindex=true`, `maintain a full address-based transaction index which makes the searchrawtransactions RPC available`, true),
		podutil.GenerateFlag(`relaynonstd`, ``, `--relaynonstd=true`, `relay non-standard transactions regardless of the default settings for the active network`, true),
		podutil.GenerateFlag(`rejectnonstd`, ``, `--rejectnonstd=false`, `reject non-standard transactions regardless of the default settings for the active network`, true),

		podutil.GenerateFlag("appdatadir", "", "--appdatadir=/path/to/appdatadir", "set app data directory for wallet, configuration and logs", true),
		podutil.GenerateFlag("testnet3", "", "--testnet=true", "use testnet", true),
		podutil.GenerateFlag("simnet", "", "--simnet=true", "use simnet", true),
		podutil.GenerateFlag("noinitialload", "", "--noinitialload=true", "defer wallet creation/opening on startup and enable loading wallets over RPC (default with --gui)", true),
		podutil.GenerateFlag("network", "", "--network=mainnet", "connect to specified network: mainnet, testnet, regtestnet or simnet", true),
		podutil.GenerateFlag("profile", "", "--profile=true", "enable HTTP profiling on given port -- NOTE port must be between 1024 and 65536", true),
		podutil.GenerateFlag("gui", "", "--gui=true", "launch GUI (wallet unlock is deferred to let GUI handle)", true),
		podutil.GenerateFlag("walletpass", "", "--walletpass=somepassword", "the public wallet password - only required if the wallet was created with one", true),
		// podutil.GenerateFlag("rpcconnect", "", "--rpcconnect=some.address.com:11048", "connect to the RPC of a parallelcoin node for chain queries", true),
		podutil.GenerateFlag("cafile", "", "--cafile=/path/to/cafile", "file containing root certificates to authenticate TLS connections with pod", true),
		// podutil.GenerateFlag("enableclienttls", "", "--enableclienttls=false", "enable TLS for the RPC client", true),
		// podutil.GenerateFlag("podusername", "", "--podusername=user", "username for node RPC authentication", true),
		// podutil.GenerateFlag("podpassword", "", "--podpassword=pa55word", "password for node RPC authentication", true),
		// podutil.GenerateFlag("proxy", "", "--proxy=127.0.0.1:9050", "address for proxy for outbound connections", true),
		// podutil.GenerateFlag("proxyuser", "", "--proxyuser=user", "username for proxy", true),
		// podutil.GenerateFlag("proxypass", "", "--proxypass=pa55word", "password for proxy", true),
		podutil.GenerateFlag("rpccert", "", "--rpccert=/path/to/rpccert", "file containing the RPC tls certificate", true),
		podutil.GenerateFlag("rpckey", "", "--rpckey=/path/to/rpckey", "file containing RPC tls key", true),
		podutil.GenerateFlag("onetimetlskey", "", "--onetimetlskey=true", "generate a new TLS certpair but only write certs to disk", true),
		podutil.GenerateFlag("enableservertls", "", "--enableservertls=false", "enable TLS on wallet RPC", true),
		podutil.GenerateFlag("legacyrpclisteners", "", "--legacyrpclisteners=127.0.0.1:11046", "add a listener for the legacy RPC", true),
		podutil.GenerateFlag("legacyrpcmaxclients", "", "--legacyrpcmaxclients=10", "maximum number of connections for legacy RPC", true),
		podutil.GenerateFlag("legacyrpcmaxwebsockets", "", "--legacyrpcmaxwebsockets=10", "maximum number of websockets for legacy RPC", true),
		podutil.GenerateFlag("username", "-u", "--username=user", "username for wallet RPC, used also for node if podusername is empty", true),
		podutil.GenerateFlag("password", "-P", "--password=pa55word", "password for wallet RPC, also used for node if podpassord", true),
		podutil.GenerateFlag("experimentalrpclisteners", "", "--experimentalrpclisteners=127.0.0.1:11045", "enable experimental RPC service on this address", true),

		podutil.GenerateFlag("debuglevel", "d", "--debuglevel=trace", "sets debuglevel, default info, sets the baseline for others not specified below (logging is per-library)", true),

		podutil.GenerateFlag("lib-blockchain", "", "--lib-blockchain=info", "", true),
		podutil.GenerateFlag("lib-connmgr", "", "--lib-connmgr=info", "", true),
		podutil.GenerateFlag("lib-database-ffldb", "", "--lib-database-ffldb=info", "", true),
		podutil.GenerateFlag("lib-database", "", "--lib-database=info", "", true),
		podutil.GenerateFlag("lib-mining-cpuminer", "", "--lib-mining-cpuminer=info", "", true),
		podutil.GenerateFlag("lib-mining", "", "--lib-mining=info", "", true),
		podutil.GenerateFlag("lib-netsync", "", "--lib-netsync=info", "", true),
		podutil.GenerateFlag("lib-peer", "", "--lib-peer=info", "", true),
		podutil.GenerateFlag("lib-rpcclient", "", "--lib-rpcclient=info", "", true),
		podutil.GenerateFlag("lib-txscript", "", "--lib-txscript=info", "", true),
		podutil.GenerateFlag("node", "", "--node=info", "", true),
		podutil.GenerateFlag("node-mempool", "", "--node-mempool=info", "", true),
		podutil.GenerateFlag("spv", "", "--spv=info", "", true),
		podutil.GenerateFlag("wallet", "", "--wallet=info", "", true),
		podutil.GenerateFlag("wallet-chain", "", "--wallet-chain=info", "", true),
		podutil.GenerateFlag("wallet-legacyrpc", "", "--wallet-legacyrpc=info", "", true),
		podutil.GenerateFlag("wallet-rpcserver", "", "--wallet-rpcserver=info", "", true),
		podutil.GenerateFlag("wallet-tx", "", "--wallet-tx=info", "", true),
		podutil.GenerateFlag("wallet-votingpool", "", "--wallet-votingpool=info", "", true),
		podutil.GenerateFlag("wallet-waddrmgr", "", "--wallet-waddrmgr=info", "", true),
		podutil.GenerateFlag("wallet-wallet", "", "--wallet-wallet=info", "", true),
		podutil.GenerateFlag("wallet-wtxmgr", "", "--wallet-wtxmgr=info", "", true),
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
			writeDefaultConfig(cfgFile)
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
				writeDefaultConfig(cfgFile)
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
		case "fatal", "error", "info", "debug", "trace":
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
		podutil.NormalizeAddresses(*r, n.DefaultPort, &Config.Node.AddPeers)
	}
	if getIfIs(ctx, "connectpeers", r) {
		podutil.NormalizeAddresses(*r, n.DefaultPort, &Config.Node.ConnectPeers)
	}
	if getIfIs(ctx, "disablelisten", r) {
		Config.Node.DisableListen = *r == "true"
	}
	if getIfIs(ctx, "listeners", r) {
		podutil.NormalizeAddresses(*r, n.DefaultPort, &Config.Node.Listeners)
	}
	if getIfIs(ctx, "maxpeers", r) {
		if err := podutil.ParseInteger(*r, "maxpeers", &Config.Node.MaxPeers); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if getIfIs(ctx, "disablebanning", r) {
		Config.Node.DisableBanning = *r == "true"
	}
	if getIfIs(ctx, "banduration", r) {
		if err := podutil.ParseDuration(*r, "banduration", &Config.Node.BanDuration); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if getIfIs(ctx, "banthreshold", r) {
		var bt int
		if err := podutil.ParseInteger(*r, "banthtreshold", &bt); err != nil {
			log <- cl.Wrn(err.Error())
		} else {
			Config.Node.BanThreshold = uint32(bt)
		}
	}
	if getIfIs(ctx, "whitelists", r) {
		podutil.NormalizeAddresses(*r, n.DefaultPort, &Config.Node.Whitelists)
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
	podutil.NormalizeAddresses(n.DefaultRPCListener, n.DefaultRPCPort, &Config.Node.RPCListeners)
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
		podutil.NormalizeAddresses(*r, n.DefaultPort, &Config.Node.ExternalIPs)
	}
	if getIfIs(ctx, "proxy", r) {
		podutil.NormalizeAddress(*r, "9050", &Config.Node.Proxy)
	}
	if getIfIs(ctx, "proxyuser", r) {
		Config.Node.ProxyUser = *r
	}
	if getIfIs(ctx, "proxypass", r) {
		Config.Node.ProxyPass = *r
	}
	if getIfIs(ctx, "onion", r) {
		podutil.NormalizeAddress(*r, "9050", &Config.Node.OnionProxy)
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
		if err := podutil.ParseFloat(*r, "minrelaytxfee", &Config.Node.MinRelayTxFee); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if getIfIs(ctx, "freetxrelaylimit", r) {
		if err := podutil.ParseFloat(*r, "freetxrelaylimit", &Config.Node.FreeTxRelayLimit); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if getIfIs(ctx, "norelaypriority", r) {
		Config.Node.NoRelayPriority = *r == "true"
	}
	if getIfIs(ctx, "trickleinterval", r) {
		if err := podutil.ParseDuration(*r, "trickleinterval", &Config.Node.TrickleInterval); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if getIfIs(ctx, "maxorphantxs", r) {
		if err := podutil.ParseInteger(*r, "maxorphantxs", &Config.Node.MaxOrphanTxs); err != nil {
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
		if err := podutil.ParseInteger(*r, "genthreads", &gt); err != nil {
			log <- cl.Wrn(err.Error())
		} else {
			Config.Node.GenThreads = int32(gt)
		}
	}
	if getIfIs(ctx, "miningaddrs", r) {
		Config.Node.MiningAddrs = strings.Split(*r, " ")
	}
	if getIfIs(ctx, "minerlistener", r) {
		podutil.NormalizeAddress(*r, n.DefaultRPCPort, &Config.Node.MinerListener)
	}
	if getIfIs(ctx, "minerpass", r) {
		Config.Node.MinerPass = *r
	}
	if getIfIs(ctx, "blockminsize", r) {
		if err := podutil.ParseUint32(*r, "blockminsize", &Config.Node.BlockMinSize); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if getIfIs(ctx, "blockmaxsize", r) {
		if err := podutil.ParseUint32(*r, "blockmaxsize", &Config.Node.BlockMaxSize); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if getIfIs(ctx, "blockminweight", r) {
		if err := podutil.ParseUint32(*r, "blockminweight", &Config.Node.BlockMinWeight); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if getIfIs(ctx, "blockmaxweight", r) {
		if err := podutil.ParseUint32(*r, "blockmaxweight", &Config.Node.BlockMaxWeight); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if getIfIs(ctx, "blockprioritysize", r) {
		if err := podutil.ParseUint32(*r, "blockmaxweight", &Config.Node.BlockPrioritySize); err != nil {
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
		if err := podutil.ParseInteger(*r, "sigcachemaxsize", &scms); err != nil {
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
		podutil.NormalizeAddress(*r, "3131", &Config.Wallet.Profile)
	}
	if getIfIs(ctx, "gui", r) {
		Config.Wallet.GUI = *r == "true"
	}
	if getIfIs(ctx, "walletpass", r) {
		Config.Wallet.WalletPass = *r
	}
	// if getIfIs(ctx, "rpcconnect", r) {
	podutil.NormalizeAddress(n.DefaultRPCListener, "11048", &Config.Wallet.RPCConnect)
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
		podutil.NormalizeAddresses(*r, "11046", &Config.Wallet.LegacyRPCListeners)
	}
	if getIfIs(ctx, "legacyrpcmaxclients", r) {
		var bt int
		if err := podutil.ParseInteger(*r, "legacyrpcmaxclients", &bt); err != nil {
			log <- cl.Wrn(err.Error())
		} else {
			Config.Wallet.LegacyRPCMaxClients = int64(bt)
		}
	}
	if getIfIs(ctx, "legacyrpcmaxwebsockets", r) {
		_, err := fmt.Sscanf(*r, "%d", Config.Wallet.LegacyRPCMaxWebsockets)
		if err != nil {
			log <- cl.Errorf{
				"malformed legacyrpcmaxwebsockets: `%s` leaving set at `%d`",
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
		podutil.NormalizeAddresses(*r, "11045", &Config.Wallet.ExperimentalRPCListeners)
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

func writeDefaultConfig(cfgFile string) {
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
	rpcusername := podutil.GenerateKey()
	rpcpassword := podutil.GenerateKey()
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
		Levels: logger.GetDefault(),
	}
}
