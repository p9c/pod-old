package app

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	netparams "git.parallelcoin.io/dev/pod/pkg/chain/config/params"
	"gopkg.in/urfave/cli.v1"
	"gopkg.in/yaml.v1"
)

func walletHandleSave() {

	appConfigCommon.Save = false
	*walletConfig.LogDir = *walletConfig.DataDir
	yn, e := yaml.Marshal(walletConfig)

	if e == nil {

		EnsureDir(*walletConfig.ConfigFile)

		if e := ioutil.WriteFile(*walletConfig.ConfigFile, yn, 0600); e != nil {

			panic(e)
		}

	} else {
		panic(e)
	}

}

func walletHandle(c *cli.Context) error {

	*walletConfig.DataDir = appConfigCommon.Datadir
	*walletConfig.AppDataDir = filepath.Join(
		appConfigCommon.Datadir,
		walletAppName,
	)
	*walletConfig.ConfigFile = filepath.Join(
		*walletConfig.AppDataDir,
		walletConfigFilename,
	)

	if !c.Parent().Bool("useproxy") {

		*nodeConfig.Proxy = ""
	}

	loglevel := c.Parent().String("loglevel")

	switch loglevel {

	case "trace", "debug", "info", "warn", "error", "fatal":
	default:
		*walletConfig.LogLevel = "info"
	}

	network := c.Parent().String("network")

	switch network {

	case "testnet", "testnet3", "t":
		*walletConfig.TestNet3 = true
		*walletConfig.SimNet = false
		activeNetParams = &netparams.TestNet3Params
	case "simnet", "s":
		*walletConfig.TestNet3 = false
		*walletConfig.SimNet = true
		activeNetParams = &netparams.SimNetParams
	default:

		if network != "mainnet" && network != "m" {

			fmt.Println("using mainnet for wallet")
		}

		*walletConfig.TestNet3 = false
		*walletConfig.SimNet = false
		activeNetParams = &netparams.MainNetParams
	}

	_ = podHandle(c)

	if appConfigCommon.Save {

		appConfigCommon.Save = false
		podHandleSave()
		walletHandleSave()
		return nil
	}

	return launchWallet(c)
}
