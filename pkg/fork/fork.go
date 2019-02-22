// Package fork handles tracking the hard fork status and is used to determine which consensus rules apply on a block
// TODO: add trailing auto-checkpoint system
package fork

import (
	"crypto/rand"
	"encoding/hex"
	"math/big"
	"time"
)

// HardForks is the details related to a hard fork, number, name and activation height
type HardForks struct {
	// Number is the sequence number of the hardfork
	Number uint32
	// Name is the human readable name of the hardfork
	Name string
	// ActivationHeight is the height at which the hardfork takes effect
	ActivationHeight int32
	// Algos is the list of algorithms used in the hardfork
	Algos map[string]AlgoParams
	// AlgoVers is a reverse lookup to find the name from the block version
	AlgoVers map[int32]string
	// WorkBase is the average time per hash calculated from the algo's nsperop value
	WorkBase int64
	// TargetTimePerBlock is the amount of time the network targets for between blocks
	TargetTimePerBlock time.Duration
	// TestNetStart is the activation height when in testnet
	TestnetStart int64
	// AveragingInterval is the number of blocks in the 'trailing'  simple average
	AveragingInterval int64
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

	// SecondPowLimit is
	SecondPowLimit = func() big.Int {
		mplb, _ := hex.DecodeString("07fffffff0000000000000000000000000000000000000000000000000000000")
		return *big.NewInt(0).SetBytes(mplb)
	}()
	// SecondPowLimitBits is
	SecondPowLimitBits = BigToCompact(&SecondPowLimit)
	// FirstPowLimit is
	FirstPowLimit = func() big.Int {
		mplb, _ := hex.DecodeString("0fffffff00000000000000000000000000000000000000000000000000000000")
		return *big.NewInt(0).SetBytes(mplb)
	}()
	// FirstPowLimitBits is
	FirstPowLimitBits = BigToCompact(&FirstPowLimit)

	// Algos are the specifications identifying the algorithm used in the block proof
	Algos = map[string]AlgoParams{
		"sha256d": {2, mainPowLimitBits, 0, 1},   //824 ns/op
		"scrypt":  {514, mainPowLimitBits, 1, 1}, //740839 ns/op
	}
	// P9Algos is the algorithm specifications after the hard fork
	P9Algos = map[string]AlgoParams{
		"blake2b":   {0, FirstPowLimitBits, 0, 69495444},
		"blake14lr": {1, FirstPowLimitBits, 1, 79734306},
		"blake2s":   {2, FirstPowLimitBits, 2, 69968425},
		"keccak":    {3, FirstPowLimitBits, 3, 71988313},
		"scrypt":    {4, FirstPowLimitBits, 4, 68395274},
		"sha256d":   {5, FirstPowLimitBits, 5, 67460443},
		"skein":     {6, FirstPowLimitBits, 7, 64433603},
		"stribog":   {7, FirstPowLimitBits, 6, 69987634},
		"x11":       {8, FirstPowLimitBits, 8, 64936544},
	}
	// AlgoVers is the lookup for pre hardfork
	AlgoVers = map[int32]string{
		2:   "sha256d",
		514: "scrypt",
	}
	// P9AlgoVers is the lookup for after 1st hardfork
	P9AlgoVers = map[int32]string{
		0: "blake2b",
		1: "blake14lr",
		2: "blake2s",
		3: "keccak",
		4: "scrypt",
		5: "sha256d",
		6: "skein",
		7: "stribog",
		8: "x11",
	}
	// List is the list of existing hard forks and when they activate
	List = []HardForks{
		{
			Number:             0,
			Name:               "Halcyon days",
			ActivationHeight:   0,
			Algos:              Algos,
			AlgoVers:           AlgoVers,
			WorkBase:           1,
			TargetTimePerBlock: 3 * time.Minute,
			AveragingInterval:  10, // 50 minutes
			TestnetStart:       0,
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
			TargetTimePerBlock: 9 * time.Second,
			AveragingInterval:  9600, // 24 hours
			TestnetStart:       36,
		},
	}
)

// GetAlgoVer returns the version number for a given algorithm (by string name) at a given height. If "random" is given, a random number is taken from the system secure random source (for randomised cpu mining)
func GetAlgoVer(
	name string, height int32) (version int32) {

	n := "sha256d"
	hf := GetCurrent(height)
	if name == "random" {
		rn, _ := rand.Int(rand.Reader, big.NewInt(8))
		randomalgover := int32(rn.Uint64())
		switch hf {
		case 0:
			switch randomalgover & 1 {
			case 0:
				version = 2
			case 1:
				version = 514
			}
			return
		case 1:
			rndalgo := List[1].AlgoVers[randomalgover]
			algo := List[1].Algos[rndalgo].Version
			return algo
		}
	} else {
		n = name
	}
	version = List[hf].Algos[n].Version
	return
}

// GetAlgoName returns the string identifier of an algorithm depending on hard fork activation status
func GetAlgoName(
	algoVer int32, height int32) (name string) {

	hf := GetCurrent(height)
	name = List[hf].AlgoVers[algoVer]
	return
}

// GetAlgoID returns the 'algo_id' which in pre-hardfork is not the same as the block version number, but is afterwards
func GetAlgoID(
	algoname string, height int32) uint32 {
	if GetCurrent(height) > 1 {
		return P9Algos[algoname].AlgoID
	}
	return Algos[algoname].AlgoID
}

// GetCurrent returns the hardfork number code
func GetCurrent(
	height int32) (curr int) {

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
func GetMinBits(
	algoname string, height int32) uint32 {
	curr := GetCurrent(height)
	r := List[curr].Algos[algoname].MinBits
	return r
}

// GetMinDiff returns the minimum difficulty in uint256 form
func GetMinDiff(
	algoname string, height int32) *big.Int {
	return CompactToBig(GetMinBits(algoname, height))
}
