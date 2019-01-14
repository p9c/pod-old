package main

import (
	"fmt"
	"os"
	"time"
)

type walletCfgLaunchGroup struct {
	ConfigFile    string `short:"C" long:"configfile" description:"path to configuration file"`
	ShowVersion   bool   `short:"V" long:"version" description:"Display version information and exit"`
	Create        bool   `long:"create" description:"Create the wallet if it does not exist"`
	CreateTemp    bool   `long:"createtemp" description:"Create a temporary simulation wallet (pass=password) in the data directory indicated; must call with --datadir"`
	NoInitialLoad bool   `long:"noinitialload" description:"Defer wallet creation/opening on startup and enable loading wallets over RPC"`
	LogDir        string `long:"logdir" description:"Directory to log output."`
	Profile       string `long:"profile" description:"Enable HTTP profiling on given port -- NOTE port must be between 1024 and 65536"`
}

type walletCfgChainGroup struct {
	TestNet3       bool `long:"testnet" description:"Connect to testnet"`
	SimNet         bool `long:"simnet" description:"Connect to the simulation test network"`
	RegressionTest bool `long:"regtest" description:"Connect to the regression test network"`
}

type walletNodeCfg struct {
	RPCConnect      string        `short:"c" long:"rpcconnect" description:"Hostname/IP and port of pod RPC server to connect to (default localhost:11048, testnet: localhost:21048, simnet: localhost:41048)"`
	CAFile          string        `long:"cafile" description:"File containing root certificates to authenticate a TLS connections with pod"`
	EnableClientTLS bool          `long:"clienttls" description:"Enable TLS for the RPC client"`
	PodUsername     string        `long:"podusername" description:"Username for pod authentication"`
	PodPassword     string        `long:"podpassword" default-mask:"-" description:"Password for pod authentication"`
	Proxy           string        `long:"proxy" description:"Connect via SOCKS5 proxy (eg. 127.0.0.1:9050)"`
	ProxyUser       string        `long:"proxyuser" description:"Username for proxy server"`
	ProxyPass       string        `long:"proxypass" default-mask:"-" description:"Password for proxy server"`
	AddPeers        []string      `short:"a" long:"addpeer" description:"Add a peer to connect with at startup"`
	ConnectPeers    []string      `long:"connect" description:"Connect only to the specified peers at startup"`
	MaxPeers        int           `long:"maxpeers" description:"Max number of inbound and outbound peers"`
	BanDuration     time.Duration `long:"banduration" description:"How long to ban misbehaving peers.  Valid time units are {s, m, h}.  Minimum 1 second"`
	BanThreshold    uint32        `long:"banthreshold" description:"Maximum allowed ban score before disconnecting and banning misbehaving peers."`
}

type walletRPCCfgGroup struct {
	RPCCert                  string   `long:"rpccert" description:"File containing the certificate file"`
	RPCKey                   string   `long:"rpckey" description:"File containing the certificate key"`
	OneTimeTLSKey            bool     `long:"onetimetlskey" description:"Generate a new TLS certpair at startup, but only write the certificate to disk"`
	EnableServerTLS          bool     `long:"servertls" description:"Enable TLS for the RPC server"`
	LegacyRPCListeners       []string `long:"rpclisten" description:"Listen for legacy RPC connections on this interface/port (default port: 11046, testnet: 21046, simnet: 41046)"`
	LegacyRPCMaxClients      int64    `long:"rpcmaxclients" description:"Max number of legacy RPC clients for standard connections"`
	LegacyRPCMaxWebsockets   int64    `long:"rpcmaxwebsockets" description:"Max number of legacy RPC websocket connections"`
	Username                 string   `short:"u" long:"username" description:"Username for legacy RPC and pod authentication (if podusername is unset)"`
	Password                 string   `short:"P" long:"password" default-mask:"-" description:"Password for legacy RPC and pod authentication (if podpassword is unset)"`
	ExperimentalRPCListeners []string `long:"experimentalrpclisten" description:"Listen for RPC connections on this interface/port"`
}

type walletCfg struct {
	LaunchGroup       walletCfgLaunchGroup `group:"Launch options"`
	NodeChainGroup    walletCfgChainGroup  `group:"Chain options"`
	NodeCfgGroup      walletNodeCfg        `group:"Node connection options"`
	WalletRPCCfgGroup walletRPCCfgGroup    `group:"Wallet RPC configuration"`
}

var wallet walletCfg

func (n *walletCfg) Execute(args []string) (err error) {
	fmt.Println("running wallet")
	fmt.Println("not implemented - quitting")
	os.Exit(1)
	return
}
