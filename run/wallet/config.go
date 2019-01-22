package walletrun

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
var Log = cl.NewSubSystem("run/wallet", "info")
var log = Log.Ch

// UseLogger uses a specified Logger to output package logging info. This should be used in preference to SetLogWriter if the caller is also using log.
func UseLogger(logger *cl.SubSystem) {
	Log = logger
	log = Log.Ch
}

// ConfigAndLog is the combined app and logging configuration data
type ConfigAndLog struct {
	Wallet *w.Config
	Levels map[string]string
}

// Config is the combined app and log levels configuration
var Config = DefaultConfig()

// Command is a command to send RPC queries to bitcoin RPC protocol server for node and wallet queries
var Command = climax.Command{
	Name:  "wallet",
	Brief: "parallelcoin wallet",
	Help:  "check balances, make payments, manage contacts",
	Flags: []climax.Flag{

		pu.GenFlag("version", "V", `--version`, `show version number and quit`, false),

		pu.GenFlag("configfile", "C", "--configfile=/path/to/conf", "path to configuration file", true),
		pu.GenFlag("datadir", "D", "--datadir=/home/user/.pod", "set the base directory for elements shared between modules", true),

		pu.GenFlag("init", "", "--init", "resets configuration to defaults", false),
		pu.GenFlag("save", "", "--save", "saves current configuration", false),

		pu.GenFlag("create", "", "--create", "create a new wallet if it does not exist", false),
		pu.GenFlag("createtemp", "", "--createtemp", "create temporary wallet (pass=password), must call with --datadir", false),

		pu.GenFlag("appdatadir", "", "--appdatadir=/path/to/appdatadir", "set app data directory for wallet, configuration and logs", true),
		pu.GenFlag("testnet3", "", "--testnet=true", "use testnet", true),
		pu.GenFlag("simnet", "", "--simnet=true", "use simnet", true),
		pu.GenFlag("noinitialload", "", "--noinitialload=true", "defer wallet creation/opening on startup and enable loading wallets over RPC (default with --gui)", true),
		pu.GenFlag("network", "", "--network=mainnet", "connect to specified network: mainnet, testnet, regtestnet or simnet", true),
		pu.GenFlag("profile", "", "--profile=true", "enable HTTP profiling on given port -- NOTE port must be between 1024 and 65536", true),
		pu.GenFlag("gui", "", "--gui=true", "launch GUI (wallet unlock is deferred to let GUI handle)", true),
		pu.GenFlag("walletpass", "", "--walletpass=somepassword", "the public wallet password - only required if the wallet was created with one", true),
		pu.GenFlag("rpcconnect", "", "--rpcconnect=some.address.com:11048", "connect to the RPC of a parallelcoin node for chain queries", true),
		pu.GenFlag("cafile", "", "--cafile=/path/to/cafile", "file containing root certificates to authenticate TLS connections with pod", true),
		pu.GenFlag("enableclienttls", "", "--enableclienttls=false", "enable TLS for the RPC client", true),
		pu.GenFlag("podusername", "", "--podusername=user", "username for node RPC authentication", true),
		pu.GenFlag("podpassword", "", "--podpassword=pa55word", "password for node RPC authentication", true),
		pu.GenFlag("proxy", "", "--proxy=127.0.0.1:9050", "address for proxy for outbound connections", true),
		pu.GenFlag("proxyuser", "", "--proxyuser=user", "username for proxy", true),
		pu.GenFlag("proxypass", "", "--proxypass=pa55word", "password for proxy", true),
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
		{
			Usecase:     "--init --rpcuser=user --rpcpass=pa55word --save",
			Description: "resets the configuration file to default, sets rpc username and password and saves the changes to config after parsing",
		},
	},
	Handle: func(ctx climax.Context) int {
		var dl string
		var ok bool
		if dl, ok = ctx.Get("debuglevel"); ok {
			log <- cl.Tracef{"setting debug level %s", dl}
			Log.SetLevel(dl)
			for i := range logger.Levels {
				logger.Levels[i].SetLevel(dl)
			}
		}
		log <- cl.Debugf{"pod/wallet version %s", w.Version()}
		if ctx.Is("version") {
			fmt.Println("pod/wallet version", w.Version())
			cl.Shutdown()
		}
		var cfgFile string
		if cfgFile, ok = ctx.Get("configfile"); !ok {
			cfgFile = w.DefaultConfigFile
		}
		if ctx.Is("init") {
			log <- cl.Debug{"writing default configuration to", cfgFile}
			WriteDefaultConfig(cfgFile)
		}
		log <- cl.Info{"loading configuration from", cfgFile}
		if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
			log <- cl.Wrn("configuration file does not exist, creating new one")
			WriteDefaultConfig(cfgFile)
		} else {
			log <- cl.Debug{"reading app configuration from", cfgFile}
			cfgData, err := ioutil.ReadFile(cfgFile)
			if err != nil {
				log <- cl.Error{"reading app config file", err.Error()}
				cl.Shutdown()
			}
			log <- cl.Tracef{"parsing app configuration\n%s", cfgData}
			err = json.Unmarshal(cfgData, &Config)
			if err != nil {
				log <- cl.Error{"parsing app config file", err.Error()}
				cl.Shutdown()
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
		pu.NormalizeAddress(*r, "3131", &Config.Wallet.Profile)
	}
	if getIfIs(ctx, "gui", r) {
		Config.Wallet.GUI = *r == "true"
	}
	if getIfIs(ctx, "walletpass", r) {
		Config.Wallet.WalletPass = *r
	}
	if getIfIs(ctx, "rpcconnect", r) {
		pu.NormalizeAddress(*r, "11048", &Config.Wallet.RPCConnect)
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
		pu.NormalizeAddress(*r, "11048", &Config.Wallet.Proxy)
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
		log <- cl.Info{"saving config file to", cfgFile}
		j, err := json.MarshalIndent(Config, "", "  ")
		if err != nil {
			log <- cl.Error{"writing app config file", err.Error()}
		}
		j = append(j, '\n')
		log <- cl.Trace{"JSON formatted config file\n", string(j)}
		ioutil.WriteFile(cfgFile, j, 0600)
	}
}

// WriteConfig creates and writes the config file in the requested location
func WriteConfig(cfgFile string, c *ConfigAndLog) {
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

// WriteDefaultConfig creates and writes a default config to the requested location
func WriteDefaultConfig(cfgFile string) {
	defCfg := DefaultConfig()
	defCfg.Wallet.ConfigFile = cfgFile
	j, err := json.MarshalIndent(defCfg, "", "  ")
	if err != nil {
		log <- cl.Error{"marshalling configuration", err.Error()}
	}
	j = append(j, '\n')
	log <- cl.Trace{"JSON formatted config file\n", string(j)}
	err = ioutil.WriteFile(cfgFile, j, 0600)
	if err != nil {
		log <- cl.Error{"writing app config file", err.Error()}
	}
	// if we are writing default config we also want to use it
	Config = defCfg
}

func DefaultConfig() *ConfigAndLog {
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
		Levels: logger.GetDefaultConfig(),
	}
}
