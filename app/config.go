package app

import (
	"git.parallelcoin.io/pod/cmd/ctl"
	"git.parallelcoin.io/pod/cmd/node"
	"git.parallelcoin.io/pod/cmd/wallet"
	"github.com/urfave/cli"
	"time"
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
	Onion        bool
	OnionProxy   string
	Onionuser    string
	Onionpass    string
	Torisolation bool
}

var True, False = true, false

var appConfigCommon = &ConfigCommon{
	Subsystems: make(cli.StringSlice, 0),
}

var ctlConfig = ctl.Config{

	ShowVersion:   new(bool),
	ListCommands:  new(bool),
	ConfigFile:    new(string),
	DebugLevel:    new(string),
	RPCUser:       new(string),
	RPCPass:       new(string),
	RPCServer:     new(string),
	TestNet3:      new(bool),
	SimNet:        new(bool),
	TLSSkipVerify: new(bool),
	Wallet:        new(string),

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

	RPCCert:              &appConfigCommon.RPCcert,
	RPCKey:               &appConfigCommon.RPCkey,
	TLS:                  &appConfigCommon.ServerTLS,
	Proxy:                &appConfigCommon.Proxy,
	ProxyUser:            &appConfigCommon.Proxyuser,
	ProxyPass:            &appConfigCommon.Proxypass,
	OnionProxy:           &appConfigCommon.OnionProxy,
	OnionProxyUser:       &appConfigCommon.Onionuser,
	OnionProxyPass:       &appConfigCommon.Onionpass,
	Onion:                &appConfigCommon.Onion,
	TorIsolation:         &appConfigCommon.Torisolation,
	Listeners:            &cli.StringSlice{node.DefaultListener},
	RPCListeners:         &cli.StringSlice{node.DefaultRPCListener},
	ShowVersion:          new(bool),
	ConfigFile:           new(string),
	DataDir:              new(string),
	LogDir:               new(string),
	DebugLevel:           new(string),
	AddPeers:             new(cli.StringSlice),
	ConnectPeers:         new(cli.StringSlice),
	DisableListen:        new(bool),
	MaxPeers:             new(int),
	DisableBanning:       new(bool),
	BanDuration:          new(time.Duration),
	BanThreshold:         new(int),
	Whitelists:           new(cli.StringSlice),
	RPCUser:              new(string),
	RPCPass:              new(string),
	RPCLimitUser:         new(string),
	RPCLimitPass:         new(string),
	RPCMaxClients:        new(int),
	RPCMaxWebsockets:     new(int),
	RPCMaxConcurrentReqs: new(int),
	RPCQuirks:            new(bool),
	DisableRPC:           new(bool),
	DisableDNSSeed:       new(bool),
	ExternalIPs:          new(cli.StringSlice),
	TestNet3:             new(bool),
	RegressionTest:       new(bool),
	SimNet:               new(bool),
	AddCheckpoints:       new(cli.StringSlice),
	DisableCheckpoints:   new(bool),
	DbType:               new(string),
	Profile:              new(string),
	CPUProfile:           new(string),
	Upnp:                 new(bool),
	MinRelayTxFee:        new(float64),
	FreeTxRelayLimit:     new(float64),
	NoRelayPriority:      new(bool),
	TrickleInterval:      new(time.Duration),
	MaxOrphanTxs:         new(int),
	Algo:                 new(string),
	Generate:             new(bool),
	GenThreads:           new(int),
	MiningAddrs:          new(cli.StringSlice),
	MinerListener:        new(string),
	MinerPass:            new(string),
	BlockMinSize:         new(int),
	BlockMaxSize:         new(int),
	BlockMinWeight:       new(int),
	BlockMaxWeight:       new(int),
	BlockPrioritySize:    new(int),
	UserAgentComments:    new(cli.StringSlice),
	NoPeerBloomFilters:   new(bool),
	NoCFilters:           new(bool),
	DropCfIndex:          new(bool),
	SigCacheMaxSize:      new(int),
	BlocksOnly:           new(bool),
	TxIndex:              new(bool),
	DropTxIndex:          new(bool),
	AddrIndex:            new(bool),
	DropAddrIndex:        new(bool),
	RelayNonStd:          new(bool),
	RejectNonStd:         new(bool),
}

var nodeDataDir = "/node"

var shellDataDir = "/shell"

var walletConfig = walletmain.Config{
	CAFile:          &appConfigCommon.CAfile,
	EnableClientTLS: &appConfigCommon.ClientTLS,
	Proxy:           &appConfigCommon.Proxy,
	ProxyUser:       &appConfigCommon.Proxyuser,
	ProxyPass:       &appConfigCommon.Proxypass,
	UseSPV:          new(bool),
	RPCCert:         &appConfigCommon.RPCcert,
	RPCKey:          &appConfigCommon.RPCkey,
	EnableServerTLS: &appConfigCommon.ServerTLS,
}

var walletDataDir = "/wallet"
