package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	n "git.parallelcoin.io/pod/cmd/node"
	w "git.parallelcoin.io/pod/cmd/wallet"
	cl "git.parallelcoin.io/pod/pkg/clog"
	"git.parallelcoin.io/pod/pkg/netparams"
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
			w.ActiveNet = &netparams.TestNet3Params
		case "simnet":
			sc.Wallet.TestNet3, sc.Wallet.SimNet = false, true
			sc.Node.TestNet3, sc.Node.SimNet, sc.Node.RegressionTest = false, true, false
			w.ActiveNet = &netparams.SimNetParams
		default:
			sc.Wallet.TestNet3, sc.Wallet.SimNet = false, false
			sc.Node.TestNet3, sc.Node.SimNet, sc.Node.RegressionTest = false, false, false
			w.ActiveNet = &netparams.MainNetParams
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
	}
	if r, ok := getIfIs(ctx, "password"); ok {
		sc.Wallet.Password = r
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
}
