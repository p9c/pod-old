package pod

import (
	"net"
	"time"

	"git.parallelcoin.io/pod/lib/chaincfg"
	"git.parallelcoin.io/pod/lib/util"
)

// Config defines the configuration options for pod. See loadConfig for details on the configuration load process.
type Config struct {
	General    generalCfg    `group:"general options"`
	Network    networkGroup  `group:"network options"`
	Ctl        ctlCfg        `command:"ctl" description:"send RPC queries to a node/wallet"`
	Node       nodeCfg       `command:"node" description:"run a core node"`
	Wallet     walletCfg     `command:"wallet" description:"run a wallet server"`
	WalletNode walletnodeCfg `command:"walletnode" description:"run a combo core/wallet server"`
	WalletGUI  walletGUICfg  `command:"walletgui" description:"run the full wallet GUI"`
	Reset      factoryReset  `command:"reset" description:"resets to factory settings"`
	// Explorer   explorerCfg   `command:"explorer" description:"run a block explorer webserver"`
	Miner minerCfg `command:"miner" description:"run the kopach miner"`
}

type minerCfg struct {
	Algo       string `long:"algo" description:"specify which algorithm to mine with"`
	Controller string `long:"controller" description:"address of miner controller to connect to"`
	Password   string `long:"password" description:"password for miner controller"`
}

type factoryReset struct {
	Really     bool `long:"really" description:"confirm factory reset"`
	Purge      bool `long:"purge" description:"delete also the data of the app"`
	OnlyPurge  bool `long:"onlypurge" description:"delete only the data of the app"`
	Ctl        bool `long:"ctl" description:"reset only ctl settings"`
	Node       bool `long:"node" description:"reset only node settings"`
	Wallet     bool `long:"wallet" description:"reset only wallet settings"`
	WalletNode bool `long:"walletnode" description:"reset only wallet/node settings"`
	WalletGUI  bool `long:"walletgui" description:"reset only wallet/node/gui settings"`
}

type generalCfg struct {
	ShowVersion bool   `long:"version" short:"V" description:"display version information and exit"`
	ConfigFile  string `long:"configfile" description:"path to configuration file"`
	DataDir     string `long:"datadir" short:"d" description:"directory to store data"`
	LogDir      string `long:"logdir" description:"directory to log output"`
	SaveConfig  bool   `long:"savecfg" description:"writes the current in-force configuration"`
}
type logTopLevel struct {
	LogLevel string `long:"debuglevel" description:"base log level that applies if no other is specified"`
}

type ctlCfg struct {
	CtlLaunch ctlCfgLaunchGroup `group:"ctl launch options"`
	CtlRPC    clientRPCGroup    `group:"ctl RPC options"`
}

type explorerCfg struct {
	LogBase        logTopLevel            `group:"logging options"`
	Logging        logSubSystems          `group:"logger subsystem options"`
	Explorer       explorerCfgLaunchGroup `group:"explorer options"`
	NodeLaunch     nodeLaunchGroup        `group:"node launch options"`
	NodeRPC        nodeCfgRPCGroup        `group:"node RPC options"`
	NodeP2P        nodeCfgP2PGroup        `group:"node P2P options"`
	NodeChain      nodeCfgChainGroup      `group:"node chain options"`
	NodeMining     nodeCfgMiningGroup     `group:"node mining options"`
	lookup         func(string) ([]net.IP, error)
	oniondial      func(string, string, time.Duration) (net.Conn, error)
	dial           func(string, string, time.Duration) (net.Conn, error)
	addCheckpoints []chaincfg.Checkpoint
	miningAddrs    []util.Address
	minerKey       []byte
	minRelayTxFee  util.Amount
	whitelists     []*net.IPNet
}

