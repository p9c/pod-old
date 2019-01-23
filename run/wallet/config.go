package walletrun

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"git.parallelcoin.io/pod/lib/clog"
	n "git.parallelcoin.io/pod/module/node"
	w "git.parallelcoin.io/pod/module/wallet"
	"git.parallelcoin.io/pod/module/wallet/netparams"
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

var f = pu.GenFlag
var t = pu.GenTrig
var s = pu.GenShort
var l = pu.GenLog

// Command is a command to send RPC queries to bitcoin RPC protocol server for node and wallet queries
var Command = climax.Command{
	Name:  "wallet",
	Brief: "parallelcoin wallet",
	Help:  "check balances, make payments, manage contacts",
	Flags: []climax.Flag{
		t("version", "V", "show version number and quit"),

		s("configfile", "C", "path to configuration file"),
		s("datadir", "D", "set the pod base directory"),
		f("appdatadir", "set app data directory for wallet, configuration and logs"),

		t("init", "i", "resets configuration to defaults"),
		t("save", "S", "saves current flags into configuration"),

		f("create", "create a new wallet if it does not exist"),
		f("createtemp", "create temporary wallet (pass=walletpass) requires --datadir"),

		f("gui", "launch GUI"),
		f("rpcconnect", "connect to the RPC of a parallelcoin node for chain queries"),

		f("podusername", "username for node RPC authentication"),
		f("podpassword", "password for node RPC authentication"),

		f("walletpass", "the public wallet password - only required if the wallet was created with one"),

		f("noinitialload", "defer wallet load to be triggered by RPC"),
		f("network", "connect to (mainnet|testnet|regtestnet|simnet)"),

		f("profile", "enable HTTP profiling on given port (1024-65536)"),

		f("rpccert", "file containing the RPC tls certificate"),
		f("rpckey", "file containing RPC TLS key"),
		f("onetimetlskey", "generate a new TLS certpair don't save key"),
		f("cafile", "certificate authority for custom TLS CA"),
		f("enableclienttls", "enable TLS for the RPC client"),
		f("enableservertls", "enable TLS on wallet RPC server"),

		f("proxy", "proxy address for outbound connections"),
		f("proxyuser", "username for proxy server"),
		f("proxypass", "password for proxy server"),

		f("legacyrpclisteners", "add a listener for the legacy RPC"),
		f("legacyrpcmaxclients", "max connections for legacy RPC"),
		f("legacyrpcmaxwebsockets", "max websockets for legacy RPC"),

		f("username", "username for wallet RPC when podusername is empty"),
		f("password", "password for wallet RPC when podpassword is omitted"),
		f("experimentalrpclisteners", "listener for experimental rpc"),

		s("debuglevel", "d", "sets debuglevel, specify per-library below"),

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
				WriteDefaultConfig(cfgFile)
			}
			log <- cl.Tracef{"parsing app configuration\n%s", cfgData}
			err = json.Unmarshal(cfgData, &Config)
			if err != nil {
				log <- cl.Error{"parsing app config file", err.Error()}
				WriteDefaultConfig(cfgFile)
			}
		}
		configWallet(&ctx, cfgFile)
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

