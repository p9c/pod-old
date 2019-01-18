package node

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"git.parallelcoin.io/pod/lib/clog"
	n "git.parallelcoin.io/pod/module/node"
	"git.parallelcoin.io/pod/module/node/mempool"
	"git.parallelcoin.io/pod/run/logger"
	"github.com/tucnak/climax"
)

var log = clog.NewSubSystem("Node", clog.Ninf)

// Config is the default configuration native to ctl
var Config = new(n.Config)

// ConfigAndLog is the combined app and logging configuration data
type ConfigAndLog struct {
	Node   *n.Config
	Levels map[string]string
}

// CombinedCfg is the combined app and log levels configuration
var CombinedCfg = ConfigAndLog{
	Node:   Config,
	Levels: logger.Levels,
}

// Command is a command to send RPC queries to bitcoin RPC protocol server for node and wallet queries
var Command = climax.Command{
	Name:  "node",
	Brief: "parallelcoin full node",
	Help:  "distrubutes, verifies and mines blocks for the parallelcoin duo cryptocurrency, as well as optionally providing search indexes for transactions in the database",
	Flags: []climax.Flag{
		{
			Name:     "version",
			Short:    "V",
			Usage:    `--version`,
			Help:     `show version number and quit`,
			Variable: false,
		},
		{
			Name:     "configfile",
			Short:    "C",
			Usage:    "--configfile=/path/to/conf",
			Help:     "path to configuration file",
			Variable: true,
		},
		{
			Name:     "datadir",
			Short:    "D",
			Usage:    "--configfile=/path/to/conf",
			Help:     "path to configuration file",
			Variable: true,
		},
		{
			Name:     "init",
			Usage:    "--init",
			Help:     "resets configuration to defaults",
			Variable: false,
		},
		{
			Name:     "save",
			Usage:    "--save",
			Help:     "saves current configuration",
			Variable: false,
		},
		{
			Name:     "debuglevel",
			Short:    "d",
			Usage:    "--debuglevel=trace",
			Help:     "sets debuglevel, default info, sets the baseline for others not specified",
			Variable: true,
		},
		{
			Name:     "log-database",
			Usage:    "--log-database=info",
			Help:     "sets log level for database",
			Variable: true,
		},
		{
			Name:     "log-txscript",
			Usage:    "--log-txscript=info",
			Help:     "sets log level for txscript",
			Variable: true,
		},
		{
			Name:     "log-peer",
			Usage:    "--log-peer=info",
			Help:     "sets log level for peer",
			Variable: true,
		},
		{
			Name:     "log-netsync",
			Usage:    "--log-netsync=info",
			Help:     "sets log level for netsync",
			Variable: true,
		},
		{
			Name:     "log-rpcclient",
			Usage:    "--log-rpcclient=info",
			Help:     "sets log level for rpcclient",
			Variable: true,
		},
		{
			Name:     "addrmgr",
			Usage:    "--log-addrmgr=info",
			Help:     "sets log level for addrmgr",
			Variable: true,
		},
		{
			Name:     "log-blockchain-indexers",
			Usage:    "--log-blockchain-indexers=info",
			Help:     "sets log level for blockchain-indexers",
			Variable: true,
		},
		{
			Name:     "log-blockchain",
			Usage:    "--log-blockchain=info",
			Help:     "sets log level for blockchain",
			Variable: true,
		},
		{
			Name:     "log-mining-cpuminer",
			Usage:    "--log-mining-cpuminer=info",
			Help:     "sets log level for mining-cpuminer",
			Variable: true,
		},
		{
			Name:     "log-mining",
			Usage:    "--log-mining=info",
			Help:     "sets log level for mining",
			Variable: true,
		},
		{
			Name:     "log-mining-controller",
			Usage:    "--log-mining-controller=info",
			Help:     "sets log level for mining-controller",
			Variable: true,
		},
		{
			Name:     "log-connmgr",
			Usage:    "--log-connmgr=info",
			Help:     "sets log level for connmgr",
			Variable: true,
		},
		{
			Name:     "log-spv",
			Usage:    "--log-spv=info",
			Help:     "sets log level for spv",
			Variable: true,
		},
		{
			Name:     "log-node-mempool",
			Usage:    "--log-node-mempool=info",
			Help:     "sets log level for node-mempool",
			Variable: true,
		},
		{
			Name:     "log-node",
			Usage:    "--log-node=info",
			Help:     "sets log level for node",
			Variable: true,
		},
		{
			Name:     "log-wallet-wallet",
			Usage:    "--log-wallet-wallet=info",
			Help:     "sets log level for wallet-wallet",
			Variable: true,
		},
		{
			Name:     "log-wallet-tx",
			Usage:    "--log-wallet-tx=info",
			Help:     "sets log level for wallet-tx",
			Variable: true,
		},
		{
			Name:     "log-wallet-votingpool",
			Usage:    "--log-wallet-votingpool=info",
			Help:     "sets log level for wallet-votingpool",
			Variable: true,
		},
		{
			Name:     "log-wallet",
			Usage:    "--log-wallet=info",
			Help:     "sets log level for wallet",
			Variable: true,
		},
		{
			Name:     "log-wallet-chain",
			Usage:    "--log-wallet-chain=info",
			Help:     "sets log level for wallet-chain",
			Variable: true,
		},
		{
			Name:     "log-wallet-rpc-rpcserver",
			Usage:    "--log-wallet-rpc-rpcserver=info",
			Help:     "sets log level for wallet-rpc-rpcserver",
			Variable: true,
		},
		{
			Name:     "log-wallet-rpc-legacyrpc",
			Usage:    "--log-wallet-rpc-legacyrpc=info",
			Help:     "sets log level for wallet-rpc-legacyrpc",
			Variable: true,
		},
		{
			Name:     "log-wallet-wtxmgr",
			Usage:    "--log-wallet-wtxmgr=info",
			Help:     "sets log level for wallet-wtxmgr",
			Variable: true,
		},
		{
			Name:     "addpeers",
			Usage:    "--addpeers=some.peer.com:11047",
			Help:     "adds a peer to the peers database to try to connect to",
			Variable: true,
		},
		{
			Name:     "connectpeers",
			Usage:    "--connectpeers=some.peer.com:11047",
			Help:     "adds a peer to a connect-only whitelist",
			Variable: true,
		},
		{
			Name:     "disablelisten",
			Usage:    "--disablelisten=true",
			Help:     "disables the P2P listener",
			Variable: true,
		},
		{
			Name:     "listeners",
			Short:    "S",
			Usage:    "--listeners=127.0.0.1:11047",
			Help:     "sets an address to listen for P2P connections",
			Variable: true,
		},
		{
			Name:     "maxpeers",
			Usage:    "--maxpeers=100",
			Help:     "sets max number of peers to open connections to at once",
			Variable: true,
		},
		{
			Name:     "disablebanning",
			Usage:    "--disablebanning",
			Help:     "disable banning of misbehaving peers",
			Variable: false,
		},
		{
			Name:     "banduration",
			Usage:    "--banduration=1h",
			Help:     "how long to ban misbehaving peers - valid time units are {s, m, h},  minimum 1 second",
			Variable: true,
		},
		{
			Name:     "banthreshold",
			Usage:    "--banthreshold=100",
			Help:     "maximum allowed ban score before disconnecting and banning misbehaving peers",
			Variable: true,
		},
		{
			Name:     "whitelists",
			Usage:    "--whitelists=127.0.0.1:11047",
			Help:     "add an IP network or IP that will not be banned - eg. 192.168.1.0/24 or ::1",
			Variable: true,
		},
		{
			Name:     "rpcuser",
			Short:    "u",
			Usage:    "--rpcuser=username",
			Help:     "RPC username",
			Variable: true,
		},
		{
			Name:     "rpcpass",
			Short:    "P",
			Usage:    "--rpcpass=password",
			Help:     "RPC password",
			Variable: true,
		},
		{
			Name:     "rpclimituser",
			Short:    "u",
			Usage:    "--rpclimituser=username",
			Help:     "limited user RPC username",
			Variable: true,
		},
		{
			Name:     "rpclimitpass",
			Short:    "P",
			Usage:    "--rpclimitpass=password",
			Help:     "limited user RPC password",
			Variable: true,
		},
		{
			Name:     "rpclisteners",
			Short:    "s",
			Usage:    "--rpclisteners=127.0.0.1:11048",
			Help:     "RPC server to connect to",
			Variable: true,
		},
		{
			Name:     "rpccert",
			Short:    "c",
			Usage:    "--rpccert=/path/to/rpn.cert",
			Help:     "RPC server tls certificate chain for validation",
			Variable: true,
		},
		{
			Name:     "rpckey",
			Short:    "c",
			Usage:    "--rpccert=/path/to/rpn.key",
			Help:     "RPC server tls key for validation",
			Variable: true,
		},
		{
			Name:     "tls",
			Usage:    "--tls=false",
			Help:     "enable TLS",
			Variable: true,
		},
		{
			Name:     "disablednsseed",
			Usage:    "--disablednsseed=false",
			Help:     "disable dns seeding",
			Variable: true,
		},
		{
			Name:     "externalips",
			Usage:    "--externalips=192.168.0.1:11048",
			Help:     "set additional listeners on different address/interfaces",
			Variable: true,
		},
		{
			Name:     "proxy",
			Usage:    "--proxy 127.0.0.1:9050",
			Help:     "connect via SOCKS5 proxy (eg. 127.0.0.1:9050)",
			Variable: true,
		},
		{
			Name:     "proxyuser",
			Usage:    "--proxyuser username",
			Help:     "username for proxy server",
			Variable: true,
		},
		{
			Name:     "proxypass",
			Usage:    "--proxypass password",
			Help:     "password for proxy server",
			Variable: true,
		},
		{
			Name:     "onion",
			Usage:    "--onion 127.0.0.1:9050",
			Help:     "connect via onion proxy (eg. 127.0.0.1:9050)",
			Variable: true,
		},
		{
			Name:     "onionuser",
			Usage:    "--onionuser username",
			Help:     "username for onion proxy server",
			Variable: true,
		},
		{
			Name:     "onionpass",
			Usage:    "--onionpass password",
			Help:     "password for onion proxy server",
			Variable: true,
		},
		{
			Name:     "noonion",
			Usage:    "--noonion=true",
			Help:     "disable onion proxy",
			Variable: true,
		},
		{
			Name:     "torisolation",
			Usage:    "--torisolation=true",
			Help:     "enable tor stream isolation by randomising user credentials for each connection",
			Variable: true,
		},
		{
			Name:     "network",
			Usage:    "--network=mainnet",
			Help:     "connect to specified network: mainnet, testnet, regtestnet or simnet",
			Variable: true,
		},
		{
			Name:     "skipverify",
			Usage:    "--skipverify=false",
			Help:     "do not verify tls certificates (not recommended!)",
			Variable: true,
		},
		{
			Name:     "addcheckpoints",
			Usage:    "--addcheckpoints <height>:<hash>",
			Help:     "add custom checkpoints",
			Variable: true,
		},
		{
			Name:     "disablecheckpoints",
			Usage:    "--disablecheckpoints=true",
			Help:     "disable all checkpoints",
			Variable: true,
		},
		{
			Name:     "dbtype",
			Usage:    "--dbtype=ffldb",
			Help:     "set database backend type",
			Variable: true,
		},
		{
			Name:     "profile",
			Usage:    "--profile=127.0.0.1:3131",
			Help:     "start HTTP profiling server on given address",
			Variable: true,
		},
		{
			Name:     "cpuprofile",
			Usage:    "--cpuprofile=127.0.0.1:3232",
			Help:     "start cpu profiling server on given address",
			Variable: true,
		},
		{
			Name:     "upnp",
			Usage:    "--upnp=true",
			Help:     "enables the use of UPNP to establish inbound port redirections",
			Variable: true,
		},
		{
			Name:     "minrelaytxfee",
			Usage:    "--minrelaytxfee=1",
			Help:     "the minimum transaction fee in DUO/Kb to be considered a nonzero fee",
			Variable: true,
		},
		{
			Name:     "freetxrelaylimit",
			Usage:    "--freetxrelaylimit=100",
			Help:     "limit amount of free transactions relayed in thousand bytes per minute",
			Variable: true,
		},
		{
			Name:     "norelaypriority",
			Usage:    "--norelaypriority=true",
			Help:     "do not require free or low-fee transactions to have high priority for relaying",
			Variable: true,
		},
		{
			Name:     "trickleinterval",
			Usage:    "--trickleinterval=1",
			Help:     "time in seconds between attempts to send new inventory to a connected peer",
			Variable: true,
		},
		{
			Name:     "maxorphantxs",
			Usage:    "--maxorphantxs=100",
			Help:     "set maximum number of orphans transactions to keep in memory",
			Variable: true,
		},
		{
			Name:     "algo",
			Usage:    "--algo=random",
			Help:     "set algorithm to be used by cpu miner",
			Variable: true,
		},
		{
			Name:     "generate",
			Usage:    "--generate=true",
			Help:     "set CPU miner to generate blocks",
			Variable: true,
		},
		{
			Name:     "genthreads",
			Usage:    "--genthreads=-1",
			Help:     "set number of threads to generate using CPU, -1 = all available",
			Variable: true,
		},
		{
			Name:     "miningaddrs",
			Usage:    "--miningaddrs=aoeuaoe0760oeu0",
			Help:     "add an address to the list of addresses to make block payments to from miners",
			Variable: true,
		},
		{
			Name:     "minerlistener",
			Usage:    "--minerlistener=127.0.0.1:11011",
			Help:     "set the port for a miner work dispatch server to listen on",
			Variable: true,
		},
		{
			Name:     "minerpass",
			Usage:    "--minerpass=pa55word",
			Help:     "set the encryption password to prevent leaking or MiTM attacks on miners",
			Variable: true,
		},
		{
			Name:     "blockminsize",
			Usage:    "--blockminsize=80",
			Help:     "mininum block size in bytes to be used when creating a block",
			Variable: true,
		},
		{
			Name:     "blockmaxsize",
			Usage:    "--blockmaxsize=1024000",
			Help:     "maximum block size in bytes to be used when creating a block",
			Variable: true,
		},
		{
			Name:     "blockminweight",
			Usage:    "--blockminweight=500",
			Help:     "mininum block weight to be used when creating a block",
			Variable: true,
		},
		{
			Name:     "blockmaxweight",
			Usage:    "--blockmaxweight=10000",
			Help:     "maximum block weight to be used when creating a block",
			Variable: true,
		},
		{
			Name:     "blockprioritysize",
			Usage:    "--blockprioritysize=256",
			Help:     "size in bytes for high-priority/low-fee transactions when creating a block",
			Variable: true,
		},
		{
			Name:     "uacomment",
			Usage:    "--uacomment=joeblogsminers",
			Help:     "comment to add to the user agent - see BIP 14 for more information.",
			Variable: true,
		},
		{
			Name:     "nopeerbloomfilters",
			Usage:    "--nopeerbloomfilters=false",
			Help:     "disable bloom filtering support",
			Variable: true,
		},
		{
			Name:     "nocfilters",
			Usage:    "--nocfilters=false",
			Help:     "disable committed filtering (CF) support",
			Variable: true,
		},
		{
			Name:     "dropcfindex",
			Usage:    "--dropcfindex",
			Help:     "deletes the index used for committed filtering (CF) support from the database on start up and then exits",
			Variable: false,
		},
		{
			Name:     "sigcachemaxsize",
			Usage:    "--sigcachemaxsize=1000",
			Help:     "the maximum number of entries in the signature verification cache",
			Variable: true,
		},
		{
			Name:     "blocksonly",
			Usage:    "--blocksonly=true",
			Help:     "do not accept transactions from remote peers",
			Variable: true,
		},
		{
			Name:     "txindex",
			Usage:    "--txindex=true",
			Help:     "maintain a full hash-based transaction index which makes all transactions available via the getrawtransaction RPC",
			Variable: true,
		},
		{
			Name:     "droptxindex",
			Usage:    "--droptxindex",
			Help:     "deletes the hash-based transaction index from the database on start up and then exits.",
			Variable: false,
		},
		{
			Name:     "addrindex",
			Usage:    "--addrindex=true",
			Help:     "maintain a full address-based transaction index which makes the searchrawtransactions RPC available",
			Variable: true,
		},
		{
			Name:     "dropaddrindex",
			Usage:    "--dropaddrindex",
			Help:     "deletes the address-based transaction index from the database on start up and then exits",
			Variable: false,
		},
		{
			Name:     "relaynonstd",
			Usage:    "--relaynonstd=true",
			Help:     "relay non-standard transactions regardless of the default settings for the active network",
			Variable: true,
		},
		{
			Name:     "rejectnonstd",
			Usage:    "--rejectnonstd=false",
			Help:     "reject non-standard transactions regardless of the default settings for the active network",
			Variable: true,
		},
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
			log.Tracef.Print("setting debug level %s", dl)
			log.SetLevel(dl)
			for i := range logger.Levels {
				logger.Levels[i] = dl
			}
		}
		log.Debugf.Print("node version %s", n.Version())
		if ctx.Is("version") {
			fmt.Println("node version", n.Version())
			clog.Shutdown()
		}
		log.Trace.Print("running command")

		var cfgFile string
		if cfgFile, ok = ctx.Get("configfile"); !ok {
			cfgFile = n.DefaultConfigFile
		}
		if ctx.Is("init") {
			log.Debugf.Print("writing default configuration to %s", cfgFile)
			writeDefaultConfig(cfgFile)
			configNode(&ctx, cfgFile)
		} else {
			log.Infof.Print("loading configuration from %s", cfgFile)
			if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
				log.Warn.Print("configuration file does not exist, creating new one")
				writeDefaultConfig(cfgFile)
				configNode(&ctx, cfgFile)
			} else {
				log.Debug.Print("reading app configuration from", cfgFile)
				cfgData, err := ioutil.ReadFile(cfgFile)
				if err != nil {
					log.Error.Print(err.Error())
					clog.Shutdown()
				}
				log.Tracef.Print("parsing app configuration\n%s", cfgData)
				err = json.Unmarshal(cfgData, CombinedCfg)
				if err != nil {
					log.Error.Print(err.Error())
					clog.Shutdown()
				}
				configNode(&ctx, cfgFile)
			}
		}
		runNode()
		clog.Shutdown()
		return 0
	},
}

