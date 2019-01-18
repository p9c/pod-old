package ctl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"git.parallelcoin.io/pod/lib/clog"
	c "git.parallelcoin.io/pod/module/ctl"
	"github.com/tucnak/climax"
)

var log = clog.NewSubSystem("Ctl", clog.Nwrn)

// Config is the default configuration native to ctl
var Config = new(c.Config)

// Command is a command to send RPC queries to bitcoin RPC protocol server for node and wallet queries
var Command = climax.Command{
	Name:  "ctl",
	Brief: "sends RPC commands and prints the reply",
	Help:  "Send queries to bitcoin JSON-RPC servers using command line shell and prints the reply to stdout",
	Flags: []climax.Flag{
		{
			Name:     "listcommands",
			Short:    "l",
			Usage:    `--listcommands`,
			Help:     `list available commands`,
			Variable: false,
		},
		{
			Name:     "version",
			Short:    "V",
			Usage:    `--version`,
			Help:     `show version number and quit`,
			Variable: false,
		},
		{
			Name:     "configfile",
			Short:    "C",
			Usage:    "--configfile=/path/to/conf",
			Help:     "Path to configuration file",
			Variable: true,
		},
		{
			Name:     "init",
			Usage:    "--init",
			Help:     "resets configuration to defaults",
			Variable: false,
		},
		{
			Name:     "save",
			Usage:    "--save",
			Help:     "saves current configuration",
			Variable: false,
		},
		{
			Name:     "debuglevel",
			Short:    "d",
			Usage:    "--debuglevel=trace",
			Help:     "sets debuglevel, default is error to keep stdout clean",
			Variable: true,
		},

		{
			Name:     "rpcuser",
			Short:    "u",
			Usage:    "--rpcuser=username",
			Help:     "RPC username",
			Variable: true,
		},
		{
			Name:     "rpcpass",
			Short:    "P",
			Usage:    "--rpcpass=password",
			Help:     "RPC password",
			Variable: true,
		},
		{
			Name:     "rpcserver",
			Short:    "s",
			Usage:    "--rpcserver=127.0.0.1:11048",
			Help:     "RPC server to connect to",
			Variable: true,
		},
		{
			Name:     "rpccert",
			Short:    "c",
			Usage:    "--rpccert=/path/to/rpc.cert",
			Help:     "RPC server certificate chain for validation",
			Variable: true,
		},
		{
			Name:     "tls",
			Usage:    "--tls=false",
			Help:     "Enable/disable TLS",
			Variable: true,
		},
		{
			Name:     "proxy",
			Usage:    "--proxy 127.0.0.1:9050",
			Help:     "Connect via SOCKS5 proxy (eg. 127.0.0.1:9050)",
			Variable: true,
		},
		{
			Name:     "proxyuser",
			Usage:    "--proxyuser username",
			Help:     "Username for proxy server",
			Variable: true,
		},
		{
			Name:     "proxypass",
			Usage:    "--proxypass password",
			Help:     "Password for proxy server",
			Variable: true,
		},
		{
			Name:     "testnet",
			Usage:    "--testnet=true",
			Help:     "Connect to testnet",
			Variable: true,
		},
		{
			Name:     "simnet",
			Usage:    "--simnet=true",
			Help:     "Connect to the simulation test network",
			Variable: true,
		},
		{
			Name:     "skipverify",
			Usage:    "--skipverify=false",
			Help:     "Do not verify tls certificates (not recommended!)",
			Variable: true,
		},
		{
			Name:     "wallet",
			Usage:    "--wallet=true",
			Help:     "Connect to wallet",
			Variable: true,
		},
	},
	Examples: []climax.Example{
		{
			Usecase:     "-l",
			Description: "lists available commands",
		},
	},
	Handle: func(ctx climax.Context) int {
		if dl, ok := ctx.Get("debuglevel"); ok {
			log.Tracef.Print("setting debug level %s", dl)
			log.SetLevel(dl)
		}
		log.Debugf.Print("ctl version %s", c.Version())
		if ctx.Is("version") {
			fmt.Println("pod version", c.Version())
			clog.Shutdown()
		}
		if ctx.Is("listcommands") {
			log.Trace.Print("listing commands")
			c.ListCommands()
		} else {
			log.Trace.Print("running command")

			var cfgFile string
			var ok bool
			if cfgFile, ok = ctx.Get("configfile"); !ok {
				cfgFile = c.DefaultConfigFile
			}
			if ctx.Is("init") {
				log.Debugf.Print("writing default configuration to %s", cfgFile)
				writeDefaultConfig(cfgFile)
				// then run from this config
				configCtl(&ctx, cfgFile)
			} else {
				log.Infof.Print("loading configuration from %s", cfgFile)
				if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
					log.Warn.Print("configuration file does not exist, creating new one")
					writeDefaultConfig(cfgFile)
					// then run from this config
					configCtl(&ctx, cfgFile)
				} else {
					log.Debug.Print("reading from", cfgFile)
					cfgData, err := ioutil.ReadFile(cfgFile)
					if err != nil {
						log.Error.Print(err.Error())
						clog.Shutdown()
					}
					log.Tracef.Print("read in config file\n%s", cfgData)
					err = json.Unmarshal(cfgData, Config)
					if err != nil {
						log.Error.Print(err.Error())
						clog.Shutdown()
					}
					// then run from this config
					configCtl(&ctx, cfgFile)
				}
			}
		}
		runCtl()
		clog.Shutdown()
		return 0
	},
}

