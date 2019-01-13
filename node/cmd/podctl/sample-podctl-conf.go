package main
var samplePodCtlConf = `;;;  Directory to store data (default: ~/.pod/data)
;datadir=              
;;;  Directory to log output. (default: ~/.pod/logs)
;logdir=               
;;;  Add a peer to connect with at startup
;addpeer=              
;;;  Connect only to the specified peers at startup
;connect=              
;;;  Disable listening for incoming connections -- NOTE: Listening is automatically disabled if the --connect or --proxy options are used without also specifying listen interfaces via --listen
;nolisten              
;;;  Add an interface/port to listen for connections (default all interfaces port: 11047, testnet: 21047)
;listen=               
;;;  Max number of inbound and outbound peers (default: 125)
;maxpeers=             
;;;  Disable banning of misbehaving peers
;nobanning             
;;;  How long to ban misbehaving peers.  Valid time units are {s, m, h}.  Minimum 1 second (default: 24h0m0s)
;banduration=          
;;;  Maximum allowed ban score before disconnecting and banning misbehaving peers. (default: 100)
;banthreshold=         
;;;  Add an IP network or IP that will not be banned. (eg. 192.168.1.0/24 or ::1)
;whitelist=            
;;;  Username for RPC connections
;rpcuser=
;;;  Password for RPC connections
;rpcpass=
;;;  Username for limited RPC connections
;rpclimituser=         
;;;  Password for limited RPC connections
;rpclimitpass=         
;;;  Add an interface/port to listen for RPC connections (default port: 11048, testnet: 21048)
;rpclisten=            
;;;  File containing the certificate file (default: ~/.pod/rpc.cert)
;rpccert=              
;;;  File containing the certificate key (default: ~/.pod/rpc.key)
;rpckey=               
;;;  Max number of RPC clients for standard connections (default: 10)
;rpcmaxclients=        
;;;  Max number of RPC websocket connections (default: 25)
;rpcmaxwebsockets=     
;;;  Max number of concurrent RPC requests that may be processed concurrently (default: 20)
;rpcmaxconcurrentreqs= 
;;;  Mirror some JSON-RPC quirks of Bitcoin Core -- NOTE: Discouraged unless interoperability issues need to be worked around
;rpcquirks             
;;;  Disable built-in RPC server -- NOTE: The RPC server is disabled by default if no rpcuser/rpcpass or rpclimituser/rpclimitpass is specified
;norpc                 
;;;  Enable TLS for the RPC server
;tls=1                 
;;;  Disable DNS seeding for peers
;nodnsseed             
;;;  Add an ip to the list of local addresses we claim to listen on to peers
;externalip=           
;;;  Connect via SOCKS5 proxy (eg. 127.0.0.1:9050)
;proxy=                
;;;  Username for proxy server
;proxyuser=            
;;;  Password for proxy server
;proxypass=            
;;;  Connect to tor hidden services via SOCKS5 proxy (eg. 127.0.0.1:9050)
;onion=                
;;;  Username for onion proxy server
;onionuser=            
;;;  Password for onion proxy server
;onionpass=            
;;;  Disable connecting to tor hidden services
;noonion               
;;;  Enable Tor stream isolation by randomizing user credentials for each connection.
;torisolation          
;;;  Use the test network
;testnet               
;;;  Use the regression test network
;regtest               
;;;  Use the simulation test network
;simnet                
;;;  Add a custom checkpoint.  Format: '<height>:<hash>'
;addcheckpoint=        
;;;  Disable built-in checkpoints.  Don't do this unless you know what you're doing.
;nocheckpoints         
;;;  Database backend to use for the Block Chain (default: ffldb)
;dbtype=               
;;;  Enable HTTP profiling on given port -- NOTE port must be between 1024 and 65536
;profile=              
;;;  Write CPU profile to the specified file
;cpuprofile=           
;;;  Logging level for all subsystems {trace, debug, info, warn, error, critical} -- You may also specify <subsystem>=<level>,<subsystem2>=<level>,... to set the log level for individual subsystems -- Use show to list available subsystems (default: info)
;debuglevel=           
;;;  Use UPnP to map our listening port outside of NAT
;upnp                  
;;;  The minimum transaction fee in DUO/kB to be considered a non-zero fee. (default: 1e-05)
;minrelaytxfee=        
;;;  Limit relay of transactions with no transaction fee to the given amount in thousands of bytes per minute (default: 15)
;limitfreerelay=       
;;;  Do not require free or low-fee transactions to have high priority for relaying
;norelaypriority       
;;;  Minimum time between attempts to send new inventory to a connected peer (default: 10s)
;trickleinterval=      
;;;  Max number of orphan transactions to keep in memory (default: 100)
;maxorphantx=          
;;;  Sets the algorithm for the CPU miner (sha256d/scrypt default sha256d) (default: sha256d)
;algo=                 
;;;  Generate (mine) bitcoins using the CPU
;generate              
;;;  Add the specified payment address to the list of addresses to use for generated blocks -- At least one address is required if the generate option is set
;miningaddr=           
;;;  Mininum block size in bytes to be used when creating a block (default: 80)
;blockminsize=         
;;;  Maximum block size in bytes to be used when creating a block (default: 200000)
;blockmaxsize=         
;;;  Mininum block weight to be used when creating a block (default: 10)
;blockminweight=       
;;;  Maximum block weight to be used when creating a block (default: 3000000)
;blockmaxweight=       
;;;  Size in bytes for high-priority/low-fee transactions when creating a block (default: 50000)
;blockprioritysize=    
;;;  Comment to add to the user agent -- See BIP 14 for more information.
;uacomment=            
;;;  Disable bloom filtering support
;nopeerbloomfilters    
;;;  Disable committed filtering (CF) support
;nocfilters            
;;;  Deletes the index used for committed filtering (CF) support from the database on start up and then exits.
;dropcfindex           
;;;  The maximum number of entries in the signature verification cache (default: 100000)
;sigcachemaxsize=      
;;;  Do not accept transactions from remote peers.
;blocksonly            
;;;  Maintain a full hash-based transaction index which makes all transactions available via the getrawtransaction RPC
;txindex               
;;;  Deletes the hash-based transaction index from the database on start up and then exits.
;droptxindex           
;;;  Maintain a full address-based transaction index which makes the searchrawtransactions RPC available
;addrindex             
;;;  Deletes the address-based transaction index from the database on start up and then exits.
;dropaddrindex         
;;;  Relay non-standard transactions regardless of the default settings for the active network.
;relaynonstd           
;;;  Reject non-standard transactions regardless of the default settings for the active network.
;rejectnonstd          
`
