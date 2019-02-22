package shell

import (
	"git.parallelcoin.io/pod/cmd/node"
	walletmain "git.parallelcoin.io/pod/cmd/wallet"
	"git.parallelcoin.io/pod/pkg/netparams"
)

// Config is the combined app and logging configuration data
type Config struct {
	ConfigFile      string
	DataDir         string
	AppDataDir      string
	Node            *node.Config
	Wallet          *walletmain.Config
	Levels          map[string]string
	nodeActiveNet   *node.Params
	walletActiveNet *netparams.Params
}

// GetNodeActiveNet returns the activenet params
func (r *Config) GetNodeActiveNet() *node.Params {
	return r.nodeActiveNet
}

// GetWalletActiveNet returns the activenet params
func (r *Config) GetWalletActiveNet() *netparams.Params {
	return r.walletActiveNet
}

// SetNodeActiveNet returns the activenet params
func (r *Config) SetNodeActiveNet(in *node.Params) {

	r.nodeActiveNet = in
}

// SetWalletActiveNet returns the activenet params
func (r *Config) SetWalletActiveNet(in *netparams.Params) {

	r.walletActiveNet = in
}
