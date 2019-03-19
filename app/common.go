package app

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"git.parallelcoin.io/dev/pod/cmd/node"
	"gopkg.in/urfave/cli.v1"
	"gopkg.in/yaml.v1"
)

func podHandleSave() {

	podCfg :=
		filepath.Join(
			node.CleanAndExpandPath(*podConfig.DataDir),
			podConfigFilename,
		)

	if yp, e := yaml.Marshal(podCfg); e == nil {

		EnsureDir(podCfg)

		if e := ioutil.WriteFile(podCfg, yp, 0600); e != nil {

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
