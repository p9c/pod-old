package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"git.parallelcoin.io/pod/cmd/node"
	n "git.parallelcoin.io/pod/cmd/node"
	blockchain "git.parallelcoin.io/pod/pkg/chain"
	cl "git.parallelcoin.io/pod/pkg/clog"
	"git.parallelcoin.io/pod/pkg/connmgr"
	"git.parallelcoin.io/pod/pkg/fork"
	"git.parallelcoin.io/pod/pkg/netparams"
	"git.parallelcoin.io/pod/pkg/util"
	"github.com/btcsuite/go-socks/socks"
	"github.com/tucnak/climax"
)

func configShell(ctx *climax.Context, cfgFile string) int {
	ShellConfig.Wallet.AppDataDir = ShellConfig.Wallet.DataDir
	if r, ok := getIfIs(ctx, "appdatadir"); ok {
		ShellConfig.Wallet.AppDataDir = r
	}
	ShellConfig.SetNodeActiveNet(&node.MainNetParams)
	ShellConfig.SetWalletActiveNet(&netparams.MainNetParams)
	var ok bool
	var r string
	if ShellConfig.Node.TestNet3 {
		r = "testnet"
		ShellConfig.SetNodeActiveNet(&node.TestNet3Params)
		ShellConfig.SetWalletActiveNet(&netparams.TestNet3Params)
	}
	if ShellConfig.Node.SimNet {
		r = "simnet"
		ShellConfig.SetNodeActiveNet(&node.SimNetParams)
		ShellConfig.SetWalletActiveNet(&netparams.SimNetParams)
	}
	fmt.Println("nodeActiveNet.Name", r)
	if r, ok = getIfIs(ctx, "network"); ok {
		switch r {
		case "testnet":
			fork.IsTestnet = true
			ShellConfig.Wallet.TestNet3, ShellConfig.Wallet.SimNet = true, false
			ShellConfig.Node.TestNet3, ShellConfig.Node.SimNet, ShellConfig.Node.RegressionTest = true, false, false
			ShellConfig.SetNodeActiveNet(&node.TestNet3Params)
			ShellConfig.SetWalletActiveNet(&netparams.TestNet3Params)
		case "simnet":
			ShellConfig.Wallet.TestNet3, ShellConfig.Wallet.SimNet = false, true
			ShellConfig.Node.TestNet3, ShellConfig.Node.SimNet, ShellConfig.Node.RegressionTest = false, true, false
			ShellConfig.SetNodeActiveNet(&node.SimNetParams)
			ShellConfig.SetWalletActiveNet(&netparams.SimNetParams)
		default:
			ShellConfig.Wallet.TestNet3, ShellConfig.Wallet.SimNet = false, false
			ShellConfig.Node.TestNet3, ShellConfig.Node.SimNet, ShellConfig.Node.RegressionTest = false, false, false
			ShellConfig.SetNodeActiveNet(&node.MainNetParams)
			ShellConfig.SetWalletActiveNet(&netparams.MainNetParams)
		}
	}

	if ctx.Is("createtemp") {
		ShellConfig.Wallet.CreateTemp = true
	}
	if r, ok := getIfIs(ctx, "walletpass"); ok {
		ShellConfig.Wallet.WalletPass = r
	}
	if r, ok := getIfIs(ctx, "listeners"); ok {
		NormalizeAddresses(
			r, ShellConfig.GetNodeActiveNet().DefaultPort,
			&ShellConfig.Node.Listeners)
		log <- cl.Debug{"node listeners", ShellConfig.Node.Listeners}
	}
	if r, ok := getIfIs(ctx, "externalips"); ok {
		NormalizeAddresses(
			r, ShellConfig.GetNodeActiveNet().DefaultPort,
			&ShellConfig.Node.ExternalIPs)
		log <- cl.Debug{ShellConfig.Node.Listeners}
	}
	if r, ok := getIfIs(ctx, "disablelisten"); ok {
		ShellConfig.Node.DisableListen = strings.ToLower(r) == "true"
	}
	if r, ok := getIfIs(ctx, "rpclisteners"); ok {
		NormalizeAddresses(
			r, ShellConfig.GetWalletActiveNet().RPCServerPort,
			&ShellConfig.Wallet.LegacyRPCListeners)
	}
	if r, ok := getIfIs(ctx, "rpcmaxclients"); ok {
		var bt int
		if err := ParseInteger(r, "legacyrpcmaxclients", &bt); err != nil {
			log <- cl.Wrn(err.Error())
		} else {
			ShellConfig.Wallet.LegacyRPCMaxClients = int64(bt)
		}
	}
	if r, ok := getIfIs(ctx, "rpcmaxwebsockets"); ok {
		_, err := fmt.Sscanf(r, "%d", ShellConfig.Wallet.LegacyRPCMaxWebsockets)
		if err != nil {
			log <- cl.Errorf{
				"malformed legacyrpcmaxwebsockets: `%s` leaving set at `%d`",
				r, ShellConfig.Wallet.LegacyRPCMaxWebsockets,
			}
		}
	}
	if r, ok := getIfIs(ctx, "username"); ok {
		ShellConfig.Wallet.Username = r
		ShellConfig.Wallet.PodPassword = r
		ShellConfig.Node.RPCUser = r
	}
	if r, ok := getIfIs(ctx, "password"); ok {
		ShellConfig.Wallet.Password = r
		ShellConfig.Wallet.PodPassword = r
		ShellConfig.Node.RPCPass = r
	}
	if r, ok := getIfIs(ctx, "rpccert"); ok {
		ShellConfig.Wallet.RPCCert = n.CleanAndExpandPath(r)
		ShellConfig.Node.RPCCert = ShellConfig.Wallet.RPCCert
	}
	if r, ok := getIfIs(ctx, "rpckey"); ok {
		ShellConfig.Wallet.RPCKey = n.CleanAndExpandPath(r)
		ShellConfig.Node.RPCKey = ShellConfig.Wallet.RPCKey
	}
	if r, ok := getIfIs(ctx, "onetimetlskey"); ok {
		ShellConfig.Wallet.OneTimeTLSKey = strings.ToLower(r) == "true"
	}
	if r, ok := getIfIs(ctx, "cafile"); ok {
		ShellConfig.Wallet.CAFile = n.CleanAndExpandPath(r)
	}
	if r, ok := getIfIs(ctx, "tls"); ok {
		ShellConfig.Wallet.EnableServerTLS = strings.ToLower(r) == "true"
	}
	if r, ok := getIfIs(ctx, "txindex"); ok {
		ShellConfig.Node.TxIndex = strings.ToLower(r) == "true"
	}
	if r, ok := getIfIs(ctx, "addrindex"); ok {
		ShellConfig.Node.AddrIndex = strings.ToLower(r) == "true"
	}
	if ctx.Is("dropcfindex") {
		ShellConfig.Node.DropCfIndex = true
	}
	if ctx.Is("droptxindex") {
		ShellConfig.Node.DropTxIndex = true
	}
	if ctx.Is("dropaddrindex") {
		ShellConfig.Node.DropAddrIndex = true
	}
	if r, ok := getIfIs(ctx, "proxy"); ok {
		NormalizeAddress(r, "9050", &ShellConfig.Node.Proxy)
		ShellConfig.Wallet.Proxy = ShellConfig.Node.Proxy
	}
	if r, ok := getIfIs(ctx, "proxyuser"); ok {
		ShellConfig.Node.ProxyUser = r
		ShellConfig.Wallet.ProxyUser = r
	}
	if r, ok := getIfIs(ctx, "proxypass"); ok {
		ShellConfig.Node.ProxyPass = r
		ShellConfig.Node.ProxyPass = r
	}
	if r, ok := getIfIs(ctx, "onion"); ok {
		NormalizeAddress(r, "9050", &ShellConfig.Node.OnionProxy)
	}
	if r, ok := getIfIs(ctx, "onionuser"); ok {
		ShellConfig.Node.OnionProxyUser = r
	}
	if r, ok := getIfIs(ctx, "onionpass"); ok {
		ShellConfig.Node.OnionProxyPass = r
	}
	if r, ok := getIfIs(ctx, "noonion"); ok {
		ShellConfig.Node.NoOnion = r == "true"
	}
	if r, ok := getIfIs(ctx, "torisolation"); ok {
		ShellConfig.Node.TorIsolation = r == "true"
	}
	if r, ok := getIfIs(ctx, "addpeers"); ok {
		NormalizeAddresses(r, n.DefaultPort, &ShellConfig.Node.AddPeers)
	}
	if r, ok := getIfIs(ctx, "connectpeers"); ok {
		NormalizeAddresses(r, n.DefaultPort, &ShellConfig.Node.ConnectPeers)
	}
	if r, ok := getIfIs(ctx, "maxpeers"); ok {
		if err := ParseInteger(
			r, "maxpeers", &ShellConfig.Node.MaxPeers); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if r, ok := getIfIs(ctx, "disablebanning"); ok {
		ShellConfig.Node.DisableBanning = r == "true"
	}
	if r, ok := getIfIs(ctx, "banduration"); ok {
		if err := ParseDuration(r, "banduration", &ShellConfig.Node.BanDuration); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if r, ok := getIfIs(ctx, "banthreshold"); ok {
		var bt int
		if err := ParseInteger(r, "banthtreshold", &bt); err != nil {
			log <- cl.Wrn(err.Error())
		} else {
			ShellConfig.Node.BanThreshold = uint32(bt)
		}
	}
	if r, ok := getIfIs(ctx, "whitelists"); ok {
		NormalizeAddresses(r, n.DefaultPort, &ShellConfig.Node.Whitelists)
	}
	if r, ok := getIfIs(ctx, "trickleinterval"); ok {
		if err := ParseDuration(
			r, "trickleinterval", &ShellConfig.Node.TrickleInterval); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if r, ok := getIfIs(ctx, "minrelaytxfee"); ok {
		if err := ParseFloat(
			r, "minrelaytxfee", &ShellConfig.Node.MinRelayTxFee); err != nil {
			log <- cl.Wrn(err.Error())
		}

	}
	if r, ok := getIfIs(ctx, "freetxrelaylimit"); ok {
		if err := ParseFloat(
			r, "freetxrelaylimit", &ShellConfig.Node.FreeTxRelayLimit); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if r, ok := getIfIs(ctx, "norelaypriority"); ok {
		ShellConfig.Node.NoRelayPriority = r == "true"
	}
	if r, ok := getIfIs(ctx, "nopeerbloomfilters"); ok {
		ShellConfig.Node.NoPeerBloomFilters = r == "true"
	}
	if r, ok := getIfIs(ctx, "nocfilters"); ok {
		ShellConfig.Node.NoCFilters = r == "true"
	}
	if r, ok := getIfIs(ctx, "blocksonly"); ok {
		ShellConfig.Node.BlocksOnly = r == "true"
	}
	if r, ok := getIfIs(ctx, "relaynonstd"); ok {
		ShellConfig.Node.RelayNonStd = r == "true"
	}
	if r, ok := getIfIs(ctx, "rejectnonstd"); ok {
		ShellConfig.Node.RejectNonStd = r == "true"
	}
	if r, ok := getIfIs(ctx, "maxorphantxs"); ok {
		if err := ParseInteger(r, "maxorphantxs", &ShellConfig.Node.MaxOrphanTxs); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if r, ok := getIfIs(ctx, "sigcachemaxsize"); ok {
		var scms int
		if err := ParseInteger(r, "sigcachemaxsize", &scms); err != nil {
			log <- cl.Wrn(err.Error())
		} else {
			ShellConfig.Node.SigCacheMaxSize = uint(scms)
		}
	}
	if r, ok := getIfIs(ctx, "generate"); ok {
		ShellConfig.Node.Generate = r == "true"
	}
	if r, ok := getIfIs(ctx, "genthreads"); ok {
		var gt int
		if err := ParseInteger(r, "genthreads", &gt); err != nil {
			log <- cl.Wrn(err.Error())
		} else {
			ShellConfig.Node.GenThreads = int32(gt)
		}
	}
	if r, ok := getIfIs(ctx, "algo"); ok {
		ShellConfig.Node.Algo = r
	}
	if r, ok := getIfIs(ctx, "miningaddrs"); ok {
		ShellConfig.Node.MiningAddrs = strings.Split(r, " ")
	}
	if r, ok := getIfIs(ctx, "minerlistener"); ok {
		NormalizeAddress(r, n.DefaultRPCPort, &ShellConfig.Node.MinerListener)
	}
	if r, ok := getIfIs(ctx, "minerpass"); ok {
		ShellConfig.Node.MinerPass = r
	}
	if r, ok := getIfIs(ctx, "addcheckpoints"); ok {
		ShellConfig.Node.AddCheckpoints = strings.Split(r, " ")
	}
	if r, ok := getIfIs(ctx, "disablecheckpoints"); ok {
		ShellConfig.Node.DisableCheckpoints = r == "true"
	}
	if r, ok := getIfIs(ctx, "blockminsize"); ok {
		if err := ParseUint32(r, "blockminsize", &ShellConfig.Node.BlockMinSize); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if r, ok := getIfIs(ctx, "blockmaxsize"); ok {
		if err := ParseUint32(r, "blockmaxsize", &ShellConfig.Node.BlockMaxSize); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if r, ok := getIfIs(ctx, "blockminweight"); ok {
		if err := ParseUint32(r, "blockminweight", &ShellConfig.Node.BlockMinWeight); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if r, ok := getIfIs(ctx, "blockmaxweight"); ok {
		if err := ParseUint32(
			r, "blockmaxweight", &ShellConfig.Node.BlockMaxWeight); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if r, ok := getIfIs(ctx, "blockprioritysize"); ok {
		if err := ParseUint32(
			r, "blockmaxweight", &ShellConfig.Node.BlockPrioritySize); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if r, ok := getIfIs(ctx, "uacomment"); ok {
		ShellConfig.Node.UserAgentComments = strings.Split(r, " ")
	}
	if r, ok := getIfIs(ctx, "upnp"); ok {
		ShellConfig.Node.Upnp = r == "true"
	}
	if r, ok := getIfIs(ctx, "dbtype"); ok {
		ShellConfig.Node.DbType = r
	}
	if r, ok := getIfIs(ctx, "disablednsseed"); ok {
		ShellConfig.Node.DisableDNSSeed = r == "true"
	}
	if r, ok := getIfIs(ctx, "profile"); ok {
		var p int
		if err := ParseInteger(r, "profile", &p); err == nil {
			ShellConfig.Node.Profile = fmt.Sprint(p)
		}
	}
	if r, ok := getIfIs(ctx, "cpuprofile"); ok {
		ShellConfig.Node.CPUProfile = r
	}

	// finished configuration

	SetLogging(ctx)

	// Service options which are only added on Windows.
	serviceOpts := serviceOptions{}
	// Perform service command and exit if specified.  Invalid service commands show an appropriate error.  Only runs on Windows since the runServiceCommand function will be nil when not on Windows.
	if serviceOpts.ServiceCommand != "" && runServiceCommand != nil {
		err := runServiceCommand(serviceOpts.ServiceCommand)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		return 0
	}
	// Don't add peers from the config file when in regression test mode.
	if ShellConfig.Node.RegressionTest && len(ShellConfig.Node.AddPeers) > 0 {
		ShellConfig.Node.AddPeers = nil
	}
	// Set the mining algorithm correctly, default to random if unrecognised
	switch ShellConfig.Node.Algo {
	case "blake14lr", "cryptonight7v2", "keccak", "lyra2rev2", "scrypt", "skein", "x11", "stribog", "random", "easy":
	default:
		ShellConfig.Node.Algo = "random"
	}
	relayNonStd := n.ActiveNetParams.RelayNonStdTxs
	funcName := "loadConfig"
	switch {
	case ShellConfig.Node.RelayNonStd && ShellConfig.Node.RejectNonStd:
		str := "%s: rejectnonstd and relaynonstd cannot be used together -- choose only one"
		err := fmt.Errorf(str, funcName)
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, usageMessage)
		return 1
	case ShellConfig.Node.RejectNonStd:
		relayNonStd = false
	case ShellConfig.Node.RelayNonStd:
		relayNonStd = true
	}
	ShellConfig.Node.RelayNonStd = relayNonStd
	// Append the network type to the data directory so it is "namespaced" per network.  In addition to the block database, there are other pieces of data that are saved to disk such as address manager state. All data is specific to a network, so namespacing the data directory means each individual piece of serialized data does not have to worry about changing names per network and such.
	ShellConfig.Node.DataDir = n.CleanAndExpandPath(ShellConfig.Node.DataDir)
	ShellConfig.Node.DataDir = filepath.Join(ShellConfig.Node.DataDir, n.NetName(ShellConfig.GetNodeActiveNet()))
	// Append the network type to the log directory so it is "namespaced" per network in the same fashion as the data directory.
	ShellConfig.Node.LogDir = n.CleanAndExpandPath(ShellConfig.Node.LogDir)
	ShellConfig.Node.LogDir = filepath.Join(ShellConfig.Node.LogDir, n.NetName(ShellConfig.GetNodeActiveNet()))

	// Initialize log rotation.  After log rotation has been initialized, the logger variables may be used.
	// initLogRotator(filepath.Join(ShellConfig.Node.LogDir, DefaultLogFilename))
	// Validate database type.
	if !n.ValidDbType(ShellConfig.Node.DbType) {
		str := "%s: The specified database type [%v] is invalid -- supported types %v"
		err := fmt.Errorf(str, funcName, ShellConfig.Node.DbType, n.KnownDbTypes)
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, usageMessage)
		return 1
	}
	// Validate profile port number
	if ShellConfig.Node.Profile != "" {
		profilePort, err := strconv.Atoi(ShellConfig.Node.Profile)
		if err != nil || profilePort < 1024 || profilePort > 65535 {
			str := "%s: The profile port must be between 1024 and 65535"
			err := fmt.Errorf(str, funcName)
			fmt.Fprintln(os.Stderr, err)
			fmt.Fprintln(os.Stderr, usageMessage)
			return 1
		}
	}
	// Don't allow ban durations that are too short.
	if ShellConfig.Node.BanDuration < time.Second {
		str := "%s: The banduration option may not be less than 1s -- parsed [%v]"
		err := fmt.Errorf(str, funcName, ShellConfig.Node.BanDuration)
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, usageMessage)
		return 1
	}
	// Validate any given whitelisted IP addresses and networks.
	if len(ShellConfig.Node.Whitelists) > 0 {
		var ip net.IP
		StateCfg.ActiveWhitelists = make([]*net.IPNet, 0, len(ShellConfig.Node.Whitelists))
		for _, addr := range ShellConfig.Node.Whitelists {
			_, ipnet, err := net.ParseCIDR(addr)
			if err != nil {
				ip = net.ParseIP(addr)
				if ip == nil {
					str := "%s: The whitelist value of '%s' is invalid"
					err = fmt.Errorf(str, funcName, addr)
					log <- cl.Err(err.Error())
					fmt.Fprintln(os.Stderr, usageMessage)
					return 1
				}
				var bits int
				if ip.To4() == nil {
					// IPv6
					bits = 128
				} else {
					bits = 32
				}
				ipnet = &net.IPNet{
					IP:   ip,
					Mask: net.CIDRMask(bits, bits),
				}
			}
			StateCfg.ActiveWhitelists = append(StateCfg.ActiveWhitelists, ipnet)
		}
	}
	// --addPeer and --connect do not mix.
	if len(ShellConfig.Node.AddPeers) > 0 && len(ShellConfig.Node.ConnectPeers) > 0 {
		str := "%s: the --addpeer and --connect options can not be " +
			"mixed"
		err := fmt.Errorf(str, funcName)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
	}
	// --proxy or --connect without --listen disables listening.
	if (ShellConfig.Node.Proxy != "" || len(ShellConfig.Node.ConnectPeers) > 0) &&
		len(ShellConfig.Node.Listeners) == 0 {
		ShellConfig.Node.DisableListen = true
	}
	// Connect means no DNS seeding.
	if len(ShellConfig.Node.ConnectPeers) > 0 {
		ShellConfig.Node.DisableDNSSeed = true
	}
	// Add the default listener if none were specified. The default listener is all addresses on the listen port for the network we are to connect to.
	if len(ShellConfig.Node.Listeners) == 0 {
		ShellConfig.Node.Listeners = []string{
			net.JoinHostPort("localhost", ShellConfig.GetNodeActiveNet().DefaultPort),
		}
	}
	// Check to make sure limited and admin users don't have the same username
	if ShellConfig.Node.RPCUser == ShellConfig.Node.RPCLimitUser && ShellConfig.Node.RPCUser != "" {
		str := "%s: --rpcuser and --rpclimituser must not specify the same username"
		err := fmt.Errorf(str, funcName)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		return 1
	}
	// Check to make sure limited and admin users don't have the same password
	if ShellConfig.Node.RPCPass == ShellConfig.Node.RPCLimitPass && ShellConfig.Node.RPCPass != "" {
		str := "%s: --rpcpass and --rpclimitpass must not specify the same password"
		err := fmt.Errorf(str, funcName)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		return 1
	}
	// The RPC server is disabled if no username or password is provided.
	if (ShellConfig.Node.RPCUser == "" || ShellConfig.Node.RPCPass == "") &&
		(ShellConfig.Node.RPCLimitUser == "" || ShellConfig.Node.RPCLimitPass == "") {
		ShellConfig.Node.DisableRPC = true
	}
	if ShellConfig.Node.DisableRPC {
		log <- cl.Inf("RPC service is disabled")
	}
	// Default RPC to listen on localhost only.
	if !ShellConfig.Node.DisableRPC && len(ShellConfig.Node.RPCListeners) == 0 {
		addrs, err := net.LookupHost(n.DefaultRPCListener)
		if err != nil {
			log <- cl.Err(err.Error())
			return 1
		}
		ShellConfig.Node.RPCListeners = make([]string, 0, len(addrs))
		for _, addr := range addrs {
			addr = net.JoinHostPort(addr, n.ActiveNetParams.RPCPort)
			ShellConfig.Node.RPCListeners = append(ShellConfig.Node.RPCListeners, addr)
		}
	}
	if ShellConfig.Node.RPCMaxConcurrentReqs < 0 {
		str := "%s: The rpcmaxwebsocketconcurrentrequests option may not be less than 0 -- parsed [%d]"
		err := fmt.Errorf(str, funcName, ShellConfig.Node.RPCMaxConcurrentReqs)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		return 1
	}
	var err error
	// Validate the the minrelaytxfee.
	StateCfg.ActiveMinRelayTxFee, err = util.NewAmount(ShellConfig.Node.MinRelayTxFee)
	if err != nil {
		str := "%s: invalid minrelaytxfee: %v"
		err := fmt.Errorf(str, funcName, err)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		return 1
	}
	// Limit the max block size to a sane value.
	if ShellConfig.Node.BlockMaxSize < n.BlockMaxSizeMin || ShellConfig.Node.BlockMaxSize >
		n.BlockMaxSizeMax {
		str := "%s: The blockmaxsize option must be in between %d and %d -- parsed [%d]"
		err := fmt.Errorf(str, funcName, n.BlockMaxSizeMin,
			n.BlockMaxSizeMax, ShellConfig.Node.BlockMaxSize)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		return 1
	}
	// Limit the max block weight to a sane value.
	if ShellConfig.Node.BlockMaxWeight < n.BlockMaxWeightMin ||
		ShellConfig.Node.BlockMaxWeight > n.BlockMaxWeightMax {
		str := "%s: The blockmaxweight option must be in between %d and %d -- parsed [%d]"
		err := fmt.Errorf(str, funcName, n.BlockMaxWeightMin,
			n.BlockMaxWeightMax, ShellConfig.Node.BlockMaxWeight)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		return 1
	}
	// Limit the max orphan count to a sane vlue.
	if ShellConfig.Node.MaxOrphanTxs < 0 {
		str := "%s: The maxorphantx option may not be less than 0 -- parsed [%d]"
		err := fmt.Errorf(str, funcName, ShellConfig.Node.MaxOrphanTxs)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		return 1
	}
	// Limit the block priority and minimum block sizes to max block size.
	ShellConfig.Node.BlockPrioritySize = minUint32(ShellConfig.Node.BlockPrioritySize, ShellConfig.Node.BlockMaxSize)
	ShellConfig.Node.BlockMinSize = minUint32(ShellConfig.Node.BlockMinSize, ShellConfig.Node.BlockMaxSize)
	ShellConfig.Node.BlockMinWeight = minUint32(ShellConfig.Node.BlockMinWeight, ShellConfig.Node.BlockMaxWeight)
	switch {
	// If the max block size isn't set, but the max weight is, then we'll set the limit for the max block size to a safe limit so weight takes precedence.
	case ShellConfig.Node.BlockMaxSize == n.DefaultBlockMaxSize &&
		ShellConfig.Node.BlockMaxWeight != n.DefaultBlockMaxWeight:
		ShellConfig.Node.BlockMaxSize = blockchain.MaxBlockBaseSize - 1000
	// If the max block weight isn't set, but the block size is, then we'll scale the set weight accordingly based on the max block size value.
	case ShellConfig.Node.BlockMaxSize != n.DefaultBlockMaxSize &&
		ShellConfig.Node.BlockMaxWeight == n.DefaultBlockMaxWeight:
		ShellConfig.Node.BlockMaxWeight = ShellConfig.Node.BlockMaxSize * blockchain.WitnessScaleFactor
	}
	// Look for illegal characters in the user agent comments.
	for _, uaComment := range ShellConfig.Node.UserAgentComments {
		if strings.ContainsAny(uaComment, "/:()") {
			err := fmt.Errorf("%s: The following characters must not "+
				"appear in user agent comments: '/', ':', '(', ')'",
				funcName)
			log <- cl.Err(err.Error())
			fmt.Fprintln(os.Stderr, usageMessage)
			return 1

		}
	}
	// --txindex and --droptxindex do not mix.
	if ShellConfig.Node.TxIndex && ShellConfig.Node.DropTxIndex {
		err := fmt.Errorf("%s: the --txindex and --droptxindex options may  not be activated at the same time",
			funcName)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		return 1

	}
	// --addrindex and --dropaddrindex do not mix.
	if ShellConfig.Node.AddrIndex && ShellConfig.Node.DropAddrIndex {
		err := fmt.Errorf("%s: the --addrindex and --dropaddrindex "+
			"options may not be activated at the same time",
			funcName)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		return 1
	}
	// --addrindex and --droptxindex do not mix.
	if ShellConfig.Node.AddrIndex && ShellConfig.Node.DropTxIndex {
		err := fmt.Errorf("%s: the --addrindex and --droptxindex options may not be activated at the same time "+
			"because the address index relies on the transaction index",
			funcName)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		return 1
	}
	// Check mining addresses are valid and saved parsed versions.
	StateCfg.ActiveMiningAddrs = make([]util.Address, 0, len(ShellConfig.Node.MiningAddrs))
	for _, strAddr := range ShellConfig.Node.MiningAddrs {
		addr, err := util.DecodeAddress(strAddr, n.ActiveNetParams.Params)
		if err != nil {
			str := "%s: mining address '%s' failed to decode: %v"
			err := fmt.Errorf(str, funcName, strAddr, err)
			log <- cl.Err(err.Error())
			fmt.Fprintln(os.Stderr, usageMessage)
			return 1
		}
		if !addr.IsForNet(n.ActiveNetParams.Params) {
			str := "%s: mining address '%s' is on the wrong network"
			err := fmt.Errorf(str, funcName, strAddr)
			log <- cl.Err(err.Error())
			fmt.Fprintln(os.Stderr, usageMessage)
			return 1
		}
		StateCfg.ActiveMiningAddrs = append(StateCfg.ActiveMiningAddrs, addr)
	}
	// Ensure there is at least one mining address when the generate flag is set.
	if (ShellConfig.Node.Generate || ShellConfig.Node.MinerListener != "") && len(ShellConfig.Node.MiningAddrs) == 0 {
		str := "%s: the generate flag is set, but there are no mining addresses specified "
		err := fmt.Errorf(str, funcName)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		os.Exit(1)
	}
	if ShellConfig.Node.MinerPass != "" {
		StateCfg.ActiveMinerKey = fork.Argon2i([]byte(ShellConfig.Node.MinerPass))
	}
	// Add default port to all listener addresses if needed and remove duplicate addresses.
	ShellConfig.Node.Listeners = n.NormalizeAddresses(ShellConfig.Node.Listeners,
		ShellConfig.GetNodeActiveNet().DefaultPort)
	// Add default port to all rpc listener addresses if needed and remove duplicate addresses.
	ShellConfig.Node.RPCListeners = n.NormalizeAddresses(ShellConfig.Node.RPCListeners,
		ShellConfig.GetNodeActiveNet().RPCPort)
	if !ShellConfig.Node.DisableRPC && !ShellConfig.Node.TLS {
		for _, addr := range ShellConfig.Node.RPCListeners {
			if err != nil {
				str := "%s: RPC listen interface '%s' is invalid: %v"
				err := fmt.Errorf(str, funcName, addr, err)
				log <- cl.Err(err.Error())
				fmt.Fprintln(os.Stderr, usageMessage)
				return 1
			}
		}
	}
	// Add default port to all added peer addresses if needed and remove duplicate addresses.
	ShellConfig.Node.AddPeers = n.NormalizeAddresses(ShellConfig.Node.AddPeers,
		ShellConfig.GetNodeActiveNet().DefaultPort)
	ShellConfig.Node.ConnectPeers = n.NormalizeAddresses(ShellConfig.Node.ConnectPeers,
		ShellConfig.GetNodeActiveNet().DefaultPort)
	// --noonion and --onion do not mix.
	if ShellConfig.Node.NoOnion && ShellConfig.Node.OnionProxy != "" {
		err := fmt.Errorf("%s: the --noonion and --onion options may not be activated at the same time", funcName)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		return 1
	}
	// Check the checkpoints for syntax errors.
	StateCfg.AddedCheckpoints, err = n.ParseCheckpoints(ShellConfig.Node.AddCheckpoints)
	if err != nil {
		str := "%s: Error parsing checkpoints: %v"
		err := fmt.Errorf(str, funcName, err)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		return 1
	}
	// Tor stream isolation requires either proxy or onion proxy to be set.
	if ShellConfig.Node.TorIsolation && ShellConfig.Node.Proxy == "" && ShellConfig.Node.OnionProxy == "" {
		str := "%s: Tor stream isolation requires either proxy or onionproxy to be set"
		err := fmt.Errorf(str, funcName)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		return 1
	}
	// Setup dial and DNS resolution (lookup) functions depending on the specified options.  The default is to use the standard net.DialTimeout function as well as the system DNS resolver.  When a proxy is specified, the dial function is set to the proxy specific dial function and the lookup is set to use tor (unless --noonion is specified in which case the system DNS resolver is used).
	StateCfg.Dial = net.DialTimeout
	StateCfg.Lookup = net.LookupIP
	if ShellConfig.Node.Proxy != "" {
		_, _, err := net.SplitHostPort(ShellConfig.Node.Proxy)
		if err != nil {
			str := "%s: Proxy address '%s' is invalid: %v"
			err := fmt.Errorf(str, funcName, ShellConfig.Node.Proxy, err)
			log <- cl.Err(err.Error())
			fmt.Fprintln(os.Stderr, usageMessage)
			return 1
		}
		// Tor isolation flag means proxy credentials will be overridden unless there is also an onion proxy configured in which case that one will be overridden.
		torIsolation := false
		if ShellConfig.Node.TorIsolation && ShellConfig.Node.OnionProxy == "" &&
			(ShellConfig.Node.ProxyUser != "" || ShellConfig.Node.ProxyPass != "") {
			torIsolation = true
			fmt.Fprintln(os.Stderr, "Tor isolation set -- "+
				"overriding specified proxy user credentials")
		}
		proxy := &socks.Proxy{
			Addr:         ShellConfig.Node.Proxy,
			Username:     ShellConfig.Node.ProxyUser,
			Password:     ShellConfig.Node.ProxyPass,
			TorIsolation: torIsolation,
		}
		StateCfg.Dial = proxy.DialTimeout
		// Treat the proxy as tor and perform DNS resolution through it unless the --noonion flag is set or there is an onion-specific proxy configured.
		if !ShellConfig.Node.NoOnion && ShellConfig.Node.OnionProxy == "" {
			StateCfg.Lookup = func(host string) ([]net.IP, error) {
				return connmgr.TorLookupIP(host, ShellConfig.Node.Proxy)
			}
		}
	}
	// Setup onion address dial function depending on the specified options. The default is to use the same dial function selected above.  However, when an onion-specific proxy is specified, the onion address dial function is set to use the onion-specific proxy while leaving the normal dial function as selected above.  This allows .onion address traffic to be routed through a different proxy than normal traffic.
	if ShellConfig.Node.OnionProxy != "" {
		_, _, err := net.SplitHostPort(ShellConfig.Node.OnionProxy)
		if err != nil {
			str := "%s: Onion proxy address '%s' is invalid: %v"
			err := fmt.Errorf(str, funcName, ShellConfig.Node.OnionProxy, err)
			log <- cl.Err(err.Error())
			fmt.Fprintln(os.Stderr, usageMessage)
			return 1
		}
		// Tor isolation flag means onion proxy credentials will be overridden.
		if ShellConfig.Node.TorIsolation &&
			(ShellConfig.Node.OnionProxyUser != "" || ShellConfig.Node.OnionProxyPass != "") {
			fmt.Fprintln(os.Stderr, "Tor isolation set -- "+
				"overriding specified onionproxy user "+
				"credentials ")
		}
		StateCfg.Oniondial = func(network, addr string, timeout time.Duration) (net.Conn, error) {
			proxy := &socks.Proxy{
				Addr:         ShellConfig.Node.OnionProxy,
				Username:     ShellConfig.Node.OnionProxyUser,
				Password:     ShellConfig.Node.OnionProxyPass,
				TorIsolation: ShellConfig.Node.TorIsolation,
			}
			return proxy.DialTimeout(network, addr, timeout)
		}
		// When configured in bridge mode (both --onion and --proxy are configured), it means that the proxy configured by --proxy is not a tor proxy, so override the DNS resolution to use the onion-specific proxy.
		if ShellConfig.Node.Proxy != "" {
			StateCfg.Lookup = func(host string) ([]net.IP, error) {
				return connmgr.TorLookupIP(host, ShellConfig.Node.OnionProxy)
			}
		}
	} else {
		StateCfg.Oniondial = StateCfg.Dial
	}
	// Specifying --noonion means the onion address dial function results in an error.
	if ShellConfig.Node.NoOnion {
		StateCfg.Oniondial = func(a, b string, t time.Duration) (net.Conn, error) {
			return nil, errors.New("tor has been disabled")
		}
	}

	ShellConfig.Wallet.PodUsername = ShellConfig.Node.RPCUser
	ShellConfig.Wallet.PodPassword = ShellConfig.Node.RPCPass

	if ctx.Is("save") {
		log <- cl.Info{"saving config file to", cfgFile}
		j, err := json.MarshalIndent(ShellConfig, "", "  ")
		if err != nil {
			log <- cl.Error{"writing app config file", err}
		}
		j = append(j, '\n')
		log <- cl.Trace{"JSON formatted config file\n", string(j)}
		ioutil.WriteFile(cfgFile, j, 0600)
	}
	return 0
}
