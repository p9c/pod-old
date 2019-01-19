package node

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"git.parallelcoin.io/pod/lib/clog"
	n "git.parallelcoin.io/pod/module/node"
	"git.parallelcoin.io/pod/module/node/mempool"
	"git.parallelcoin.io/pod/run/logger"
	"git.parallelcoin.io/pod/run/util"
	"github.com/tucnak/climax"
)

// Log is thte main node logger
var Log = clog.NewSubSystem("pod/node", clog.Ninf)

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

		podutil.GenerateFlag(`log-database`, ``, `--log-database=info`, `sets log level for database`, true),
		podutil.GenerateFlag(`log-txscript`, ``, `--log-txscript=info`, `sets log level for txscript`, true),
		podutil.GenerateFlag(`log-peer`, ``, `--log-peer=info`, `sets log level for peer`, true),
		podutil.GenerateFlag(`log-netsync`, ``, `--log-netsync=info`, `sets log level for netsync`, true),
		podutil.GenerateFlag(`log-rpcclient`, ``, `--log-rpcclient=info`, `sets log level for rpcclient`, true),
		podutil.GenerateFlag(`addrmgr`, ``, `--log-addrmgr=info`, `sets log level for addrmgr`, true),
		podutil.GenerateFlag(`log-blockchain-indexers`, ``, `--log-blockchain-indexers=info`, `sets log level for blockchain-indexers`, true),
		podutil.GenerateFlag(`log-blockchain`, ``, `--log-blockchain=info`, `sets log level for blockchain`, true),
		podutil.GenerateFlag(`log-mining-cpuminer`, ``, `--log-mining-cpuminer=info`, `sets log level for mining-cpuminer`, true),
		podutil.GenerateFlag(`log-mining`, ``, `--log-mining=info`, `sets log level for mining`, true),
		podutil.GenerateFlag(`log-mining-controller`, ``, `--log-mining-controller=info`, `sets log level for mining-controller`, true),
		podutil.GenerateFlag(`log-connmgr`, ``, `--log-connmgr=info`, `sets log level for connmgr`, true),
		podutil.GenerateFlag(`log-spv`, ``, `--log-spv=info`, `sets log level for spv`, true),
		podutil.GenerateFlag(`log-node-mempool`, ``, `--log-node-mempool=info`, `sets log level for node-mempool`, true),
		podutil.GenerateFlag(`log-node`, ``, `--log-node=info`, `sets log level for node`, true),
		podutil.GenerateFlag(`log-wallet-wallet`, ``, `--log-wallet-wallet=info`, `sets log level for wallet-wallet`, true),
		podutil.GenerateFlag(`log-wallet-tx`, ``, `--log-wallet-tx=info`, `sets log level for wallet-tx`, true),
		podutil.GenerateFlag(`log-wallet-votingpool`, ``, `--log-wallet-votingpool=info`, `sets log level for wallet-votingpool`, true),
		podutil.GenerateFlag(`log-wallet`, ``, `--log-wallet=info`, `sets log level for wallet`, true),
		podutil.GenerateFlag(`log-wallet-chain`, ``, `--log-wallet-chain=info`, `sets log level for wallet-chain`, true),
		podutil.GenerateFlag(`log-wallet-rpc-rpcserver`, ``, `--log-wallet-rpc-rpcserver=info`, `sets log level for wallet-rpc-rpcserver`, true),
		podutil.GenerateFlag(`log-wallet-rpc-legacyrpc`, ``, `--log-wallet-rpc-legacyrpc=info`, `sets log level for wallet-rpc-legacyrpc`, true),
		podutil.GenerateFlag(`log-wallet-wtxmgr`, ``, `--log-wallet-wtxmgr=info`, `sets log level for wallet-wtxmgr`, true),
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
			Log.Tracef.Print("setting debug level %s", dl)
			Config.Node.DebugLevel = dl
			Log.SetLevel(dl)
			for i := range logger.Levels {
				logger.Levels[i] = dl
			}
		}
		Log.Debugf.Print("pod/node version %s", n.Version())
		if ctx.Is("version") {
			fmt.Println("pod/node version", n.Version())
			clog.Shutdown()
		}
		Log.Trace.Print("running command")

		var cfgFile string
		if cfgFile, ok = ctx.Get("configfile"); !ok {
			cfgFile = n.DefaultConfigFile
		}
		if ctx.Is("init") {
			Log.Debugf.Print("writing default configuration to %s", cfgFile)
			writeDefaultConfig(cfgFile)
		} else {
			Log.Infof.Print("loading configuration from %s", cfgFile)
			if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
				Log.Warn.Print("configuration file does not exist, creating new one")
				writeDefaultConfig(cfgFile)
			} else {
				Log.Debug.Print("reading app configuration from", cfgFile)
				cfgData, err := ioutil.ReadFile(cfgFile)
				if err != nil {
					Log.Error.Print("reading app config file:", err.Error())
					clog.Shutdown()
				}
				Log.Tracef.Print("parsing app configuration\n%s", cfgData)
				err = json.Unmarshal(cfgData, &Config)
				if err != nil {
					Log.Error.Print("parsing app config file:", err.Error())
					clog.Shutdown()
				}
			}
		}
		configNode(&ctx, cfgFile)
		runNode()
		clog.Shutdown()
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
			Log.Warn <- err.Error()
		}
	}
	if getIfIs(ctx, "disablebanning", r) {
		Config.Node.DisableBanning = *r == "true"
	}
	if getIfIs(ctx, "banduration", r) {
		if err := podutil.ParseDuration(*r, "banduration", &Config.Node.BanDuration); err != nil {
			Log.Warn <- err.Error()
		}
	}
	if getIfIs(ctx, "banthreshold", r) {
		var bt int
		if err := podutil.ParseInteger(*r, "banthtreshold", &bt); err != nil {
			Log.Warn <- err.Error()
		} else {
			Config.Node.BanThreshold = uint32(bt)
		}
	}
	if getIfIs(ctx, "whitelists", r) {
		podutil.NormalizeAddresses(*r, n.DefaultPort, &Config.Node.Whitelists)
	}
	if getIfIs(ctx, "rpcuser", r) {
		Config.Node.RPCUser = *r
	}
	if getIfIs(ctx, "rpcpass", r) {
		Config.Node.RPCPass = *r
	}
	if getIfIs(ctx, "rpclimituser", r) {
		Config.Node.RPCLimitUser = *r
	}
	if getIfIs(ctx, "rpclimitpass", r) {
		Config.Node.RPCLimitPass = *r
	}
	if getIfIs(ctx, "rpclisteners", r) {
		podutil.NormalizeAddresses(*r, n.DefaultRPCPort, &Config.Node.RPCListeners)
	}
	if getIfIs(ctx, "rpccert", r) {
		Config.Node.RPCCert = n.CleanAndExpandPath(*r)
	}
	if getIfIs(ctx, "rpckey", r) {
		Config.Node.RPCKey = n.CleanAndExpandPath(*r)
	}
	if getIfIs(ctx, "tls", r) {
		Config.Node.TLS = *r == "true"
	}
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
			Log.Warn <- err.Error()
		}
	}
	if getIfIs(ctx, "freetxrelaylimit", r) {
		if err := podutil.ParseFloat(*r, "freetxrelaylimit", &Config.Node.FreeTxRelayLimit); err != nil {
			Log.Warn <- err.Error()
		}
	}
	if getIfIs(ctx, "norelaypriority", r) {
		Config.Node.NoRelayPriority = *r == "true"
	}
	if getIfIs(ctx, "trickleinterval", r) {
		if err := podutil.ParseDuration(*r, "trickleinterval", &Config.Node.TrickleInterval); err != nil {
			Log.Warn <- err.Error()
		}
	}
	if getIfIs(ctx, "maxorphantxs", r) {
		if err := podutil.ParseInteger(*r, "maxorphantxs", &Config.Node.MaxOrphanTxs); err != nil {
			Log.Warn <- err.Error()
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
			Log.Warn <- err.Error()
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
			Log.Warn <- err.Error()
		}
	}
	if getIfIs(ctx, "blockmaxsize", r) {
		if err := podutil.ParseUint32(*r, "blockmaxsize", &Config.Node.BlockMaxSize); err != nil {
			Log.Warn <- err.Error()
		}
	}
	if getIfIs(ctx, "blockminweight", r) {
		if err := podutil.ParseUint32(*r, "blockminweight", &Config.Node.BlockMinWeight); err != nil {
			Log.Warn <- err.Error()
		}
	}
	if getIfIs(ctx, "blockmaxweight", r) {
		if err := podutil.ParseUint32(*r, "blockmaxweight", &Config.Node.BlockMaxWeight); err != nil {
			Log.Warn <- err.Error()
		}
	}
	if getIfIs(ctx, "blockprioritysize", r) {
		if err := podutil.ParseUint32(*r, "blockmaxweight", &Config.Node.BlockPrioritySize); err != nil {
			Log.Warn <- err.Error()
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
			Log.Warn <- err.Error()
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
	logger.SetLogging(ctx)
	if ctx.Is("save") {
		Log.Infof.Print("saving config file to %s", cfgFile)
		j, err := json.MarshalIndent(Config, "", "  ")
		if err != nil {
			Log.Error.Print("saving config file:", err.Error())
		}
		j = append(j, '\n')
		Log.Tracef.Print("JSON formatted config file\n%s", j)
		err = ioutil.WriteFile(cfgFile, j, 0600)
		if err != nil {
			Log.Error.Print("writing app config file:", err.Error())
		}
	}
}

func writeDefaultConfig(cfgFile string) {
	defCfg := DefaultConfig()
	defCfg.Node.ConfigFile = cfgFile
	j, err := json.MarshalIndent(defCfg, "", "  ")
	if err != nil {
		Log.Error.Print("marshalling default app config file:", err.Error())
	}
	j = append(j, '\n')
	Log.Tracef.Print("JSON formatted config file\n%s", j)
	err = ioutil.WriteFile(cfgFile, j, 0600)
	if err != nil {
		Log.Error.Print("writing default app config file:", err.Error())
	}
	// if we are writing default config we also want to use it
	Config = defCfg
}

func DefaultConfig() *Cfg {
	return &Cfg{
		Node: &n.Config{
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
		LogLevels: logger.GetDefault(),
	}
}
