package app

import (
	"io/ioutil"
	"path/filepath"
	// "github.com/davecgh/go-spew/spew"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"gopkg.in/urfave/cli.v1"
	"gopkg.in/yaml.v1"
)

func podHandleSave() {
	podconfig := filepath.Join(appConfigCommon.Datadir, podConfigFilename)
	yp, e := yaml.Marshal(appConfigCommon)
	if e == nil {
		EnsureDir(podconfig)
		ioutil.WriteFile(
			podconfig,
			yp, 0600)
	} else {
		panic(e)
	}
}

func nodeHandleSave() {
	appConfigCommon.Save = false
	*nodeConfig.LogDir = *nodeConfig.DataDir
	podHandleSave()
	yn, e := yaml.Marshal(nodeConfig)
	if e == nil {
		EnsureDir(*nodeConfig.ConfigFile)
		e = ioutil.WriteFile(
			*nodeConfig.ConfigFile,
			yn,
			0600)
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
	if !c.Parent().IsSet("useproxy") {
		*nodeConfig.Proxy = ""
	}
	if !*nodeConfig.Onion {
		*nodeConfig.OnionProxy = ""
	}
	if appConfigCommon.Save {
		nodeHandleSave()
	}
	return nil
}

func ctlHandle(c *cli.Context) error {
	datadir := filepath.Join(
		appConfigCommon.Datadir,
		ctlAppName)
	*ctlConfig.ConfigFile = filepath.Join(
		datadir,
		ctlConfigFilename)
	if !c.Parent().IsSet("useproxy") {
		*ctlConfig.Proxy = ""
	}
	if appConfigCommon.Save {
		appConfigCommon.Save = false
		podHandleSave()
		yn, e := yaml.Marshal(ctlConfig)
		if e == nil {
			EnsureDir(*ctlConfig.ConfigFile)
			e = ioutil.WriteFile(
				*ctlConfig.ConfigFile,
				yn,
				0600)
			if e != nil {
				panic(e)
			}
		} else {
			panic(e)
		}
	}
	return nil
}

func ctlHandleList(c *cli.Context) error {
	fmt.Println("running ctl listcommands")
	*ctlConfig.ListCommands = true
	if !c.IsSet("wallet") {
		*ctlConfig.Wallet = ""
	}
	if !c.Parent().IsSet("useproxy") {
		*ctlConfig.Proxy = ""
	}
	spew.Dump(ctlConfig)
	return nil
}
