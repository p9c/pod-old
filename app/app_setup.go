package app

import (
	"fmt"
	"os"
	"path/filepath"

	w "git.parallelcoin.io/pod/cmd/wallet"
	walletmain "git.parallelcoin.io/pod/cmd/wallet"
	"git.parallelcoin.io/pod/pkg/netparams"
	"github.com/tucnak/climax"
)

// SetupCfg is the type for the default config data
type SetupCfg struct {
	DataDir string
	Network string
	Config  *walletmain.Config
}

// SetupConfig is
var SetupConfig = SetupCfg{
	DataDir: walletmain.DefaultAppDataDir,
	Network: "mainnet",
}

// SetupCommand is a command to send RPC queries to bitcoin RPC protocol server for node and wallet queries
var SetupCommand = climax.Command{
	Name:  "setup",
	Brief: "initialises configuration and creates a new wallet",
	Help:  "initialises configuration and creates a new wallet in specified data directory for a specified network",
	Flags: []climax.Flag{
		t("help", "h", "show help text"),
		s("datadir", "D", walletmain.DefaultAppDataDir, "specify where the wallet will be created"),
		f("network", "mainnet", "connect to (mainnet|testnet|simnet)"),
	},
	Handle: func(ctx climax.Context) int {
		if ctx.Is("help") {
			fmt.Print(`Usage: create [-h] [-D] [--network]

creates a new wallet given CLI flags, or interactively
			
Available options:

	-h, --help
		show help text
	-D, --datadir="~/.pod/wallet"
		specify where the wallet will be created
	--network="mainnet"
		connect to (mainnet|testnet|simnet)

`)
			os.Exit(0)
		}
		SetupConfig.DataDir = w.DefaultDataDir
		if r, ok := getIfIs(&ctx, "datadir"); ok {
			SetupConfig.DataDir = r
		}
		WriteDefaultConfConfig(SetupConfig.DataDir)
		WriteDefaultCtlConfig(SetupConfig.DataDir)
		WriteDefaultNodeConfig(SetupConfig.DataDir)
		WriteDefaultWalletConfig(SetupConfig.DataDir)
		WriteDefaultShellConfig(SetupConfig.DataDir)
		activeNet := walletmain.ActiveNet
		if r, ok := getIfIs(&ctx, "network"); ok {
			switch r {
			case "testnet":
				activeNet = &netparams.TestNet3Params
			case "simnet":
				activeNet = &netparams.SimNetParams
			default:
				activeNet = &netparams.MainNetParams
			}
			SetupConfig.Network = r
		}

		SetupConfig.Config = WalletConfig.Wallet
		if SetupConfig.Config.TestNet3 {
			fmt.Println("using testnet")
			activeNet = &netparams.TestNet3Params
		}
		if SetupConfig.Config.SimNet {
			fmt.Println("using simnet")
			activeNet = &netparams.SimNetParams
		}
		SetupConfig.Config.AppDataDir = filepath.Join(
			SetupConfig.DataDir, "wallet")
		// fmt.Println(activeNet.Name)
		// spew.Dump(SetupConfig)
		walletmain.CreateWallet(SetupConfig.Config, activeNet)
		fmt.Print("\nYou can now open the wallet\n")
		return 0
	},
}
