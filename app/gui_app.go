package app

import (
	"fmt"
	"os"
	"time"

	walletmain "git.parallelcoin.io/pod/cmd/wallet"
	"git.parallelcoin.io/pod/pkg/gui"
	"git.parallelcoin.io/pod/pkg/netparams"
	"git.parallelcoin.io/pod/pkg/util/hdkeychain"
	"git.parallelcoin.io/pod/pkg/wallet"
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
		s("appdatadir", "d", walletmain.DefaultAppDataDir, "specify where the wallet will be created"),
		s("seed", "s", "", "input pre-existing seed"),
		s("password", "p", "", "specify password for private data"),
		s("publicpass", "P", "", "specify password for public data"),
		t("cli", "c", "use commandline interface interactive input"),
		f("network", "mainnet", "connect to (mainnet|testnet|simnet)"),
	},
	Handle: func(ctx climax.Context) int {
		if ctx.Is("help") {
			fmt.Print(`Usage: createwallet [-h] [-d] [-s] [-p] [-P] [-c] [--network]

creates a new wallet given CLI flags, or launches the GUI which creates a new wallet if necessary
			
Available options:

	-h, --help
		show help text
	-d, --appdatadir="/loki/.pod/wallet"
		specify where the wallet will be created
	-s, --seed=""
		input pre-existing seed
	-p, --password=""
		specify password for private data
	-P, --publicpass=""
		specify password for public data
	-c, --cli
		use commandline interface interactive input
	--network="mainnet"
		connect to (mainnet|testnet|simnet)

`)
			os.Exit(0)
		}
		argsGiven := false
		if r, ok := getIfIs(&ctx, "appdatadir"); ok {
			GUIConfig.AppDataDir = r
			argsGiven = true
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
			argsGiven = true
		}
		if ctx.Is("cli") {
			walletmain.CreateWallet(&walletmain.Config{
				AppDataDir: GUIConfig.AppDataDir,
				WalletPass: GUIConfig.PublicPass,
			})
			fmt.Print("\nYou can now open the wallet\n")
			os.Exit(0)
		}
		if r, ok := getIfIs(&ctx, "seed"); ok {
			GUIConfig.Seed = []byte(r)
			argsGiven = true
		}
		if r, ok := getIfIs(&ctx, "password"); ok {
			GUIConfig.Password = r
			argsGiven = true
		}
		if r, ok := getIfIs(&ctx, "publicpass"); ok {
			GUIConfig.PublicPass = r
			argsGiven = true
		}
		if argsGiven {
			dbDir := walletmain.NetworkDir(
				GUIConfig.AppDataDir, walletmain.ActiveNet.Params)
			loader := wallet.NewLoader(
				walletmain.ActiveNet.Params, dbDir, 250)
			if GUIConfig.Password == "" {
				fmt.Println("no password given")
				return 1
			}
			if GUIConfig.Seed == nil {
				seed, err := hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen)
				if err != nil {
					fmt.Println("failed to generate new seed")
					return 1
				}
				fmt.Println("Your wallet generation seed is:")
				fmt.Printf("\n%x\n\n", seed)
				fmt.Print("IMPORTANT: Keep the seed in a safe place as you will NOT be able to restore your wallet without it.\n\n")
				fmt.Print("Please keep in mind that anyone who has access to the seed can also restore your wallet thereby giving them access to all your funds, so it is imperative that you keep it in a secure location.\n\n")
				GUIConfig.Seed = []byte(seed)
			}
			w, err := loader.CreateNewWallet(
				[]byte(GUIConfig.PublicPass),
				[]byte(GUIConfig.Password),
				GUIConfig.Seed,
				time.Now())
			if err != nil {
				fmt.Println(err)
				return 1
			}
			fmt.Println("Wallet creation completed")
			fmt.Println("Seed:", string(GUIConfig.Seed))
			fmt.Println("Password: '" + string(GUIConfig.Password) + "'")
			fmt.Println("Public Password: '" + string(GUIConfig.PublicPass) + "'")
			w.Manager.Close()
			return 0

		} else {
			fmt.Println("launching GUI")
			gui.GUI()
		}
		return 0
	},
}
