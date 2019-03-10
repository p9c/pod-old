package app

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/urfave/cli"
)

func nodeHandle(c *cli.Context) error {
	fmt.Println("running node")
	if !c.Parent().IsSet("useproxy") {
		*nodeConfig.Proxy = ""
	}
	if !*nodeConfig.Onion {
		*nodeConfig.OnionProxy = ""
	}
	spew.Dump(nodeConfig)
	spew.Dump(c.Args(), c.FlagNames())
	return nil
}
