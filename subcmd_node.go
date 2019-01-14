package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	"git.parallelcoin.io/pod/util"
	"github.com/EXCCoin/exccd/chaincfg"
)

type nodeCfgRPCGroup struct {
	RPCUser              string   `long:"rpcuser" description:"Username for RPC connections"`
	RPCPass              string   `long:"rpcpass" default-mask:"-" description:"Password for RPC connections"`
	RPCLimitUser         string   `long:"rpclimituser" description:"Username for limited RPC connections"`
	RPCLimitPass         string   `long:"rpclimitpass" default-mask:"-" description:"Password for limited RPC connections"`
	RPCListeners         []string `long:"rpclisten" description:"Add an interface/port to listen for RPC connections (default port: 11048, testnet: 21048) gives sha256d block templates"`
	RPCCert              string   `long:"rpccert" description:"File containing the certificate file"`
	RPCKey               string   `long:"rpckey" description:"File containing the certificate key"`
	RPCMaxClients        int64    `long:"rpcmaxclients" description:"Max number of RPC clients for standard connections"`
	RPCMaxWebsockets     int64    `long:"rpcmaxwebsockets" description:"Max number of RPC websocket connections"`
	RPCMaxConcurrentReqs int64    `long:"rpcmaxconcurrentreqs" description:"Max number of concurrent RPC requests that may be processed concurrently"`
	RPCQuirks            bool     `long:"rpcquirks" description:"Mirror some JSON-RPC quirks of Bitcoin Core -- NOTE: Discouraged unless interoperability issues need to be worked around"`
	DisableRPC           bool     `long:"norpc" description:"Disable built-in RPC server -- NOTE: The RPC server is disabled by default if no rpcuser/rpcpass or rpclimituser/rpclimitpass is specified"`
	TLS                  bool     `long:"tls" description:"Enable TLS for the RPC server"`
}

type nodeCfgP2PGroup struct {
	AddPeers           []string      `short:"a" long:"addpeer" description:"Add a peer to connect with at startup"`
	ConnectPeers       []string      `long:"connect" description:"Connect only to the specified peers at startup"`
	DisableListen      bool          `long:"nolisten" description:"Disable listening for incoming connections -- NOTE: Listening is automatically disabled if the --connect or --proxy options are used without also specifying listen interfaces via --listen"`
	Listeners          []string      `long:"listen" description:"Add an interface/port to listen for connections (default all interfaces port: 11047, testnet: 21047)"`
	MaxPeers           int           `long:"maxpeers" description:"Max number of inbound and outbound peers"`
	DisableBanning     bool          `long:"nobanning" description:"Disable banning of misbehaving peers"`
	BanDuration        time.Duration `long:"banduration" description:"How long to ban misbehaving peers.  Valid time units are {s, m, h}.  Minimum 1 second"`
	BanThreshold       uint32        `long:"banthreshold" description:"Maximum allowed ban score before disconnecting and banning misbehaving peers."`
	Whitelists         []string      `long:"whitelist" description:"Add an IP network or IP that will not be banned. (eg. 192.168.1.0/24 or ::1)"`
	DisableDNSSeed     bool          `long:"nodnsseed" description:"Disable DNS seeding for peers"`
	ExternalIPs        []string      `long:"externalip" description:"Add an ip to the list of local addresses we claim to listen on to peers"`
	Proxy              string        `long:"proxy" description:"Connect via SOCKS5 proxy (eg. 127.0.0.1:9050)"`
	ProxyUser          string        `long:"proxyuser" description:"Username for proxy server"`
	ProxyPass          string        `long:"proxypass" default-mask:"-" description:"Password for proxy server"`
	OnionProxy         string        `long:"onion" description:"Connect to tor hidden services via SOCKS5 proxy (eg. 127.0.0.1:9050)"`
	OnionProxyUser     string        `long:"onionuser" description:"Username for onion proxy server"`
	OnionProxyPass     string        `long:"onionpass" default-mask:"-" description:"Password for onion proxy server"`
	NoOnion            bool          `long:"noonion" description:"Disable connecting to tor hidden services"`
	TorIsolation       bool          `long:"torisolation" description:"Enable Tor stream isolation by randomizing user credentials for each connection."`
	Upnp               bool          `long:"upnp" description:"Use UPnP to map our listening port outside of NAT"`
	UserAgentComments  []string      `long:"uacomment" description:"Comment to add to the user agent -- See BIP 14 for more information."`
	NoPeerBloomFilters bool          `long:"nopeerbloomfilters" description:"Disable bloom filtering support"`
	NoCFilters         bool          `long:"nocfilters" description:"Disable committed filtering (CF) support"`
	BlocksOnly         bool          `long:"blocksonly" description:"Do not accept transactions from remote peers."`
	RelayNonStd        bool          `long:"relaynonstd" description:"Relay non-standard transactions regardless of the default settings for the active network."`
	RejectNonStd       bool          `long:"rejectnonstd" description:"Reject non-standard transactions regardless of the default settings for the active network."`
}

