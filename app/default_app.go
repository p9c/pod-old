package app

import (
	"fmt"
	"os"
	"time"

	walletmain "git.parallelcoin.io/pod/cmd/wallet"
	"git.parallelcoin.io/pod/pkg/netparams"
	"git.parallelcoin.io/pod/pkg/util/hdkeychain"
	"git.parallelcoin.io/pod/pkg/wallet"
	"github.com/tucnak/climax"
)

// DefaultCfg is the type for the default config data
type DefaultCfg struct {
	AppDataDir string
	Password   string
	PublicPass string
	Seed       []byte
	Network    string
}

// DefaultConfig is
var DefaultConfig = DefaultCfg{
	AppDataDir: walletmain.DefaultAppDataDir,
	Network:    "mainnet",
}

// DefaultCommand is a command to send RPC queries to bitcoin RPC protocol server for node and wallet queries
var DefaultCommand = climax.Command{
	Name:  "createwallet",
	Brief: "creates a new wallet",
	Help:  "creates a new wallet using GUI or CLI flags, if no parameters are given, launches the GUI input, if wallet exists, launch the shell with GUI",
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

creates a new wallet using GUI or CLI flags, if no parameters are given, launches the GUI input, if wallet exists, launch the shell with GUI
			
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
			DefaultConfig.AppDataDir = r
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
			DefaultConfig.Network = r
			argsGiven = true
		}
		if ctx.Is("cli") {
			walletmain.CreateWallet(&walletmain.Config{
				AppDataDir: DefaultConfig.AppDataDir,
				WalletPass: DefaultConfig.PublicPass,
			})
			fmt.Print("\nYou can now open the wallet\n")
			os.Exit(0)
		}
		if r, ok := getIfIs(&ctx, "seed"); ok {
			DefaultConfig.Seed = []byte(r)
			argsGiven = true
		}
		if r, ok := getIfIs(&ctx, "password"); ok {
			DefaultConfig.Password = r
			argsGiven = true
		}
		if r, ok := getIfIs(&ctx, "publicpass"); ok {
			DefaultConfig.PublicPass = r
			argsGiven = true
		}
		if argsGiven {
			dbDir := walletmain.NetworkDir(
				DefaultConfig.AppDataDir, walletmain.ActiveNet.Params)
			loader := wallet.NewLoader(
				walletmain.ActiveNet.Params, dbDir, 250)
			if DefaultConfig.Password == "" {
				fmt.Println("no password given")
				return 1
			}
			if DefaultConfig.Seed == nil {
				seed, err := hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen)
				if err != nil {
					fmt.Println("failed to generate new seed")
					return 1
				}
				fmt.Println("Your wallet generation seed is:")
				fmt.Printf("\n%x\n\n", seed)
				fmt.Print("IMPORTANT: Keep the seed in a safe place as you will NOT be able to restore your wallet without it.\n\n")
				fmt.Print("Please keep in mind that anyone who has access to the seed can also restore your wallet thereby giving them access to all your funds, so it is imperative that you keep it in a secure location.\n\n")
				DefaultConfig.Seed = []byte(seed)
			}
			w, err := loader.CreateNewWallet(
				[]byte(DefaultConfig.PublicPass),
				[]byte(DefaultConfig.Password),
				DefaultConfig.Seed,
				time.Now())
			if err != nil {
				fmt.Println(err)
				return 1
			}
			fmt.Println("Wallet creation completed")
			fmt.Println("Seed:", string(DefaultConfig.Seed))
			fmt.Println("Password: '" + string(DefaultConfig.Password) + "'")
			fmt.Println("Public Password: '" + string(DefaultConfig.PublicPass) + "'")
			w.Manager.Close()
			return 0

		} else {
			fmt.Println("launching GUI")
			// Start GUI actually!
		}
		return 0
	},
}
