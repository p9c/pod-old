// Package fork handles tracking the hard fork status and is used to determine which consensus rules apply on a block
// TODO: add trailing auto-checkpoint system
package fork

import (
	"crypto/rand"
	"encoding/hex"
	"math/big"

	"git.parallelcoin.io/pod/chaincfg/chainhash"
)

// HardForks is the details related to a hard fork, number, name and activation height
type HardForks struct {
	Number           uint32
	Name             string
	ActivationHeight int32
	Algos            map[string]AlgoParams
	AlgoVers         map[int32]string
	WorkBase         int64
}

// AlgoParams are the identifying block version number and their minimum target bits
type AlgoParams struct {
	Version int32
	MinBits uint32
	AlgoID  uint32
	NSperOp int64
}

var (
	// IsTestnet is set at startup here to be accessible to all other libraries
	IsTestnet bool
	// List is the list of existing hard forks and when they activate
	List = []HardForks{
		{
			Number:           0,
			Name:             "Halcyon days",
			ActivationHeight: 0, // Approximately 18 Jan 2019
			Algos:            Algos,
			AlgoVers:         AlgoVers,
			WorkBase:         1,
		},
		{
			Number:           1,
			Name:             "Plan 9 from Crypto Space",
			ActivationHeight: 199999,
			Algos:            P9Algos,
			AlgoVers:         P9AlgoVers,
			WorkBase: func() (out int64) {
				for i := range P9Algos {
					out += P9Algos[i].NSperOp
				}
				out /= int64(len(P9Algos))
				return
			}(),
		},
	}
	mainPowLimit = func() big.Int {
		mplb, _ := hex.DecodeString("00000fffff000000000000000000000000000000000000000000000000000000")
		return *big.NewInt(0).SetBytes(mplb)
	}()
	mainPowLimitBits = BigToCompact(&mainPowLimit)

	p9PowLimit = func() big.Int {
		mplb, _ := hex.DecodeString("000ffffff0000000000000000000000000000000000000000000000000000000")
		return *big.NewInt(0).SetBytes(mplb)
	}()
	p9PowLimitBits = BigToCompact(&p9PowLimit)

	SecondPowLimit = func() big.Int {
		mplb, _ := hex.DecodeString("07fffffff0000000000000000000000000000000000000000000000000000000")
		return *big.NewInt(0).SetBytes(mplb)
	}()
	SecondPowLimitBits = BigToCompact(&SecondPowLimit)

	FirstPowLimit = func() big.Int {
		mplb, _ := hex.DecodeString("0fffffff00000000000000000000000000000000000000000000000000000000")
		return *big.NewInt(0).SetBytes(mplb)
	}()
	FirstPowLimitBits = BigToCompact(&FirstPowLimit)

	// Algos are the specifications identifying the algorithm used in the block proof
	Algos = map[string]AlgoParams{
		"sha256d": {2, mainPowLimitBits, 0, 1},   //824 ns/op
		"scrypt":  {514, mainPowLimitBits, 1, 1}, //740839 ns/op
	}
	// P9Algos is the algorithm specifications after the hard fork
	P9Algos = map[string]AlgoParams{
		"blake14lr":      {0, FirstPowLimitBits, 0, 43935943},
		"cryptonight7v2": {1, FirstPowLimitBits, 1, 44195890},
		"keccak":         {2, FirstPowLimitBits, 2, 42804256},
		"lyra2rev2":      {3, FirstPowLimitBits, 3, 76719207},
		"scrypt":         {4, FirstPowLimitBits, 4, 43898224},
		"sha256d":        {5, FirstPowLimitBits, 5, 43418857},
		"skein":          {6, FirstPowLimitBits, 7, 44523156},
		"stribog":        {7, FirstPowLimitBits, 6, 46297969},
		"x11":            {8, FirstPowLimitBits, 8, 44318830},
	}
	// AlgoVers is the lookup for pre hardfork
	AlgoVers = map[int32]string{
		2:   "sha256d",
		514: "scrypt",
	}
	// P9AlgoVers is the lookup for after 1st hardfork
	P9AlgoVers = map[int32]string{
		0: "blake14lr",
		1: "cryptonight7v2",
		2: "keccak",
		3: "lyra2rev2",
		4: "scrypt",
		5: "sha256d",
		6: "skein",
		7: "stribog",
		8: "x11",
	}
)

const (
	Blake14lrReps = 1
	Lyra2rev2Reps = 1
)

