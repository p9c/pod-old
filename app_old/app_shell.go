package app

import (
	"fmt"
	"path/filepath"

	n "git.parallelcoin.io/pod/cmd/node"
	w "git.parallelcoin.io/pod/cmd/wallet"
	"github.com/tucnak/climax"
)

// DefaultShellAppDataDir is the default app data dir
var DefaultShellAppDataDir = filepath.Join(w.DefaultDataDir, "shell")

// DefaultShellConfigFile is the default configfile for shell
var DefaultShellConfigFile = filepath.Join(DefaultShellAppDataDir, "conf.json")

// ShellCommand is a command to send RPC queries to bitcoin RPC protocol server for node and wallet queries
var ShellCommand = climax.Command{
	Name:  "shell",
	Brief: "parallelcoin shell",
	Help:  "check balances, make payments, manage contacts, search the chain, it slices, it dices",
	Flags: []climax.Flag{
		t("version", "V", "show version number and quit"),

		s("configfile", "C", DefaultShellConfFileName, "path to configuration file"),

		s("datadir", "D", n.DefaultDataDir, "set the pod base directory"),
		f("appdatadir", "shell", "set app data directory for wallet, configuration and logs"),

		t("init", "i", "resets configuration to defaults"),
		t("save", "S", "saves current flags into configuration"),

		f("network", "mainnet",
			"connect to (mainnet|testnet|regtestnet|simnet)"),

		f("createtemp", "false",
			"create temporary wallet (pass=walletpass) requires --datadir"),

		f("walletpass", "",
			"the public wallet password - only required if the wallet was created with one"),

		s("listeners", "S", n.DefaultListener,
			"sets an address to listen for P2P connections"),
		f("externalips", "", "additional P2P listeners"),
		f("disablelisten", "false", "disables the P2P listener"),

		f("rpclisteners", "127.0.0.1:11046",
			"add a listener for the wallet RPC"),
		f("rpcmaxclients", fmt.Sprint(n.DefaultMaxRPCClients),
			"max connections for wallet RPC"),
		f("rpcmaxwebsockets", fmt.Sprint(n.DefaultMaxRPCWebsockets),
			"max websockets for wallet RPC"),

		f("username", "user", "username for wallet RPC"),
		f("password", "pa55word", "password for wallet RPC"),

		f("rpccert", n.DefaultRPCCertFile,
			"file containing the RPC tls certificate"),
		f("rpckey", n.DefaultRPCKeyFile,
			"file containing RPC TLS key"),
		f("onetimetlskey", "false",
			"generate a new TLS certpair don't save key"),
		f("cafile", w.DefaultCAFile,
			"certificate authority for custom TLS CA"),
		f("tls", "false", "enable TLS on wallet RPC server"),

		f("txindex", "true", "enable transaction index"),
		f("addrindex", "true", "enable address index"),
		t("dropcfindex", "", "delete committed filtering (CF) index then exit"),
		t("droptxindex", "", "deletes transaction index then exit"),
		t("dropaddrindex", "", "deletes the address index then exits"),

		f("proxy", "", "proxy address for outbound connections"),
		f("proxyuser", "", "username for proxy server"),
		f("proxypass", "", "password for proxy server"),

		f("onion", "", "connect via tor proxy relay"),
		f("onionuser", "", "username for onion proxy server"),
		f("onionpass", "", "password for onion proxy server"),
		f("noonion", "false", "disable onion proxy"),
		f("torisolation", "false", "use a different user/pass for each peer"),

		f("addpeers", "",
			"adds a peer to the peers database to try to connect to"),
		f("connectpeers", "",
			"adds a peer to a connect-only whitelist"),
		f(`maxpeers`, fmt.Sprint(n.DefaultMaxPeers),
			"sets max number of peers to connect to to at once"),
		f(`disablebanning`, "false",
			"disable banning of misbehaving peers"),
		f("banduration", "1d",
			"time to ban misbehaving peers (d/h/m/s)"),
		f("banthreshold", fmt.Sprint(n.DefaultBanThreshold),
			"banscore that triggers a ban"),
		f("whitelists", "", "addresses and networks immune to banning"),

		f("trickleinterval", fmt.Sprint(n.DefaultTrickleInterval),
			"time between sending inventory batches to peers"),
		f("minrelaytxfee", "0",
			"min fee in DUO/kb to relay transaction"),
		f("freetxrelaylimit", fmt.Sprint(n.DefaultFreeTxRelayLimit),
			"limit below min fee transactions in kb/bin"),
		f("norelaypriority", "false",
			"do not discriminate transactions for relaying"),

		f("nopeerbloomfilters", "false",
			"disable bloom filtering support"),
		f("nocfilters", "false",
			"disable committed filtering (CF) support"),
		f("blocksonly", "false", "do not accept transactions from peers"),
		f("relaynonstd", "false", "relay nonstandard transactions"),
		f("rejectnonstd", "false", "reject nonstandard transactions"),

		f("maxorphantxs", fmt.Sprint(n.DefaultMaxOrphanTransactions),
			"max number of orphan transactions to store"),
		f("sigcachemaxsize", fmt.Sprint(n.DefaultSigCacheMaxSize),
			"maximum number of signatures to store in memory"),

		f("generate", fmt.Sprint(n.DefaultGenerate),
			"set CPU miner to generate blocks"),
		f("genthreads", fmt.Sprint(n.DefaultGenThreads),
			"set number of threads to generate using CPU, -1 = all"),
		f("algo", n.DefaultAlgo, "set algorithm to be used by cpu miner"),
		f("miningaddrs", "", "add address to pay block rewards to"),
		f("minerlistener", n.DefaultMinerListener,
			"address to listen for mining work subscriptions"),
		f("minerpass", "",
			"PSK to prevent snooping/spoofing of miner traffic"),

		f("addcheckpoints", "false", `add custom checkpoints "height:hash"`),
		f("disablecheckpoints", "false", "disable all checkpoints"),

		f("blockminsize", fmt.Sprint(n.DefaultBlockMinSize),
			"min block size for miners"),
		f("blockmaxsize", fmt.Sprint(n.DefaultBlockMaxSize),
			"max block size for miners"),
		f("blockminweight", fmt.Sprint(n.DefaultBlockMinWeight),
			"min block weight for miners"),
		f("blockmaxweight", fmt.Sprint(n.DefaultBlockMaxWeight),
			"max block weight for miners"),
		f("blockprioritysize", fmt.Sprint(),
			"size in bytes of high priority blocks"),

		f("uacomment", "", "comment to add to the P2P network user agent"),
		f("upnp", "false", "use UPNP to automatically port forward to node"),
		f("dbtype", "ffldb", "set database backend type"),
		f("disablednsseed", "false", "disable dns seeding"),

		f("profile", "false", "start HTTP profiling server on given address"),
		f("cpuprofile", "false", "start cpu profiling server on given address"),

		s("debuglevel", "d", "info", "sets debuglevel, specify per-library below"),

		l("lib-addrmgr"), l("lib-blockchain"), l("lib-connmgr"), l("lib-database-ffldb"), l("lib-database"), l("lib-mining-cpuminer"), l("lib-mining"), l("lib-netsync"), l("lib-peer"), l("lib-rpcclient"), l("lib-txscript"), l("node"), l("node-mempool"), l("spv"), l("wallet"), l("wallet-chain"), l("wallet-legacyrpc"), l("wallet-rpcserver"), l("wallet-tx"), l("wallet-votingpool"), l("wallet-waddrmgr"), l("wallet-wallet"), l("wallet-wtxmgr"),
	},
	Examples: []climax.Example{
		{
			Usecase:     "--init --rpcuser=user --rpcpass=pa55word --save",
			Description: "resets the configuration file to default, sets rpc username and password and saves the changes to config after parsing",
		},
	},
	Handle: shellHandle,
}

// ShellConfig is the combined app and log levels configuration
var ShellConfig = DefaultShellConfig(
	w.DefaultDataDir,
)
