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
var Log = clog.NewSubSystem("ctl", clog.Ndbg)

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
		podutil.GenerateFlag("rpcpass", "P", "--rpcpass=password", "RPC password", true),
		podutil.GenerateFlag("rpcserver", "s", "--rpcserver=127.0.0.1:11048", "RPC server to connect to", true),
		podutil.GenerateFlag("rpccert", "c", "--rpccert=/path/to/rpc.cert", "RPC server certificate chain for validation", true),
		podutil.GenerateFlag("tls", "", "--tls=false", "Enable/disable TLS", false),
		podutil.GenerateFlag("proxy", "", "--proxy 127.0.0.1:9050", "Connect via SOCKS5 proxy (eg. 127.0.0.1:9050)", true),
		podutil.GenerateFlag("proxyuser", "", "--proxyuser username", "Username for proxy server", true),
		podutil.GenerateFlag("proxypass", "", "--proxypass password", "Password for proxy server", true),
		podutil.GenerateFlag("testnet", "", "--testnet=true", "Connect to testnet", true),
		podutil.GenerateFlag("simnet", "", "--simnet=true", "Connect to the simulation test network", true),
		podutil.GenerateFlag("skipverify", "", "--skipverify=false", "Do not verify tls certificates (not recommended!)", true),
		podutil.GenerateFlag("wallet", "", "--wallet=true", "Connect to wallet", true),
	},
	Examples: []climax.Example{
		{
			Usecase:     "-l",
			Description: "lists available commands",
		},
	},
	Handle: func(ctx climax.Context) int {
		if dl, ok := ctx.Get("debuglevel"); ok {
			Log.Tracef.Print("setting debug level %s", dl)
			Log.SetLevel(dl)
		}
		Log.Debugf.Print("pod/ctl version %s", c.Version())
		if ctx.Is("version") {
			fmt.Println("pod/ctl version", c.Version())
			clog.Shutdown()
		}
		if ctx.Is("listcommands") {
			Log.Trace.Print("listing commands")
			c.ListCommands()
		} else {
			Log.Trace.Print("running command")

			var cfgFile string
			var ok bool
			if cfgFile, ok = ctx.Get("configfile"); !ok {
				cfgFile = c.DefaultConfigFile
			}
			if ctx.Is("init") {
				Log.Debugf.Print("writing default configuration to %s", cfgFile)
				writeDefaultConfig(cfgFile)
				// then run from this config
				configCtl(&ctx, cfgFile)
			} else {
				Log.Infof.Print("loading configuration from %s", cfgFile)
				if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
					Log.Warn.Print("configuration file does not exist, creating new one")
					writeDefaultConfig(cfgFile)
					// then run from this config
					configCtl(&ctx, cfgFile)
				} else {
					Log.Debug.Print("reading from", cfgFile)
					cfgData, err := ioutil.ReadFile(cfgFile)
					if err != nil {
						Log.Error.Print(err.Error())
						clog.Shutdown()
					}
					Log.Tracef.Print("read in config file\n%s", cfgData)
					err = json.Unmarshal(cfgData, Config)
					if err != nil {
						Log.Error.Print(err.Error())
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
		Log.Tracef.Print("set %s to %s", "rpcuser", r)
	}
	if ctx.Is("rpcpass") {
		r, _ := ctx.Get("rpcpass")
		Config.RPCPassword = r
		Log.Tracef.Print("set %s to %s", "rpcpass", r)
	}
	if ctx.Is("rpcserver") {
		r, _ := ctx.Get("rpcserver")
		Config.RPCServer = r
		Log.Tracef.Print("set %s to %s", "rpcserver", r)
	}
	if ctx.Is("rpccert") {
		r, _ := ctx.Get("rpccert")
		Config.RPCCert = r
		Log.Tracef.Print("set %s to %s", "rpccert", r)
	}
	if ctx.Is("tls") {
		r, _ := ctx.Get("tls")
		Config.TLS = r == "true"
		Log.Tracef.Print("set %s to %s", "tls", r)
	}
	if ctx.Is("proxy") {
		r, _ := ctx.Get("proxy")
		Config.Proxy = r
		Log.Tracef.Print("set %s to %s", "proxy", r)
	}
	if ctx.Is("proxyuser") {
		r, _ := ctx.Get("proxyuser")
		Config.ProxyUser = r
		Log.Tracef.Print("set %s to %s", "proxyuser", r)
	}
	if ctx.Is("proxypass") {
		r, _ := ctx.Get("proxypass")
		Config.ProxyPass = r
		Log.Tracef.Print("set %s to %s", "proxypass", r)
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
		Log.Error.Print("cannot enable simnet and testnet at the same time. current settings testnet =", otn, "simnet =", osn)
	}
	if ctx.Is("skipverify") {
		Config.TLSSkipVerify = true
		Log.Tracef.Print("set %s to true", "skipverify")
	}
	if ctx.Is("wallet") {
		Config.Wallet = true
		Log.Tracef.Print("set %s to true", "wallet")
	}
	if ctx.Is("save") {
		Log.Infof.Print("saving config file to %s", cfgFile)
		j, err := json.MarshalIndent(Config, "", "  ")
		if err != nil {
			Log.Error.Print(err.Error())
		}
		j = append(j, '\n')
		Log.Tracef.Print("JSON formatted config file\n%s", j)
		ioutil.WriteFile(cfgFile, j, 0600)
	}
}

func writeDefaultConfig(cfgFile string) {
	defCfg := defaultConfig()
	defCfg.ConfigFile = cfgFile
	j, err := json.MarshalIndent(defCfg, "", "  ")
	if err != nil {
		Log.Error.Print(err.Error())
	}
	j = append(j, '\n')
	Log.Tracef.Print("JSON formatted config file\n%s", j)
	err = ioutil.WriteFile(cfgFile, j, 0600)
	if err != nil {
		Log.Fatalf.Print("unable to write config file %s", err.Error())
		clog.Shutdown()
	}
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
