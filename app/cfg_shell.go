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

func configShell(sc *ShellCfg, ctx *climax.Context, cfgFile string) {
	log <- cl.Trace{"configuring from command line flags ", os.Args}

	if ctx.Is("version") {
		//
	}
	if r, ok := getIfIs(ctx, "configfile"); ok {
		ShellConfig.ConfigFile = r
	}
	if r, ok := getIfIs(ctx, "datadir"); ok {
		ShellConfig.Node.DataDir = r
		ShellConfig.Wallet.DataDir = r
	}
	if r, ok := getIfIs(ctx, "appdatadir"); ok {
		ShellConfig.Wallet.AppDataDir = r
	}
	if ctx.Is("init") {

	}
	if r, ok := getIfIs(ctx, "network"); ok {
		switch r {
		case "testnet":
			sc.Wallet.TestNet3, sc.Wallet.SimNet = true, false
			sc.Node.TestNet3, sc.Node.SimNet, sc.Node.RegressionTest = true, false, false
			sc.nodeActiveNet = &node.TestNet3Params
			sc.walletActiveNet = &netparams.TestNet3Params
		case "simnet":
			sc.Wallet.TestNet3, sc.Wallet.SimNet = false, true
			sc.Node.TestNet3, sc.Node.SimNet, sc.Node.RegressionTest = false, true, false
			sc.nodeActiveNet = &node.SimNetParams
			sc.walletActiveNet = &netparams.SimNetParams
		default:
			sc.Wallet.TestNet3, sc.Wallet.SimNet = false, false
			sc.Node.TestNet3, sc.Node.SimNet, sc.Node.RegressionTest = false, false, false
			sc.nodeActiveNet = &node.MainNetParams
			sc.walletActiveNet = &netparams.MainNetParams
		}
	}

	if ctx.Is("createtemp") {
		sc.Wallet.CreateTemp = true
	}
	if r, ok := getIfIs(ctx, "walletpass"); ok {
		sc.Wallet.WalletPass = r
	}
	if r, ok := getIfIs(ctx, "listeners"); ok {
		NormalizeAddresses(r, "11047", &sc.Node.Listeners)
	}
	if r, ok := getIfIs(ctx, "externalips"); ok {
		NormalizeAddresses(r, "11047", &sc.Node.ExternalIPs)
	}
	if r, ok := getIfIs(ctx, "disablelisten"); ok {
		sc.Node.DisableListen = strings.ToLower(r) == "true"
	}
	if r, ok := getIfIs(ctx, "rpclisteners"); ok {
		NormalizeAddresses(r, "11046", &sc.Wallet.LegacyRPCListeners)
	}
	if r, ok := getIfIs(ctx, "rpcmaxclients"); ok {
		var bt int
		if err := ParseInteger(r, "legacyrpcmaxclients", &bt); err != nil {
			log <- cl.Wrn(err.Error())
		} else {
			sc.Wallet.LegacyRPCMaxClients = int64(bt)
		}
	}
	if r, ok := getIfIs(ctx, "rpcmaxwebsockets"); ok {
		_, err := fmt.Sscanf(r, "%d", sc.Wallet.LegacyRPCMaxWebsockets)
		if err != nil {
			log <- cl.Errorf{
				"malformed legacyrpcmaxwebsockets: `%s` leaving set at `%d`",
				r, sc.Wallet.LegacyRPCMaxWebsockets,
			}
		}
	}
	if r, ok := getIfIs(ctx, "username"); ok {
		sc.Wallet.Username = r
		sc.Wallet.PodPassword = r
		sc.Node.RPCUser = r
	}
	if r, ok := getIfIs(ctx, "password"); ok {
		sc.Wallet.Password = r
		sc.Wallet.PodPassword = r
		sc.Node.RPCPass = r
	}
	if r, ok := getIfIs(ctx, "rpccert"); ok {
		sc.Wallet.RPCCert = n.CleanAndExpandPath(r)
		sc.Node.RPCCert = sc.Wallet.RPCCert
	}
	if r, ok := getIfIs(ctx, "rpckey"); ok {
		sc.Wallet.RPCKey = n.CleanAndExpandPath(r)
		sc.Node.RPCKey = sc.Wallet.RPCKey
	}
	if r, ok := getIfIs(ctx, "onetimetlskey"); ok {
		sc.Wallet.OneTimeTLSKey = strings.ToLower(r) == "true"
	}
	if r, ok := getIfIs(ctx, "cafile"); ok {
		sc.Wallet.CAFile = n.CleanAndExpandPath(r)
	}
	if r, ok := getIfIs(ctx, "tls"); ok {
		sc.Wallet.EnableServerTLS = strings.ToLower(r) == "true"
	}
	if r, ok := getIfIs(ctx, "txindex"); ok {
		sc.Node.TxIndex = strings.ToLower(r) == "true"
	}
	if r, ok := getIfIs(ctx, "addrindex"); ok {
		sc.Node.AddrIndex = strings.ToLower(r) == "true"
	}
	if ctx.Is("dropcfindex") {
		sc.Node.DropCfIndex = true
	}
	if ctx.Is("droptxindex") {
		sc.Node.DropTxIndex = true
	}
	if ctx.Is("dropaddrindex") {
		sc.Node.DropAddrIndex = true
	}
	if r, ok := getIfIs(ctx, "proxy"); ok {
		NormalizeAddress(r, "11048", &sc.Node.Proxy)
		sc.Wallet.Proxy = sc.Node.Proxy
	}
	if r, ok := getIfIs(ctx, "proxyuser"); ok {
		sc.Node.ProxyUser = r
		sc.Wallet.ProxyUser = r
	}
	if r, ok := getIfIs(ctx, "proxypass"); ok {
		sc.Node.ProxyPass = r
		sc.Node.ProxyPass = r
	}
	if r, ok := getIfIs(ctx, "onion"); ok {
		NormalizeAddress(r, "9050", &sc.Node.OnionProxy)
	}
	if r, ok := getIfIs(ctx, "onionuser"); ok {
		sc.Node.OnionProxyUser = r
	}
	if r, ok := getIfIs(ctx, "onionpass"); ok {
		sc.Node.OnionProxyPass = r
	}
	if r, ok := getIfIs(ctx, "noonion"); ok {
		sc.Node.NoOnion = r == "true"
	}
	if r, ok := getIfIs(ctx, "torisolation"); ok {
		sc.Node.TorIsolation = r == "true"
	}
	if r, ok := getIfIs(ctx, "addpeers"); ok {
		NormalizeAddresses(r, n.DefaultPort, &sc.Node.AddPeers)
	}
	if r, ok := getIfIs(ctx, "connectpeers"); ok {
		NormalizeAddresses(r, n.DefaultPort, &sc.Node.ConnectPeers)
	}
	if r, ok := getIfIs(ctx, "maxpeers"); ok {
		if err := ParseInteger(
			r, "maxpeers", &sc.Node.MaxPeers); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if r, ok := getIfIs(ctx, "disablebanning"); ok {
		sc.Node.DisableBanning = r == "true"
	}
	if r, ok := getIfIs(ctx, "banduration"); ok {
		if err := ParseDuration(r, "banduration", &sc.Node.BanDuration); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if r, ok := getIfIs(ctx, "banthreshold"); ok {
		var bt int
		if err := ParseInteger(r, "banthtreshold", &bt); err != nil {
			log <- cl.Wrn(err.Error())
		} else {
			sc.Node.BanThreshold = uint32(bt)
		}
	}
	if r, ok := getIfIs(ctx, "whitelists"); ok {
		NormalizeAddresses(r, n.DefaultPort, &sc.Node.Whitelists)
	}
	if r, ok := getIfIs(ctx, "trickleinterval"); ok {
		if err := ParseDuration(
			r, "trickleinterval", &sc.Node.TrickleInterval); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if r, ok := getIfIs(ctx, "minrelaytxfee"); ok {
		if err := ParseFloat(
			r, "minrelaytxfee", &sc.Node.MinRelayTxFee); err != nil {
			log <- cl.Wrn(err.Error())
		}

	}
	if r, ok := getIfIs(ctx, "freetxrelaylimit"); ok {
		if err := ParseFloat(
			r, "freetxrelaylimit", &sc.Node.FreeTxRelayLimit); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if r, ok := getIfIs(ctx, "norelaypriority"); ok {
		sc.Node.NoRelayPriority = r == "true"
	}
	if r, ok := getIfIs(ctx, "nopeerbloomfilters"); ok {
		sc.Node.NoPeerBloomFilters = r == "true"
	}
	if r, ok := getIfIs(ctx, "nocfilters"); ok {
		sc.Node.NoCFilters = r == "true"
	}
	if r, ok := getIfIs(ctx, "blocksonly"); ok {
		sc.Node.BlocksOnly = r == "true"
	}
	if r, ok := getIfIs(ctx, "relaynonstd"); ok {
		sc.Node.RelayNonStd = r == "true"
	}
	if r, ok := getIfIs(ctx, "rejectnonstd"); ok {
		sc.Node.RejectNonStd = r == "true"
	}
	if r, ok := getIfIs(ctx, "maxorphantxs"); ok {
		if err := ParseInteger(r, "maxorphantxs", &sc.Node.MaxOrphanTxs); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if r, ok := getIfIs(ctx, "sigcachemaxsize"); ok {
		var scms int
		if err := ParseInteger(r, "sigcachemaxsize", &scms); err != nil {
			log <- cl.Wrn(err.Error())
		} else {
			sc.Node.SigCacheMaxSize = uint(scms)
		}
	}
	if r, ok := getIfIs(ctx, "generate"); ok {
		sc.Node.Generate = r == "true"
	}
	if r, ok := getIfIs(ctx, "genthreads"); ok {
		var gt int
		if err := ParseInteger(r, "genthreads", &gt); err != nil {
			log <- cl.Wrn(err.Error())
		} else {
			sc.Node.GenThreads = int32(gt)
		}
	}
	if r, ok := getIfIs(ctx, "algo"); ok {
		sc.Node.Algo = r
	}
	if r, ok := getIfIs(ctx, "miningaddrs"); ok {
		sc.Node.MiningAddrs = strings.Split(r, " ")
	}
	if r, ok := getIfIs(ctx, "minerlistener"); ok {
		NormalizeAddress(r, n.DefaultRPCPort, &sc.Node.MinerListener)
	}
	if r, ok := getIfIs(ctx, "minerpass"); ok {
		sc.Node.MinerPass = r
	}
	if r, ok := getIfIs(ctx, "addcheckpoints"); ok {
		sc.Node.AddCheckpoints = strings.Split(r, " ")
	}
	if r, ok := getIfIs(ctx, "disablecheckpoints"); ok {
		sc.Node.DisableCheckpoints = r == "true"
	}
	if r, ok := getIfIs(ctx, "blockminsize"); ok {
		if err := ParseUint32(r, "blockminsize", &sc.Node.BlockMinSize); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if r, ok := getIfIs(ctx, "blockmaxsize"); ok {
		if err := ParseUint32(r, "blockmaxsize", &sc.Node.BlockMaxSize); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if r, ok := getIfIs(ctx, "blockminweight"); ok {
		if err := ParseUint32(r, "blockminweight", &sc.Node.BlockMinWeight); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if r, ok := getIfIs(ctx, "blockmaxweight"); ok {
		if err := ParseUint32(
			r, "blockmaxweight", &sc.Node.BlockMaxWeight); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if r, ok := getIfIs(ctx, "blockprioritysize"); ok {
		if err := ParseUint32(
			r, "blockmaxweight", &sc.Node.BlockPrioritySize); err != nil {
			log <- cl.Wrn(err.Error())
		}
	}
	if r, ok := getIfIs(ctx, "uacomment"); ok {
		sc.Node.UserAgentComments = strings.Split(r, " ")
	}
	if r, ok := getIfIs(ctx, "upnp"); ok {
		sc.Node.Upnp = r == "true"
	}
	if r, ok := getIfIs(ctx, "dbtype"); ok {
		sc.Node.DbType = r
	}
	if r, ok := getIfIs(ctx, "disablednsseed"); ok {
		sc.Node.DisableDNSSeed = r == "true"
	}
	if r, ok := getIfIs(ctx, "profile"); ok {
		var p int
		if err := ParseInteger(r, "profile", &p); err == nil {
			sc.Node.Profile = fmt.Sprint(p)
		}
	}
	if r, ok := getIfIs(ctx, "cpuprofile"); ok {
		sc.Node.CPUProfile = r
	}

	// finished configuration

	SetLogging(ctx)

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

	// Service options which are only added on Windows.
	serviceOpts := serviceOptions{}
	// Perform service command and exit if specified.  Invalid service commands show an appropriate error.  Only runs on Windows since the runServiceCommand function will be nil when not on Windows.
	if serviceOpts.ServiceCommand != "" && runServiceCommand != nil {
		err := runServiceCommand(serviceOpts.ServiceCommand)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(0)
	}
	// Don't add peers from the config file when in regression test mode.
	if sc.Node.RegressionTest && len(sc.Node.AddPeers) > 0 {
		sc.Node.AddPeers = nil
	}
	// Set the mining algorithm correctly, default to random if unrecognised
	switch sc.Node.Algo {
	case "blake14lr", "cryptonight7v2", "keccak", "lyra2rev2", "scrypt", "skein", "x11", "stribog", "random", "easy":
	default:
		sc.Node.Algo = "random"
	}
	relayNonStd := n.ActiveNetParams.RelayNonStdTxs
	funcName := "loadConfig"
	switch {
	case sc.Node.RelayNonStd && sc.Node.RejectNonStd:
		str := "%s: rejectnonstd and relaynonstd cannot be used together -- choose only one"
		err := fmt.Errorf(str, funcName)
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	case sc.Node.RejectNonStd:
		relayNonStd = false
	case sc.Node.RelayNonStd:
		relayNonStd = true
	}
	sc.Node.RelayNonStd = relayNonStd
	// Append the network type to the data directory so it is "namespaced" per network.  In addition to the block database, there are other pieces of data that are saved to disk such as address manager state. All data is specific to a network, so namespacing the data directory means each individual piece of serialized data does not have to worry about changing names per network and such.
	sc.Node.DataDir = n.CleanAndExpandPath(sc.Node.DataDir)
	sc.Node.DataDir = filepath.Join(sc.Node.DataDir, n.NetName(n.ActiveNetParams))
	// Append the network type to the log directory so it is "namespaced" per network in the same fashion as the data directory.
	sc.Node.LogDir = n.CleanAndExpandPath(sc.Node.LogDir)
	sc.Node.LogDir = filepath.Join(sc.Node.LogDir, n.NetName(n.ActiveNetParams))

	// Initialize log rotation.  After log rotation has been initialized, the logger variables may be used.
	// initLogRotator(filepath.Join(sc.Node.LogDir, DefaultLogFilename))
	// Validate database type.
	if !n.ValidDbType(sc.Node.DbType) {
		str := "%s: The specified database type [%v] is invalid -- " +
			"supported types %v"
		err := fmt.Errorf(str, funcName, sc.Node.DbType, n.KnownDbTypes)
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// Validate profile port number
	if sc.Node.Profile != "" {
		profilePort, err := strconv.Atoi(sc.Node.Profile)
		if err != nil || profilePort < 1024 || profilePort > 65535 {
			str := "%s: The profile port must be between 1024 and 65535"
			err := fmt.Errorf(str, funcName)
			fmt.Fprintln(os.Stderr, err)
			fmt.Fprintln(os.Stderr, usageMessage)
			cl.Shutdown()
		}
	}
	// Don't allow ban durations that are too short.
	if sc.Node.BanDuration < time.Second {
		str := "%s: The banduration option may not be less than 1s -- parsed [%v]"
		err := fmt.Errorf(str, funcName, sc.Node.BanDuration)
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// Validate any given whitelisted IP addresses and networks.
	if len(sc.Node.Whitelists) > 0 {
		var ip net.IP
		StateCfg.ActiveWhitelists = make([]*net.IPNet, 0, len(sc.Node.Whitelists))
		for _, addr := range sc.Node.Whitelists {
			_, ipnet, err := net.ParseCIDR(addr)
			if err != nil {
				ip = net.ParseIP(addr)
				if ip == nil {
					str := "%s: The whitelist value of '%s' is invalid"
					err = fmt.Errorf(str, funcName, addr)
					log <- cl.Err(err.Error())
					fmt.Fprintln(os.Stderr, usageMessage)
					cl.Shutdown()
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
	if len(sc.Node.AddPeers) > 0 && len(sc.Node.ConnectPeers) > 0 {
		str := "%s: the --addpeer and --connect options can not be " +
			"mixed"
		err := fmt.Errorf(str, funcName)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
	}
	// --proxy or --connect without --listen disables listening.
	if (sc.Node.Proxy != "" || len(sc.Node.ConnectPeers) > 0) &&
		len(sc.Node.Listeners) == 0 {
		sc.Node.DisableListen = true
	}
	// Connect means no DNS seeding.
	if len(sc.Node.ConnectPeers) > 0 {
		sc.Node.DisableDNSSeed = true
	}
	// Add the default listener if none were specified. The default listener is all addresses on the listen port for the network we are to connect to.
	if len(sc.Node.Listeners) == 0 {
		sc.Node.Listeners = []string{
			net.JoinHostPort("", n.ActiveNetParams.DefaultPort),
		}
	}
	// Check to make sure limited and admin users don't have the same username
	if sc.Node.RPCUser == sc.Node.RPCLimitUser && sc.Node.RPCUser != "" {
		str := "%s: --rpcuser and --rpclimituser must not specify the same username"
		err := fmt.Errorf(str, funcName)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// Check to make sure limited and admin users don't have the same password
	if sc.Node.RPCPass == sc.Node.RPCLimitPass && sc.Node.RPCPass != "" {
		str := "%s: --rpcpass and --rpclimitpass must not specify the same password"
		err := fmt.Errorf(str, funcName)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// The RPC server is disabled if no username or password is provided.
	if (sc.Node.RPCUser == "" || sc.Node.RPCPass == "") &&
		(sc.Node.RPCLimitUser == "" || sc.Node.RPCLimitPass == "") {
		sc.Node.DisableRPC = true
	}
	if sc.Node.DisableRPC {
		log <- cl.Inf("RPC service is disabled")
	}
	// Default RPC to listen on localhost only.
	if !sc.Node.DisableRPC && len(sc.Node.RPCListeners) == 0 {
		addrs, err := net.LookupHost(n.DefaultRPCListener)
		if err != nil {
			log <- cl.Err(err.Error())
			cl.Shutdown()
		}
		sc.Node.RPCListeners = make([]string, 0, len(addrs))
		for _, addr := range addrs {
			addr = net.JoinHostPort(addr, n.ActiveNetParams.RPCPort)
			sc.Node.RPCListeners = append(sc.Node.RPCListeners, addr)
		}
	}
	if sc.Node.RPCMaxConcurrentReqs < 0 {
		str := "%s: The rpcmaxwebsocketconcurrentrequests option may not be less than 0 -- parsed [%d]"
		err := fmt.Errorf(str, funcName, sc.Node.RPCMaxConcurrentReqs)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	var err error
	// Validate the the minrelaytxfee.
	StateCfg.ActiveMinRelayTxFee, err = util.NewAmount(sc.Node.MinRelayTxFee)
	if err != nil {
		str := "%s: invalid minrelaytxfee: %v"
		err := fmt.Errorf(str, funcName, err)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// Limit the max block size to a sane value.
	if sc.Node.BlockMaxSize < n.BlockMaxSizeMin || sc.Node.BlockMaxSize >
		n.BlockMaxSizeMax {
		str := "%s: The blockmaxsize option must be in between %d and %d -- parsed [%d]"
		err := fmt.Errorf(str, funcName, n.BlockMaxSizeMin,
			n.BlockMaxSizeMax, sc.Node.BlockMaxSize)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// Limit the max block weight to a sane value.
	if sc.Node.BlockMaxWeight < n.BlockMaxWeightMin ||
		sc.Node.BlockMaxWeight > n.BlockMaxWeightMax {
		str := "%s: The blockmaxweight option must be in between %d and %d -- parsed [%d]"
		err := fmt.Errorf(str, funcName, n.BlockMaxWeightMin,
			n.BlockMaxWeightMax, sc.Node.BlockMaxWeight)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// Limit the max orphan count to a sane vlue.
	if sc.Node.MaxOrphanTxs < 0 {
		str := "%s: The maxorphantx option may not be less than 0 -- parsed [%d]"
		err := fmt.Errorf(str, funcName, sc.Node.MaxOrphanTxs)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// Limit the block priority and minimum block sizes to max block size.
	sc.Node.BlockPrioritySize = minUint32(sc.Node.BlockPrioritySize, sc.Node.BlockMaxSize)
	sc.Node.BlockMinSize = minUint32(sc.Node.BlockMinSize, sc.Node.BlockMaxSize)
	sc.Node.BlockMinWeight = minUint32(sc.Node.BlockMinWeight, sc.Node.BlockMaxWeight)
	switch {
	// If the max block size isn't set, but the max weight is, then we'll set the limit for the max block size to a safe limit so weight takes precedence.
	case sc.Node.BlockMaxSize == n.DefaultBlockMaxSize &&
		sc.Node.BlockMaxWeight != n.DefaultBlockMaxWeight:
		sc.Node.BlockMaxSize = blockchain.MaxBlockBaseSize - 1000
	// If the max block weight isn't set, but the block size is, then we'll scale the set weight accordingly based on the max block size value.
	case sc.Node.BlockMaxSize != n.DefaultBlockMaxSize &&
		sc.Node.BlockMaxWeight == n.DefaultBlockMaxWeight:
		sc.Node.BlockMaxWeight = sc.Node.BlockMaxSize * blockchain.WitnessScaleFactor
	}
	// Look for illegal characters in the user agent comments.
	for _, uaComment := range sc.Node.UserAgentComments {
		if strings.ContainsAny(uaComment, "/:()") {
			err := fmt.Errorf("%s: The following characters must not "+
				"appear in user agent comments: '/', ':', '(', ')'",
				funcName)
			log <- cl.Err(err.Error())
			fmt.Fprintln(os.Stderr, usageMessage)
			cl.Shutdown()

		}
	}
	// --txindex and --droptxindex do not mix.
	if sc.Node.TxIndex && sc.Node.DropTxIndex {
		err := fmt.Errorf("%s: the --txindex and --droptxindex options may  not be activated at the same time",
			funcName)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()

	}
	// --addrindex and --dropaddrindex do not mix.
	if sc.Node.AddrIndex && sc.Node.DropAddrIndex {
		err := fmt.Errorf("%s: the --addrindex and --dropaddrindex "+
			"options may not be activated at the same time",
			funcName)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// --addrindex and --droptxindex do not mix.
	if sc.Node.AddrIndex && sc.Node.DropTxIndex {
		err := fmt.Errorf("%s: the --addrindex and --droptxindex options may not be activated at the same time "+
			"because the address index relies on the transaction index",
			funcName)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// Check mining addresses are valid and saved parsed versions.
	StateCfg.ActiveMiningAddrs = make([]util.Address, 0, len(sc.Node.MiningAddrs))
	for _, strAddr := range sc.Node.MiningAddrs {
		addr, err := util.DecodeAddress(strAddr, n.ActiveNetParams.Params)
		if err != nil {
			str := "%s: mining address '%s' failed to decode: %v"
			err := fmt.Errorf(str, funcName, strAddr, err)
			log <- cl.Err(err.Error())
			fmt.Fprintln(os.Stderr, usageMessage)
			cl.Shutdown()
		}
		if !addr.IsForNet(n.ActiveNetParams.Params) {
			str := "%s: mining address '%s' is on the wrong network"
			err := fmt.Errorf(str, funcName, strAddr)
			log <- cl.Err(err.Error())
			fmt.Fprintln(os.Stderr, usageMessage)
			cl.Shutdown()
		}
		StateCfg.ActiveMiningAddrs = append(StateCfg.ActiveMiningAddrs, addr)
	}
	// Ensure there is at least one mining address when the generate flag is set.
	if (sc.Node.Generate || sc.Node.MinerListener != "") && len(sc.Node.MiningAddrs) == 0 {
		str := "%s: the generate flag is set, but there are no mining addresses specified "
		err := fmt.Errorf(str, funcName)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()

	}
	if sc.Node.MinerPass != "" {
		StateCfg.ActiveMinerKey = fork.Argon2i([]byte(sc.Node.MinerPass))
	}
	// Add default port to all listener addresses if needed and remove duplicate addresses.
	sc.Node.Listeners = n.NormalizeAddresses(sc.Node.Listeners,
		n.ActiveNetParams.DefaultPort)
	// Add default port to all rpc listener addresses if needed and remove duplicate addresses.
	sc.Node.RPCListeners = n.NormalizeAddresses(sc.Node.RPCListeners,
		n.ActiveNetParams.RPCPort)
	if !sc.Node.DisableRPC && !sc.Node.TLS {
		for _, addr := range sc.Node.RPCListeners {
			if err != nil {
				str := "%s: RPC listen interface '%s' is invalid: %v"
				err := fmt.Errorf(str, funcName, addr, err)
				log <- cl.Err(err.Error())
				fmt.Fprintln(os.Stderr, usageMessage)
				cl.Shutdown()
			}
		}
	}
	// Add default port to all added peer addresses if needed and remove duplicate addresses.
	sc.Node.AddPeers = n.NormalizeAddresses(sc.Node.AddPeers,
		n.ActiveNetParams.DefaultPort)
	sc.Node.ConnectPeers = n.NormalizeAddresses(sc.Node.ConnectPeers,
		n.ActiveNetParams.DefaultPort)
	// --noonion and --onion do not mix.
	if sc.Node.NoOnion && sc.Node.OnionProxy != "" {
		err := fmt.Errorf("%s: the --noonion and --onion options may not be activated at the same time", funcName)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// Check the checkpoints for syntax errors.
	StateCfg.AddedCheckpoints, err = n.ParseCheckpoints(sc.Node.AddCheckpoints)
	if err != nil {
		str := "%s: Error parsing checkpoints: %v"
		err := fmt.Errorf(str, funcName, err)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// Tor stream isolation requires either proxy or onion proxy to be set.
	if sc.Node.TorIsolation && sc.Node.Proxy == "" && sc.Node.OnionProxy == "" {
		str := "%s: Tor stream isolation requires either proxy or onionproxy to be set"
		err := fmt.Errorf(str, funcName)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		cl.Shutdown()
	}
	// Setup dial and DNS resolution (lookup) functions depending on the specified options.  The default is to use the standard net.DialTimeout function as well as the system DNS resolver.  When a proxy is specified, the dial function is set to the proxy specific dial function and the lookup is set to use tor (unless --noonion is specified in which case the system DNS resolver is used).
	StateCfg.Dial = net.DialTimeout
	StateCfg.Lookup = net.LookupIP
	if sc.Node.Proxy != "" {
		_, _, err := net.SplitHostPort(sc.Node.Proxy)
		if err != nil {
			str := "%s: Proxy address '%s' is invalid: %v"
			err := fmt.Errorf(str, funcName, sc.Node.Proxy, err)
			log <- cl.Err(err.Error())
			fmt.Fprintln(os.Stderr, usageMessage)
			cl.Shutdown()
		}
		// Tor isolation flag means proxy credentials will be overridden unless there is also an onion proxy configured in which case that one will be overridden.
		torIsolation := false
		if sc.Node.TorIsolation && sc.Node.OnionProxy == "" &&
			(sc.Node.ProxyUser != "" || sc.Node.ProxyPass != "") {
			torIsolation = true
			fmt.Fprintln(os.Stderr, "Tor isolation set -- "+
				"overriding specified proxy user credentials")
		}
		proxy := &socks.Proxy{
			Addr:         sc.Node.Proxy,
			Username:     sc.Node.ProxyUser,
			Password:     sc.Node.ProxyPass,
			TorIsolation: torIsolation,
		}
		StateCfg.Dial = proxy.DialTimeout
		// Treat the proxy as tor and perform DNS resolution through it unless the --noonion flag is set or there is an onion-specific proxy configured.
		if !sc.Node.NoOnion && sc.Node.OnionProxy == "" {
			StateCfg.Lookup = func(host string) ([]net.IP, error) {
				return connmgr.TorLookupIP(host, sc.Node.Proxy)
			}
		}
	}
	// Setup onion address dial function depending on the specified options. The default is to use the same dial function selected above.  However, when an onion-specific proxy is specified, the onion address dial function is set to use the onion-specific proxy while leaving the normal dial function as selected above.  This allows .onion address traffic to be routed through a different proxy than normal traffic.
	if sc.Node.OnionProxy != "" {
		_, _, err := net.SplitHostPort(sc.Node.OnionProxy)
		if err != nil {
			str := "%s: Onion proxy address '%s' is invalid: %v"
			err := fmt.Errorf(str, funcName, sc.Node.OnionProxy, err)
			log <- cl.Err(err.Error())
			fmt.Fprintln(os.Stderr, usageMessage)
			cl.Shutdown()
		}
		// Tor isolation flag means onion proxy credentials will be overridden.
		if sc.Node.TorIsolation &&
			(sc.Node.OnionProxyUser != "" || sc.Node.OnionProxyPass != "") {
			fmt.Fprintln(os.Stderr, "Tor isolation set -- "+
				"overriding specified onionproxy user "+
				"credentials ")
		}
		StateCfg.Oniondial = func(network, addr string, timeout time.Duration) (net.Conn, error) {
			proxy := &socks.Proxy{
				Addr:         sc.Node.OnionProxy,
				Username:     sc.Node.OnionProxyUser,
				Password:     sc.Node.OnionProxyPass,
				TorIsolation: sc.Node.TorIsolation,
			}
			return proxy.DialTimeout(network, addr, timeout)
		}
		// When configured in bridge mode (both --onion and --proxy are configured), it means that the proxy configured by --proxy is not a tor proxy, so override the DNS resolution to use the onion-specific proxy.
		if sc.Node.Proxy != "" {
			StateCfg.Lookup = func(host string) ([]net.IP, error) {
				return connmgr.TorLookupIP(host, sc.Node.OnionProxy)
			}
		}
	} else {
		StateCfg.Oniondial = StateCfg.Dial
	}
	// Specifying --noonion means the onion address dial function results in an error.
	if sc.Node.NoOnion {
		StateCfg.Oniondial = func(a, b string, t time.Duration) (net.Conn, error) {
			return nil, errors.New("tor has been disabled")
		}
	}
	sc.Wallet.PodUsername = sc.Node.RPCUser
	sc.Wallet.PodPassword = sc.Node.RPCPass
}
