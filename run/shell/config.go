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

// Command is a command to send RPC queries to bitcoin RPC protocol server for node and wallet queries
var Command = climax.Command{
	Name:  "shell",
	Brief: "parallelcoin combined full node and wallet",
	Help:  "distrubutes, verifies and mines blocks for the parallelcoin duo cryptocurrency, as well as optionally providing search indexes for transactions in the database, and provides RPC and GUI interfaces for a built-in wallet",
	Flags: []climax.Flag{
		pu.GenFlag("version", "V", `--version`, `show version number and quit`, false),
		pu.GenFlag("configfile", "C", "--configfile=/path/to/conf", "path to configuration file", true),
		pu.GenFlag("datadir", "D", "--datadir=/home/user/.pod", "set the base directory for elements shared between modules", true),

		pu.GenFlag("init", "", "--init", "resets configuration to defaults", false),
		pu.GenFlag("save", "", "--save", "saves current configuration", false),

		pu.GenFlag("create", "", "--create", "create a new wallet if it does not exist", false),
		pu.GenFlag("createtemp", "", "--createtemp", "create temporary wallet (pass=password), must call with --datadir", false),

		pu.GenFlag(`dropcfindex`, ``, `--dropcfindex`, `deletes the index used for committed filtering (CF) support from the database on start up and then exits`, false),
		pu.GenFlag(`droptxindex`, ``, `--droptxindex`, `deletes the hash-based transaction index from the database on start up and then exits.`, false),
		pu.GenFlag(`dropaddrindex`, ``, `--dropaddrindex`, `deletes the address-based transaction index from the database on start up and then exits`, false),

		pu.GenFlag(`addpeers`, ``, `--addpeers=some.peer.com:11047`, `adds a peer to the peers database to try to connect to`, true),
		pu.GenFlag(`connectpeers`, ``, `--connectpeers=some.peer.com:11047`, `adds a peer to a connect-only whitelist`, true),
		pu.GenFlag(`disablelisten`, ``, `--disablelisten=true`, `disables the P2P listener`, true),
		pu.GenFlag(`listeners`, `S`, `--listeners=127.0.0.1:11047`, `sets an address to listen for P2P connections`, true),
		pu.GenFlag(`maxpeers`, ``, `--maxpeers=100`, `sets max number of peers to open connections to at once`, true),
		pu.GenFlag(`disablebanning`, ``, `--disablebanning`, `disable banning of misbehaving peers`, false),
		pu.GenFlag(`banduration`, ``, `--banduration=1h`, `how long to ban misbehaving peers - valid time units are {s, m, h},  minimum 1s`, true),
		pu.GenFlag(`banthreshold`, ``, `--banthreshold=100`, `maximum allowed ban score before disconnecting and banning misbehaving peers`, true),
		pu.GenFlag(`whitelists`, ``, `--whitelists=127.0.0.1:11047`, `add an IP network or IP that will not be banned - eg. 192.168.1.0/24 or ::1`, true),
		// pu.GenFlag(`rpcuser`, `u`, `--rpcuser=username`, `RPC username`, true),
		// pu.GenFlag(`rpcpass`, `P`, `--rpcpass=password`, `RPC password`, true),
		// pu.GenFlag(`rpclimituser`, `u`, `--rpclimituser=username`, `limited user RPC username`, true),
		// pu.GenFlag(`rpclimitpass`, `P`, `--rpclimitpass=password`, `limited user RPC password`, true),
		// pu.GenFlag(`rpclisteners`, `s`, `--rpclisteners=127.0.0.1:11048`, `RPC server to connect to`, true),
		// pu.GenFlag(`rpccert`, `c`, `--rpccert=/path/to/rpn.cert`, `RPC server tls certificate chain for validation`, true),
		// pu.GenFlag(`rpckey`, `c`, `--rpccert=/path/to/rpn.key`, `RPC server tls key for validation`, true),
		// pu.GenFlag(`tls`, ``, `--tls=false`, `enable TLS`, true),
		pu.GenFlag(`disablednsseed`, ``, `--disablednsseed=false`, `disable dns seeding`, true),
		pu.GenFlag(`externalips`, ``, `--externalips=192.168.0.1:11048`, `set additional listeners on different address/interfaces`, true),
		pu.GenFlag(`proxy`, ``, `--proxy 127.0.0.1:9050`, `connect via SOCKS5 proxy (eg. 127.0.0.1:9050)`, true),
		pu.GenFlag(`proxyuser`, ``, `--proxyuser username`, `username for proxy server`, true),
		pu.GenFlag(`proxypass`, ``, `--proxypass password`, `password for proxy server`, true),
		pu.GenFlag(`onion`, ``, `--onion 127.0.0.1:9050`, `connect via onion proxy (eg. 127.0.0.1:9050)`, true),
		pu.GenFlag(`onionuser`, ``, `--onionuser username`, `username for onion proxy server`, true),
		pu.GenFlag(`onionpass`, ``, `--onionpass password`, `password for onion proxy server`, true),
		pu.GenFlag(`noonion`, ``, `--noonion=true`, `disable onion proxy`, true),
		pu.GenFlag(`torisolation`, ``, `--torisolation=true`, `enable tor stream isolation by randomising user credentials for each connection`, true),
		pu.GenFlag(`network`, ``, `--network=mainnet`, `connect to specified network: mainnet, testnet, regtestnet or simnet`, true),
		pu.GenFlag(`skipverify`, ``, `--skipverify=false`, `do not verify tls certificates (not recommended!)`, true),
		pu.GenFlag(`addcheckpoints`, ``, `--addcheckpoints <height>:<hash>`, `add custom checkpoints`, true),
		pu.GenFlag(`disablecheckpoints`, ``, `--disablecheckpoints=true`, `disable all checkpoints`, true),
		pu.GenFlag(`dbtype`, ``, `--dbtype=ffldb`, `set database backend type`, true),
		pu.GenFlag(`profile`, ``, `--profile=127.0.0.1:3131`, `start HTTP profiling server on given address`, true),
		pu.GenFlag(`cpuprofile`, ``, `--cpuprofile=127.0.0.1:3232`, `start cpu profiling server on given address`, true),
		pu.GenFlag(`upnp`, ``, `--upnp=true`, `enables the use of UPNP to establish inbound port redirections`, true),
		pu.GenFlag(`minrelaytxfee`, ``, `--minrelaytxfee=1`, `the minimum transaction fee in DUO/Kb to be considered a nonzero fee`, true),
		pu.GenFlag(`freetxrelaylimit`, ``, `--freetxrelaylimit=100`, `limit amount of free transactions relayed in thousand bytes per minute`, true),
		pu.GenFlag(`norelaypriority`, ``, `--norelaypriority=true`, `do not require free or low-fee transactions to have high priority for relaying`, true),
		pu.GenFlag(`trickleinterval`, ``, `--trickleinterval=1`, `time in seconds between attempts to send new inventory to a connected peer`, true),
		pu.GenFlag(`maxorphantxs`, ``, `--maxorphantxs=100`, `set maximum number of orphans transactions to keep in memory`, true),
		pu.GenFlag(`algo`, ``, `--algo=random`, `set algorithm to be used by cpu miner`, true),
		pu.GenFlag(`generate`, ``, `--generate=true`, `set CPU miner to generate blocks`, true),
		pu.GenFlag(`genthreads`, ``, `--genthreads=-1`, `set number of threads to generate using CPU, -1 = all available`, true),
		pu.GenFlag(`miningaddrs`, ``, `--miningaddrs=aoeuaoe0760oeu0`, `add an address to the list of addresses to make block payments to from miners`, true),
		pu.GenFlag(`minerlistener`, ``, `--minerlistener=127.0.0.1:11011`, `set the port for a miner work dispatch server to listen on`, true),
		pu.GenFlag(`minerpass`, ``, `--minerpass=pa55word`, `set the encryption password to prevent leaking or MiTM attacks on miners`, true),
		pu.GenFlag(`blockminsize`, ``, `--blockminsize=80`, `mininum block size in bytes to be used when creating a block`, true),
		pu.GenFlag(`blockmaxsize`, ``, `--blockmaxsize=1024000`, `maximum block size in bytes to be used when creating a block`, true),
		pu.GenFlag(`blockminweight`, ``, `--blockminweight=500`, `mininum block weight to be used when creating a block`, true),
		pu.GenFlag(`blockmaxweight`, ``, `--blockmaxweight=10000`, `maximum block weight to be used when creating a block`, true),
		pu.GenFlag(`blockprioritysize`, ``, `--blockprioritysize=256`, `size in bytes for high-priority/low-fee transactions when creating a block`, true),
		pu.GenFlag(`uacomment`, ``, `--uacomment=joeblogsminers`, `comment to add to the user agent - see BIP 14 for more information.`, true),
		pu.GenFlag(`nopeerbloomfilters`, ``, `--nopeerbloomfilters=false`, `disable bloom filtering support`, true),
		pu.GenFlag(`nocfilters`, ``, `--nocfilters=false`, `disable committed filtering (CF) support`, true),
		pu.GenFlag(`sigcachemaxsize`, ``, `--sigcachemaxsize=1000`, `the maximum number of entries in the signature verification cache`, true),
		pu.GenFlag(`blocksonly`, ``, `--blocksonly=true`, `do not accept transactions from remote peers`, true),
		pu.GenFlag(`txindex`, ``, `--txindex=true`, `maintain a full hash-based transaction index which makes all transactions available via the getrawtransaction`, true),
		pu.GenFlag(`addrindex`, ``, `--addrindex=true`, `maintain a full address-based transaction index which makes the searchrawtransactions RPC available`, true),
		pu.GenFlag(`relaynonstd`, ``, `--relaynonstd=true`, `relay non-standard transactions regardless of the default settings for the active network`, true),
		pu.GenFlag(`rejectnonstd`, ``, `--rejectnonstd=false`, `reject non-standard transactions regardless of the default settings for the active network`, true),

		pu.GenFlag("appdatadir", "", "--appdatadir=/path/to/appdatadir", "set app data directory for wallet, configuration and logs", true),
		pu.GenFlag("testnet3", "", "--testnet=true", "use testnet", true),
		pu.GenFlag("simnet", "", "--simnet=true", "use simnet", true),
		pu.GenFlag("noinitialload", "", "--noinitialload=true", "defer wallet creation/opening on startup and enable loading wallets over RPC (default with --gui)", true),
		pu.GenFlag("network", "", "--network=mainnet", "connect to specified network: mainnet, testnet, regtestnet or simnet", true),
		pu.GenFlag("profile", "", "--profile=true", "enable HTTP profiling on given port -- NOTE port must be between 1024 and 65536", true),
		pu.GenFlag("gui", "", "--gui=true", "launch GUI (wallet unlock is deferred to let GUI handle)", true),
		pu.GenFlag("walletpass", "", "--walletpass=somepassword", "the public wallet password - only required if the wallet was created with one", true),
		// pu.GenFlag("rpcconnect", "", "--rpcconnect=some.address.com:11048", "connect to the RPC of a parallelcoin node for chain queries", true),
		pu.GenFlag("cafile", "", "--cafile=/path/to/cafile", "file containing root certificates to authenticate TLS connections with pod", true),
		// pu.GenFlag("enableclienttls", "", "--enableclienttls=false", "enable TLS for the RPC client", true),
		// pu.GenFlag("podusername", "", "--podusername=user", "username for node RPC authentication", true),
		// pu.GenFlag("podpassword", "", "--podpassword=pa55word", "password for node RPC authentication", true),
		// pu.GenFlag("proxy", "", "--proxy=127.0.0.1:9050", "address for proxy for outbound connections", true),
		// pu.GenFlag("proxyuser", "", "--proxyuser=user", "username for proxy", true),
		// pu.GenFlag("proxypass", "", "--proxypass=pa55word", "password for proxy", true),
		pu.GenFlag("rpccert", "", "--rpccert=/path/to/rpccert", "file containing the RPC tls certificate", true),
		pu.GenFlag("rpckey", "", "--rpckey=/path/to/rpckey", "file containing RPC tls key", true),
		pu.GenFlag("onetimetlskey", "", "--onetimetlskey=true", "generate a new TLS certpair but only write certs to disk", true),
		pu.GenFlag("enableservertls", "", "--enableservertls=false", "enable TLS on wallet RPC", true),
		pu.GenFlag("legacyrpclisteners", "", "--legacyrpclisteners=127.0.0.1:11046", "add a listener for the legacy RPC", true),
		pu.GenFlag("legacyrpcmaxclients", "", "--legacyrpcmaxclients=10", "maximum number of connections for legacy RPC", true),
		pu.GenFlag("legacyrpcmaxwebsockets", "", "--legacyrpcmaxwebsockets=10", "maximum number of websockets for legacy RPC", true),
		pu.GenFlag("username", "-u", "--username=user", "username for wallet RPC, used also for node if podusername is empty", true),
		pu.GenFlag("password", "-P", "--password=pa55word", "password for wallet RPC, also used for node if podpassord", true),
		pu.GenFlag("experimentalrpclisteners", "", "--experimentalrpclisteners=127.0.0.1:11045", "enable experimental RPC service on this address", true),

		pu.GenFlag("debuglevel", "d", "--debuglevel=trace", "sets debuglevel, default info, sets the baseline for others not specified below (logging is per-library)", true),

		pu.GenFlag("lib-addrmgr", "", "--lib-addrmg=info", "", true),
		pu.GenFlag("lib-blockchain", "", "--lib-blockchain=info", "", true),
		pu.GenFlag("lib-connmgr", "", "--lib-connmgr=info", "", true),
		pu.GenFlag("lib-database-ffldb", "", "--lib-database-ffldb=info", "", true),
		pu.GenFlag("lib-database", "", "--lib-database=info", "", true),
		pu.GenFlag("lib-mining-cpuminer", "", "--lib-mining-cpuminer=info", "", true),
		pu.GenFlag("lib-mining", "", "--lib-mining=info", "", true),
		pu.GenFlag("lib-netsync", "", "--lib-netsync=info", "", true),
		pu.GenFlag("lib-peer", "", "--lib-peer=info", "", true),
		pu.GenFlag("lib-rpcclient", "", "--lib-rpcclient=info", "", true),
		pu.GenFlag("lib-txscript", "", "--lib-txscript=info", "", true),
		pu.GenFlag("node", "", "--node=info", "", true),
		pu.GenFlag("node-mempool", "", "--node-mempool=info", "", true),
		pu.GenFlag("spv", "", "--spv=info", "", true),
		pu.GenFlag("wallet", "", "--wallet=info", "", true),
		pu.GenFlag("wallet-chain", "", "--wallet-chain=info", "", true),
		pu.GenFlag("wallet-legacyrpc", "", "--wallet-legacyrpc=info", "", true),
		pu.GenFlag("wallet-rpcserver", "", "--wallet-rpcserver=info", "", true),
		pu.GenFlag("wallet-tx", "", "--wallet-tx=info", "", true),
		pu.GenFlag("wallet-votingpool", "", "--wallet-votingpool=info", "", true),
		pu.GenFlag("wallet-waddrmgr", "", "--wallet-waddrmgr=info", "", true),
		pu.GenFlag("wallet-wallet", "", "--wallet-wallet=info", "", true),
		pu.GenFlag("wallet-wtxmgr", "", "--wallet-wtxmgr=info", "", true),
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