func configWallet(ctx *climax.Context, cfgFile string) {
	log <- cl.Debug{"configuring from command line flags ", os.Args}
	var r *string
	t := ""
	r = &t
	if ctx.Is("create") {
		log <- cl.Dbg("")
		Config.Wallet.Create = true
	}
	if ctx.Is("createtemp") {
		log <- cl.Dbg("")
		Config.Wallet.CreateTemp = true
	}
	if getIfIs(ctx, "appdatadir", r) {
		log <- cl.Dbg("")
		Config.Wallet.AppDataDir = n.CleanAndExpandPath(*r)
	}
	if getIfIs(ctx, "noinitialload", r) {
		log <- cl.Dbg("")
		Config.Wallet.NoInitialLoad = *r == "true"
	}
	if getIfIs(ctx, "logdir", r) {
		log <- cl.Dbg("")
		Config.Wallet.LogDir = n.CleanAndExpandPath(*r)
	}
	if getIfIs(ctx, "profile", r) {
		log <- cl.Dbg("")
		pu.NormalizeAddress(*r, "3131", &Config.Wallet.Profile)
	}
	if getIfIs(ctx, "gui", r) {
		log <- cl.Dbg("")
		Config.Wallet.GUI = *r == "true"
	}
	if getIfIs(ctx, "walletpass", r) {
		log <- cl.Dbg("")
		Config.Wallet.WalletPass = *r
	}
	if getIfIs(ctx, "rpcconnect", r) {
		log <- cl.Dbg("")
		pu.NormalizeAddress(*r, "11048", &Config.Wallet.RPCConnect)
	}
	if getIfIs(ctx, "cafile", r) {
		log <- cl.Dbg("")
		Config.Wallet.CAFile = n.CleanAndExpandPath(*r)
	}
	if getIfIs(ctx, "enableclienttls", r) {
		log <- cl.Dbg("")
		Config.Wallet.EnableClientTLS = *r == "true"
	}
	if getIfIs(ctx, "podusername", r) {
		log <- cl.Dbg("")
		Config.Wallet.PodUsername = *r
	}
	if getIfIs(ctx, "podpassword", r) {
		log <- cl.Dbg("")
		Config.Wallet.PodPassword = *r
	}
	if getIfIs(ctx, "proxy", r) {
		log <- cl.Dbg("")
		pu.NormalizeAddress(*r, "11048", &Config.Wallet.Proxy)
	}
	if getIfIs(ctx, "proxyuser", r) {
		log <- cl.Dbg("")
		Config.Wallet.ProxyUser = *r
	}
	if getIfIs(ctx, "proxypass", r) {
		log <- cl.Dbg("")
		Config.Wallet.ProxyPass = *r
	}
	if getIfIs(ctx, "rpccert", r) {
		log <- cl.Dbg("")
		Config.Wallet.RPCCert = n.CleanAndExpandPath(*r)
	}
	if getIfIs(ctx, "rpckey", r) {
		log <- cl.Dbg("")
		Config.Wallet.RPCKey = n.CleanAndExpandPath(*r)
	}
	if getIfIs(ctx, "onetimetlskey", r) {
		log <- cl.Dbg("")
		Config.Wallet.OneTimeTLSKey = *r == "true"
	}
	if getIfIs(ctx, "enableservertls", r) {
		log <- cl.Dbg("")
		Config.Wallet.EnableServerTLS = *r == "true"
	}
	if getIfIs(ctx, "legacyrpclisteners", r) {
		log <- cl.Dbg("")
		pu.NormalizeAddresses(*r, "11046", &Config.Wallet.LegacyRPCListeners)
	}
	if getIfIs(ctx, "legacyrpcmaxclients", r) {
		log <- cl.Dbg("")
		var bt int
		if err := pu.ParseInteger(*r, "legacyrpcmaxclients", &bt); err != nil {
			log <- cl.Wrn(err.Error())
		} else {
			Config.Wallet.LegacyRPCMaxClients = int64(bt)
		}
	}
	if getIfIs(ctx, "legacyrpcmaxwebsockets", r) {
		log <- cl.Dbg("")
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
		log <- cl.Dbg("")
		Config.Wallet.Username = *r
	}
	if getIfIs(ctx, "password", r) {
		log <- cl.Dbg("")
		Config.Wallet.Password = *r
	}
	if getIfIs(ctx, "experimentalrpclisteners", r) {
		log <- cl.Dbg("")
		pu.NormalizeAddresses(*r, "11045", &Config.Wallet.ExperimentalRPCListeners)
	}
	if getIfIs(ctx, "datadir", r) {
		log <- cl.Dbg("")
		Config.Wallet.DataDir = *r
	}
	if getIfIs(ctx, "network", r) {
		log <- cl.Dbg("")
		switch *r {
		case "testnet":
			Config.Wallet.TestNet3, Config.Wallet.SimNet = true, false
			w.ActiveNet = &netparams.TestNet3Params
		case "simnet":
			Config.Wallet.TestNet3, Config.Wallet.SimNet = false, true
			w.ActiveNet = &netparams.SimNetParams
		default:
			Config.Wallet.TestNet3, Config.Wallet.SimNet = false, false
			w.ActiveNet = &netparams.MainNetParams
		}
	}

	// finished configuration

	logger.SetLogging(ctx)

	if ctx.Is("save") {
		log <- cl.Info{"saving config file to", cfgFile}
		j, err := json.MarshalIndent(Config, "", "  ")
		if err != nil {
			log <- cl.Error{"writing app config file", err}
		}
		j = append(j, '\n')
		log <- cl.Trace{"JSON formatted config file\n", string(j)}
		ioutil.WriteFile(cfgFile, j, 0600)
	}
}

// WriteConfig creates and writes the config file in the requested location
func WriteConfig(cfgFile string, c *ConfigAndLog) {
	log <- cl.Dbg("writing config")
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
	log <- cl.Dbg("writing default config")
	defCfg := DefaultConfig()
	defCfg.Wallet.ConfigFile = cfgFile
	j, err := json.MarshalIndent(defCfg, "", "  ")
	if err != nil {
		log <- cl.Error{"marshalling configuration", err}
		panic(err)
	}
	j = append(j, '\n')
	log <- cl.Trace{"JSON formatted config file\n", string(j)}
	err = ioutil.WriteFile(cfgFile, j, 0600)
	if err != nil {
		log <- cl.Error{"writing app config file", err}
		panic(err)
	}
	// if we are writing default config we also want to use it
	Config = defCfg
}

// DefaultConfig returns a default configuration
func DefaultConfig() *ConfigAndLog {
	log <- cl.Dbg("getting default config")
	return &ConfigAndLog{
		Wallet: &w.Config{
			NoInitialLoad:          false,
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

// FileExists reports whether the named file or directory exists.
func FileExists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
