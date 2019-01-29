package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"git.parallelcoin.io/pod/cmd/ctl"
	"git.parallelcoin.io/pod/cmd/node"
	walletmain "git.parallelcoin.io/pod/cmd/wallet"
	"git.parallelcoin.io/pod/pkg/wallet"
	"github.com/tucnak/climax"
)

var confFile = DefaultDataDir + "/conf"

// ConfCfg is the settings that can be set to synchronise across all pod modules
type ConfCfg struct {
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

// ConfConfig is the configuration for this tool
var ConfConfig ConfCfg

// Confs is the central repository of all the other app configurations
var Confs ConfConfigs

// ConfConfigs are the configurations for each app that are applied
type ConfConfigs struct {
	Ctl    ctl.Config
	Node   node.Config
	Wallet walletmain.Config
	Shell  ShellCfg
}

// ConfCommand is a command to send RPC queries to bitcoin RPC protocol server for node and wallet queries
var ConfCommand = climax.Command{
	Name:  "conf",
	Brief: "sets configurations common across modules",
	Help:  "automates synchronising common settings between servers and clients",
	Flags: []climax.Flag{
		t("version", "V", "show version number and quit"),
		t("init", "i", "resets configuration to defaults"),
		t("show", "s", "prints currently configuration"),

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
	},
	// Examples: []climax.Example{
	// 	{
	// 		Usecase:     "--nodeuser=user --nodepass=pa55word",
	// 		Description: "set the username and password for the node RPC",
	// 	},
	// },
	// Handle:
}

func init() {
	ConfCommand.Handle = func(ctx climax.Context) int {
		if ctx.Is("init") {
			WriteDefaultConfConfig(confFile)
		} else {
			if _, err := os.Stat(confFile); os.IsNotExist(err) {
				WriteDefaultConfConfig(confFile)
			} else {
				cfgData, err := ioutil.ReadFile(confFile)
				if err != nil {
					WriteDefaultConfConfig(confFile)
				}
				err = json.Unmarshal(cfgData, &ConfConfig)
				if err != nil {
					WriteDefaultConfConfig(confFile)
				}
			}
		}
		configConf(&ctx, confFile)
		runConf()
		return 0
	}
}

// cf is the list of flags and the default values stored in the Usage field
var cf = GetFlags(ConfCommand)

func configConf(ctx *climax.Context, cfgFile string) {
	// First load all of the module configurations and unmarshal into their structs
	confs := []string{
		DefaultDataDir + "/ctl/conf",
		DefaultDataDir + "/node/conf",
		DefaultDataDir + "/wallet/conf",
		DefaultDataDir + "/shell/conf",
	}

	// If we can't parse the config files we just reset them to default

	ctlCfg := *DefaultCtlConfig()
	if _, err := os.Stat(confs[0]); os.IsNotExist(err) {
		WriteDefaultConfConfig(confs[0])
	} else {
		ctlCfgData, err := ioutil.ReadFile(confs[0])
		if err != nil {
			WriteDefaultCtlConfig(confs[0])
		} else {
			err = json.Unmarshal(ctlCfgData, &ctlCfg)
			if err != nil {
				WriteDefaultCtlConfig(confs[0])
			}
		}
	}

	nodeCfg := *DefaultNodeConfig()
	if _, err := os.Stat(confs[1]); os.IsNotExist(err) {
		WriteDefaultNodeConfig(confs[1])
	} else {
		nodeCfgData, err := ioutil.ReadFile(confs[1])
		if err != nil {
			WriteDefaultNodeConfig(confs[1])
		} else {
			err = json.Unmarshal(nodeCfgData, &nodeCfg)
			if err != nil {
				WriteDefaultNodeConfig(confs[1])
			}
		}
	}

	walletCfg := *DefaultWalletConfig()
	if _, err := os.Stat(confs[2]); os.IsNotExist(err) {
		WriteDefaultWalletConfig(confs[2])
	} else {
		walletCfgData, err := ioutil.ReadFile(confs[2])
		if err != nil {
			WriteDefaultWalletConfig(confs[2])
		} else {
			err = json.Unmarshal(walletCfgData, &walletCfg)
			if err != nil {
				WriteDefaultWalletConfig(confs[2])
			}
		}
	}

	shellCfg := *DefaultShellConfig()
	if _, err := os.Stat(confs[3]); os.IsNotExist(err) {
		WriteDefaultShellConfig(confs[3])
	} else {
		shellCfgData, err := ioutil.ReadFile(confs[3])
		if err != nil {
			WriteDefaultShellConfig(confs[3])
		} else {
			err = json.Unmarshal(shellCfgData, &shellCfg)
			if err != nil {
				WriteDefaultShellConfig(confs[3])
			}
		}
	}

	// push all current settings as from the conf configuration to the module configs
	nodeCfg.Node.Listeners = ConfConfig.NodeListeners
	shellCfg.Node.Listeners = ConfConfig.NodeListeners
	nodeCfg.Node.RPCListeners = ConfConfig.NodeRPCListeners
	walletCfg.Wallet.RPCConnect = ConfConfig.NodeRPCListeners[0]
	shellCfg.Node.RPCListeners = ConfConfig.NodeRPCListeners
	shellCfg.Wallet.RPCConnect = ConfConfig.NodeRPCListeners[0]
	ctlCfg.RPCServer = ConfConfig.NodeRPCListeners[0]

	walletCfg.Wallet.LegacyRPCListeners = ConfConfig.WalletListeners
	ctlCfg.Wallet = ConfConfig.NodeRPCListeners[0]
	shellCfg.Wallet.LegacyRPCListeners = ConfConfig.NodeRPCListeners
	walletCfg.Wallet.LegacyRPCListeners = ConfConfig.WalletListeners
	ctlCfg.Wallet = ConfConfig.WalletListeners[0]
	shellCfg.Wallet.LegacyRPCListeners = ConfConfig.WalletListeners

	nodeCfg.Node.RPCUser = ConfConfig.NodeUser
	walletCfg.Wallet.PodUsername = ConfConfig.NodeUser
	walletCfg.Wallet.Username = ConfConfig.NodeUser
	shellCfg.Node.RPCUser = ConfConfig.NodeUser
	shellCfg.Wallet.PodUsername = ConfConfig.NodeUser
	shellCfg.Wallet.Username = ConfConfig.NodeUser
	ctlCfg.RPCUser = ConfConfig.NodeUser

	nodeCfg.Node.RPCPass = ConfConfig.NodePass
	walletCfg.Wallet.PodPassword = ConfConfig.NodePass
	walletCfg.Wallet.Password = ConfConfig.NodePass
	shellCfg.Node.RPCPass = ConfConfig.NodePass
	shellCfg.Wallet.PodPassword = ConfConfig.NodePass
	shellCfg.Wallet.Password = ConfConfig.NodePass
	ctlCfg.RPCPass = ConfConfig.NodePass

	nodeCfg.Node.RPCKey = ConfConfig.RPCKey
	walletCfg.Wallet.RPCKey = ConfConfig.RPCKey
	shellCfg.Node.RPCKey = ConfConfig.RPCKey
	shellCfg.Wallet.RPCKey = ConfConfig.RPCKey

	nodeCfg.Node.RPCCert = ConfConfig.RPCCert
	walletCfg.Wallet.RPCCert = ConfConfig.RPCCert
	shellCfg.Node.RPCCert = ConfConfig.RPCCert
	shellCfg.Wallet.RPCCert = ConfConfig.RPCCert

	walletCfg.Wallet.CAFile = ConfConfig.CAFile
	shellCfg.Wallet.CAFile = ConfConfig.CAFile

	nodeCfg.Node.TLS = ConfConfig.TLS
	walletCfg.Wallet.EnableClientTLS = ConfConfig.TLS
	shellCfg.Node.TLS = ConfConfig.TLS
	shellCfg.Wallet.EnableClientTLS = ConfConfig.TLS
	walletCfg.Wallet.EnableServerTLS = ConfConfig.TLS
	shellCfg.Wallet.EnableServerTLS = ConfConfig.TLS
	ctlCfg.TLSSkipVerify = ConfConfig.SkipVerify

	ctlCfg.Proxy = ConfConfig.Proxy
	nodeCfg.Node.Proxy = ConfConfig.Proxy
	walletCfg.Wallet.Proxy = ConfConfig.Proxy
	shellCfg.Node.Proxy = ConfConfig.Proxy
	shellCfg.Wallet.Proxy = ConfConfig.Proxy

	ctlCfg.ProxyUser = ConfConfig.ProxyUser
	nodeCfg.Node.ProxyUser = ConfConfig.ProxyUser
	walletCfg.Wallet.ProxyUser = ConfConfig.ProxyUser
	shellCfg.Node.ProxyUser = ConfConfig.ProxyUser
	shellCfg.Wallet.ProxyUser = ConfConfig.ProxyUser

	ctlCfg.ProxyPass = ConfConfig.ProxyPass
	nodeCfg.Node.ProxyPass = ConfConfig.ProxyPass
	walletCfg.Wallet.ProxyPass = ConfConfig.ProxyPass
	shellCfg.Node.ProxyPass = ConfConfig.ProxyPass
	shellCfg.Wallet.ProxyPass = ConfConfig.ProxyPass

	walletCfg.Wallet.WalletPass = ConfConfig.WalletPass
	shellCfg.Wallet.WalletPass = ConfConfig.WalletPass

	var r string
	var ok bool
	var listeners []string
	if r, ok = getIfIs(ctx, "nodelistener"); ok {
		NormalizeAddresses(r, node.DefaultPort, &listeners)
		ConfConfig.NodeListeners = listeners
		nodeCfg.Node.Listeners = listeners
		shellCfg.Node.Listeners = listeners
	}
	if r, ok = getIfIs(ctx, "noderpclistener"); ok {
		NormalizeAddresses(r, node.DefaultRPCPort, &listeners)
		ConfConfig.NodeRPCListeners = listeners
		nodeCfg.Node.RPCListeners = listeners
		walletCfg.Wallet.RPCConnect = r
		shellCfg.Node.RPCListeners = listeners
		shellCfg.Wallet.RPCConnect = r
		ctlCfg.RPCServer = r
	}
	if r, ok = getIfIs(ctx, "walletlistener"); ok {
		NormalizeAddresses(r, node.DefaultRPCPort, &listeners)
		ConfConfig.WalletListeners = listeners
		walletCfg.Wallet.LegacyRPCListeners = listeners
		ctlCfg.Wallet = r
		shellCfg.Wallet.LegacyRPCListeners = listeners
	}
	if r, ok = getIfIs(ctx, "user"); ok {
		ConfConfig.NodeUser = r
		nodeCfg.Node.RPCUser = r
		walletCfg.Wallet.PodUsername = r
		walletCfg.Wallet.Username = r
		shellCfg.Node.RPCUser = r
		shellCfg.Wallet.PodUsername = r
		shellCfg.Wallet.Username = r
		ctlCfg.RPCUser = r
	}
	if r, ok = getIfIs(ctx, "pass"); ok {
		ConfConfig.NodePass = r
		nodeCfg.Node.RPCPass = r
		walletCfg.Wallet.PodPassword = r
		walletCfg.Wallet.Password = r
		shellCfg.Node.RPCPass = r
		shellCfg.Wallet.PodPassword = r
		shellCfg.Wallet.Password = r
		ctlCfg.RPCPass = r
	}
	if r, ok = getIfIs(ctx, "walletpass"); ok {
		ConfConfig.WalletPass = r
		walletCfg.Wallet.WalletPass = ConfConfig.WalletPass
		shellCfg.Wallet.WalletPass = ConfConfig.WalletPass
	}

	if r, ok = getIfIs(ctx, "rpckey"); ok {
		r = node.CleanAndExpandPath(r)
		ConfConfig.RPCKey = r
		nodeCfg.Node.RPCKey = r
		walletCfg.Wallet.RPCKey = r
		shellCfg.Node.RPCKey = r
		shellCfg.Wallet.RPCKey = r
	}
	if r, ok = getIfIs(ctx, "rpccert"); ok {
		r = node.CleanAndExpandPath(r)
		ConfConfig.RPCCert = r
		nodeCfg.Node.RPCCert = r
		walletCfg.Wallet.RPCCert = r
		shellCfg.Node.RPCCert = r
		shellCfg.Wallet.RPCCert = r
	}
	if r, ok = getIfIs(ctx, "cafile"); ok {
		r = node.CleanAndExpandPath(r)
		ConfConfig.CAFile = r
		walletCfg.Wallet.CAFile = r
		shellCfg.Wallet.CAFile = r
	}
	if r, ok = getIfIs(ctx, "tls"); ok {
		ConfConfig.TLS = r == "true"
		nodeCfg.Node.TLS = ConfConfig.TLS
		walletCfg.Wallet.EnableClientTLS = ConfConfig.TLS
		shellCfg.Node.TLS = ConfConfig.TLS
		shellCfg.Wallet.EnableClientTLS = ConfConfig.TLS
		walletCfg.Wallet.EnableServerTLS = ConfConfig.TLS
		shellCfg.Wallet.EnableServerTLS = ConfConfig.TLS
	}
	if r, ok = getIfIs(ctx, "skipverify"); ok {
		ConfConfig.SkipVerify = r == "true"
		ctlCfg.TLSSkipVerify = r == "true"
	}
	if r, ok = getIfIs(ctx, "proxy"); ok {
		NormalizeAddresses(r, node.DefaultRPCPort, &listeners)
		ConfConfig.Proxy = r
		ctlCfg.Proxy = ConfConfig.Proxy
		nodeCfg.Node.Proxy = ConfConfig.Proxy
		walletCfg.Wallet.Proxy = ConfConfig.Proxy
		shellCfg.Node.Proxy = ConfConfig.Proxy
		shellCfg.Wallet.Proxy = ConfConfig.Proxy
	}
	if r, ok = getIfIs(ctx, "proxyuser"); ok {
		ConfConfig.ProxyUser = r
		ctlCfg.ProxyUser = ConfConfig.ProxyUser
		nodeCfg.Node.ProxyUser = ConfConfig.ProxyUser
		walletCfg.Wallet.ProxyUser = ConfConfig.ProxyUser
		shellCfg.Node.ProxyUser = ConfConfig.ProxyUser
		shellCfg.Wallet.ProxyUser = ConfConfig.ProxyUser
	}
	if r, ok = getIfIs(ctx, "proxypass"); ok {
		ConfConfig.ProxyPass = r
		ctlCfg.ProxyPass = ConfConfig.ProxyPass
		nodeCfg.Node.ProxyPass = ConfConfig.ProxyPass
		walletCfg.Wallet.ProxyPass = ConfConfig.ProxyPass
		shellCfg.Node.ProxyPass = ConfConfig.ProxyPass
		shellCfg.Wallet.ProxyPass = ConfConfig.ProxyPass
	}
	if r, ok = getIfIs(ctx, "network"); ok {
		r = strings.ToLower(r)
		switch r {
		case "mainnet", "testnet", "regtestnet", "simnet":
		default:
			r = "mainnet"
		}
		ConfConfig.Network = r
		switch r {
		case "mainnet":
			ctlCfg.TestNet3 = false
			ctlCfg.SimNet = false
			nodeCfg.Node.TestNet3 = false
			nodeCfg.Node.SimNet = false
			nodeCfg.Node.RegressionTest = false
			walletCfg.Wallet.SimNet = false
			walletCfg.Wallet.TestNet3 = false
			shellCfg.Node.TestNet3 = false
			shellCfg.Node.RegressionTest = false
			shellCfg.Node.SimNet = false
			shellCfg.Wallet.TestNet3 = false
			shellCfg.Wallet.SimNet = false
		case "testnet":
			ctlCfg.TestNet3 = true
			ctlCfg.SimNet = false
			nodeCfg.Node.TestNet3 = true
			nodeCfg.Node.SimNet = false
			nodeCfg.Node.RegressionTest = false
			walletCfg.Wallet.SimNet = false
			walletCfg.Wallet.TestNet3 = true
			shellCfg.Node.TestNet3 = true
			shellCfg.Node.RegressionTest = false
			shellCfg.Node.SimNet = false
			shellCfg.Wallet.TestNet3 = true
			shellCfg.Wallet.SimNet = false
		case "regtestnet":
			ctlCfg.TestNet3 = false
			ctlCfg.SimNet = false
			nodeCfg.Node.TestNet3 = false
			nodeCfg.Node.SimNet = false
			nodeCfg.Node.RegressionTest = true
			walletCfg.Wallet.SimNet = false
			walletCfg.Wallet.TestNet3 = false
			shellCfg.Node.TestNet3 = false
			shellCfg.Node.RegressionTest = true
			shellCfg.Node.SimNet = false
			shellCfg.Wallet.TestNet3 = false
			shellCfg.Wallet.SimNet = false
		case "simnet":
			ctlCfg.TestNet3 = false
			ctlCfg.SimNet = true
			nodeCfg.Node.TestNet3 = false
			nodeCfg.Node.SimNet = true
			nodeCfg.Node.RegressionTest = false
			walletCfg.Wallet.SimNet = true
			walletCfg.Wallet.TestNet3 = false
			shellCfg.Node.TestNet3 = false
			shellCfg.Node.RegressionTest = false
			shellCfg.Node.SimNet = true
			shellCfg.Wallet.TestNet3 = false
			shellCfg.Wallet.SimNet = true
		}
	}
	WriteConfConfig(cfgFile, ConfConfig)
	// Now write the configs for all the others reading them and overwriting the changed values
	WriteCtlConfig(confs[0], &ctlCfg)
	WriteNodeConfig(confs[1], &nodeCfg)
	WriteWalletConfig(confs[2], &walletCfg)
	WriteShellConfig(confs[3], &shellCfg)
	if ctx.Is("show") {
		j, err := json.MarshalIndent(ConfConfig, "", "  ")
		if err != nil {
			panic(err.Error())
		}
		fmt.Println(string(j))
	}
}

// WriteConfConfig creates and writes the config file in the requested location
func WriteConfConfig(cfgFile string, cfg ConfCfg) {
	j, err := json.MarshalIndent(ConfConfig, "", "  ")
	if err != nil {
		panic(err.Error())
	}
	j = append(j, '\n')
	err = ioutil.WriteFile(cfgFile, j, 0600)
	if err != nil {
		panic(err.Error())
	}
}

// WriteDefaultConfConfig creates and writes a default config file in the requested location
func WriteDefaultConfConfig(cfgFile string) {
	defCfg := DefaultConfConfig()
	j, err := json.MarshalIndent(defCfg, "", "  ")
	if err != nil {
		panic(err.Error())
	}
	j = append(j, '\n')
	err = ioutil.WriteFile(cfgFile, j, 0600)
	if err != nil {
		panic(err.Error())
	}
	// if we are writing default config we also want to use it
	ConfConfig = *defCfg
}

// DefaultConfConfig returns a crispy fresh default conf configuration
func DefaultConfConfig() *ConfCfg {
	u := GenKey()
	p := GenKey()
	return &ConfCfg{
		NodeListeners:    []string{"127.0.0.1:11047"},
		NodeRPCListeners: []string{"127.0.0.1:11048"},
		WalletListeners:  []string{"127.0.0.1:11046"},
		NodeUser:         u,
		NodePass:         p,
		WalletPass:       wallet.InsecurePubPassphrase,
		RPCKey:           walletmain.DefaultRPCKeyFile,
		RPCCert:          walletmain.DefaultRPCCertFile,
		CAFile:           walletmain.DefaultCAFile,
		TLS:              false,
		SkipVerify:       false,
		Proxy:            "",
		ProxyUser:        "",
		ProxyPass:        "",
		Network:          "mainnet",
	}
}
