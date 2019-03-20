package app

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"git.parallelcoin.io/dev/pod/cmd/node"
	"git.parallelcoin.io/dev/pod/pkg/pod"
	cl "git.parallelcoin.io/dev/pod/pkg/util/cl"
	"github.com/pelletier/go-toml"
	"gopkg.in/urfave/cli.v1"
)

func Configure(c *pod.Config) {

	log <- cl.Debug{"checking configurations"}

	if *c.ConfigFile == "" {
		*c.ConfigFile = filepath.Join(*c.DataDir, podConfigFilename)
	}

	if *c.LogDir == "" {

		*c.LogDir = *c.DataDir
	}

	if len(*c.Listeners) < 1 {

		*c.Listeners = append(*c.Listeners, "127.0.0.1:11047")
	}

	if *c.RPCCert == "" {

		*c.RPCCert = filepath.Join(*c.DataDir, "rpc.cert")
	}

	if *c.RPCKey == "" {

		*c.RPCKey = filepath.Join(*c.DataDir, "rpc.key")
	}

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
