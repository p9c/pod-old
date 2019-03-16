package app_old

import (
	"fmt"
	"path/filepath"

	w "git.parallelcoin.io/dev/pod/cmd/wallet"
	walletmain "git.parallelcoin.io/dev/pod/cmd/wallet"
	netparams "git.parallelcoin.io/dev/pod/pkg/chain/config/params"
	"git.parallelcoin.io/dev/pod/pkg/wallet"
	"github.com/tucnak/climax"
)

// SetupCfg is the type for the default config data
type SetupCfg struct {
	DataDir string
	Network string
	Config  *walletmain.Config
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
		fmt.Println("pod wallet setup")
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
			return 0
		}
		SetupConfig.DataDir = w.DefaultDataDir
		if r, ok := getIfIs(&ctx, "datadir"); ok {
			SetupConfig.DataDir = r
		}
		activeNet := walletmain.ActiveNet
		wc := DefaultWalletConfig(SetupConfig.DataDir)
		SetupConfig.Config = wc.Wallet
		SetupConfig.Config.TestNet3 = false
		SetupConfig.Config.SimNet = false
		if r, ok := getIfIs(&ctx, "network"); ok {
			switch r {
			case "testnet":
				activeNet = &netparams.TestNet3Params
				SetupConfig.Config.TestNet3 = true
				SetupConfig.Config.SimNet = false
			case "simnet":
				activeNet = &netparams.SimNetParams
				SetupConfig.Config.TestNet3 = false
				SetupConfig.Config.SimNet = true
			default:
				activeNet = &netparams.MainNetParams
			}
			SetupConfig.Network = r
		}
		dbDir := walletmain.NetworkDir(
			filepath.Join(SetupConfig.DataDir, "wallet"), activeNet.Params)
		loader := wallet.NewLoader(
			walletmain.ActiveNet.Params, dbDir, 250)
		exists, err := loader.WalletExists()
		if err != nil {
			fmt.Println("ERROR", err)
			return 1
		}
		if exists {
			fmt.Print("\n!!! A wallet already exists at '" + dbDir + "/wallet.db' !!! \n")
			fmt.Println(`if you are sure it isn't valuable you can delete it before running this again:

	rm ` + dbDir + `/wallet.db			
`)
			return 1
		}
		SetupConfig.Config.AppDataDir = filepath.Join(
			SetupConfig.DataDir, "wallet")
		WriteDefaultConfConfig(SetupConfig.DataDir)
		WriteDefaultCtlConfig(SetupConfig.DataDir)
		WriteDefaultNodeConfig(SetupConfig.DataDir)
		WriteDefaultWalletConfig(SetupConfig.DataDir)
		WriteDefaultShellConfig(SetupConfig.DataDir)
		e := walletmain.CreateWallet(SetupConfig.Config, activeNet)
		if e != nil {
			panic(e)
		}
		fmt.Print("\nYou can now open the wallet\n")
		return 0
	},
}

// SetupConfig is
var SetupConfig = SetupCfg{
	DataDir: walletmain.DefaultAppDataDir,
	Network: "mainnet",
}
