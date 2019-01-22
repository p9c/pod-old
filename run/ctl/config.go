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

// Command is a command to send RPC queries to bitcoin RPC protocol server for node and wallet queries
var Command = climax.Command{
	Name:  "ctl",
	Brief: "sends RPC commands and prints the reply",
	Help:  "Send queries to bitcoin JSON-RPC servers using command line shell and prints the reply to stdout",
	Flags: []climax.Flag{
		podutil.GenerateFlag("version", "V", `--version`, `show version number and quit`, false),
		podutil.GenerateFlag("configfile", "C", "--configfile=/path/to/conf", "Path to configuration file", true),

		podutil.GenerateFlag("listcommands", "l", `--listcommands`, `list available commands`, false),
		podutil.GenerateFlag("init", "", "--init", "resets configuration to defaults", false),
		podutil.GenerateFlag("save", "", "--save", "saves current configuration", false),

		podutil.GenerateFlag("debuglevel", "d", "--debuglevel=trace", "sets debuglevel, default is error to keep stdout clean", true),

		podutil.GenerateFlag("rpcuser", "u", "--rpcuser=username", "RPC username", true),
		podutil.GenerateFlag("rpcpass", "P", "--rpcpass=pa55word", "RPC password", true),
		podutil.GenerateFlag("rpcserver", "s", "--rpcserver=127.0.0.1:11048", "RPC server to connect to", true),
		podutil.GenerateFlag("rpccert", "c", "--rpccert=/path/to/rpc.cert", "RPC server certificate chain for validation", true),
		podutil.GenerateFlag("tls", "", "--tls=false", "enable/disable TLS", false),
		podutil.GenerateFlag("proxy", "", "--proxy 127.0.0.1:9050", "connect via SOCKS5 proxy (eg. 127.0.0.1:9050)", true),
		podutil.GenerateFlag("proxyuser", "", "--proxyuser=username", "username for proxy server", true),
		podutil.GenerateFlag("proxypass", "", "--proxypass=password", "password for proxy server", true),
		podutil.GenerateFlag("testnet", "", "--testnet=true", "connect to testnet", true),
		podutil.GenerateFlag("simnet", "", "--simnet=true", "connect to the simulation test network", true),
		podutil.GenerateFlag("skipverify", "", "--skipverify=false", "do not verify tls certificates (not recommended!)", true),
		podutil.GenerateFlag("wallet", "", "--wallet=true", "connect to wallet", true),
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
				writeDefaultConfig(cfgFile)
				// then run from this config
				configCtl(&ctx, cfgFile)
			} else {
				log <- cl.Info{
					"loading configuration from", cfgFile,
				}
				if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
					log <- cl.Wrn("configuration file does not exist, creating new one")
					writeDefaultConfig(cfgFile)
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
	if ctx.Is("wallet") {
		Config.Wallet = true
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

func writeDefaultConfig(cfgFile string) {
	defCfg := defaultConfig()
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

func defaultConfig() *c.Config {
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
		Wallet:        false,
	}
}
