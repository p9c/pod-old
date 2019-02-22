package app

import (
	"git.parallelcoin.io/pod/cmd/node"
	"git.parallelcoin.io/pod/cmd/node/mempool"
	"git.parallelcoin.io/pod/cmd/spv"
	walletmain "git.parallelcoin.io/pod/cmd/wallet"
	"git.parallelcoin.io/pod/pkg/addrmgr"
	blockchain "git.parallelcoin.io/pod/pkg/chain"
	cl "git.parallelcoin.io/pod/pkg/clog"
	"git.parallelcoin.io/pod/pkg/connmgr"
	database "git.parallelcoin.io/pod/pkg/db"
	"git.parallelcoin.io/pod/pkg/db/ffldb"
	"git.parallelcoin.io/pod/pkg/mining"
	"git.parallelcoin.io/pod/pkg/mining/cpuminer"
	"git.parallelcoin.io/pod/pkg/netsync"
	"git.parallelcoin.io/pod/pkg/peer"
	"git.parallelcoin.io/pod/pkg/rpc/legacyrpc"
	"git.parallelcoin.io/pod/pkg/rpc/rpcserver"
	"git.parallelcoin.io/pod/pkg/rpcclient"
	"git.parallelcoin.io/pod/pkg/txscript"
	"git.parallelcoin.io/pod/pkg/votingpool"
	"git.parallelcoin.io/pod/pkg/waddrmgr"
	"git.parallelcoin.io/pod/pkg/wallet"
	"git.parallelcoin.io/pod/pkg/wallettx"
	chain "git.parallelcoin.io/pod/pkg/wchain"
	"git.parallelcoin.io/pod/pkg/wtxmgr"
	"github.com/tucnak/climax"
)

// LogLevels are the configured log level settings
var LogLevels = GetDefaultLogLevelsConfig()

// GetAllSubSystems returns a map with all the SubSystems in Parallelcoin Pod
func GetAllSubSystems() map[string]*cl.SubSystem {
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

// GetDefaultLogLevelsConfig returns a fresh shiny new default levels map
func GetDefaultLogLevelsConfig() map[string]string {

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

// SetAllLogging sets all the logging to a particular level
func SetAllLogging(
	level string,
) {

	ss := GetAllSubSystems()
	for i := range ss {

		ss[i].SetLevel(level)
	}
}

// SetLogging sets the logging settings according to the provided context
func SetLogging(
	ctx *climax.Context,
) {

	ss := GetAllSubSystems()
	var baselevel = "info"
	if r, ok := getIfIs(ctx, "debuglevel"); ok {

		baselevel = r
	}
	for i := range ss {

		if lvl, ok := ctx.Get(i); ok {

			ss[i].SetLevel(lvl)
		} else {
			ss[i].SetLevel(baselevel)
		}
	}
}