type nodeCfg struct {
	LogBase        logTopLevel        `group:"logging options"`
	Logging        logSubSystems      `group:"logger subsystem options"`
	NodeLaunch     nodeLaunchGroup    `group:"node launch options"`
	NodeRPC        nodeCfgRPCGroup    `group:"node RPC options"`
	NodeP2P        nodeCfgP2PGroup    `group:"node P2P options"`
	NodeChain      nodeCfgChainGroup  `group:"node chain options"`
	NodeMining     nodeCfgMiningGroup `group:"node mining options"`
	lookup         func(string) ([]net.IP, error)
	oniondial      func(string, string, time.Duration) (net.Conn, error)
	dial           func(string, string, time.Duration) (net.Conn, error)
	addCheckpoints []chaincfg.Checkpoint
	miningAddrs    []util.Address
	minerKey       []byte
	minRelayTxFee  util.Amount
	whitelists     []*net.IPNet
}

type walletCfg struct {
	LogBase      logTopLevel          `group:"logging options"`
	Logging      logSubSystems        `group:"logger subsystem options"`
	WalletLaunch walletCfgLaunchGroup `group:"launch options"`
	WalletRPC    walletRPCCfgGroup    `group:"wallet RPC configuration"`
	WalletNode   walletNodeCfg        `group:"wallet to node connection options"`
}

type walletGUICfg struct {
	LogBase        logTopLevel            `group:"logging options"`
	Logging        logSubSystems          `group:"logger subsystem options"`
	Explorer       explorerCfgLaunchGroup `group:"explorer options"`
	NodeLaunch     nodeLaunchGroup        `group:"node launch options"`
	NodeRPC        nodecomboCfgRPCGroup   `group:"node RPC options"`
	NodeP2P        nodeCfgP2PGroup        `group:"node P2P options"`
	NodeChain      nodeCfgChainGroup      `group:"node chain options"`
	NodeMining     nodeCfgMiningGroup     `group:"node mining options"`
	WalletLaunch   walletCfgLaunchGroup   `group:"wallet launch options"`
	WalletNode     walletNodeCfg          `group:"wallet to node connection options"`
	WalletRPC      walletcomboRPCCfgGroup `group:"wallet RPC configuration"`
	lookup         func(string) ([]net.IP, error)
	oniondial      func(string, string, time.Duration) (net.Conn, error)
	dial           func(string, string, time.Duration) (net.Conn, error)
	addCheckpoints []chaincfg.Checkpoint
	miningAddrs    []util.Address
	minerKey       []byte
	minRelayTxFee  util.Amount
	whitelists     []*net.IPNet
}

type walletnodeCfg struct {
	LogBase        logTopLevel            `group:"logging options"`
	Logging        logSubSystems          `group:"logger subsystem options"`
	Explorer       explorerCfgLaunchGroup `group:"explorer options"`
	NodeLaunch     nodeLaunchGroup        `group:"node launch options"`
	NodeRPC        nodecomboCfgRPCGroup   `group:"node RPC options"`
	NodeP2P        nodeCfgP2PGroup        `group:"node P2P options"`
	NodeChain      nodeCfgChainGroup      `group:"node chain options"`
	NodeMining     nodeCfgMiningGroup     `group:"node mining options"`
	WalletLaunch   walletCfgLaunchGroup   `group:"wallet launch options"`
	WalletNode     walletNodeCfg          `group:"wallet to node connection options"`
	WalletRPC      walletcomboRPCCfgGroup `group:"wallet RPC configuration"`
	lookup         func(string) ([]net.IP, error)
	oniondial      func(string, string, time.Duration) (net.Conn, error)
	dial           func(string, string, time.Duration) (net.Conn, error)
	addCheckpoints []chaincfg.Checkpoint
	miningAddrs    []util.Address
	minerKey       []byte
	minRelayTxFee  util.Amount
	whitelists     []*net.IPNet
}