// Hash computes the hash of bytes using the named hash
func Hash(bytes []byte, name string, height int32) (out chainhash.Hash) {
	switch name {
	case "blake14lr":
		out.SetBytes(Argon2i(Blake14lr(Cryptonight7v2(bytes))))
	case "cryptonight7v2":
		out.SetBytes(Argon2i(Cryptonight7v2(bytes)))
	case "lyra2rev2":
		out.SetBytes(Argon2i(Lyra2REv2(Cryptonight7v2(bytes))))
	case "scrypt":
		if GetCurrent(height) > 0 {
			bytes = Argon2i(Scrypt(Cryptonight7v2(bytes)))
		} else {
			bytes = Scrypt(bytes)
		}
		out.SetBytes(bytes)
	case "sha256d": // sha256d
		if GetCurrent(height) > 0 {
			bytes = Argon2i(chainhash.DoubleHashB(Cryptonight7v2(bytes)))
		} else {
			bytes = chainhash.DoubleHashB(bytes)
		}
		out.SetBytes(bytes)
	case "stribog":
		out.SetBytes(Argon2i(Stribog(Cryptonight7v2(bytes))))
	case "skein":
		out.SetBytes(Argon2i(Skein(Cryptonight7v2(bytes))))
	case "x11":
		out.SetBytes(Argon2i(X11(Cryptonight7v2(bytes))))
	case "keccak":
		out.SetBytes(Argon2i(Keccak(Cryptonight7v2(bytes))))
	}
	return
}

// GetAlgoVer returns the version number for a given algorithm (by string name) at a given height. If "random" is given, a random number is taken from the system secure random source (for randomised cpu mining)
func GetAlgoVer(name string, height int32) (version int32) {
	n := "sha256d"
	hf := GetCurrent(height)
	if name == "random" {
		rn, _ := rand.Int(rand.Reader, big.NewInt(8))
		randomalgover := int32(rn.Uint64())
		switch hf {
		case 0:
			rndalgo := List[0].AlgoVers[randomalgover&1]
			algo := List[0].Algos[rndalgo].Version
			return algo
		case 1:
			rndalgo := List[1].AlgoVers[randomalgover]
			algo := List[1].Algos[rndalgo].Version
			return algo
		}
	} else {
		n = name
	}
	if IsTestnet {
		return List[len(List)-1].Algos[n].Version
	}
	for i := range List {
		if height > List[i].ActivationHeight {
			version = List[i].Algos[n].Version
		}
	}
	return
}

// GetAlgoName returns the string identifier of an algorithm depending on hard fork activation status
func GetAlgoName(algoVer int32, height int32) (name string) {
	if IsTestnet {
		return List[len(List)-1].AlgoVers[algoVer]
	}
	for i := range List {
		if height > List[i].ActivationHeight {
			name = List[i].AlgoVers[algoVer]
		}
	}
	return
}

// GetAlgoID returns the 'algo_id' which in pre-hardfork is not the same as the block version number, but is afterwards
func GetAlgoID(algoname string, height int32) uint32 {
	if GetCurrent(height) > 1 {
		return P9Algos[algoname].AlgoID
	}
	return Algos[algoname].AlgoID
}

// GetCurrent returns the hardfork number code
func GetCurrent(height int32) (curr int) {
	if IsTestnet {
		return len(List) - 1
	}
	for i := range List {
		if height > List[i].ActivationHeight {
			curr = i
		}
	}
	return
}

// GetMinBits returns the minimum diff bits based on height and testnet
func GetMinBits(algoname string, height int32) uint32 {
	curr := GetCurrent(height)
	return List[curr].Algos[algoname].MinBits
}

// GetMinDiff returns the minimum difficulty in uint256 form
func GetMinDiff(algoname string, height int32) *big.Int {
	return CompactToBig(GetMinBits(algoname, height))
}

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
func CompactToBig(compact uint32) *big.Int {
	// Extract the mantissa, sign bit, and exponent.
	mantissa := compact & 0x007fffff
	isNegative := compact&0x00800000 != 0
	exponent := uint(compact >> 24)
	// Since the base for the exponent is 256, the exponent can be treated as the number of bytes to represent the full 256-bit number.  So, treat the exponent as the number of bytes and shift the mantissa right or left accordingly.  This is equivalent to N = mantissa * 256^(exponent-3)
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
func BigToCompact(n *big.Int) uint32 {
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
