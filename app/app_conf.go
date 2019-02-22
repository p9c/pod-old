package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"git.parallelcoin.io/pod/cmd/ctl"
	"git.parallelcoin.io/pod/cmd/node"
	"git.parallelcoin.io/pod/cmd/shell"
	walletmain "git.parallelcoin.io/pod/cmd/wallet"
	cl "git.parallelcoin.io/pod/pkg/clog"
	"git.parallelcoin.io/pod/pkg/fork"
	"github.com/tucnak/climax"
)

// ConfCfg is the settings that can be set to synchronise across all pod modules
type ConfCfg struct {
	DataDir          string
	ConfigFile       string
	NodeListeners    []string
	NodeRPCListeners []string
	WalletListeners  []string
	NodeUser         string
	NodePass         string
	WalletPass       string
	RPCKey           string
	RPCCert          string
	CAFile           string
	TLS              bool
	SkipVerify       bool
	Proxy            string
	ProxyUser        string
	ProxyPass        string
	Network          string
}

// ConfConfigs are the configurations for each app that are applied
type ConfConfigs struct {
	Ctl    ctl.Config
	Node   node.Config
	Wallet walletmain.Config
	Shell  shell.Config
}

const lH = "127.0.0.1:"

// ConfCommand is a command to send RPC queries to bitcoin RPC protocol server for node and wallet queries
var ConfCommand = climax.Command{
	Name:  "conf",
	Brief: "sets configurations common across modules",
	Help:  "automates synchronising common settings between servers and clients",
	Flags: []climax.Flag{
		t("version", "V", "show version number and quit"),
		t("init", "i", "resets configuration to defaults"),
		t("show", "s", "prints currently configuration"),

		f("createtest", "1", "create test configuration (set number to create max 10)"),
		f("testname", "test", "base name for test configurations"),
		f("testportbase", "21047", "base port number for test configurations"),

		s("datadir", "D", "~/.pod", "where to create the new profile"),

		f("nodelistener", node.DefaultListener,
			"main peer to peer address for apps that connect to the parallelcoin peer to peer network"),
		f("noderpclistener", node.DefaultRPCListener,
			"address where node listens for RPC"),
		f("walletlistener", walletmain.DefaultListener, "address where wallet listens for RPC"),
		s("user", "u", "user", "username for all the things"),
		s("pass", "P", "pa55word", "password for all the things"),
		s("walletpass", "", "w", "public password for wallet"),
		f("rpckey", walletmain.DefaultRPCKeyFile,
			"RPC server certificate key"),
		f("rpccert", walletmain.DefaultRPCCertFile,
			"RPC server certificate"),
		f("cafile", walletmain.DefaultCAFile,
			"RPC server certificate chain for validation"),
		f("tls", "false", "enable TLS"),
		f("skipverify", "false", "do not verify TLS certificates (not recommended!)"),
		f("proxy", "127.0.0.1:9050", "connect via SOCKS5 proxy"),
		f("proxyuser", "user", "username for proxy"),
		f("proxypass", "pa55word", "password for proxy"),

		f("network", "mainnet", "connect to [mainnet|testnet|regtestnet|simnet]"),
		s("debuglevel", "d", "info", "sets log level for those unspecified below"),
	},
	Examples: []climax.Example{
		{
			Usecase:     "--D test --init",
			Description: "creates a new data directory at test",
		},
	},
	Handle: func(ctx climax.Context) int {
		var dl, ct, tpb string
		var ok bool
		if dl, ok = ctx.Get("debuglevel"); ok {
			log <- cl.Tracef{
				"setting debug level %s",
				dl,
			}
			Log.SetLevel(dl)
			ll := GetAllSubSystems()
			for i := range ll {
				ll[i].SetLevel(dl)
			}
		}

		if ct, ok = ctx.Get("createtest"); ok {
			testname := "test"
			testnum := 1
			testportbase := 21047
			if err := ParseInteger(
				ct, "createtest", &testnum,
			); err != nil {
				log <- cl.Wrn(err.Error())
			}
			if tn, ok := ctx.Get("testname"); ok {
				testname = tn
			}
			if tpb, ok = ctx.Get("testportbase"); ok {
				if err := ParseInteger(
					tpb, "testportbase", &testportbase,
				); err != nil {
					log <- cl.Wrn(err.Error())
				}
			}
			// Generate a full set of default configs first
			var testConfigSet []ConfigSet
			for i := 0; i < testnum; i++ {
				tn := fmt.Sprintf("%s%d", testname, i)
				cs := GetDefaultConfs(tn)
				SyncToConfs(cs)
				testConfigSet = append(testConfigSet, *cs)
			}
			var ps []PortSet
			for i := 0; i < testnum; i++ {
				p := GenPortSet(testportbase + 100*i)
				ps = append(ps, *p)
			}
			// Set the correct listeners and add the correct addpeers entries
			for i, ts := range testConfigSet {
				// conf
				tc := ts.Conf
				tc.NodeListeners = []string{
					lH + ps[i].P2P,
				}
				tc.NodeRPCListeners = []string{
					lH + ps[i].NodeRPC,
				}
				tc.WalletListeners = []string{
					lH + ps[i].WalletRPC,
				}
				tc.TLS = false
				tc.Network = "testnet"
				// ctl
				tcc := ts.Ctl
				tcc.SimNet = false
				tcc.RPCServer = ts.Conf.NodeRPCListeners[0]
				tcc.TestNet3 = true
				tcc.TLS = false
				tcc.Wallet = ts.Conf.WalletListeners[0]
				// node
				tnn := ts.Node.Node
				for j := range ps {
					// add all other peers in the portset list
					if j != i {
						tnn.AddPeers = append(
							tnn.AddPeers,
							lH+ps[j].P2P,
						)
					}
				}
				tnn.Listeners = tc.NodeListeners
				tnn.RPCListeners = tc.NodeRPCListeners
				tnn.SimNet = false
				tnn.TestNet3 = true
				tnn.RegressionTest = false
				tnn.TLS = false
				// wallet
				tw := ts.Wallet.Wallet
				tw.EnableClientTLS = false
				tw.EnableServerTLS = false
				tw.LegacyRPCListeners = ts.Conf.WalletListeners
				tw.RPCConnect = tc.NodeRPCListeners[0]
				tw.SimNet = false
				tw.TestNet3 = true
				// shell
				tss := ts.Shell
				// shell/node
				tsn := tss.Node
				tsn.Listeners = tnn.Listeners
				tsn.RPCListeners = tnn.RPCListeners
				tsn.TestNet3 = true
				tsn.SimNet = true
				for j := range ps {
					// add all other peers in the portset list
					if j != i {
						tsn.AddPeers = append(
							tsn.AddPeers,
							lH+ps[j].P2P,
						)
					}
				}
				tsn.SimNet = false
				tsn.TestNet3 = true
				tsn.RegressionTest = false
				tsn.TLS = false
				// shell/wallet
				tsw := tss.Wallet
				tsw.EnableClientTLS = false
				tsw.EnableServerTLS = false
				tsw.LegacyRPCListeners = ts.Conf.WalletListeners
				tsw.RPCConnect = tcc.RPCServer
				tsw.SimNet = false
				tsw.TestNet3 = true
				// write to disk
				WriteConfigSet(&ts)
			}
			os.Exit(0)
		}

		confFile = DefaultDataDir + "/conf.json"
		if r, ok := ctx.Get("datadir"); ok {
			DefaultDataDir = r
			confFile = DefaultDataDir + "/conf.json"
		}
		confs = []string{
			DefaultDataDir + "/ctl/conf.json",
			DefaultDataDir + "/node/conf.json",
			DefaultDataDir + "/wallet/conf.json",
			DefaultDataDir + "/shell/conf.json",
		}
		for i := range confs {
			EnsureDir(confs[i])
		}
		EnsureDir(confFile)
		if ctx.Is("init") {
			WriteDefaultConfConfig(DefaultDataDir)
			WriteDefaultCtlConfig(DefaultDataDir)
			WriteDefaultNodeConfig(DefaultDataDir)
			WriteDefaultWalletConfig(DefaultDataDir)
			WriteDefaultShellConfig(DefaultDataDir)
		} else {
			if _, err := os.Stat(confFile); os.IsNotExist(err) {
				WriteDefaultConfConfig(DefaultDataDir)
			} else {
				cfgData, err := ioutil.ReadFile(confFile)
				if err != nil {
					WriteDefaultConfConfig(DefaultDataDir)
				}
				err = json.Unmarshal(cfgData, &ConfConfig)
				if err != nil {
					WriteDefaultConfConfig(DefaultDataDir)
				}
			}
		}
		configConf(&ctx, DefaultDataDir, node.DefaultPort)
		runConf()
		return 0
	},
	// Examples: []climax.Example{
	// 	{
	// 		Usecase:     "--nodeuser=user --nodepass=pa55word",
	// 		Description: "set the username and password for the node RPC",
	// 	},
	// },
	// Handle:
}

