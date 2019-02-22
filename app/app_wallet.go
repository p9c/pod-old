package app

import (
	"fmt"

	n "git.parallelcoin.io/pod/cmd/node"
	w "git.parallelcoin.io/pod/cmd/wallet"
	"git.parallelcoin.io/pod/pkg/netparams"
	"github.com/tucnak/climax"
)


// WalletCfg is the combined app and logging configuration data
type WalletCfg struct {
	Wallet    *w.Config
	Levels    map[string]string
	activeNet *netparams.Params
}


// WalletCommand is a command to send RPC queries to bitcoin RPC protocol server for node and wallet queries
var WalletCommand = climax.Command{
	Name:  "wallet",
	Brief: "parallelcoin wallet",
	Help:  "check balances, make payments, manage contacts",
	Flags: []climax.Flag{
		t("version", "V", "show version number and quit"),

		s("configfile", "C", w.DefaultConfigFilename,
			"path to configuration file"),
		s("datadir", "D", w.DefaultDataDir,
			"set the pod base directory"),
		f("appdatadir", w.DefaultAppDataDir, "set app data directory for wallet, configuration and logs"),

		t("init", "i", "resets configuration to defaults"),
		t("save", "S", "saves current flags into configuration"),

		t("createtemp", "", "create temporary wallet (pass=walletpass) requires --datadir"),

		t("gui", "G", "launch GUI"),
		f("rpcconnect", n.DefaultRPCListener, "connect to the RPC of a parallelcoin node for chain queries"),

		f("podusername", "user", "username for node RPC authentication"),
		f("podpassword", "pa55word", "password for node RPC authentication"),

		f("walletpass", "", "the public wallet password - only required if the wallet was created with one"),

		f("noinitialload", "false", "defer wallet load to be triggered by RPC"),
		f("network", "mainnet", "connect to (mainnet|testnet|regtestnet|simnet)"),

		f("profile", "false", "enable HTTP profiling on given port (1024-65536)"),

		f("rpccert", w.DefaultRPCCertFile,
			"file containing the RPC tls certificate"),
		f("rpckey", w.DefaultRPCKeyFile,
			"file containing RPC TLS key"),
		f("onetimetlskey", "false", "generate a new TLS certpair don't save key"),
		f("cafile", w.DefaultCAFile, "certificate authority for custom TLS CA"),
		f("enableclienttls", "false", "enable TLS for the RPC client"),
		f("enableservertls", "false", "enable TLS on wallet RPC server"),

		f("proxy", "", "proxy address for outbound connections"),
		f("proxyuser", "", "username for proxy server"),
		f("proxypass", "", "password for proxy server"),

		f("legacyrpclisteners", w.DefaultListener, "add a listener for the legacy RPC"),
		f("legacyrpcmaxclients", fmt.Sprint(w.DefaultRPCMaxClients),
			"max connections for legacy RPC"),
		f("legacyrpcmaxwebsockets", fmt.Sprint(w.DefaultRPCMaxWebsockets),
			"max websockets for legacy RPC"),

		f("username", "user",
			"username for wallet RPC when podusername is empty"),
		f("password", "pa55word",
			"password for wallet RPC when podpassword is omitted"),
		f("experimentalrpclisteners", "",
			"listener for experimental rpc"),

		s("debuglevel", "d", "info", "sets debuglevel, specify per-library below"),

		l("lib-addrmgr"), l("lib-blockchain"), l("lib-connmgr"), l("lib-database-ffldb"), l("lib-database"), l("lib-mining-cpuminer"), l("lib-mining"), l("lib-netsync"), l("lib-peer"), l("lib-rpcclient"), l("lib-txscript"), l("node"), l("node-mempool"), l("spv"), l("wallet"), l("wallet-chain"), l("wallet-legacyrpc"), l("wallet-rpcserver"), l("wallet-tx"), l("wallet-votingpool"), l("wallet-waddrmgr"), l("wallet-wallet"), l("wallet-wtxmgr"),
	},

	// Examples: []climax.Example{

	// 	{

	// 		Usecase:     "--init --rpcuser=user --rpcpass=pa55word --save",

	// 		Description: "resets the configuration file to default, sets rpc username and password and saves the changes to config after parsing",

	// 	},

	// },
}


// WalletConfig is the combined app and log levels configuration
var WalletConfig = DefaultWalletConfig(w.DefaultConfigFile)


// wf is the list of flags and the default values stored in the Usage field
var wf = GetFlags(WalletCommand)
