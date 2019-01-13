/*
pod is a full-node Parallelcoin implementation written in Go.

The default options are sane for most users.  This means pod will work 'out of the box' for most users.  However, there are also a wide variety of flags that can be used to control it.

The following section provides a usage overview which enumerates the flags.  An interesting point to note is that the long form of all of these options (except -C) can be specified in a configuration file that is automatically parsed when pod starts up.  By default, the configuration file is located at ~/.pod/pod.conf on POSIX-style operating systems and %LOCALAPPDATA%\pod\pod.conf on Windows.  The -C (--configfile) flag, as shown below, can be used to override this location.

Usage:
  pod [OPTIONS]

Application Options:
  -V, --version               Display version information and exit
  -C, --configfile=           Path to configuration file (default: /home/loki/.pod/pod.conf)
  -b, --datadir=              Directory to store data (default: /home/loki/.pod/data)
      --logdir=               Directory to log output. (default: /home/loki/.pod/logs)
  -a, --addpeer=              Add a peer to connect with at startup
      --connect=              Connect only to the specified peers at startup
      --nolisten              Disable listening for incoming connections -- NOTE: Listening is automatically disabled if the --connect or --proxy options are used without also specifying listen interfaces via --listen
      --listen=               Add an interface/port to listen for connections (default all interfaces port: 11047, testnet: 21047)
      --maxpeers=             Max number of inbound and outbound peers (default: 125)
      --nobanning             Disable banning of misbehaving peers
      --banduration=          How long to ban misbehaving peers.  Valid time units are {s, m, h}.  Minimum 1 second (default: 24h0m0s)
      --banthreshold=         Maximum allowed ban score before disconnecting and banning misbehaving peers. (default: 100)
      --whitelist=            Add an IP network or IP that will not be banned. (eg. 192.168.1.0/24 or ::1)
  -u, --rpcuser=              Username for RPC connections
  -P, --rpcpass=              Password for RPC connections
      --rpclimituser=         Username for limited RPC connections
      --rpclimitpass=         Password for limited RPC connections
      --rpclisten=            Add an interface/port to listen for RPC connections (default port: 11048, testnet: 21048) gives sha256d block templates
      --blake14lrlisten=      Additional RPC port that delivers blake14lr versioned block templates
      --cryptonight7v2=      Additional RPC port that delivers cryptonight7v2 versioned block templates
      --keccaklisten=         Additional RPC port that delivers keccak versioned block templates
      --lyra2rev2listen=      Additional RPC port that delivers lyra2rev2 versioned block templates
      --scryptlisten=         Additional RPC port that delivers scrypt versioned block templates
      --striboglisten=           Additional RPC port that delivers stribog versioned block templates
      --skeinlisten=          Additional RPC port that delivers skein versioned block templates
      --x11listen=            Additional RPC port that delivers x11 versioned block templates
      --rpccert=              File containing the certificate file (default: /home/loki/.pod/rpc.cert)
      --rpckey=               File containing the certificate key (default: /home/loki/.pod/rpc.key)
      --rpcmaxclients=        Max number of RPC clients for standard connections (default: 10)
      --rpcmaxwebsockets=     Max number of RPC websocket connections (default: 25)
      --rpcmaxconcurrentreqs= Max number of concurrent RPC requests that may be processed concurrently (default: 20)
      --rpcquirks             Mirror some JSON-RPC quirks of Bitcoin Core -- NOTE: Discouraged unless interoperability issues need to be worked around
      --norpc                 Disable built-in RPC server -- NOTE: The RPC server is disabled by default if no rpcuser/rpcpass or rpclimituser/rpclimitpass is specified
      --tls                   Enable TLS for the RPC server
      --nodnsseed             Disable DNS seeding for peers
      --externalip=           Add an ip to the list of local addresses we claim to listen on to peers
      --proxy=                Connect via SOCKS5 proxy (eg. 127.0.0.1:9050)
      --proxyuser=            Username for proxy server
      --proxypass=            Password for proxy server
      --onion=                Connect to tor hidden services via SOCKS5 proxy (eg. 127.0.0.1:9050)
      --onionuser=            Username for onion proxy server
      --onionpass=            Password for onion proxy server
      --noonion               Disable connecting to tor hidden services
      --torisolation          Enable Tor stream isolation by randomizing user credentials for each connection.
      --testnet               Use the test network
      --regtest               Use the regression test network
      --simnet                Use the simulation test network
      --addcheckpoint=        Add a custom checkpoint.  Format: '<height>:<hash>'
      --nocheckpoints         Disable built-in checkpoints.  Don't do this unless you know what you're doing.
      --dbtype=               Database backend to use for the Block Chain (default: ffldb)
      --profile=              Enable HTTP profiling on given port -- NOTE port must be between 1024 and 65536
      --cpuprofile=           Write CPU profile to the specified file
  -d, --debuglevel=           Logging level for all subsystems {trace, debug, info, warn, error, critical} -- You may also specify <subsystem>=<level>,<subsystem2>=<level>,... to set the log level for individual subsystems -- Use show to list available
                              subsystems (default: info)
      --upnp                  Use UPnP to map our listening port outside of NAT
      --minrelaytxfee=        The minimum transaction fee in DUO/kB to be considered a non-zero fee. (default: 1e-05)
      --limitfreerelay=       Limit relay of transactions with no transaction fee to the given amount in thousands of bytes per minute (default: 15)
      --norelaypriority       Do not require free or low-fee transactions to have high priority for relaying
      --trickleinterval=      Minimum time between attempts to send new inventory to a connected peer (default: 10s)
      --maxorphantx=          Max number of orphan transactions to keep in memory (default: 100)
      --algo=                 Sets the algorithm for the CPU miner ( blake14lr, cryptonight7v2, keccak, lyra2rev2, sha256d, scrypt, stribog, skein, x11,default sha256d) (default: sha256d)
      --generate              Generate (mine) bitcoins using the CPU
      --genthreads=           Number of CPU threads to use with CPU miner -1 = all cores (default: 1)
      --miningaddr=           Add the specified payment address to the list of addresses to use for generated blocks -- At least one address is required if the generate option is set
      --blockminsize=         Mininum block size in bytes to be used when creating a block (default: 80)
      --blockmaxsize=         Maximum block size in bytes to be used when creating a block (default: 200000)
      --blockminweight=       Mininum block weight to be used when creating a block (default: 10)
      --blockmaxweight=       Maximum block weight to be used when creating a block (default: 3000000)
      --blockprioritysize=    Size in bytes for high-priority/low-fee transactions when creating a block (default: 50000)
      --uacomment=            Comment to add to the user agent -- See BIP 14 for more information.
      --nopeerbloomfilters    Disable bloom filtering support
      --nocfilters            Disable committed filtering (CF) support
      --dropcfindex           Deletes the index used for committed filtering (CF) support from the database on start up and then exits.
      --sigcachemaxsize=      The maximum number of entries in the signature verification cache (default: 100000)
      --blocksonly            Do not accept transactions from remote peers.
      --txindex               Maintain a full hash-based transaction index which makes all transactions available via the getrawtransaction RPC
      --droptxindex           Deletes the hash-based transaction index from the database on start up and then exits.
      --addrindex             Maintain a full address-based transaction index which makes the searchrawtransactions RPC available
      --dropaddrindex         Deletes the address-based transaction index from the database on start up and then exits.
      --relaynonstd           Relay non-standard transactions regardless of the default settings for the active network.
      --rejectnonstd          Reject non-standard transactions regardless of the default settings for the active network.

Help Options:
  -h, --help                  Show this help message
*/
package node
