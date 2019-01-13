package addrmgr

import (
	"github.com/parallelcointeam/pod/wire"
	"time"
)

func TstKnownAddressIsBad(ka *KnownAddress) bool {
	return ka.isBad()
}
func TstKnownAddressChance(ka *KnownAddress) float64 {
	return ka.chance()
}
func TstNewKnownAddress(na *wire.NetAddress, attempts int, lastattempt, lastsuccess time.Time, tried bool, refs int) *KnownAddress {
	return &KnownAddress{na: na, attempts: attempts, lastattempt: lastattempt, lastsuccess: lastsuccess, tried: tried, refs: refs}
}
