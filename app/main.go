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
	Commands: []cli.Command{
		{
			Name:    "ctl",
			Aliases: []string{"c"},
			Usage:   "send RPC commands to a node or wallet and print the result",
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
				fmt.Println("calling ctl")
				return nil
			},
		},
		{
			Name:    "wallet",
			Aliases: []string{"w"},
			Usage:   "start parallelcoin wallet server",
			Action: func(c *cli.Context) error {
				fmt.Println("calling ctl")
				return nil
			},
		},
		{
			Name:    "shell",
			Aliases: []string{"s"},
			Usage:   "start combined wallet/node shell",
			Action: func(c *cli.Context) error {
				fmt.Println("calling ctl")
				return nil
			},
		},
		{
			Name:    "gui",
			Aliases: []string{"g"},
			Usage:   "start GUI (TODO: should ultimately be default)",
			Action: func(c *cli.Context) error {
				fmt.Println("calling ctl")
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
