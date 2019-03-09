package app

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
)

var App = cli.NewApp()

func Main() int {
	e := App.Run(os.Args)
	if e != nil {
		return 1
	}
	return 0
}

func init() {
	*App = cli.App{
		Name:                 "pod",
		Version:              "v0.0.1",
		Description:          "Parallelcoin Pod Suite -- All-in-one everything for Parallelcoin!",
		Copyright:            "Legacy portions derived from btcsuite/btcd under ISC licence. The remainder is already in your possession. Use it wisely.",
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "datadir",
				Value:  "~/.pod",
				Usage:  "sets the data directory base for a pod instance",
				EnvVar: "POD_DATADIR",
			},
			cli.StringFlag{
				Name:   "loglevel",
				Value:  "info",
				Usage:  "sets the base for all subsystem logging",
				EnvVar: "POD_LOGLEVEL",
			},
			cli.StringSliceFlag{
				Name:  "subsystems",
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
				Name:    "ctl",
				Aliases: []string{"c"},
				Usage:   "send RPC commands to a node or wallet and print the result",
				Action:  ctlHandle,
				Subcommands: []cli.Command{
					{
						Name:    "listcommands",
						Aliases: []string{"list", "l"},
						Usage:   "list commands available at endpoint",
						Action:  ctlHandleList,
					},
				},
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:        "rpcserver, server, s",
						Value:       "localhost:11048",
						Usage:       "set rpc password",
						Destination: &ctlConfig.RPCServer,
					},
					cli.StringFlag{
						Name:        "walletserver, ws",
						Value:       "localhost:11046",
						Usage:       "set wallet server to use",
						Destination: &ctlConfig.Wallet,
					},
					cli.StringFlag{
						Name:        "rpcusername, username, user, u",
						Value:       "user",
						Usage:       "set rpc username",
						Destination: &ctlConfig.RPCUser,
					},
					cli.StringFlag{
						Name:        "rpcpassword, password, pass, p",
						Value:       "pa55word",
						Usage:       "set rpc password",
						Destination: &ctlConfig.RPCPass,
					},
					cli.StringFlag{
						Name:        "rpccert, cert, C",
						Value:       "~/.pod/rpc.cert",
						Usage:       "set rpc password",
						Destination: &ctlConfig.RPCCert,
					},
					cli.StringFlag{
						Name:        "proxyserver, S",
						Value:       "localhost:9050",
						Usage:       "set proxy server address",
						Destination: &ctlConfig.Proxy,
					},
					cli.StringFlag{
						Name:        "proxyusername, proxyuser, U",
						Value:       "user",
						Usage:       "set proxy username",
						Destination: &ctlConfig.ProxyUser,
					},
					cli.StringFlag{
						Name:        "proxypassword, proxypass, P",
						Value:       "pa55word",
						Usage:       "set proxy password",
						Destination: &ctlConfig.ProxyPass,
					},
					cli.BoolFlag{
						Name:        "tls, T",
						Usage:       "enable tls on connections",
						Destination: &ctlConfig.TLS,
					},
					cli.BoolFlag{
						Name:        "skipverify",
						Usage:       "disable TLS certificate verification (not recommended)",
						Destination: &ctlConfig.TLSSkipVerify,
					},
					cli.BoolFlag{
						Name:        "testnet3, testnet",
						Usage:       "connect to testnet",
						Destination: &ctlConfig.TestNet3,
					},
					cli.BoolFlag{
						Name:        "simnet",
						Usage:       "connect to simnet",
						Destination: &ctlConfig.SimNet,
					},
					cli.BoolFlag{
						Name:  "useproxy, r",
						Usage: "use configured proxy",
					},
					cli.BoolFlag{
						Name:  "wallet, w",
						Usage: "use configured proxy",
					},
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
}
