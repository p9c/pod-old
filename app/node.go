package app

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"git.parallelcoin.io/dev/pod/cmd/node"
	blockchain "git.parallelcoin.io/dev/pod/pkg/chain"
	"git.parallelcoin.io/dev/pod/pkg/chain/fork"
	"git.parallelcoin.io/dev/pod/pkg/peer/connmgr"
	"git.parallelcoin.io/dev/pod/pkg/util"
	cl "git.parallelcoin.io/dev/pod/pkg/util/cl"
	"github.com/btcsuite/go-socks/socks"
	"gopkg.in/urfave/cli.v1"
)

// StateCfg is a reference to the main node state configuration struct
var StateCfg = node.StateCfg

func nodeHandle(c *cli.Context) error {

	Configure()

	// serviceOptions defines the configuration options for the daemon as a service on Windows.
	type serviceOptions struct {
		ServiceCommand string `short:"s" long:"service" description:"Service command {install, remove, start, stop}"`
	}

	var usageMessage = fmt.Sprintf("use `%s help node` to show usage", appName)

	// runServiceCommand is only set to a real function on Windows.  It is used to parse and execute service commands specified via the -s flag.
	var runServiceCommand func(string) error

	// Service options which are only added on Windows.
	serviceOpts := serviceOptions{}

	// Perform service command and exit if specified.  Invalid service commands show an appropriate error.
	// Only runs on Windows since the runServiceCommand function will be nil when not on Windows.
	if serviceOpts.ServiceCommand != "" && runServiceCommand != nil {

		err := runServiceCommand(serviceOpts.ServiceCommand)

		if err != nil {

			log <- cl.Error{err}

			return err
		}

		return nil
	}

	// Don't add peers from the config file when in regression test mode.

	if *podConfig.RegressionTest && len(*podConfig.AddPeers) > 0 {

		podConfig.AddPeers = nil
	}

	// Set the mining algorithm correctly, default to random if unrecognised
	switch *podConfig.Algo {

	case fork.P9AlgoVers[0], fork.P9AlgoVers[1], fork.P9AlgoVers[2], fork.P9AlgoVers[3], fork.P9AlgoVers[4], fork.P9AlgoVers[5], fork.P9AlgoVers[6], fork.P9AlgoVers[7], fork.P9AlgoVers[8], "random", "easy":

	default:
		*podConfig.Algo = "random"
	}
	log <- cl.Debug{"mining algorithm", *podConfig.Algo}

	relayNonStd := *podConfig.RelayNonStd
	funcName := "loadConfig"

	switch {

	case *podConfig.RelayNonStd && *podConfig.RejectNonStd:
		errf := "%s: rejectnonstd and relaynonstd cannot be used together -- choose only one"

		log <- cl.Errorf{errf, funcName}

		// log <- cl.Err(usageMessage)
		return fmt.Errorf(errf, funcName)

	case *podConfig.RejectNonStd:
		relayNonStd = false

	case *podConfig.RelayNonStd:
		relayNonStd = true
	}

	*podConfig.RelayNonStd = relayNonStd

	// Append the network type to the data directory so it is "namespaced" per network.
	// In addition to the block database, there are other pieces of data that are saved
	// to disk such as address manager state. All data is specific to a network, so
	//namespacing the data directory means each individual piece of serialized data
	// does not have to worry about changing names per network and such.
	log <- cl.Debug{"netname", activeNetParams.Name}

	*podConfig.DataDir = CleanAndExpandPath(*podConfig.DataDir)
	*podConfig.DataDir = filepath.Join(
		*podConfig.DataDir, activeNetParams.Name)
	*podConfig.LogDir = CleanAndExpandPath(*podConfig.DataDir)
	*podConfig.LogDir = filepath.Join(
		*podConfig.DataDir, activeNetParams.Name)

	// Validate database type.
	log <- cl.Debug{"validating database type"}
	if !node.ValidDbType(*podConfig.DbType) {

		str := "%s: The specified database type [%v] is invalid -- " +
			"supported types %v"
		err := fmt.Errorf(str, funcName, *podConfig.DbType, node.KnownDbTypes)
		log <- cl.Error{err}
		return err
	}

	// Validate profile port number
	log <- cl.Debug{"validating profile port number"}
	if *podConfig.Profile != "" {

		profilePort, err := strconv.Atoi(*podConfig.Profile)

		if err != nil || profilePort < 1024 || profilePort > 65535 {

			str := "%s: The profile port must be between 1024 and 65535"
			err := fmt.Errorf(str, funcName)

			log <- cl.Error{err}

			return err
		}

	}

	// Don't allow ban durations that are too short.
	log <- cl.Debug{"validating ban duration"}
	if *podConfig.BanDuration < time.Second {

		err := fmt.Errorf("%s: The banduration option may not be less than 1s -- parsed [%v]", funcName, *podConfig.BanDuration)

		log <- cl.Error{err}

		return err
	}

	// Validate any given whitelisted IP addresses and networks.
	log <- cl.Debug{"validating whitelists"}
	if len(*podConfig.Whitelists) > 0 {

		var ip net.IP
		StateCfg.ActiveWhitelists = make([]*net.IPNet, 0, len(*podConfig.Whitelists))

		for _, addr := range *podConfig.Whitelists {

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

	log <- cl.Debug{"checking addpeer and connectpeer lists"}
	if len(*podConfig.AddPeers) > 0 && len(*podConfig.ConnectPeers) > 0 {

		err := fmt.Errorf(
			"%s: the --addpeer and --connect options can not be mixed",
			funcName)

		log <- cl.Error{err}

		return err
	}

	// --proxy or --connect without --listen disables listening.
	log <- cl.Debug{"checking proxy/conneect for disabling listening"}
	if (*podConfig.Proxy != "" || len(*podConfig.ConnectPeers) > 0) &&
		len(*podConfig.Listeners) == 0 {

		*podConfig.DisableListen = true
	}

	// Connect means no DNS seeding.
	if len(*podConfig.ConnectPeers) > 0 {

		*podConfig.DisableDNSSeed = true
	}

	// Add the default listener if none were specified. The default listener is all addresses on the listen port for the network we are to connect to.
	log <- cl.Debug{"checking if listener was set"}
	if len(*podConfig.Listeners) == 0 {

		*podConfig.Listeners = []string{

			net.JoinHostPort("127.0.0.1", activeNetParams.DefaultPort),
		}

	}

	// Check to make sure limited and admin users don't have the same username
	log <- cl.Debug{"checking admin and limited username is different"}
	if *podConfig.Username != "" &&
		*podConfig.Username == *podConfig.LimitUser {

		str := "%s: --username and --limituser must not specify the same username"
		err := fmt.Errorf(str, funcName)

		log <- cl.Error{err}

		return err
	}

	// Check to make sure limited and admin users don't have the same password
	log <- cl.Debug{"checking limited and admin passwords are not the same"}
	if *podConfig.Password != "" &&
		*podConfig.Password == *podConfig.LimitPass {

		str := "%s: --password and --limitpass must not specify the same password"
		err := fmt.Errorf(str, funcName)

		log <- cl.Error{err}

		// fmt.Fprintln(os.Stderr, usageMessage)
		return err
	}

	// The RPC server is disabled if no username or password is provided.
	log <- cl.Debug{"checking rpc server has a login enabled"}
	if (*podConfig.Username == "" || *podConfig.Password == "") &&
		(*podConfig.LimitUser == "" || *podConfig.LimitPass == "") {

		*podConfig.DisableRPC = true
	}

	if *podConfig.DisableRPC {

		log <- cl.Inf("RPC service is disabled")

	}

	log <- cl.Debug{"checking rpc server has listeners set"}
	if !*podConfig.DisableRPC && len(*podConfig.RPCListeners) == 0 {

		log <- cl.Debug{"looking up default listener"}
		addrs, err := net.LookupHost(node.DefaultRPCListener)

		if err != nil {

			log <- cl.Debug{err}

			return err
		}

		*podConfig.RPCListeners = make([]string, 0, len(addrs))

		log <- cl.Debug{"setting listeners"}
		for _, addr := range addrs {

			addr = net.JoinHostPort(addr, activeNetParams.RPCClientPort)
			*podConfig.RPCListeners = append(*podConfig.RPCListeners, addr)
		}

	}

	log <- cl.Debug{"checking rpc max concurrent requests"}
	if *podConfig.RPCMaxConcurrentReqs < 0 {

		str := "%s: The rpcmaxwebsocketconcurrentrequests option may not be less than 0 -- parsed [%d]"
		err := fmt.Errorf(str, funcName, *podConfig.RPCMaxConcurrentReqs)

		log <- cl.Error{err}

		// fmt.Fprintln(os.Stderr, usageMessage)
		return err
	}

	var err error

	// Validate the the minrelaytxfee.
	log <- cl.Debug{"checking min relay tx fee"}
	StateCfg.ActiveMinRelayTxFee, err = util.NewAmount(*podConfig.MinRelayTxFee)

	if err != nil {

		str := "%s: invalid minrelaytxfee: %v"
		err := fmt.Errorf(str, funcName, err)

		log <- cl.Error{err}

		// fmt.Fprintln(os.Stderr, usageMessage)
		return err
	}

	// Limit the max block size to a sane value.
	log <- cl.Debug{"checking max block size"}
	if *podConfig.BlockMaxSize < node.BlockMaxSizeMin ||
		*podConfig.BlockMaxSize > node.BlockMaxSizeMax {

		str := "%s: The blockmaxsize option must be in between %d and %d -- parsed [%d]"
		err := fmt.Errorf(str, funcName, node.BlockMaxSizeMin,
			node.BlockMaxSizeMax, *podConfig.BlockMaxSize)

		log <- cl.Error{err}

		// fmt.Fprintln(os.Stderr, usageMessage)
		return err
	}

	// Limit the max block weight to a sane value.
	log <- cl.Debug{"checking max block weight"}
	if *podConfig.BlockMaxWeight < node.BlockMaxWeightMin ||
		*podConfig.BlockMaxWeight > node.BlockMaxWeightMax {

		str := "%s: The blockmaxweight option must be in between %d and %d -- parsed [%d]"
		err := fmt.Errorf(str, funcName, node.BlockMaxWeightMin,
			node.BlockMaxWeightMax, *podConfig.BlockMaxWeight)

		log <- cl.Error{err}

		// fmt.Fprintln(os.Stderr, usageMessage)
		return err
	}

	// Limit the max orphan count to a sane vlue.
	log <- cl.Debug{"checking max orphan limit"}
	if *podConfig.MaxOrphanTxs < 0 {

		str := "%s: The maxorphantx option may not be less than 0 -- parsed [%d]"
		err := fmt.Errorf(str, funcName, *podConfig.MaxOrphanTxs)

		log <- cl.Error{err}

		// fmt.Fprintln(os.Stderr, usageMessage)
		return err
	}

	// Limit the block priority and minimum block sizes to max block size.
	log <- cl.Debug{"checking validating block priority and minimium size/weight"}
	*podConfig.BlockPrioritySize = int(minUint32(
		uint32(*podConfig.BlockPrioritySize),
		uint32(*podConfig.BlockMaxSize)))
	*podConfig.BlockMinSize = int(minUint32(
		uint32(*podConfig.BlockMinSize),
		uint32(*podConfig.BlockMaxSize)))
	*podConfig.BlockMinWeight = int(minUint32(
		uint32(*podConfig.BlockMinWeight),
		uint32(*podConfig.BlockMaxWeight)))

	switch {

	// If the max block size isn't set, but the max weight is, then we'll set the limit for the max block size to a safe limit so weight takes precedence.
	case *podConfig.BlockMaxSize == node.DefaultBlockMaxSize &&
		*podConfig.BlockMaxWeight != node.DefaultBlockMaxWeight:
		*podConfig.BlockMaxSize = blockchain.MaxBlockBaseSize - 1000

	// If the max block weight isn't set, but the block size is, then we'll scale the set weight accordingly based on the max block size value.
	case *podConfig.BlockMaxSize != node.DefaultBlockMaxSize &&
		*podConfig.BlockMaxWeight == node.DefaultBlockMaxWeight:
		*podConfig.BlockMaxWeight = *podConfig.BlockMaxSize * blockchain.WitnessScaleFactor
	}

	// Look for illegal characters in the user agent comments.
	log <- cl.Debug{"checking user agent comments"}
	for _, uaComment := range *podConfig.UserAgentComments {

		if strings.ContainsAny(uaComment, "/:()") {

			err := fmt.Errorf("%s: The following characters must not "+
				"appear in user agent comments: '/', ':', '(', ')'",
				funcName)

			log <- cl.Err(err.Error())

			// fmt.Fprintln(os.Stderr, usageMessage)
			return err
		}

	}

	// Check mining addresses are valid and saved parsed versions.
	log <- cl.Debug{"checking mining addresses"}
	StateCfg.ActiveMiningAddrs = make([]util.Address, 0, len(*podConfig.MiningAddrs))

	for _, strAddr := range *podConfig.MiningAddrs {

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
	if (*podConfig.Generate || *podConfig.MinerListener != "") && len(*podConfig.MiningAddrs) == 0 {

		str := "%s: the generate flag is set, but there are no mining addresses specified "
		err := fmt.Errorf(str, funcName)

		log <- cl.Err(err.Error())

		fmt.Fprintln(os.Stderr, usageMessage)
		os.Exit(1)
	}

	if *podConfig.MinerPass != "" {

		StateCfg.ActiveMinerKey = fork.Argon2i([]byte(*podConfig.MinerPass))
	}

	// Add default port to all rpc listener addresses if needed and remove duplicate addresses.
	log <- cl.Debug{"checking rpc listener addresses"}
	*podConfig.RPCListeners = node.NormalizeAddresses(*podConfig.RPCListeners,
		activeNetParams.RPCClientPort)

	// Add default port to all listener addresses if needed and remove duplicate addresses.
	*podConfig.Listeners = node.NormalizeAddresses(*podConfig.Listeners,
		activeNetParams.DefaultPort)

	// Add default port to all added peer addresses if needed and remove duplicate addresses.
	*podConfig.AddPeers = node.NormalizeAddresses(*podConfig.AddPeers,
		activeNetParams.DefaultPort)
	*podConfig.ConnectPeers = node.NormalizeAddresses(*podConfig.ConnectPeers,
		activeNetParams.DefaultPort)

	// --onionproxy and not --onion are contradictory (TODO: this is kinda stupid hm? switch *and* toggle by presence of flag value, one should be enough)
	if !*podConfig.Onion && *podConfig.OnionProxy != "" {

		err := fmt.Errorf("%s: the --onionproxy and --onion options may not be activated at the same time", funcName)

		log <- cl.Error{err}

		// fmt.Fprintln(os.Stderr, usageMessage)
		return err
	}

	// Check the checkpoints for syntax errors.
	log <- cl.Debug{"checking the checkpoints"}
	StateCfg.AddedCheckpoints, err = node.ParseCheckpoints(*podConfig.AddCheckpoints)

	if err != nil {

		str := "%s: Error parsing checkpoints: %v"
		err := fmt.Errorf(str, funcName, err)

		log <- cl.Err(err.Error())

		// fmt.Fprintln(os.Stderr, usageMessage)
		return err
	}

	// Tor stream isolation requires either proxy or onion proxy to be set.
	if *podConfig.TorIsolation &&
		*podConfig.Proxy == "" &&
		*podConfig.OnionProxy == "" {

		str := "%s: Tor stream isolation requires either proxy or onionproxy to be set"
		err := fmt.Errorf(str, funcName)

		log <- cl.Error{err}

		// fmt.Fprintln(os.Stderr, usageMessage)
		return err
	}

	// Setup dial and DNS resolution (lookup) functions depending on the specified options.  The default is to use the standard net.DialTimeout function as well as the system DNS resolver.  When a proxy is specified, the dial function is set to the proxy specific dial function and the lookup is set to use tor (unless --noonion is specified in which case the system DNS resolver is used).
	log <- cl.Debug{"setting network dialer and lookup"}
	StateCfg.Dial = net.DialTimeout
	StateCfg.Lookup = net.LookupIP

	if *podConfig.Proxy != "" {

		log <- cl.Debug{"we are loading a proxy!"}

		_, _, err := net.SplitHostPort(*podConfig.Proxy)

		if err != nil {

			str := "%s: Proxy address '%s' is invalid: %v"
			err := fmt.Errorf(str, funcName, *podConfig.Proxy, err)

			log <- cl.Error{err}

			// fmt.Fprintln(os.Stderr, usageMessage)
			return err
		}

		// Tor isolation flag means proxy credentials will be overridden unless there is also an onion proxy configured in which case that one will be overridden.
		torIsolation := false

		if *podConfig.TorIsolation &&
			*podConfig.OnionProxy == "" &&
			(*podConfig.ProxyUser != "" ||
				*podConfig.ProxyPass != "") {

			torIsolation = true

			log <- cl.Warn{
				"Tor isolation set -- overriding specified proxy user credentials"}
		}

		proxy := &socks.Proxy{

			Addr:         *podConfig.Proxy,
			Username:     *podConfig.ProxyUser,
			Password:     *podConfig.ProxyPass,
			TorIsolation: torIsolation,
		}

		StateCfg.Dial = proxy.DialTimeout

		// Treat the proxy as tor and perform DNS resolution through it unless the --noonion flag is set or there is an onion-specific proxy configured.
		if *podConfig.Onion &&
			*podConfig.OnionProxy == "" {

			StateCfg.Lookup = func(host string) ([]net.IP, error) {

				return connmgr.TorLookupIP(host, *podConfig.Proxy)
			}

		}

	}

	// Setup onion address dial function depending on the specified options. The default is to use the same dial function selected above.  However, when an onion-specific proxy is specified, the onion address dial function is set to use the onion-specific proxy while leaving the normal dial function as selected above.  This allows .onion address traffic to be routed through a different proxy than normal traffic.
	log <- cl.Debug{"setting up tor proxy if enabled"}
	if *podConfig.OnionProxy != "" {

		_, _, err := net.SplitHostPort(*podConfig.OnionProxy)

		if err != nil {

			str := "%s: Onion proxy address '%s' is invalid: %v"
			err := fmt.Errorf(str, funcName, *podConfig.OnionProxy, err)

			log <- cl.Error{err}

			// fmt.Fprintln(os.Stderr, usageMessage)
			return err
		}

		// Tor isolation flag means onion proxy credentials will be overriddenode.
		if *podConfig.TorIsolation &&
			(*podConfig.OnionProxyUser != "" || *podConfig.OnionProxyPass != "") {

			log <- cl.Warn{

				"Tor isolation set - overriding specified onionproxy user credentials "}
		}

		StateCfg.Oniondial =

			func(network, addr string, timeout time.Duration) (net.Conn, error) {

				proxy := &socks.Proxy{

					Addr:         *podConfig.OnionProxy,
					Username:     *podConfig.OnionProxyUser,
					Password:     *podConfig.OnionProxyPass,
					TorIsolation: *podConfig.TorIsolation,
				}

				return proxy.DialTimeout(network, addr, timeout)
			}

		// When configured in bridge mode (both --onion and --proxy are configured), it means that the proxy configured by --proxy is not a tor proxy, so override the DNS resolution to use the onion-specific proxy.
		if *podConfig.Proxy != "" {

			StateCfg.Lookup = func(host string) ([]net.IP, error) {

				return connmgr.TorLookupIP(host, *podConfig.OnionProxy)
			}

		}

	} else {

		StateCfg.Oniondial = StateCfg.Dial
	}

	// Specifying --noonion means the onion address dial function results in an error.
	if !*podConfig.Onion {

		StateCfg.Oniondial = func(a, b string, t time.Duration) (net.Conn, error) {

			return nil, errors.New("tor has been disabled")
		}

	}

	if StateCfg.Save {

		StateCfg.Save = false
		podHandleSave()
	}

	log <- cl.Debug{"finished nodeHandle"}
	node.Main(&podConfig, activeNetParams, nil)
	return nil
}

func NormalizeStringSliceAddresses(a *cli.StringSlice, port string) {

	variable := []string(*a)
	NormalizeAddresses(
		strings.Join(variable, " "),
		port, &variable)
	*a = variable
}
