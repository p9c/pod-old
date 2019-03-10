package app

import (
	"encoding/json"
	"fmt"
	// "github.com/davecgh/go-spew/spew"
	"gopkg.in/urfave/cli.v1"
)

func nodeHandle(c *cli.Context) error {
	fmt.Println("running node")
	if !c.Parent().IsSet("useproxy") {
		*nodeConfig.Proxy = ""
	}
	if !*nodeConfig.Onion {
		*nodeConfig.OnionProxy = ""
	}
	// TODO: get user input datadir flag to set file paths
	j, e := json.MarshalIndent(nodeConfig, "", "  ")
	if e == nil {
		fmt.Println(string(j))
	}
	// spew.Dump(nodeConfig)
	// spew.Dump(c.Args(), c.FlagNames())
	return nil
}
