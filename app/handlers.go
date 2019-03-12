package app

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"git.parallelcoin.io/pod/cmd/node"
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
	case "simnet", "s":
		*ctlConfig.TestNet3 = false
		*ctlConfig.SimNet = true
	default:
		if network != "mainnet" && network != "m" {
			fmt.Println("using mainnet for ctl")
		}
		*ctlConfig.TestNet3 = false
		*ctlConfig.SimNet = false
	}
	if appConfigCommon.Save {
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
		activeNetParams = &node.TestNet3Params
	case "regtestnet", "regressiontest", "r":
		*nodeConfig.TestNet3 = false
		*nodeConfig.SimNet = false
		*nodeConfig.RegressionTest = true
		activeNetParams = &node.SimNetParams
	case "simnet", "s":
		*nodeConfig.TestNet3 = false
		*nodeConfig.SimNet = true
		*nodeConfig.RegressionTest = false
		activeNetParams = &node.MainNetParams
	default:
		if network != "mainnet" && network != "m" {
			fmt.Println("using mainnet for node")
		}
		*nodeConfig.TestNet3 = false
		*nodeConfig.SimNet = false
		*nodeConfig.RegressionTest = false
	}
	if !*nodeConfig.Onion {
		*nodeConfig.OnionProxy = ""
	}
	if appConfigCommon.Save {
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
		activeNetParams = &node.TestNet3Params
	case "simnet", "s":
		*walletConfig.TestNet3 = false
		*walletConfig.SimNet = true
		activeNetParams = &node.SimNetParams
	default:
		if network != "mainnet" && network != "m" {
			fmt.Println("using mainnet for wallet")
		}
		*walletConfig.TestNet3 = false
		*walletConfig.SimNet = false
		activeNetParams = &node.MainNetParams
	}
	if appConfigCommon.Save {
		podHandleSave()
		walletHandleSave()
		return nil
	}
	return launchWallet(c)
}

func confHandle(c *cli.Context) error {
	appConfigCommon.Save = true
	_ = ctlHandle(c)
	appConfigCommon.Save = true
	_ = nodeHandle(c)
	appConfigCommon.Save = true
	_ = walletHandle(c)
	appConfigCommon.Save = true
	return nil
}
