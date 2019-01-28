package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"git.parallelcoin.io/pod/cmd/ctl"
	cl "git.parallelcoin.io/pod/pkg/clog"
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

		t("listcommands", "l", "list available commands"),
		t("init", "", "resets configuration to defaults"),
		t("save", "", "saves current configuration"),
		t("wallet", "w", "uses configured walletrpc instead of full node rpc"),

		f("walletrpc", "wallet RPC address to try when given wallet RPC queries"),
		f("rpcuser", "RPC username"),
		s("rpcpass", "P", "RPC password"),
		s("rpcserver", "s", "RPC server to connect to"),
		f("tls", "enable/disable (true|false)"),
		s("rpccert", "c", "RPC server certificate chain for validation"),
		f("skipverify", "do not verify tls certificates"),

		s("configfile", "C", "Path to configuration file"),
		s("debuglevel", "d", "sets logging level"),

		f("proxy", "connect via SOCKS5 proxy"),
		f("proxyuser", "username for proxy server"),
		f("proxypass", "password for proxy server"),

		f("network", "connect to (mainnet|testnet|simnet"),
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
			cl.Shutdown()
		}
		if ctx.Is("listcommands") {
			ctl.ListCommands()
		} else {
			var cfgFile string
			var ok bool
			if cfgFile, ok = ctx.Get("configfile"); !ok {
				cfgFile = ctl.DefaultConfigFile
			}
			if ctx.Is("init") {
				log <- cl.Debug{
					"writing default configuration to", cfgFile,
				}
				WriteDefaultCtlConfig(cfgFile)
				// then run from this config
				configCtl(&ctx, cfgFile)
			} else {
				log <- cl.Info{
					"loading configuration from", cfgFile,
				}
				if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
					log <- cl.Wrn("configuration file does not exist, creating new one")
					WriteDefaultCtlConfig(cfgFile)
					// then run from this config
					configCtl(&ctx, cfgFile)
				} else {
					log <- cl.Debug{"reading from", cfgFile}
					cfgData, err := ioutil.ReadFile(cfgFile)
					if err != nil {
						log <- cl.Error{err}
						cl.Shutdown()
					}
					log <- cl.Trace{"read in config file\n", string(cfgData)}
					err = json.Unmarshal(cfgData, CtlCfg)
					if err != nil {
						log <- cl.Err(err.Error())
						cl.Shutdown()
					}
					// then run from this config
					configCtl(&ctx, cfgFile)
				}
			}
		}
		log <- cl.Trace{ctx.Args}
		runCtl(ctx.Args)
		cl.Shutdown()
		return 0
	},
}