// ConfConfig is the configuration for this tool
var ConfConfig ConfCfg

// Confs is the central repository of all the other app configurations
var Confs ConfConfigs

var confFile = DefaultDataDir + "/conf"

var confs []string

// DefaultConfConfig returns a crispy fresh default conf configuration
func DefaultConfConfig(
	datadir string,
) *ConfCfg {

	u := GenKey()
	p := GenKey()
	return &ConfCfg{
		DataDir:          datadir,
		ConfigFile:       filepath.Join(datadir, "conf.json"),
		NodeListeners:    []string{node.DefaultListener},
		NodeRPCListeners: []string{node.DefaultRPCListener},
		WalletListeners:  []string{walletmain.DefaultListener},
		NodeUser:         u,
		NodePass:         p,
		WalletPass:       "",
		RPCCert:          filepath.Join(datadir, "rpc.cert"),
		RPCKey:           filepath.Join(datadir, "rpc.key"),
		CAFile: filepath.Join(
			datadir, walletmain.DefaultCAFilename),
		TLS:        false,
		SkipVerify: false,
		Proxy:      "",
		ProxyUser:  "",
		ProxyPass:  "",
		Network:    "mainnet",
	}
}

// WriteConfConfig creates and writes the config file in the requested location
func WriteConfConfig(
	cfg *ConfCfg,
) {

	j, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		panic(err.Error())
	}
	j = append(j, '\n')
	EnsureDir(cfg.ConfigFile)
	err = ioutil.WriteFile(cfg.ConfigFile, j, 0600)
	if err != nil {
		panic(err.Error())
	}
}

