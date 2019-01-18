package logger

import (
	"github.com/tucnak/climax"
)

// Levels are the configured log level settings
var Levels = GetDefault()

// GetDefault returns a fresh shiny new default levels map
func GetDefault() map[string]string {
	return map[string]string{
		"log-database":             "info",
		"log-txscript":             "info",
		"log-peer":                 "info",
		"log-netsync":              "info",
		"log-rpcclient":            "info",
		"log-addrmgr":              "info",
		"log-blockchain-indexers":  "info",
		"log-blockchain":           "info",
		"log-mining-cpuminer":      "info",
		"log-mining":               "info",
		"log-mining-controller":    "info",
		"log-connmgr":              "info",
		"log-spv":                  "info",
		"log-node-mempool":         "info",
		"log-node":                 "info",
		"log-wallet-wallet":        "info",
		"log-wallet-tx":            "info",
		"log-wallet-votingpool":    "info",
		"log-wallet":               "info",
		"log-wallet-chain":         "info",
		"log-wallet-rpc-rpcserver": "info",
		"log-wallet-rpc-legacyrpc": "info",
		"log-wallet-wtxmgr":        "info",
	}
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
	if ctx.Is("log-wallet-wallet") {
		r, _ := ctx.Get("log-wallet-wallet")
		Levels["log-wallet-wallet"] = r
	}
	if ctx.Is("log-wallet-tx") {
		r, _ := ctx.Get("log-wallet-tx")
		Levels["log-wallet-tx"] = r
	}
	if ctx.Is("log-wallet-votingpool") {
		r, _ := ctx.Get("log-wallet-votingpool")
		Levels["log-wallet-votingpool"] = r
	}
	if ctx.Is("log-wallet") {
		r, _ := ctx.Get("log-wallet")
		Levels["log-wallet"] = r
	}
	if ctx.Is("log-wallet-chain") {
		r, _ := ctx.Get("log-wallet-chain")
		Levels["log-wallet-chain"] = r
	}
	if ctx.Is("log-wallet-rpc-rpcserver") {
		r, _ := ctx.Get("log-wallet-rpc")
		Levels["log-wallet-rpc"] = r
	}
	if ctx.Is("log-wallet-rpc-legacyrpc") {
		r, _ := ctx.Get("log-wallet-rpc")
		Levels["log-wallet-rpc"] = r
	}
	if ctx.Is("log-wallet-wtxmgr") {
		r, _ := ctx.Get("log-wallet-wtxmgr")
		Levels["log-wallet-wtxmgr"] = r
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
		Name:     "log-wallet-wallet",
		Usage:    "--log-wallet-wallet",
		Help:     "",
		Variable: true,
	},
	{
		Name:     "log-wallet-tx",
		Usage:    "--log-wallet-tx",
		Help:     "",
		Variable: true,
	},
	{
		Name:     "log-wallet-votingpool",
		Usage:    "--log-wallet-votingpool",
		Help:     "",
		Variable: true,
	},
	{
		Name:     "log-wallet",
		Usage:    "--log-wallet",
		Help:     "",
		Variable: true,
	},
	{
		Name:     "log-wallet-chain",
		Usage:    "--log-wallet-chain",
		Help:     "",
		Variable: true,
	},
	{
		Name:     "log-wallet-rpc-rpcserver",
		Usage:    "--log-wallet-rpc-rpcserver",
		Help:     "",
		Variable: true,
	},
	{
		Name:     "log-wallet-rpc-legacyrpc",
		Usage:    "--log-wallet-rpc-legacyrpc",
		Help:     "",
		Variable: true,
	},
	{
		Name:     "log-wallet-wtxmgr",
		Usage:    "--log-wallet-wtxmgr",
		Help:     "",
		Variable: true,
	},
}
