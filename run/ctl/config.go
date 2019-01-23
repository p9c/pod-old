package ctl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"git.parallelcoin.io/pod/lib/clog"
	c "git.parallelcoin.io/pod/module/ctl"
	"git.parallelcoin.io/pod/run/util"
	"github.com/tucnak/climax"
)

// Log is the ctl main logger
var Log = cl.NewSubSystem("run/ctl", "off")
var log = Log.Ch

// Config is the default configuration native to ctl
var Config = new(c.Config)

var f = pu.GenFlag
var t = pu.GenTrig
var s = pu.GenShort
var l = pu.GenLog

// Command is a command to send RPC queries to bitcoin RPC protocol server for node and wallet queries
var Command = climax.Command{
	Name:  "ctl",
	Brief: "sends RPC commands and prints the reply",
	Help:  "Send queries to bitcoin JSON-RPC servers using command line shell and prints the reply to stdout",
	Flags: []climax.Flag{
		t("version", "V", "show version number and quit"),

		t("listcommands", "l", "list available commands"),
		t("init", "", "resets configuration to defaults"),
		t("save", "", "saves current configuration"),

		f("wallet", "wallet RPC address to try when given wallet RPC queries"),
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
		if dl, ok := ctx.Get("debuglevel"); ok {
			log <- cl.Trace{
				"setting debug level", dl,
			}
			Log.SetLevel(dl)
		}
		log <- cl.Debug{
			"pod/ctl version", c.Version(),
		}
		if ctx.Is("version") {
			fmt.Println("pod/ctl version", c.Version())
			cl.Shutdown()
		}
		if ctx.Is("listcommands") {
			c.ListCommands()
		} else {
			var cfgFile string
			var ok bool
			if cfgFile, ok = ctx.Get("configfile"); !ok {
				cfgFile = c.DefaultConfigFile
			}
			if ctx.Is("init") {
				log <- cl.Debug{
					"writing default configuration to", cfgFile,
				}
				WriteDefaultConfig(cfgFile)
				// then run from this config
				configCtl(&ctx, cfgFile)
			} else {
				log <- cl.Info{
					"loading configuration from", cfgFile,
				}
				if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
					log <- cl.Wrn("configuration file does not exist, creating new one")
					WriteDefaultConfig(cfgFile)
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
					err = json.Unmarshal(cfgData, Config)
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

func getIfIs(ctx *climax.Context, name string) (out string, ok bool) {
	if ctx.Is(name) {
		return ctx.Get(name)
	}
	return
}

func configCtl(ctx *climax.Context, cfgFile string) {
	var r string
	var ok bool
	// Apply all configurations specified on commandline
	if r, ok = getIfIs(ctx, "debuglevel"); ok {
		Config.DebugLevel = r
		log <- cl.Trace{
			"set", "debuglevel", "to", r,
		}
	}
	if r, ok = getIfIs(ctx, "rpcuser"); ok {
		Config.RPCUser = r
		log <- cl.Tracef{
			"set %s to %s", "rpcuser", r,
		}
	}
	if r, ok = getIfIs(ctx, "rpcpass"); ok {
		Config.RPCPass = r
		log <- cl.Tracef{
			"set %s to %s", "rpcpass", r,
		}
	}
	if r, ok = getIfIs(ctx, "rpcserver"); ok {
		Config.RPCServer = r
		log <- cl.Tracef{
			"set %s to %s", "rpcserver", r,
		}
	}
	if r, ok = getIfIs(ctx, "rpccert"); ok {
		Config.RPCCert = r
		log <- cl.Tracef{"set %s to %s", "rpccert", r}
	}
	if r, ok = getIfIs(ctx, "tls"); ok {
		Config.TLS = r == "true"
		log <- cl.Tracef{"set %s to %s", "tls", r}
	}
	if r, ok = getIfIs(ctx, "proxy"); ok {
		Config.Proxy = r
		log <- cl.Tracef{"set %s to %s", "proxy", r}
	}
	if r, ok = getIfIs(ctx, "proxyuser"); ok {
		Config.ProxyUser = r
		log <- cl.Tracef{"set %s to %s", "proxyuser", r}
	}
	if r, ok = getIfIs(ctx, "proxypass"); ok {
		Config.ProxyPass = r
		log <- cl.Tracef{"set %s to %s", "proxypass", r}
	}
	otn, osn := "false", "false"
	if Config.TestNet3 {
		otn = "true"
	}
	if Config.SimNet {
		osn = "true"
	}
	tn, ts := ctx.Get("testnet")
	sn, ss := ctx.Get("simnet")
	if ts {
		Config.TestNet3 = tn == "true"
	}
	if ss {
		Config.SimNet = sn == "true"
	}
	if Config.TestNet3 && Config.SimNet {
		log <- cl.Error{
			"cannot enable simnet and testnet at the same time. current settings testnet =", otn,
			"simnet =", osn,
		}
	}
	if ctx.Is("skipverify") {
		Config.TLSSkipVerify = true
		log <- cl.Tracef{
			"set %s to true", "skipverify",
		}
	}
	if r, ok = getIfIs(ctx, "wallet"); ok {
		Config.Wallet = r
		log <- cl.Tracef{
			"set %s to true", "wallet",
		}
	}
	if ctx.Is("save") {
		log <- cl.Info{
			"saving config file to",
			cfgFile,
		}
		j, err := json.MarshalIndent(Config, "", "  ")
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

// WriteConfig writes the current config in the requested location
func WriteConfig(cfgFile string, cc *c.Config) {
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

// WriteDefaultConfig writes a default config in the requested location
func WriteDefaultConfig(cfgFile string) {
	defCfg := DefaultConfig()
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
	Config = defCfg
}

func DefaultConfig() *c.Config {
	return &c.Config{
		DebugLevel:    "off",
		RPCUser:       "user",
		RPCPass:       "pa55word",
		RPCServer:     c.DefaultRPCServer,
		RPCCert:       c.DefaultRPCCertFile,
		TLS:           false,
		Proxy:         "",
		ProxyUser:     "",
		ProxyPass:     "",
		TestNet3:      false,
		SimNet:        false,
		TLSSkipVerify: false,
		Wallet:        "",
	}
}
