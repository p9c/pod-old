package app_old

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

	"git.parallelcoin.io/dev/pod/cmd/node"
	"git.parallelcoin.io/dev/pod/cmd/node/mempool"
	blockchain "git.parallelcoin.io/dev/pod/pkg/chain"
	"git.parallelcoin.io/dev/pod/pkg/chain/fork"
	"git.parallelcoin.io/dev/pod/pkg/peer/connmgr"
	"git.parallelcoin.io/dev/pod/pkg/util"
	cl "git.parallelcoin.io/dev/pod/pkg/util/cl"
	"github.com/btcsuite/go-socks/socks"
	"github.com/davecgh/go-spew/spew"
	"github.com/tucnak/climax"
)

// DefaultNodeConfig is the default configuration for node
func DefaultNodeConfig(datadir string) *NodeCfg {

	user := GenKey()
	pass := GenKey()
	params := node.MainNetParams

	switch node.ActiveNetParams.Name {

	case "testnet3":
		params = node.TestNet3Params
	case "simnet":
		params = node.SimNetParams
	}

	appdir := filepath.Join(datadir, "node")

	return &NodeCfg{

		Node: &node.Config{

			RPCUser:      user,
			RPCPass:      pass,
			Listeners:    []string{node.DefaultListener},
			RPCListeners: []string{node.DefaultRPCListener},
			DebugLevel:   "info",
			ConfigFile: filepath.Join(appdir,
				node.DefaultConfigFilename),
			MaxPeers:             node.DefaultMaxPeers,
			BanDuration:          node.DefaultBanDuration,
			BanThreshold:         node.DefaultBanThreshold,
			RPCMaxClients:        node.DefaultMaxRPCClients,
			RPCMaxWebsockets:     node.DefaultMaxRPCWebsockets,
			RPCMaxConcurrentReqs: node.DefaultMaxRPCConcurrentReqs,
			DataDir:              appdir,
			LogDir:               appdir,
			DbType:               node.DefaultDbType,
			RPCKey:               filepath.Join(datadir, "rpc.key"),
			RPCCert:              filepath.Join(datadir, "rpc.cert"),
			MinRelayTxFee:        mempool.DefaultMinRelayTxFee.ToDUO(),
			FreeTxRelayLimit:     node.DefaultFreeTxRelayLimit,
			TrickleInterval:      node.DefaultTrickleInterval,
			BlockMinSize:         node.DefaultBlockMinSize,
			BlockMaxSize:         node.DefaultBlockMaxSize,
			BlockMinWeight:       node.DefaultBlockMinWeight,
			BlockMaxWeight:       node.DefaultBlockMaxWeight,
			BlockPrioritySize:    mempool.DefaultBlockPrioritySize,
			MaxOrphanTxs:         node.DefaultMaxOrphanTransactions,
			SigCacheMaxSize:      node.DefaultSigCacheMaxSize,
			Generate:             node.DefaultGenerate,
			GenThreads:           1,
			TxIndex:              node.DefaultTxIndex,
			AddrIndex:            node.DefaultAddrIndex,
			Algo:                 node.DefaultAlgo,
		},

		LogLevels: GetDefaultLogLevelsConfig(),
		params:    &params,
	}

}

// WriteDefaultNodeConfig creates a default config and writes it to the requested location
func WriteDefaultNodeConfig(
	datadir string,

) {

	log <- cl.Dbg("writing default config")

	defCfg := DefaultNodeConfig(datadir)
	j, err := json.MarshalIndent(defCfg, "", "  ")

	if err != nil {

		log <- cl.Error{`marshalling default app config file: "`, err, `"`}

		log <- cl.Err(spew.Sdump(defCfg))

		return
	}

	j = append(j, '\n')

	log <- cl.Tracef{

		"JSON formatted config file\n%s",
		j,
	}

	EnsureDir(defCfg.Node.ConfigFile)
	err = ioutil.WriteFile(defCfg.Node.ConfigFile, j, 0600)

	if err != nil {

		log <- cl.Error{"writing default app config file:", err.Error()}

		return
	}

	// if we are writing default config we also want to use it
	NodeConfig = defCfg
}

