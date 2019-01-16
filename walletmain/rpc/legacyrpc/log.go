package legacyrpc

import l "git.parallelcoin.io/pod/log"

var log = l.Disabled

// UseLogger sets the package-wide logger.  Any calls to this function must be
// made before a server is created and used (it is not concurrent safe).
func UseLogger(logger l.Logger) {
	log = logger
}
