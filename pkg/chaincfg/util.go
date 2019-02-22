package chaincfg

import "math/big"

// CompactToBig converts a compact representation of a whole number N to an unsigned 32-bit number.  The representation is similar to IEEE754 floating point numbers.
// Like IEEE754 floating point, there are three basic components: the sign, the exponent, and the mantissa.  They are broken out as follows:
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
func compactToBig(
	compact uint32) *big.Int {

	// Extract the mantissa, sign bit, and exponent.
	mantissa := compact & 0x007fffff
	isNegative := compact&0x00800000 != 0
	exponent := uint(compact >> 24)

	// Since the base for the exponent is 256, the exponent can be treated as the number of bytes to represent the full 256-bit number.  So, treat the exponent as the number of bytes and shift the mantissa right or left accordingly.  This is equivalent to `N = mantissa * 256^(exponent-3)``
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
