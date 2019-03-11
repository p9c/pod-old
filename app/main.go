package app

import (
	"fmt"
	"git.parallelcoin.io/pod/cmd/node"
	"git.parallelcoin.io/pod/cmd/node/mempool"
	"gopkg.in/urfave/cli.v1"
	"gopkg.in/urfave/cli.v1/altsrc"
	"os"
	"time"
)

// TODO: create common base with all common common fields linked as pointers

var App = cli.NewApp()

func Main() int {
	App.Before = altsrc.InitInputSourceWithContext(App.Flags,
		NewYamlSourceFromFlagAndNameFunc(podConfigFilename, "datadir"))
	e := App.Run(os.Args)
	if e != nil {
		return 1
	}
	return 0
}

func init() {

	*App = cli.App{
		Name:        "pod",
		Version:     "v0.0.1",
		Description: "Parallelcoin Pod Suite -- All-in-one everything for Parallelcoin!",
		Copyright:   "Legacy portions derived from btcsuite/btcd under ISC licence. The remainder is already in your possession. Use it wisely.",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "datadir, D",
				Value:       appDatadir,
				Usage:       "sets the data directory base for a pod instance",
				EnvVar:      "POD_DATADIR",
				Destination: &appConfigCommon.Datadir,
			},
			cli.BoolFlag{
				Name:        "save, i",
				Usage:       "save settings as effective from invocation",
				Destination: &appConfigCommon.Save,
			},
			altsrc.NewStringFlag(cli.StringFlag{
				Name:        "loglevel",
				Value:       "info",
				Usage:       "sets the base for all subsystem logging",
				EnvVar:      "POD_LOGLEVEL",
				Destination: &appConfigCommon.Loglevel,
			}),
			cli.StringSliceFlag{
				Name:  "subsystems",
				Usage: "sets individual subsystems log levels, use 'listsubsystems' to list available with list syntax",
				Value: &appConfigCommon.Subsystems,
			},
			cli.StringFlag{
				Name:        "network",
				Value:       "mainnet",
				Usage:       "connect to mainnet/testnet3/simnet",
				Destination: &appConfigCommon.Network,
			},
			cli.StringFlag{
				Name:        "rpccert",
				Value:       defaultDatadir + "/rpc.cert",
				Usage:       "File containing the certificate file",
				Destination: &appConfigCommon.RPCcert,
			},
			cli.StringFlag{
				Name:        "rpckey",
				Value:       defaultDatadir + "/rpc.key",
				Usage:       "File containing the certificate key",
				Destination: &appConfigCommon.RPCkey,
			},
			cli.StringFlag{
				Name:        "cafile",
				Value:       defaultDatadir + "/cafile",
				Usage:       "File containing root certificates to authenticate a TLS connections with pod",
				Destination: &appConfigCommon.CAfile,
			},
			cli.BoolFlag{
				Name:        "tls, clienttls",
				Usage:       "Enable TLS for client connections",
				Destination: &appConfigCommon.ClientTLS,
			},
			cli.BoolFlag{
				Name:        "servertls",
				Usage:       "Enable TLS for server connections",
				Destination: &appConfigCommon.ServerTLS,
			},
			cli.BoolFlag{
				Name:        "useproxy, r",
				Usage:       "use configured proxy",
				Destination: &appConfigCommon.Useproxy,
			},
			cli.StringFlag{
				Name:        "proxy",
				Value:       "localhost:9050",
				Usage:       "Connect via SOCKS5 proxy",
				Destination: &appConfigCommon.Proxy,
			},
			cli.StringFlag{
				Name:        "proxyuser",
				Value:       "user",
				Usage:       "Username for proxy server",
				Destination: &appConfigCommon.Proxyuser,
			},
			cli.StringFlag{
				Name:        "proxypass",
				Value:       "pa55word",
				Usage:       "Password for proxy server",
				Destination: &appConfigCommon.Proxypass,
			},
			cli.BoolFlag{
				Name:        "onion",
				Usage:       "Enable connecting to tor hidden services",
				Destination: &appConfigCommon.Onion,
			},
			cli.StringFlag{
				Name:        "onionproxy",
				Value:       "localhost:9050",
				Usage:       "Connect to tor hidden services via SOCKS5 proxy (eg. 127.0.0.1:9050)",
				Destination: &appConfigCommon.OnionProxy,
			},
			cli.StringFlag{
				Name:        "onionuser",
				Value:       "user",
				Usage:       "Username for onion proxy server",
				Destination: &appConfigCommon.Onionuser,
			},
			cli.StringFlag{
				Name:        "onionpass",
				Value:       "pa55word",
				Usage:       "Password for onion proxy server",
				Destination: &appConfigCommon.Onionpass,
			},
			cli.BoolFlag{
				Name:        "torisolation",
				Usage:       "Enable Tor stream isolation by randomizing user credentials for each connection.",
				Destination: &appConfigCommon.Torisolation,
			},
		},
		Commands: []cli.Command{
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "print version and exit",
				Action: func(c *cli.Context) error {
					fmt.Println(c.App.Name, c.App.Version)
					return nil
				},
			},
			{
				Name:    "listsubsystems",
				Aliases: []string{"l"},
				Usage:   "list available logging subsystems",
				Action: func(c *cli.Context) error {
					fmt.Println("todo list logging subsystems")
					return nil
				},
			},
			{
				Name:    "ctl",
				Aliases: []string{"c"},
				Usage:   "send RPC commands to a node or wallet and print the result",
				Action:  ctlHandle,
				Subcommands: []cli.Command{
					{
						Name:    "listcommands",
						Aliases: []string{"list", "l"},
						Usage:   "list commands available at endpoint",
						Action:  ctlHandleList,
					},
				},
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:        "rpcserver, server, s",
						Value:       "localhost:11048",
						Usage:       "set rpc password",
						Destination: ctlConfig.RPCServer,
					},
					cli.StringFlag{
						Name:        "walletserver, ws",
						Value:       "localhost:11046",
						Usage:       "set wallet server to use",
						Destination: ctlConfig.Wallet,
					},
					cli.StringFlag{
						Name:        "rpcusername, username, user, u",
						Value:       "user",
						Usage:       "set rpc username",
						Destination: ctlConfig.RPCUser,
					},
					cli.StringFlag{
						Name:        "rpcpassword, password, pass, p",
						Value:       "pa55word",
						Usage:       "set rpc password",
						Destination: ctlConfig.RPCPass,
					},
					cli.BoolFlag{
						Name:  "wallet, w",
						Usage: "use configured wallet rpc address",
					},
				},
			},
			nodeCommand,
			walletCommand,
			{
				Name:    "init",
				Aliases: []string{"i"},
				Usage:   "resets configuration to factory",
				Action: func(c *cli.Context) error {
					fmt.Println("resetting configuration")
					return nil
				},
			},
			{
				Name:    "conf",
				Aliases: []string{"C"},
				Usage:   "automate configuration setup for testnets etc",
				Action: func(c *cli.Context) error {
					fmt.Println("calling conf")
					return nil
				},
			},

			{
				Name:    "shell",
				Aliases: []string{"s"},
				Usage:   "start combined wallet/node shell",
				Action: func(c *cli.Context) error {
					fmt.Println("calling shell")
					return nil
				},
			},
			{
				Name:    "gui",
				Aliases: []string{"g"},
				Usage:   "start GUI (TODO: should ultimately be default)",
				Action: func(c *cli.Context) error {
					fmt.Println("calling gui")
					return nil
				},
			},
		},
	}
}

