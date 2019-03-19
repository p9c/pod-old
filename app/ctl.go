package app

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	netparams "git.parallelcoin.io/dev/pod/pkg/chain/config/params"
	"gopkg.in/urfave/cli.v1"
	"gopkg.in/yaml.v1"
)

func ctlHandleSave() {

	appConfigCommon.Save = false
	yn, e := yaml.Marshal(ctlConfig)

	if e == nil {

		EnsureDir(*ctlConfig.ConfigFile)
		e =
			ioutil.WriteFile(*ctlConfig.ConfigFile, yn, 0600)

		if e != nil {

			panic(e)
		}

	} else {

		panic(e)
	}

}
func ctlHandle(c *cli.Context) error {

	datadir :=
		filepath.Join(
			appConfigCommon.Datadir,
			ctlAppName,
		)
	*ctlConfig.ConfigFile = filepath.Join(
		datadir,
		ctlConfigFilename,
	)

	if !c.Parent().Bool("useproxy") {

		*ctlConfig.Proxy = ""
	}

	loglevel := c.Parent().String("loglevel")

	switch loglevel {

	case "trace", "debug", "info", "warn", "error", "fatal":

	default:
		*ctlConfig.DebugLevel = "warn"
	}

	network := c.Parent().String("network")

	switch network {

	case "testnet", "testnet3", "t":
		*ctlConfig.TestNet3 = true
		*ctlConfig.SimNet = false
		activeNetParams = &netparams.TestNet3Params

	case "simnet", "s":
		*ctlConfig.TestNet3 = false
		*ctlConfig.SimNet = true
		activeNetParams = &netparams.SimNetParams

	default:

		if network != "mainnet" && network != "m" {

			fmt.Println("using mainnet for ctl")
		}

		*ctlConfig.TestNet3 = false
		*ctlConfig.SimNet = false
		activeNetParams = &netparams.MainNetParams
	}

	_ = podHandle(c)

	if appConfigCommon.Save {

		appConfigCommon.Save = false
		podHandleSave()
		ctlHandleSave()
		return nil
	}

	return launchCtl(c)
}
