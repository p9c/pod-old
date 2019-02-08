package app

import (
	"fmt"

	"git.parallelcoin.io/pod/cmd/ctl"
	"git.parallelcoin.io/pod/cmd/shell"
	"github.com/davecgh/go-spew/spew"
)

func getConfs(datadir string) {
	confs = []string{
		datadir + "/ctl/conf.json",
		datadir + "/node/conf.json",
		datadir + "/wallet/conf.json",
		datadir + "/shell/conf.json",
	}
}

// ConfigSet is a full set of configuration structs
type ConfigSet struct {
	Conf   *ConfCfg
	Ctl    *ctl.Config
	Node   *NodeCfg
	Wallet *WalletCfg
	Shell  *shell.Config
}

// WriteConfigSet writes a set of configurations to disk
func WriteConfigSet(in *ConfigSet) {
	WriteConfConfig(in.Conf)
	WriteCtlConfig(in.Ctl)
	WriteNodeConfig(in.Node)
	WriteWalletConfig(in.Wallet)
	WriteShellConfig(in.Shell)
	return
}

// GetDefaultConfs returns all of the configurations in their default state
func GetDefaultConfs(datadir string) (out *ConfigSet) {
	out = new(ConfigSet)
	out.Conf = DefaultConfConfig(datadir)
	out.Ctl = DefaultCtlConfig(datadir)
	out.Node = DefaultNodeConfig(datadir)
	out.Wallet = DefaultWalletConfig(datadir)
	out.Shell = DefaultShellConfig(datadir)
	return
}

// SyncToConfs takes a ConfigSet and synchronises the values according to the ConfCfg settings
func SyncToConfs(in *ConfigSet) {
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

// PortSet is a single set of ports for a configuration
type PortSet struct {
	P2P       string
	NodeRPC   string
	WalletRPC string
}

// GenPortSet creates a set of ports for testnet configuration
func GenPortSet(portbase int) (ps *PortSet) {
	// From the base, each element is as follows:
	// - P2P = portbase
	// - NodeRPC = portbase + 1
	// - WalletRPC =  portbase -1
	// For each set, the base is incremented by 100
	// so from 21047, you get 21047, 21048, 21046
	// and next would be 21147, 21148, 21146
	t := portbase
	ps = &PortSet{
		P2P:       fmt.Sprint(t),
		NodeRPC:   fmt.Sprint(t + 1),
		WalletRPC: fmt.Sprint(t - 1),
	}
	return
}
