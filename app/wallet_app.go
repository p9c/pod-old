package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	n "git.parallelcoin.io/pod/cmd/node"
	w "git.parallelcoin.io/pod/cmd/wallet"
	"git.parallelcoin.io/pod/cmd/wallet/netparams"
	cl "git.parallelcoin.io/pod/pkg/clog"
	"github.com/tucnak/climax"
)

// WalletCfg is the combined app and logging configuration data
type WalletCfg struct {
	Wallet *w.Config
	Levels map[string]string
}

// WalletConfig is the combined app and log levels configuration
var WalletConfig = DefaultWalletConfig()

// WalletCommand is a command to send RPC queries to bitcoin RPC protocol server for node and wallet queries
var WalletCommand = climax.Command{
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

		t("create", "", "create a new wallet if it does not exist"),
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
			ll := GetAllSubSystems()
			for i := range ll {
				ll[i].SetLevel(dl)
			}
		}
		log <- cl.Trc("starting wallet app")
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
			WriteDefaultWalletConfig(cfgFile)
		}
		log <- cl.Info{"loading configuration from", cfgFile}
		if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
			log <- cl.Wrn("configuration file does not exist, creating new one")
			WriteDefaultWalletConfig(cfgFile)
		} else {
			log <- cl.Debug{"reading app configuration from", cfgFile}
			cfgData, err := ioutil.ReadFile(cfgFile)
			if err != nil {
				log <- cl.Error{"reading app config file", err.Error()}
				WriteDefaultWalletConfig(cfgFile)
			}
			log <- cl.Tracef{"parsing app configuration\n%s", cfgData}
			err = json.Unmarshal(cfgData, &WalletConfig)
			if err != nil {
				log <- cl.Error{"parsing app config file", err.Error()}
				WriteDefaultWalletConfig(cfgFile)
			}
		}
		configWallet(&ctx, cfgFile)
		runWallet(ctx.Args)
		cl.Shutdown()
		return 0
	},
}

