package logger

import (
	"github.com/tucnak/climax"
)

// Levels are the configured log level settings
var Levels = map[string]string{
	"log-database":            "info",
	"log-txscript":            "info",
	"log-peer":                "info",
	"log-netsync":             "info",
	"log-rpcclient":           "info",
	"log-addrmgr":             "info",
	"log-blockchain-indexers": "info",
	"log-blockchain":          "info",
	"log-mining-cpuminer":     "info",
	"log-mining":              "info",
	"log-mining-controller":   "info",
	"log-connmgr":             "info",
	"log-spv":                 "info",
	"log-node-mempool":        "info",
	"log-node":                "info",
	"log-shell-wallet":        "info",
	"log-shell-tx":            "info",
	"log-shell-votingpool":    "info",
	"log-shell":               "info",
	"log-shell-chain":         "info",
	"log-shell-rpc-rpcserver": "info",
	"log-shell-rpc-legacyrpc": "info",
	"log-shell-wtxmgr":        "info",
}

// SetLogging sets the logging settings according to the provided context
func SetLogging(ctx *climax.Context) {
	if ctx.Is("log-database") {
		r, _ := ctx.Get("log-database")
		Levels["log-database"] = r
	}
	if ctx.Is("log-txscript") {
		r, _ := ctx.Get("log-txscript")
		Levels["log-txscript"] = r
	}
	if ctx.Is("log-peer") {
		r, _ := ctx.Get("log-peer")
		Levels["log-peer"] = r
	}
	if ctx.Is("log-netsync") {
		r, _ := ctx.Get("log-netsync")
		Levels["log-netsync"] = r
	}
	if ctx.Is("log-rpcclient") {
		r, _ := ctx.Get("log-rpcclient")
		Levels["log-rpcclient"] = r
	}
	if ctx.Is("log-addrmgr") {
		r, _ := ctx.Get("log-addrmgr")
		Levels["log-addrmgr"] = r
	}
	if ctx.Is("log-blockchain-indexers") {
		r, _ := ctx.Get("log-blockchain-indexers")
		Levels["log-blockchain-indexers"] = r
	}
	if ctx.Is("log-blockchain") {
		r, _ := ctx.Get("log-blockchain")
		Levels["log-blockchain"] = r
	}
	if ctx.Is("log-mining-cpuminer") {
		r, _ := ctx.Get("log-mining-cpuminer")
		Levels["log-mining-cpuminer"] = r
	}
	if ctx.Is("log-mining") {
		r, _ := ctx.Get("log-mining")
		Levels["log-mining"] = r
	}
	if ctx.Is("log-mining-controller") {
		r, _ := ctx.Get("log-mining-controller")
		Levels["log-mining-controller"] = r
	}
	if ctx.Is("log-connmgr") {
		r, _ := ctx.Get("log-connmgr")
		Levels["log-connmgr"] = r
	}
	if ctx.Is("log-spv") {
		r, _ := ctx.Get("log-spv")
		Levels["log-log-spv"] = r
	}
	if ctx.Is("log-node-mempool") {
		r, _ := ctx.Get("log-node-mempool")
		Levels["log-node-mempool"] = r
	}
	if ctx.Is("log-node") {
		r, _ := ctx.Get("log-node")
		Levels["log-node"] = r
	}
	if ctx.Is("log-shell-wallet") {
		r, _ := ctx.Get("log-shell-wallet")
		Levels["log-shell-wallet"] = r
	}
	if ctx.Is("log-shell-tx") {
		r, _ := ctx.Get("log-shell-tx")
		Levels["log-shell-tx"] = r
	}
	if ctx.Is("log-shell-votingpool") {
		r, _ := ctx.Get("log-shell-votingpool")
		Levels["log-shell-votingpool"] = r
	}
	if ctx.Is("log-shell") {
		r, _ := ctx.Get("log-shell")
		Levels["log-shell"] = r
	}
	if ctx.Is("log-shell-chain") {
		r, _ := ctx.Get("log-shell-chain")
		Levels["log-shell-chain"] = r
	}
	if ctx.Is("log-shell-rpc-rpcserver") {
		r, _ := ctx.Get("log-shell-rpc")
		Levels["log-shell-rpc"] = r
	}
	if ctx.Is("log-shell-rpc-legacyrpc") {
		r, _ := ctx.Get("log-shell-rpc")
		Levels["log-shell-rpc"] = r
	}
	if ctx.Is("log-shell-wtxmgr") {
		r, _ := ctx.Get("log-shell-wtxmgr")
		Levels["log-shell-wtxmgr"] = r
	}
}

