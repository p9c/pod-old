package wallet

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"git.parallelcoin.io/pod/lib/clog"
	w "git.parallelcoin.io/pod/module/wallet"
	"git.parallelcoin.io/pod/module/wallet/wallet"
	"git.parallelcoin.io/pod/run/logger"
	"github.com/tucnak/climax"
)

var log = clog.NewSubSystem("Wallet", clog.Ninf)

// Config is the default configuration native to ctl
var Config = new(w.Config)

// ConfigAndLog is the combined app and logging configuration data
type ConfigAndLog struct {
	Wallet *w.Config
	Levels map[string]string
}

// CombinedCfg is the combined app and log levels configuration
var CombinedCfg = ConfigAndLog{
	Wallet: Config,
	Levels: logger.Levels,
}

// Command is a command to send RPC queries to bitcoin RPC protocol server for node and wallet queries
var Command = climax.Command{
	Name:  "wallet",
	Brief: "parallelcoin wallet",
	Help:  "check balances, make payments, manage contacts",
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
			Help:     "sets debuglevel, default info, sets the baseline for others not specified below (logging is per-library basis",
			Variable: true,
		},
		{
			Name:     "log-database",
			Usage:    "--log-database=debug",
			Help:     "sets log level for database",
			Variable: true,
		},
		{
			Name:     "log-txscript",
			Usage:    "--log-txscript=debug",
			Help:     "sets log level for txscript",
			Variable: true,
		},
		{
			Name:     "log-peer",
			Usage:    "--log-peer=debug",
			Help:     "sets log level for peer",
			Variable: true,
		},
		{
			Name:     "log-netsync",
			Usage:    "--log-netsync=debug",
			Help:     "sets log level for netsync",
			Variable: true,
		},
		{
			Name:     "log-rpcclient",
			Usage:    "--log-rpcclient=debug",
			Help:     "sets log level for rpcclient",
			Variable: true,
		},
		{
			Name:     "addrmgr",
			Usage:    "--log-addrmgr=debug",
			Help:     "sets log level for mgr",
			Variable: true,
		},
		{
			Name:     "log-blockchain-indexers",
			Usage:    "--log-blockchain-indexers=debug",
			Help:     "sets log level for blockchain-indexers",
			Variable: true,
		},
		{
			Name:     "log-blockchain",
			Usage:    "--log-blockchain=debug",
			Help:     "sets log level for blockchain",
			Variable: true,
		},
		{
			Name:     "log-mining-cpuminer",
			Usage:    "--log-mining-cpuminer=debug",
			Help:     "sets log level for mining-cpuminer",
			Variable: true,
		},
		{
			Name:     "log-mining",
			Usage:    "--log-mining=debug",
			Help:     "sets log level for mining",
			Variable: true,
		},
		{
			Name:     "log-mining-controller",
			Usage:    "--log-mining-controller=debug",
			Help:     "sets log level for mining-controller",
			Variable: true,
		},
		{
			Name:     "log-connmgr",
			Usage:    "--log-connmgr=debug",
			Help:     "sets log level for connmgr",
			Variable: true,
		},
		{
			Name:     "log-spv",
			Usage:    "--log-spv=debug",
			Help:     "sets log level for spv",
			Variable: true,
		},
		{
			Name:     "log-node-mempool",
			Usage:    "--log-node-mempool=debug",
			Help:     "sets log level for node-mempool",
			Variable: true,
		},
		{
			Name:     "log-node",
			Usage:    "--log-node=debug",
			Help:     "sets log level for node",
			Variable: true,
		},
		{
			Name:     "log-wallet-wallet",
			Usage:    "--log-wallet-wallet=debug",
			Help:     "sets log level for wallet-wallet",
			Variable: true,
		},
		{
			Name:     "log-wallet-tx",
			Usage:    "--log-wallet-tx=debug",
			Help:     "sets log level for wallet-tx",
			Variable: true,
		},
		{
			Name:     "log-wallet-votingpool",
			Usage:    "--log-wallet-votingpool=debug",
			Help:     "sets log level for wallet-votingpool",
			Variable: true,
		},
		{
			Name:     "log-wallet",
			Usage:    "--log-wallet=debug",
			Help:     "sets log level for wallet",
			Variable: true,
		},
		{
			Name:     "log-wallet-chain",
			Usage:    "--log-wallet-chain=debug",
			Help:     "sets log level for wallet-chain",
			Variable: true,
		},
		{
			Name:     "log-wallet-rpc-rpcserver",
			Usage:    "--log-wallet-rpc-rpcserver=debug",
			Help:     "sets log level for wallet-rpc-rpcserver",
			Variable: true,
		},
		{
			Name:     "log-wallet-rpc-legacyrpc",
			Usage:    "--log-wallet-rpc-legacyrpc=debug",
			Help:     "sets log level for wallet-rpc-legacyrpc",
			Variable: true,
		},
		{
			Name:     "log-wallet-wtxmgr",
			Usage:    "--log-wallet-wtxmgr=debug",
			Help:     "sets log level for wallet-wtxmgr",
			Variable: true,
		},
		{
			Name:     "create",
			Usage:    "--create",
			Help:     "create a new wallet if it does not exist",
			Variable: false,
		},
		{
			Name:     "createtemp",
			Usage:    "--createtemp",
			Help:     "create temporary wallet (pass=password), must call with --datadir",
			Variable: false,
		},
		{
			Name:     "appdatadir",
			Usage:    "--appdatadir=/path/to/appdatadir",
			Help:     "set app data directory for wallet, configuration and logs",
			Variable: true,
		},
		{
			Name:     "testnet3",
			Usage:    "--testnet=true",
			Help:     "use testnet",
			Variable: true,
		},
		{
			Name:     "simnet",
			Usage:    "--simnet=true",
			Help:     "use simnet",
			Variable: true,
		},
		{
			Name:     "noinitialload",
			Usage:    "--noinitialload=true",
			Help:     "defer wallet creation/opening on startup and enable loading wallets over RPC (default with --gui)",
			Variable: true,
		},
		{
			Name:     "network",
			Usage:    "--network=mainnet",
			Help:     "connect to specified network: mainnet, testnet, regtestnet or simnet",
			Variable: true,
		},
		{
			Name:     "profile",
			Usage:    "--profile=true",
			Help:     "enable HTTP profiling on given port -- NOTE port must be between 1024 and 65536",
			Variable: true,
		},
		{
			Name:     "gui",
			Usage:    "--gui=true",
			Help:     "launch GUI (wallet unlock is deferred to let GUI handle)",
			Variable: true,
		},
		{
			Name:     "walletpass",
			Usage:    "--walletpass=somepassword",
			Help:     "the public wallet password - only required if the wallet was created with one",
			Variable: true,
		},
		{
			Name:     "rpcconnect",
			Usage:    "--rpcconnect=some.address.com:11048",
			Help:     "connect to the RPC of a parallelcoin node for chain queries",
			Variable: true,
		},
		{
			Name:     "cafile",
			Usage:    "--cafile=/path/to/cafile",
			Help:     "file containing root certificates to authenticate TLS connections with pod",
			Variable: true,
		},
		{
			Name:     "enableclienttls",
			Usage:    "--enableclienttls=false",
			Help:     "enable TLS for the RPC client",
			Variable: true,
		},
		{
			Name:     "podusername",
			Usage:    "--podusername=user",
			Help:     "username for node RPC authentication",
			Variable: true,
		},
		{
			Name:     "podpassword",
			Usage:    "--podpassword=pa55word",
			Help:     "password for node RPC authentication",
			Variable: true,
		},
		{
			Name:     "proxy",
			Usage:    "--proxy=127.0.0.1:9050",
			Help:     "address for proxy for outbound connections",
			Variable: true,
		},
		{
			Name:     "proxyuser",
			Usage:    "--proxyuser=user",
			Help:     "username for proxy",
			Variable: true,
		},
		{
			Name:     "proxypass",
			Usage:    "--proxypass=pa55word",
			Help:     "password for proxy",
			Variable: true,
		},
		{
			Name:     "rpccert",
			Usage:    "--rpccert=/path/to/rpccert",
			Help:     "file containing the RPC tls certificate",
			Variable: true,
		},
		{
			Name:     "rpckey",
			Usage:    "--rpckey=/path/to/rpckey",
			Help:     "file containing RPC tls key",
			Variable: true,
		},
		{
			Name:     "onetimetlskey",
			Usage:    "--onetimetlskey=true",
			Help:     "generate a new TLS certpair but only write certs to disk",
			Variable: true,
		},
		{
			Name:     "enableservertls",
			Usage:    "--enableservertls=false",
			Help:     "enable TLS on wallet RPC",
			Variable: true,
		},
		{
			Name:     "legacyrpclisteners",
			Usage:    "--legacyrpclisteners=127.0.0.1:11046",
			Help:     "add a listener for the legacy RPC",
			Variable: true,
		},
		{
			Name:     "legacyrpcmaxclients",
			Usage:    "--legacyrpcmaxclients=10",
			Help:     "maximum number of connections for legacy RPC",
			Variable: true,
		},
		{
			Name:     "legacyrpcmaxwebsockets",
			Usage:    "--legacyrpcmaxwebsockets=10",
			Help:     "maximum number of websockets for legacy RPC",
			Variable: true,
		},
		{
			Name:     "username",
			Short:    "-u",
			Usage:    "--username=user",
			Help:     "username for wallet RPC, used also for node if podusername is empty",
			Variable: true,
		},
		{
			Name:     "password",
			Short:    "-P",
			Usage:    "--password=pa55word",
			Help:     "password for wallet RPC, also used for node if podpassord",
			Variable: true,
		},
		{
			Name:     "experimentalrpclisteners",
			Usage:    "--experimentalrpclisteners=127.0.0.1:11045",
			Help:     "enable experimental RPC service on this address",
			Variable: true,
		},
		{
			Name:     "datadir",
			Usage:    "--datadir=/home/user/.pod",
			Help:     "set the base directory for elements shared between modules",
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
		log.Debugf.Print("node version %s", w.Version())
		if ctx.Is("version") {
			fmt.Println("node version", w.Version())
			clog.Shutdown()
		}
		log.Trace.Print("running command")

		var cfgFile string
		if cfgFile, ok = ctx.Get("configfile"); !ok {
			cfgFile = w.DefaultConfigFile
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
				// logCfgFile := getLogConfFileName()
				// log.Debug.Print("reading logger configuration from", logCfgFile)
				// logCfgData, err := ioutil.ReadFile(logCfgFile)
				// if err != nil {
				// 	log.Error.Print(err.Error())
				// 	clog.Shutdown()
				// }
				// log.Tracef.Print("parsing logger configuration\n%s", logCfgData)
				// err = json.Unmarshal(logCfgData, &CombinedCfg.Levels)
				// if err != nil {
				// 	log.Error.Print(err.Error())
				// 	clog.Shutdown()
				// }
				configNode(&ctx, cfgFile)
			}
		}
		runNode()
		clog.Shutdown()
		return 0
	},
}