func configCtl(ctx *climax.Context, cfgFile string) {
	// Apply all configurations specified on commandline
	if ctx.Is("rpcuser") {
		r, _ := ctx.Get("rpcuser")
		Config.RPCUser = r
		log.Tracef.Print("set %s to %s", "rpcuser", r)
	}
	if ctx.Is("rpcpass") {
		r, _ := ctx.Get("rpcpass")
		Config.RPCPassword = r
		log.Tracef.Print("set %s to %s", "rpcpass", r)
	}
	if ctx.Is("rpcserver") {
		r, _ := ctx.Get("rpcserver")
		Config.RPCServer = r
		log.Tracef.Print("set %s to %s", "rpcserver", r)
	}
	if ctx.Is("rpccert") {
		r, _ := ctx.Get("rpccert")
		Config.RPCCert = r
		log.Tracef.Print("set %s to %s", "rpccert", r)
	}
	if ctx.Is("tls") {
		r, _ := ctx.Get("tls")
		Config.TLS = r == "true"
		log.Tracef.Print("set %s to %s", "tls", r)
	}
	if ctx.Is("proxy") {
		r, _ := ctx.Get("proxy")
		Config.Proxy = r
		log.Tracef.Print("set %s to %s", "proxy", r)
	}
	if ctx.Is("proxyuser") {
		r, _ := ctx.Get("proxyuser")
		Config.ProxyUser = r
		log.Tracef.Print("set %s to %s", "proxyuser", r)
	}
	if ctx.Is("proxypass") {
		r, _ := ctx.Get("proxypass")
		Config.ProxyPass = r
		log.Tracef.Print("set %s to %s", "proxypass", r)
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
		log.Error.Print("cannot enable simnet and testnet at the same time. current settings testnet =", otn, "simnet =", osn)
	}
	if ctx.Is("skipverify") {
		Config.TLSSkipVerify = true
		log.Tracef.Print("set %s to true", "skipverify")
	}
	if ctx.Is("wallet") {
		Config.Wallet = true
		log.Tracef.Print("set %s to true", "wallet")
	}
	if ctx.Is("save") {
		log.Infof.Print("saving config file to %s", cfgFile)
		j, err := json.MarshalIndent(Config, "", "  ")
		if err != nil {
			log.Error.Print(err.Error())
		}
		j = append(j, '\n')
		log.Tracef.Print("JSON formatted config file\n%s", j)
		ioutil.WriteFile(cfgFile, j, 0600)
	}
}

func writeDefaultConfig(cfgFile string) {
	defCfg := defaultConfig()
	defCfg.ConfigFile = cfgFile
	j, err := json.MarshalIndent(defCfg, "", "  ")
	if err != nil {
		log.Error.Print(err.Error())
	}
	j = append(j, '\n')
	log.Tracef.Print("JSON formatted config file\n%s", j)
	ioutil.WriteFile(cfgFile, j, 0600)
	// if we are writing default config we also want to use it
	Config = defCfg
}

func defaultConfig() *c.Config {
	return &c.Config{
		RPCUser:       "user",
		RPCPassword:   "pa55word",
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
