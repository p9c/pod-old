package app

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"

	netparams "git.parallelcoin.io/pod/pkg/chain/config/params"
	"github.com/davecgh/go-spew/spew"
	"gopkg.in/urfave/cli.v1"
	"gopkg.in/yaml.v1"
)

func ctlHandleSave() {
	appConfigCommon.Save = false
	yn, e := yaml.Marshal(ctlConfig)
	if e == nil {
		EnsureDir(*ctlConfig.ConfigFile)
		e = ioutil.WriteFile(
			*ctlConfig.ConfigFile, yn, 0600)
		if e != nil {
			panic(e)
		}
	} else {
		panic(e)
	}
}

func nodeHandleSave() {
	appConfigCommon.Save = false
	*nodeConfig.LogDir = *nodeConfig.DataDir
	yn, e := yaml.Marshal(nodeConfig)
	if e == nil {
		EnsureDir(*nodeConfig.ConfigFile)
		e = ioutil.WriteFile(
			*nodeConfig.ConfigFile, yn, 0600)
		if e != nil {
			panic(e)
		}
	} else {
		panic(e)
	}
}

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

func podHandleSave() {
	podconfig := filepath.Join(appConfigCommon.Datadir, podConfigFilename)
	fmt.Println("saving to", podconfig)
	if yp, e := yaml.Marshal(appConfigCommon); e == nil {
		EnsureDir(podconfig)
		if e := ioutil.WriteFile(podconfig, yp, 0600); e != nil {
			panic(e)
		}
	} else {
		panic(e)
	}
}

func ctlHandle(c *cli.Context) error {
	datadir := filepath.Join(
		appConfigCommon.Datadir,
		ctlAppName)
	*ctlConfig.ConfigFile = filepath.Join(
		datadir,
		ctlConfigFilename)
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
	if appConfigCommon.Save {
		appConfigCommon.Save = false
		podHandleSave()
		ctlHandleSave()
		return nil
	}
	return launchCtl(c)
}

func ctlHandleList(c *cli.Context) error {
	fmt.Println("running ctl listcommands")
	_ = ctlHandle(c)
	spew.Dump(ctlConfig)
	return nil
}

func nodeHandle(c *cli.Context) error {
	*nodeConfig.DataDir = filepath.Join(
		appConfigCommon.Datadir,
		nodeAppName)
	*nodeConfig.ConfigFile = filepath.Join(
		*nodeConfig.DataDir,
		nodeConfigFilename)
	*nodeConfig.LogDir = *nodeConfig.DataDir
	if !c.Parent().Bool("useproxy") {
		*nodeConfig.Proxy = ""
	}
	loglevel := c.Parent().String("loglevel")
	switch loglevel {
	case "trace", "debug", "info", "warn", "error", "fatal":
	default:
		*nodeConfig.DebugLevel = "info"
	}
	network := c.Parent().String("network")
	switch network {
	case "testnet", "testnet3", "t":
		*nodeConfig.TestNet3 = true
		*nodeConfig.SimNet = false
		*nodeConfig.RegressionTest = false
		activeNetParams = &netparams.TestNet3Params
	case "regtestnet", "regressiontest", "r":
		*nodeConfig.TestNet3 = false
		*nodeConfig.SimNet = false
		*nodeConfig.RegressionTest = true
		activeNetParams = &netparams.RegressionTestParams
	case "simnet", "s":
		*nodeConfig.TestNet3 = false
		*nodeConfig.SimNet = true
		*nodeConfig.RegressionTest = false
		activeNetParams = &netparams.SimNetParams
	default:
		if network != "mainnet" && network != "m" {
			fmt.Println("using mainnet for node")
		}
		*nodeConfig.TestNet3 = false
		*nodeConfig.SimNet = false
		*nodeConfig.RegressionTest = false
		activeNetParams = &netparams.MainNetParams

	}
	if !*nodeConfig.Onion {
		*nodeConfig.OnionProxy = ""
	}
	if appConfigCommon.Save {
		appConfigCommon.Save = false
		podHandleSave()
		nodeHandleSave()
		return nil
	}
	return launchNode(c)
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
	if appConfigCommon.Save {
		appConfigCommon.Save = false
		podHandleSave()
		walletHandleSave()
		return nil
	}
	return launchWallet(c)
}

func confHandle(c *cli.Context) error {
	appendNum := false
	number := 1
	if c.IsSet("number") {
		appendNum = true
		number = c.Int("number")
		if number > 10 {
			return errors.New("cannot make more than 10 (0-9) test profiles")
		}
	}
	base := c.String("base")
	var working string
	fmt.Println("base:", base)
	for i := 0; i < number; i++ {
		working = "" + base
		if appendNum {
			working += fmt.Sprint(i)
		}
		apps := []string{"c", "n", "w", "s", "g"}
		for _, x := range apps {
			e := App.Run([]string{"pod", "-i", "-D", working, x})
			if e != nil {
				panic(e)
			}
		}
	}
	return nil
}
