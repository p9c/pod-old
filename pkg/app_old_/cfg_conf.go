package app_old

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"git.parallelcoin.io/dev/pod/cmd/ctl"
	"git.parallelcoin.io/dev/pod/cmd/node"
	"git.parallelcoin.io/dev/pod/cmd/shell"
	walletmain "git.parallelcoin.io/dev/pod/cmd/wallet"
	"git.parallelcoin.io/dev/pod/pkg/chain/fork"
	"github.com/davecgh/go-spew/spew"
	"github.com/tucnak/climax"
)

// ConfigSet is a full set of configuration structs

type ConfigSet struct {
	Conf   *ConfCfg
	Ctl    *ctl.Config
	Node   *NodeCfg
	Wallet *WalletCfg
	Shell  *shell.Config
}

// PortSet is a single set of ports for a configuration

type PortSet struct {
	P2P       string
	NodeRPC   string
	WalletRPC string
}

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

// GenPortSet creates a set of ports for testnet configuration
func GenPortSet(
	portbase int,
) (
	ps *PortSet,

) {

	/* From the base, each element is as follows:

	- P2P = portbase

	- NodeRPC = portbase + 1

	- WalletRPC =  portbase -1

	For each set, the base is incremented by 100
	so from 21047, you get 21047, 21048, 21046
	and next would be 21147, 21148, 21146
	*/
	t := portbase

	ps = &PortSet{

		P2P:       fmt.Sprint(t),
		NodeRPC:   fmt.Sprint(t + 1),
		WalletRPC: fmt.Sprint(t - 1),
	}

	return
}

// GetDefaultConfs returns all of the configurations in their default state
func GetDefaultConfs(
	datadir string,
) (
	out *ConfigSet,

) {

	out = new(ConfigSet)
	out.Conf = DefaultConfConfig(datadir)
	out.Ctl = DefaultCtlConfig(datadir)
	out.Node = DefaultNodeConfig(datadir)
	out.Wallet = DefaultWalletConfig(datadir)
	out.Shell = DefaultShellConfig(datadir)
	return
}

// SyncToConfs takes a ConfigSet and synchronises the values according to the ConfCfg settings
func SyncToConfs(
	in *ConfigSet,

) {

	if in == nil {

		panic("received nil configset")
	}

	if in.Conf == nil ||
		in.Ctl == nil ||
		in.Node == nil ||
		in.Wallet == nil ||

		in.Shell == nil {

		panic("configset had a nil element\n" + spew.Sdump(in))
	}

	// push all current settings as from the conf configuration to the module configs
	in.Node.Node.Listeners = in.Conf.NodeListeners
	in.Shell.Node.Listeners = in.Conf.NodeListeners
	in.Node.Node.RPCListeners = in.Conf.NodeRPCListeners
	in.Wallet.Wallet.RPCConnect = in.Conf.NodeRPCListeners[0]
	in.Shell.Node.RPCListeners = in.Conf.NodeRPCListeners
	in.Shell.Wallet.RPCConnect = in.Conf.NodeRPCListeners[0]
	in.Ctl.RPCServer = in.Conf.NodeRPCListeners[0]

	in.Wallet.Wallet.LegacyRPCListeners = in.Conf.WalletListeners
	in.Ctl.Wallet = in.Conf.NodeRPCListeners[0]
	in.Shell.Wallet.LegacyRPCListeners = in.Conf.NodeRPCListeners
	in.Wallet.Wallet.LegacyRPCListeners = in.Conf.WalletListeners
	in.Ctl.Wallet = in.Conf.WalletListeners[0]
	in.Shell.Wallet.LegacyRPCListeners = in.Conf.WalletListeners

	in.Node.Node.RPCUser = in.Conf.NodeUser
	in.Wallet.Wallet.PodUsername = in.Conf.NodeUser
	in.Wallet.Wallet.Username = in.Conf.NodeUser
	in.Shell.Node.RPCUser = in.Conf.NodeUser
	in.Shell.Wallet.PodUsername = in.Conf.NodeUser
	in.Shell.Wallet.Username = in.Conf.NodeUser
	in.Ctl.RPCUser = in.Conf.NodeUser

	in.Node.Node.RPCPass = in.Conf.NodePass
	in.Wallet.Wallet.PodPassword = in.Conf.NodePass
	in.Wallet.Wallet.Password = in.Conf.NodePass
	in.Shell.Node.RPCPass = in.Conf.NodePass
	in.Shell.Wallet.PodPassword = in.Conf.NodePass
	in.Shell.Wallet.Password = in.Conf.NodePass
	in.Ctl.RPCPass = in.Conf.NodePass

	in.Node.Node.RPCKey = in.Conf.RPCKey
	in.Wallet.Wallet.RPCKey = in.Conf.RPCKey
	in.Shell.Node.RPCKey = in.Conf.RPCKey
	in.Shell.Wallet.RPCKey = in.Conf.RPCKey

	in.Node.Node.RPCCert = in.Conf.RPCCert
	in.Wallet.Wallet.RPCCert = in.Conf.RPCCert
	in.Shell.Node.RPCCert = in.Conf.RPCCert
	in.Shell.Wallet.RPCCert = in.Conf.RPCCert

	in.Wallet.Wallet.CAFile = in.Conf.CAFile
	in.Shell.Wallet.CAFile = in.Conf.CAFile

	in.Node.Node.TLS = in.Conf.TLS
	in.Wallet.Wallet.EnableClientTLS = in.Conf.TLS
	in.Shell.Node.TLS = in.Conf.TLS
	in.Shell.Wallet.EnableClientTLS = in.Conf.TLS
	in.Wallet.Wallet.EnableServerTLS = in.Conf.TLS
	in.Shell.Wallet.EnableServerTLS = in.Conf.TLS
	in.Ctl.TLSSkipVerify = in.Conf.SkipVerify

	in.Ctl.Proxy = in.Conf.Proxy
	in.Node.Node.Proxy = in.Conf.Proxy
	in.Wallet.Wallet.Proxy = in.Conf.Proxy
	in.Shell.Node.Proxy = in.Conf.Proxy
	in.Shell.Wallet.Proxy = in.Conf.Proxy

	in.Ctl.ProxyUser = in.Conf.ProxyUser
	in.Node.Node.ProxyUser = in.Conf.ProxyUser
	in.Wallet.Wallet.ProxyUser = in.Conf.ProxyUser
	in.Shell.Node.ProxyUser = in.Conf.ProxyUser
	in.Shell.Wallet.ProxyUser = in.Conf.ProxyUser

	in.Ctl.ProxyPass = in.Conf.ProxyPass
	in.Node.Node.ProxyPass = in.Conf.ProxyPass
	in.Wallet.Wallet.ProxyPass = in.Conf.ProxyPass
	in.Shell.Node.ProxyPass = in.Conf.ProxyPass
	in.Shell.Wallet.ProxyPass = in.Conf.ProxyPass

	in.Wallet.Wallet.WalletPass = in.Conf.WalletPass
	in.Shell.Wallet.WalletPass = in.Conf.WalletPass
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

// WriteConfigSet writes a set of configurations to disk
func WriteConfigSet(
	in *ConfigSet,

) {

	WriteConfConfig(in.Conf)
	WriteCtlConfig(in.Ctl)
	WriteNodeConfig(in.Node)
	WriteWalletConfig(in.Wallet)
	WriteShellConfig(in.Shell)
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

/*
func getConfs(
	datadir string,


) {



	confs = []string{

		datadir + "/ctl/conf.json",
		datadir + "/node/conf.json",
		datadir + "/wallet/conf.json",
		datadir + "/shell/conf.json",
	}

}

*/
