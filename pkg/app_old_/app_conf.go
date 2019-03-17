package app_old

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"git.parallelcoin.io/dev/pod/cmd/ctl"
	"git.parallelcoin.io/dev/pod/cmd/node"
	"git.parallelcoin.io/dev/pod/cmd/shell"
	walletmain "git.parallelcoin.io/dev/pod/cmd/wallet"
	cl "git.parallelcoin.io/dev/pod/pkg/util/cl"
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
