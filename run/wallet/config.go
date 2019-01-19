package wallet

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"git.parallelcoin.io/pod/lib/clog"
	n "git.parallelcoin.io/pod/module/node"
	w "git.parallelcoin.io/pod/module/wallet"
	"git.parallelcoin.io/pod/module/wallet/wallet"
	"git.parallelcoin.io/pod/run/logger"
	"git.parallelcoin.io/pod/run/util"
	"github.com/tucnak/climax"
)

// Log is the main logger for wallet
var Log = clog.NewSubSystem("pod/wallet", clog.Ninf)

// ConfigAndLog is the combined app and logging configuration data
type ConfigAndLog struct {
	Wallet *w.Config
	Levels map[string]string
}

// Config is the combined app and log levels configuration
var Config = defaultConfig()

// Command is a command to send RPC queries to bitcoin RPC protocol server for node and wallet queries
var Command = climax.Command{
	Name:  "wallet",
	Brief: "parallelcoin wallet",
	Help:  "check balances, make payments, manage contacts",
	Flags: []climax.Flag{

		podutil.GenerateFlag("version", "V", `--version`, `show version number and quit`, false),

		podutil.GenerateFlag("configfile", "C", "--configfile=/path/to/conf", "path to configuration file", true),
		podutil.GenerateFlag("datadir", "D", "--datadir=/home/user/.pod", "set the base directory for elements shared between modules", true),

		podutil.GenerateFlag("init", "", "--init", "resets configuration to defaults", false),
		podutil.GenerateFlag("save", "", "--save", "saves current configuration", false),

		podutil.GenerateFlag("create", "", "--create", "create a new wallet if it does not exist", false),
		podutil.GenerateFlag("createtemp", "", "--createtemp", "create temporary wallet (pass=password), must call with --datadir", false),

		podutil.GenerateFlag("appdatadir", "", "--appdatadir=/path/to/appdatadir", "set app data directory for wallet, configuration and logs", true),
		podutil.GenerateFlag("testnet3", "", "--testnet=true", "use testnet", true),
		podutil.GenerateFlag("simnet", "", "--simnet=true", "use simnet", true),
		podutil.GenerateFlag("noinitialload", "", "--noinitialload=true", "defer wallet creation/opening on startup and enable loading wallets over RPC (default with --gui)", true),
		podutil.GenerateFlag("network", "", "--network=mainnet", "connect to specified network: mainnet, testnet, regtestnet or simnet", true),
		podutil.GenerateFlag("profile", "", "--profile=true", "enable HTTP profiling on given port -- NOTE port must be between 1024 and 65536", true),
		podutil.GenerateFlag("gui", "", "--gui=true", "launch GUI (wallet unlock is deferred to let GUI handle)", true),
		podutil.GenerateFlag("walletpass", "", "--walletpass=somepassword", "the public wallet password - only required if the wallet was created with one", true),
		podutil.GenerateFlag("rpcconnect", "", "--rpcconnect=some.address.com:11048", "connect to the RPC of a parallelcoin node for chain queries", true),
		podutil.GenerateFlag("cafile", "", "--cafile=/path/to/cafile", "file containing root certificates to authenticate TLS connections with pod", true),
		podutil.GenerateFlag("enableclienttls", "", "--enableclienttls=false", "enable TLS for the RPC client", true),
		podutil.GenerateFlag("podusername", "", "--podusername=user", "username for node RPC authentication", true),
		podutil.GenerateFlag("podpassword", "", "--podpassword=pa55word", "password for node RPC authentication", true),
		podutil.GenerateFlag("proxy", "", "--proxy=127.0.0.1:9050", "address for proxy for outbound connections", true),
		podutil.GenerateFlag("proxyuser", "", "--proxyuser=user", "username for proxy", true),
		podutil.GenerateFlag("proxypass", "", "--proxypass=pa55word", "password for proxy", true),
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

		podutil.GenerateFlag("log-database", "", "--log-database=debug", "sets log level for database", true),
		podutil.GenerateFlag("log-txscript", "", "--log-txscript=debug", "sets log level for txscript", true),
		podutil.GenerateFlag("log-peer", "", "--log-peer=debug", "sets log level for peer", true),
		podutil.GenerateFlag("log-netsync", "", "--log-netsync=debug", "sets log level for netsync", true),
		podutil.GenerateFlag("log-rpcclient", "", "--log-rpcclient=debug", "sets log level for rpcclient", true),
		podutil.GenerateFlag("addrmgr", "", "--log-addrmgr=debug", "sets log level for mgr", true),
		podutil.GenerateFlag("log-blockchain-indexers", "", "--log-blockchain-indexers=debug", "sets log level for blockchain-indexers", true),
		podutil.GenerateFlag("log-blockchain", "", "--log-blockchain=debug", "sets log level for blockchain", true),
		podutil.GenerateFlag("log-mining-cpuminer", "", "--log-mining-cpuminer=debug", "sets log level for mining-cpuminer", true),
		podutil.GenerateFlag("log-mining", "", "--log-mining=debug", "sets log level for mining", true),
		podutil.GenerateFlag("log-mining-controller", "", "--log-mining-controller=debug", "sets log level for mining-controller", true),
		podutil.GenerateFlag("log-connmgr", "", "--log-connmgr=debug", "sets log level for connmgr", true),
		podutil.GenerateFlag("log-spv", "", "--log-spv=debug", "sets log level for spv", true),
		podutil.GenerateFlag("log-node-mempool", "", "--log-node-mempool=debug", "sets log level for node-mempool", true),
		podutil.GenerateFlag("log-node", "", "--log-node=debug", "sets log level for node", true),
		podutil.GenerateFlag("log-wallet-wallet", "", "--log-wallet-wallet=debug", "sets log level for wallet-wallet", true),
		podutil.GenerateFlag("log-wallet-tx", "", "--log-wallet-tx=debug", "sets log level for wallet-tx", true),
		podutil.GenerateFlag("log-wallet-votingpool", "", "--log-wallet-votingpool=debug", "sets log level for wallet-votingpool", true),
		podutil.GenerateFlag("log-wallet", "", "--log-wallet=debug", "sets log level for wallet", true),
		podutil.GenerateFlag("log-wallet-chain", "", "--log-wallet-chain=debug", "sets log level for wallet-chain", true),
		podutil.GenerateFlag("log-wallet-rpc-rpcserver", "", "--log-wallet-rpc-rpcserver=debug", "sets log level for wallet-rpc-rpcserver", true),
		podutil.GenerateFlag("log-wallet-rpc-legacyrpc", "", "--log-wallet-rpc-legacyrpc=debug", "sets log level for wallet-rpc-legacyrpc", true),
		podutil.GenerateFlag("log-wallet-wtxmgr", "", "--log-wallet-wtxmgr=debug", "sets log level for wallet-wtxmgr", true),
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
			Log.SetLevel(dl)
			for i := range logger.Levels {
				logger.Levels[i] = dl
			}
		}
		Log.Debugf.Print("pod/wallet version %s", w.Version())
		if ctx.Is("version") {
			fmt.Println("pod/wallet version", w.Version())
			clog.Shutdown()
		}
		Log.Trace.Print("running command wallet")
		var cfgFile string
		if cfgFile, ok = ctx.Get("configfile"); !ok {
			cfgFile = w.DefaultConfigFile
		}
		if ctx.Is("init") {
			Log.Debugf.Print("writing default configuration to %s", cfgFile)
			writeDefaultConfig(cfgFile)
			configNode(&ctx, cfgFile)
		} else {
			Log.Infof.Print("loading configuration from %s", cfgFile)
			if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
				Log.Warn.Print("configuration file does not exist, creating new one")
				writeDefaultConfig(cfgFile)
				configNode(&ctx, cfgFile)
			} else {
				Log.Debug.Print("reading app configuration from", cfgFile)
				cfgData, err := ioutil.ReadFile(cfgFile)
				if err != nil {
					Log.Error.Print("reading app config file", err.Error())
					clog.Shutdown()
				}
				Log.Tracef.Print("parsing app configuration\n%s", cfgData)
				err = json.Unmarshal(cfgData, &Config)
				if err != nil {
					Log.Error.Print("parsing app config file", err.Error())
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
	if getIfIs(ctx, "rpcconnect", r) {
		podutil.NormalizeAddress(*r, "11048", &Config.Wallet.RPCConnect)
	}
	if getIfIs(ctx, "cafile", r) {
		Config.Wallet.CAFile = n.CleanAndExpandPath(*r)
	}
	if getIfIs(ctx, "enableclienttls", r) {
		Config.Wallet.EnableClientTLS = *r == "true"
	}
	if getIfIs(ctx, "podusername", r) {
		Config.Wallet.PodUsername = *r
	}
	if getIfIs(ctx, "podpassword", r) {
		Config.Wallet.PodPassword = *r
	}
	if getIfIs(ctx, "proxy", r) {
		podutil.NormalizeAddress(*r, "11048", &Config.Wallet.Proxy)
	}
	if getIfIs(ctx, "proxyuser", r) {
		Config.Wallet.ProxyUser = *r
	}
	if getIfIs(ctx, "proxypass", r) {
		Config.Wallet.ProxyPass = *r
	}
	if getIfIs(ctx, "rpccert", r) {
		Config.Wallet.RPCCert = n.CleanAndExpandPath(*r)
	}
	if getIfIs(ctx, "rpckey", r) {
		Config.Wallet.RPCKey = n.CleanAndExpandPath(*r)
	}
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
			Log.Warn <- err.Error()
		} else {
			Config.Wallet.LegacyRPCMaxClients = int64(bt)
		}
	}
	if getIfIs(ctx, "legacyrpcmaxwebsockets", r) {
		_, err := fmt.Sscanf(*r, "%d", Config.Wallet.LegacyRPCMaxWebsockets)
		if err != nil {
			Log.Errorf.Print("malformed legacyrpcmaxwebsockets: `%s` leaving set at `%d`",
				r, Config.Wallet.LegacyRPCMaxWebsockets)
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
	if getIfIs(ctx, "datadir", r) {
		Config.Wallet.DataDir = *r
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
		Log.Infof.Print("saving config file to %s", cfgFile)
		j, err := json.MarshalIndent(Config, "", "  ")
		if err != nil {
			Log.Error.Print("writing app config file", err.Error())
		}
		j = append(j, '\n')
		Log.Tracef.Print("JSON formatted config file\n%s", j)
		ioutil.WriteFile(cfgFile, j, 0600)
	}
}

func writeDefaultConfig(cfgFile string) {
	defCfg := defaultConfig()
	defCfg.Wallet.ConfigFile = cfgFile
	j, err := json.MarshalIndent(defCfg, "", "  ")
	if err != nil {
		Log.Error.Print("marshalling configuration", err.Error())
	}
	j = append(j, '\n')
	Log.Tracef.Print("JSON formatted config file\n%s", j)
	err = ioutil.WriteFile(cfgFile, j, 0600)
	if err != nil {
		Log.Error.Print("writing app config file", err.Error())
	}
	// if we are writing default config we also want to use it
	Config = defCfg
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
