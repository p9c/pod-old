package app

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
)

var App = cli.NewApp()

func Main() int {
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
				Name:   "datadir",
				Value:  "~/.pod",
				Usage:  "sets the data directory base for a pod instance",
				EnvVar: "POD_DATADIR",
			},
			cli.StringFlag{
				Name:   "loglevel",
				Value:  "info",
				Usage:  "sets the base for all subsystem logging",
				EnvVar: "POD_LOGLEVEL",
			},
			cli.StringSliceFlag{
				Name:  "subsystems",
				Value: &cli.StringSlice{""},
				Usage: "sets individual subsystems log levels, use 'help' to list available with list syntax",
			},
			cli.StringFlag{
				Name:  "network",
				Value: "mainnet",
				Usage: "connect to mainnet/testnet3/simnet",
				// Destination: nil,
			},
			cli.StringFlag{
				Name:  "rpccert",
				Value: "",
				Usage: "File containing the certificate file",
				// Destination: nil,
			},
			cli.StringFlag{
				Name:  "rpckey",
				Value: "",
				Usage: "File containing the certificate key",
				// Destination: nil,
			},
			cli.StringFlag{
				Name:  "tls, clienttls",
				Value: "",
				Usage: "Enable TLS for client connections",
				// Destination: nil,
			},
			cli.StringFlag{
				Name:  "servertls",
				Value: "",
				Usage: "Enable TLS for server connections",
				// Destination: nil,
			},
			cli.BoolFlag{
				Name:  "useproxy, r",
				Usage: "use configured proxy",
			},
			cli.StringFlag{
				Name:  "proxy",
				Value: "",
				Usage: "Connect via SOCKS5 proxy (eg. 127.0.0.1:9050)",
				// Destination: nil,
			},
			cli.StringFlag{
				Name:  "proxyuser",
				Value: "",
				Usage: "Username for proxy server",
				// Destination: nil,
			},
			cli.StringFlag{
				Name:  "proxypass",
				Value: "",
				Usage: "Password for proxy server",
				// Destination: nil,
			},
			cli.StringFlag{
				Name:  "onion",
				Value: "",
				Usage: "Connect to tor hidden services via SOCKS5 proxy (eg. 127.0.0.1:9050)",
				// Destination: nil,
			},
			cli.StringFlag{
				Name:  "onionuser",
				Value: "",
				Usage: "Username for onion proxy server",
				// Destination: nil,
			},
			cli.StringFlag{
				Name:  "onionpass",
				Value: "",
				Usage: "Password for onion proxy server",
				// Destination: nil,
			},
			cli.StringFlag{
				Name:  "noonion",
				Value: "",
				Usage: "Disable connecting to tor hidden services",
				// Destination: nil,
			},
			cli.StringFlag{
				Name:  "torisolation",
				Value: "",
				Usage: "Enable Tor stream isolation by randomizing user credentials for each connection.",
				// Destination: nil,
			},
			cli.StringFlag{
				Name:  "cafile",
				Value: "~/.pod/cafile",
				Usage: "File containing root certificates to authenticate a TLS connections with pod",
				// Destination: nil,
			},
		},
		Commands: []cli.Command{
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "print version and exit",
				Action: func(c *cli.Context) error {
					fmt.Println(c.App.Version)
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
						Destination: &ctlConfig.RPCServer,
					},
					cli.StringFlag{
						Name:        "walletserver, ws",
						Value:       "localhost:11046",
						Usage:       "set wallet server to use",
						Destination: &ctlConfig.Wallet,
					},
					cli.StringFlag{
						Name:        "rpcusername, username, user, u",
						Value:       "user",
						Usage:       "set rpc username",
						Destination: &ctlConfig.RPCUser,
					},
					cli.StringFlag{
						Name:        "rpcpassword, password, pass, p",
						Value:       "pa55word",
						Usage:       "set rpc password",
						Destination: &ctlConfig.RPCPass,
					},
					cli.BoolFlag{
						Name:  "wallet, w",
						Usage: "use configured wallet rpc address",
					},
				},
			},
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
				Name:    "node",
				Aliases: []string{"n"},
				Usage:   "start parallelcoin full node",
				Subcommands: []cli.Command{
					{
						Name:    "listcommands",
						Aliases: []string{"list", "l"},
						Usage:   "list commands available at endpoint",
						Action:  ctlHandleList,
					},
				},
				Flags: []cli.Flag{
					cli.StringSliceFlag{
						Name:  "addpeer",
						Usage: "Add a peer to connect with at startup",
						// Destination: nil,
					},
					cli.StringSliceFlag{
						Name:  "connect",
						Usage: "Connect only to the specified peers at startup",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "nolisten",
						Value: "",
						Usage: "Disable listening for incoming connections -- NOTE: Listening is automatically disabled if the --connect or --proxy options are used without also specifying listen interfaces via --listen",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "listen",
						Value: "",
						Usage: "Add an interface/port to listen for connections (default all interfaces port: 11047, testnet: 21047)",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "maxpeers",
						Value: "",
						Usage: "Max number of inbound and outbound peers",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "nobanning",
						Value: "",
						Usage: "Disable banning of misbehaving peers",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "banduration",
						Value: "",
						Usage: "How long to ban misbehaving peers.  Valid time units are {s, m, h, d}.  Minimum 1 second",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "banthreshold",
						Value: "",
						Usage: "Maximum allowed ban score before disconnecting and banning misbehaving peers.",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "whitelist",
						Value: "",
						Usage: "Add an IP network or IP that will not be banned. (eg. 192.168.1.0/24 or ::1)",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "rpcuser",
						Value: "",
						Usage: "Username for RPC connections",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "rpcpass",
						Value: "",
						Usage: "Password for RPC connections",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "rpclimituser",
						Value: "",
						Usage: "Username for limited RPC connections",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "rpclimitpass",
						Value: "",
						Usage: "Password for limited RPC connections",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "rpclisten",
						Value: "",
						Usage: "Add an interface/port to listen for RPC connections (default port: 11048, testnet: 21048) gives sha256d block templates",
						// Destination: nil,
					},

					cli.StringFlag{
						Name:  "rpcmaxclients",
						Value: "",
						Usage: "Max number of RPC clients for standard connections",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "rpcmaxwebsockets",
						Value: "",
						Usage: "Max number of RPC websocket connections",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "rpcmaxconcurrentreqs",
						Value: "",
						Usage: "Max number of concurrent RPC requests that may be processed concurrently",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "rpcquirks",
						Value: "",
						Usage: "Mirror some JSON-RPC quirks of Bitcoin Core -- NOTE: Discouraged unless interoperability issues need to be worked around",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "norpc",
						Value: "",
						Usage: "Disable built-in RPC server -- NOTE: The RPC server is disabled by default if no rpcuser/rpcpass or rpclimituser/rpclimitpass is specified",
						// Destination: nil,
					},

					cli.StringFlag{
						Name:  "nodnsseed",
						Value: "",
						Usage: "Disable DNS seeding for peers",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "externalip",
						Value: "",
						Usage: "Add an ip to the list of local addresses we claim to listen on to peers",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "addcheckpoint",
						Value: "",
						Usage: "Add a custom checkpoint.  Format: '<height>:<hash>'",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "nocheckpoints",
						Value: "",
						Usage: "Disable built-in checkpoints.  Don't do this unless you know what you're doing.",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "dbtype",
						Value: "",
						Usage: "Database backend to use for the Block Chain",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "profile",
						Value: "",
						Usage: "Enable HTTP profiling on given port -- NOTE port must be between 1024 and 65536",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "cpuprofile",
						Value: "",
						Usage: "Write CPU profile to the specified file",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "upnp",
						Value: "",
						Usage: "Use UPnP to map our listening port outside of NAT",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "minrelaytxfee",
						Value: "",
						Usage: "The minimum transaction fee in DUO/kB to be considered a non-zero fee.",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "limitfreerelay",
						Value: "",
						Usage: "Limit relay of transactions with no transaction fee to the given amount in thousands of bytes per minute",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "norelaypriority",
						Value: "",
						Usage: "Do not require free or low-fee transactions to have high priority for relaying",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "trickleinterval",
						Value: "",
						Usage: "Minimum time between attempts to send new inventory to a connected peer",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "maxorphantx",
						Value: "",
						Usage: "Max number of orphan transactions to keep in memory",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "algo",
						Value: "",
						Usage: "Sets the algorithm for the CPU miner ( blake14lr, cryptonight7v2, keccak, lyra2rev2, scrypt, sha256d, stribog, skein, x11 default is 'random')",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "generate",
						Value: "",
						Usage: "Generate (mine) DUO using the CPU",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "genthreads",
						Value: "",
						Usage: "Number of CPU threads to use with CPU miner -1 = all cores",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "miningaddr",
						Value: "",
						Usage: "Add the specified payment address to the list of addresses to use for generated blocks, at least one is required if generate or minerport are set",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "minerlistener",
						Value: "",
						Usage: "listen address for miner controller",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "minerpass",
						Value: "",
						Usage: "Encryption password required for miner clients to subscribe to work updates, for use over insecure connections",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "blockminsize",
						Value: "",
						Usage: "Mininum block size in bytes to be used when creating a block",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "blockmaxsize",
						Value: "",
						Usage: "Maximum block size in bytes to be used when creating a block",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "blockminweight",
						Value: "",
						Usage: "Mininum block weight to be used when creating a block",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "blockmaxweight",
						Value: "",
						Usage: "Maximum block weight to be used when creating a block",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "blockprioritysize",
						Value: "",
						Usage: "Size in bytes for high-priority/low-fee transactions when creating a block",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "uacomment",
						Value: "",
						Usage: "Comment to add to the user agent -- See BIP 14 for more information.",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "nopeerbloomfilters",
						Value: "",
						Usage: "Disable bloom filtering support",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "nocfilters",
						Value: "",
						Usage: "Disable committed filtering (CF) support",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "dropcfindex",
						Value: "",
						Usage: "Deletes the index used for committed filtering (CF) support from the database on start up and then exits.",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "sigcachemaxsize",
						Value: "",
						Usage: "The maximum number of entries in the signature verification cache",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "blocksonly",
						Value: "",
						Usage: "Do not accept transactions from remote peers.",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "txindex",
						Value: "",
						Usage: "Maintain a full hash-based transaction index which makes all transactions available via the getrawtransaction RPC",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "droptxindex",
						Value: "",
						Usage: "Deletes the hash-based transaction index from the database on start up and then exits.",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "addrindex",
						Value: "",
						Usage: "Maintain a full address-based transaction index which makes the searchrawtransactions RPC available",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "dropaddrindex",
						Value: "",
						Usage: "Deletes the address-based transaction index from the database on start up and then exits.",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "relaynonstd",
						Value: "",
						Usage: "Relay non-standard transactions regardless of the default settings for the active network.",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "rejectnonstd",
						Value: "",
						Usage: "Reject non-standard transactions regardless of the default settings for the active network.",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "service",
						Value: "",
						Usage: "Service command {install, remove, start, stop}",
						// Destination: nil,
					},
				},
				Action: func(c *cli.Context) error {
					fmt.Println("calling node")
					return nil
				},
			},
			{
				Name:    "wallet",
				Aliases: []string{"w"},
				Usage:   "start parallelcoin wallet server",
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
						Name:  "noinitialload",
						Usage: "Defer wallet creation/opening on startup and enable loading wallets over RPC",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "rpcconnect",
						Value: "localhost:11048",
						Usage: "Hostname/IP and port of pod RPC server to connect to (default localhost:11048, testnet: localhost:21048, simnet: localhost:41048)",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "walletpass",
						Value: "public",
						Usage: "The public wallet password -- Only required if the wallet was created with one",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "podusername",
						Value: "user",
						Usage: "Username for pod authentication",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "podpassword",
						Value: "pa55word",
						Usage: "Password for pod authentication",
						// Destination: nil,
					},
					cli.BoolFlag{
						Name:  "onetimetlskey",
						Usage: "Generate a new TLS certpair at startup, but only write the certificate to disk",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "rpclisten",
						Value: "",
						Usage: "Listen for legacy RPC connections on this interface/port (default port: 11046, testnet: 21046, simnet: 41046)",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "rpcmaxclients",
						Value: "",
						Usage: "Max number of legacy RPC clients for standard connections",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "rpcmaxwebsockets",
						Value: "",
						Usage: "Max number of legacy RPC websocket connections",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "username",
						Value: "",
						Usage: "Username for legacy RPC and pod authentication (if podusername is unset)",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "password",
						Value: "",
						Usage: "Password for legacy RPC and pod authentication (if podpassword is unset)",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "profile",
						Value: "",
						Usage: "Enable HTTP profiling on given port -- NOTE port must be between 1024 and 65536",
						// Destination: nil,
					},
					cli.StringFlag{
						Name:  "experimentalrpclisten",
						Value: "",
						Usage: "Listen for RPC connections on this interface/port",
						// Destination: nil,
					},
				},
				Action: func(c *cli.Context) error {
					fmt.Println("calling wallet")
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
