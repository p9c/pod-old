package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	n "git.parallelcoin.io/pod/cmd/node"
	w "git.parallelcoin.io/pod/cmd/wallet"
	walletmain "git.parallelcoin.io/pod/cmd/wallet"
	cl "git.parallelcoin.io/pod/pkg/clog"
	"git.parallelcoin.io/pod/pkg/netparams"
	"github.com/tucnak/climax"
)

// WalletCfg is the combined app and logging configuration data
type WalletCfg struct {
	Wallet    *w.Config
	Levels    map[string]string
	activeNet *netparams.Params
}

// WalletCommand is a command to send RPC queries to bitcoin RPC protocol server for node and wallet queries
var WalletCommand = climax.Command{
	Name:  "wallet",
	Brief: "parallelcoin wallet",
	Help:  "check balances, make payments, manage contacts",
	Flags: []climax.Flag{
		t("version", "V", "show version number and quit"),

		s("configfile", "C", w.DefaultConfigFilename,
			"path to configuration file"),
		s("datadir", "D", w.DefaultDataDir,
			"set the pod base directory"),
		f("appdatadir", w.DefaultAppDataDir, "set app data directory for wallet, configuration and logs"),

		t("init", "i", "resets configuration to defaults"),
		t("save", "S", "saves current flags into configuration"),

		t("createtemp", "", "create temporary wallet (pass=walletpass) requires --datadir"),

		t("gui", "G", "launch GUI"),
		f("rpcconnect", n.DefaultRPCListener, "connect to the RPC of a parallelcoin node for chain queries"),

		f("podusername", "user", "username for node RPC authentication"),
		f("podpassword", "pa55word", "password for node RPC authentication"),

		f("walletpass", "", "the public wallet password - only required if the wallet was created with one"),

		f("noinitialload", "false", "defer wallet load to be triggered by RPC"),
		f("network", "mainnet", "connect to (mainnet|testnet|regtestnet|simnet)"),

		f("profile", "false", "enable HTTP profiling on given port (1024-65536)"),

		f("rpccert", w.DefaultRPCCertFile,
			"file containing the RPC tls certificate"),
		f("rpckey", w.DefaultRPCKeyFile,
			"file containing RPC TLS key"),
		f("onetimetlskey", "false", "generate a new TLS certpair don't save key"),
		f("cafile", w.DefaultCAFile, "certificate authority for custom TLS CA"),
		f("enableclienttls", "false", "enable TLS for the RPC client"),
		f("enableservertls", "false", "enable TLS on wallet RPC server"),

		f("proxy", "", "proxy address for outbound connections"),
		f("proxyuser", "", "username for proxy server"),
		f("proxypass", "", "password for proxy server"),

		f("legacyrpclisteners", w.DefaultListener, "add a listener for the legacy RPC"),
		f("legacyrpcmaxclients", fmt.Sprint(w.DefaultRPCMaxClients),
			"max connections for legacy RPC"),
		f("legacyrpcmaxwebsockets", fmt.Sprint(w.DefaultRPCMaxWebsockets),
			"max websockets for legacy RPC"),

		f("username", "user",
			"username for wallet RPC when podusername is empty"),
		f("password", "pa55word",
			"password for wallet RPC when podpassword is omitted"),
		f("experimentalrpclisteners", "",
			"listener for experimental rpc"),

		s("debuglevel", "d", "info", "sets debuglevel, specify per-library below"),

		l("lib-addrmgr"), l("lib-blockchain"), l("lib-connmgr"), l("lib-database-ffldb"), l("lib-database"), l("lib-mining-cpuminer"), l("lib-mining"), l("lib-netsync"), l("lib-peer"), l("lib-rpcclient"), l("lib-txscript"), l("node"), l("node-mempool"), l("spv"), l("wallet"), l("wallet-chain"), l("wallet-legacyrpc"), l("wallet-rpcserver"), l("wallet-tx"), l("wallet-votingpool"), l("wallet-waddrmgr"), l("wallet-wallet"), l("wallet-wtxmgr"),
	},
	// Examples: []climax.Example{
	// 	{
	// 		Usecase:     "--init --rpcuser=user --rpcpass=pa55word --save",
	// 		Description: "resets the configuration file to default, sets rpc username and password and saves the changes to config after parsing",
	// 	},
	// },
}

// WalletConfig is the combined app and log levels configuration
var WalletConfig = DefaultWalletConfig(w.DefaultConfigFile)

// wf is the list of flags and the default values stored in the Usage field
var wf = GetFlags(WalletCommand)

