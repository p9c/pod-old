package app

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"git.parallelcoin.io/dev/pod/cmd/node"
	netparams "git.parallelcoin.io/dev/pod/pkg/chain/config/params"
	"git.parallelcoin.io/dev/pod/pkg/chain/fork"
	cl "git.parallelcoin.io/dev/pod/pkg/util/cl"
	"github.com/pelletier/go-toml"
	"gopkg.in/urfave/cli.v1"
)

func Configure() {

	log <- cl.Debug{"checking configurations"}

	if *podConfig.ConfigFile == "" {
		*podConfig.ConfigFile = filepath.Join(*podConfig.DataDir, podConfigFilename)
	}

	if *podConfig.LogDir == "" {

		*podConfig.LogDir = *podConfig.DataDir
	}

	if len(*podConfig.Listeners) < 1 {

		*podConfig.Listeners = append(*podConfig.Listeners, "127.0.0.1:11047")
	}

	if len(*podConfig.LegacyRPCListeners) < 1 {
		*podConfig.LegacyRPCListeners = append(*podConfig.LegacyRPCListeners, "127.0.0.1:11046")
	}

	if *podConfig.RPCCert == "" {

		*podConfig.RPCCert = filepath.Join(*podConfig.DataDir, "rppodConfig.cert")
	}

	if *podConfig.RPCKey == "" {

		*podConfig.RPCKey = filepath.Join(*podConfig.DataDir, "rppodConfig.key")
	}

	loglevel := *podConfig.LogLevel

	switch loglevel {

	case "trace", "debug", "info", "warn", "error", "fatal":
		log <- cl.Info{"log level", loglevel}
	default:
		log <- cl.Info{"unrecognised loglevel", loglevel, "setting default info"}
		*podConfig.LogLevel = "info"
	}

	cl.Register.SetAllLevels(*podConfig.LogLevel)

	if !*podConfig.Onion {

		*podConfig.OnionProxy = ""
	}

	network := "mainnet"
	if podConfig.Network != nil {
		network = *podConfig.Network
	}
	switch network {

	case "testnet", "testnet3", "t":
		log <- cl.Debug{"on testnet"}
		*podConfig.TestNet3 = true
		*podConfig.SimNet = false
		*podConfig.RegressionTest = false
		activeNetParams = &netparams.TestNet3Params
		fork.IsTestnet = true

	case "regtestnet", "regressiontest", "r":
		log <- cl.Debug{"on regression testnet"}
		*podConfig.TestNet3 = false
		*podConfig.SimNet = false
		*podConfig.RegressionTest = true
		activeNetParams = &netparams.RegressionTestParams

	case "simnet", "s":
		log <- cl.Debug{"on simnet"}
		*podConfig.TestNet3 = false
		*podConfig.SimNet = true
		*podConfig.RegressionTest = false
		activeNetParams = &netparams.SimNetParams

	default:

		if network != "mainnet" && network != "m" {

			log <- cl.Warn{"using mainnet for node"}
		}

		log <- cl.Debug{"on mainnet"}
		*podConfig.TestNet3 = false
		*podConfig.SimNet = false
		*podConfig.RegressionTest = false
		activeNetParams = &netparams.MainNetParams
	}

	log <- cl.Debug{"normalising addresses"}
	port := node.DefaultPort
	NormalizeStringSliceAddresses(podConfig.AddPeers, port)
	NormalizeStringSliceAddresses(podConfig.ConnectPeers, port)
	NormalizeStringSliceAddresses(podConfig.Listeners, port)
	NormalizeStringSliceAddresses(podConfig.Whitelists, port)
	NormalizeStringSliceAddresses(podConfig.RPCListeners, port)

}

func podHandleSave() {

	StateCfg.Save = false
	*podConfig.ConfigFile =
		filepath.Join(
			node.CleanAndExpandPath(*podConfig.DataDir),
			podConfigFilename,
		)

	if yp, e := toml.Marshal(podConfig); e == nil {

		EnsureDir(*podConfig.ConfigFile)

		if e := ioutil.WriteFile(*podConfig.ConfigFile, yp, 0600); e != nil {

			panic(e)
		}

	} else {

		panic(e)
	}

}

func podHandle(c *cli.Context) error {

	*podConfig.RPCCert = node.CleanAndExpandPath(*podConfig.RPCCert)
	*podConfig.RPCKey = node.CleanAndExpandPath(*podConfig.RPCKey)
	*podConfig.CAFile = node.CleanAndExpandPath(*podConfig.CAFile)

	NormalizeAddress(
		*podConfig.Proxy, "9050", podConfig.Proxy)
	NormalizeAddress(
		*podConfig.OnionProxy, "9050", podConfig.OnionProxy)
	return nil
}

func confHandle(c *cli.Context) error {

	appendNum := false
	number := 1

	if c.IsSet("number") {

		appendNum = true
		number = c.Int("number")

		if number > 100 {

			return errors.New("cannot make more than 100 (0-99) test profiles")
		}

	}

	base := c.String("base")
	var working string
	fmt.Println("base:", base)

	for i := 0; i < number; i++ {

		working = "" + base

		if appendNum {

			working += fmt.Sprintf("%02d", i)
		}

		// if e != nil {

		// 	panic(e)
		// }

	}

	return nil
}