type logSubSystems struct {
	AddrMgr    string `long:"addrmgrlog" description:"address manager log"`
	BlockChain string `long:"blockchainlog" description:"blockchain log"`
	Indexers   string `long:"indexlog" description:"indexers log"`
	ConnMgr    string `long:"connmgrlog" description:"connection manager log"`
	Database   string `long:"dblog" description:"database log"`
	Mining     string `long:"mininglog" description:"mining log"`
	Controller string `long:"controllerlog" description:"mining controller log"`
	CPUMiner   string `long:"cpuminerlog" description:"cpu miner log"`
	NetSync    string `long:"netsynclog" description:"netsync log"`
	Node       string `long:"nodelog" description:"node log"`
	Peer       string `long:"peerlog" description:"peer log"`
	RPCClient  string `long:"rpcclientlog" description:"rpc client log"`
	SPV        string `long:"spvlog" description:"light wallet log"`
	TxScript   string `long:"txscriptlog" description:"transaction script log"`
	Wallet     string `long:"walletlog" description:"wallet log"`
	WChain     string `long:"wchainlog" description:"wallet chain log"`
	LegacyRPC  string `long:"legacyrpclog" description:"legacy rpc log"`
	RPCServer  string `long:"rpcserverlog" description:"rpc server log"`
	TxMgr      string `long:"txmgrlog" description:"transaction manager log"`
	VotingPool string `long:"votelog" description:"voting pool log"`
	WTxMgr     string `long:"wtxlog" description:"wallet transaction manager log"`
}

type nodeCfgRPCGroup struct {
	RPCUser              string   `long:"rpcuser" description:"username for RPC connections"`
	RPCPass              string   `long:"rpcpass" default-mask:"-" description:"password for RPC connections"`
	RPCLimitUser         string   `long:"rpclimituser" description:"username for limited RPC connections"`
	RPCLimitPass         string   `long:"rpclimitpass" default-mask:"-" description:"password for limited RPC connections"`
	RPCListeners         []string `long:"rpclisten" description:"add an interface/port to listen for RPC connections (default port: 11048, testnet: 21048) gives sha256d block templates"`
	RPCCert              string   `long:"rpccert" description:"file containing the certificate file"`
	RPCKey               string   `long:"rpckey" description:"file containing the certificate key"`
	RPCMaxClients        int64    `long:"rpcmaxclients" description:"max number of RPC clients for standard connections"`
	RPCMaxWebsockets     int64    `long:"rpcmaxwebsockets" description:"max number of RPC websocket connections"`
	RPCMaxConcurrentReqs int64    `long:"rpcmaxconcurrentreqs" description:"max number of concurrent RPC requests that may be processed concurrently"`
	RPCQuirks            bool     `long:"rpcquirks" description:"mirror some JSON-RPC quirks of Bitcoin Core -- NOTE: Discouraged unless interoperability issues need to be worked around"`
	DisableRPC           bool     `long:"norpc" description:"misable built-in RPC server -- the RPC server is disabled by default if no rpcuser/rpcpass or rpclimituser/rpclimitpass is specified"`
	TLS                  bool     `long:"tls" description:"enable TLS for the RPC server"`
}

type nodecomboCfgRPCGroup struct {
	RPCUser              string   `long:"noderpcuser" description:"username for RPC connections"`
	RPCPass              string   `long:"noderpcpass" default-mask:"-" description:"password for RPC connections"`
	RPCLimitUser         string   `long:"noderpclimituser" description:"username for limited RPC connections"`
	RPCLimitPass         string   `long:"noderpclimitpass" default-mask:"-" description:"password for limited RPC connections"`
	RPCListeners         []string `long:"noderpclisten" description:"add an interface/port to listen for RPC connections (default port: 11048, testnet: 21048) gives sha256d block templates"`
	RPCCert              string   `long:"noderpccert" description:"file containing the certificate file"`
	RPCKey               string   `long:"noderpckey" description:"file containing the certificate key"`
	RPCMaxClients        int64    `long:"noderpcmaxclients" description:"fmx number of RPC clients for standard connections"`
	RPCMaxWebsockets     int64    `long:"noderpcmaxwebsockets" description:"max number of RPC websocket connections"`
	RPCMaxConcurrentReqs int64    `long:"noderpcmaxconcurrentreqs" description:"max number of concurrent RPC requests that may be processed concurrently"`
	RPCQuirks            bool     `long:"noderpcquirks" description:"mirror some JSON-RPC quirks of Bitcoin Core -- NOTE: Discouraged unless interoperability issues need to be worked around"`
	DisableRPC           bool     `long:"nodenorpc" description:"disable built-in RPC server -- NOTE: The RPC server is disabled by default if no rpcuser/rpcpass or rpclimituser/rpclimitpass is specified"`
	TLS                  bool     `long:"nodetls" description:"enable TLS for the RPC server"`
}

