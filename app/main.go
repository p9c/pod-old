package app

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/urfave/cli.v1/altsrc"

	"git.parallelcoin.io/dev/pod/cmd/node"
	"git.parallelcoin.io/dev/pod/cmd/node/mempool"
	walletmain "git.parallelcoin.io/dev/pod/cmd/walletmain"
	"git.parallelcoin.io/dev/pod/pkg/util/cl"
	"gopkg.in/urfave/cli.v1"
)

func Main() int {

	App = GetApp()

	log <- cl.Debug{"running App"}

	e := App.Run(os.Args)

	if e != nil {

		fmt.Println("Pod ERROR:", e)
		return 1
	}
	return 0
}

func GetApp() (a *cli.App) {

	a = &cli.App{

		Name:        "pod",
		Version:     "v0.0.1",
		Description: "Parallelcoin Pod Suite -- All-in-one everything for Parallelcoin!",
		Copyright:   "Legacy portions derived from btcsuite/btcd under ISC licence. The remainder is already in your possession. Use it wisely.",
		Action: func(c *cli.Context) error {

			Configure()

			fmt.Println("no subcommand requested")
			if StateCfg.Save {
				podHandleSave()
			}
			cli.ShowAppHelpAndExit(c, 1)
			return nil
		},
		Before: func(c *cli.Context) error {

			if FileExists(*podConfig.ConfigFile) {

				inputSource, err := altsrc.NewTomlSourceFromFile(*podConfig.ConfigFile)

				if err != nil {
					fmt.Println("error -", err)
					panic(err)
				}
				return altsrc.ApplyInputSourceValues(c, inputSource, c.App.Flags)
			}
			return nil
		},

		Commands: []cli.Command{
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "print version and exit",

				Action: func(c *cli.Context) error {

					fmt.Println(c.App.Name, c.App.Version)
					return nil
				}},
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
			},
			{

				Name:    "node",
				Aliases: []string{"n"},
				Usage:   "start parallelcoin full node",
				Action:  nodeHandle,
				Subcommands: []cli.Command{
					{
						Name:  "dropaddrindex",
						Usage: "drop the address search index",
						Action: func(c *cli.Context) error {
							StateCfg.DropAddrIndex = true
							return nodeHandle(c)
						},
					},
					{
						Name:  "droptxindex",
						Usage: "drop the address search index",
						Action: func(c *cli.Context) error {
							StateCfg.DropTxIndex = true
							return nodeHandle(c)
						},
					},
					{
						Name:  "dropcfindex",
						Usage: "drop the address search index",
						Action: func(c *cli.Context) error {
							StateCfg.DropCfIndex = true
							return nodeHandle(c)
						},
					},
				},
			},
			{

				Name:    "wallet",
				Aliases: []string{"w"},
				Usage:   "start parallelcoin wallet server",
				Action:  walletHandle,
				Subcommands: []cli.Command{

					{
						Name:  "create",
						Usage: "Create the wallet if it does not exist",
						Action: func(c *cli.Context) error {

							Configure()
							if err := walletmain.CreateWallet(&podConfig, activeNetParams); err != nil {

								log <- cl.Error{"failed to create wallet", err}

								return err
							}

							return nil
						},
					},

					// {
					// 	Name:  "createtemp",
					// 	Usage: "Create a temporary simulation wallet (pass=password) in the data directory indicated; must call with --datadir",
					// 	Action: func(c *cli.Context) error {
					// 		Configure()
					// 		if err := walletmain.CreateWallet(&podConfig, activeNetParams); err != nil {

					// 			log <- cl.Error{"failed to create wallet", err}
					// 			return err
					// 		}
					// 		return nil
					// 	},
					// },

				},
			},
			{
				Name:    "conf",
				Aliases: []string{"C"},
				Usage:   "populate all of the initial default configuration of a new data directory, all set globals will also apply. Exits after saving",
				Action:  confHandle,
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "base, b",
						Usage: "base name to extend with two number characters for testnet configurations",
						Value: "./test",
					}, cli.IntFlag{
						Name:  "number, n",
						Usage: "number of test? profiles to make based ",
					}},
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
		},
		Flags: []cli.Flag{altsrc.NewStringFlag(cli.StringFlag{
			Name:        "datadir, D",
			Value:       DefaultDataDir,
			Usage:       "sets the data directory base for a pod instance",
			EnvVar:      "POD_DATADIR",
			Destination: podConfig.DataDir,
		}), cli.BoolFlag{
			Name:        "save, i",
			Usage:       "save settings as effective from invocation",
			Destination: &StateCfg.Save,
		}, altsrc.NewStringFlag(cli.StringFlag{
			Name:        "loglevel, l",
			Value:       "info",
			Usage:       "sets the base for all subsystem logging",
			EnvVar:      "POD_LOGLEVEL",
			Destination: podConfig.LogLevel,
		}), altsrc.NewStringSliceFlag(cli.StringSliceFlag{
			Name:  "subsystem",
			Usage: "sets individual subsystems log levels, use 'listsubsystems' to list available",
			Value: podConfig.Subsystems,
		}), altsrc.NewStringFlag(cli.StringFlag{
			Name:        "network, n",
			Value:       "mainnet",
			Usage:       "connect to mainnet/testnet3/simnet",
			Destination: podConfig.Network,
		}), altsrc.NewStringFlag(cli.StringFlag{
			Name:        "username",
			Value:       "server",
			Usage:       "sets the username for services",
			Destination: podConfig.Username,
		}), altsrc.NewStringFlag(cli.StringFlag{
			Name:        "password",
			Value:       "pa55word",
			Usage:       "sets the password for services",
			Destination: podConfig.Password,
		}), altsrc.NewStringFlag(cli.StringFlag{
			Name:        "serveruser",
			Value:       "client",
			Usage:       "sets the username for clients of services",
			Destination: podConfig.ServerUser,
		}), altsrc.NewStringFlag(cli.StringFlag{
			Name:        "serverpass",
			Usage:       "sets the password for clients of services",
			Destination: podConfig.ServerPass,
		}), altsrc.NewStringFlag(cli.StringFlag{
			Name:        "limituser",
			Value:       "limit",
			Usage:       "sets the limited rpc username",
			Destination: podConfig.LimitUser,
		}), altsrc.NewStringFlag(cli.StringFlag{
			Name:        "limitpass",
			Usage:       "sets the password for clients of services",
			Destination: podConfig.LimitPass,
		}), altsrc.NewStringFlag(cli.StringFlag{
			Name:        "rpccert",
			Usage:       "File containing the certificate file",
			Destination: podConfig.RPCCert,
		}), altsrc.NewStringFlag(cli.StringFlag{
			Name:        "rpckey",
			Usage:       "File containing the certificate key",
			Destination: podConfig.RPCKey,
		}), altsrc.NewStringFlag(cli.StringFlag{
			Name:        "cafile",
			Value:       filepath.Join(DefaultDataDir, "cafile"),
			Usage:       "File containing root certificates to authenticate a TLS connections with pod",
			Destination: podConfig.CAFile,
		}), altsrc.NewBoolFlag(cli.BoolFlag{
			Name:        "tls, clienttls",
			Usage:       "Enable TLS for client connections",
			Destination: podConfig.TLS,
		}), altsrc.NewBoolFlag(cli.BoolFlag{
			Name:        "servertls",
			Usage:       "Enable TLS for server connections",
			Destination: podConfig.ServerTLS,
		}), altsrc.NewStringFlag(cli.StringFlag{
			Name:        "proxy",
			Usage:       "Connect via SOCKS5 proxy",
			Destination: podConfig.Proxy,
		}), altsrc.NewStringFlag(cli.StringFlag{
			Name:        "proxyuser",
			Value:       "user",
			Usage:       "Username for proxy server",
			Destination: podConfig.ProxyUser,
		}), altsrc.NewStringFlag(cli.StringFlag{
			Name:        "proxypass",
			Value:       "pa55word",
			Usage:       "Password for proxy server",
			Destination: podConfig.ProxyPass,
		}), altsrc.NewBoolFlag(cli.BoolFlag{
			Name:        "onion",
			Usage:       "Enable connecting to tor hidden services",
			Destination: podConfig.Onion,
		}), altsrc.NewStringFlag(cli.StringFlag{
			Name:        "onionproxy",
			Value:       "127.0.0.1:9050",
			Usage:       "Connect to tor hidden services via SOCKS5 proxy (eg. 127.0.0.1:9050)",
			Destination: podConfig.OnionProxy,
		}), altsrc.NewStringFlag(cli.StringFlag{
			Name:        "onionuser",
			Value:       "user",
			Usage:       "Username for onion proxy server",
			Destination: podConfig.OnionProxyUser,
		}), altsrc.NewStringFlag(cli.StringFlag{
			Name:        "onionpass",
			Value:       "pa55word",
			Usage:       "Password for onion proxy server",
			Destination: podConfig.OnionProxyPass,
		}), altsrc.NewBoolFlag(cli.BoolFlag{
			Name:        "torisolation",
			Usage:       "Enable Tor stream isolation by randomizing user credentials for each connection.",
			Destination: podConfig.TorIsolation,
		}), altsrc.NewStringFlag(cli.StringFlag{
			Name:        "walletserver, ws",
			Usage:       "set wallet server to connect to",
			Destination: podConfig.Wallet,
		}), altsrc.NewStringSliceFlag(cli.StringSliceFlag{
			Name:  "addpeer",
			Value: podConfig.AddPeers,
			Usage: "Add a peer to connect with at startup",
		}), altsrc.NewStringSliceFlag(cli.StringSliceFlag{
			Name:  "connect",
			Value: podConfig.ConnectPeers,
			Usage: "Connect only to the specified peers at startup",
		}), altsrc.NewBoolFlag(cli.BoolFlag{
			Name:        "nolisten",
			Usage:       "Disable listening for incoming connections -- NOTE: Listening is automatically disabled if the --connect or --proxy options are used without also specifying listen interfaces via --listen",
			Destination: podConfig.DisableListen,
		}), altsrc.NewStringSliceFlag(cli.StringSliceFlag{
			Name:  "listen",
			Value: podConfig.Listeners,
			Usage: "Add an interface/port to listen for connections",
		}), altsrc.NewIntFlag(cli.IntFlag{
			Name:        "maxpeers",
			Value:       node.DefaultMaxPeers,
			Usage:       "Max number of inbound and outbound peers",
			Destination: podConfig.MaxPeers,
		}), altsrc.NewBoolFlag(cli.BoolFlag{
			Name:        "nobanning",
			Usage:       "Disable banning of misbehaving peers",
			Destination: podConfig.DisableBanning,
		}), altsrc.NewDurationFlag(cli.DurationFlag{
			Name:        "banduration",
			Value:       time.Hour * 24,
			Usage:       "How long to ban misbehaving peers",
			Destination: podConfig.BanDuration,
		}), altsrc.NewIntFlag(cli.IntFlag{
			Name:        "banthreshold",
			Value:       node.DefaultBanThreshold,
			Usage:       "Maximum allowed ban score before disconnecting and banning misbehaving peers.",
			Destination: podConfig.BanThreshold,
		}), altsrc.NewStringSliceFlag(cli.StringSliceFlag{
			Name:  "whitelist",
			Usage: "Add an IP network or IP that will not be banned. (eg. 192.168.1.0/24 or ::1)",
			Value: podConfig.Whitelists,
		}), altsrc.NewStringFlag(cli.StringFlag{
			Name:        "rpcconnect",
			Value:       "127.0.0.1:11048",
			Usage:       "Hostname/IP and port of pod RPC server to connect to (default 127.0.0.1:11048, testnet: 127.0.0.1:21048, simnet: 127.0.0.1:41048)",
			Destination: podConfig.RPCConnect,
		}), altsrc.NewStringSliceFlag(cli.StringSliceFlag{
			Name:  "rpclisten",
			Value: podConfig.RPCListeners,
			Usage: "Add an interface/port to listen for RPC connections (default port: 11048, testnet: 21048) gives sha256d block templates",
		}), altsrc.NewIntFlag(cli.IntFlag{
			Name:        "rpcmaxclients",
			Value:       node.DefaultMaxRPCClients,
			Usage:       "Max number of RPC clients for standard connections",
			Destination: podConfig.RPCMaxClients,
		}), altsrc.NewIntFlag(cli.IntFlag{
			Name:        "rpcmaxwebsockets",
			Value:       node.DefaultMaxRPCWebsockets,
			Usage:       "Max number of RPC websocket connections",
			Destination: podConfig.RPCMaxWebsockets,
		}), altsrc.NewIntFlag(cli.IntFlag{
			Name:        "rpcmaxconcurrentreqs",
			Value:       node.DefaultMaxRPCConcurrentReqs,
			Usage:       "Max number of concurrent RPC requests that may be processed concurrently",
			Destination: podConfig.RPCMaxConcurrentReqs,
		}), altsrc.NewBoolFlag(cli.BoolFlag{
			Name:        "rpcquirks",
			Usage:       "Mirror some JSON-RPC quirks of Bitcoin Core -- NOTE: Discouraged unless interoperability issues need to be worked around",
			Destination: podConfig.RPCQuirks,
		}), altsrc.NewBoolFlag(cli.BoolFlag{
			Name:        "norpc",
			Usage:       "Disable built-in RPC server -- NOTE: The RPC server is disabled by default if no rpcuser/rpcpass or rpclimituser/rpclimitpass is specified",
			Destination: podConfig.DisableRPC,
		}), altsrc.NewBoolFlag(cli.BoolFlag{
			Name:        "nodnsseed",
			Usage:       "Disable DNS seeding for peers",
			Destination: podConfig.DisableDNSSeed,
		}), altsrc.NewStringSliceFlag(cli.StringSliceFlag{
			Name:  "externalip",
			Value: podConfig.ExternalIPs,
			Usage: "Add an ip to the list of local addresses we claim to listen on to peers",
		}), altsrc.NewStringSliceFlag(cli.StringSliceFlag{
			Name:  "addcheckpoint",
			Value: podConfig.AddCheckpoints,
			Usage: "Add a custom checkpoint.  Format: '<height>:<hash>'",
		}), altsrc.NewBoolFlag(cli.BoolFlag{
			Name:        "nocheckpoints",
			Usage:       "Disable built-in checkpoints.  Don't do this unless you know what you're doing.",
			Destination: podConfig.DisableCheckpoints,
		}), altsrc.NewStringFlag(cli.StringFlag{
			Name:        "dbtype",
			Value:       node.DefaultDbType,
			Usage:       "Database backend to use for the Block Chain",
			Destination: podConfig.DbType,
		}), altsrc.NewStringFlag(cli.StringFlag{
			Name:        "profile",
			Usage:       "Enable HTTP profiling on given port -- NOTE port must be between 1024 and 65536",
			Destination: podConfig.Profile,
		}), altsrc.NewStringFlag(cli.StringFlag{
			Name:        "cpuprofile",
			Usage:       "Write CPU profile to the specified file",
			Destination: podConfig.CPUProfile,
		}), altsrc.NewBoolFlag(cli.BoolFlag{
			Name:        "upnp",
			Usage:       "Use UPnP to map our listening port outside of NAT",
			Destination: podConfig.Upnp,
		}), altsrc.NewFloat64Flag(cli.Float64Flag{
			Name:        "minrelaytxfee",
			Value:       mempool.DefaultMinRelayTxFee.ToDUO(),
			Usage:       "The minimum transaction fee in DUO/kB to be considered a non-zero fee.",
			Destination: podConfig.MinRelayTxFee,
		}), altsrc.NewFloat64Flag(cli.Float64Flag{
			Name:        "limitfreerelay",
			Value:       node.DefaultFreeTxRelayLimit,
			Usage:       "Limit relay of transactions with no transaction fee to the given amount in thousands of bytes per minute",
			Destination: podConfig.FreeTxRelayLimit,
		}), altsrc.NewBoolFlag(cli.BoolFlag{
			Name:        "norelaypriority",
			Usage:       "Do not require free or low-fee transactions to have high priority for relaying",
			Destination: podConfig.NoRelayPriority,
		}), altsrc.NewDurationFlag(cli.DurationFlag{
			Name:        "trickleinterval",
			Value:       node.DefaultTrickleInterval,
			Usage:       "Minimum time between attempts to send new inventory to a connected peer",
			Destination: podConfig.TrickleInterval,
		}), altsrc.NewIntFlag(cli.IntFlag{
			Name:        "maxorphantx",
			Value:       node.DefaultMaxOrphanTransactions,
			Usage:       "Max number of orphan transactions to keep in memory",
			Destination: podConfig.MaxOrphanTxs,
		}), altsrc.NewStringFlag(cli.StringFlag{
			Name:        "algo",
			Value:       "random",
			Usage:       "Sets the algorithm for the CPU miner ( blake14lr, cryptonight7v2, keccak, lyra2rev2, scrypt, sha256d, stribog, skein, x11 default is 'random')",
			Destination: podConfig.Algo,
		}), altsrc.NewBoolFlag(cli.BoolFlag{
			Name:        "generate",
			Usage:       "Generate (mine) DUO using the CPU",
			Destination: podConfig.Generate,
		}), altsrc.NewIntFlag(cli.IntFlag{
			Name:        "genthreads",
			Value:       -1,
			Usage:       "Number of CPU threads to use with CPU miner -1 = all cores",
			Destination: podConfig.GenThreads,
		}), altsrc.NewStringSliceFlag(cli.StringSliceFlag{
			Name:  "miningaddr",
			Value: podConfig.MiningAddrs,
			Usage: "Add the specified payment address to the list of addresses to use for generated blocks, at least one is required if generate or minerlistener are set",
		}), altsrc.NewStringFlag(cli.StringFlag{
			Name:        "minerlistener",
			Usage:       "listen address for miner controller",
			Destination: podConfig.MinerListener,
		}), altsrc.NewStringFlag(cli.StringFlag{
			Name:        "minerpass",
			Usage:       "Encryption password required for miner clients to subscribe to work updates, for use over insecure connections",
			Destination: podConfig.MinerPass,
		}), altsrc.NewIntFlag(cli.IntFlag{
			Name:        "blockminsize",
			Value:       node.BlockMaxSizeMin,
			Usage:       "Mininum block size in bytes to be used when creating a block",
			Destination: podConfig.BlockMinSize,
		}), altsrc.NewIntFlag(cli.IntFlag{
			Name:        "blockmaxsize",
			Value:       node.BlockMaxSizeMax,
			Usage:       "Maximum block size in bytes to be used when creating a block",
			Destination: podConfig.BlockMaxSize,
		}), altsrc.NewIntFlag(cli.IntFlag{
			Name:        "blockminweight",
			Value:       node.BlockMaxWeightMin,
			Usage:       "Mininum block weight to be used when creating a block",
			Destination: podConfig.BlockMinWeight,
		}), altsrc.NewIntFlag(cli.IntFlag{
			Name:        "blockmaxweight",
			Value:       node.BlockMaxWeightMax,
			Usage:       "Maximum block weight to be used when creating a block",
			Destination: podConfig.BlockMaxWeight,
		}), altsrc.NewIntFlag(cli.IntFlag{
			Name:        "blockprioritysize",
			Usage:       "Size in bytes for high-priority/low-fee transactions when creating a block",
			Destination: podConfig.BlockPrioritySize,
		}), altsrc.NewStringSliceFlag(cli.StringSliceFlag{
			Name:  "uacomment",
			Usage: "Comment to add to the user agent -- See BIP 14 for more information.",
			Value: podConfig.UserAgentComments,
		}), altsrc.NewBoolFlag(cli.BoolFlag{
			Name:        "nopeerbloomfilters",
			Usage:       "Disable bloom filtering support",
			Destination: podConfig.NoPeerBloomFilters,
		}), altsrc.NewBoolFlag(cli.BoolFlag{
			Name:        "nocfilters",
			Usage:       "Disable committed filtering (CF) support",
			Destination: podConfig.NoCFilters,
		}), altsrc.NewIntFlag(cli.IntFlag{
			Name:        "sigcachemaxsize",
			Value:       node.DefaultSigCacheMaxSize,
			Usage:       "The maximum number of entries in the signature verification cache",
			Destination: podConfig.SigCacheMaxSize,
		}), altsrc.NewBoolFlag(cli.BoolFlag{
			Name:        "blocksonly",
			Usage:       "Do not accept transactions from remote peers.",
			Destination: podConfig.BlocksOnly,
		}), altsrc.NewBoolTFlag(cli.BoolTFlag{
			Name:        "notxindex",
			Usage:       "Disable the transaction index which makes all transactions available via the getrawtransaction RPC",
			Destination: podConfig.TxIndex,
		}), altsrc.NewBoolTFlag(cli.BoolTFlag{
			Name:        "noaddrindex",
			Usage:       "Disable address-based transaction index which makes the searchrawtransactions RPC available",
			Destination: podConfig.AddrIndex,
		}), altsrc.NewBoolFlag(cli.BoolFlag{
			Name:        "relaynonstd",
			Usage:       "Relay non-standard transactions regardless of the default settings for the active network.",
			Destination: podConfig.RelayNonStd,
		}), altsrc.NewBoolFlag(cli.BoolFlag{
			Name:        "rejectnonstd",
			Usage:       "Reject non-standard transactions regardless of the default settings for the active network.",
			Destination: podConfig.RejectNonStd,
		}), altsrc.NewBoolFlag(cli.BoolFlag{
			Name:        "noinitialload",
			Usage:       "Defer wallet creation/opening on startup and enable loading wallets over RPC",
			Destination: podConfig.NoInitialLoad,
		}), altsrc.NewStringFlag(cli.StringFlag{
			Name:        "walletpass",
			Usage:       "The public wallet password -- Only required if the wallet was created with one",
			Destination: podConfig.WalletPass,
		}), altsrc.NewBoolFlag(cli.BoolFlag{
			Name:        "onetimetlskey",
			Usage:       "Generate a new TLS certpair at startup, but only write the certificate to disk",
			Destination: podConfig.OneTimeTLSKey,
		}), altsrc.NewStringSliceFlag(cli.StringSliceFlag{
			Name:  "walletrpclisten",
			Usage: "Listen for wallet RPC connections on this interface/port (default port: 11046, testnet: 21046, simnet: 41046)",
			Value: podConfig.LegacyRPCListeners,
		}), altsrc.NewIntFlag(cli.IntFlag{
			Name:        "walletrpcmaxclients",
			Value:       8,
			Usage:       "Max number of legacy RPC clients for standard connections",
			Destination: podConfig.LegacyRPCMaxClients,
		}), altsrc.NewIntFlag(cli.IntFlag{
			Name:        "walletrpcmaxwebsockets",
			Value:       8,
			Usage:       "Max number of legacy RPC websocket connections",
			Destination: podConfig.LegacyRPCMaxWebsockets,
		}), altsrc.NewStringSliceFlag(cli.StringSliceFlag{
			Name:  "experimentalrpclisten",
			Usage: "Listen for RPC connections on this interface/port",
			Value: podConfig.ExperimentalRPCListeners,
		}),
		},
	}
	return
}
