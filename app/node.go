package app

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"git.parallelcoin.io/pod/cmd/node"
	netparams "git.parallelcoin.io/pod/pkg/chain/config/params"
	"gopkg.in/urfave/cli.v1"
	"gopkg.in/yaml.v1"
)

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
	// TODO: now to sanitize the rest
	port := node.DefaultPort
	NormalizeStringSliceAddresses(nodeConfig.AddPeers, port)
	NormalizeStringSliceAddresses(nodeConfig.ConnectPeers, port)
	NormalizeStringSliceAddresses(nodeConfig.Listeners, port)
	NormalizeStringSliceAddresses(nodeConfig.Whitelists, port)
	NormalizeStringSliceAddresses(nodeConfig.RPCListeners, port)

	_ = podHandle(c)
	if appConfigCommon.Save {
		appConfigCommon.Save = false
		podHandleSave()
		nodeHandleSave()
		return nil
	}
	return launchNode(c)
}

func NormalizeStringSliceAddresses(a *cli.StringSlice, port string) {
	variable := []string(*a)
	NormalizeAddresses(
		strings.Join(variable, " "),
		port, &variable)
	*a = variable
}
