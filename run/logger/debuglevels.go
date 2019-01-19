package logger

import (
	"git.parallelcoin.io/pod/run/util"
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

func getIfIs(ctx *climax.Context, name string, r *string) (ok bool) {
	if ctx.Is(name) {
		var s string
		s, ok = ctx.Get(name)
		r = &s
	}
	return
}

// SetLogging sets the logging settings according to the provided context
func SetLogging(ctx *climax.Context) {
	var r *string
	t := ""
	r = &t
	if getIfIs(ctx, "log-database", r) {
		Levels["log-database"] = *r
	}
	if getIfIs(ctx, "log-txscript", r) {
		Levels["log-txscript"] = *r
	}
	if getIfIs(ctx, "log-peer", r) {
		Levels["log-peer"] = *r
	}
	if getIfIs(ctx, "log-netsync", r) {
		Levels["log-netsync"] = *r
	}
	if getIfIs(ctx, "log-rpcclient", r) {
		Levels["log-rpcclient"] = *r
	}
	if getIfIs(ctx, "log-addrmgr", r) {
		Levels["log-addrmgr"] = *r
	}
	if getIfIs(ctx, "log-blockchain-indexers", r) {
		Levels["log-blockchain-indexers"] = *r
	}
	if getIfIs(ctx, "log-blockchain", r) {
		Levels["log-blockchain"] = *r
	}
	if getIfIs(ctx, "log-mining-cpuminer", r) {
		Levels["log-mining-cpuminer"] = *r
	}
	if getIfIs(ctx, "log-mining", r) {
		Levels["log-mining"] = *r
	}
	if getIfIs(ctx, "log-mining-controller", r) {
		Levels["log-mining-controller"] = *r
	}
	if getIfIs(ctx, "log-connmgr", r) {
		Levels["log-connmgr"] = *r
	}
	if getIfIs(ctx, "log-spv", r) {
		Levels["log-log-spv"] = *r
	}
	if getIfIs(ctx, "log-node-mempool", r) {
		Levels["log-node-mempool"] = *r
	}
	if getIfIs(ctx, "log-node", r) {
		Levels["log-node"] = *r
	}
	if getIfIs(ctx, "log-wallet-wallet", r) {
		Levels["log-wallet-wallet"] = *r
	}
	if getIfIs(ctx, "log-wallet-tx", r) {
		Levels["log-wallet-tx"] = *r
	}
	if getIfIs(ctx, "log-wallet-votingpool", r) {
		Levels["log-wallet-votingpool"] = *r
	}
	if getIfIs(ctx, "log-wallet", r) {
		Levels["log-wallet"] = *r
	}
	if getIfIs(ctx, "log-wallet-chain", r) {
		Levels["log-wallet-chain"] = *r
	}
	if getIfIs(ctx, "log-wallet-rpc", r) {
		Levels["log-wallet-rpc"] = *r
	}
	if getIfIs(ctx, "log-wallet-rpc", r) {
		Levels["log-wallet-rpc"] = *r
	}
	if getIfIs(ctx, "log-wallet-wtxmgr", r) {
		Levels["log-wallet-wtxmgr"] = *r
	}
}

var debugLevels = []climax.Flag{
	podutil.GenerateFlag("log-database", "", "--log-database", "", true),
	podutil.GenerateFlag("log-txscript", "", "--log-txscript", "", true),
	podutil.GenerateFlag("log-peer", "", "--log-peer", "", true),
	podutil.GenerateFlag("log-netsync", "", "--log-netsync", "", true),
	podutil.GenerateFlag("log-rpcclient", "", "--log-rpcclient", "", true),
	podutil.GenerateFlag("addrmgr", "", "--log-addrmgr", "", true),
	podutil.GenerateFlag("log-blockchain-indexers", "", "--log-blockchain-indexers", "", true),
	podutil.GenerateFlag("log-blockchain", "", "--log-blockchain", "", true),
	podutil.GenerateFlag("log-mining-cpuminer", "", "--log-mining-cpuminer", "", true),
	podutil.GenerateFlag("log-mining", "", "--log-mining", "", true),
	podutil.GenerateFlag("log-mining-controller", "", "--log-mining-controller", "", true),
	podutil.GenerateFlag("log-connmgr", "", "--log-connmgr", "", true),
	podutil.GenerateFlag("log-spv", "", "--log-spv", "", true),
	podutil.GenerateFlag("log-node-mempool", "", "--log-node-mempool", "", true),
	podutil.GenerateFlag("log-node", "", "--log-node", "", true),
	podutil.GenerateFlag("log-wallet-wallet", "", "--log-wallet-wallet", "", true),
	podutil.GenerateFlag("log-wallet-tx", "", "--log-wallet-tx", "", true),
	podutil.GenerateFlag("log-wallet-votingpool", "", "--log-wallet-votingpool", "", true),
	podutil.GenerateFlag("log-wallet", "", "--log-wallet", "", true),
	podutil.GenerateFlag("log-wallet-chain", "", "--log-wallet-chain", "", true),
	podutil.GenerateFlag("log-wallet-rpc-rpcserver", "", "--log-wallet-rpc-rpcserver", "", true),
	podutil.GenerateFlag("log-wallet-rpc-legacyrpc", "", "--log-wallet-rpc-legacyrpc", "", true),
	podutil.GenerateFlag("log-wallet-wtxmgr", "", "--log-wallet-wtxmgr", "", true),
}