type nodeCfgP2PGroup struct {
	AddPeers           []string      `long:"addpeer" description:"add a peer to connect with at startup"`
	ConnectPeers       []string      `long:"connect" description:"connect only to the specified peers at startup"`
	DisableListen      bool          `long:"nolisten" description:"disable listening for incoming connections -- NOTE: Listening is automatically disabled if the --connect or --proxy options are used without also specifying listen interfaces via --listen"`
	Listeners          []string      `long:"listen" description:"add an interface/port to listen for connections (default all interfaces port: 11047, testnet: 21047)"`
	MaxPeers           int           `long:"maxpeers" description:"max number of inbound and outbound peers"`
	DisableBanning     bool          `long:"nobanning" description:"disable banning of misbehaving peers"`
	BanDuration        time.Duration `long:"banduration" description:"how long to ban misbehaving peers.  Valid time units are {s, m, h}.  Minimum 1 second"`
	BanThreshold       uint32        `long:"banthreshold" description:"maximum allowed ban score before disconnecting and banning misbehaving peers."`
	Whitelists         []string      `long:"whitelist" description:"add an IP network or IP that will not be banned. (eg. 192.168.1.0/24 or ::1)"`
	DisableDNSSeed     bool          `long:"nodnsseed" description:"disable DNS seeding for peers"`
	ExternalIPs        []string      `long:"externalip" description:"add an ip to the list of local addresses we claim to listen on to peers"`
	Proxy              string        `long:"proxy" description:"connect via SOCKS5 proxy (eg. 127.0.0.1:9050)"`
	ProxyUser          string        `long:"proxyuser" description:"username for proxy server"`
	ProxyPass          string        `long:"proxypass" default-mask:"-" description:"password for proxy server"`
	OnionProxy         string        `long:"onion" description:"Connect to tor hidden services via SOCKS5 proxy (eg. 127.0.0.1:9050)"`
	OnionProxyUser     string        `long:"onionuser" description:"username for onion proxy server"`
	OnionProxyPass     string        `long:"onionpass" default-mask:"-" description:"password for onion proxy server"`
	NoOnion            bool          `long:"noonion" description:"disable connecting to tor hidden services"`
	TorIsolation       bool          `long:"torisolation" description:"enable Tor stream isolation by randomizing user credentials for each connection."`
	Upnp               bool          `long:"upnp" description:"use UPnP to map our listening port outside of NAT"`
	UserAgentComments  []string      `long:"uacomment" description:"comment to add to the user agent -- See BIP 14 for more information."`
	NoPeerBloomFilters bool          `long:"nopeerbloomfilters" description:"disable bloom filtering support"`
	NoCFilters         bool          `long:"nocfilters" description:"disable committed filtering (CF) support"`
	BlocksOnly         bool          `long:"blocksonly" description:"do not accept transactions from remote peers."`
	RelayNonStd        bool          `long:"relaynonstd" description:"relay non-standard transactions regardless of the default settings for the active network."`
	RejectNonStd       bool          `long:"rejectnonstd" description:"reject non-standard transactions regardless of the default settings for the active network."`
}

