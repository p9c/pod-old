package pod

import (
	"time"

	"gopkg.in/urfave/cli.v1"
)

type Config struct {
	ConfigFile               *string          `toml:"configfile"`
	DataDir                  *string          `toml:"datadir"`
	LogDir                   *string          `toml:"logdir"`
	LogLevel                 *string          `toml:"loglevel"`
	Subsystems               *cli.StringSlice `toml:"subsystems"`
	Network                  *string          `toml:"network"`
	AddPeers                 *cli.StringSlice `toml:"addpeers"`
	ConnectPeers             *cli.StringSlice `toml:"connectpeers"`
	MaxPeers                 *int             `toml:"maxpeers"`
	Listeners                *cli.StringSlice `toml:"listeners"`
	DisableListen            *bool            `toml:"disablelisten"`
	DisableBanning           *bool            `toml:"disablebanning"`
	BanDuration              *time.Duration   `toml:"banduration"`
	BanThreshold             *int             `toml:"banthreshold"`
	Whitelists               *cli.StringSlice `toml:"whitelists"`
	Username                 *string          `toml:"username"`
	Password                 *string          `toml:"password"`
	ServerUser               *string          `toml:"serveruser"`
	ServerPass               *string          `toml:"serverpass"`
	LimitUser                *string          `toml:"limituser"`
	LimitPass                *string          `toml:"limitpass"`
	RPCConnect               *string          `toml:"rpcconnect"`
	RPCListeners             *cli.StringSlice `toml:"rpclisteners"`
	RPCCert                  *string          `toml:"rpccert"`
	RPCKey                   *string          `toml:"rpckey"`
	RPCMaxClients            *int             `toml:"rpcmaxclients"`
	RPCMaxWebsockets         *int             `toml:"rpcmaxwebsockets"`
	RPCMaxConcurrentReqs     *int             `toml:"rpcmaxconcurrentreqs"`
	RPCQuirks                *bool            `toml:"rpcquirks"`
	DisableRPC               *bool            `toml:"disablerpc"`
	TLS                      *bool            `toml:"tls"`
	DisableDNSSeed           *bool            `toml:"disablednsseed"`
	ExternalIPs              *cli.StringSlice `toml:"externalips"`
	Proxy                    *string          `toml:"proxy"`
	ProxyUser                *string          `toml:"proxyuser"`
	ProxyPass                *string          `toml:"proxypass"`
	OnionProxy               *string          `toml:"onionproxy"`
	OnionProxyUser           *string          `toml:"onionproxyuser"`
	OnionProxyPass           *string          `toml:"onionproxypass"`
	Onion                    *bool            `toml:"onion"`
	TorIsolation             *bool            `toml:"torisolation"`
	TestNet3                 *bool            `toml:"testnet3" omitempty:"true"`
	RegressionTest           *bool            `toml:"regressiontest" omitempty:"true"`
	SimNet                   *bool            `toml:"simnet" omitempty:"true"`
	AddCheckpoints           *cli.StringSlice `toml:"addcheckpoints"`
	DisableCheckpoints       *bool            `toml:"disablecheckpoints"`
	DbType                   *string          `toml:"dbtype"`
	Profile                  *string          `toml:"profile"`
	CPUProfile               *string          `toml:"cpuprofile"`
	Upnp                     *bool            `toml:"upnp"`
	MinRelayTxFee            *float64         `toml:"minrelaytxfee"`
	FreeTxRelayLimit         *float64         `toml:"freetxrelaylimit"`
	NoRelayPriority          *bool            `toml:"norelaypriority"`
	TrickleInterval          *time.Duration   `toml:"trickleinterval"`
	MaxOrphanTxs             *int             `toml:"maxorphantxs"`
	Algo                     *string          `toml:"algo"`
	Generate                 *bool            `toml:"generate"`
	GenThreads               *int             `toml:"genthreads"`
	MiningAddrs              *cli.StringSlice `toml:"miningaddrs"`
	MinerListener            *string          `toml:"minerlistener"`
	MinerPass                *string          `toml:"minerpass"`
	BlockMinSize             *int             `toml:"blockminsize"`
	BlockMaxSize             *int             `toml:"blockmaxsize"`
	BlockMinWeight           *int             `toml:"blockminweight"`
	BlockMaxWeight           *int             `toml:"blockmaxweight"`
	BlockPrioritySize        *int             `toml:"blockprioritysize"`
	UserAgentComments        *cli.StringSlice `toml:"useragentcomments"`
	NoPeerBloomFilters       *bool            `toml:"nopeerbloomfilters"`
	NoCFilters               *bool            `toml:"nocfilters"`
	SigCacheMaxSize          *int             `toml:"sigcachemaxsize"`
	BlocksOnly               *bool            `toml:"blocksonly"`
	TxIndex                  *bool            `toml:"txindex"`
	AddrIndex                *bool            `toml:"addrindex"`
	RelayNonStd              *bool            `toml:"relaynonstd"`
	RejectNonStd             *bool            `toml:"rejectnonstd"`
	TLSSkipVerify            *bool            `toml:"tlsskipverify"`
	Wallet                   *bool            `toml:"wallet"`
	NoInitialLoad            *bool            `toml:"noinitialload"`
	WalletPass               *string          `toml:"walletpass"`
	WalletServer             *string          `toml:"walletserver"`
	CAFile                   *string          `toml:"cafile"`
	OneTimeTLSKey            *bool            `toml:"onetimetlskey"`
	ServerTLS                *bool            `toml:"servertls"`
	LegacyRPCListeners       *cli.StringSlice `toml:"legacyrpclisteners"`
	LegacyRPCMaxClients      *int             `toml:"legacyrpcmaxclients"`
	LegacyRPCMaxWebsockets   *int             `toml:"legacyrpcmaxwebsockets"`
	ExperimentalRPCListeners *cli.StringSlice `toml:"experimentalrpclisteners"`
}
