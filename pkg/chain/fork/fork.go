// Package fork handles tracking the hard fork status and is used to determine which consensus rules apply on a block
// TODO: add trailing auto-checkpoint system and hard fork block time change
package fork

import (
	"crypto/rand"
	"encoding/hex"
	"math/big"
	"time"
)

// AlgoParams are the identifying block version number and their minimum target bits

type AlgoParams struct {
	Version int32
	MinBits uint32
	AlgoID  uint32
	NSperOp int64
}

// HardForks is the details related to a hard fork, number, name and activation height

type HardForks struct {
	Number             uint32
	Name               string
	ActivationHeight   int32
	Algos              map[string]AlgoParams
	AlgoVers           map[int32]string
	WorkBase           int64
	TargetTimePerBlock time.Duration
	TestnetStart       int32
	AveragingInterval  int64
}

// AlgoVers is the lookup for pre hardfork
var AlgoVers = map[int32]string{
	2:   "sha256d",
	514: "scrypt",
}

// Algos are the specifications identifying the algorithm used in the block proof
var Algos = map[string]AlgoParams{
	AlgoVers[2]:   {2, mainPowLimitBits, 0, 824},      //824 ns/op
	AlgoVers[514]: {514, mainPowLimitBits, 1, 740839}, //740839 ns/op
}

// FirstPowLimit is
var FirstPowLimit = func() big.Int {

	mplb, _ := hex.DecodeString("0fffffff00000000000000000000000000000000000000000000000000000000")
	return *big.NewInt(0).SetBytes(mplb)
}()

// FirstPowLimitBits is
var FirstPowLimitBits = BigToCompact(&FirstPowLimit)

// IsTestnet is set at startup here to be accessible to all other libraries
var IsTestnet bool

// List is the list of existing hard forks and when they activate
var List = []HardForks{
	{
		Number:           0,
		Name:             "Halcyon days",
		ActivationHeight: 0,
		Algos:            Algos,
		AlgoVers:         AlgoVers,
		WorkBase: func() (out int64) {

			for i := range Algos {

				out += Algos[i].NSperOp
			}
			out /= int64(len(Algos))
			return
		}(),
		TargetTimePerBlock: 3 * time.Minute,
		AveragingInterval:  10, // 50 minutes
		TestnetStart:       0,
	},
	{
		Number:           1,
		Name:             "Plan 9 from Crypto Space",
		ActivationHeight: 250000,
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
		TestnetStart:       100,
	},
}

// P9AlgoVers is the lookup for after 1st hardfork
var P9AlgoVers = map[int32]string{
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

// P9Algos is the algorithm specifications after the hard fork
var P9Algos = map[string]AlgoParams{
	P9AlgoVers[0]: {0, FirstPowLimitBits, 0, 69495444},
	P9AlgoVers[1]: {1, FirstPowLimitBits, 1, 79734306},
	P9AlgoVers[2]: {2, FirstPowLimitBits, 2, 69968425},
	P9AlgoVers[3]: {3, FirstPowLimitBits, 3, 71988313},
	P9AlgoVers[4]: {4, FirstPowLimitBits, 4, 68395274},
	P9AlgoVers[5]: {5, FirstPowLimitBits, 5, 67460443},
	P9AlgoVers[6]: {6, FirstPowLimitBits, 7, 64433603},
	P9AlgoVers[7]: {7, FirstPowLimitBits, 6, 69987634},
	P9AlgoVers[8]: {8, FirstPowLimitBits, 8, 64936544},
}

// SecondPowLimit is
var SecondPowLimit = func() big.Int {

	mplb, _ := hex.DecodeString("07fffffff0000000000000000000000000000000000000000000000000000000")
	return *big.NewInt(0).SetBytes(mplb)
}()

// SecondPowLimitBits is
var SecondPowLimitBits = BigToCompact(&SecondPowLimit)

var mainPowLimit = func() big.Int {

	mplb, _ := hex.DecodeString("00000fffff000000000000000000000000000000000000000000000000000000")
	return *big.NewInt(0).SetBytes(mplb)
}()

var mainPowLimitBits = BigToCompact(&mainPowLimit)

// GetAlgoID returns the 'algo_id' which in pre-hardfork is not the same as the block version number, but is afterwards
func GetAlgoID(
	algoname string,
	height int32,
) uint32 {

	if GetCurrent(height) > 1 {

		return P9Algos[algoname].AlgoID
	}
	return Algos[algoname].AlgoID
}

// GetAlgoName returns the string identifier of an algorithm depending on hard fork activation status
func GetAlgoName(
	algoVer int32,
	height int32,
) (
	name string,
) {

	hf := GetCurrent(height)
	name = List[hf].AlgoVers[algoVer]
	return
}

// GetAlgoVer returns the version number for a given algorithm (by string name) at a given height. If "random" is given, a random number is taken from the system secure random source (for randomised cpu mining)
func GetAlgoVer(
	name string,
	height int32,
) (
	version int32,
) {

	n := "sha256d"
	hf := GetCurrent(height)

	if name == "random" {

		rn, _ := rand.Int(rand.Reader,
			big.NewInt(int64(len(P9AlgoVers)-1)))
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

// GetAveragingInterval returns the active block interval target based on hard fork status
func GetAveragingInterval(
	height int32,
) (
	r int64,
) {

	r = int64(List[GetCurrent(height)].AveragingInterval)
	return
}

// GetCurrent returns the hardfork number code
func GetCurrent(
	height int32,
) (
	curr int,
) {

	if IsTestnet {

		for i := range List {

			if height > List[i].TestnetStart {

				curr = i
			}
		}
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
	algoname string,
	height int32,
) (
	mb uint32,
) {

	curr := GetCurrent(height)
	mb = List[curr].Algos[algoname].MinBits
	return
}

// GetMinDiff returns the minimum difficulty in uint256 form
func GetMinDiff(
	algoname string,
	height int32,
) (
	md *big.Int,
) {

	return CompactToBig(GetMinBits(algoname, height))
}

// GetTargetTimePerBlock returns the active block interval target based on hard fork status
func GetTargetTimePerBlock(
	height int32,
) (
	r int64,
) {

	r = int64(List[GetCurrent(height)].TargetTimePerBlock)
	return
}