type nodeCfgChainGroup struct {
	AddCheckpoints     []string      `long:"addcheckpoint" description:"add a custom checkpoint.  Format: '<height>:<hash>'"`
	DisableCheckpoints bool          `long:"nocheckpoints" description:"disable built-in checkpoints.  Don't do this unless you know what you're doing."`
	DbType             string        `long:"dbtype" description:"database backend to use for the blockchain"`
	MinRelayTxFee      float64       `long:"minrelaytxfee" description:"the minimum transaction fee in DUO/kB to be considered a non-zero fee."`
	FreeTxRelayLimit   float64       `long:"limitfreerelay" description:"limit relay of transactions with no transaction fee to the given amount in thousands of bytes per minute"`
	NoRelayPriority    bool          `long:"norelaypriority" description:"do not require free or low-fee transactions to have high priority for relaying"`
	TrickleInterval    time.Duration `long:"trickleinterval" description:"minimum time between attempts to send new inventory to a connected peer"`
	MaxOrphanTxs       int           `long:"maxorphantx" description:"max number of orphan transactions to keep in memory"`
	BlockMinSize       uint32        `long:"blockminsize" description:"mininum block size in bytes to be used when creating a block"`
	BlockMaxSize       uint32        `long:"blockmaxsize" description:"maximum block size in bytes to be used when creating a block"`
	BlockMinWeight     uint32        `long:"blockminweight" description:"mininum block weight to be used when creating a block"`
	BlockMaxWeight     uint32        `long:"blockmaxweight" description:"maximum block weight to be used when creating a block"`
	BlockPrioritySize  uint32        `long:"blockprioritysize" description:"size in bytes for high-priority/low-fee transactions when creating a block"`
	SigCacheMaxSize    uint          `long:"sigcachemaxsize" description:"the maximum number of entries in the signature verification cache"`
	TxIndex            bool          `long:"txindex" description:"maintain a full hash-based transaction index which makes all transactions available via the getrawtransaction RPC"`
	AddrIndex          bool          `long:"addrindex" description:"maintain a full address-based transaction index which makes the searchrawtransactions RPC available"`
}

type nodeCfgMiningGroup struct {
	Algo          string   `long:"algo" description:"sets the algorithm for the CPU miner (blake14lr, cryptonight7v2, keccak, lyra2rev2, scrypt, sha256d, stribog, skein, x11, easy, random)"`
	Generate      bool     `long:"generate" description:"generate (mine) bitcoins using the CPU"`
	GenThreads    int32    `long:"genthreads" description:"number of CPU threads to use with CPU miner -1 = all cores"`
	MiningAddrs   []string `long:"miningaddr" description:"add the specified payment address to the list of addresses to use for generated blocks, at least one is required if generate or minerport are set"`
	MinerListener string   `long:"minerlistener" description:"listen for miner work subscription requests and such"`
	MinerPass     string   `long:"minerpass" description:"encryption password required for miner clients to subscribe to work updates, for use over insecure connections"`
}

type ctlCfgLaunchGroup struct {
	ListCommands bool `long:"listcommands" short:"l" description:"list available commands"`
	Wallet       bool `long:"wallet" description:"connect to wallet"`
}

type clientRPCGroup struct {
	RPCUser       string `long:"rpcuser" description:"RPC username"`
	RPCPassword   string `long:"rpcpass" default-mask:"-" description:"RPC password"`
	RPCServer     string `long:"rpcserver" description:"RPC server to connect to"`
	RPCCert       string `long:"rpccert" description:"RPC server certificate chain for validation"`
	TLS           bool   `long:"tls" description:"enable TLS"`
	Proxy         string `long:"proxy" description:"connect via SOCKS5 proxy (eg. 127.0.0.1:9050)"`
	ProxyUser     string `long:"proxyuser" description:"username for proxy server"`
	ProxyPass     string `long:"proxypass" default-mask:"-" description:"password for proxy server"`
	TLSSkipVerify bool   `long:"skipverify" description:"do not verify tls certificates (not recommended!)"`
}

