package app

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/urfave/cli"
	"fmt"
	"git.parallelcoin.io/pod/cmd/ctl"
)

var ctlConfig = ctl.Config{}

func ctlHandle(c *cli.Context) error {
	fmt.Println("running ctl")
	if !c.IsSet("wallet") {
		ctlConfig.Wallet = ""
	}
	if !c.IsSet("useproxy") {
		ctlConfig.Proxy = ""
	}
	spew.Dump(ctlConfig)
	return nil
}

func ctlHandleList(c *cli.Context) error {
	fmt.Println("running ctl listcommands")
	ctlConfig.ListCommands = true
	if !c.IsSet("wallet") {
		ctlConfig.Wallet = ""
	}
	if !c.IsSet("useproxy") {
		ctlConfig.Proxy = ""
	}
	spew.Dump(ctlConfig)
	return nil
}
