package app

import (
	"git.parallelcoin.io/pod/cmd/ctl"
	"github.com/urfave/cli"
)

var ctlConfig = ctl.Config{}

var ctlFlags = []cli.Flag{
	cli.StringFlag{
		Name:        "rpcserver, server, s",
		Value:       "localhost:11048",
		Usage:       "set rpc password",
		Destination: &ctlConfig.RPCServer,
	},
	cli.StringFlag{
		Name:        "rpcusername, rpcuser, u",
		Value:       "user",
		Usage:       "set rpc username",
		Destination: &ctlConfig.RPCUser,
	},
	cli.StringFlag{
		Name:        "rpcpassword, rpcpass, p",
		Value:       "pa55word",
		Usage:       "set rpc password",
		Destination: &ctlConfig.RPCPass,
	},
	cli.StringFlag{
		Name:        "rpccert, cert, C",
		Value:       "~/.pod/rpc.cert",
		Usage:       "set rpc password",
		Destination: &ctlConfig.RPCPass,
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
}