type ctlCfgChainGroup struct {
	TestNet3       bool `long:"testnet" description:"connect to testnet"`
	SimNet         bool `long:"simnet" description:"connect to the simulation test network"`
	RegressionTest bool `long:"regtest" description:"connect to the regression test network"`
}

type explorerCfgLaunchGroup struct {
	WebserverAddress string `long:"webserver" description:"address to listen for explorer clients"`
}

type nodeLaunchGroup struct {
	Profile       string `long:"nodeprofile" description:"enable HTTP profiling on given port - port must be between 1024 and 65536"`
	CPUProfile    string `long:"nodecpuprofile" description:"write CPU profile to the specified file"`
	DropCfIndex   bool   `long:"dropcfindex" description:"deletes the index used for committed filtering (CF) support from the database on start up and then exits."`
	DropTxIndex   bool   `long:"droptxindex" description:"deletes the hash-based transaction index from the database on start up and then exits."`
	DropAddrIndex bool   `long:"dropaddrindex" description:"deletes the address-based transaction index from the database on start up and then exits."`
}

type walletCfgLaunchGroup struct {
	Profile       string `long:"walletprofile" description:"enable HTTP profiling on given port - port must be between 1024 and 65536"`
	CPUProfile    string `long:"walletcpuprofile" description:"write CPU profile to the specified file"`
	Create        bool   `long:"createwallet" description:"create the wallet if it does not exist"`
	CreateTemp    bool   `long:"createtemp" description:"create a temporary simulation wallet (pass=password) in the data directory indicated; must call with --datadir"`
	NoInitialLoad bool   `long:"noinitialload" description:"defer wallet creation/opening on startup and enable loading wallets over RPC"`
}

type networkGroup struct {
	TestNet3       bool `long:"testnet" description:"connect to testnet"`
	SimNet         bool `long:"simnet" description:"connect to the simulation test network"`
	RegressionTest bool `long:"regtest" description:"connect to the regression test network"`
}

type walletNodeCfg struct {
	RPCConnect      string        `long:"noderpcconnect" description:"hostname/IP and port of pod RPC server to connect to (default localhost:11048, testnet: localhost:21048, simnet: localhost:41048)"`
	CAFile          string        `long:"nodecafile" description:"file containing root certificates to authenticate a TLS connections with pod"`
	EnableClientTLS bool          `long:"nodeclienttls" description:"enable TLS for the RPC client"`
	PodUsername     string        `long:"nodeusername" description:"username for pod authentication"`
	PodPassword     string        `long:"nodepassword" default-mask:"-" description:"password for pod authentication"`
	Proxy           string        `long:"nodeproxy" description:"connect via SOCKS5 proxy (eg. 127.0.0.1:9050)"`
	ProxyUser       string        `long:"nodeproxyuser" description:"username for proxy server"`
	ProxyPass       string        `long:"nodeproxypass" description:"password for proxy server"`
	AddPeers        []string      `long:"nodeaddpeer" description:"add a peer to connect with at startup"`
	ConnectPeers    []string      `long:"nodeconnect" description:"connect only to the specified peers at startup"`
	MaxPeers        int           `long:"nodemaxpeers" description:"max number of inbound and outbound peers"`
	BanDuration     time.Duration `long:"nodebanduration" description:"how long to ban misbehaving peers.  Valid time units are {s, m, h}.  Minimum 1 second"`
	BanThreshold    uint32        `long:"nodebanthreshold" description:"maximum allowed ban score before disconnecting and banning misbehaving peers."`
}

