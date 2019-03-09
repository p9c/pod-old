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
	cl "git.parallelcoin.io/clog"
	"git.parallelcoin.io/pod/pkg/chain/config/params"
	"git.parallelcoin.io/pod/pkg/util/hdkeychain"
	"git.parallelcoin.io/pod/pkg/wallet"
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
			return 0
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
		if _, err := os.Stat(cfgFile); os.IsNotExist(err) {

			fmt.Println("configuration file does not exist, creating new one")
			WriteDefaultWalletConfig(CreateConfig.DataDir)
		} else {
			fmt.Println("reading app configuration from", cfgFile)
			cfgData, err := ioutil.ReadFile(cfgFile)
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
		activeNet := walletmain.ActiveNet
		CreateConfig.Config.TestNet3 = false
		CreateConfig.Config.SimNet = false
		if r, ok := getIfIs(&ctx, "network"); ok {
			switch r {
			case "testnet":
				activeNet = &netparams.TestNet3Params
				CreateConfig.Config.TestNet3 = true
				CreateConfig.Config.SimNet = false
			case "simnet":
				activeNet = &netparams.SimNetParams
				CreateConfig.Config.TestNet3 = false
				CreateConfig.Config.SimNet = true
			default:
				activeNet = &netparams.MainNetParams
			}
			CreateConfig.Network = r
			argsGiven = true
		}

		if CreateConfig.Config.TestNet3 {
			activeNet = &netparams.TestNet3Params
			CreateConfig.Config.TestNet3 = true
			CreateConfig.Config.SimNet = false
		}
		if CreateConfig.Config.SimNet {
			activeNet = &netparams.SimNetParams
			CreateConfig.Config.TestNet3 = false
			CreateConfig.Config.SimNet = true
		}
		CreateConfig.Config.AppDataDir = filepath.Join(
			CreateConfig.DataDir, "wallet")
		// spew.Dump(CreateConfig)
		dbDir := walletmain.NetworkDir(
			filepath.Join(CreateConfig.DataDir, "wallet"), activeNet.Params)
		loader := wallet.NewLoader(
			walletmain.ActiveNet.Params, dbDir, 250)
		exists, err := loader.WalletExists()
		if err != nil {
			fmt.Println("ERROR", err)
			return 1
		}
		if !exists {

		} else {
			fmt.Println("\n!!! A wallet already exists at '" + dbDir + "' !!! \n")
			fmt.Println("if you are sure it isn't valuable you can delete it before running this again")
			return 1
		}
		if ctx.Is("cli") {

			e := walletmain.CreateWallet(CreateConfig.Config, activeNet)
			if e != nil {
				fmt.Println("\nerror creating wallet:", e)
			}
			fmt.Print("\nYou can now open the wallet\n")
			return 0
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

		}
		return 0
	},
}

// CreateConfig is
var CreateConfig = CreateCfg{
	DataDir: walletmain.DefaultAppDataDir,
	Network: "mainnet",
}
