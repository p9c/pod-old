package pod

import (
	"git.parallelcoin.io/pod/lib/clog"
	"git.parallelcoin.io/pod/run/logger"
)

var Log = cl.NewSubSystem("pod", "info")
var log = Log.Ch

// SetLoggers sets the loggers according to the current configuration
func SetLoggers() {
	for i := range logger.Levels {
		switch i {
		case "log-database":

		case "log-txscript":

		case "log-peer":

		case "log-netsync":

		case "log-rpcclient":

		case "log-addrmgr":

		case "log-blockchain-indexers":

		case "log-blockchain":

		case "log-mining-cpuminer":

		case "log-mining":

		case "log-mining-controller":

		case "log-connmgr":

		case "log-spv":

		case "log-node-mempool":

		case "log-node":

		case "log-wallet-wallet":

		case "log-wallet-tx":

		case "log-wallet-votingpool":

		case "log-wallet":

		case "log-wallet-chain":

		case "log-wallet-rpc-rpcserver":

		case "log-wallet-rpc-legacyrpc":

		case "log-wallet-wtxmgr":

		}
	}
}