// WriteDefaultConfConfig creates and writes a default config file in the requested location
func WriteDefaultConfConfig(
	datadir string,
) {

	defCfg := DefaultConfConfig(datadir)
	j, err := json.MarshalIndent(defCfg, "", "  ")
	if err != nil {
		panic(err.Error())
	}
	j = append(j, '\n')
	EnsureDir(defCfg.ConfigFile)
	err = ioutil.WriteFile(defCfg.ConfigFile, j, 0600)
	if err != nil {
		panic(err.Error())
	}
	// if we are writing default config we also want to use it
	ConfConfig = *defCfg
}

// // cf is the list of flags and the default values stored in the Usage field
// var cf = GetFlags(ConfCommand)

func configConf(
	ctx *climax.Context,
	datadir,
	portbase string,
) {

	cs := GetDefaultConfs(datadir)
	SyncToConfs(cs)
	var r string
	var ok bool
	var listeners []string
	if r, ok = getIfIs(ctx, "nodelistener"); ok {
		NormalizeAddresses(r, portbase, &listeners)
		fmt.Println("nodelistener set to", listeners)
		ConfConfig.NodeListeners = listeners
		cs.Node.Node.Listeners = listeners
		cs.Shell.Node.Listeners = listeners
	}
	if r, ok = getIfIs(ctx, "noderpclistener"); ok {
		NormalizeAddresses(r, node.DefaultRPCPort, &listeners)
		fmt.Println("noderpclistener set to", listeners)
		ConfConfig.NodeRPCListeners = listeners
		cs.Node.Node.RPCListeners = listeners
		cs.Wallet.Wallet.RPCConnect = r
		cs.Shell.Node.RPCListeners = listeners
		cs.Shell.Wallet.RPCConnect = r
		cs.Ctl.RPCServer = r
	}
	if r, ok = getIfIs(ctx, "walletlistener"); ok {
		NormalizeAddresses(r, node.DefaultRPCPort, &listeners)
		fmt.Println("walletlistener set to", listeners)
		ConfConfig.WalletListeners = listeners
		cs.Wallet.Wallet.LegacyRPCListeners = listeners
		cs.Ctl.Wallet = r
		cs.Shell.Wallet.LegacyRPCListeners = listeners
	}
	if r, ok = getIfIs(ctx, "user"); ok {
		ConfConfig.NodeUser = r
		cs.Node.Node.RPCUser = r
		cs.Wallet.Wallet.PodUsername = r
		cs.Wallet.Wallet.Username = r
		cs.Shell.Node.RPCUser = r
		cs.Shell.Wallet.PodUsername = r
		cs.Shell.Wallet.Username = r
		cs.Ctl.RPCUser = r
	}
	if r, ok = getIfIs(ctx, "pass"); ok {
		ConfConfig.NodePass = r
		cs.Node.Node.RPCPass = r
		cs.Wallet.Wallet.PodPassword = r
		cs.Wallet.Wallet.Password = r
		cs.Shell.Node.RPCPass = r
		cs.Shell.Wallet.PodPassword = r
		cs.Shell.Wallet.Password = r
		cs.Ctl.RPCPass = r
	}
	if r, ok = getIfIs(ctx, "walletpass"); ok {
		ConfConfig.WalletPass = r
		cs.Wallet.Wallet.WalletPass = ConfConfig.WalletPass
		cs.Shell.Wallet.WalletPass = ConfConfig.WalletPass
	}

	if r, ok = getIfIs(ctx, "rpckey"); ok {
		r = node.CleanAndExpandPath(r)
		ConfConfig.RPCKey = r
		cs.Node.Node.RPCKey = r
		cs.Wallet.Wallet.RPCKey = r
		cs.Shell.Node.RPCKey = r
		cs.Shell.Wallet.RPCKey = r
	}
	if r, ok = getIfIs(ctx, "rpccert"); ok {
		r = node.CleanAndExpandPath(r)
		ConfConfig.RPCCert = r
		cs.Node.Node.RPCCert = r
		cs.Wallet.Wallet.RPCCert = r
		cs.Shell.Node.RPCCert = r
		cs.Shell.Wallet.RPCCert = r
	}
	if r, ok = getIfIs(ctx, "cafile"); ok {
		r = node.CleanAndExpandPath(r)
		ConfConfig.CAFile = r
		cs.Wallet.Wallet.CAFile = r
		cs.Shell.Wallet.CAFile = r
	}
	if r, ok = getIfIs(ctx, "tls"); ok {
		ConfConfig.TLS = r == "true"
		cs.Node.Node.TLS = ConfConfig.TLS
		cs.Wallet.Wallet.EnableClientTLS = ConfConfig.TLS
		cs.Shell.Node.TLS = ConfConfig.TLS
		cs.Shell.Wallet.EnableClientTLS = ConfConfig.TLS
		cs.Wallet.Wallet.EnableServerTLS = ConfConfig.TLS
		cs.Shell.Wallet.EnableServerTLS = ConfConfig.TLS
	}
	if r, ok = getIfIs(ctx, "skipverify"); ok {
		ConfConfig.SkipVerify = r == "true"
		cs.Ctl.TLSSkipVerify = r == "true"
	}
	if r, ok = getIfIs(ctx, "proxy"); ok {
		NormalizeAddresses(r, node.DefaultRPCPort, &listeners)
		ConfConfig.Proxy = r
		cs.Ctl.Proxy = ConfConfig.Proxy
		cs.Node.Node.Proxy = ConfConfig.Proxy
		cs.Wallet.Wallet.Proxy = ConfConfig.Proxy
		cs.Shell.Node.Proxy = ConfConfig.Proxy
		cs.Shell.Wallet.Proxy = ConfConfig.Proxy
	}
	if r, ok = getIfIs(ctx, "proxyuser"); ok {
		ConfConfig.ProxyUser = r
		cs.Ctl.ProxyUser = ConfConfig.ProxyUser
		cs.Node.Node.ProxyUser = ConfConfig.ProxyUser
		cs.Wallet.Wallet.ProxyUser = ConfConfig.ProxyUser
		cs.Shell.Node.ProxyUser = ConfConfig.ProxyUser
		cs.Shell.Wallet.ProxyUser = ConfConfig.ProxyUser
	}
	if r, ok = getIfIs(ctx, "proxypass"); ok {
		ConfConfig.ProxyPass = r
		cs.Ctl.ProxyPass = ConfConfig.ProxyPass
		cs.Node.Node.ProxyPass = ConfConfig.ProxyPass
		cs.Wallet.Wallet.ProxyPass = ConfConfig.ProxyPass
		cs.Shell.Node.ProxyPass = ConfConfig.ProxyPass
		cs.Shell.Wallet.ProxyPass = ConfConfig.ProxyPass
	}
	if r, ok = getIfIs(ctx, "network"); ok {
		r = strings.ToLower(r)
		switch r {
		case "mainnet", "testnet", "regtestnet", "simnet":
		default:
			r = "mainnet"
		}
		ConfConfig.Network = r
		fmt.Println("configured for", r, "network")
		switch r {
		case "mainnet":
			cs.Ctl.TestNet3 = false
			cs.Ctl.SimNet = false
			cs.Node.Node.TestNet3 = false
			cs.Node.Node.SimNet = false
			cs.Node.Node.RegressionTest = false
			cs.Wallet.Wallet.SimNet = false
			cs.Wallet.Wallet.TestNet3 = false
			cs.Shell.Node.TestNet3 = false
			cs.Shell.Node.RegressionTest = false
			cs.Shell.Node.SimNet = false
			cs.Shell.Wallet.TestNet3 = false
			cs.Shell.Wallet.SimNet = false
		case "testnet":
			fork.IsTestnet = true
			cs.Ctl.TestNet3 = true
			cs.Ctl.SimNet = false
			cs.Node.Node.TestNet3 = true
			cs.Node.Node.SimNet = false
			cs.Node.Node.RegressionTest = false
			cs.Wallet.Wallet.SimNet = false
			cs.Wallet.Wallet.TestNet3 = true
			cs.Shell.Node.TestNet3 = true
			cs.Shell.Node.RegressionTest = false
			cs.Shell.Node.SimNet = false
			cs.Shell.Wallet.TestNet3 = true
			cs.Shell.Wallet.SimNet = false
		case "regtestnet":
			cs.Ctl.TestNet3 = false
			cs.Ctl.SimNet = false
			cs.Node.Node.TestNet3 = false
			cs.Node.Node.SimNet = false
			cs.Node.Node.RegressionTest = true
			cs.Wallet.Wallet.SimNet = false
			cs.Wallet.Wallet.TestNet3 = false
			cs.Shell.Node.TestNet3 = false
			cs.Shell.Node.RegressionTest = true
			cs.Shell.Node.SimNet = false
			cs.Shell.Wallet.TestNet3 = false
			cs.Shell.Wallet.SimNet = false
		case "simnet":
			cs.Ctl.TestNet3 = false
			cs.Ctl.SimNet = true
			cs.Node.Node.TestNet3 = false
			cs.Node.Node.SimNet = true
			cs.Node.Node.RegressionTest = false
			cs.Wallet.Wallet.SimNet = true
			cs.Wallet.Wallet.TestNet3 = false
			cs.Shell.Node.TestNet3 = false
			cs.Shell.Node.RegressionTest = false
			cs.Shell.Node.SimNet = true
			cs.Shell.Wallet.TestNet3 = false
			cs.Shell.Wallet.SimNet = true
		}
	}

	WriteConfConfig(cs.Conf)
	// Now write the configs for all the others reading them and overwriting the changed values
	WriteCtlConfig(cs.Ctl)
	WriteNodeConfig(cs.Node)
	WriteWalletConfig(cs.Wallet)
	WriteShellConfig(cs.Shell)
	if ctx.Is("show") {
		j, err := json.MarshalIndent(cs.Conf, "", "  ")
		if err != nil {
			panic(err.Error())
		}
		fmt.Println(string(j))
	}
}