type nodeCfgChainGroup struct {
	TestNet3           bool          `long:"testnet" description:"Use the test network"`
	RegressionTest     bool          `long:"regtest" description:"Use the regression test network"`
	SimNet             bool          `long:"simnet" description:"Use the simulation test network"`
	AddCheckpoints     []string      `long:"addcheckpoint" description:"Add a custom checkpoint.  Format: '<height>:<hash>'"`
	DisableCheckpoints bool          `long:"nocheckpoints" description:"Disable built-in checkpoints.  Don't do this unless you know what you're doing."`
	DbType             string        `long:"dbtype" description:"Database backend to use for the Block Chain"`
	MinRelayTxFee      float64       `long:"minrelaytxfee" description:"The minimum transaction fee in DUO/kB to be considered a non-zero fee."`
	FreeTxRelayLimit   float64       `long:"limitfreerelay" description:"Limit relay of transactions with no transaction fee to the given amount in thousands of bytes per minute"`
	NoRelayPriority    bool          `long:"norelaypriority" description:"Do not require free or low-fee transactions to have high priority for relaying"`
	TrickleInterval    time.Duration `long:"trickleinterval" description:"Minimum time between attempts to send new inventory to a connected peer"`
	MaxOrphanTxs       int           `long:"maxorphantx" description:"Max number of orphan transactions to keep in memory"`
	BlockMinSize       uint32        `long:"blockminsize" description:"Mininum block size in bytes to be used when creating a block"`
	BlockMaxSize       uint32        `long:"blockmaxsize" description:"Maximum block size in bytes to be used when creating a block"`
	BlockMinWeight     uint32        `long:"blockminweight" description:"Mininum block weight to be used when creating a block"`
	BlockMaxWeight     uint32        `long:"blockmaxweight" description:"Maximum block weight to be used when creating a block"`
	BlockPrioritySize  uint32        `long:"blockprioritysize" description:"Size in bytes for high-priority/low-fee transactions when creating a block"`
	SigCacheMaxSize    uint          `long:"sigcachemaxsize" description:"The maximum number of entries in the signature verification cache"`
	TxIndex            bool          `long:"txindex" description:"Maintain a full hash-based transaction index which makes all transactions available via the getrawtransaction RPC"`
	AddrIndex          bool          `long:"addrindex" description:"Maintain a full address-based transaction index which makes the searchrawtransactions RPC available"`
}

type nodeCfgLaunchGroup struct {
	ShowVersion   bool   `short:"V" long:"version" description:"display version information and exit"`
	ConfigFile    string `short:"C" long:"configfile" description:"path to configuration file"`
	DataDir       string `short:"b" long:"datadir" description:"directory to store data"`
	LogDir        string `long:"logdir" description:"directory to log output"`
	Profile       string `long:"profile" description:"Enable HTTP profiling on given port -- NOTE port must be between 1024 and 65536"`
	CPUProfile    string `long:"cpuprofile" description:"Write CPU profile to the specified file"`
	DropCfIndex   bool   `long:"dropcfindex" description:"Deletes the index used for committed filtering (CF) support from the database on start up and then exits."`
	DropTxIndex   bool   `long:"droptxindex" description:"Deletes the hash-based transaction index from the database on start up and then exits."`
	DropAddrIndex bool   `long:"dropaddrindex" description:"Deletes the address-based transaction index from the database on start up and then exits."`
}

type nodeCfgMiningGroup struct {
	Algo            string   `long:"algo" description:"Sets the algorithm for the CPU miner ( blake14lr, cryptonight7v2, keccak, lyra2rev2, scrypt, sha256d, stribog, skein, x11 default is 'random')"`
	Generate        bool     `long:"generate" description:"Generate (mine) bitcoins using the CPU"`
	GenThreads      int32    `long:"genthreads" description:"Number of CPU threads to use with CPU miner -1 = all cores"`
	MiningAddrs     []string `long:"miningaddr" description:"Add the specified payment address to the list of addresses to use for generated blocks, at least one is required if generate or minerport are set"`
	MinerController bool     `long:"controller" description:"Activate the miner controller"`
	MinerPort       uint16   `long:"minerport" description:"Port to listen for miner subscribers"`
	MinerPass       string   `long:"minerpass" description:"Encryption password required for miner clients to subscribe to work updates, for use over insecure connections"`
}

type nodeCfg struct {
	LaunchGroup    nodeCfgLaunchGroup `group:"Launch options"`
	NodeRPCGroup   nodeCfgRPCGroup    `group:"RPC options"`
	NodeP2PGroup   nodeCfgP2PGroup    `group:"P2P options"`
	NodeChainGroup nodeCfgChainGroup  `group:"Chain options"`
	MiningGroup    nodeCfgMiningGroup `group:"Mining options"`
	lookup         func(string) ([]net.IP, error)
	oniondial      func(string, string, time.Duration) (net.Conn, error)
	dial           func(string, string, time.Duration) (net.Conn, error)
	addCheckpoints []chaincfg.Checkpoint
	miningAddrs    []util.Address
	minerKey       []byte
	minRelayTxFee  util.Amount
	whitelists     []*net.IPNet
}

var node nodeCfg

func (n *nodeCfg) Execute(args []string) (err error) {
	fmt.Println("running full node")
	j, _ := json.MarshalIndent(n, "", "\t")
	fmt.Println(string(j))
	fmt.Println("not implemented - quitting")
	os.Exit(1)
	return
}
