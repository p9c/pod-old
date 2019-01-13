package legacyrpc

import "github.com/parallelcointeam/pod/btclog"

var log = btclog.Disabled

// UseLogger sets the package-wide logger.  Any calls to this function must be
// made before a server is created and used (it is not concurrent safe).
func UseLogger(logger btclog.Logger) {
	log = logger
}