type walletRPCCfgGroup struct {
	Username                 string   `long:"username" description:"username for legacy RPC and pod authentication (if podusername is unset)"`
	Password                 string   `long:"password" default-mask:"-" description:"password for legacy RPC and pod authentication (if podpassword is unset)"`
	RPCCert                  string   `long:"rpccert" description:"file containing the certificate file"`
	RPCKey                   string   `long:"rpckey" description:"file containing the certificate key"`
	OneTimeTLSKey            bool     `long:"onetimetlskey" description:"generate a new TLS certpair at startup, but only write the certificate to disk"`
	EnableServerTLS          bool     `long:"servertls" description:"enable TLS for the RPC server"`
	LegacyRPCListeners       []string `long:"legacyrpclisten" description:"listen for legacy RPC connections on this interface/port (default port: 11046, testnet: 21046, simnet: 41046)"`
	LegacyRPCMaxClients      int64    `long:"legacyrpcmaxclients" description:"max number of legacy RPC clients for standard connections"`
	LegacyRPCMaxWebsockets   int64    `long:"legacyrpcmaxwebsockets" description:"max number of legacy RPC websocket connections"`
	ExperimentalRPCListeners []string `long:"experimentalrpclisten" description:"listen for RPC connections on this interface/port"`
	RPCMaxClients            int64    `long:"rpcmaxclients" description:"max number of RPC clients for standard connections"`
	RPCMaxWebsockets         int64    `long:"rpcmaxwebsockets" description:"max number of RPC websocket connections"`
	RPCMaxConcurrentReqs     int64    `long:"rpcmaxconcurrentreqs" description:"max number of concurrent RPC requests that may be processed concurrently"`
	RPCQuirks                bool     `long:"rpcquirks" description:"mirror some JSON-RPC quirks of Bitcoin Core -- NOTE: Discouraged unless interoperability issues need to be worked around"`
	DisableRPC               bool     `long:"norpc" description:"disable built-in RPC server -- NOTE: The RPC server is disabled by default if no rpcuser/rpcpass or rpclimituser/rpclimitpass is specified"`
	TLS                      bool     `long:"tls" description:"enable TLS for the RPC server"`
}

type walletcomboRPCCfgGroup struct {
	Username                 string   `long:"walletusername" description:"username for legacy RPC and pod authentication (if podusername is unset)"`
	Password                 string   `long:"walletpassword" default-mask:"-" description:"password for legacy RPC and pod authentication (if podpassword is unset)"`
	RPCCert                  string   `long:"walletrpccert" description:"file containing the certificate file"`
	RPCKey                   string   `long:"walletrpckey" description:"file containing the certificate key"`
	OneTimeTLSKey            bool     `long:"walletonetimetlskey" description:"generate a new TLS certpair at startup, but only write the certificate to disk"`
	EnableServerTLS          bool     `long:"walletservertls" description:"enable TLS for the RPC server"`
	LegacyRPCListeners       []string `long:"walletlegacyrpclisten" description:"listen for legacy RPC connections on this interface/port (default port: 11046, testnet: 21046, simnet: 41046)"`
	LegacyRPCMaxClients      int64    `long:"walletlegacyrpcmaxclients" description:"max number of legacy RPC clients for standard connections"`
	LegacyRPCMaxWebsockets   int64    `long:"walletlegacyrpcmaxwebsockets" description:"max number of legacy RPC websocket connections"`
	ExperimentalRPCListeners []string `long:"walletexperimentalrpclisten" description:"listen for RPC connections on this interface/port"`
	RPCMaxClients            int64    `long:"walletrpcmaxclients" description:"max number of RPC clients for standard connections"`
	RPCMaxWebsockets         int64    `long:"walletrpcmaxwebsockets" description:"max number of RPC websocket connections"`
	RPCMaxConcurrentReqs     int64    `long:"walletrpcmaxconcurrentreqs" description:"max number of concurrent RPC requests that may be processed concurrently"`
	RPCQuirks                bool     `long:"walletrpcquirks" description:"mirror some JSON-RPC quirks of Bitcoin Core -- NOTE: Discouraged unless interoperability issues need to be worked around"`
	DisableRPC               bool     `long:"walletnorpc" description:"disable built-in RPC server -- NOTE: The RPC server is disabled by default if no rpcuser/rpcpass or rpclimituser/rpclimitpass is specified"`
	TLS                      bool     `long:"wallettls" description:"enable TLS for the RPC server"`
}

