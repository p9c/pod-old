package app

import (
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
	blockchain "git.parallelcoin.io/pod/pkg/chain"
	netparams "git.parallelcoin.io/pod/pkg/chain/config/params"
	"git.parallelcoin.io/pod/pkg/chain/fork"
	"git.parallelcoin.io/pod/pkg/peer/connmgr"
	"git.parallelcoin.io/pod/pkg/util"
	cl "git.parallelcoin.io/pod/pkg/util/cl"
	"github.com/btcsuite/go-socks/socks"
	"github.com/davecgh/go-spew/spew"
	"gopkg.in/urfave/cli.v1"
	"gopkg.in/yaml.v1"
)

// StateCfg is a reference to the main node state configuration struct
var StateCfg = node.StateCfg

func nodeHandleSave() {
	appConfigCommon.Save = false
	*nodeConfig.LogDir = *nodeConfig.DataDir
	yn, e := yaml.Marshal(nodeConfig)
	if e == nil {
		EnsureDir(*nodeConfig.ConfigFile)
		e = ioutil.WriteFile(
			*nodeConfig.ConfigFile, yn, 0600)
		if e != nil {
			panic(e)
		}
	} else {
		panic(e)
	}
}

func nodeHandle(c *cli.Context) error {
	Log.SetLevel("trace")
	log <- cl.Debug{"running node"}

	*nodeConfig.DataDir = filepath.Join(
		appConfigCommon.Datadir,
		nodeAppName)
	*nodeConfig.ConfigFile = filepath.Join(
		*nodeConfig.DataDir,
		nodeConfigFilename)
	if FileExists(*nodeConfig.ConfigFile) {
		ncb, e := ioutil.ReadFile(*nodeConfig.ConfigFile)
		if e != nil {
			panic(e)
		}
		ncf := &node.Config{}
		e = yaml.Unmarshal(ncb, ncf)
		if e != nil {
			panic(e)
		}
		nodeConfig = ncf
	} else {
		appConfigCommon.Save = true
	}
	*nodeConfig.LogDir = *nodeConfig.DataDir
	if !c.Parent().Bool("useproxy") {
		*nodeConfig.Proxy = ""
	}
	loglevel := c.Parent().String("loglevel")
	switch loglevel {
	case "trace", "debug", "info", "warn", "error", "fatal":
	default:
		*nodeConfig.DebugLevel = "info"
	}
	network := c.Parent().String("network")
	switch network {
	case "testnet", "testnet3", "t":
		nodeConfig.TestNet3 = &True
		nodeConfig.SimNet = &False
		nodeConfig.RegressionTest = &False
		activeNetParams = &netparams.TestNet3Params
	case "regtestnet", "regressiontest", "r":
		nodeConfig.TestNet3 = &False
		nodeConfig.SimNet = &False
		nodeConfig.RegressionTest = &True
		activeNetParams = &netparams.RegressionTestParams
	case "simnet", "s":
		nodeConfig.TestNet3 = &False
		nodeConfig.SimNet = &True
		nodeConfig.RegressionTest = &False
		activeNetParams = &netparams.SimNetParams
	default:
		if network != "mainnet" && network != "m" {
			fmt.Println("using mainnet for node")
		}
		nodeConfig.TestNet3 = &False
		nodeConfig.SimNet = &False
		nodeConfig.RegressionTest = &False
		activeNetParams = &netparams.MainNetParams

	}
	if !*nodeConfig.Onion {
		*nodeConfig.OnionProxy = ""
	}
	// TODO: now to sanitize the rest
	port := node.DefaultPort
	NormalizeStringSliceAddresses(nodeConfig.AddPeers, port)
	NormalizeStringSliceAddresses(nodeConfig.ConnectPeers, port)
	NormalizeStringSliceAddresses(nodeConfig.Listeners, port)
	NormalizeStringSliceAddresses(nodeConfig.Whitelists, port)
	NormalizeStringSliceAddresses(nodeConfig.RPCListeners, port)

	log <- cl.Debug{spew.Sdump(nodeConfig)}
	cl.Register.SetAllLevels(*nodeConfig.DebugLevel)
	_ = podHandle(c)
	if appConfigCommon.Save {
		appConfigCommon.Save = false
		podHandleSave()
		nodeHandleSave()
		return nil
	}

	// serviceOptions defines the configuration options for the daemon as a service on Windows.
	type serviceOptions struct {
		ServiceCommand string `short:"s" long:"service" description:"Service command {install, remove, start, stop}"`
	}

	var usageMessage = fmt.Sprintf("use `%s help node` to show usage", appName)

	// runServiceCommand is only set to a real function on Windows.  It is used to parse and execute service commands specified via the -s flag.
	var runServiceCommand func(string) error

	// Service options which are only added on Windows.
	serviceOpts := serviceOptions{}

	// Perform service command and exit if specified.  Invalid service commands show an appropriate error.  Only runs on Windows since the runServiceCommand function will be nil when not on Windows.
	if serviceOpts.ServiceCommand != "" && runServiceCommand != nil {
		err := runServiceCommand(serviceOpts.ServiceCommand)
		if err != nil {
			log <- cl.Error{err}
			return err
		}
		return nil
	}

	// Don't add peers from the config file when in regression test mode.
	if *nodeConfig.RegressionTest && len(*nodeConfig.AddPeers) > 0 {
		nodeConfig.AddPeers = nil
	}

	// Set the mining algorithm correctly, default to random if unrecognised
	switch *nodeConfig.Algo {
	case fork.P9AlgoVers[0], fork.P9AlgoVers[1], fork.P9AlgoVers[2], fork.P9AlgoVers[3], fork.P9AlgoVers[4], fork.P9AlgoVers[5], fork.P9AlgoVers[6], fork.P9AlgoVers[7], fork.P9AlgoVers[8], "random", "easy":
	default:
		*nodeConfig.Algo = "random"
	}
	relayNonStd := *nodeConfig.RelayNonStd
	funcName := "loadConfig"
	switch {
	case *nodeConfig.RelayNonStd && *nodeConfig.RejectNonStd:
		errf := "%s: rejectnonstd and relaynonstd cannot be used together -- choose only one"
		log <- cl.Errorf{errf, funcName}
		// log <- cl.Err(usageMessage)
		return fmt.Errorf(errf, funcName)

	case *nodeConfig.RejectNonStd:
		relayNonStd = false

	case *nodeConfig.RelayNonStd:
		relayNonStd = true
	}
	*nodeConfig.RelayNonStd = relayNonStd

	// Append the network type to the data directory so it is "namespaced" per network.  In addition to the block database, there are other pieces of data that are saved to disk such as address manager state. All data is specific to a network, so namespacing the data directory means each individual piece of serialized data does not have to worry about changing names per network and such.
	log <- cl.Debug{"netname", activeNetParams.Name}
	*nodeConfig.DataDir = CleanAndExpandPath(*nodeConfig.DataDir)
	*nodeConfig.DataDir = filepath.Join(
		*nodeConfig.DataDir, activeNetParams.Name)
	*nodeConfig.LogDir = CleanAndExpandPath(*nodeConfig.DataDir)
	*nodeConfig.LogDir = filepath.Join(
		*nodeConfig.DataDir, activeNetParams.Name)

	// Validate database type.
	if !node.ValidDbType(*nodeConfig.DbType) {
		str := "%s: The specified database type [%v] is invalid -- " +
			"supported types %v"
		err := fmt.Errorf(str, funcName, *nodeConfig.DbType, node.KnownDbTypes)
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, usageMessage)
		return err
	}

	// Validate profile port number
	if *nodeConfig.Profile != "" {
		profilePort, err := strconv.Atoi(*nodeConfig.Profile)
		if err != nil || profilePort < 1024 || profilePort > 65535 {
			str := "%s: The profile port must be between 1024 and 65535"
			err := fmt.Errorf(str, funcName)
			log <- cl.Error{err}
			return err
		}
	}

	// Don't allow ban durations that are too short.
	if *nodeConfig.BanDuration < time.Second {
		err := fmt.Errorf("%s: The banduration option may not be less than 1s -- parsed [%v]", funcName, *nodeConfig.BanDuration)
		log <- cl.Error{err}
		return err
	}

	// Validate any given whitelisted IP addresses and networks.
	if len(*nodeConfig.Whitelists) > 0 {
		var ip net.IP
		StateCfg.ActiveWhitelists = make([]*net.IPNet, 0, len(*nodeConfig.Whitelists))
		for _, addr := range *nodeConfig.Whitelists {
			_, ipnet, err := net.ParseCIDR(addr)
			if err != nil {
				err = fmt.Errorf("%s '%s'", cl.Ine(), err.Error())
				ip = net.ParseIP(addr)
				if ip == nil {
					str := err.Error() + " %s: The whitelist value of '%s' is invalid"
					err = fmt.Errorf(str, funcName, addr)
					log <- cl.Err(err.Error())
					fmt.Fprintln(os.Stderr, usageMessage)
					return err
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
	if len(*nodeConfig.AddPeers) > 0 && len(*nodeConfig.ConnectPeers) > 0 {
		err := fmt.Errorf(
			"%s: the --addpeer and --connect options can not be mixed",
			funcName)
		log <- cl.Error{err}
		return err
	}

	// --proxy or --connect without --listen disables listening.
	if (*nodeConfig.Proxy != "" || len(*nodeConfig.ConnectPeers) > 0) &&
		len(*nodeConfig.Listeners) == 0 {
		*nodeConfig.DisableListen = true
	}

	// Connect means no DNS seeding.
	if len(*nodeConfig.ConnectPeers) > 0 {
		*nodeConfig.DisableDNSSeed = true
	}

	// Add the default listener if none were specified. The default listener is all addresses on the listen port for the network we are to connect to.
	if len(*nodeConfig.Listeners) == 0 {
		*nodeConfig.Listeners = []string{
			net.JoinHostPort("127.0.0.1", activeNetParams.DefaultPort),
		}
	}

	// Check to make sure limited and admin users don't have the same username
	if *nodeConfig.RPCUser != "" &&
		*nodeConfig.RPCUser == *nodeConfig.RPCLimitUser {
		str := "%s: --rpcuser and --rpclimituser must not specify the same username"
		err := fmt.Errorf(str, funcName)
		log <- cl.Error{err}
		return err
	}

	// Check to make sure limited and admin users don't have the same password
	if *nodeConfig.RPCPass != "" &&
		*nodeConfig.RPCPass == *nodeConfig.RPCLimitPass {
		str := "%s: --rpcpass and --rpclimitpass must not specify the same password"
		err := fmt.Errorf(str, funcName)
		log <- cl.Error{err}
		// fmt.Fprintln(os.Stderr, usageMessage)
		return err
	}

	// The RPC server is disabled if no username or password is provided.
	if (*nodeConfig.RPCUser == "" || *nodeConfig.RPCPass == "") &&
		(*nodeConfig.RPCLimitUser == "" || *nodeConfig.RPCLimitPass == "") {
		*nodeConfig.DisableRPC = true
	}
	if *nodeConfig.DisableRPC {
		log <- cl.Inf("RPC service is disabled")
	}

	// Default RPC to listen on localhost only.
	if !*nodeConfig.DisableRPC && len(*nodeConfig.RPCListeners) == 0 {
		addrs, err := net.LookupHost(node.DefaultRPCListener)
		if err != nil {
			log <- cl.Error{err}
			return err
		}
		*nodeConfig.RPCListeners = make([]string, 0, len(addrs))
		for _, addr := range addrs {
			addr = net.JoinHostPort(addr, activeNetParams.RPCClientPort)
			*nodeConfig.RPCListeners = append(*nodeConfig.RPCListeners, addr)
		}
	}

	if *nodeConfig.RPCMaxConcurrentReqs < 0 {
		str := "%s: The rpcmaxwebsocketconcurrentrequests option may not be less than 0 -- parsed [%d]"
		err := fmt.Errorf(str, funcName, *nodeConfig.RPCMaxConcurrentReqs)
		log <- cl.Error{err}
		// fmt.Fprintln(os.Stderr, usageMessage)
		return err
	}

	var err error
	// Validate the the minrelaytxfee.
	StateCfg.ActiveMinRelayTxFee, err = util.NewAmount(*nodeConfig.MinRelayTxFee)
	if err != nil {
		str := "%s: invalid minrelaytxfee: %v"
		err := fmt.Errorf(str, funcName, err)
		log <- cl.Error{err}
		// fmt.Fprintln(os.Stderr, usageMessage)
		return err
	}

	// Limit the max block size to a sane value.
	if *nodeConfig.BlockMaxSize < node.BlockMaxSizeMin ||
		*nodeConfig.BlockMaxSize > node.BlockMaxSizeMax {
		str := "%s: The blockmaxsize option must be in between %d and %d -- parsed [%d]"
		err := fmt.Errorf(str, funcName, node.BlockMaxSizeMin,
			node.BlockMaxSizeMax, *nodeConfig.BlockMaxSize)
		log <- cl.Error{err}
		// fmt.Fprintln(os.Stderr, usageMessage)
		return err
	}

	// Limit the max block weight to a sane value.
	if *nodeConfig.BlockMaxWeight < node.BlockMaxWeightMin ||
		*nodeConfig.BlockMaxWeight > node.BlockMaxWeightMax {
		str := "%s: The blockmaxweight option must be in between %d and %d -- parsed [%d]"
		err := fmt.Errorf(str, funcName, node.BlockMaxWeightMin,
			node.BlockMaxWeightMax, *nodeConfig.BlockMaxWeight)
		log <- cl.Error{err}
		// fmt.Fprintln(os.Stderr, usageMessage)
		return err
	}

	// Limit the max orphan count to a sane vlue.
	if *nodeConfig.MaxOrphanTxs < 0 {
		str := "%s: The maxorphantx option may not be less than 0 -- parsed [%d]"
		err := fmt.Errorf(str, funcName, *nodeConfig.MaxOrphanTxs)
		log <- cl.Error{err}
		// fmt.Fprintln(os.Stderr, usageMessage)
		return err
	}

	// Limit the block priority and minimum block sizes to max block size.
	*nodeConfig.BlockPrioritySize = int(minUint32(
		uint32(*nodeConfig.BlockPrioritySize),
		uint32(*nodeConfig.BlockMaxSize)))
	*nodeConfig.BlockMinSize = int(minUint32(
		uint32(*nodeConfig.BlockMinSize),
		uint32(*nodeConfig.BlockMaxSize)))
	*nodeConfig.BlockMinWeight = int(minUint32(
		uint32(*nodeConfig.BlockMinWeight),
		uint32(*nodeConfig.BlockMaxWeight)))
	switch {

	// If the max block size isn't set, but the max weight is, then we'll set the limit for the max block size to a safe limit so weight takes precedence.
	case *nodeConfig.BlockMaxSize == node.DefaultBlockMaxSize &&
		*nodeConfig.BlockMaxWeight != node.DefaultBlockMaxWeight:
		*nodeConfig.BlockMaxSize = blockchain.MaxBlockBaseSize - 1000

	// If the max block weight isn't set, but the block size is, then we'll scale the set weight accordingly based on the max block size value.
	case *nodeConfig.BlockMaxSize != node.DefaultBlockMaxSize &&
		*nodeConfig.BlockMaxWeight == node.DefaultBlockMaxWeight:
		*nodeConfig.BlockMaxWeight = *nodeConfig.BlockMaxSize * blockchain.WitnessScaleFactor
	}

	// Look for illegal characters in the user agent comments.
	for _, uaComment := range *nodeConfig.UserAgentComments {
		if strings.ContainsAny(uaComment, "/:()") {
			err := fmt.Errorf("%s: The following characters must not "+
				"appear in user agent comments: '/', ':', '(', ')'",
				funcName)
			log <- cl.Err(err.Error())
			// fmt.Fprintln(os.Stderr, usageMessage)
			return err
		}
	}
	// --addrindex and --dropaddrindex do not mix.
	if *nodeConfig.AddrIndex && *nodeConfig.DropAddrIndex {
		err := fmt.Errorf("%s: the --addrindex and --dropaddrindex options may not be activated at the same time",
			funcName)
		log <- cl.Error{err}
		// fmt.Fprintln(os.Stderr, usageMessage)
		return err
	}

	// Check mining addresses are valid and saved parsed versions.
	StateCfg.ActiveMiningAddrs = make([]util.Address, 0, len(*nodeConfig.MiningAddrs))
	for _, strAddr := range *nodeConfig.MiningAddrs {
		addr, err := util.DecodeAddress(strAddr, activeNetParams.Params)
		if err != nil {
			str := "%s: mining address '%s' failed to decode: %v"
			err := fmt.Errorf(str, funcName, strAddr, err)
			log <- cl.Err(err.Error())
			// fmt.Fprintln(os.Stderr, usageMessage)
			return err
		}
		if !addr.IsForNet(activeNetParams.Params) {
			str := "%s: mining address '%s' is on the wrong network"
			err := fmt.Errorf(str, funcName, strAddr)
			log <- cl.Error{err}
			// fmt.Fprintln(os.Stderr, usageMessage)
			return err
		}
		StateCfg.ActiveMiningAddrs = append(StateCfg.ActiveMiningAddrs, addr)
	}

	// Ensure there is at least one mining address when the generate flag is set.
	if (*nodeConfig.Generate || *nodeConfig.MinerListener != "") && len(*nodeConfig.MiningAddrs) == 0 {
		str := "%s: the generate flag is set, but there are no mining addresses specified "
		err := fmt.Errorf(str, funcName)
		log <- cl.Err(err.Error())
		fmt.Fprintln(os.Stderr, usageMessage)
		os.Exit(1)
	}
	if *nodeConfig.MinerPass != "" {
		StateCfg.ActiveMinerKey = fork.Argon2i([]byte(*nodeConfig.MinerPass))
	}

	// Add default port to all rpc listener addresses if needed and remove duplicate addresses.
	*nodeConfig.RPCListeners = node.NormalizeAddresses(*nodeConfig.RPCListeners,
		activeNetParams.RPCClientPort)
	if !*nodeConfig.DisableRPC && !*nodeConfig.TLS {
		for _, addr := range *nodeConfig.RPCListeners {
			if err != nil {
				str := "%s: RPC listen interface '%s' is invalid: %v"
				err := fmt.Errorf(str, funcName, addr, err)
				log <- cl.Error{err}
				// fmt.Fprintln(os.Stderr, usageMessage)
				return err
			}
		}
	}

	// Add default port to all listener addresses if needed and remove duplicate addresses.
	*nodeConfig.Listeners = node.NormalizeAddresses(*nodeConfig.Listeners,
		activeNetParams.DefaultPort)

	// Add default port to all added peer addresses if needed and remove duplicate addresses.
	*nodeConfig.AddPeers = node.NormalizeAddresses(*nodeConfig.AddPeers,
		activeNetParams.DefaultPort)
	*nodeConfig.ConnectPeers = node.NormalizeAddresses(*nodeConfig.ConnectPeers,
		activeNetParams.DefaultPort)

	// --onionproxy and not --onion are contradictory (TODO: this is kinda stupid hm? switch *and* toggle by presence of flag value, one should be enough)
	if !*nodeConfig.Onion && *nodeConfig.OnionProxy != "" {
		err := fmt.Errorf("%s: the --onionproxy and --onion options may not be activated at the same time", funcName)
		log <- cl.Error{err}
		// fmt.Fprintln(os.Stderr, usageMessage)
		return err
	}

	// Check the checkpoints for syntax errors.
	StateCfg.AddedCheckpoints, err = node.ParseCheckpoints(*nodeConfig.AddCheckpoints)
	if err != nil {
		str := "%s: Error parsing checkpoints: %v"
		err := fmt.Errorf(str, funcName, err)
		log <- cl.Err(err.Error())
		// fmt.Fprintln(os.Stderr, usageMessage)
		return err
	}

	// Tor stream isolation requires either proxy or onion proxy to be set.
	if *nodeConfig.TorIsolation &&
		*nodeConfig.Proxy == "" &&
		*nodeConfig.OnionProxy == "" {
		str := "%s: Tor stream isolation requires either proxy or onionproxy to be set"
		err := fmt.Errorf(str, funcName)
		log <- cl.Error{err}
		// fmt.Fprintln(os.Stderr, usageMessage)
		return err
	}

	// Setup dial and DNS resolution (lookup) functions depending on the specified options.  The default is to use the standard net.DialTimeout function as well as the system DNS resolver.  When a proxy is specified, the dial function is set to the proxy specific dial function and the lookup is set to use tor (unless --noonion is specified in which case the system DNS resolver is used).
	StateCfg.Dial = net.DialTimeout
	StateCfg.Lookup = net.LookupIP
	if *nodeConfig.Proxy != "" {
		log <- cl.Info{"we are loading a proxy!"}
		_, _, err := net.SplitHostPort(*nodeConfig.Proxy)
		if err != nil {
			str := "%s: Proxy address '%s' is invalid: %v"
			err := fmt.Errorf(str, funcName, *nodeConfig.Proxy, err)
			log <- cl.Error{err}
			// fmt.Fprintln(os.Stderr, usageMessage)
			return err
		}

		// Tor isolation flag means proxy credentials will be overridden unless there is also an onion proxy configured in which case that one will be overridden.
		torIsolation := false
		if *nodeConfig.TorIsolation &&
			*nodeConfig.OnionProxy == "" &&
			(*nodeConfig.ProxyUser != "" ||
				*nodeConfig.ProxyPass != "") {
			torIsolation = true
			log <- cl.Warn{
				"Tor isolation set -- overriding specified proxy user credentials"}
		}
		proxy := &socks.Proxy{
			Addr:         *nodeConfig.Proxy,
			Username:     *nodeConfig.ProxyUser,
			Password:     *nodeConfig.ProxyPass,
			TorIsolation: torIsolation,
		}
		StateCfg.Dial = proxy.DialTimeout
		// Treat the proxy as tor and perform DNS resolution through it unless the --noonion flag is set or there is an onion-specific proxy configured.
		if *nodeConfig.Onion &&
			*nodeConfig.OnionProxy == "" {
			StateCfg.Lookup = func(host string) ([]net.IP, error) {
				return connmgr.TorLookupIP(host, *nodeConfig.Proxy)
			}
		}
	}

	// Setup onion address dial function depending on the specified options. The default is to use the same dial function selected above.  However, when an onion-specific proxy is specified, the onion address dial function is set to use the onion-specific proxy while leaving the normal dial function as selected above.  This allows .onion address traffic to be routed through a different proxy than normal traffic.
	if *nodeConfig.OnionProxy != "" {
		_, _, err := net.SplitHostPort(*nodeConfig.OnionProxy)
		if err != nil {
			str := "%s: Onion proxy address '%s' is invalid: %v"
			err := fmt.Errorf(str, funcName, *nodeConfig.OnionProxy, err)
			log <- cl.Error{err}
			// fmt.Fprintln(os.Stderr, usageMessage)
			return err
		}

		// Tor isolation flag means onion proxy credentials will be overriddenode.
		if *nodeConfig.TorIsolation &&
			(*nodeConfig.OnionProxyUser != "" || *nodeConfig.OnionProxyPass != "") {
			fmt.Fprintln(os.Stderr, "Tor isolation set -- "+
				"overriding specified onionproxy user "+
				"credentials ")
		}
		StateCfg.Oniondial = func(network, addr string, timeout time.Duration) (net.Conn, error) {
			proxy := &socks.Proxy{
				Addr:         *nodeConfig.OnionProxy,
				Username:     *nodeConfig.OnionProxyUser,
				Password:     *nodeConfig.OnionProxyPass,
				TorIsolation: *nodeConfig.TorIsolation,
			}
			return proxy.DialTimeout(network, addr, timeout)
		}

		// When configured in bridge mode (both --onion and --proxy are configured), it means that the proxy configured by --proxy is not a tor proxy, so override the DNS resolution to use the onion-specific proxy.
		if *nodeConfig.Proxy != "" {
			StateCfg.Lookup = func(host string) ([]net.IP, error) {
				return connmgr.TorLookupIP(host, *nodeConfig.OnionProxy)
			}
		}
	} else {
		StateCfg.Oniondial = StateCfg.Dial
	}

	// Specifying --noonion means the onion address dial function results in an error.
	if !*nodeConfig.Onion {
		StateCfg.Oniondial = func(a, b string, t time.Duration) (net.Conn, error) {
			return nil, errors.New("tor has been disabled")
		}
	}

	if appConfigCommon.Save {
		appConfigCommon.Save = false
		podHandleSave()
		nodeHandleSave()
		return nil
	}

	return launchNode(c)
}

func NormalizeStringSliceAddresses(a *cli.StringSlice, port string) {
	variable := []string(*a)
	NormalizeAddresses(
		strings.Join(variable, " "),
		port, &variable)
	*a = variable
}
