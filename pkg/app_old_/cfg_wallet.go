package app_old

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	n "git.parallelcoin.io/dev/pod/cmd/node"
	walletmain "git.parallelcoin.io/dev/pod/cmd/wallet"
	netparams "git.parallelcoin.io/dev/pod/pkg/chain/config/params"
	"git.parallelcoin.io/dev/pod/pkg/chain/fork"
	cl "git.parallelcoin.io/dev/pod/pkg/util/cl"
	"github.com/tucnak/climax"
)

// DefaultWalletConfig returns a default configuration
func DefaultWalletConfig(
	datadir string,
) (
	wc *WalletCfg,

) {

	log <- cl.Dbg("getting default config")

	appdatadir := filepath.Join(datadir, walletmain.DefaultAppDataDirname)

	return &WalletCfg{

		Wallet: &walletmain.Config{

			ConfigFile: filepath.Join(
				appdatadir, walletmain.DefaultConfigFilename),
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

				walletmain.DefaultListener,
			},

			LegacyRPCMaxClients:      walletmain.DefaultRPCMaxClients,
			LegacyRPCMaxWebsockets:   walletmain.DefaultRPCMaxWebsockets,
			Username:                 "user",
			Password:                 "pa55word",
			ExperimentalRPCListeners: []string{},
		},

		Levels:    GetDefaultLogLevelsConfig(),
		activeNet: &netparams.MainNetParams,
	}

}

// WriteDefaultWalletConfig creates and writes a default config to the requested location
func WriteDefaultWalletConfig(
	datadir string,

) {

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

// WriteWalletConfig creates and writes the config file in the requested location
func WriteWalletConfig(
	c *WalletCfg,

) {

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

func configWallet(
	wc *walletmain.Config,
	ctx *climax.Context,
	cfgFile string,

) {

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

			fork.IsTestnet = true

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

		e := ioutil.WriteFile(cfgFile, j, 0600)

		if e != nil {

			log <- cl.Error{

				"error writing configuration file:", e,
			}

		}

	}

}

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

		log <- cl.Debugf{"pod/wallet version %s", walletmain.Version()}

		if ctx.Is("version") {

			fmt.Println("pod/wallet version", walletmain.Version())
			return 0
		}

		var datadir, cfgFile string

		if datadir, ok = ctx.Get("datadir"); !ok {

			datadir = walletmain.DefaultDataDir
		}

		cfgFile = filepath.Join(filepath.Join(datadir, "node"), "conf.json")

		log <- cl.Debug{"DataDir", datadir, "cfgFile", cfgFile}

		if cfgFile, ok = ctx.Get("configfile"); !ok {

			cfgFile = filepath.Join(
				filepath.Join(datadir, "wallet"), walletmain.DefaultConfigFilename)
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
