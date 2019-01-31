package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	w "git.parallelcoin.io/pod/cmd/wallet"
	walletmain "git.parallelcoin.io/pod/cmd/wallet"
	cl "git.parallelcoin.io/pod/pkg/clog"
	"git.parallelcoin.io/pod/pkg/gui"
	"git.parallelcoin.io/pod/pkg/netparams"
	"git.parallelcoin.io/pod/pkg/util/hdkeychain"
	"git.parallelcoin.io/pod/pkg/wallet"
	"github.com/davecgh/go-spew/spew"
	"github.com/tucnak/climax"
)

// CreateCfg is the type for the default config data
type CreateCfg struct {
	DataDir    string
	Password   string
	PublicPass string
	Seed       []byte
	Network    string
	Config     *walletmain.Config
}

// CreateConfig is
var CreateConfig = CreateCfg{
	DataDir: walletmain.DefaultAppDataDir,
	Network: "mainnet",
}

// CreateCommand is a command to send RPC queries to bitcoin RPC protocol server for node and wallet queries
var CreateCommand = climax.Command{
	Name:  "create",
	Brief: "creates a new wallet",
	Help:  "creates a new wallet in specified data directory for a specified network",
	Flags: []climax.Flag{
		t("help", "h", "show help text"),
		s("datadir", "D", walletmain.DefaultAppDataDir, "specify where the wallet will be created"),
		s("seed", "s", "", "input pre-existing seed"),
		s("password", "p", "", "specify password for private data"),
		s("publicpass", "P", "", "specify password for public data"),
		t("cli", "c", "use commandline interface interactive input"),
		f("network", "mainnet", "connect to (mainnet|testnet|simnet)"),
	},
	Handle: func(ctx climax.Context) int {
		if ctx.Is("help") {
			fmt.Print(`Usage: create [-h] [-d] [-s] [-p] [-P] [-c] [--network]

creates a new wallet given CLI flags, or interactively
			
Available options:

	-h, --help
		show help text
	-D, --datadir="~/.pod/wallet"
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
		CreateConfig.DataDir = w.DefaultDataDir
		if r, ok := getIfIs(&ctx, "datadir"); ok {
			CreateConfig.DataDir = r
		}
		var cfgFile string
		var ok bool
		if cfgFile, ok = ctx.Get("configfile"); !ok {
			cfgFile = filepath.Join(
				filepath.Join(CreateConfig.DataDir, "wallet"),
				w.DefaultConfigFilename)
			argsGiven = true
		}
		log <- cl.Info{"loading configuration from", cfgFile}
		if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
			fmt.Println("configuration file does not exist, creating new one")
			WriteDefaultWalletConfig(CreateConfig.DataDir)
		} else {
			fmt.Println("reading app configuration from", cfgFile)
			cfgData, err := ioutil.ReadFile(cfgFile)
			fmt.Println(string(cfgData))
			if err != nil {
				fmt.Println("reading app config file", err.Error())
				WriteDefaultWalletConfig(CreateConfig.DataDir)
			}
			log <- cl.Tracef{"parsing app configuration\n%s", cfgData}
			err = json.Unmarshal(cfgData, &WalletConfig)
			if err != nil {
				fmt.Println("parsing app config file", err.Error())
				WriteDefaultWalletConfig(CreateConfig.DataDir)
			}
		}
		CreateConfig.Config = WalletConfig.Wallet
		spew.Dump(CreateConfig.Config)
		activeNet := walletmain.ActiveNet
		if CreateConfig.Config.TestNet3 {
			fmt.Println("using testnet")
			activeNet = &netparams.TestNet3Params
		}
		if CreateConfig.Config.SimNet {
			fmt.Println("using simnet")
			activeNet = &netparams.SimNetParams
		}
		if r, ok := getIfIs(&ctx, "network"); ok {
			switch r {
			case "testnet":
				activeNet = &netparams.TestNet3Params
			case "simnet":
				activeNet = &netparams.SimNetParams
			default:
				activeNet = &netparams.MainNetParams
			}
			CreateConfig.Network = r
			argsGiven = true
		}
		CreateConfig.Config.AppDataDir = filepath.Join(
			CreateConfig.DataDir, "wallet")
		fmt.Println(activeNet.Name)
		// spew.Dump(CreateConfig)
		if ctx.Is("cli") {
			walletmain.CreateWallet(CreateConfig.Config, activeNet)
			fmt.Print("\nYou can now open the wallet\n")
			os.Exit(0)
		}
		if r, ok := getIfIs(&ctx, "seed"); ok {
			CreateConfig.Seed = []byte(r)
			argsGiven = true
		}
		if r, ok := getIfIs(&ctx, "password"); ok {
			CreateConfig.Password = r
			argsGiven = true
		}
		if r, ok := getIfIs(&ctx, "publicpass"); ok {
			CreateConfig.PublicPass = r
			argsGiven = true
		}
		if argsGiven {
			dbDir := walletmain.NetworkDir(
				CreateConfig.DataDir, walletmain.ActiveNet.Params)
			loader := wallet.NewLoader(
				walletmain.ActiveNet.Params, dbDir, 250)
			if CreateConfig.Password == "" {
				fmt.Println("no password given")
				return 1
			}
			if CreateConfig.Seed == nil {
				seed, err := hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen)
				if err != nil {
					fmt.Println("failed to generate new seed")
					return 1
				}
				fmt.Println("Your wallet generation seed is:")
				fmt.Printf("\n%x\n\n", seed)
				fmt.Print("IMPORTANT: Keep the seed in a safe place as you will NOT be able to restore your wallet without it.\n\n")
				fmt.Print("Please keep in mind that anyone who has access to the seed can also restore your wallet thereby giving them access to all your funds, so it is imperative that you keep it in a secure location.\n\n")
				CreateConfig.Seed = []byte(seed)
			}
			w, err := loader.CreateNewWallet(
				[]byte(CreateConfig.PublicPass),
				[]byte(CreateConfig.Password),
				CreateConfig.Seed,
				time.Now())
			if err != nil {
				fmt.Println(err)
				return 1
			}
			fmt.Println("Wallet creation completed")
			fmt.Println("Seed:", string(CreateConfig.Seed))
			fmt.Println("Password: '" + string(CreateConfig.Password) + "'")
			fmt.Println("Public Password: '" + string(CreateConfig.PublicPass) + "'")
			w.Manager.Close()
			return 0

		} else {
			fmt.Println("launching GUI")
			gui.GUI()
		}
		return 0
	},
}