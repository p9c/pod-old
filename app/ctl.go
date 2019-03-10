package app

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"gopkg.in/urfave/cli.v1"
)

func ctlHandle(c *cli.Context) error {
	fmt.Println("running ctl")
	if !c.IsSet("wallet") {
		*ctlConfig.Wallet = ""
	}
	if !c.Parent().IsSet("useproxy") {
		*ctlConfig.Proxy = ""
	}
	spew.Dump(ctlConfig)
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
