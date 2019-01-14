package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type ctlCfgLaunchGroup struct {
	ShowVersion  bool   `short:"V" long:"version" description:"Display version information and exit"`
	ListCommands bool   `short:"l" long:"listcommands" description:"List all of the supported commands and exit"`
	ConfigFile   string `short:"C" long:"configfile" description:"Path to configuration file"`
	Wallet       bool   `long:"wallet" description:"Connect to wallet"`
}

type ctlCfgRPCGroup struct {
	RPCUser       string `short:"u" long:"rpcuser" description:"RPC username"`
	RPCPassword   string `short:"P" long:"rpcpass" default-mask:"-" description:"RPC password"`
	RPCServer     string `short:"s" long:"rpcserver" description:"RPC server to connect to"`
	RPCCert       string `short:"c" long:"rpccert" description:"RPC server certificate chain for validation"`
	TLS           bool   `long:"tls" description:"Enable TLS"`
	Proxy         string `long:"proxy" description:"Connect via SOCKS5 proxy (eg. 127.0.0.1:9050)"`
	ProxyUser     string `long:"proxyuser" description:"Username for proxy server"`
	ProxyPass     string `long:"proxypass" default-mask:"-" description:"Password for proxy server"`
	TLSSkipVerify bool   `long:"skipverify" description:"Do not verify tls certificates (not recommended!)"`
}

type ctlCfgChainGroup struct {
	TestNet3       bool `long:"testnet" description:"Connect to testnet"`
	SimNet         bool `long:"simnet" description:"Connect to the simulation test network"`
	RegressionTest bool `long:"regtest" description:"Connect to the regression test network"`
}

type ctlCfg struct {
	LaunchOptions ctlCfgLaunchGroup `group:"Launch options"`
	RPCOptions    ctlCfgRPCGroup    `group:"RPC options"`
	ChainOptions  ctlCfgChainGroup  `group:"Chain options"`
}

var ctl ctlCfg

func (n *ctlCfg) Execute(args []string) (err error) {
	fmt.Println("running ctl")
	j, _ := json.MarshalIndent(n, "", "\t")
	fmt.Println(string(j))
	fmt.Println("not implemented - quitting")
	os.Exit(1)
	return
}
