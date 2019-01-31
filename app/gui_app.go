package app

import (
	"fmt"
	"os"

	walletmain "git.parallelcoin.io/pod/cmd/wallet"
	"git.parallelcoin.io/pod/pkg/gui"
	"git.parallelcoin.io/pod/pkg/netparams"
	"github.com/tucnak/climax"
)

// GUICfg is the type for the default config data
type GUICfg struct {
	AppDataDir string
	Password   string
	PublicPass string
	Seed       []byte
	Network    string
}

// GUIConfig is
var GUIConfig = GUICfg{
	AppDataDir: walletmain.DefaultAppDataDir,
	Network:    "mainnet",
}

// GUICommand is a command to send RPC queries to bitcoin RPC protocol server for node and wallet queries
var GUICommand = climax.Command{
	Name:  "gui",
	Brief: "runs the GUI",
	Help:  "if given CLI parameters, creates a new wallet, otherwise launches the GUI",
	Flags: []climax.Flag{
		t("help", "h", "show help text"),
		s("datadir", "D", walletmain.DefaultAppDataDir, "specify where the wallet will be created"),
		f("network", "mainnet", "connect to (mainnet|testnet|simnet)"),
	},
	Handle: func(ctx climax.Context) int {
		if ctx.Is("help") {
			fmt.Print(`Usage: createwallet [-h] [-D] [--network]

launches the GUI given the data directory provided
			
Available options:

	-h, --help
		show help text
	-D, --datadir="~/.pod/wallet"
		data directory to use
	--network="mainnet"
		connect to (mainnet|testnet|simnet)

`)
			os.Exit(0)
		}
		if r, ok := getIfIs(&ctx, "datadir"); ok {
			CreateConfig.DataDir = r
		}
		if r, ok := getIfIs(&ctx, "network"); ok {
			switch r {
			case "testnet":
				walletmain.ActiveNet = &netparams.TestNet3Params
			case "simnet":
				walletmain.ActiveNet = &netparams.SimNetParams
			default:
				walletmain.ActiveNet = &netparams.MainNetParams
			}
			GUIConfig.Network = r
		}
		fmt.Println("launching GUI")
		gui.GUI()
		return 0
	},
}
