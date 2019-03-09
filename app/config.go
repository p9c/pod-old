package app

import (
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
var nodeDataDir = "/node"
var shellDataDir = "/shell"
var walletDataDir = "/wallet"