var debugLevels = []climax.Flag{
	{
		Name:     "log-database",
		Usage:    "--log-database",
		Help:     "",
		Variable: true,
	},
	{
		Name:     "log-txscript",
		Usage:    "--log-txscript",
		Help:     "",
		Variable: true,
	},
	{
		Name:     "log-peer",
		Usage:    "--log-peer",
		Help:     "",
		Variable: true,
	},
	{
		Name:     "log-netsync",
		Usage:    "--log-netsync",
		Help:     "",
		Variable: true,
	},
	{
		Name:     "log-rpcclient",
		Usage:    "--log-rpcclient",
		Help:     "",
		Variable: true,
	},
	{
		Name:     "addrmgr",
		Usage:    "--log-addrmgr",
		Help:     "",
		Variable: true,
	},
	{
		Name:     "log-blockchain-indexers",
		Usage:    "--log-blockchain-indexers",
		Help:     "",
		Variable: true,
	},
	{
		Name:     "log-blockchain",
		Usage:    "--log-blockchain",
		Help:     "",
		Variable: true,
	},
	{
		Name:     "log-mining-cpuminer",
		Usage:    "--log-mining-cpuminer",
		Help:     "",
		Variable: true,
	},
	{
		Name:     "log-mining",
		Usage:    "--log-mining",
		Help:     "",
		Variable: true,
	},
	{
		Name:     "log-mining-controller",
		Usage:    "--log-mining-controller",
		Help:     "",
		Variable: true,
	},
	{
		Name:     "log-connmgr",
		Usage:    "--log-connmgr",
		Help:     "",
		Variable: true,
	},
	{
		Name:     "log-spv",
		Usage:    "--log-spv",
		Help:     "",
		Variable: true,
	},
	{
		Name:     "log-node-mempool",
		Usage:    "--log-node-mempool",
		Help:     "",
		Variable: true,
	},
	{
		Name:     "log-node",
		Usage:    "--log-node",
		Help:     "",
		Variable: true,
	},
	{
		Name:     "log-shell-wallet",
		Usage:    "--log-shell-wallet",
		Help:     "",
		Variable: true,
	},
	{
		Name:     "log-shell-tx",
		Usage:    "--log-shell-tx",
		Help:     "",
		Variable: true,
	},
	{
		Name:     "log-shell-votingpool",
		Usage:    "--log-shell-votingpool",
		Help:     "",
		Variable: true,
	},
	{
		Name:     "log-shell",
		Usage:    "--log-shell",
		Help:     "",
		Variable: true,
	},
	{
		Name:     "log-shell-chain",
		Usage:    "--log-shell-chain",
		Help:     "",
		Variable: true,
	},
	{
		Name:     "log-shell-rpc-rpcserver",
		Usage:    "--log-shell-rpc-rpcserver",
		Help:     "",
		Variable: true,
	},
	{
		Name:     "log-shell-rpc-legacyrpc",
		Usage:    "--log-shell-rpc-legacyrpc",
		Help:     "",
		Variable: true,
	},
	{
		Name:     "log-shell-wtxmgr",
		Usage:    "--log-shell-wtxmgr",
		Help:     "",
		Variable: true,
	},
}
