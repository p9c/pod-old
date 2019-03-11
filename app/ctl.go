package app

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"gopkg.in/urfave/cli.v1"
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"path/filepath"
)

func ctlHandle(c *cli.Context) error {
	if !c.Parent().IsSet("useproxy") {
		*ctlConfig.Proxy = ""
	}
	if appConfigCommon.Save {
		appConfigCommon.Save = false
		datadir := filepath.Join(
			appConfigCommon.Datadir,
			ctlAppName)
		*ctlConfig.ConfigFile = filepath.Join(
			datadir,
			ctlConfigFilename)
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
		yn, e := yaml.Marshal(ctlConfig)
		if e == nil {
			// fmt.Println(*nodeConfig.ConfigFile, string(yn))
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