func init() {
	// Loads after the var clauses run
	WalletCommand.Handle = func(ctx climax.Context) int {

		Log.SetLevel("off")
		var dl string
		var ok bool
		if dl, ok = ctx.Get("debuglevel"); ok {
			Log.SetLevel(dl)
			ll := GetAllSubSystems()
			for i := range ll {
				ll[i].SetLevel(dl)
			}
		}
		log <- cl.Tracef{"setting debug level %s", dl}
		log <- cl.Trc("starting wallet app")
		log <- cl.Debugf{"pod/wallet version %s", w.Version()}
		if ctx.Is("version") {
			fmt.Println("pod/wallet version", w.Version())
			cl.Shutdown()
		}
		var datadir, cfgFile string
		if datadir, ok = ctx.Get("datadir"); !ok {
			datadir = w.DefaultDataDir
		}
		cfgFile = filepath.Join(filepath.Join(datadir, "node"), "conf.json")
		log <- cl.Debug{"DataDir", datadir, "cfgFile", cfgFile}
		if cfgFile, ok = ctx.Get("configfile"); !ok {
			cfgFile = filepath.Join(
				filepath.Join(datadir, "wallet"), w.DefaultConfigFilename)
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
			WalletConfig.activeNet = &netparams.MainNetParams
			if WalletConfig.Wallet.TestNet3 {
				WalletConfig.activeNet = &netparams.TestNet3Params
			}
			if WalletConfig.Wallet.SimNet {
				WalletConfig.activeNet = &netparams.SimNetParams
			}
		}

		configWallet(WalletConfig.Wallet, &ctx, cfgFile)
		if dl, ok = ctx.Get("debuglevel"); ok {
			for i := range WalletConfig.Levels {
				WalletConfig.Levels[i] = dl
			}
		}
		fmt.Println("running wallet on", WalletConfig.activeNet.Name)
		runWallet(WalletConfig.Wallet, WalletConfig.activeNet)
		return 0
	}
}