var nodeCommand = cli.Command{
	Name:    "node",
	Aliases: []string{"n"},
	Usage:   "start parallelcoin full node",
	Action:  nodeHandle,
	Subcommands: []cli.Command{
		{
			Name:  "droptxindex",
			Usage: "Deletes the hash-based transaction index from the database on start up and exits.",
		},
		{
			Name:  "dropaddrindex",
			Usage: "Deletes the address-based transaction index from the database on start up and exits.",
		},
		{
			Name:  "dropcfindex",
			Usage: "Deletes the index used for committed filtering (CF) support from the database on start up and exits.",
		},
	},
	Flags: []cli.Flag{
		cli.StringSliceFlag{
			Name:  "addpeer",
			Value: nodeConfig.AddPeers,
			Usage: "Add a peer to connect with at startup",
		},
		cli.StringSliceFlag{
			Name:  "connect",
			Value: nodeConfig.ConnectPeers,
			Usage: "Connect only to the specified peers at startup",
		},
		cli.BoolFlag{
			Name:        "nolisten",
			Usage:       "Disable listening for incoming connections -- NOTE: Listening is automatically disabled if the --connect or --proxy options are used without also specifying listen interfaces via --listen",
			Destination: nodeConfig.DisableListen,
		},
		cli.StringSliceFlag{
			Name:  "listen",
			Value: nodeConfig.Listeners,
			Usage: "Add an interface/port to listen for connections",
		},
		cli.IntFlag{
			Name:        "maxpeers",
			Value:       node.DefaultMaxPeers,
			Usage:       "Max number of inbound and outbound peers",
			Destination: nodeConfig.MaxPeers,
		},
		cli.BoolFlag{
			Name:        "nobanning",
			Usage:       "Disable banning of misbehaving peers",
			Destination: nodeConfig.DisableBanning,
		},
		cli.DurationFlag{
			Name:        "banduration",
			Value:       time.Hour * 24,
			Usage:       "How long to ban misbehaving peers",
			Destination: nodeConfig.BanDuration,
		},
		cli.IntFlag{
			Name:        "banthreshold",
			Value:       node.DefaultBanThreshold,
			Usage:       "Maximum allowed ban score before disconnecting and banning misbehaving peers.",
			Destination: nodeConfig.BanThreshold,
		},
		cli.StringSliceFlag{
			Name:  "whitelist",
			Usage: "Add an IP network or IP that will not be banned. (eg. 192.168.1.0/24 or ::1)",
			Value: nodeConfig.Whitelists,
		},
		cli.StringFlag{
			Name:        "rpcuser",
			Value:       "user",
			Usage:       "Username for RPC connections",
			Destination: nodeConfig.RPCUser,
		},
		cli.StringFlag{
			Name:        "rpcpass",
			Value:       "pa55word",
			Usage:       "Password for RPC connections",
			Destination: nodeConfig.RPCPass,
		},
		cli.StringFlag{
			Name:        "rpclimituser",
			Value:       "user",
			Usage:       "Username for limited RPC connections",
			Destination: nodeConfig.RPCLimitUser,
		},
		cli.StringFlag{
			Name:        "rpclimitpass",
			Value:       "pa55word",
			Usage:       "Password for limited RPC connections",
			Destination: nodeConfig.RPCLimitPass,
		},
		cli.StringSliceFlag{
			Name:  "rpclisten",
			Value: nodeConfig.RPCListeners,
			Usage: "Add an interface/port to listen for RPC connections (default port: 11048, testnet: 21048) gives sha256d block templates",
		},

		cli.IntFlag{
			Name:        "rpcmaxclients",
			Value:       node.DefaultMaxRPCClients,
			Usage:       "Max number of RPC clients for standard connections",
			Destination: nodeConfig.RPCMaxClients,
		},
		cli.IntFlag{
			Name:        "rpcmaxwebsockets",
			Value:       node.DefaultMaxRPCWebsockets,
			Usage:       "Max number of RPC websocket connections",
			Destination: nodeConfig.RPCMaxWebsockets,
		},
		cli.IntFlag{
			Name:        "rpcmaxconcurrentreqs",
			Value:       node.DefaultMaxRPCConcurrentReqs,
			Usage:       "Max number of concurrent RPC requests that may be processed concurrently",
			Destination: nodeConfig.RPCMaxConcurrentReqs,
		},
		cli.BoolFlag{
			Name:        "rpcquirks",
			Usage:       "Mirror some JSON-RPC quirks of Bitcoin Core -- NOTE: Discouraged unless interoperability issues need to be worked around",
			Destination: nodeConfig.RPCQuirks,
		},
		cli.BoolFlag{
			Name:        "norpc",
			Usage:       "Disable built-in RPC server -- NOTE: The RPC server is disabled by default if no rpcuser/rpcpass or rpclimituser/rpclimitpass is specified",
			Destination: nodeConfig.DisableRPC,
		},

		cli.BoolFlag{
			Name:        "nodnsseed",
			Usage:       "Disable DNS seeding for peers",
			Destination: nodeConfig.DisableDNSSeed,
		},
		cli.StringSliceFlag{
			Name:  "externalip",
			Value: nodeConfig.ExternalIPs,
			Usage: "Add an ip to the list of local addresses we claim to listen on to peers",
		},
		cli.StringSliceFlag{
			Name:  "addcheckpoint",
			Value: nodeConfig.AddCheckpoints,
			Usage: "Add a custom checkpoint.  Format: '<height>:<hash>'",
		},
		cli.BoolFlag{
			Name:        "nocheckpoints",
			Usage:       "Disable built-in checkpoints.  Don't do this unless you know what you're doing.",
			Destination: nodeConfig.DisableCheckpoints,
		},
		cli.StringFlag{
			Name:        "dbtype",
			Value:       node.DefaultDbType,
			Usage:       "Database backend to use for the Block Chain",
			Destination: nodeConfig.DbType,
		},
		cli.StringFlag{
			Name:        "profile",
			Usage:       "Enable HTTP profiling on given port -- NOTE port must be between 1024 and 65536",
			Destination: nodeConfig.Profile,
		},
		cli.StringFlag{
			Name:        "cpuprofile",
			Usage:       "Write CPU profile to the specified file",
			Destination: nodeConfig.CPUProfile,
		},
		cli.BoolFlag{
			Name:        "upnp",
			Usage:       "Use UPnP to map our listening port outside of NAT",
			Destination: nodeConfig.Upnp,
		},
		cli.Float64Flag{
			Name:        "minrelaytxfee",
			Value:       mempool.DefaultMinRelayTxFee.ToDUO(),
			Usage:       "The minimum transaction fee in DUO/kB to be considered a non-zero fee.",
			Destination: nodeConfig.MinRelayTxFee,
		},
		cli.Float64Flag{
			Name:        "limitfreerelay",
			Value:       node.DefaultFreeTxRelayLimit,
			Usage:       "Limit relay of transactions with no transaction fee to the given amount in thousands of bytes per minute",
			Destination: nodeConfig.FreeTxRelayLimit,
		},
		cli.BoolFlag{
			Name:        "norelaypriority",
			Usage:       "Do not require free or low-fee transactions to have high priority for relaying",
			Destination: nodeConfig.NoRelayPriority,
		},
		cli.DurationFlag{
			Name:        "trickleinterval",
			Value:       node.DefaultTrickleInterval,
			Usage:       "Minimum time between attempts to send new inventory to a connected peer",
			Destination: nodeConfig.TrickleInterval,
		},
		cli.IntFlag{
			Name:        "maxorphantx",
			Value:       node.DefaultMaxOrphanTransactions,
			Usage:       "Max number of orphan transactions to keep in memory",
			Destination: nodeConfig.MaxOrphanTxs,
		},
		cli.StringFlag{
			Name:        "algo",
			Value:       "random",
			Usage:       "Sets the algorithm for the CPU miner ( blake14lr, cryptonight7v2, keccak, lyra2rev2, scrypt, sha256d, stribog, skein, x11 default is 'random')",
			Destination: nodeConfig.Algo,
		},
		cli.BoolFlag{
			Name:        "generate",
			Usage:       "Generate (mine) DUO using the CPU",
			Destination: nodeConfig.Generate,
		},
		cli.IntFlag{
			Name:        "genthreads",
			Value:       -1,
			Usage:       "Number of CPU threads to use with CPU miner -1 = all cores",
			Destination: nodeConfig.GenThreads,
		},
		cli.StringSliceFlag{
			Name:  "miningaddr",
			Value: nodeConfig.MiningAddrs,
			Usage: "Add the specified payment address to the list of addresses to use for generated blocks, at least one is required if generate or minerlistener are set",
		},
		cli.StringFlag{
			Name:        "minerlistener",
			Usage:       "listen address for miner controller",
			Destination: nodeConfig.MinerListener,
		},
		cli.StringFlag{
			Name:        "minerpass",
			Value:       "pa55word",
			Usage:       "Encryption password required for miner clients to subscribe to work updates, for use over insecure connections",
			Destination: nodeConfig.MinerPass,
		},
		cli.IntFlag{
			Name:        "blockminsize",
			Value:       node.BlockMaxSizeMin,
			Usage:       "Mininum block size in bytes to be used when creating a block",
			Destination: nodeConfig.BlockMinSize,
		},
		cli.IntFlag{
			Name:        "blockmaxsize",
			Value:       node.BlockMaxSizeMax,
			Usage:       "Maximum block size in bytes to be used when creating a block",
			Destination: nodeConfig.BlockMaxSize,
		},
		cli.IntFlag{
			Name:        "blockminweight",
			Value:       node.BlockMaxWeightMin,
			Usage:       "Mininum block weight to be used when creating a block",
			Destination: nodeConfig.BlockMinWeight,
		},
		cli.IntFlag{
			Name:        "blockmaxweight",
			Value:       node.BlockMaxWeightMax,
			Usage:       "Maximum block weight to be used when creating a block",
			Destination: nodeConfig.BlockMaxWeight,
		},
		cli.IntFlag{
			Name:        "blockprioritysize",
			Usage:       "Size in bytes for high-priority/low-fee transactions when creating a block",
			Destination: nodeConfig.BlockPrioritySize,
		},
		cli.StringSliceFlag{
			Name:  "uacomment",
			Usage: "Comment to add to the user agent -- See BIP 14 for more information.",
			Value: nodeConfig.UserAgentComments,
		},
		cli.BoolFlag{
			Name:        "nopeerbloomfilters",
			Usage:       "Disable bloom filtering support",
			Destination: nodeConfig.NoPeerBloomFilters,
		},
		cli.BoolFlag{
			Name:        "nocfilters",
			Usage:       "Disable committed filtering (CF) support",
			Destination: nodeConfig.NoCFilters,
		},
		cli.IntFlag{
			Name:        "sigcachemaxsize",
			Value:       node.DefaultSigCacheMaxSize,
			Usage:       "The maximum number of entries in the signature verification cache",
			Destination: nodeConfig.SigCacheMaxSize,
		},
		cli.BoolFlag{
			Name:        "blocksonly",
			Usage:       "Do not accept transactions from remote peers.",
			Destination: nodeConfig.BlocksOnly,
		},
		cli.BoolTFlag{
			Name:        "notxindex",
			Usage:       "Disable the transaction index which makes all transactions available via the getrawtransaction RPC",
			Destination: nodeConfig.TxIndex,
		},
		cli.BoolTFlag{
			Name:        "noaddrindex",
			Usage:       "Disable address-based transaction index which makes the searchrawtransactions RPC available",
			Destination: nodeConfig.AddrIndex,
		},
		cli.BoolFlag{
			Name:        "relaynonstd",
			Usage:       "Relay non-standard transactions regardless of the default settings for the active network.",
			Destination: nodeConfig.RelayNonStd,
		},
		cli.BoolFlag{
			Name:        "rejectnonstd",
			Usage:       "Reject non-standard transactions regardless of the default settings for the active network.",
			Destination: nodeConfig.RejectNonStd,
		},
	},
}

