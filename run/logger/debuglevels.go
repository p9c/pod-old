package logger

import (
	"git.parallelcoin.io/pod/lib/addrmgr"
	"git.parallelcoin.io/pod/lib/blockchain"
	"git.parallelcoin.io/pod/lib/clog"
	"git.parallelcoin.io/pod/lib/connmgr"
	"git.parallelcoin.io/pod/lib/database"
	"git.parallelcoin.io/pod/lib/database/ffldb"
	"git.parallelcoin.io/pod/lib/mining"
	"git.parallelcoin.io/pod/lib/mining/cpuminer"
	"git.parallelcoin.io/pod/lib/netsync"
	"git.parallelcoin.io/pod/lib/peer"
	"git.parallelcoin.io/pod/lib/rpcclient"
	"git.parallelcoin.io/pod/lib/txscript"
	"git.parallelcoin.io/pod/module/node"
	"git.parallelcoin.io/pod/module/node/mempool"
	"git.parallelcoin.io/pod/module/spv"
	"git.parallelcoin.io/pod/module/wallet"
	"git.parallelcoin.io/pod/module/wallet/chain"
	"git.parallelcoin.io/pod/module/wallet/rpc/legacyrpc"
	"git.parallelcoin.io/pod/module/wallet/rpc/rpcserver"
	"git.parallelcoin.io/pod/module/wallet/tx"
	"git.parallelcoin.io/pod/module/wallet/votingpool"
	"git.parallelcoin.io/pod/module/wallet/waddrmgr"
	"git.parallelcoin.io/pod/module/wallet/wallet"
	"git.parallelcoin.io/pod/module/wallet/wtxmgr"
	"git.parallelcoin.io/pod/run/util"
	"github.com/tucnak/climax"
)

// Levels are the configured log level settings
var Levels = GetDefault()

// GetDefaultConfig returns a fresh shiny new default levels map
func GetDefaultConfig() map[string]string {
	return map[string]string{
		"lib-addrmgr":         "info",
		"lib-blockchain":      "info",
		"lib-connmgr":         "info",
		"lib-database-ffldb":  "info",
		"lib-database":        "info",
		"lib-mining-cpuminer": "info",
		"lib-mining":          "info",
		"lib-netsync":         "info",
		"lib-peer":            "info",
		"lib-rpcclient":       "info",
		"lib-txscript":        "info",
		"node":                "info",
		"node-mempool":        "info",
		"spv":                 "info",
		"wallet":              "info",
		"wallet-chain":        "info",
		"wallet-legacyrpc":    "info",
		"wallet-rpcserver":    "info",
		"wallet-tx":           "info",
		"wallet-votingpool":   "info",
		"wallet-waddrmgr":     "info",
		"wallet-wallet":       "info",
		"wallet-wtxmgr":       "info",
	}
}

// GetDefault returns a fresh shiny new default levels map
func GetDefault() map[string]*cl.SubSystem {
	return map[string]*cl.SubSystem{
		"lib-addrmgr":         addrmgr.Log,
		"lib-blockchain":      blockchain.Log,
		"lib-connmgr":         connmgr.Log,
		"lib-database-ffldb":  ffldb.Log,
		"lib-database":        database.Log,
		"lib-mining-cpuminer": cpuminer.Log,
		"lib-mining":          mining.Log,
		"lib-netsync":         netsync.Log,
		"lib-peer":            peer.Log,
		"lib-rpcclient":       rpcclient.Log,
		"lib-txscript":        txscript.Log,
		"node":                node.Log,
		"node-mempool":        mempool.Log,
		"spv":                 spv.Log,
		"wallet":              walletmain.Log,
		"wallet-chain":        chain.Log,
		"wallet-legacyrpc":    legacyrpc.Log,
		"wallet-rpcserver":    rpcserver.Log,
		"wallet-tx":           wallettx.Log,
		"wallet-votingpool":   votingpool.Log,
		"wallet-waddrmgr":     waddrmgr.Log,
		"wallet-wallet":       wallet.Log,
		"wallet-wtxmgr":       wtxmgr.Log,
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

func setIfIs(ctx *climax.Context, name string) {
	var r *string
	t := ""
	r = &t
	if getIfIs(ctx, name, r) {
		Levels[name].SetLevel(*r)
	}
}

// SetLogging sets the logging settings according to the provided context
func SetLogging(ctx *climax.Context) {
	for i := range Levels {
		setIfIs(ctx, i)
	}
}

// SetAllLogging sets all the logging to a particular level
func SetAllLogging(level string) {
	for i := range Levels {
		Levels[i].SetLevel(level)
	}
}

var l = pu.GenLog

var debugLevels = []climax.Flag{
	l("lib-blockchain"), l("lib-connmgr"), l("lib-database-ffldb"), l("lib-database"), l("lib-mining-cpuminer"), l("lib-mining"), l("lib-netsync"), l("lib-peer"), l("lib-rpcclient"), l("lib-txscript"), l("node"), l("node-mempool"), l("spv"), l("wallet"), l("wallet-chain"), l("wallet-legacyrpc"), l("wallet-rpcserver"), l("wallet-tx"), l("wallet-votingpool"), l("wallet-waddrmgr"), l("wallet-wallet"), l("wallet-wtxmgr")}