// WriteNodeConfig writes the current config to the requested location
func WriteNodeConfig(
	c *NodeCfg,

) {

	log <- cl.Dbg("writing config")

	j, err := json.MarshalIndent(c, "", "  ")

	if err != nil {

		log <- cl.Error{`marshalling default app config file: "`, err, `"`}

		log <- cl.Err(spew.Sdump(c))

		return
	}

	j = append(j, '\n')

	log <- cl.Tracef{

		"JSON formatted config file\n%s",
		j,
	}

	EnsureDir(c.Node.ConfigFile)
	err = ioutil.WriteFile(c.Node.ConfigFile, j, 0600)

	if err != nil {

		log <- cl.Error{"writing default app config file:", err.Error()}

		return
	}

}
func configNode(
	nc *node.Config,
	ctx *climax.Context,
	cfgFile string,

) int {

	var err error

	if r, ok := getIfIs(ctx, "datadir"); ok {

		nc.DataDir = filepath.Join(node.CleanAndExpandPath(r), "node")
		nc.ConfigFile = filepath.Join(node.CleanAndExpandPath(r), "conf.json")
	}

	if r, ok := getIfIs(ctx, "addpeers"); ok {

		NormalizeAddresses(r, node.DefaultPort, &nc.AddPeers)
	}

	if r, ok := getIfIs(ctx, "connectpeers"); ok {

		NormalizeAddresses(r, node.DefaultPort, &nc.ConnectPeers)
	}

	if r, ok := getIfIs(ctx, "disablelisten"); ok {

		nc.DisableListen = r == "true"
	}

	if r, ok := getIfIs(ctx, "listeners"); ok {

		NormalizeAddresses(r, node.DefaultPort, &nc.Listeners)
	}

	if r, ok := getIfIs(ctx, "maxpeers"); ok {

		if err := ParseInteger(r, "maxpeers", &nc.MaxPeers); err != nil {

			log <- cl.Wrn(err.Error())

		}

	}

	if r, ok := getIfIs(ctx, "disablebanning"); ok {

		nc.DisableBanning = r == "true"
	}

	if r, ok := getIfIs(ctx, "banduration"); ok {

		if err := ParseDuration(r, "banduration", &nc.BanDuration); err != nil {

			log <- cl.Wrn(err.Error())

		}

	}

	if r, ok := getIfIs(ctx, "banthreshold"); ok {

		var bt int

		if err := ParseInteger(r, "banthtreshold", &bt); err != nil {

			log <- cl.Wrn(err.Error())

		} else {

			nc.BanThreshold = uint32(bt)
		}

	}

	if r, ok := getIfIs(ctx, "whitelists"); ok {

		NormalizeAddresses(r, node.DefaultPort, &nc.Whitelists)
	}

	if r, ok := getIfIs(ctx, "rpcuser"); ok {

		nc.RPCUser = r
	}

	if r, ok := getIfIs(ctx, "rpcpass"); ok {

		nc.RPCPass = r
	}

	if r, ok := getIfIs(ctx, "rpclimituser"); ok {

		nc.RPCLimitUser = r
	}

	if r, ok := getIfIs(ctx, "rpclimitpass"); ok {

		nc.RPCLimitPass = r
	}

	if r, ok := getIfIs(ctx, "rpclisteners"); ok {

		NormalizeAddresses(r, node.DefaultRPCPort, &nc.RPCListeners)
	}

	if r, ok := getIfIs(ctx, "rpccert"); ok {

		nc.RPCCert = node.CleanAndExpandPath(r)
	}

	if r, ok := getIfIs(ctx, "rpckey"); ok {

		nc.RPCKey = node.CleanAndExpandPath(r)
	}

	if r, ok := getIfIs(ctx, "tls"); ok {

		nc.TLS = r == "true"
	}

	if r, ok := getIfIs(ctx, "disablednsseed"); ok {

		nc.DisableDNSSeed = r == "true"
	}

	if r, ok := getIfIs(ctx, "externalips"); ok {

		NormalizeAddresses(r, node.DefaultPort, &nc.ExternalIPs)
	}

	if r, ok := getIfIs(ctx, "proxy"); ok {

		NormalizeAddress(r, "9050", &nc.Proxy)
	}

	if r, ok := getIfIs(ctx, "proxyuser"); ok {

		nc.ProxyUser = r
	}

	if r, ok := getIfIs(ctx, "proxypass"); ok {

		nc.ProxyPass = r
	}

	if r, ok := getIfIs(ctx, "onion"); ok {

		NormalizeAddress(r, "9050", &nc.OnionProxy)
	}

	if r, ok := getIfIs(ctx, "onionuser"); ok {

		nc.OnionProxyUser = r
	}

	if r, ok := getIfIs(ctx, "onionpass"); ok {

		nc.OnionProxyPass = r
	}

	if r, ok := getIfIs(ctx, "noonion"); ok {

		nc.NoOnion = r == "true"
	}

	if r, ok := getIfIs(ctx, "torisolation"); ok {

		nc.TorIsolation = r == "true"
	}

	if r, ok := getIfIs(ctx, "network"); ok {

		switch r {

		case "testnet":

			fork.IsTestnet = true

			nc.TestNet3, nc.RegressionTest, nc.SimNet = true, false, false
			NodeConfig.params = &node.TestNet3Params
		case "regtest":
			nc.TestNet3, nc.RegressionTest, nc.SimNet = false, true, false
			NodeConfig.params = &node.RegressionNetParams
		case "simnet":
			nc.TestNet3, nc.RegressionTest, nc.SimNet = false, false, true
			NodeConfig.params = &node.SimNetParams
		default:
			nc.TestNet3, nc.RegressionTest, nc.SimNet = false, false, false
			NodeConfig.params = &node.MainNetParams
		}

		log <- cl.Debug{NodeConfig.params.Name, r}

	}

	if r, ok := getIfIs(ctx, "addcheckpoints"); ok {

		nc.AddCheckpoints = strings.Split(r, " ")
	}

	if r, ok := getIfIs(ctx, "disablecheckpoints"); ok {

		nc.DisableCheckpoints = r == "true"
	}

	if r, ok := getIfIs(ctx, "dbtype"); ok {

		nc.DbType = r
	}

	if r, ok := getIfIs(ctx, "profile"); ok {

		var p int

		if err = ParseInteger(r, "profile", &p); err == nil {

			nc.Profile = fmt.Sprint(p)
		}

	}

	if r, ok := getIfIs(ctx, "cpuprofile"); ok {

		nc.CPUProfile = r
	}

	if r, ok := getIfIs(ctx, "upnp"); ok {

		nc.Upnp = r == "true"
	}

	if r, ok := getIfIs(ctx, "minrelaytxfee"); ok {

		if err := ParseFloat(r, "minrelaytxfee", &nc.MinRelayTxFee); err != nil {

			log <- cl.Wrn(err.Error())

		}

	}

	if r, ok := getIfIs(ctx, "freetxrelaylimit"); ok {

		if err := ParseFloat(r, "freetxrelaylimit", &nc.FreeTxRelayLimit); err != nil {

			log <- cl.Wrn(err.Error())

		}

	}

	if r, ok := getIfIs(ctx, "norelaypriority"); ok {

		nc.NoRelayPriority = r == "true"
	}

	if r, ok := getIfIs(ctx, "trickleinterval"); ok {

		if err := ParseDuration(r, "trickleinterval", &nc.TrickleInterval); err != nil {

			log <- cl.Wrn(err.Error())

		}

	}

	if r, ok := getIfIs(ctx, "maxorphantxs"); ok {

		if err := ParseInteger(r, "maxorphantxs", &nc.MaxOrphanTxs); err != nil {

			log <- cl.Wrn(err.Error())

		}

	}

	if r, ok := getIfIs(ctx, "algo"); ok {

		nc.Algo = r
	}

	if r, ok := getIfIs(ctx, "generate"); ok {

		nc.Generate = r == "true"
	}

	if r, ok := getIfIs(ctx, "genthreads"); ok {

		var gt int

		if err := ParseInteger(r, "genthreads", &gt); err != nil {

			log <- cl.Wrn(err.Error())

		} else {

			nc.GenThreads = int32(gt)
		}

	}

	if r, ok := getIfIs(ctx, "miningaddrs"); ok {

		nc.MiningAddrs = strings.Split(r, " ")
	}

	if r, ok := getIfIs(ctx, "minerlistener"); ok {

		NormalizeAddress(r, node.DefaultRPCPort, &nc.MinerListener)
	}

	if r, ok := getIfIs(ctx, "minerpass"); ok {

		nc.MinerPass = r
	}

	if r, ok := getIfIs(ctx, "blockminsize"); ok {

		if err := ParseUint32(r, "blockminsize", &nc.BlockMinSize); err != nil {

			log <- cl.Wrn(err.Error())

		}

	}

	if r, ok := getIfIs(ctx, "blockmaxsize"); ok {

		if err := ParseUint32(r, "blockmaxsize", &nc.BlockMaxSize); err != nil {

			log <- cl.Wrn(err.Error())

		}

	}

	if r, ok := getIfIs(ctx, "blockminweight"); ok {

		if err := ParseUint32(r, "blockminweight", &nc.BlockMinWeight); err != nil {

			log <- cl.Wrn(err.Error())

		}

	}

	if r, ok := getIfIs(ctx, "blockmaxweight"); ok {

		if err := ParseUint32(r, "blockmaxweight", &nc.BlockMaxWeight); err != nil {

			log <- cl.Wrn(err.Error())

		}

	}

	if r, ok := getIfIs(ctx, "blockprioritysize"); ok {

		if err := ParseUint32(r, "blockmaxweight", &nc.BlockPrioritySize); err != nil {

			log <- cl.Wrn(err.Error())

		}

	}

	if r, ok := getIfIs(ctx, "uacomment"); ok {

		nc.UserAgentComments = strings.Split(r, " ")
	}

	if r, ok := getIfIs(ctx, "nopeerbloomfilters"); ok {

		nc.NoPeerBloomFilters = r == "true"
	}

	if r, ok := getIfIs(ctx, "nocfilters"); ok {

		nc.NoCFilters = r == "true"
	}

	if ctx.Is("dropcfindex") {

		nc.DropCfIndex = true
	}

	if r, ok := getIfIs(ctx, "sigcachemaxsize"); ok {

		var scms int

		if err := ParseInteger(r, "sigcachemaxsize", &scms); err != nil {

			log <- cl.Wrn(err.Error())

		} else {

			nc.SigCacheMaxSize = uint(scms)
		}

	}

	if r, ok := getIfIs(ctx, "blocksonly"); ok {

		nc.BlocksOnly = r == "true"
	}

	if r, ok := getIfIs(ctx, "txindex"); ok {

		nc.TxIndex = r == "true"
	}

	if ctx.Is("droptxindex") {

		nc.DropTxIndex = true
	}

	if ctx.Is("addrindex") {

		r, _ := ctx.Get("addrindex")
		nc.AddrIndex = r == "true"
	}

	if ctx.Is("dropaddrindex") {

		nc.DropAddrIndex = true
	}

	if r, ok := getIfIs(ctx, "relaynonstd"); ok {

		nc.RelayNonStd = r == "true"
	}

	if r, ok := getIfIs(ctx, "rejectnonstd"); ok {

		nc.RejectNonStd = r == "true"
	}

	// finished configuration

	SetLogging(ctx)

	if ctx.Is("save") {

		log <- cl.Infof{

			"saving config file to %s",
			cfgFile,
		}

		j, err := json.MarshalIndent(NodeConfig, "", "  ")

		if err != nil {

			log <- cl.Error{

				"saving config file:",
				err.Error(),
			}

		}

		j = append(j, '\n')

		log <- cl.Tracef{

			"JSON formatted config file\n%s",
			j,
		}

		err = ioutil.WriteFile(cfgFile, j, 0600)

		if err != nil {

			log <- cl.Error{"writing app config file:", err.Error()}

		}

	}

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

	if nc.RegressionTest && len(nc.AddPeers) > 0 {

		nc.AddPeers = nil
	}

	// Set the mining algorithm correctly, default to random if unrecognised

	switch nc.Algo {

	case "blake14lr", "cryptonight7v2", "keccak", "lyra2rev2", "scrypt", "skein", "x11", "stribog", "random", "easy":
	default:
		nc.Algo = "random"
	}

	relayNonStd := NodeConfig.params.RelayNonStdTxs
	funcName := "loadConfig"

	switch {

	case nc.RelayNonStd && nc.RejectNonStd:
		str := "%s: rejectnonstd and relaynonstd cannot be used together -- choose only one"
		err := fmt.Errorf(str, funcName)
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, usageMessage)
		return 1
	case nc.RejectNonStd:
		relayNonStd = false
	case nc.RelayNonStd:
		relayNonStd = true
	}

	nc.RelayNonStd = relayNonStd

	// Append the network type to the data directory so it is "namespaced" per network.  In addition to the block database, there are other pieces of data that are saved to disk such as address manager state. All data is specific to a network, so namespacing the data directory means each individual piece of serialized data does not have to worry about changing names per network and such.
	nc.DataDir = node.CleanAndExpandPath(nc.DataDir)

	log <- cl.Debug{"netname", NodeConfig.params.Name, node.NetName(NodeConfig.params)}

	nc.DataDir = filepath.Join(nc.DataDir, node.NetName(NodeConfig.params))

	// Validate database type.

	if !node.ValidDbType(nc.DbType) {

		str := "%s: The specified database type [%v] is invalid -- " +
			"supported types %v"
		err := fmt.Errorf(str, funcName, nc.DbType, node.KnownDbTypes)
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, usageMessage)
		return 1
	}

	// Validate profile port number

	if nc.Profile != "" {

		profilePort, err := strconv.Atoi(nc.Profile)

		if err != nil || profilePort < 1024 || profilePort > 65535 {

			str := "%s: The profile port must be between 1024 and 65535"
			err := fmt.Errorf(str, funcName)
			fmt.Fprintln(os.Stderr, err)
			fmt.Fprintln(os.Stderr, usageMessage)
			return 1
		}

	}

	// Don't allow ban durations that are too short.

	if nc.BanDuration < time.Second {

		str := "%s: The banduration option may not be less than 1s -- parsed [%v]"
		err := fmt.Errorf(str, funcName, nc.BanDuration)
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, usageMessage)
		return 1
	}

	// Validate any given whitelisted IP addresses and networks.

	if len(nc.Whitelists) > 0 {

		var ip net.IP
		StateCfg.ActiveWhitelists = make([]*net.IPNet, 0, len(nc.Whitelists))

		for _, addr := range nc.Whitelists {

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

	if len(nc.AddPeers) > 0 && len(nc.ConnectPeers) > 0 {

		str := "%s: the --addpeer and --connect options can not be " +
			"mixed"
		err := fmt.Errorf(str, funcName)

		log <- cl.Err(err.Error())

		fmt.Fprintln(os.Stderr, usageMessage)
	}

	// --proxy or --connect without --listen disables listening.

	if (nc.Proxy != "" || len(nc.ConnectPeers) > 0) &&

		len(nc.Listeners) == 0 {

		nc.DisableListen = true
	}

	// Connect means no DNS seeding.

	if len(nc.ConnectPeers) > 0 {

		nc.DisableDNSSeed = true
	}

	// Add the default listener if none were specified. The default listener is all addresses on the listen port for the network we are to connect to.

	if len(nc.Listeners) == 0 {

		nc.Listeners = []string{

			net.JoinHostPort("", NodeConfig.params.DefaultPort),
		}

	}

	// Check to make sure limited and admin users don't have the same username

	if nc.RPCUser == nc.RPCLimitUser && nc.RPCUser != "" {

		str := "%s: --rpcuser and --rpclimituser must not specify the same username"
		err := fmt.Errorf(str, funcName)

		log <- cl.Err(err.Error())

		fmt.Fprintln(os.Stderr, usageMessage)
		return 1
	}

	// Check to make sure limited and admin users don't have the same password

	if nc.RPCPass == nc.RPCLimitPass && nc.RPCPass != "" {

		str := "%s: --rpcpass and --rpclimitpass must not specify the " +
			"same password"
		err := fmt.Errorf(str, funcName)

		log <- cl.Err(err.Error())

		fmt.Fprintln(os.Stderr, usageMessage)
		return 1
	}

	// The RPC server is disabled if no username or password is provided.

	if (nc.RPCUser == "" || nc.RPCPass == "") &&

		(nc.RPCLimitUser == "" || nc.RPCLimitPass == "") {

		nc.DisableRPC = true
	}

	if nc.DisableRPC {

		log <- cl.Inf("RPC service is disabled")

	}

	// Default RPC to listen on localhost only.

	if !nc.DisableRPC && len(nc.RPCListeners) == 0 {

		addrs, err := net.LookupHost(node.DefaultRPCListener)

		if err != nil {

			log <- cl.Err(err.Error())

			return 1
		}

		nc.RPCListeners = make([]string, 0, len(addrs))

		for _, addr := range addrs {

			addr = net.JoinHostPort(addr, NodeConfig.params.RPCPort)
			nc.RPCListeners = append(nc.RPCListeners, addr)
		}

	}

	if nc.RPCMaxConcurrentReqs < 0 {

		str := "%s: The rpcmaxwebsocketconcurrentrequests option may not be less than 0 -- parsed [%d]"
		err := fmt.Errorf(str, funcName, nc.RPCMaxConcurrentReqs)

		log <- cl.Err(err.Error())

		fmt.Fprintln(os.Stderr, usageMessage)
		return 1
	}

	// Validate the the minrelaytxfee.
	StateCfg.ActiveMinRelayTxFee, err = util.NewAmount(nc.MinRelayTxFee)

	if err != nil {

		str := "%s: invalid minrelaytxfee: %v"
		err := fmt.Errorf(str, funcName, err)

		log <- cl.Err(err.Error())

		fmt.Fprintln(os.Stderr, usageMessage)
		return 1
	}

	// Limit the max block size to a sane value.

	if nc.BlockMaxSize < node.BlockMaxSizeMin || nc.BlockMaxSize >

		node.BlockMaxSizeMax {

		str := "%s: The blockmaxsize option must be in between %d and %d -- parsed [%d]"
		err := fmt.Errorf(str, funcName, node.BlockMaxSizeMin,
			node.BlockMaxSizeMax, nc.BlockMaxSize)

		log <- cl.Err(err.Error())

		fmt.Fprintln(os.Stderr, usageMessage)
		return 1
	}

	// Limit the max block weight to a sane value.

	if nc.BlockMaxWeight < node.BlockMaxWeightMin ||

		nc.BlockMaxWeight > node.BlockMaxWeightMax {

		str := "%s: The blockmaxweight option must be in between %d and %d -- parsed [%d]"
		err := fmt.Errorf(str, funcName, node.BlockMaxWeightMin,
			node.BlockMaxWeightMax, nc.BlockMaxWeight)

		log <- cl.Err(err.Error())

		fmt.Fprintln(os.Stderr, usageMessage)
		return 1
	}

	// Limit the max orphan count to a sane vlue.

	if nc.MaxOrphanTxs < 0 {

		str := "%s: The maxorphantx option may not be less than 0 -- parsed [%d]"
		err := fmt.Errorf(str, funcName, nc.MaxOrphanTxs)

		log <- cl.Err(err.Error())

		fmt.Fprintln(os.Stderr, usageMessage)
		return 1
	}

	// Limit the block priority and minimum block sizes to max block size.
	nc.BlockPrioritySize = minUint32(nc.BlockPrioritySize, nc.BlockMaxSize)
	nc.BlockMinSize = minUint32(nc.BlockMinSize, nc.BlockMaxSize)
	nc.BlockMinWeight = minUint32(nc.BlockMinWeight, nc.BlockMaxWeight)

	switch {

	// If the max block size isn't set, but the max weight is, then we'll set the limit for the max block size to a safe limit so weight takes precedence.
	case nc.BlockMaxSize == node.DefaultBlockMaxSize &&
		nc.BlockMaxWeight != node.DefaultBlockMaxWeight:
		nc.BlockMaxSize = blockchain.MaxBlockBaseSize - 1000

	// If the max block weight isn't set, but the block size is, then we'll scale the set weight accordingly based on the max block size value.
	case nc.BlockMaxSize != node.DefaultBlockMaxSize &&
		nc.BlockMaxWeight == node.DefaultBlockMaxWeight:
		nc.BlockMaxWeight = nc.BlockMaxSize * blockchain.WitnessScaleFactor
	}

	// Look for illegal characters in the user agent comments.

	for _, uaComment := range nc.UserAgentComments {

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

	if nc.TxIndex && nc.DropTxIndex {

		err := fmt.Errorf("%s: the --txindex and --droptxindex options may  not be activated at the same time",
			funcName)

		log <- cl.Err(err.Error())

		fmt.Fprintln(os.Stderr, usageMessage)
		return 1

	}

	// --addrindex and --dropaddrindex do not mix.

	if nc.AddrIndex && nc.DropAddrIndex {

		err := fmt.Errorf("%s: the --addrindex and --dropaddrindex "+
			"options may not be activated at the same time",
			funcName)

		log <- cl.Err(err.Error())

		fmt.Fprintln(os.Stderr, usageMessage)
		return 1
	}

	// Check mining addresses are valid and saved parsed versions.
	StateCfg.ActiveMiningAddrs = make([]util.Address, 0, len(nc.MiningAddrs))

	for _, strAddr := range nc.MiningAddrs {

		addr, err := util.DecodeAddress(strAddr, NodeConfig.params.Params)

		if err != nil {

			str := "%s: mining address '%s' failed to decode: %v"
			err := fmt.Errorf(str, funcName, strAddr, err)

			log <- cl.Err(err.Error())

			fmt.Fprintln(os.Stderr, usageMessage)
			return 1
		}

		if !addr.IsForNet(NodeConfig.params.Params) {

			str := "%s: mining address '%s' is on the wrong network"
			err := fmt.Errorf(str, funcName, strAddr)

			log <- cl.Err(err.Error())

			fmt.Fprintln(os.Stderr, usageMessage)
			return 1
		}

		StateCfg.ActiveMiningAddrs = append(StateCfg.ActiveMiningAddrs, addr)
	}

	// Ensure there is at least one mining address when the generate flag is set.

	if (nc.Generate || nc.MinerListener != "") && len(nc.MiningAddrs) == 0 {

		str := "%s: the generate flag is set, but there are no mining addresses specified "
		err := fmt.Errorf(str, funcName)

		log <- cl.Err(err.Error())

		fmt.Fprintln(os.Stderr, usageMessage)
		os.Exit(1)
	}

	if nc.MinerPass != "" {

		StateCfg.ActiveMinerKey = fork.Argon2i([]byte(nc.MinerPass))
	}

	// Add default port to all listener addresses if needed and remove duplicate addresses.
	nc.Listeners = node.NormalizeAddresses(nc.Listeners,
		NodeConfig.params.DefaultPort)

	// Add default port to all rpc listener addresses if needed and remove duplicate addresses.
	nc.RPCListeners = node.NormalizeAddresses(nc.RPCListeners,
		NodeConfig.params.RPCPort)

	if !nc.DisableRPC && !nc.TLS {

		for _, addr := range nc.RPCListeners {

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
	nc.AddPeers = node.NormalizeAddresses(nc.AddPeers,
		NodeConfig.params.DefaultPort)
	nc.ConnectPeers = node.NormalizeAddresses(nc.ConnectPeers,
		NodeConfig.params.DefaultPort)

	// --noonion and --onion do not mix.

	if nc.NoOnion && nc.OnionProxy != "" {

		err := fmt.Errorf("%s: the --noonion and --onion options may not be activated at the same time", funcName)

		log <- cl.Err(err.Error())

		fmt.Fprintln(os.Stderr, usageMessage)
		return 1
	}

	// Check the checkpoints for syntax errors.
	StateCfg.AddedCheckpoints, err = node.ParseCheckpoints(nc.AddCheckpoints)

	if err != nil {

		str := "%s: Error parsing checkpoints: %v"
		err := fmt.Errorf(str, funcName, err)

		log <- cl.Err(err.Error())

		fmt.Fprintln(os.Stderr, usageMessage)
		return 1
	}

	// Tor stream isolation requires either proxy or onion proxy to be set.

	if nc.TorIsolation && nc.Proxy == "" && nc.OnionProxy == "" {

		str := "%s: Tor stream isolation requires either proxy or onionproxy to be set"
		err := fmt.Errorf(str, funcName)

		log <- cl.Err(err.Error())

		fmt.Fprintln(os.Stderr, usageMessage)
		return 1
	}

	// Setup dial and DNS resolution (lookup) functions depending on the specified options.  The default is to use the standard net.DialTimeout function as well as the system DNS resolver.  When a proxy is specified, the dial function is set to the proxy specific dial function and the lookup is set to use tor (unless --noonion is specified in which case the system DNS resolver is used).
	StateCfg.Dial = net.DialTimeout
	StateCfg.Lookup = net.LookupIP

	if nc.Proxy != "" {

		_, _, err := net.SplitHostPort(nc.Proxy)

		if err != nil {

			str := "%s: Proxy address '%s' is invalid: %v"
			err := fmt.Errorf(str, funcName, nc.Proxy, err)

			log <- cl.Err(err.Error())

			fmt.Fprintln(os.Stderr, usageMessage)
			return 1
		}

		// Tor isolation flag means proxy credentials will be overridden unless there is also an onion proxy configured in which case that one will be overriddenode.
		torIsolation := false

		if nc.TorIsolation && nc.OnionProxy == "" &&

			(nc.ProxyUser != "" || nc.ProxyPass != "") {

			torIsolation = true
			fmt.Fprintln(os.Stderr, "Tor isolation set -- "+
				"overriding specified proxy user credentials")
		}

		proxy := &socks.Proxy{

			Addr:         nc.Proxy,
			Username:     nc.ProxyUser,
			Password:     nc.ProxyPass,
			TorIsolation: torIsolation,
		}

		StateCfg.Dial = proxy.DialTimeout

		// Treat the proxy as tor and perform DNS resolution through it unless the --noonion flag is set or there is an onion-specific proxy configured.

		if !nc.NoOnion && nc.OnionProxy == "" {

			StateCfg.Lookup = func(host string) ([]net.IP, error) {

				return connmgr.TorLookupIP(host, nc.Proxy)
			}

		}

	}

	// Setup onion address dial function depending on the specified options. The default is to use the same dial function selected above.  However, when an onion-specific proxy is specified, the onion address dial function is set to use the onion-specific proxy while leaving the normal dial function as selected above.  This allows .onion address traffic to be routed through a different proxy than normal traffic.

	if nc.OnionProxy != "" {

		_, _, err := net.SplitHostPort(nc.OnionProxy)

		if err != nil {

			str := "%s: Onion proxy address '%s' is invalid: %v"
			err := fmt.Errorf(str, funcName, nc.OnionProxy, err)

			log <- cl.Err(err.Error())

			fmt.Fprintln(os.Stderr, usageMessage)
			return 1
		}

		// Tor isolation flag means onion proxy credentials will be overriddenode.

		if nc.TorIsolation &&

			(nc.OnionProxyUser != "" || nc.OnionProxyPass != "") {

			fmt.Fprintln(os.Stderr, "Tor isolation set -- "+
				"overriding specified onionproxy user "+
				"credentials ")
		}

		StateCfg.Oniondial = func(network, addr string, timeout time.Duration) (net.Conn, error) {

			proxy := &socks.Proxy{

				Addr:         nc.OnionProxy,
				Username:     nc.OnionProxyUser,
				Password:     nc.OnionProxyPass,
				TorIsolation: nc.TorIsolation,
			}

			return proxy.DialTimeout(network, addr, timeout)
		}

		// When configured in bridge mode (both --onion and --proxy are configured), it means that the proxy configured by --proxy is not a tor proxy, so override the DNS resolution to use the onion-specific proxy.

		if nc.Proxy != "" {

			StateCfg.Lookup = func(host string) ([]net.IP, error) {

				return connmgr.TorLookupIP(host, nc.OnionProxy)
			}

		}

	} else {

		StateCfg.Oniondial = StateCfg.Dial
	}

	// Specifying --noonion means the onion address dial function results in an error.

	if nc.NoOnion {

		StateCfg.Oniondial = func(a, b string, t time.Duration) (net.Conn, error) {

			return nil, errors.New("tor has been disabled")
		}

	}

	return 0
}
