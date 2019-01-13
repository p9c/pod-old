package main

var sampleModConf = []byte(`[Application Options]

; ------------------------------------------------------------------------------
; Bitcoin wallet settings
; ------------------------------------------------------------------------------

; Use testnet (cannot be used with simnet=1).
; testnet=0

; Use simnet (cannot be used with testnet=1).
; simnet=0

; The directory to open and save wallet, transaction, and unspent transaction
; output files.  Two directories, "mainnet" and "testnet" are used in this
; directory for mainnet and testnet wallets, respectively.
; appdata=~/.mod


; ------------------------------------------------------------------------------
; RPC client settings
; ------------------------------------------------------------------------------

; Connect via a SOCKS5 proxy.  NOTE: Specifying a proxy will disable listening
; for incoming connections unless listen addresses are provided via the
; 'rpclisten' option.
; proxy=127.0.0.1:9050
; proxyuser=
; proxypass=

; The server and port used for pod websocket connections.
; rpcconnect=localhost:21048

; File containing root certificates to authenticate a TLS connections with pod
; cafile=~/.mod/mod.cert



; ------------------------------------------------------------------------------
; RPC server settings
; ------------------------------------------------------------------------------

; Enable TLS to full node RPC
; servertls=1

; Enable TLS for the clients of mod
; clienttls=1

; TLS certificate and key file locations
; rpccert=~/.mod/rpc.cert
; rpckey=~/.mod/rpc.key

; Enable one time TLS keys.  This option results in the process generating
; a new certificate pair each startup, writing only the certificate file
; to disk.  This is a more secure option for clients that only interact with
; a local wallet process where persistent certs are not needed.
;
; This option will error at startup if the key specified by the rpckey option
; already exists.
; onetimetlskey=0

; Specify the interfaces for the RPC server listen on.  One rpclisten address
; per line.  Multiple rpclisten options may be set in the same configuration,
; and each will be used to listen for connections.  NOTE: The default port is
; modified by some options such as 'testnet', so it is recommended to not
; specify a port and allow a proper default to be chosen unless you have a
; specific reason to do otherwise.
; rpclisten=                ; all interfaces on default port
; rpclisten=0.0.0.0         ; all ipv4 interfaces on default port
; rpclisten=::              ; all ipv6 interfaces on default port
; rpclisten=:11046           ; all interfaces on port 11046
; rpclisten=0.0.0.0:11046    ; all ipv4 interfaces on port 11046
; rpclisten=[::]:11046       ; all ipv6 interfaces on port 11046
; rpclisten=127.0.0.1:11046  ; only ipv4 localhost on port 11046 (this is a default)
; rpclisten=[::1]:11046      ; only ipv6 localhost on port 11046 (this is a default)
; rpclisten=127.0.0.1:8337  ; only ipv4 localhost on non-standard port 8337
; rpclisten=:8337           ; all interfaces on non-standard port 8337
; rpclisten=0.0.0.0:8337    ; all ipv4 interfaces on non-standard port 8337
; rpclisten=[::]:8337       ; all ipv6 interfaces on non-standard port 8337

; Legacy (Bitcoin Core-compatible) RPC listener addresses.  Addresses without a
; port specified use the same default port as the new server.  Listeners cannot
; be shared between both RPC servers.
;
; Adding any legacy RPC listen addresses disable all default rpclisten options.
; If both servers must run, all listen addresses must be manually specified for
; each.
; legacyrpclisten=



; ------------------------------------------------------------------------------
; RPC settings (both client and server)
; ------------------------------------------------------------------------------

; Username and password to authenticate to pod a RPC server and authenticate
; new client connections
; username=
; password=

; Alternative username and password for pod.  If set, these will be used
; instead of the username and password set above for authentication to a
; pod RPC server.
; podusername=
; podpassword=


; ------------------------------------------------------------------------------
; Debug
; ------------------------------------------------------------------------------

; Debug logging level.
; Valid options are {trace, debug, info, warn, error, critical}
; debuglevel=info

; The port used to listen for HTTP profile requests.  The profile server will   
; be disabled if this option is not specified.  The profile information can be
; accessed at http://localhost:<profileport>/debug/pprof once running.
; profile=6062
`)
