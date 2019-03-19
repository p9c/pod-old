package app

import (
	"io/ioutil"
	"path/filepath"
	"time"

	"gopkg.in/urfave/cli.v1/altsrc"
	netparams "git.parallelcoin.io/dev/pod/pkg/chain/config/params"
	"git.parallelcoin.io/dev/pod/pkg/pod"
	"git.parallelcoin.io/dev/pod/pkg/util"
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

var DefaultDataDir = util.AppDataDir(appName, false)

var activeNetParams *netparams.Params
var appDatadir = util.AppDataDir(appName, false)
var podConfig = podDefConfig()

func podDefConfig() pod.Config {

	return pod.Config{
		ShowVersion:              new(bool),
		ConfigFile:               new(string),
		DataDir:                  new(string),
		LogDir:                   new(string),
		LogLevel:                 new(string),
		Subsystems:               new(cli.StringSlice),
		AddPeers:                 new(cli.StringSlice),
		ConnectPeers:             new(cli.StringSlice),
		MaxPeers:                 new(int),
		Listeners:                new(cli.StringSlice),
		DisableListen:            new(bool),
		DisableBanning:           new(bool),
		BanDuration:              new(time.Duration),
		BanThreshold:             new(int),
		Whitelists:               new(cli.StringSlice),
		Username:                 new(string),
		Password:                 new(string),
		ServerUser:               new(string),
		ServerPass:               new(string),
		LimitUser:                new(string),
		LimitPass:                new(string),
		RPCConnect:               new(string),
		RPCListeners:             new(cli.StringSlice),
		RPCCert:                  new(string),
		RPCKey:                   new(string),
		RPCMaxClients:            new(int),
		RPCMaxWebsockets:         new(int),
		RPCMaxConcurrentReqs:     new(int),
		RPCQuirks:                new(bool),
		DisableRPC:               new(bool),
		TLS:                      new(bool),
		DisableDNSSeed:           new(bool),
		ExternalIPs:              new(cli.StringSlice),
		Proxy:                    new(string),
		ProxyUser:                new(string),
		ProxyPass:                new(string),
		OnionProxy:               new(string),
		OnionProxyUser:           new(string),
		OnionProxyPass:           new(string),
		Onion:                    new(bool),
		TorIsolation:             new(bool),
		TestNet3:                 new(bool),
		RegressionTest:           new(bool),
		SimNet:                   new(bool),
		AddCheckpoints:           new(cli.StringSlice),
		DisableCheckpoints:       new(bool),
		DbType:                   new(string),
		Profile:                  new(string),
		CPUProfile:               new(string),
		Upnp:                     new(bool),
		MinRelayTxFee:            new(float64),
		FreeTxRelayLimit:         new(float64),
		NoRelayPriority:          new(bool),
		TrickleInterval:          new(time.Duration),
		MaxOrphanTxs:             new(int),
		Algo:                     new(string),
		Generate:                 new(bool),
		GenThreads:               new(int),
		MiningAddrs:              new(cli.StringSlice),
		MinerListener:            new(string),
		MinerPass:                new(string),
		BlockMinSize:             new(int),
		BlockMaxSize:             new(int),
		BlockMinWeight:           new(int),
		BlockMaxWeight:           new(int),
		BlockPrioritySize:        new(int),
		UserAgentComments:        new(cli.StringSlice),
		NoPeerBloomFilters:       new(bool),
		NoCFilters:               new(bool),
		SigCacheMaxSize:          new(int),
		BlocksOnly:               new(bool),
		TxIndex:                  new(bool),
		AddrIndex:                new(bool),
		RelayNonStd:              new(bool),
		RejectNonStd:             new(bool),
		ListCommands:             new(bool),
		TLSSkipVerify:            new(bool),
		Wallet:                   new(string),
		NoInitialLoad:            new(bool),
		WalletPass:               new(string),
		CAFile:                   new(string),
		OneTimeTLSKey:            new(bool),
		ServerTLS:                new(bool),
		LegacyRPCListeners:       new(cli.StringSlice),
		LegacyRPCMaxClients:      new(int),
		LegacyRPCMaxWebsockets:   new(int),
		ExperimentalRPCListeners: new(cli.StringSlice),
	}
}

// NewSourceFromFlagAndBase creates a new Yaml
// InputSourceContext from a provided flag name and source context.
// If file doesn't exist, make one, empty is same as whatever is default
func NewSourceFromFlagAndBase(
	c *cli.Context, confName, flagFileName string,
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
