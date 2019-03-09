package app

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
)

var App = cli.App{
	Name:                 "pod",
	Version:              "v0.9.9",
	Description:          "Parallelcoin Pod Suite -- All-in-one everything for Parallelcoin!",
	Copyright:            "Legacy portions derived from btcsuite/btcd under ISC licence. The remainder is already in your possession. Use it wisely.",
	EnableBashCompletion: true,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:   "datadir, d",
			Value:  "~/.pod",
			Usage:  "sets the data directory base for a pod instance",
			EnvVar: "POD_DATADIR",
		},
		cli.StringFlag{
			Name:   "loglevel, l",
			Value:  "info",
			Usage:  "sets the base for all subsystem logging",
			EnvVar: "POD_LOGLEVEL",
		},
		cli.StringSliceFlag{
			Name:  "subsystems, S",
			Value: &cli.StringSlice{""},
			Usage: "sets individual subsystems log levels, use 'help' to list available with list syntax",
		},
	},
	Commands: []cli.Command{
		{
			Name:    "version",
			Aliases: []string{"v"},
			Usage:   "print version and exit",
			Action: func(c *cli.Context) error {
				fmt.Println(c.App.Version)
				return nil
			},
		},
		{
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   "resets configuration to factory",
			Action: func(c *cli.Context) error {
				fmt.Println("resetting configuration")
				return nil
			},
		},
		{
			Name:    "conf",
			Aliases: []string{"C"},
			Usage:   "automate configuration setup for testnets etc",
			Action: func(c *cli.Context) error {
				fmt.Println("calling conf")
				return nil
			},
		},
		{
			Name:        "ctl",
			Aliases:     []string{"c"},
			Usage:       "send RPC commands to a node or wallet and print the result",
			Subcommands: ctlCommands,
			Flags:       ctlFlags,
			Action: func(c *cli.Context) error {
				fmt.Println("calling ctl")
				return nil
			},
		},
		{
			Name:    "node",
			Aliases: []string{"n"},
			Usage:   "start parallelcoin full node",
			Action: func(c *cli.Context) error {
				fmt.Println("calling node")
				return nil
			},
		},
		{
			Name:    "wallet",
			Aliases: []string{"w"},
			Usage:   "start parallelcoin wallet server",
			Action: func(c *cli.Context) error {
				fmt.Println("calling wallet")
				return nil
			},
		},
		{
			Name:    "shell",
			Aliases: []string{"s"},
			Usage:   "start combined wallet/node shell",
			Action: func(c *cli.Context) error {
				fmt.Println("calling shell")
				return nil
			},
		},
		{
			Name:    "gui",
			Aliases: []string{"g"},
			Usage:   "start GUI (TODO: should ultimately be default)",
			Action: func(c *cli.Context) error {
				fmt.Println("calling gui")
				return nil
			},
		},
	},
}

func Main() int {
	e := App.Run(os.Args)
	if e != nil {
		return 1
	}
	return 0
}