func configNode(ctx *climax.Context, cfgFile string) {
	var err error
	// Apply all configurations specified on commandline
	if ctx.Is("datadir") {
		r, _ := ctx.Get("datadir")
		Config.DataDir = r
	}
	if ctx.Is("addpeers") {
		r, _ := ctx.Get("addpeers")
		Config.AddPeers = strings.Split(r, " ")
	}
	if ctx.Is("connectpeers") {
		r, _ := ctx.Get("connectpeers")
		Config.ConnectPeers = strings.Split(r, " ")
	}
	if ctx.Is("disablelisten") {
		r, _ := ctx.Get("disablelisten")
		Config.DisableListen = r == "true"
	}
	if ctx.Is("listeners") {
		r, _ := ctx.Get("listeners")
		Config.Listeners = strings.Split(r, " ")
	}
	if ctx.Is("maxpeers") {
		r, _ := ctx.Get("maxpeers")
		Config.MaxPeers, err = strconv.Atoi(r)
		if err != nil {
			log.Error.Print(err.Error())
		}
	}
	if ctx.Is("disablebanning") {
		r, _ := ctx.Get("disablebanning")
		Config.DisableBanning = r == "true"
	}
	if ctx.Is("banduration") {
		r, _ := ctx.Get("banduration")
		error := false
		var bd time.Duration
		switch r[len(r)-1] {
		case 's':
			ts, err := strconv.Atoi(r[:len(r)-1])
			error = err != nil
			bd = time.Duration(ts) * time.Second
		case 'm':
			tm, err := strconv.Atoi(r[:len(r)-1])
			error = err != nil
			bd = time.Duration(tm) * time.Minute
		case 'h':
			th, err := strconv.Atoi(r[:len(r)-1])
			error = err != nil
			bd = time.Duration(th) * time.Hour
		case 'd':
			td, err := strconv.Atoi(r[:len(r)-1])
			error = err != nil
			bd = time.Duration(td) * 24 * time.Hour
		}
		if error {
			log.Errorf.Print("malformed banduration `%s` leaving set at `%s` err: %s", r, Config.BanDuration, err.Error())
		}
		Config.BanDuration = bd
	}
	if ctx.Is("banthreshold") {
		r, _ := ctx.Get("banthreshold")
		bt, err := strconv.Atoi(r)
		if err != nil {
			log.Errorf.Print("malformed banthreshold `%s` leaving set at `%s` err: %s", r, Config.BanThreshold, err.Error())
		} else {
			Config.BanThreshold = uint32(bt)
		}
	}
	if ctx.Is("whitelists") {
		r, _ := ctx.Get("whitelists")
		Config.Whitelists = strings.Split(r, " ")
	}
	if ctx.Is("rpcuser") {
		r, _ := ctx.Get("rpcuser")
		Config.RPCUser = r
	}
	if ctx.Is("rpcpass") {
		r, _ := ctx.Get("rpcpass")
		Config.RPCPass = r
	}
	if ctx.Is("rpclimituser") {
		r, _ := ctx.Get("rpclimituser")
		Config.RPCLimitUser = r
	}
	if ctx.Is("rpclimitpass") {
		r, _ := ctx.Get("rpclimitpass")
		Config.RPCLimitPass = r
	}
	if ctx.Is("rpclisteners") {
		r, _ := ctx.Get("rpclisteners")
		Config.RPCListeners = strings.Split(r, " ")
	}
	if ctx.Is("rpccert") {
		r, _ := ctx.Get("rpccert")
		Config.RPCCert = r
	}
	if ctx.Is("rpckey") {
		r, _ := ctx.Get("rpckey")
		Config.RPCKey = r
	}
	if ctx.Is("tls") {
		r, _ := ctx.Get("tls")
		Config.TLS = r == "true"
	}
	if ctx.Is("disablednsseed") {
		r, _ := ctx.Get("disablednsseed")
		Config.DisableDNSSeed = r == "true"
	}
	if ctx.Is("externalips") {
		r, _ := ctx.Get("externalips")
		Config.ExternalIPs = strings.Split(r, " ")
	}
	if ctx.Is("proxy") {
		r, _ := ctx.Get("proxy")
		Config.Proxy = r
	}
	if ctx.Is("proxyuser") {
		r, _ := ctx.Get("proxyuser")
		Config.ProxyUser = r
	}
	if ctx.Is("proxypass") {
		r, _ := ctx.Get("proxypass")
		Config.ProxyPass = r
	}
	if ctx.Is("onion") {
		r, _ := ctx.Get("onion")
		Config.OnionProxy = r
	}
	if ctx.Is("onionuser") {
		r, _ := ctx.Get("onionuser")
		Config.OnionProxyUser = r
	}
	if ctx.Is("onionpass") {
		r, _ := ctx.Get("onionpass")
		Config.OnionProxyPass = r
	}
	if ctx.Is("noonion") {
		r, _ := ctx.Get("noonion")
		Config.NoOnion = r == "true"
	}
	if ctx.Is("torisolation") {
		r, _ := ctx.Get("torisolation")
		Config.TorIsolation = r == "true"
	}
	if ctx.Is("network") {
		r, _ := ctx.Get("network")
		switch r {
		case "testnet":
			Config.TestNet3, Config.RegressionTest, Config.SimNet = true, false, false
		case "regtest":
			Config.TestNet3, Config.RegressionTest, Config.SimNet = false, true, false
		case "simnet":
			Config.TestNet3, Config.RegressionTest, Config.SimNet = false, false, true
		default:
			Config.TestNet3, Config.RegressionTest, Config.SimNet = false, false, false
		}
	}
	if ctx.Is("addcheckpoints") {
		r, _ := ctx.Get("")
		Config.AddCheckpoints = strings.Split(r, " ")
	}
	if ctx.Is("disablecheckpoints") {
		r, _ := ctx.Get("disablecheckpoints")
		Config.DisableCheckpoints = r == "true"
	}
	if ctx.Is("dbtype") {
		r, _ := ctx.Get("dbtype")
		Config.DbType = r
	}
	if ctx.Is("profile") {
		r, _ := ctx.Get("profile")
		Config.Profile = r
	}
	if ctx.Is("cpuprofile") {
		r, _ := ctx.Get("cpuprofile")
		Config.CPUProfile = r
	}
	if ctx.Is("upnp") {
		r, _ := ctx.Get("upnp")
		Config.Upnp = r == "true"
	}
	if ctx.Is("minrelaytxfee") {
		r, _ := ctx.Get("minrelaytxfee")
		_, err := fmt.Sscanf(r, "%0.f", Config.MinRelayTxFee)
		if err != nil {
			log.Errorf.Print("malformed minrelaytxfee: `%s` leaving set at `%0.f`",
				r, Config.MinRelayTxFee)
		}
	}
	if ctx.Is("freetxrelaylimit") {
		r, _ := ctx.Get("freetxrelaylimit")
		_, err = fmt.Sscanf(r, "%d", Config.FreeTxRelayLimit)
		if err != nil {
			log.Errorf.Print("malformed freetxrelaylimit: `%s` leaving set at `%d`",
				r, Config.FreeTxRelayLimit)
		}
	}
	if ctx.Is("norelaypriority") {
		r, _ := ctx.Get("norelaypriority")
		Config.NoRelayPriority = r == "true"
	}
	if ctx.Is("trickleinterval") {
		r, _ := ctx.Get("trickleinterval")
		error := false
		var ti time.Duration
		switch r[len(r)-1] {
		case 's':
			ts, err := strconv.Atoi(r[:len(r)-1])
			error = err != nil
			ti = time.Duration(ts) * time.Second
		case 'm':
			tm, err := strconv.Atoi(r[:len(r)-1])
			error = err != nil
			ti = time.Duration(tm) * time.Minute
		case 'h':
			th, err := strconv.Atoi(r[:len(r)-1])
			error = err != nil
			ti = time.Duration(th) * time.Hour
		case 'd':
			td, err := strconv.Atoi(r[:len(r)-1])
			error = err != nil
			ti = time.Duration(td) * 24 * time.Hour
		}
		if error {
			log.Errorf.Print("malformed trickleinterval `%s` leaving set at `%s` err: %s", r, Config.TrickleInterval, err.Error())
		}
		Config.TrickleInterval = ti
	}
	if ctx.Is("maxorphantxs") {
		r, _ := ctx.Get("maxorphantxs")
		mot, err := strconv.Atoi(r)
		if err != nil {
			log.Errorf.Print("malformed maxorphantxs: `%s` leaving set at `%d`",
				r, Config.MaxOrphanTxs)
		} else {
			Config.MaxOrphanTxs = mot
		}
	}
	if ctx.Is("algo") {
		r, _ := ctx.Get("algo")
		Config.Algo = r
	}
	if ctx.Is("generate") {
		r, _ := ctx.Get("generate")
		Config.Generate = r == "true"
	}
	if ctx.Is("genthreads") {
		r, _ := ctx.Get("genthreads")
		gt, err := strconv.Atoi(r)
		if err != nil {
			log.Errorf.Print("malformed freetxrelaylimit: `%s` leaving set at `%d`",
				r, Config.GenThreads)
		} else {
			Config.GenThreads = int32(gt)
		}
	}
	if ctx.Is("miningaddrs") {
		r, _ := ctx.Get("miningaddrs")
		Config.MiningAddrs = strings.Split(r, " ")
	}
	if ctx.Is("minerlistener") {
		r, _ := ctx.Get("minerlistener")
		Config.MinerListener = r
	}
	if ctx.Is("minerpass") {
		r, _ := ctx.Get("minerpass")
		Config.MinerPass = r
	}
	if ctx.Is("blockminsize") {
		r, _ := ctx.Get("blockminsize")
		bms, err := strconv.Atoi(r)
		if err != nil {
			log.Errorf.Print("malformed blockminsize: `%s` leaving set at `%d`",
				r, Config.BlockMinSize)
		} else {
			Config.BlockMinSize = uint32(bms)
		}
	}
	if ctx.Is("blockmaxsize") {
		r, _ := ctx.Get("blockmaxsize")
		bms, err := strconv.Atoi(r)
		if err != nil {
			log.Errorf.Print("malformed blockmaxsize: `%s` leaving set at `%d`",
				r, Config.BlockMaxSize)
		} else {
			Config.BlockMaxSize = uint32(bms)
		}
	}
	if ctx.Is("blockminweight") {
		r, _ := ctx.Get("blockminweight")
		bmw, err := strconv.Atoi(r)
		if err != nil {
			log.Errorf.Print("malformed blockminweight: `%s` leaving set at `%d`",
				r, Config.BlockMinWeight)
		} else {
			Config.BlockMinWeight = uint32(bmw)
		}
	}
	if ctx.Is("blockmaxweight") {
		r, _ := ctx.Get("blockmaxweight")
		bmw, err := strconv.Atoi(r)
		if err != nil {
			log.Errorf.Print("malformed blockmaxweight: `%s` leaving set at `%d`",
				r, Config.BlockMaxWeight)
		} else {
			Config.BlockMaxWeight = uint32(bmw)
		}
	}
	if ctx.Is("blockprioritysize") {
		r, _ := ctx.Get("blockprioritysize")
		bps, err := strconv.Atoi(r)
		if err != nil {
			log.Errorf.Print("malformed blockprioritysize: `%s` leaving set at `%d`",
				r, Config.BlockPrioritySize)
		} else {
			Config.BlockPrioritySize = uint32(bps)
		}
	}
	if ctx.Is("uacomment") {
		r, _ := ctx.Get("uacomment")
		Config.UserAgentComments = strings.Split(r, " ")
	}
	if ctx.Is("nopeerbloomfilters") {
		r, _ := ctx.Get("nopeerbloomfilters")
		Config.NoPeerBloomFilters = r == "true"
	}
	if ctx.Is("nocfilters") {
		r, _ := ctx.Get("nocfilters")
		Config.NoCFilters = r == "true"
	}
	if ctx.Is("dropcfindex") {
		Config.DropCfIndex = true
	}
	if ctx.Is("sigcachemaxsize") {
		r, _ := ctx.Get("sigcachemaxsize")
		sms, err := strconv.Atoi(r)
		if err != nil || sms < 0 {
			log.Errorf.Print("malformed sigcachemaxsize: `%s` leaving set at `%d`",
				r, Config.SigCacheMaxSize)
		} else {
			Config.SigCacheMaxSize = uint(sms)
		}
	}
	if ctx.Is("blocksonly") {
		r, _ := ctx.Get("blocksonly")
		Config.BlocksOnly = r == "true"
	}
	if ctx.Is("txindex") {
		r, _ := ctx.Get("txindex")
		Config.TxIndex = r == "true"
	}
	if ctx.Is("droptxindex") {
		Config.DropTxIndex = true
	}
	if ctx.Is("addrindex") {
		r, _ := ctx.Get("addrindex")
		Config.AddrIndex = r == "true"
	}
	if ctx.Is("dropaddrindex") {
		Config.DropAddrIndex = true
	}
	if ctx.Is("relaynonstd") {
		r, _ := ctx.Get("relaynonstd")
		Config.RelayNonStd = r == "true"
	}
	if ctx.Is("rejectnonstd") {
		r, _ := ctx.Get("rejectnonstd")
		Config.RejectNonStd = r == "true"
	}
	logger.SetLogging(ctx)
	if ctx.Is("save") {
		log.Infof.Print("saving config file to %s", cfgFile)
		j, err := json.MarshalIndent(CombinedCfg, "", "  ")
		if err != nil {
			log.Error.Print(err.Error())
		}
		j = append(j, '\n')
		log.Tracef.Print("JSON formatted config file\n%s", j)
		ioutil.WriteFile(cfgFile, j, 0600)
	}
}

func writeDefaultConfig(cfgFile string) {
	defCfg := defaultConfig()
	defCfg.Node.ConfigFile = cfgFile
	j, err := json.MarshalIndent(defCfg, "", "  ")
	if err != nil {
		log.Error.Print(err.Error())
	}
	j = append(j, '\n')
	log.Tracef.Print("JSON formatted config file\n%s", j)
	ioutil.WriteFile(cfgFile, j, 0600)
	// if we are writing default config we also want to use it
	CombinedCfg = *defCfg
}

func defaultConfig() *ConfigAndLog {
	return &ConfigAndLog{
		Node: &n.Config{
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
		Levels: logger.GetDefault(),
	}
}
