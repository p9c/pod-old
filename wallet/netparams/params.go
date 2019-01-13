package netparams

import "github.com/parallelcointeam/pod/chaincfg"

// Params is used to group parameters for various networks such as the main
// network and test networks.
type Params struct {
	*chaincfg.Params
	RPCClientPort string
	RPCServerPort string
}

// MainNetParams contains parameters specific running btcwallet and
// pod on the main network (wire.MainNet).
var MainNetParams = Params{
	Params:        &chaincfg.MainNetParams,
	RPCClientPort: "11048",
	RPCServerPort: "11046",
}

// TestNet3Params contains parameters specific running btcwallet and
// pod on the test network (version 3) (wire.TestNet3).
var TestNet3Params = Params{
	Params:        &chaincfg.TestNet3Params,
	RPCClientPort: "21048",
	RPCServerPort: "21046",
}

// SimNetParams contains parameters specific to the simulation test network
// (wire.SimNet).
var SimNetParams = Params{
	Params:        &chaincfg.SimNetParams,
	RPCClientPort: "41048",
	RPCServerPort: "41046",
}
