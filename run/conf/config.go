package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	c "git.parallelcoin.io/pod/module/ctl"
	n "git.parallelcoin.io/pod/module/node"
	w "git.parallelcoin.io/pod/module/wallet"
	"git.parallelcoin.io/pod/run/ctl"
	"git.parallelcoin.io/pod/run/def"
	"git.parallelcoin.io/pod/run/node"
	"git.parallelcoin.io/pod/run/shell"
	"git.parallelcoin.io/pod/run/util"
	"git.parallelcoin.io/pod/run/wallet"
	"github.com/tucnak/climax"
)

var confFile = def.DefaultDataDir + "/conf"

// Configuration is the settings that can be set to synchronise across all pod modules
type Configuration struct {
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

// Config is the configuration for this tool
var Config Configuration

// Apps is the central repository of all the other app configurations
var Apps AppConfigs

// AppConfigs are the configurations for each app that are applied
type AppConfigs struct {
	Ctl    c.Config
	Node   n.Config
	Wallet w.Config
	Shell  s.Cfg
}

var f = pu.GenFlag
var t = pu.GenTrig
var s = pu.GenShort
var l = pu.GenLog

// Command is a command to send RPC queries to bitcoin RPC protocol server for node and wallet queries
var Command = climax.Command{
	Name:  "conf",
	Brief: "sets configurations common across modules",
	Help:  "automates synchronising common settings between servers and clients",
	Flags: []climax.Flag{
		t("version", "V", "show version number and quit"),
		t("init", "i", "resets configuration to defaults"),
		t("show", "s", "prints currently configuration"),

		f("nodelistener", "main peer to peer address for apps that connect to the parallelcoin peer to peer network"),
		f("noderpclistener", "address where node listens for RPC"),
		f("walletlistener", "address where wallet listens for RPC"),
		s("user", "u", "username for all the things"),
		s("pass", "P", "password for all the things"),
		s("walletpass", "w", "public password for wallet"),
		f("rpckey", "RPC server certificate key"),
		f("rpccert", "RPC server certificate"),
		f("cafile", "RPC server certificate chain for validation"),
		f("tls", "enable/disable TLS"),
		f("skipverify", "do not verify TLS certificates (not recommended!)"),
		f("proxy", "connect via SOCKS5 proxy"),
		f("proxyuser", "username for proxy"),
		f("proxypass", "password for proxy"),

		f("network", "connect to [mainnet|testnet|regtestnet|simnet]"),
	},
	Examples: []climax.Example{
		{
			Usecase:     "--nodeuser=user --nodepass=pa55word",
			Description: "set the username and password for the node RPC",
		},
	},
	Handle: func(ctx climax.Context) int {
		if ctx.Is("version") {
			fmt.Println("pod conf version", Version())
			os.Exit(0)
		}
		if ctx.Is("init") {
			WriteDefaultConfig(confFile)
		} else {
			if _, err := os.Stat(confFile); os.IsNotExist(err) {
				WriteDefaultConfig(confFile)
			} else {
				cfgData, err := ioutil.ReadFile(confFile)
				if err != nil {
					WriteDefaultConfig(confFile)
				}
				err = json.Unmarshal(cfgData, &Config)
				if err != nil {
					WriteDefaultConfig(confFile)
				}
			}
		}
		configConf(&ctx, confFile)
		runCtl()
		return 0
	},
}

func getIfIs(ctx *climax.Context, name string) (out string, ok bool) {
	if ctx.Is(name) {
		return ctx.Get(name)
	}
	return
}

func configConf(ctx *climax.Context, cfgFile string) {
	// First load all of the module configurations and unmarshal into their structs
	confs := []string{
		def.DefaultDataDir + "/ctl/conf",
		def.DefaultDataDir + "/node/conf",
		def.DefaultDataDir + "/wallet/conf",
		def.DefaultDataDir + "/shell/conf",
	}

	// If we can't parse the config files we just reset them to default

	ctlCfg := *ctl.DefaultConfig()
	if _, err := os.Stat(confs[0]); os.IsNotExist(err) {
		ctl.WriteDefaultConfig(confs[0])
	} else {
		ctlCfgData, err := ioutil.ReadFile(confs[0])
		if err != nil {
			shell.WriteDefaultConfig(confs[0])
		} else {
			err = json.Unmarshal(ctlCfgData, &ctlCfg)
			if err != nil {
				shell.WriteDefaultConfig(confs[0])
			}
		}
	}

	nodeCfg := *node.DefaultConfig()
	if _, err := os.Stat(confs[1]); os.IsNotExist(err) {
		node.WriteDefaultConfig(confs[1])
	} else {
		nodeCfgData, err := ioutil.ReadFile(confs[1])
		if err != nil {
			node.WriteDefaultConfig(confs[1])
		} else {
			err = json.Unmarshal(nodeCfgData, &nodeCfg)
			if err != nil {
				node.WriteDefaultConfig(confs[1])
			}
		}
	}

	walletCfg := *walletrun.DefaultConfig()
	if _, err := os.Stat(confs[2]); os.IsNotExist(err) {
		walletrun.WriteDefaultConfig(confs[2])
	} else {
		walletCfgData, err := ioutil.ReadFile(confs[2])
		if err != nil {
			walletrun.WriteDefaultConfig(confs[2])
		} else {
			err = json.Unmarshal(walletCfgData, &walletCfg)
			if err != nil {
				walletrun.WriteDefaultConfig(confs[2])
			}
		}
	}

	shellCfg := *shell.DefaultConfig()
	if _, err := os.Stat(confs[3]); os.IsNotExist(err) {
		shell.WriteDefaultConfig(confs[3])
	} else {
		shellCfgData, err := ioutil.ReadFile(confs[3])
		if err != nil {
			shell.WriteDefaultConfig(confs[3])
		} else {
			err = json.Unmarshal(shellCfgData, &shellCfg)
			if err != nil {
				shell.WriteDefaultConfig(confs[3])
			}
		}
	}

	// push all current settings as from the conf configuration to the module configs
	nodeCfg.Node.Listeners = Config.NodeListeners
	shellCfg.Node.Listeners = Config.NodeListeners
	nodeCfg.Node.RPCListeners = Config.NodeRPCListeners
	walletCfg.Wallet.RPCConnect = Config.NodeRPCListeners[0]
	shellCfg.Node.RPCListeners = Config.NodeRPCListeners
	shellCfg.Wallet.RPCConnect = Config.NodeRPCListeners[0]
	ctlCfg.RPCServer = Config.NodeRPCListeners[0]

	walletCfg.Wallet.LegacyRPCListeners = Config.WalletListeners
	ctlCfg.Wallet = Config.NodeRPCListeners[0]
	shellCfg.Wallet.LegacyRPCListeners = Config.NodeRPCListeners
	walletCfg.Wallet.LegacyRPCListeners = Config.WalletListeners
	ctlCfg.Wallet = Config.WalletListeners[0]
	shellCfg.Wallet.LegacyRPCListeners = Config.WalletListeners

	nodeCfg.Node.RPCUser = Config.NodeUser
	walletCfg.Wallet.PodUsername = Config.NodeUser
	walletCfg.Wallet.Username = Config.NodeUser
	shellCfg.Node.RPCUser = Config.NodeUser
	shellCfg.Wallet.PodUsername = Config.NodeUser
	shellCfg.Wallet.Username = Config.NodeUser
	ctlCfg.RPCUser = Config.NodeUser

	nodeCfg.Node.RPCPass = Config.NodePass
	walletCfg.Wallet.PodPassword = Config.NodePass
	walletCfg.Wallet.Password = Config.NodePass
	shellCfg.Node.RPCPass = Config.NodePass
	shellCfg.Wallet.PodPassword = Config.NodePass
	shellCfg.Wallet.Password = Config.NodePass
	ctlCfg.RPCPass = Config.NodePass

	nodeCfg.Node.RPCKey = Config.RPCKey
	walletCfg.Wallet.RPCKey = Config.RPCKey
	shellCfg.Node.RPCKey = Config.RPCKey
	shellCfg.Wallet.RPCKey = Config.RPCKey

	nodeCfg.Node.RPCCert = Config.RPCCert
	walletCfg.Wallet.RPCCert = Config.RPCCert
	shellCfg.Node.RPCCert = Config.RPCCert
	shellCfg.Wallet.RPCCert = Config.RPCCert

	walletCfg.Wallet.CAFile = Config.CAFile
	shellCfg.Wallet.CAFile = Config.CAFile

	nodeCfg.Node.TLS = Config.TLS
	walletCfg.Wallet.EnableClientTLS = Config.TLS
	shellCfg.Node.TLS = Config.TLS
	shellCfg.Wallet.EnableClientTLS = Config.TLS
	walletCfg.Wallet.EnableServerTLS = Config.TLS
	shellCfg.Wallet.EnableServerTLS = Config.TLS
	ctlCfg.TLSSkipVerify = Config.SkipVerify

	ctlCfg.Proxy = Config.Proxy
	nodeCfg.Node.Proxy = Config.Proxy
	walletCfg.Wallet.Proxy = Config.Proxy
	shellCfg.Node.Proxy = Config.Proxy
	shellCfg.Wallet.Proxy = Config.Proxy

	ctlCfg.ProxyUser = Config.ProxyUser
	nodeCfg.Node.ProxyUser = Config.ProxyUser
	walletCfg.Wallet.ProxyUser = Config.ProxyUser
	shellCfg.Node.ProxyUser = Config.ProxyUser
	shellCfg.Wallet.ProxyUser = Config.ProxyUser

	ctlCfg.ProxyPass = Config.ProxyPass
	nodeCfg.Node.ProxyPass = Config.ProxyPass
	walletCfg.Wallet.ProxyPass = Config.ProxyPass
	shellCfg.Node.ProxyPass = Config.ProxyPass
	shellCfg.Wallet.ProxyPass = Config.ProxyPass

	walletCfg.Wallet.WalletPass = Config.WalletPass
	shellCfg.Wallet.WalletPass = Config.WalletPass

	var r string
	var ok bool
	var listeners []string
	if r, ok = getIfIs(ctx, "nodelistener"); ok {
		pu.NormalizeAddresses(r, n.DefaultPort, &listeners)
		Config.NodeListeners = listeners
		nodeCfg.Node.Listeners = listeners
		shellCfg.Node.Listeners = listeners
	}
	if r, ok = getIfIs(ctx, "noderpclistener"); ok {
		pu.NormalizeAddresses(r, n.DefaultRPCPort, &listeners)
		Config.NodeRPCListeners = listeners
		nodeCfg.Node.RPCListeners = listeners
		walletCfg.Wallet.RPCConnect = r
		shellCfg.Node.RPCListeners = listeners
		shellCfg.Wallet.RPCConnect = r
		ctlCfg.RPCServer = r
	}
	if r, ok = getIfIs(ctx, "walletlistener"); ok {
		pu.NormalizeAddresses(r, n.DefaultRPCPort, &listeners)
		Config.WalletListeners = listeners
		walletCfg.Wallet.LegacyRPCListeners = listeners
		ctlCfg.Wallet = r
		shellCfg.Wallet.LegacyRPCListeners = listeners
	}
	if r, ok = getIfIs(ctx, "user"); ok {
		Config.NodeUser = r
		nodeCfg.Node.RPCUser = r
		walletCfg.Wallet.PodUsername = r
		walletCfg.Wallet.Username = r
		shellCfg.Node.RPCUser = r
		shellCfg.Wallet.PodUsername = r
		shellCfg.Wallet.Username = r
		ctlCfg.RPCUser = r
	}
	if r, ok = getIfIs(ctx, "pass"); ok {
		Config.NodePass = r
		nodeCfg.Node.RPCPass = r
		walletCfg.Wallet.PodPassword = r
		walletCfg.Wallet.Password = r
		shellCfg.Node.RPCPass = r
		shellCfg.Wallet.PodPassword = r
		shellCfg.Wallet.Password = r
		ctlCfg.RPCPass = r
	}
	if r, ok = getIfIs(ctx, "walletpass"); ok {
		Config.WalletPass = r
		walletCfg.Wallet.WalletPass = Config.WalletPass
		shellCfg.Wallet.WalletPass = Config.WalletPass
	}

	if r, ok = getIfIs(ctx, "rpckey"); ok {
		r = n.CleanAndExpandPath(r)
		Config.RPCKey = r
		nodeCfg.Node.RPCKey = r
		walletCfg.Wallet.RPCKey = r
		shellCfg.Node.RPCKey = r
		shellCfg.Wallet.RPCKey = r
	}
	if r, ok = getIfIs(ctx, "rpccert"); ok {
		r = n.CleanAndExpandPath(r)
		Config.RPCCert = r
		nodeCfg.Node.RPCCert = r
		walletCfg.Wallet.RPCCert = r
		shellCfg.Node.RPCCert = r
		shellCfg.Wallet.RPCCert = r
	}
	if r, ok = getIfIs(ctx, "cafile"); ok {
		r = n.CleanAndExpandPath(r)
		Config.CAFile = r
		walletCfg.Wallet.CAFile = r
		shellCfg.Wallet.CAFile = r
	}
	if r, ok = getIfIs(ctx, "tls"); ok {
		Config.TLS = r == "true"
		nodeCfg.Node.TLS = Config.TLS
		walletCfg.Wallet.EnableClientTLS = Config.TLS
		shellCfg.Node.TLS = Config.TLS
		shellCfg.Wallet.EnableClientTLS = Config.TLS
		walletCfg.Wallet.EnableServerTLS = Config.TLS
		shellCfg.Wallet.EnableServerTLS = Config.TLS
	}
	if r, ok = getIfIs(ctx, "skipverify"); ok {
		Config.SkipVerify = r == "true"
		ctlCfg.TLSSkipVerify = r == "true"
	}
	if r, ok = getIfIs(ctx, "proxy"); ok {
		pu.NormalizeAddresses(r, n.DefaultRPCPort, &listeners)
		Config.Proxy = r
		ctlCfg.Proxy = Config.Proxy
		nodeCfg.Node.Proxy = Config.Proxy
		walletCfg.Wallet.Proxy = Config.Proxy
		shellCfg.Node.Proxy = Config.Proxy
		shellCfg.Wallet.Proxy = Config.Proxy
	}
	if r, ok = getIfIs(ctx, "proxyuser"); ok {
		Config.ProxyUser = r
		ctlCfg.ProxyUser = Config.ProxyUser
		nodeCfg.Node.ProxyUser = Config.ProxyUser
		walletCfg.Wallet.ProxyUser = Config.ProxyUser
		shellCfg.Node.ProxyUser = Config.ProxyUser
		shellCfg.Wallet.ProxyUser = Config.ProxyUser
	}
	if r, ok = getIfIs(ctx, "proxypass"); ok {
		Config.ProxyPass = r
		ctlCfg.ProxyPass = Config.ProxyPass
		nodeCfg.Node.ProxyPass = Config.ProxyPass
		walletCfg.Wallet.ProxyPass = Config.ProxyPass
		shellCfg.Node.ProxyPass = Config.ProxyPass
		shellCfg.Wallet.ProxyPass = Config.ProxyPass
	}
	if r, ok = getIfIs(ctx, "network"); ok {
		r = strings.ToLower(r)
		switch r {
		case "mainnet", "testnet", "regtestnet", "simnet":
		default:
			r = "mainnet"
		}
		Config.Network = r
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
	WriteConfig(cfgFile, Config)
	// Now write the configs for all the others reading them and overwriting the changed values
	ctl.WriteConfig(confs[0], &ctlCfg)
	node.WriteConfig(confs[1], &nodeCfg)
	walletrun.WriteConfig(confs[2], &walletCfg)
	shell.WriteConfig(confs[3], &shellCfg)
	if ctx.Is("show") {
		j, err := json.MarshalIndent(Config, "", "  ")
		if err != nil {
			panic(err.Error())
		}
		fmt.Println(string(j))
	}
}

// WriteConfig creates and writes the config file in the requested location
func WriteConfig(cfgFile string, cfg Configuration) {
	j, err := json.MarshalIndent(Config, "", "  ")
	if err != nil {
		panic(err.Error())
	}
	j = append(j, '\n')
	err = ioutil.WriteFile(cfgFile, j, 0600)
	if err != nil {
		panic(err.Error())
	}
}

// WriteDefaultConfig creates and writes a default config file in the requested location
func WriteDefaultConfig(cfgFile string) {
	defCfg := defaultConfig()
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
	Config = *defCfg
}

func defaultConfig() *Configuration {
	u := pu.GenKey()
	p := pu.GenKey()
	k := pu.GenKey()
	return &Configuration{
		NodeListeners:    []string{"127.0.0.1:11047"},
		NodeRPCListeners: []string{"127.0.0.1:11048"},
		WalletListeners:  []string{"127.0.0.1:11046"},
		NodeUser:         u,
		NodePass:         p,
		WalletPass:       k,
		RPCKey:           w.DefaultRPCKeyFile,
		RPCCert:          w.DefaultRPCCertFile,
		CAFile:           w.DefaultCAFile,
		TLS:              false,
		SkipVerify:       false,
		Proxy:            "",
		ProxyUser:        "",
		ProxyPass:        "",
		Network:          "mainnet",
	}
}
