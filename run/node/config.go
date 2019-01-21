package node

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

	"git.parallelcoin.io/pod/lib/blockchain"
	cl "git.parallelcoin.io/pod/lib/clog"
	"git.parallelcoin.io/pod/lib/connmgr"
	"git.parallelcoin.io/pod/lib/fork"
	"git.parallelcoin.io/pod/lib/util"
	n "git.parallelcoin.io/pod/module/node"
	"git.parallelcoin.io/pod/module/node/mempool"
	"git.parallelcoin.io/pod/run/logger"
	"git.parallelcoin.io/pod/run/util"
	"github.com/btcsuite/go-socks/socks"
	"github.com/davecgh/go-spew/spew"
	"github.com/tucnak/climax"
)

// Log is thte main node logger
var Log = cl.NewSubSystem("run/node", "info")
var log = Log.Ch

// serviceOptions defines the configuration options for the daemon as a service on Windows.
type serviceOptions struct {
	ServiceCommand string `short:"s" long:"service" description:"Service command {install, remove, start, stop}"`
}

// StateCfg is a reference to the main node state configuration struct
var StateCfg = n.StateCfg

// minUint32 is a helper function to return the minimum of two uint32s. This avoids a math import and the need to cast to floats.
func minUint32(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}

// runServiceCommand is only set to a real function on Windows.  It is used to parse and execute service commands specified via the -s flag.
var runServiceCommand func(string) error

var aN = filepath.Base(os.Args[0])
var appName = strings.TrimSuffix(aN, filepath.Ext(aN))

var usageMessage = fmt.Sprintf("use `%s help node` to show usage", appName)

// Cfg is the combined app and logging configuration data
type Cfg struct {
	Node      *n.Config
	LogLevels map[string]string
}

// Config is the combined app and log levels configuration
var Config = DefaultConfig()

