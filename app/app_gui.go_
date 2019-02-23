package app

import (
	"fmt"

	"git.parallelcoin.io/pod/cmd/gui"
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

// GUICommand is a command to send RPC queries to bitcoin RPC protocol server for node and wallet queries
var GUICommand = climax.Command{
	Name:  "gui",
	Brief: "runs the GUI",
	Help:  "launches the GUI",
	Flags: []climax.Flag{
		// t("help", "h", "show help text"),
		// s("datadir", "D", walletmain.DefaultAppDataDir, "specify where the wallet will be created"),
		// f("network", "mainnet", "connect to (mainnet|testnet|simnet)"),
	},
	Handle: func(ctx climax.Context) int {
		fmt.Println("launching GUI")
		gui.GUI(ShellConfig)
		return 0
	},
}