// type walletSpvGUICfg struct {
// 	Main           allCfgLaunch           `group:"launch options"`
// 	Net            networkGroup           `group:"network options"`
// 	Explorer       explorerCfgLaunchGroup `group:"explorer options"`
// 	NodeLaunch     nodeLaunchGroup        `group:"node launch options"`
// 	NodeRPCGroup   nodecomboCfgRPCGroup   `group:"node RPC options"`
// 	NodeP2P        nodeCfgP2PGroup        `group:"node P2P options"`
// 	NodeChain      nodeCfgChainGroup      `group:"node chain options"`
// 	LaunchWallet   walletCfgLaunchGroup   `group:"wallet launch options"`
// 	WalletCfg      walletNodeCfg          `group:"wallet to node connection options"`
// 	WalletRPCCfg   walletcomboRPCCfgGroup `group:"wallet RPC configuration"`
// 	lookup         func(string) ([]net.IP, error)
// 	oniondial      func(string, string, time.Duration) (net.Conn, error)
// 	dial           func(string, string, time.Duration) (net.Conn, error)
// 	addCheckpoints []chaincfg.Checkpoint
// 	miningAddrs    []util.Address
// 	minerKey       []byte
// 	minRelayTxFee  util.Amount
// 	whitelists     []*net.IPNet
// }

// type walletSpvCfg struct {
// 	Main           allCfgLaunch           `group:"launch options"`
// 	Net            networkGroup           `group:"network options"`
// 	Launch         explorerCfgLaunchGroup `group:"explorer options"`
// 	NodeLaunch     nodeLaunchGroup        `group:"node launch options"`
// 	NodeRPCGroup   nodecomboCfgRPCGroup   `group:"node RPC options"`
// 	NodeP2P        nodeCfgP2PGroup        `group:"node P2P options"`
// 	NodeChain      nodeCfgChainGroup      `group:"node Chain options"`
// 	LaunchWallet   walletCfgLaunchGroup   `group:"launch options"`
// 	WalletCfg      walletNodeCfg          `group:"wallet to node connection options"`
// 	WalletRPCCfg   walletcomboRPCCfgGroup `group:"wallet RPC configuration"`
// 	lookup         func(string) ([]net.IP, error)
// 	oniondial      func(string, string, time.Duration) (net.Conn, error)
// 	dial           func(string, string, time.Duration) (net.Conn, error)
// 	addCheckpoints []chaincfg.Checkpoint
// 	miningAddrs    []util.Address
// 	minerKey       []byte
// 	minRelayTxFee  util.Amount
// 	whitelists     []*net.IPNet
// }

// type spvCfg struct {
// 	Main           allCfgLaunch       `group:"general launch options"`
// 	Launch         nodeLaunchGroup    `group:"node launch options"`
// 	RPCGroup       nodeCfgRPCGroup    `group:"node RPC options"`
// 	P2PGroup       nodeCfgP2PGroup    `group:"node P2P options"`
// 	ChainGroup     nodeCfgChainGroup  `group:"node chain options"`
// 	MiningGroup    nodeCfgMiningGroup `group:"node mining options"`
// 	lookup         func(string) ([]net.IP, error)
// 	oniondial      func(string, string, time.Duration) (net.Conn, error)
// 	dial           func(string, string, time.Duration) (net.Conn, error)
// 	addCheckpoints []chaincfg.Checkpoint
// 	miningAddrs    []util.Address
// 	minerKey       []byte
// 	minRelayTxFee  util.Amount
// 	whitelists     []*net.IPNet
// }
