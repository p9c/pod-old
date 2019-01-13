package node

import (
	"github.com/parallelcointeam/pod/chaincfg"
	"github.com/parallelcointeam/pod/wire"
)

// activeNetParams is a pointer to the parameters specific to the currently active bitcoin network.
var activeNetParams = &mainNetParams

// params is used to group parameters for various networks such as the main network and test networks.
type params struct {
	*chaincfg.Params
	RPCPort            string
	ScryptRPCPort      string
	Cryptonight7v2Port string
	Blake14lrRPCPort   string
	KeccakRPCPort      string
	Lyra2rev2RPCPort   string
	SkeinRPCPort       string
	X11RPCPort         string
	StribogRPCPort     string
}

// mainNetParams contains parameters specific to the main network (wire.MainNet).  NOTE: The RPC port is intentionally different than the reference implementation because pod does not handle wallet requests.  The separate wallet process listens on the well-known port and forwards requests it does not handle on to pod.  This approach allows the wallet process to emulate the full reference implementation RPC API.
var mainNetParams = params{
	Params:             &chaincfg.MainNetParams,
	RPCPort:            "11048",
	Blake14lrRPCPort:   "11049",
	Cryptonight7v2Port: "11050",
	KeccakRPCPort:      "11051",
	Lyra2rev2RPCPort:   "11052",
	ScryptRPCPort:      "11053",
	StribogRPCPort:     "11054",
	SkeinRPCPort:       "11055",
	X11RPCPort:         "11056",
}

// regressionNetParams contains parameters specific to the regression test network (wire.TestNet).  NOTE: The RPC port is intentionally different than the reference implementation - see the mainNetParams comment for details.
var regressionNetParams = params{
	Params:             &chaincfg.RegressionNetParams,
	RPCPort:            "31048",
	Blake14lrRPCPort:   "31049",
	Cryptonight7v2Port: "31050",
	KeccakRPCPort:      "31051",
	Lyra2rev2RPCPort:   "31052",
	ScryptRPCPort:      "31053",
	StribogRPCPort:     "31054",
	SkeinRPCPort:       "31055",
	X11RPCPort:         "31056",
}

// testNet3Params contains parameters specific to the test network (version 3) (wire.TestNet3).  NOTE: The RPC port is intentionally different than the reference implementation - see the mainNetParams comment for details.
var testNet3Params = params{
	Params:             &chaincfg.TestNet3Params,
	RPCPort:            "21048",
	Blake14lrRPCPort:   "21049",
	Cryptonight7v2Port: "21050",
	KeccakRPCPort:      "21051",
	Lyra2rev2RPCPort:   "21052",
	ScryptRPCPort:      "21053",
	StribogRPCPort:     "21054",
	SkeinRPCPort:       "21055",
	X11RPCPort:         "21056",
}

// simNetParams contains parameters specific to the simulation test network (wire.SimNet).
var simNetParams = params{
	Params:             &chaincfg.SimNetParams,
	RPCPort:            "41048",
	Blake14lrRPCPort:   "41049",
	Cryptonight7v2Port: "41050",
	KeccakRPCPort:      "41051",
	Lyra2rev2RPCPort:   "41052",
	ScryptRPCPort:      "41053",
	StribogRPCPort:     "41054",
	SkeinRPCPort:       "41055",
	X11RPCPort:         "41056",
}

// netName returns the name used when referring to a bitcoin network.  At the time of writing, pod currently places blocks for testnet version 3 in the data and log directory "testnet", which does not match the Name field of the chaincfg parameters.  This function can be used to override this directory name as "testnet" when the passed active network matches wire.TestNet3. A proper upgrade to move the data and log directories for this network to "testnet3" is planned for the future, at which point this function can be removed and the network parameter's name used instead.
func netName(chainParams *params) string {
	switch chainParams.Net {
	case wire.TestNet3:
		return "testnet"
	default:
		return chainParams.Name
	}
}
