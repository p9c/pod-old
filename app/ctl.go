package app

import (
	"git.parallelcoin.io/dev/pod/cmd/ctl"
	"gopkg.in/urfave/cli.v1"
)

func ctlHandle(c *cli.Context) error {

	args := c.Args()

	if len(args) < 1 {

		return cli.ShowSubcommandHelp(c)

	}

	ctl.HelpPrint = func() {
		cli.ShowSubcommandHelp(c)
	}

	ctl.Main(args, &podConfig)

	return nil
}
