package walletrun

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"

	"git.parallelcoin.io/pod/lib/clog"
	n "git.parallelcoin.io/pod/module/node"
	"git.parallelcoin.io/pod/module/spv"
	w "git.parallelcoin.io/pod/module/wallet"
	"git.parallelcoin.io/pod/module/wallet/cfgutil"
	"git.parallelcoin.io/pod/module/wallet/legacy/keystore"
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
			w.ActiveNet = &netparams.TestNet3Params
		case "simnet":
			Config.Wallet.TestNet3, Config.Wallet.SimNet = false, true
			w.ActiveNet = &netparams.SimNetParams
		default:
			Config.Wallet.TestNet3, Config.Wallet.SimNet = false, false
			w.ActiveNet = &netparams.MainNetParams
		}
	}

	// Exit if you try to use a simulation wallet with a standard data directory.
	if !(ctx.Is("appdatadir") || ctx.Is("datadir")) && Config.Wallet.CreateTemp {
		fmt.Fprintln(os.Stderr, "Tried to create a temporary simulation wallet, but failed to specify data directory!")
		os.Exit(0)
	}

	// Exit if you try to use a simulation wallet on anything other than simnet or testnet3.
	if !Config.Wallet.SimNet && Config.Wallet.CreateTemp {
		fmt.Fprintln(os.Stderr,
			"Tried to create a temporary simulation wallet for network other than simnet!",
		)
		os.Exit(0)
	}

	// // Ensure the wallet exists or create it when the create flag is set.
	netDir := w.NetworkDir(Config.Wallet.AppDataDir, w.ActiveNet.Params)
	dbPath := filepath.Join(netDir, w.WalletDbName)

	if ctx.Is("createtemp") && ctx.Is("create") {
		err := fmt.Errorf("The flags --create and --createtemp can not " +
			"be specified together. Use --help for more information.")
		log <- cl.Error{err}
		cl.Shutdown()
	}

	dbFileExists, err := FileExists(dbPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		cl.Shutdown()
	}

	if ctx.Is("createtemp") {
		tempWalletExists := false

		if dbFileExists {
			str := fmt.Sprintf(
				"The wallet already exists. Loading this wallet instead.",
			)
			fmt.Fprintln(os.Stdout, str)
			tempWalletExists = true
		}

		if !tempWalletExists {
			// Perform the initial wallet creation wizard.
			if err := w.CreateSimulationWallet(Config.Wallet); err != nil {
				log <- cl.Error{"Unable to create wallet:", err}
			}
		}
	} else if ctx.Is("create") {
		// Error if the create flag is set and the wallet already exists.
		if dbFileExists {
			log <- cl.Fatal{
				"The wallet database file `%v` already exists.", dbPath,
			}
			cl.Shutdown()
		}
	}

	// Ensure the data directory for the network exists.
	if err := pu.CheckCreateDir(netDir); err != nil {
		log <- cl.Error{err}
		cl.Shutdown()
	}

	// Perform the initial wallet creation wizard.
	if err := w.CreateWallet(Config.Wallet); err != nil {
		log <- cl.Fatal{"Unable to create wallet:", err}
		// Created successfully, so exit now with success.
		os.Exit(0)
	} else if !dbFileExists && !Config.Wallet.NoInitialLoad {
		keystorePath := filepath.Join(netDir, keystore.Filename)
		keystoreExists, err := cfgutil.FileExists(keystorePath)
		if err != nil {
			log <- cl.Error{err}
			cl.Shutdown()
		}
		if !keystoreExists {
			// err = fmt.Errorf("The wallet does not exist.  Run with the " +
			// "--create option to initialize and create it...")
			// Ensure the data directory for the network exists.
			// fmt.Println("Existing wallet not found in", config.Wallet.ConfigFile.Value)
			if err := pu.CheckCreateDir(netDir); err != nil {
				log <- cl.Error{err}
				cl.Shutdown()
			}

			// Perform the initial wallet creation wizard.
			if err := w.CreateWallet(Config.Wallet); err != nil {
				log <- cl.Error{"Unable to create wallet:", err}
			}

			// Created successfully, so exit now with success.
			cl.Shutdown()

		} else {
			err = fmt.Errorf(
				"the wallet is in legacy format - run with the --create option to import it",
			)
		}
		log <- cl.Error{err}
		cl.Shutdown()
	}

	if Config.Wallet.UseSPV {
		spv.MaxPeers = Config.Wallet.MaxPeers
		spv.BanDuration = Config.Wallet.BanDuration
		spv.BanThreshold = Config.Wallet.BanThreshold
	} else if Config.Wallet.RPCConnect == "" {
		Config.Wallet.RPCConnect = net.JoinHostPort("localhost", w.ActiveNet.RPCClientPort)
	}

	// Add default port to connect flag if missing.
	Config.Wallet.RPCConnect, err = cfgutil.NormalizeAddress(
		Config.Wallet.RPCConnect, w.ActiveNet.RPCClientPort,
	)
	if err != nil {
		log <- cl.Error{"invalid rpcconnect network address: %v\n", err}
		cl.Shutdown()
	}

	if Config.Wallet.EnableClientTLS {
		// If CAFile is unset, choose either the copy or local pod cert.
		if !ctx.Is("cafile") {
			Config.Wallet.CAFile = filepath.Join(Config.Wallet.AppDataDir, w.DefaultCAFilename)
			// If the CA copy does not exist, check if we're connecting to
			// a local pod and switch to its RPC cert if it exists.
			certExists, err := cfgutil.FileExists(Config.Wallet.CAFile)
			if err != nil {
				log <- cl.Error{err}
				cl.Shutdown()
			}
			if !certExists {
				podCertExists, err := cfgutil.FileExists(w.DefaultCAFile)
				if err != nil {
					log <- cl.Error{err}
				}
				if podCertExists {
					Config.Wallet.CAFile = w.DefaultCAFile
				}
			}
		}
	}

	// Only set default RPC listeners when there are no listeners set for the experimental RPC server.  This is required to prevent the old RPC server from sharing listen addresses, since it is impossible to remove defaults from go-flags slice options without assigning specific behavior to a particular string.
	if len(Config.Wallet.ExperimentalRPCListeners) == 0 && len(Config.Wallet.LegacyRPCListeners) == 0 {
		addrs, err := net.LookupHost("localhost")
		if err != nil {
			cl.Shutdown()
		}
		Config.Wallet.LegacyRPCListeners = make([]string, 0, len(addrs))
		for _, addr := range addrs {
			addr = net.JoinHostPort(addr, w.ActiveNet.RPCServerPort)
			Config.Wallet.LegacyRPCListeners = append(Config.Wallet.LegacyRPCListeners, addr)
		}
	}

	// Add default port to all rpc listener addresses if needed and remove duplicate addresses.
	Config.Wallet.LegacyRPCListeners, err = cfgutil.NormalizeAddresses(
		Config.Wallet.LegacyRPCListeners, w.ActiveNet.RPCServerPort)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"Invalid network address in legacy RPC listeners: %v\n", err)
		cl.Shutdown()
	}
	Config.Wallet.ExperimentalRPCListeners, err = cfgutil.NormalizeAddresses(
		Config.Wallet.ExperimentalRPCListeners, w.ActiveNet.RPCServerPort)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"Invalid network address in RPC listeners: %v\n", err)
		cl.Shutdown()
	}

	// Both RPC servers may not listen on the same interface/port.
	if len(Config.Wallet.LegacyRPCListeners) > 0 && len(Config.Wallet.ExperimentalRPCListeners) > 0 {
		seenAddresses := make(map[string]struct{}, len(Config.Wallet.LegacyRPCListeners))
		for _, addr := range Config.Wallet.LegacyRPCListeners {
			seenAddresses[addr] = struct{}{}
		}
		for _, addr := range Config.Wallet.ExperimentalRPCListeners {
			_, seen := seenAddresses[addr]
			if seen {
				log <- cl.Errorf{
					"Address `%s` may not be used as a listener address for both RPC servers", addr}
				cl.Shutdown()
			}
		}
	}

	// Expand environment variable and leading ~ for filepaths.
	Config.Wallet.CAFile = n.CleanAndExpandPath(Config.Wallet.CAFile)
	Config.Wallet.RPCCert = n.CleanAndExpandPath(Config.Wallet.RPCCert)
	Config.Wallet.RPCKey = n.CleanAndExpandPath(Config.Wallet.RPCKey)

	// If the pod username or password are unset, use the same auth as for the client.  The two settings were previously shared for pod and client auth, so this avoids breaking backwards compatibility while allowing users to use different auth settings for pod and wallet.
	if Config.Wallet.PodUsername == "" {
		Config.Wallet.PodUsername = Config.Wallet.Username
	}
	if Config.Wallet.PodPassword == "" {
		Config.Wallet.PodPassword = Config.Wallet.Password
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
		panic(err)
	}
	j = append(j, '\n')
	log <- cl.Trace{"JSON formatted config file\n", string(j)}
	err = ioutil.WriteFile(cfgFile, j, 0600)
	if err != nil {
		log <- cl.Error{"writing app config file", err.Error()}
		panic(err)
	}
	// if we are writing default config we also want to use it
	Config = defCfg
}

// DefaultConfig returns a default configuration
func DefaultConfig() *ConfigAndLog {
	return &ConfigAndLog{
		Wallet: &w.Config{
			NoInitialLoad:          true,
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
