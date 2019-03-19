package app

import (
	"io/ioutil"
	"path/filepath"
	"time"

	"git.parallelcoin.io/dev/pod/cmd/ctl"
	"git.parallelcoin.io/dev/pod/cmd/node"
	walletmain "git.parallelcoin.io/dev/pod/cmd/wallet"
	netparams "git.parallelcoin.io/dev/pod/pkg/chain/config/params"
	"git.parallelcoin.io/dev/pod/pkg/util"
	"git.parallelcoin.io/pod/pkg/altsrc"
	"gopkg.in/urfave/cli.v1"
)

const appName = "pod"
const confExt = ".yaml"
const podConfigFilename = appName + confExt
const ctlAppName = "ctl"
const ctlConfigFilename = ctlAppName + confExt
const nodeAppName = "node"
const nodeConfigFilename = nodeAppName + confExt
const walletAppName = "wallet"
const walletConfigFilename = walletAppName + confExt

// App is the heart of the application system, this creates and initialises it.
var App = cli.NewApp()

// DefaultHomeDir is the default location where the data directory is located
// when none is specified.
var DefaultHomeDir = util.AppDataDir(appName, false)

var activeNetParams *netparams.Params
var appDatadir = util.AppDataDir(appName, false)

type ConfigCommon struct {
	Datadir      string
	Save         bool
	Loglevel     string
	Subsystem    cli.StringSlice
	Network      string
	ServerUser   string
	ServerPass   string
	ClientUser   string
	ClientPass   string
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

// These are in a form usable for getting a *bool
var True, False = true, false

var appConfigCommon = &ConfigCommon{

	Subsystem: cli.StringSlice{},
}

var ctlConfig = ctl.Config{

	ShowVersion:   new(bool),
	ListCommands:  new(bool),
	ConfigFile:    new(string),
	DebugLevel:    &appConfigCommon.Loglevel,
	RPCUser:       &appConfigCommon.ServerUser,
	RPCPass:       &appConfigCommon.ServerPass,
	RPCServer:     &(*nodeConfig.RPCListeners)[0],
	TestNet3:      new(bool),
	SimNet:        new(bool),
	TLSSkipVerify: new(bool),
	Wallet:        &(*walletConfig.LegacyRPCListeners)[0],
	RPCCert:       &appConfigCommon.RPCcert,
	TLS:           &appConfigCommon.ClientTLS,
	Proxy:         &appConfigCommon.Proxy,
	ProxyUser:     &appConfigCommon.Proxyuser,
	ProxyPass:     &appConfigCommon.Proxypass,
}

var ctlDatadir = "ctl"

var defaultDatadir = "~/.pod"

var guiDataDir = "/gui"

var nodeConfig = &node.Config{

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
	DebugLevel:           &appConfigCommon.Loglevel,
	AddPeers:             new(cli.StringSlice),
	ConnectPeers:         new(cli.StringSlice),
	DisableListen:        new(bool),
	MaxPeers:             new(int),
	DisableBanning:       new(bool),
	BanDuration:          new(time.Duration),
	BanThreshold:         new(int),
	Whitelists:           new(cli.StringSlice),
	RPCUser:              &appConfigCommon.ServerUser,
	RPCPass:              &appConfigCommon.ServerPass,
	RPCLimitUser:         &appConfigCommon.ClientUser,
	RPCLimitPass:         &appConfigCommon.ClientPass,
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

	CAFile:                   &appConfigCommon.CAfile,
	EnableClientTLS:          &appConfigCommon.ClientTLS,
	Proxy:                    &appConfigCommon.Proxy,
	ProxyUser:                &appConfigCommon.Proxyuser,
	ProxyPass:                &appConfigCommon.Proxypass,
	UseSPV:                   new(bool),
	RPCCert:                  &appConfigCommon.RPCcert,
	RPCKey:                   &appConfigCommon.RPCkey,
	EnableServerTLS:          &appConfigCommon.ServerTLS,
	ConfigFile:               new(string),
	ShowVersion:              new(bool),
	LogLevel:                 &appConfigCommon.Loglevel,
	Create:                   new(bool),
	CreateTemp:               new(bool),
	AppDataDir:               new(string),
	TestNet3:                 new(bool),
	SimNet:                   new(bool),
	NoInitialLoad:            new(bool),
	LogDir:                   new(string),
	Profile:                  new(string),
	WalletPass:               new(string),
	RPCConnect:               &(*nodeConfig.RPCListeners)[0],
	PodUsername:              &appConfigCommon.ServerUser,
	PodPassword:              &appConfigCommon.ServerPass,
	AddPeers:                 new(cli.StringSlice),
	ConnectPeers:             new(cli.StringSlice),
	MaxPeers:                 new(int),
	BanDuration:              new(time.Duration),
	BanThreshold:             new(int),
	OneTimeTLSKey:            new(bool),
	LegacyRPCListeners:       &cli.StringSlice{"127.0.0.1:11048"},
	LegacyRPCMaxClients:      new(int),
	LegacyRPCMaxWebsockets:   new(int),
	Username:                 new(string),
	Password:                 new(string),
	ExperimentalRPCListeners: new(cli.StringSlice),
	DataDir:                  new(string),
}

var walletDataDir = "/wallet"

// NewSourceFromFlagAndBase creates a new Yaml
// InputSourceContext from a provided flag name and source context.
// If file doesn't exist, make one, empty is same as whatever is default
func NewSourceFromFlagAndBase(c *cli.Context, confName, flagFileName string,

) func(context *cli.Context) (altsrc.InputSourceContext, error) {

	return func(context *cli.Context) (altsrc.InputSourceContext, error) {

		filePath := c.String(flagFileName)
		filePath = filepath.Join(filePath, confName)
		EnsureDir(filePath)

		if !FileExists(filePath) {

			err := ioutil.WriteFile(filePath, []byte{'\n'}, 0600)

			if err != nil {

				panic(err)
			}

		}

		return altsrc.NewYamlSourceFromFile(filePath)
	}

}
