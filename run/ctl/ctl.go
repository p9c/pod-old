package ctl

import (
	"fmt"

	"git.parallelcoin.io/pod/lib/clog"
	c "git.parallelcoin.io/pod/module/ctl"
	"github.com/tucnak/climax"
)

var log = clog.NewSubSystem("Ctl", clog.Ntrc)

// Command is a command to send RPC queries to bitcoin RPC protocol server for node and wallet queries
var Command = climax.Command{
	Name:  "ctl",
	Brief: "sends RPC commands and prints the reply",
	Usage: "[-l] command",
	Help:  "this tool allows you to send queries to bitcoin JSON-RPC servers and prints the reply to stdout",
	Flags: []climax.Flag{
		{
			Name:     "listcommands",
			Short:    "l",
			Usage:    `--listcommands`,
			Help:     `list available commands`,
			Variable: false,
		},
		{
			Name:     "version",
			Short:    "v",
			Usage:    `--version`,
			Help:     `show version number and quit`,
			Variable: false,
		},
	},
	Examples: []climax.Example{
		{
			Usecase:     "-l",
			Description: "lists available commands",
		},
	},
	Handle: func(ctx climax.Context) int {
		fmt.Println("ctl version", c.Version())
		if ctx.Is("version") {
			clog.Shutdown()
		}
		if ctx.Is("listcommands") {
			log.Trace.Print("listing commands")

			c.ListCommands()
		} else {
			log.Trace.Print("running command")
		}
		clog.Shutdown()
		return 0
	},
}