func configCtl(ctx *climax.Context, cfgFile string) {
	var r string
	var ok bool
	// Apply all configurations specified on commandline
	if r, ok = getIfIs(ctx, "debuglevel"); ok {
		CtlCfg.DebugLevel = r
		log <- cl.Trace{
			"set", "debuglevel", "to", r,
		}
	}
	if r, ok = getIfIs(ctx, "rpcuser"); ok {
		CtlCfg.RPCUser = r
		log <- cl.Tracef{
			"set %s to %s", "rpcuser", r,
		}
	}
	if r, ok = getIfIs(ctx, "rpcpass"); ok {
		CtlCfg.RPCPass = r
		log <- cl.Tracef{
			"set %s to %s", "rpcpass", r,
		}
	}
	if r, ok = getIfIs(ctx, "rpcserver"); ok {
		CtlCfg.RPCServer = r
		log <- cl.Tracef{
			"set %s to %s", "rpcserver", r,
		}
	}
	if r, ok = getIfIs(ctx, "rpccert"); ok {
		CtlCfg.RPCCert = r
		log <- cl.Tracef{"set %s to %s", "rpccert", r}
	}
	if r, ok = getIfIs(ctx, "tls"); ok {
		CtlCfg.TLS = r == "true"
		log <- cl.Tracef{"set %s to %s", "tls", r}
	}
	if r, ok = getIfIs(ctx, "proxy"); ok {
		CtlCfg.Proxy = r
		log <- cl.Tracef{"set %s to %s", "proxy", r}
	}
	if r, ok = getIfIs(ctx, "proxyuser"); ok {
		CtlCfg.ProxyUser = r
		log <- cl.Tracef{"set %s to %s", "proxyuser", r}
	}
	if r, ok = getIfIs(ctx, "proxypass"); ok {
		CtlCfg.ProxyPass = r
		log <- cl.Tracef{"set %s to %s", "proxypass", r}
	}
	otn, osn := "false", "false"
	if CtlCfg.TestNet3 {
		otn = "true"
	}
	if CtlCfg.SimNet {
		osn = "true"
	}
	tn, ts := ctx.Get("testnet")
	sn, ss := ctx.Get("simnet")
	if ts {
		CtlCfg.TestNet3 = tn == "true"
	}
	if ss {
		CtlCfg.SimNet = sn == "true"
	}
	if CtlCfg.TestNet3 && CtlCfg.SimNet {
		log <- cl.Error{
			"cannot enable simnet and testnet at the same time. current settings testnet =", otn,
			"simnet =", osn,
		}
	}
	if ctx.Is("skipverify") {
		CtlCfg.TLSSkipVerify = true
		log <- cl.Tracef{
			"set %s to true", "skipverify",
		}
	}
	if ctx.Is("wallet") {
		CtlCfg.RPCServer = CtlCfg.Wallet
		log <- cl.Trc("using configured wallet rpc server")
	}
	if r, ok = getIfIs(ctx, "walletrpc"); ok {
		CtlCfg.Wallet = r
		log <- cl.Tracef{
			"set %s to true", "walletrpc",
		}
	}
	if ctx.Is("save") {
		log <- cl.Info{
			"saving config file to",
			cfgFile,
		}
		j, err := json.MarshalIndent(CtlCfg, "", "  ")
		if err != nil {
			log <- cl.Err(err.Error())
		}
		j = append(j, '\n')
		log <- cl.Trace{
			"JSON formatted config file\n", string(j),
		}
		ioutil.WriteFile(cfgFile, j, 0600)
	}
}

// WriteCtlConfig writes the current config in the requested location
func WriteCtlConfig(cfgFile string, cc *ctl.Config) {
	j, err := json.MarshalIndent(cc, "", "  ")
	if err != nil {
		log <- cl.Err(err.Error())
	}
	j = append(j, '\n')
	log <- cl.Tracef{"JSON formatted config file\n%s", string(j)}
	err = ioutil.WriteFile(cfgFile, j, 0600)
	if err != nil {
		log <- cl.Fatal{
			"unable to write config file %s",
			err.Error(),
		}
		cl.Shutdown()
	}
}

// WriteDefaultCtlConfig writes a default config in the requested location
func WriteDefaultCtlConfig(cfgFile string) {
	defCfg := DefaultCtlConfig()
	defCfg.ConfigFile = cfgFile
	j, err := json.MarshalIndent(defCfg, "", "  ")
	if err != nil {
		log <- cl.Err(err.Error())
	}
	j = append(j, '\n')
	log <- cl.Tracef{"JSON formatted config file\n%s", string(j)}
	err = ioutil.WriteFile(cfgFile, j, 0600)
	if err != nil {
		log <- cl.Fatal{
			"unable to write config file %s",
			err.Error(),
		}
		cl.Shutdown()
	}
	// if we are writing default config we also want to use it
	CtlCfg = defCfg
}

// DefaultCtlConfig returns an allocated, default CtlCfg
func DefaultCtlConfig() *ctl.Config {
	return &ctl.Config{
		DebugLevel:    "off",
		RPCUser:       "user",
		RPCPass:       "pa55word",
		RPCServer:     ctl.DefaultRPCServer,
		RPCCert:       ctl.DefaultRPCCertFile,
		TLS:           false,
		Proxy:         "",
		ProxyUser:     "",
		ProxyPass:     "",
		TestNet3:      false,
		SimNet:        false,
		TLSSkipVerify: false,
		Wallet:        ctl.DefaultWallet,
	}
}
