package app

import (
	"git.parallelcoin.io/pod/cmd/ctl"
	"git.parallelcoin.io/pod/cmd/node"
	"git.parallelcoin.io/pod/cmd/wallet"
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
	NoOnion      bool
	Onion        string
	Onionuser    string
	Onionpass    string
	Torisolation bool
}

var True, False = true, false

var appConfigCommon = &ConfigCommon{}
var ctlConfig = ctl.Config{
	RPCCert:   &appConfigCommon.RPCcert,
	TLS:       &appConfigCommon.ClientTLS,
	Proxy:     &appConfigCommon.Proxy,
	ProxyUser: &appConfigCommon.Proxyuser,
	ProxyPass: &appConfigCommon.Proxypass,
}
var ctlDatadir = "ctl"
var defaultDatadir = "~/.pod"
var guiDataDir = "/gui"
var nodeConfig = node.Config{
	RPCCert:        &appConfigCommon.RPCcert,
	RPCKey:         &appConfigCommon.RPCkey,
	TLS:            &appConfigCommon.ServerTLS,
	Proxy:          &appConfigCommon.Proxy,
	ProxyUser:      &appConfigCommon.Proxyuser,
	ProxyPass:      &appConfigCommon.Proxypass,
	OnionProxy:     &appConfigCommon.Onion,
	OnionProxyUser: &appConfigCommon.Onionuser,
	OnionProxyPass: &appConfigCommon.Onionpass,
	NoOnion:        &appConfigCommon.NoOnion,
	TorIsolation:   &appConfigCommon.Torisolation,
	Listeners:      &cli.StringSlice{node.DefaultListener},
	RPCListeners:   &cli.StringSlice{node.DefaultRPCListener},
}
var nodeDataDir = "/node"
var shellDataDir = "/shell"
var walletConfig = walletmain.Config{
	CAFile:          &appConfigCommon.CAfile,
	EnableClientTLS: &appConfigCommon.ClientTLS,
	Proxy:           &appConfigCommon.Proxy,
	ProxyUser:       &appConfigCommon.Proxyuser,
	ProxyPass:       &appConfigCommon.Proxypass,
	UseSPV:          &False,
	RPCCert:         &appConfigCommon.RPCcert,
	RPCKey:          &appConfigCommon.RPCkey,
	EnableServerTLS: &appConfigCommon.ServerTLS,
}
var walletDataDir = "/wallet"
