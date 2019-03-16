package chaincfg

import (
	"math/big"
	"strings"

	"git.parallelcoin.io/dev/pod/pkg/chain/hash"
)

// String returns the hostname of the DNS seed in human-readable form.
func (d DNSSeed) String() string {
	return d.Host
}

// Register registers the network parameters for a Bitcoin network.  This may error with ErrDuplicateNet if the network is already registered (either due to a previous Register call, or the network being one of the default networks). Network parameters should be registered into this package by a main package as early as possible.  Then, library packages may lookup networks or network parameters based on inputs and work regardless of the network being standard or not.
func Register(
	params *Params) error {
	if _, ok := registeredNets[params.Net]; ok {
		return ErrDuplicateNet
	}
	registeredNets[params.Net] = struct{}{}
	pubKeyHashAddrIDs[params.PubKeyHashAddrID] = struct{}{}
	scriptHashAddrIDs[params.ScriptHashAddrID] = struct{}{}
	hdPrivToPubKeyIDs[params.HDPrivateKeyID] = params.HDPublicKeyID[:]

	// A valid Bech32 encoded segwit address always has as prefix the human-readable part for the given net followed by '1'.
	bech32SegwitPrefixes[params.Bech32HRPSegwit+"1"] = struct{}{}
	return nil
}

// mustRegister performs the same function as Register except it panics if there is an error.  This should only be called from package init functions.
func mustRegister(
	params *Params) {

	if err := Register(params); err != nil {
		panic("failed to register network: " + err.Error())
	}
}

// IsPubKeyHashAddrID returns whether the id is an identifier known to prefix a pay-to-pubkey-hash address on any default or registered network.  This is used when decoding an address string into a specific address type.  It is up to the caller to check both this and IsScriptHashAddrID and decide whether an address is a pubkey hash address, script hash address, neither, or undeterminable (if both return true).
func IsPubKeyHashAddrID(
	id byte) bool {
	_, ok := pubKeyHashAddrIDs[id]
	return ok
}

// IsScriptHashAddrID returns whether the id is an identifier known to prefix a pay-to-script-hash address on any default or registered network.  This is used when decoding an address string into a specific address type.  It is up to the caller to check both this and IsPubKeyHashAddrID and decide whether an address is a pubkey hash address, script hash address, neither, or undeterminable (if both return true).
func IsScriptHashAddrID(
	id byte) bool {
	_, ok := scriptHashAddrIDs[id]
	return ok
}

// IsBech32SegwitPrefix returns whether the prefix is a known prefix for segwit addresses on any default or registered network.  This is used when decoding an address string into a specific address type.
func IsBech32SegwitPrefix(
	prefix string) bool {
	prefix = strings.ToLower(prefix)
	_, ok := bech32SegwitPrefixes[prefix]
	return ok
}

// HDPrivateKeyToPublicKeyID accepts a private hierarchical deterministic extended key id and returns the associated public key id.  When the provided id is not registered, the ErrUnknownHDKeyID error will be returned.
func HDPrivateKeyToPublicKeyID(
	id []byte) ([]byte, error) {

	if len(id) != 4 {
		return nil, ErrUnknownHDKeyID
	}
	var key [4]byte
	copy(key[:], id)
	pubBytes, ok := hdPrivToPubKeyIDs[key]
	if !ok {
		return nil, ErrUnknownHDKeyID
	}
	return pubBytes, nil
}

// newHashFromStr converts the passed big-endian hex string into a chainhash.Hash.  It only differs from the one available in chainhash in that it panics on an error since it will only (and must only) be called with hard-coded, and therefore known good, hashes.
func newHashFromStr(
	hexStr string) *chainhash.Hash {
	hash, err := chainhash.NewHashFromStr(hexStr)
	if err != nil {
		// Ordinarily I don't like panics in library code since it can take applications down without them having a chance to recover which is extremely annoying, however an exception is being made in this case because the only way this can panic is if there is an error in the hard-coded hashes.  Thus it will only ever potentially panic on init and therefore is 100% predictable.
		// loki: Panics are good when the condition should not happen!
		panic(err)
	}
	return hash
}
func init() {


	// Register all default networks when the package is initialized.
	mustRegister(&MainNetParams)
	mustRegister(&TestNet3Params)
	mustRegister(&RegressionNetParams)
	mustRegister(&SimNetParams)
}

// CompactToBig converts a compact representation of a whole number N to an unsigned 32-bit number.  The representation is similar to IEEE754 floating point numbers. Like IEEE754 floating point, there are three basic components: the sign, the exponent, and the mantissa.  They are broken out as follows:
//	* the most significant 8 bits represent the unsigned base 256 exponent
// 	* bit 23 (the 24th bit) represents the sign bit
//	* the least significant 23 bits represent the mantissa
//	-------------------------------------------------
//	|   Exponent     |    Sign    |    Mantissa     |
//	-------------------------------------------------
//	| 8 bits [31-24] | 1 bit [23] | 23 bits [22-00] |
//	-------------------------------------------------
// The formula to calculate N is:
// 	N = (-1^sign) * mantissa * 256^(exponent-3)
// This compact form is only used in bitcoin to encode unsigned 256-bit numbers which represent difficulty targets, thus there really is not a need for a sign bit, but it is implemented here to stay consistent with bitcoind.
func CompactToBig(
	compact uint32) *big.Int {

	// Extract the mantissa, sign bit, and exponent.
	mantissa := compact & 0x007fffff
	isNegative := compact&0x00800000 != 0
	exponent := uint(compact >> 24)

	// Since the base for the exponent is 256, the exponent can be treated as the number of bytes to represent the full 256-bit number.  So, treat the exponent as the number of bytes and shift the mantissa right or left accordingly.  This is equivalent to: N = mantissa * 256^(exponent-3)
	var bn *big.Int
	if exponent <= 3 {
		mantissa >>= 8 * (3 - exponent)
		bn = big.NewInt(int64(mantissa))
	} else {
		bn = big.NewInt(int64(mantissa))
		bn.Lsh(bn, 8*(exponent-3))
	}

	// Make it negative if the sign bit is set.
	if isNegative {
		bn = bn.Neg(bn)
	}
	return bn
}

// BigToCompact converts a whole number N to a compact representation using an unsigned 32-bit number.  The compact representation only provides 23 bits of precision, so values larger than (2^23 - 1) only encode the most significant digits of the number.  See CompactToBig for details.
func BigToCompact(
	n *big.Int) uint32 {

	// No need to do any work if it's zero.
	if n.Sign() == 0 {
		return 0
	}

	// Since the base for the exponent is 256, the exponent can be treated as the number of bytes.  So, shift the number right or left accordingly.  This is equivalent to: mantissa = mantissa / 256^(exponent-3)
	var mantissa uint32
	exponent := uint(len(n.Bytes()))
	if exponent <= 3 {
		mantissa = uint32(n.Bits()[0])
		mantissa <<= 8 * (3 - exponent)
	} else {
		// Use a copy to avoid modifying the caller's original number.
		tn := new(big.Int).Set(n)
		mantissa = uint32(tn.Rsh(tn, 8*(exponent-3)).Bits()[0])
	}

	// When the mantissa already has the sign bit set, the number is too large to fit into the available 23-bits, so divide the number by 256 and increment the exponent accordingly.
	if mantissa&0x00800000 != 0 {
		mantissa >>= 8
		exponent++
	}

	// Pack the exponent, sign bit, and mantissa into an unsigned 32-bit int and return it.
	compact := uint32(exponent<<24) | mantissa
	if n.Sign() < 0 {
		compact |= 0x00800000
	}
	return compact
}