func configNode(ctx *climax.Context, cfgFile string) {
	// Apply all configurations specified on commandline
	if ctx.Is("create") {
		Config.Create = true
	}
	if ctx.Is("createtemp") {
		Config.CreateTemp = true
	}
	if ctx.Is("appdatadir") {
		r, _ := ctx.Get("appdatadir")
		Config.AppDataDir = r
	}
	if ctx.Is("noinitialload") {
		r, _ := ctx.Get("noinitialload")
		Config.NoInitialLoad = r == "true"
	}
	if ctx.Is("logdir") {
		r, _ := ctx.Get("logdir")
		Config.LogDir = r
	}
	if ctx.Is("profile") {
		r, _ := ctx.Get("profile")
		Config.Profile = r
	}
	if ctx.Is("gui") {
		r, _ := ctx.Get("gui")
		Config.GUI = r == "true"
	}
	if ctx.Is("walletpass") {
		r, _ := ctx.Get("walletpass")
		Config.WalletPass = r
	}
	if ctx.Is("rpcconnect") {
		r, _ := ctx.Get("rpcconnect")
		Config.RPCConnect = r
	}
	if ctx.Is("cafile") {
		r, _ := ctx.Get("cafile")
		Config.CAFile = r
	}
	if ctx.Is("enableclienttls") {
		r, _ := ctx.Get("enableclienttls")
		Config.EnableClientTLS = r == "true"
	}
	if ctx.Is("podusername") {
		r, _ := ctx.Get("podusername")
		Config.PodUsername = r
	}
	if ctx.Is("podpassword") {
		r, _ := ctx.Get("podpassword")
		Config.PodPassword = r
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
	if ctx.Is("rpccert") {
		r, _ := ctx.Get("rpccert")
		Config.RPCCert = r
	}
	if ctx.Is("rpckey") {
		r, _ := ctx.Get("rpckey")
		Config.RPCKey = r
	}
	if ctx.Is("onetimetlskey") {
		r, _ := ctx.Get("onetimetlskey")
		Config.OneTimeTLSKey = r == "true"
	}
	if ctx.Is("enableservertls") {
		r, _ := ctx.Get("enableservertls")
		Config.EnableServerTLS = r == "true"
	}
	if ctx.Is("legacyrpclisteners") {
		r, _ := ctx.Get("legacyrpclisteners")
		Config.LegacyRPCListeners = strings.Split(r, " ")
	}
	if ctx.Is("legacyrpcmaxclients") {
		r, _ := ctx.Get("legacyrpcmaxclients")
		_, err := fmt.Sscanf(r, "%d", Config.LegacyRPCMaxClients)
		if err != nil {
			log.Errorf.Print("malformed legacymaxclients: `%s` leaving set at `%d`",
				r, Config.LegacyRPCMaxClients)
		}
	}
	if ctx.Is("legacyrpcmaxwebsockets") {
		r, _ := ctx.Get("legacyrpcmaxwebsockets")
		_, err := fmt.Sscanf(r, "%d", Config.LegacyRPCMaxWebsockets)
		if err != nil {
			log.Errorf.Print("malformed legacyrpcmaxwebsockets: `%s` leaving set at `%d`",
				r, Config.LegacyRPCMaxWebsockets)
		}
	}
	if ctx.Is("username") {
		r, _ := ctx.Get("username")
		Config.Username = r
	}
	if ctx.Is("password") {
		r, _ := ctx.Get("password")
		Config.Password = r
	}
	if ctx.Is("experimentalrpclisteners") {
		r, _ := ctx.Get("experimentalrpclisteners")
		Config.ExperimentalRPCListeners = strings.Split(r, " ")
	}
	if ctx.Is("datadir") {
		r, _ := ctx.Get("datadir")
		Config.DataDir = r
	}
	if ctx.Is("network") {
		r, _ := ctx.Get("network")
		switch r {
		case "testnet":
			Config.TestNet3, Config.SimNet = true, false
		case "simnet":
			Config.TestNet3, Config.SimNet = false, true
		default:
			Config.TestNet3, Config.SimNet = false, false
		}
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
	defCfg.Wallet.ConfigFile = cfgFile
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
		Wallet: &w.Config{
			ConfigFile:             w.DefaultConfigFile,
			DataDir:                w.DefaultDataDir,
			AppDataDir:             w.DefaultAppDataDir,
			LogDir:                 w.DefaultLogDir,
			RPCKey:                 w.DefaultRPCKeyFile,
			RPCCert:                w.DefaultRPCCertFile,
			WalletPass:             wallet.InsecurePubPassphrase,
			CAFile:                 "",
			LegacyRPCMaxClients:    w.DefaultRPCMaxClients,
			LegacyRPCMaxWebsockets: w.DefaultRPCMaxWebsockets,
		},
		Levels: logger.GetDefault(),
	}
}
