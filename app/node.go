package app

import (
	"fmt"
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
	yp, e := yaml.Marshal(appConfigCommon)
	if e == nil {
		fmt.Println(string(yp))
	}

	yn, e := yaml.Marshal(nodeConfig)
	if e == nil {
		fmt.Println(string(yn))
	}
	// spew.Dump(nodeConfig)
	// spew.Dump(c.Args(), c.FlagNames())
	return nil
}