func configWallet(wc *w.Config, ctx *climax.Context, cfgFile string) {
	log <- cl.Trace{"configuring from command line flags ", os.Args}
	if ctx.Is("createtemp") {
		log <- cl.Dbg("request to make temp wallet")
		wc.CreateTemp = true
	}
	if r, ok := getIfIs(ctx, "appdatadir"); ok {
		log <- cl.Debug{"appdatadir set to", r}
		wc.AppDataDir = n.CleanAndExpandPath(r)
	}
	if r, ok := getIfIs(ctx, "logdir"); ok {
		wc.LogDir = n.CleanAndExpandPath(r)
	}
	if r, ok := getIfIs(ctx, "profile"); ok {
		NormalizeAddress(r, "3131", &wc.Profile)
	}
	if r, ok := getIfIs(ctx, "walletpass"); ok {
		wc.WalletPass = r
	}
	if r, ok := getIfIs(ctx, "rpcconnect"); ok {
		NormalizeAddress(r, "11048", &wc.RPCConnect)
	}
	if r, ok := getIfIs(ctx, "cafile"); ok {
		wc.CAFile = n.CleanAndExpandPath(r)
	}
	if r, ok := getIfIs(ctx, "enableclienttls"); ok {
		wc.EnableClientTLS = r == "true"
	}
	if r, ok := getIfIs(ctx, "podusername"); ok {
		wc.PodUsername = r
	}
	if r, ok := getIfIs(ctx, "podpassword"); ok {
		wc.PodPassword = r
	}
	if r, ok := getIfIs(ctx, "proxy"); ok {
		NormalizeAddress(r, "11048", &wc.Proxy)
	}
	if r, ok := getIfIs(ctx, "proxyuser"); ok {
		wc.ProxyUser = r
	}
	if r, ok := getIfIs(ctx, "proxypass"); ok {
		wc.ProxyPass = r
	}
	if r, ok := getIfIs(ctx, "rpccert"); ok {
		wc.RPCCert = n.CleanAndExpandPath(r)
	}
	if r, ok := getIfIs(ctx, "rpckey"); ok {
		wc.RPCKey = n.CleanAndExpandPath(r)
	}
	if r, ok := getIfIs(ctx, "onetimetlskey"); ok {
		wc.OneTimeTLSKey = r == "true"
	}
	if r, ok := getIfIs(ctx, "enableservertls"); ok {
		wc.EnableServerTLS = r == "true"
	}
	if r, ok := getIfIs(ctx, "legacyrpclisteners"); ok {
		NormalizeAddresses(r, "11046", &wc.LegacyRPCListeners)
	}
	if r, ok := getIfIs(ctx, "legacyrpcmaxclients"); ok {
		var bt int
		if err := ParseInteger(r, "legacyrpcmaxclients", &bt); err != nil {
			log <- cl.Wrn(err.Error())
		} else {
			wc.LegacyRPCMaxClients = int64(bt)
		}
	}
	if r, ok := getIfIs(ctx, "legacyrpcmaxwebsockets"); ok {
		_, err := fmt.Sscanf(r, "%d", wc.LegacyRPCMaxWebsockets)
		if err != nil {
			log <- cl.Errorf{
				"malformed legacyrpcmaxwebsockets: `%s` leaving set at `%d`",
				r, wc.LegacyRPCMaxWebsockets,
			}
		}
	}
	if r, ok := getIfIs(ctx, "username"); ok {
		wc.Username = r
	}
	if r, ok := getIfIs(ctx, "password"); ok {
		wc.Password = r
	}
	if r, ok := getIfIs(ctx, "experimentalrpclisteners"); ok {
		NormalizeAddresses(r, "11045", &wc.ExperimentalRPCListeners)
	}
	if r, ok := getIfIs(ctx, "datadir"); ok {
		wc.DataDir = r
	}
	if r, ok := getIfIs(ctx, "network"); ok {
		switch r {
		case "testnet":
			wc.TestNet3, wc.SimNet = true, false
			WalletConfig.activeNet = &netparams.TestNet3Params
		case "simnet":
			wc.TestNet3, wc.SimNet = false, true
			WalletConfig.activeNet = &netparams.SimNetParams
		default:
			wc.TestNet3, wc.SimNet = false, false
			WalletConfig.activeNet = &netparams.MainNetParams
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
func WriteWalletConfig(c *WalletCfg) {
	log <- cl.Dbg("writing config")
	j, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		panic(err.Error())
	}
	j = append(j, '\n')
	EnsureDir(c.Wallet.ConfigFile)
	err = ioutil.WriteFile(c.Wallet.ConfigFile, j, 0600)
	if err != nil {
		panic(err.Error())
	}
}

// WriteDefaultWalletConfig creates and writes a default config to the requested location
func WriteDefaultWalletConfig(datadir string) {
	defCfg := DefaultWalletConfig(datadir)
	j, err := json.MarshalIndent(defCfg, "", "  ")
	if err != nil {
		log <- cl.Error{"marshalling configuration", err}
		panic(err)
	}
	j = append(j, '\n')
	EnsureDir(defCfg.Wallet.ConfigFile)
	log <- cl.Trace{"JSON formatted config file\n", string(j)}
	EnsureDir(defCfg.Wallet.ConfigFile)
	err = ioutil.WriteFile(defCfg.Wallet.ConfigFile, j, 0600)
	if err != nil {
		log <- cl.Error{"writing app config file", err}
		panic(err)
	}
	// if we are writing default config we also want to use it
	WalletConfig = defCfg
}

// DefaultWalletConfig returns a default configuration
func DefaultWalletConfig(datadir string) *WalletCfg {
	log <- cl.Dbg("getting default config")
	appdatadir := filepath.Join(datadir, w.DefaultAppDataDirname)
	return &WalletCfg{
		Wallet: &w.Config{
			ConfigFile: filepath.Join(
				appdatadir, w.DefaultConfigFilename),
			DataDir:         datadir,
			AppDataDir:      appdatadir,
			RPCConnect:      n.DefaultRPCListener,
			PodUsername:     "user",
			PodPassword:     "pa55word",
			WalletPass:      "",
			NoInitialLoad:   false,
			RPCCert:         filepath.Join(datadir, "rpc.cert"),
			RPCKey:          filepath.Join(datadir, "rpc.key"),
			CAFile:          walletmain.DefaultCAFile,
			EnableClientTLS: false,
			EnableServerTLS: false,
			Proxy:           "",
			ProxyUser:       "",
			ProxyPass:       "",
			LegacyRPCListeners: []string{
				w.DefaultListener,
			},
			LegacyRPCMaxClients:      w.DefaultRPCMaxClients,
			LegacyRPCMaxWebsockets:   w.DefaultRPCMaxWebsockets,
			Username:                 "user",
			Password:                 "pa55word",
			ExperimentalRPCListeners: []string{},
		},
		Levels:    GetDefaultLogLevelsConfig(),
		activeNet: &netparams.MainNetParams,
	}
}
