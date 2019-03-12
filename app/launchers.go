package app

import (
	"github.com/davecgh/go-spew/spew"
	"gopkg.in/urfave/cli.v1"
)

func launchCtl(c *cli.Context) error {
	spew.Dump(ctlConfig)
	return nil
}

func launchGUI(c *cli.Context) error {
	return nil
}

func launchNode(c *cli.Context) error {
	spew.Dump(nodeConfig)
	return nil
}

func launchShell(c *cli.Context) error {
	return nil
}

func launchWallet(c *cli.Context) error {
	spew.Dump(walletConfig)
	return nil
}
