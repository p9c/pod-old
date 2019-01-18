package pod

import "github.com/tucnak/climax"

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