func configWallet(ctx *climax.Context, cfgFile string) {
	log <- cl.Trace{"configuring from command line flags ", os.Args}
	if ctx.Is("create") {
		log <- cl.Dbg("request to make new wallet")
		WalletConfig.Wallet.Create = true
	}
	if ctx.Is("createtemp") {
		log <- cl.Dbg("request to make temp wallet")
		WalletConfig.Wallet.CreateTemp = true
	}
	if r, ok := getIfIs(ctx, "appdatadir"); ok {
		log <- cl.Debug{"appdatadir set to", r}
		WalletConfig.Wallet.AppDataDir = n.CleanAndExpandPath(r)
	}
	if r, ok := getIfIs(ctx, "noinitialload"); ok {
		log <- cl.Dbg("no initial load requested")
		WalletConfig.Wallet.NoInitialLoad = r == "true"
	}
	if r, ok := getIfIs(ctx, "logdir"); ok {
		log <- cl.Dbg("")
		WalletConfig.Wallet.LogDir = n.CleanAndExpandPath(r)
	}
	if r, ok := getIfIs(ctx, "profile"); ok {
		log <- cl.Dbg("")
		NormalizeAddress(r, "3131", &WalletConfig.Wallet.Profile)
	}
	if r, ok := getIfIs(ctx, "gui"); ok {
		log <- cl.Dbg("")
		WalletConfig.Wallet.GUI = r == "true"
	}
	if r, ok := getIfIs(ctx, "walletpass"); ok {
		log <- cl.Dbg("")
		WalletConfig.Wallet.WalletPass = r
	}
	if r, ok := getIfIs(ctx, "rpcconnect"); ok {
		log <- cl.Dbg("")
		NormalizeAddress(r, "11048", &WalletConfig.Wallet.RPCConnect)
	}
	if r, ok := getIfIs(ctx, "cafile"); ok {
		log <- cl.Dbg("")
		WalletConfig.Wallet.CAFile = n.CleanAndExpandPath(r)
	}
	if r, ok := getIfIs(ctx, "enableclienttls"); ok {
		log <- cl.Dbg("")
		WalletConfig.Wallet.EnableClientTLS = r == "true"
	}
	if r, ok := getIfIs(ctx, "podusername"); ok {
		log <- cl.Dbg("")
		WalletConfig.Wallet.PodUsername = r
	}
	if r, ok := getIfIs(ctx, "podpassword"); ok {
		log <- cl.Dbg("")
		WalletConfig.Wallet.PodPassword = r
	}
	if r, ok := getIfIs(ctx, "proxy"); ok {
		log <- cl.Dbg("")
		NormalizeAddress(r, "11048", &WalletConfig.Wallet.Proxy)
	}
	if r, ok := getIfIs(ctx, "proxyuser"); ok {
		log <- cl.Dbg("")
		WalletConfig.Wallet.ProxyUser = r
	}
	if r, ok := getIfIs(ctx, "proxypass"); ok {
		log <- cl.Dbg("")
		WalletConfig.Wallet.ProxyPass = r
	}
	if r, ok := getIfIs(ctx, "rpccert"); ok {
		log <- cl.Dbg("")
		WalletConfig.Wallet.RPCCert = n.CleanAndExpandPath(r)
	}
	if r, ok := getIfIs(ctx, "rpckey"); ok {
		log <- cl.Dbg("")
		WalletConfig.Wallet.RPCKey = n.CleanAndExpandPath(r)
	}
	if r, ok := getIfIs(ctx, "onetimetlskey"); ok {
		log <- cl.Dbg("")
		WalletConfig.Wallet.OneTimeTLSKey = r == "true"
	}
	if r, ok := getIfIs(ctx, "enableservertls"); ok {
		log <- cl.Dbg("")
		WalletConfig.Wallet.EnableServerTLS = r == "true"
	}
	if r, ok := getIfIs(ctx, "legacyrpclisteners"); ok {
		log <- cl.Dbg("")
		NormalizeAddresses(r, "11046", &WalletConfig.Wallet.LegacyRPCListeners)
	}
	if r, ok := getIfIs(ctx, "legacyrpcmaxclients"); ok {
		log <- cl.Dbg("")
		var bt int
		if err := ParseInteger(r, "legacyrpcmaxclients", &bt); err != nil {
			log <- cl.Wrn(err.Error())
		} else {
			WalletConfig.Wallet.LegacyRPCMaxClients = int64(bt)
		}
	}
	if r, ok := getIfIs(ctx, "legacyrpcmaxwebsockets"); ok {
		log <- cl.Dbg("")
		_, err := fmt.Sscanf(r, "%d", WalletConfig.Wallet.LegacyRPCMaxWebsockets)
		if err != nil {
			log <- cl.Errorf{
				"malformed legacyrpcmaxwebsockets: `%s` leaving set at `%d`",
				r, WalletConfig.Wallet.LegacyRPCMaxWebsockets,
			}
		}
	}
	if r, ok := getIfIs(ctx, "username"); ok {
		log <- cl.Dbg("")
		WalletConfig.Wallet.Username = r
	}
	if r, ok := getIfIs(ctx, "password"); ok {
		log <- cl.Dbg("")
		WalletConfig.Wallet.Password = r
	}
	if r, ok := getIfIs(ctx, "experimentalrpclisteners"); ok {
		log <- cl.Dbg("")
		NormalizeAddresses(r, "11045", &WalletConfig.Wallet.ExperimentalRPCListeners)
	}
	if r, ok := getIfIs(ctx, "datadir"); ok {
		log <- cl.Dbg("")
		WalletConfig.Wallet.DataDir = r
	}
	if r, ok := getIfIs(ctx, "network"); ok {
		log <- cl.Dbg("")
		switch r {
		case "testnet":
			WalletConfig.Wallet.TestNet3, WalletConfig.Wallet.SimNet = true, false
			w.ActiveNet = &netparams.TestNet3Params
		case "simnet":
			WalletConfig.Wallet.TestNet3, WalletConfig.Wallet.SimNet = false, true
			w.ActiveNet = &netparams.SimNetParams
		default:
			WalletConfig.Wallet.TestNet3, WalletConfig.Wallet.SimNet = false, false
			w.ActiveNet = &netparams.MainNetParams
		}
	}

	// finished configuration

	SetLogging(ctx)

	if ctx.Is("save") {
		log <- cl.Info{"saving config file to", cfgFile}
		j, err := json.MarshalIndent(WalletConfig, "", "  ")
		if err != nil {
			log <- cl.Error{"writing app config file", err}
		}
		j = append(j, '\n')
		log <- cl.Trace{"JSON formatted config file\n", string(j)}
		ioutil.WriteFile(cfgFile, j, 0600)
	}
}

// WriteWalletConfig creates and writes the config file in the requested location
func WriteWalletConfig(cfgFile string, c *WalletCfg) {
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

// WriteDefaultWalletConfig creates and writes a default config to the requested location
func WriteDefaultWalletConfig(cfgFile string) {
	defCfg := DefaultWalletConfig()
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
	WalletConfig = defCfg
}

// DefaultWalletConfig returns a default configuration
func DefaultWalletConfig() *WalletCfg {
	log <- cl.Dbg("getting default config")
	return &WalletCfg{
		Wallet: &w.Config{
			NoInitialLoad:          false,
			ConfigFile:             w.DefaultConfigFile,
			DataDir:                w.DefaultDataDir,
			AppDataDir:             w.DefaultAppDataDir,
			LogDir:                 w.DefaultLogDir,
			RPCKey:                 w.DefaultRPCKeyFile,
			RPCCert:                w.DefaultRPCCertFile,
			WalletPass:             "",
			CAFile:                 "",
			LegacyRPCMaxClients:    w.DefaultRPCMaxClients,
			LegacyRPCMaxWebsockets: w.DefaultRPCMaxWebsockets,
		},
		Levels: GetDefaultLogLevelsConfig(),
	}
}