var walletCommand = cli.Command{
	Name:    "wallet",
	Aliases: []string{"w"},
	Usage:   "start parallelcoin wallet server",
	Action:  walletHandle,
	Subcommands: []cli.Command{
		{
			Name:  "create",
			Usage: "Create the wallet if it does not exist",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
		{
			Name:  "createtemp",
			Usage: "Create a temporary simulation wallet (pass=password) in the data directory indicated; must call with --datadir",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
	},
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:        "noinitialload",
			Usage:       "Defer wallet creation/opening on startup and enable loading wallets over RPC",
			Destination: walletConfig.NoInitialLoad,
		},
		cli.StringFlag{
			Name:        "rpcconnect",
			Usage:       "Hostname/IP and port of pod RPC server to connect to (default localhost:11048, testnet: localhost:21048, simnet: localhost:41048)",
			Destination: walletConfig.RPCConnect,
		},
		cli.StringFlag{
			Name:        "walletpass",
			Usage:       "The public wallet password -- Only required if the wallet was created with one",
			Destination: walletConfig.WalletPass,
		},
		cli.StringFlag{
			Name:        "podusername",
			Value:       "user",
			Usage:       "Username for pod authentication",
			Destination: walletConfig.PodUsername,
		},
		cli.StringFlag{
			Name:        "podpassword",
			Value:       "pa55word",
			Usage:       "Password for pod authentication",
			Destination: walletConfig.PodPassword,
		},
		cli.BoolFlag{
			Name:        "onetimetlskey",
			Usage:       "Generate a new TLS certpair at startup, but only write the certificate to disk",
			Destination: walletConfig.OneTimeTLSKey,
		},
		cli.StringSliceFlag{
			Name:  "rpclisten",
			Usage: "Listen for legacy RPC connections on this interface/port (default port: 11046, testnet: 21046, simnet: 41046)",
			Value: walletConfig.LegacyRPCListeners,
		},
		cli.IntFlag{
			Name:        "rpcmaxclients",
			Value:       8,
			Usage:       "Max number of legacy RPC clients for standard connections",
			Destination: walletConfig.LegacyRPCMaxClients,
		},
		cli.IntFlag{
			Name:        "rpcmaxwebsockets",
			Value:       8,
			Usage:       "Max number of legacy RPC websocket connections",
			Destination: walletConfig.LegacyRPCMaxWebsockets,
		},
		cli.StringFlag{
			Name:        "username",
			Value:       "user",
			Usage:       "Username for legacy RPC and pod authentication (if podusername is unset)",
			Destination: walletConfig.Username,
		},
		cli.StringFlag{
			Name:        "password",
			Value:       "pa55word",
			Usage:       "Password for legacy RPC and pod authentication (if podpassword is unset)",
			Destination: walletConfig.Password,
		},
		cli.StringFlag{
			Name:        "profile",
			Usage:       "Enable HTTP profiling on given port -- NOTE port must be between 1024 and 65536",
			Destination: walletConfig.Profile,
		},
		cli.StringSliceFlag{
			Name:  "experimentalrpclisten",
			Usage: "Listen for RPC connections on this interface/port",
			Value: walletConfig.ExperimentalRPCListeners,
		},
	},
}
