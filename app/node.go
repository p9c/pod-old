package app

import (
	"io/ioutil"
	"path/filepath"
	// "github.com/davecgh/go-spew/spew"
	"gopkg.in/urfave/cli.v1"
	"gopkg.in/yaml.v1"
)

func nodeHandle(c *cli.Context) error {
	// fmt.Println("running node")
	if !c.Parent().IsSet("useproxy") {
		*nodeConfig.Proxy = ""
	}
	if !*nodeConfig.Onion {
		*nodeConfig.OnionProxy = ""
	}
	if appConfigCommon.Save {
		appConfigCommon.Save = false
		*nodeConfig.DataDir = filepath.Join(
			appConfigCommon.Datadir,
			nodePath)
		*nodeConfig.ConfigFile = filepath.Join(
			*nodeConfig.DataDir,
			nodeConfigFilename)
		*nodeConfig.LogDir = *nodeConfig.DataDir
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
		yn, e := yaml.Marshal(nodeConfig)
		if e == nil {
			// fmt.Println(*nodeConfig.ConfigFile, string(yn))
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

	// spew.Dump(nodeConfig)
	// spew.Dump(c.Args(), c.FlagNames())
	return nil
}
