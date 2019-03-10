package app

import (
	"git.parallelcoin.io/pod/cmd/node"
	"github.com/urfave/cli"
)

type ConfigCommon struct {
	Datadir      string
	Save         bool
	Loglevel     string
	Subsystems   cli.StringSlice
	Network      string
	RPCcert      string
	RPCkey       string
	CAfile       string
	ClientTLS    bool
	ServerTLS    bool
	Useproxy     bool
	Proxy        string
	Proxyuser    string
	Proxypass    string
	Useonion     bool
	Onion        string
	Onionuser    string
	Onionpass    string
	Torisolation bool
}

var appConfigCommon = ConfigCommon{}
var ctlDatadir = "ctl"
var defaultDatadir = "~/.pod"
var guiDataDir = "/gui"

var nodeConfig = node.Config{
	Listeners:    cli.StringSlice{node.DefaultListener},
	RPCListeners: cli.StringSlice{node.DefaultRPCListener},
}
var nodeDataDir = "/node"
var shellDataDir = "/shell"
var walletDataDir = "/wallet"