// Command is a command to send RPC queries to bitcoin RPC protocol server for node and wallet queries
var Command = climax.Command{
	Name:  "node",
	Brief: "parallelcoin full node",
	Help:  "distrubutes, verifies and mines blocks for the parallelcoin duo cryptocurrency, as well as optionally providing search indexes for transactions in the database",
	Flags: []climax.Flag{

		podutil.GenerateFlag(`version`, `V`, `--version`, `show version number and quit`, false),

		podutil.GenerateFlag(`configfile`, `C`, `--configfile=/path/to/conf`, `path to configuration file`, true),
		podutil.GenerateFlag(`datadir`, `D`, `--configfile=/path/to/conf`, `path to configuration file`, true),

		podutil.GenerateFlag(`init`, ``, `--init`, `resets configuration to defaults`, false),
		podutil.GenerateFlag(`save`, ``, `--save`, `saves current configuration`, false),

		podutil.GenerateFlag(`dropcfindex`, ``, `--dropcfindex`, `deletes the index used for committed filtering (CF) support from the database on start up and then exits`, false),
		podutil.GenerateFlag(`droptxindex`, ``, `--droptxindex`, `deletes the hash-based transaction index from the database on start up and then exits.`, false),
		podutil.GenerateFlag(`dropaddrindex`, ``, `--dropaddrindex`, `deletes the address-based transaction index from the database on start up and then exits`, false),

		podutil.GenerateFlag(`addpeers`, ``, `--addpeers=some.peer.com:11047`, `adds a peer to the peers database to try to connect to`, true),
		podutil.GenerateFlag(`connectpeers`, ``, `--connectpeers=some.peer.com:11047`, `adds a peer to a connect-only whitelist`, true),
		podutil.GenerateFlag(`disablelisten`, ``, `--disablelisten=true`, `disables the P2P listener`, true),
		podutil.GenerateFlag(`listeners`, `S`, `--listeners=127.0.0.1:11047`, `sets an address to listen for P2P connections`, true),
		podutil.GenerateFlag(`maxpeers`, ``, `--maxpeers=100`, `sets max number of peers to open connections to at once`, true),
		podutil.GenerateFlag(`disablebanning`, ``, `--disablebanning`, `disable banning of misbehaving peers`, true),
		podutil.GenerateFlag(`banduration`, ``, `--banduration=1h`, `how long to ban misbehaving peers - valid time units are {s, m, h},  minimum 1s`, true),
		podutil.GenerateFlag(`banthreshold`, ``, `--banthreshold=100`, `maximum allowed ban score before disconnecting and banning misbehaving peers`, true),
		podutil.GenerateFlag(`whitelists`, ``, `--whitelists=127.0.0.1:11047`, `add an IP network or IP that will not be banned - eg. 192.168.1.0/24 or ::1`, true),
		podutil.GenerateFlag(`rpcuser`, `u`, `--rpcuser=username`, `RPC username`, true),
		podutil.GenerateFlag(`rpcpass`, `P`, `--rpcpass=password`, `RPC password`, true),
		podutil.GenerateFlag(`rpclimituser`, `u`, `--rpclimituser=username`, `limited user RPC username`, true),
		podutil.GenerateFlag(`rpclimitpass`, `P`, `--rpclimitpass=password`, `limited user RPC password`, true),
		podutil.GenerateFlag(`rpclisteners`, `s`, `--rpclisteners=127.0.0.1:11048`, `RPC server to connect to`, true),
		podutil.GenerateFlag(`rpccert`, `c`, `--rpccert=/path/to/rpn.cert`, `RPC server tls certificate chain for validation`, true),
		podutil.GenerateFlag(`rpckey`, `c`, `--rpccert=/path/to/rpn.key`, `RPC server tls key for validation`, true),
		podutil.GenerateFlag(`tls`, ``, `--tls=false`, `enable TLS`, true),
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

		podutil.GenerateFlag(`debuglevel`, `d`, `--debuglevel=trace`, `sets debuglevel, default info, sets the baseline for others not specified`, true),

		podutil.GenerateFlag("lib-addrmgr", "", "--lib-addrmg=info", "", true),
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
			writeDefaultConfig(cfgFile)
		} else {
			log <- cl.Infof{
				"loading configuration from %s",
				cfgFile,
			}
			if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
				log <- cl.Warn{"configuration file does not exist, creating new one"}
				writeDefaultConfig(cfgFile)
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
				log <- cl.Tracef{
					"parsing app configuration\n%s",
					cfgData,
				}
				err = json.Unmarshal(cfgData, &Config)
				if err != nil {
					log <- cl.Error{
						"parsing app config file:",
						err.Error(),
					}
					writeDefaultConfig(cfgFile)
				}
			}
		}
		configNode(&ctx, cfgFile)
		runNode()
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
	cfg := Config.Node
	var err error
	var r *string
	t := ""
	r = &t
	if getIfIs(ctx, "debuglevel", r) {
		switch *r {
		case "fatal", "error", "warn", "info", "debug", "trace":
			cfg.DebugLevel = *r
		default:
			cfg.DebugLevel = "info"
		}
		Log.SetLevel(cfg.DebugLevel)
	}
	if getIfIs(ctx, "datadir", r) {
		cfg.DataDir = n.CleanAndExpandPath(*r)
	}
	if getIfIs(ctx, "addpeers", r) {
		podutil.NormalizeAddresses(*r, n.DefaultPort, &cfg.AddPeers)
	}
	if getIfIs(ctx, "connectpeers", r) {
		podutil.NormalizeAddresses(*r, n.DefaultPort, &cfg.ConnectPeers)
	}
	if getIfIs(ctx, "disablelisten", r) {
		cfg.DisableListen = *r == "true"
	}
	if getIfIs(ctx, "listeners", r) {
		podutil.NormalizeAddresses(*r, n.DefaultPort, &cfg.Listeners)
	}
	if getIfIs(ctx, "maxpeers", r) {
		if err := podutil.ParseInteger(*r, "maxpeers", &cfg.MaxPeers); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if getIfIs(ctx, "disablebanning", r) {
		cfg.DisableBanning = *r == "true"
	}
	if getIfIs(ctx, "banduration", r) {
		if err := podutil.ParseDuration(*r, "banduration", &cfg.BanDuration); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if getIfIs(ctx, "banthreshold", r) {
		var bt int
		if err := podutil.ParseInteger(*r, "banthtreshold", &bt); err != nil {
			log <- cl.Wrn(err.Error())
		} else {
			cfg.BanThreshold = uint32(bt)
		}
	}
	if getIfIs(ctx, "whitelists", r) {
		podutil.NormalizeAddresses(*r, n.DefaultPort, &cfg.Whitelists)
	}
	if getIfIs(ctx, "rpcuser", r) {
		cfg.RPCUser = *r
	}
	if getIfIs(ctx, "rpcpass", r) {
		cfg.RPCPass = *r
	}
	if getIfIs(ctx, "rpclimituser", r) {
		cfg.RPCLimitUser = *r
	}
	if getIfIs(ctx, "rpclimitpass", r) {
		cfg.RPCLimitPass = *r
	}
	if getIfIs(ctx, "rpclisteners", r) {
		podutil.NormalizeAddresses(*r, n.DefaultRPCPort, &cfg.RPCListeners)
	}
	if getIfIs(ctx, "rpccert", r) {
		cfg.RPCCert = n.CleanAndExpandPath(*r)
	}
	if getIfIs(ctx, "rpckey", r) {
		cfg.RPCKey = n.CleanAndExpandPath(*r)
	}
	if getIfIs(ctx, "tls", r) {
		cfg.TLS = *r == "true"
	}
	if getIfIs(ctx, "disablednsseed", r) {
		cfg.DisableDNSSeed = *r == "true"
	}
	if getIfIs(ctx, "externalips", r) {
		podutil.NormalizeAddresses(*r, n.DefaultPort, &cfg.ExternalIPs)
	}
	if getIfIs(ctx, "proxy", r) {
		podutil.NormalizeAddress(*r, "9050", &cfg.Proxy)
	}
	if getIfIs(ctx, "proxyuser", r) {
		cfg.ProxyUser = *r
	}
	if getIfIs(ctx, "proxypass", r) {
		cfg.ProxyPass = *r
	}
	if getIfIs(ctx, "onion", r) {
		podutil.NormalizeAddress(*r, "9050", &cfg.OnionProxy)
	}
	if getIfIs(ctx, "onionuser", r) {
		cfg.OnionProxyUser = *r
	}
	if getIfIs(ctx, "onionpass", r) {
		cfg.OnionProxyPass = *r
	}
	if getIfIs(ctx, "noonion", r) {
		cfg.NoOnion = *r == "true"
	}
	if getIfIs(ctx, "torisolation", r) {
		cfg.TorIsolation = *r == "true"
	}
	if getIfIs(ctx, "network", r) {
		switch *r {
		case "testnet":
			cfg.TestNet3, cfg.RegressionTest, cfg.SimNet = true, false, false
		case "regtest":
			cfg.TestNet3, cfg.RegressionTest, cfg.SimNet = false, true, false
		case "simnet":
			cfg.TestNet3, cfg.RegressionTest, cfg.SimNet = false, false, true
		default:
			cfg.TestNet3, cfg.RegressionTest, cfg.SimNet = false, false, false
		}
	}
	if getIfIs(ctx, "addcheckpoints", r) {
		cfg.AddCheckpoints = strings.Split(*r, " ")
	}
	if getIfIs(ctx, "disablecheckpoints", r) {
		cfg.DisableCheckpoints = *r == "true"
	}
	if getIfIs(ctx, "dbtype", r) {
		cfg.DbType = *r
	}
	if getIfIs(ctx, "profile", r) {
		var p int
		if err = podutil.ParseInteger(*r, "profile", &p); err == nil {
			cfg.Profile = fmt.Sprint(p)
		}
	}
	if getIfIs(ctx, "cpuprofile", r) {
		cfg.CPUProfile = *r
	}
	if getIfIs(ctx, "upnp", r) {
		cfg.Upnp = *r == "true"
	}
	if getIfIs(ctx, "minrelaytxfee", r) {
		if err := podutil.ParseFloat(*r, "minrelaytxfee", &cfg.MinRelayTxFee); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if getIfIs(ctx, "freetxrelaylimit", r) {
		if err := podutil.ParseFloat(*r, "freetxrelaylimit", &cfg.FreeTxRelayLimit); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if getIfIs(ctx, "norelaypriority", r) {
		cfg.NoRelayPriority = *r == "true"
	}
	if getIfIs(ctx, "trickleinterval", r) {
		if err := podutil.ParseDuration(*r, "trickleinterval", &cfg.TrickleInterval); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if getIfIs(ctx, "maxorphantxs", r) {
		if err := podutil.ParseInteger(*r, "maxorphantxs", &cfg.MaxOrphanTxs); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if getIfIs(ctx, "algo", r) {
		cfg.Algo = *r
	}
	if getIfIs(ctx, "generate", r) {
		cfg.Generate = *r == "true"
	}
	if getIfIs(ctx, "genthreads", r) {
		var gt int
		if err := podutil.ParseInteger(*r, "genthreads", &gt); err != nil {
			log <- cl.Wrn(err.Error())
		} else {
			cfg.GenThreads = int32(gt)
		}
	}
	if getIfIs(ctx, "miningaddrs", r) {
		cfg.MiningAddrs = strings.Split(*r, " ")
	}
	if getIfIs(ctx, "minerlistener", r) {
		podutil.NormalizeAddress(*r, n.DefaultRPCPort, &cfg.MinerListener)
	}
	if getIfIs(ctx, "minerpass", r) {
		cfg.MinerPass = *r
	}
	if getIfIs(ctx, "blockminsize", r) {
		if err := podutil.ParseUint32(*r, "blockminsize", &cfg.BlockMinSize); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if getIfIs(ctx, "blockmaxsize", r) {
		if err := podutil.ParseUint32(*r, "blockmaxsize", &cfg.BlockMaxSize); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if getIfIs(ctx, "blockminweight", r) {
		if err := podutil.ParseUint32(*r, "blockminweight", &cfg.BlockMinWeight); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if getIfIs(ctx, "blockmaxweight", r) {
		if err := podutil.ParseUint32(*r, "blockmaxweight", &cfg.BlockMaxWeight); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if getIfIs(ctx, "blockprioritysize", r) {
		if err := podutil.ParseUint32(*r, "blockmaxweight", &cfg.BlockPrioritySize); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if getIfIs(ctx, "uacomment", r) {
		cfg.UserAgentComments = strings.Split(*r, " ")
	}
	if getIfIs(ctx, "nopeerbloomfilters", r) {
		cfg.NoPeerBloomFilters = *r == "true"
	}
	if getIfIs(ctx, "nocfilters", r) {
		cfg.NoCFilters = *r == "true"
	}
	if ctx.Is("dropcfindex") {
		cfg.DropCfIndex = true
	}
	if getIfIs(ctx, "sigcachemaxsize", r) {
		var scms int
		if err := podutil.ParseInteger(*r, "sigcachemaxsize", &scms); err != nil {
			log <- cl.Wrn(err.Error())
		} else {
			cfg.SigCacheMaxSize = uint(scms)
		}
	}
	if getIfIs(ctx, "blocksonly", r) {
		cfg.BlocksOnly = *r == "true"
	}
	if getIfIs(ctx, "txindex", r) {
		cfg.TxIndex = *r == "true"
	}
	if ctx.Is("droptxindex") {
		cfg.DropTxIndex = true
	}
	if ctx.Is("addrindex") {
		r, _ := ctx.Get("addrindex")
		cfg.AddrIndex = r == "true"
	}
	if ctx.Is("dropaddrindex") {
		cfg.DropAddrIndex = true
	}
	if getIfIs(ctx, "relaynonstd", r) {
		cfg.RelayNonStd = *r == "true"
	}
	if getIfIs(ctx, "rejectnonstd", r) {
		cfg.RejectNonStd = *r == "true"
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
	if cfg.RegressionTest && len(cfg.AddPeers) > 0 {
		cfg.AddPeers = nil
	}
	// Set the mining algorithm correctly, default to random if unrecognised
	switch cfg.Algo {
	case "blake14lr", "cryptonight7v2", "keccak", "lyra2rev2", "scrypt", "skein", "x11", "stribog", "random", "easy":
	default:
		cfg.Algo = "random"
	}
	relayNonStd := ActiveNetParams.RelayNonStdTxs
	funcName := "loadConfig"
	switch {
	case cfg.RelayNonStd && cfg.RejectNonStd:
		str := "%s: rejectnonstd and relaynonstd cannot be used together -- choose only one"
		err := fmt.Errorf(str, funcName)
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	case cfg.RejectNonStd:
		relayNonStd = false
	case cfg.RelayNonStd:
		relayNonStd = true
	}
	cfg.RelayNonStd = relayNonStd
	// Append the network type to the data directory so it is "namespaced" per network.  In addition to the block database, there are other pieces of data that are saved to disk such as address manager state. All data is specific to a network, so namespacing the data directory means each individual piece of serialized data does not have to worry about changing names per network and such.
	cfg.DataDir = n.CleanAndExpandPath(cfg.DataDir)
	cfg.DataDir = filepath.Join(cfg.DataDir, netName(ActiveNetParams))
	// Append the network type to the log directory so it is "namespaced" per network in the same fashion as the data directory.
	cfg.LogDir = n.CleanAndExpandPath(cfg.LogDir)
	cfg.LogDir = filepath.Join(cfg.LogDir, netName(ActiveNetParams))

	// Initialize log rotation.  After log rotation has been initialized, the logger variables may be used.
	// initLogRotator(filepath.Join(cfg.LogDir, DefaultLogFilename))
	// Validate database type.
	if !n.ValidDbType(cfg.DbType) {
		str := "%s: The specified database type [%v] is invalid -- " +
			"supported types %v"
		err := fmt.Errorf(str, funcName, cfg.DbType, n.KnownDbTypes)
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// Validate profile port number
	if cfg.Profile != "" {
		profilePort, err := strconv.Atoi(cfg.Profile)
		if err != nil || profilePort < 1024 || profilePort > 65535 {
			str := "%s: The profile port must be between 1024 and 65535"
			err := fmt.Errorf(str, funcName)
			fmt.Fprintln(os.Stderr, err)
			fmt.Fprintln(os.Stderr, usageMessage)
			cl.Shutdown()
		}
	}
	// Don't allow ban durations that are too short.
	if cfg.BanDuration < time.Second {
		str := "%s: The banduration option may not be less than 1s -- parsed [%v]"
		err := fmt.Errorf(str, funcName, cfg.BanDuration)
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// Validate any given whitelisted IP addresses and networks.
	if len(cfg.Whitelists) > 0 {
		var ip net.IP
		StateCfg.ActiveWhitelists = make([]*net.IPNet, 0, len(cfg.Whitelists))
		for _, addr := range cfg.Whitelists {
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
	if len(cfg.AddPeers) > 0 && len(cfg.ConnectPeers) > 0 {
		str := "%s: the --addpeer and --connect options can not be " +
			"mixed"
		err := fmt.Errorf(str, funcName)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
	}
	// --proxy or --connect without --listen disables listening.
	if (cfg.Proxy != "" || len(cfg.ConnectPeers) > 0) &&
		len(cfg.Listeners) == 0 {
		cfg.DisableListen = true
	}
	// Connect means no DNS seeding.
	if len(cfg.ConnectPeers) > 0 {
		cfg.DisableDNSSeed = true
	}
	// Add the default listener if none were specified. The default listener is all addresses on the listen port for the network we are to connect to.
	if len(cfg.Listeners) == 0 {
		cfg.Listeners = []string{
			net.JoinHostPort("", ActiveNetParams.DefaultPort),
		}
	}
	// Check to make sure limited and admin users don't have the same username
	if cfg.RPCUser == cfg.RPCLimitUser && cfg.RPCUser != "" {
		str := "%s: --rpcuser and --rpclimituser must not specify the same username"
		err := fmt.Errorf(str, funcName)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// Check to make sure limited and admin users don't have the same password
	if cfg.RPCPass == cfg.RPCLimitPass && cfg.RPCPass != "" {
		str := "%s: --rpcpass and --rpclimitpass must not specify the " +
			"same password"
		err := fmt.Errorf(str, funcName)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// The RPC server is disabled if no username or password is provided.
	if (cfg.RPCUser == "" || cfg.RPCPass == "") &&
		(cfg.RPCLimitUser == "" || cfg.RPCLimitPass == "") {
		cfg.DisableRPC = true
	}
	if cfg.DisableRPC {
		log <- cl.Inf("RPC service is disabled")
	}
	// Default RPC to listen on localhost only.
	if !cfg.DisableRPC && len(cfg.RPCListeners) == 0 {
		addrs, err := net.LookupHost(n.DefaultRPCListener)
		if err != nil {
			log <- cl.Err(err.Error())
			cl.Shutdown()
		}
		cfg.RPCListeners = make([]string, 0, len(addrs))
		for _, addr := range addrs {
			addr = net.JoinHostPort(addr, ActiveNetParams.RPCPort)
			cfg.RPCListeners = append(cfg.RPCListeners, addr)
		}
	}
	if cfg.RPCMaxConcurrentReqs < 0 {
		str := "%s: The rpcmaxwebsocketconcurrentrequests option may not be less than 0 -- parsed [%d]"
		err := fmt.Errorf(str, funcName, cfg.RPCMaxConcurrentReqs)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// Validate the the minrelaytxfee.
	StateCfg.ActiveMinRelayTxFee, err = util.NewAmount(cfg.MinRelayTxFee)
	if err != nil {
		str := "%s: invalid minrelaytxfee: %v"
		err := fmt.Errorf(str, funcName, err)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// Limit the max block size to a sane value.
	if cfg.BlockMaxSize < n.BlockMaxSizeMin || cfg.BlockMaxSize >
		n.BlockMaxSizeMax {
		str := "%s: The blockmaxsize option must be in between %d and %d -- parsed [%d]"
		err := fmt.Errorf(str, funcName, n.BlockMaxSizeMin,
			n.BlockMaxSizeMax, cfg.BlockMaxSize)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// Limit the max block weight to a sane value.
	if cfg.BlockMaxWeight < n.BlockMaxWeightMin ||
		cfg.BlockMaxWeight > n.BlockMaxWeightMax {
		str := "%s: The blockmaxweight option must be in between %d and %d -- parsed [%d]"
		err := fmt.Errorf(str, funcName, n.BlockMaxWeightMin,
			n.BlockMaxWeightMax, cfg.BlockMaxWeight)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// Limit the max orphan count to a sane vlue.
	if cfg.MaxOrphanTxs < 0 {
		str := "%s: The maxorphantx option may not be less than 0 -- parsed [%d]"
		err := fmt.Errorf(str, funcName, cfg.MaxOrphanTxs)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// Limit the block priority and minimum block sizes to max block size.
	cfg.BlockPrioritySize = minUint32(cfg.BlockPrioritySize, cfg.BlockMaxSize)
	cfg.BlockMinSize = minUint32(cfg.BlockMinSize, cfg.BlockMaxSize)
	cfg.BlockMinWeight = minUint32(cfg.BlockMinWeight, cfg.BlockMaxWeight)
	switch {
	// If the max block size isn't set, but the max weight is, then we'll set the limit for the max block size to a safe limit so weight takes precedence.
	case cfg.BlockMaxSize == n.DefaultBlockMaxSize &&
		cfg.BlockMaxWeight != n.DefaultBlockMaxWeight:
		cfg.BlockMaxSize = blockchain.MaxBlockBaseSize - 1000
	// If the max block weight isn't set, but the block size is, then we'll scale the set weight accordingly based on the max block size value.
	case cfg.BlockMaxSize != n.DefaultBlockMaxSize &&
		cfg.BlockMaxWeight == n.DefaultBlockMaxWeight:
		cfg.BlockMaxWeight = cfg.BlockMaxSize * blockchain.WitnessScaleFactor
	}
	// Look for illegal characters in the user agent comments.
	for _, uaComment := range cfg.UserAgentComments {
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
	if cfg.TxIndex && cfg.DropTxIndex {
		err := fmt.Errorf("%s: the --txindex and --droptxindex options may  not be activated at the same time",
			funcName)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()

	}
	// --addrindex and --dropaddrindex do not mix.
	if cfg.AddrIndex && cfg.DropAddrIndex {
		err := fmt.Errorf("%s: the --addrindex and --dropaddrindex "+
			"options may not be activated at the same time",
			funcName)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// --addrindex and --droptxindex do not mix.
	if cfg.AddrIndex && cfg.DropTxIndex {
		err := fmt.Errorf("%s: the --addrindex and --droptxindex options may not be activated at the same time "+
			"because the address index relies on the transaction index",
			funcName)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// Check mining addresses are valid and saved parsed versions.
	StateCfg.ActiveMiningAddrs = make([]util.Address, 0, len(cfg.MiningAddrs))
	for _, strAddr := range cfg.MiningAddrs {
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
	if (cfg.Generate || cfg.MinerListener != "") && len(cfg.MiningAddrs) == 0 {
		str := "%s: the generate flag is set, but there are no mining addresses specified "
		err := fmt.Errorf(str, funcName)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()

	}
	if cfg.MinerPass != "" {
		StateCfg.ActiveMinerKey = fork.Argon2i([]byte(cfg.MinerPass))
	}
	// Add default port to all listener addresses if needed and remove duplicate addresses.
	cfg.Listeners = n.NormalizeAddresses(cfg.Listeners,
		ActiveNetParams.DefaultPort)
	// Add default port to all rpc listener addresses if needed and remove duplicate addresses.
	cfg.RPCListeners = n.NormalizeAddresses(cfg.RPCListeners,
		ActiveNetParams.RPCPort)
	if !cfg.DisableRPC && !cfg.TLS {
		for _, addr := range cfg.RPCListeners {
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
	cfg.AddPeers = n.NormalizeAddresses(cfg.AddPeers,
		ActiveNetParams.DefaultPort)
	cfg.ConnectPeers = n.NormalizeAddresses(cfg.ConnectPeers,
		ActiveNetParams.DefaultPort)
	// --noonion and --onion do not mix.
	if cfg.NoOnion && cfg.OnionProxy != "" {
		err := fmt.Errorf("%s: the --noonion and --onion options may not be activated at the same time", funcName)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// Check the checkpoints for syntax errors.
	StateCfg.AddedCheckpoints, err = n.ParseCheckpoints(cfg.AddCheckpoints)
	if err != nil {
		str := "%s: Error parsing checkpoints: %v"
		err := fmt.Errorf(str, funcName, err)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// Tor stream isolation requires either proxy or onion proxy to be set.
	if cfg.TorIsolation && cfg.Proxy == "" && cfg.OnionProxy == "" {
		str := "%s: Tor stream isolation requires either proxy or onionproxy to be set"
		err := fmt.Errorf(str, funcName)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// Setup dial and DNS resolution (lookup) functions depending on the specified options.  The default is to use the standard net.DialTimeout function as well as the system DNS resolver.  When a proxy is specified, the dial function is set to the proxy specific dial function and the lookup is set to use tor (unless --noonion is specified in which case the system DNS resolver is used).
	StateCfg.Dial = net.DialTimeout
	StateCfg.Lookup = net.LookupIP
	if cfg.Proxy != "" {
		_, _, err := net.SplitHostPort(cfg.Proxy)
		if err != nil {
			str := "%s: Proxy address '%s' is invalid: %v"
			err := fmt.Errorf(str, funcName, cfg.Proxy, err)
			log <- cl.Err(err.Error())
			fmt.Fprintln(os.Stderr, usageMessage)
			cl.Shutdown()
		}
		// Tor isolation flag means proxy credentials will be overridden unless there is also an onion proxy configured in which case that one will be overridden.
		torIsolation := false
		if cfg.TorIsolation && cfg.OnionProxy == "" &&
			(cfg.ProxyUser != "" || cfg.ProxyPass != "") {
			torIsolation = true
			fmt.Fprintln(os.Stderr, "Tor isolation set -- "+
				"overriding specified proxy user credentials")
		}
		proxy := &socks.Proxy{
			Addr:         cfg.Proxy,
			Username:     cfg.ProxyUser,
			Password:     cfg.ProxyPass,
			TorIsolation: torIsolation,
		}
		StateCfg.Dial = proxy.DialTimeout
		// Treat the proxy as tor and perform DNS resolution through it unless the --noonion flag is set or there is an onion-specific proxy configured.
		if !cfg.NoOnion && cfg.OnionProxy == "" {
			StateCfg.Lookup = func(host string) ([]net.IP, error) {
				return connmgr.TorLookupIP(host, cfg.Proxy)
			}
		}
	}
	// Setup onion address dial function depending on the specified options. The default is to use the same dial function selected above.  However, when an onion-specific proxy is specified, the onion address dial function is set to use the onion-specific proxy while leaving the normal dial function as selected above.  This allows .onion address traffic to be routed through a different proxy than normal traffic.
	if cfg.OnionProxy != "" {
		_, _, err := net.SplitHostPort(cfg.OnionProxy)
		if err != nil {
			str := "%s: Onion proxy address '%s' is invalid: %v"
			err := fmt.Errorf(str, funcName, cfg.OnionProxy, err)
			log <- cl.Err(err.Error())
			fmt.Fprintln(os.Stderr, usageMessage)
			cl.Shutdown()
		}
		// Tor isolation flag means onion proxy credentials will be overridden.
		if cfg.TorIsolation &&
			(cfg.OnionProxyUser != "" || cfg.OnionProxyPass != "") {
			fmt.Fprintln(os.Stderr, "Tor isolation set -- "+
				"overriding specified onionproxy user "+
				"credentials ")
		}
		StateCfg.Oniondial = func(network, addr string, timeout time.Duration) (net.Conn, error) {
			proxy := &socks.Proxy{
				Addr:         cfg.OnionProxy,
				Username:     cfg.OnionProxyUser,
				Password:     cfg.OnionProxyPass,
				TorIsolation: cfg.TorIsolation,
			}
			return proxy.DialTimeout(network, addr, timeout)
		}
		// When configured in bridge mode (both --onion and --proxy are configured), it means that the proxy configured by --proxy is not a tor proxy, so override the DNS resolution to use the onion-specific proxy.
		if cfg.Proxy != "" {
			StateCfg.Lookup = func(host string) ([]net.IP, error) {
				return connmgr.TorLookupIP(host, cfg.OnionProxy)
			}
		}
	} else {
		StateCfg.Oniondial = StateCfg.Dial
	}
	// Specifying --noonion means the onion address dial function results in an error.
	if cfg.NoOnion {
		StateCfg.Oniondial = func(a, b string, t time.Duration) (net.Conn, error) {
			return nil, errors.New("tor has been disabled")
		}
	}
}

func writeDefaultConfig(cfgFile string) {
	defCfg := DefaultConfig()
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

// DefaultConfig is the default configuration for node
func DefaultConfig() *Cfg {
	user := podutil.GenerateKey()
	pass := podutil.GenerateKey()
	return &Cfg{
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
