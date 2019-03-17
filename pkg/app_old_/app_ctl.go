package app_old

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"git.parallelcoin.io/dev/pod/cmd/ctl"
	w "git.parallelcoin.io/dev/pod/cmd/wallet"
	cl "git.parallelcoin.io/dev/pod/pkg/util/cl"
	"github.com/tucnak/climax"
)

// CtlCfg is the default configuration native to ctl
var CtlCfg = new(ctl.Config)

// CtlCommand is a command to send RPC queries to bitcoin RPC protocol server for node and wallet queries
var CtlCommand = climax.Command{
	Name:  "ctl",
	Brief: "sends RPC commands and prints the reply",
	Help:  "Send queries to bitcoin JSON-RPC servers using command line shell and prints the reply to stdout",
	Flags: []climax.Flag{
		t("version", "V", "show version number and quit"),
		s("configfile", "C", ctl.DefaultConfigFile, "Path to configuration file"),
		s("datadir", "D", w.DefaultDataDir,
			"set the pod base directory"),

		t("init", "", "resets configuration to defaults"),
		t("save", "", "saves current configuration"),
		t("wallet", "w", "uses configured walletrpc instead of full node rpc"),

		f("walletrpc", ctl.DefaultRPCServer,
			"wallet RPC address to try when given wallet RPC queries"),
		s("rpcuser", "u", "user", "RPC username"),
		s("rpcpass", "P", "pa55word", "RPC password"),
		s("rpcserver", "s", "127.0.0.1:11048", "RPC server to connect to"),
		f("tls", "false", "enable/disable (true|false)"),
		s("rpccert", "c", "rpc.cert", "RPC server certificate chain for validation"),
		f("skipverify", "false", "do not verify tls certificates"),

		s("debuglevel", "d", "off", "sets logging level (off|fatal|error|info|debug|trace)"),

		f("proxy", "", "connect via SOCKS5 proxy"),
		f("proxyuser", "user", "username for proxy server"),
		f("proxypass", "pa55word", "password for proxy server"),

		f("network", "mainnet", "connect to (mainnet|testnet|simnet)"),
	},

	Examples: []climax.Example{
		{
			Usecase:     "-l",
			Description: "lists available commands",
		},
	},

	Handle: func(ctx climax.Context) int {
		Log.SetLevel("off")

		if dl, ok := ctx.Get("debuglevel"); ok {

			log <- cl.Trace{

				"setting debug level", dl,
			}

			Log.SetLevel(dl)
		}

		log <- cl.Debug{

			"pod/ctl version", ctl.Version(),
		}

		if ctx.Is("version") {

			fmt.Println("pod/ctl version", ctl.Version())
			return 0
		}

		if ctx.Is("listcommands") {

			ctl.ListCommands()

		} else {
			var cfgFile, datadir string
			var ok bool

			if cfgFile, ok = ctx.Get("configfile"); !ok {

				cfgFile = ctl.DefaultConfigFile
			}

			if datadir, ok = ctx.Get("datadir"); ok {
				cfgFile = filepath.Join(filepath.Join(datadir, "ctl"), "conf.json")
				CtlCfg.ConfigFile = cfgFile

			} else {
				datadir = w.DefaultDataDir
			}

			if ctx.Is("init") {

				log <- cl.Debug{

					"writing default configuration to", cfgFile,
				}

				WriteDefaultCtlConfig(datadir)

			} else {
				log <- cl.Info{

					"loading configuration from", cfgFile,
				}

				if _, err := os.Stat(cfgFile); os.IsNotExist(err) {

					log <- cl.Wrn("configuration file does not exist, creating new one")

					WriteDefaultCtlConfig(datadir)
					// then run from this config
					configCtl(&ctx, cfgFile)

				} else {
					log <- cl.Debug{"reading from", cfgFile}

					cfgData, err := ioutil.ReadFile(cfgFile)

					if err != nil {

						WriteDefaultCtlConfig(datadir)
						log <- cl.Error{err}

					}

					log <- cl.Trace{"read in config file\n", string(cfgData)}

					err = json.Unmarshal(cfgData, CtlCfg)

					if err != nil {

						log <- cl.Err(err.Error())

						return 1
					}

				}

				// then run from this config
				configCtl(&ctx, cfgFile)
			}

		}

		log <- cl.Trace{ctx.Args}

		runCtl(ctx.Args, CtlCfg)
		return 0
	},
}
