package app

import (
	"fmt"

	"git.parallelcoin.io/dev/pod/cmd/node"
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

	// spew.Dump(nodeConfig)
	err := node.Main(nodeConfig, activeNetParams, nil)
	if err != nil {
		fmt.Print(err)
		return err
	}

	return nil
}

func launchShell(c *cli.Context) error {

	return nil
}

func launchWallet(c *cli.Context) error {

	spew.Dump(walletConfig)
	return nil
}
