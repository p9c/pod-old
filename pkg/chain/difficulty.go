package blockchain

import (
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"time"

	"git.parallelcoin.io/dev/pod/pkg/chain/fork"
	chainhash "git.parallelcoin.io/dev/pod/pkg/chain/hash"
	cl "git.parallelcoin.io/dev/pod/pkg/util/cl"
)

var ScryptPowLimit = scryptPowLimit

var ScryptPowLimitBits = BigToCompact(&scryptPowLimit)

// bigOne is 1 represented as a big.Int.  It is defined here to avoid the overhead of creating it multiple times.
var bigOne = big.NewInt(1)

// oneLsh256 is 1 shifted left 256 bits.  It is defined here to avoid the overhead of creating it multiple times.
var oneLsh256 = new(big.Int).Lsh(bigOne, 256)

var scryptPowLimit = func() big.Int {

	mplb, _ := hex.DecodeString("000000039fcaa04ac30b6384471f337748ef5c87c7aeffce5e51770ce6283137,")
	return *big.NewInt(0).SetBytes(mplb) //AllOnes.Rsh(&AllOnes, 0)
}()

// CalcNextRequiredDifficulty calculates the required difficulty for the block after the end of the current best chain based on the difficulty retarget rules. This function is safe for concurrent access.
func (
	b *BlockChain,
) CalcNextRequiredDifficulty(
	timestamp time.Time,
	algo string,
) (
	difficulty uint32,
	err error,

) {

	b.chainLock.Lock()
	difficulty, err = b.calcNextRequiredDifficulty(b.bestChain.Tip(), timestamp, algo, true)
	b.chainLock.Unlock()
	return
}

// calcEasiestDifficulty calculates the easiest possible difficulty that a block can have given starting difficulty bits and a duration.  It is mainly used to verify that claimed proof of work by a block is sane as compared to a known good checkpoint.
func (
	b *BlockChain,
) calcEasiestDifficulty(
	bits uint32,
	duration time.Duration,

) uint32 {

	// Convert types used in the calculations below.
	durationVal := int64(duration / time.Second)
	adjustmentFactor := big.NewInt(b.chainParams.RetargetAdjustmentFactor)

	// Since easier difficulty equates to higher numbers, the easiest difficulty for a given duration is the largest value possible given the number of retargets for the duration and starting difficulty multiplied by the max adjustment factor.
	newTarget := CompactToBig(bits)

	for durationVal > 0 && newTarget.Cmp(b.chainParams.PowLimit) < 0 {

		newTarget.Mul(newTarget, adjustmentFactor)
		durationVal -= b.maxRetargetTimespan
	}

	// Limit new value to the proof of work limit.

	if newTarget.Cmp(b.chainParams.PowLimit) > 0 {

		newTarget.Set(b.chainParams.PowLimit)
	}

	return BigToCompact(newTarget)
}

// calcNextRequiredDifficulty calculates the required difficulty for the block after the passed previous block node based on the difficulty retarget rules. This function differs from the exported  CalcNextRequiredDifficulty in that the exported version uses the current best chain as the previous block node while this function accepts any block node.
func (
	b *BlockChain,
) calcNextRequiredDifficulty(
	lastNode *blockNode,
	newBlockTime time.Time,
	algoname string,
	l bool,
) (
	newTargetBits uint32,
	err error,

) {

	nH := lastNode.height + 1

	switch fork.GetCurrent(nH) {

	// Legacy difficulty adjustment
	case 0:

		log <- cl.Debug{"on pre-hardfork"}

		algo := fork.GetAlgoVer(algoname, nH)
		algoName := fork.GetAlgoName(algo, nH)
		newTargetBits = fork.GetMinBits(algoName, nH)

		// log <- cl.Debugf{
		// 	"newTargetBits %064x", CompactToBig(newTargetBits),
		// }

		if lastNode == nil {

			return newTargetBits, nil
		}

		prevNode := lastNode.RelativeAncestor(1)
		log <- cl.Debug{"prevNode version", prevNode.version, prevNode.height}
		prevversion := prevNode.version

		if fork.GetCurrent(prevNode.height) == 0 {

			if prevversion != 514 &&

				prevversion != 2 {

				log <- cl.Warn{"irregular block version, assuming 2 (sha256d)"}
				prevversion = 2
			}

		}

		if prevversion != algo {

			prevNode = prevNode.GetLastWithAlgo(algo)
		}

		if prevNode == nil {

			return newTargetBits, nil
		}

		prevversion = prevNode.version
		log <- cl.Debugf{
			"found version %d corrected %d height %d bits %8x",
			prevNode.version, prevversion, prevNode.height, prevNode.bits}
		firstNode := prevNode // .GetPrevWithAlgo(algo)
		i := int64(1)

		for ; firstNode != nil && i < fork.GetAveragingInterval(nH); i++ {

			firstNode = firstNode.RelativeAncestor(1).GetLastWithAlgo(algo)

			if firstNode == nil {

				log <- cl.Debug{"passed genesis block"}
				return newTargetBits, nil
			}

			log <- cl.Debugf{"found a prev %d %d %8x", firstNode.version, firstNode.height, firstNode.bits}
		}

		if firstNode == nil {

			return newTargetBits, nil
		}

		actualTimespan := prevNode.timestamp - firstNode.timestamp
		adjustedTimespan := actualTimespan

		if actualTimespan < b.chainParams.MinActualTimespan {

			adjustedTimespan = b.chainParams.MinActualTimespan

		} else if actualTimespan > b.chainParams.MaxActualTimespan {

			adjustedTimespan = b.chainParams.MaxActualTimespan
		}

		oldTarget := CompactToBig(prevNode.bits)
		newTarget := new(big.Int).
			Mul(oldTarget, big.NewInt(adjustedTimespan))
		newTarget = newTarget.
			Div(newTarget, big.NewInt(b.chainParams.AveragingTargetTimespan))

		if newTarget.Cmp(CompactToBig(newTargetBits)) > 0 {

			newTarget.Set(CompactToBig(newTargetBits))
		}

		newTargetBits = BigToCompact(newTarget)
		log <- cl.Debugf{
			"difficulty retarget at block height %d, old %08x new %08x",
			lastNode.height + 1,
			prevNode.bits,
			newTargetBits,
		}

		Log.Trcc(func() string {

			return fmt.Sprintf(
				"actual timespan %v, adjusted timespan %v, target timespan %v"+
					"\nOld %064x\nNew %064x",
				actualTimespan,
				adjustedTimespan,
				b.chainParams.AveragingTargetTimespan,
				oldTarget,
				CompactToBig(newTargetBits),
			)
		})

		return newTargetBits, nil

	case 1: // Plan 9 from Crypto Space

		log <- cl.Debug{"on plan 9 hardfork"}

		if lastNode.height == 0 {

			return fork.FirstPowLimitBits, nil
		}

		nH := lastNode.height + 1
		algo := fork.GetAlgoVer(algoname, nH)
		newTargetBits = fork.GetMinBits(algoname, nH)
		last := lastNode

		// find the most recent block of the same algo

		if last.version != algo {

			l := last.RelativeAncestor(1)
			l = l.GetLastWithAlgo(algo)

			// ignore the first block as its time is not a normal timestamp

			if l.height < 1 {

				break
			}

			last = l
		}

		counter := 1
		var timestamps []float64
		timestamps = append(timestamps, float64(last.timestamp))
		pb := last

		// collect the timestamps of all the blocks of the same algo until we pass genesis block or get AveragingInterval blocks

		for ; counter < int(fork.GetAveragingInterval(nH)) && pb.height > 2; counter++ {

			p := pb.RelativeAncestor(1)

			if p != nil {

				if p.height == 0 {

					return fork.SecondPowLimitBits, nil
				}

				pb = p.GetLastWithAlgo(algo)

			} else {

				break
			}

			if pb != nil && pb.height > 0 {

				// only add the timestamp if is not the same as the previous
				timestamps = append(timestamps, float64(pb.timestamp))

			} else {

				break
			}

		}

		allTimeAverage, trailTimeAverage := float64(fork.GetTargetTimePerBlock(nH)), float64(fork.GetTargetTimePerBlock(nH))
		startHeight := fork.List[1].ActivationHeight

		if b.chainParams.Name == "testnet" {

			startHeight = 1
		}

		trailHeight := int32(int64(lastNode.height) -
			fork.GetAveragingInterval(nH)*int64(len(fork.List[1].Algos)))

		if trailHeight < 0 {

			trailHeight = 1
		}

		firstBlock, _ := b.BlockByHeight(startHeight)
		trailBlock, _ := b.BlockByHeight(trailHeight)
		lastTime := lastNode.timestamp

		if firstBlock != nil {

			firstTime := firstBlock.MsgBlock().Header.Timestamp.Unix()
			allTimeAverage = (float64(lastTime) - float64(firstTime)) / (float64(lastNode.height) - float64(firstBlock.Height()))
		}

		if trailBlock != nil {

			trailTime := trailBlock.MsgBlock().Header.Timestamp.Unix()
			trailTimeAverage = (float64(lastTime) - float64(trailTime)) / (float64(lastNode.height) - float64(trailBlock.Height()))
		}

		if len(timestamps) < 2 {

			return fork.SecondPowLimitBits, nil
		}

		var adjusted, targetAdjusted, adjustment float64

		if len(timestamps) > 1 {

			numalgos := int64(len(fork.List[1].Algos))
			target := fork.GetTargetTimePerBlock(nH) * numalgos
			counter = 0

			for i := 0; i < len(timestamps)-1; i++ {

				factor := 0.75

				if i == 0 {

					f := factor

					for j := 0; j < i; j++ {

						f *= factor
					}

					factor = f

				} else {

					factor = 1.0
				}

				adjustment = timestamps[i] - timestamps[i+1]
				adjustment *= factor

				switch {

				case math.IsNaN(adjustment):
					break
				case adjustment == 0.0:
					break
				}

				adjusted += adjustment
				targetAdjusted += float64(target) * factor
				counter++
			}

		} else {

			targetAdjusted = 100
			adjusted = 100
		}

		var trailingTimestamps []float64

		pb = lastNode
		trailingTimestamps = append(
			trailingTimestamps, float64(pb.timestamp))
		counter = 1
		for ; counter < int(fork.GetAveragingInterval(nH)) &&

			pb.height > 2; counter++ {

			pb = pb.RelativeAncestor(1)
			trailingTimestamps = append(
				trailingTimestamps, float64(pb.timestamp))
			counter++
		}

		var trailingAdjusted,
			trailingTargetAdjusted,
			trailingAdjustment float64

		if len(trailingTimestamps) > 1 {

			target := fork.GetTargetTimePerBlock(nH)
			counter = 0

			for i := 0; i < len(trailingTimestamps)-1; i++ {

				factor := 0.81

				if i == 0 {

					f := factor

					for j := 0; j < i; j++ {

						f *= factor
					}

					factor = f

				} else {

					factor = 1.0
				}

				trailingAdjustment = trailingTimestamps[i] - trailingTimestamps[i+1]
				trailingAdjustment *= factor

				switch {

				case math.IsNaN(trailingAdjustment):
					break
				case trailingAdjustment == 0.0:
					break
				}

				trailingAdjusted += trailingAdjustment
				trailingTargetAdjusted += float64(target) * factor
				counter++
			}

		} else {

			trailingTargetAdjusted = 100
			trailingAdjusted = 100
		}

		ttpb := float64(fork.GetTargetTimePerBlock(nH))
		allTimeDivergence := allTimeAverage / ttpb
		trailTimeDivergence := trailTimeAverage / ttpb
		trailingTimeDivergence := trailingAdjusted / trailingTargetAdjusted
		log <- cl.Trace{
			"trailingtimedivergence",
			trailingTimeDivergence,
			trailingAdjusted,
			trailingTargetAdjusted}
		weighted := adjusted / targetAdjusted
		adjustment = (weighted*weighted*weighted +
			trailingTimeDivergence*trailingTimeDivergence*trailingTimeDivergence +
			trailTimeDivergence*trailTimeDivergence*trailTimeDivergence +
			allTimeDivergence*allTimeDivergence*allTimeDivergence) / 4.0

		if adjustment < 0 {

			fmt.Println("negative weight adjustment")
			adjustment = allTimeDivergence
		}

		if math.IsNaN(adjustment) {

			return lastNode.bits, nil
		}

		// Bias adjustment for difficulty reductions to reduce incidence of sub 1 second blocks

		if adjustment < 0 {

			adjustment = (1 - adjustment) * adjustment
		}

		bigadjustment := big.NewFloat(adjustment)
		bigoldtarget := big.NewFloat(1.0).SetInt(CompactToBig(last.bits))
		bigfnewtarget := big.NewFloat(1.0).Mul(bigadjustment, bigoldtarget)
		newtarget, _ := bigfnewtarget.Int(nil)

		if newtarget == nil {

			return newTargetBits, nil
		}

		mintarget := CompactToBig(newTargetBits)

		if newtarget.Cmp(mintarget) < 0 {

			newTargetBits = BigToCompact(newtarget)
			b.DifficultyAdjustments[algoname] = adjustment

			if l {

				log <- cl.Infof{
					"%d: old %08x, new %08x, av %3.2f, tr %3.2f, tr wgtd %3.2f, alg wgtd %3.2f, blks %d, adj %0.1f%%, alg %s",
					lastNode.height + 1, last.bits,
					newTargetBits,
					allTimeAverage,
					trailTimeAverage,
					trailingTimeDivergence * ttpb,
					weighted * ttpb,
					counter,
					(1 - adjustment) * 100,
					fork.List[1].AlgoVers[algo],
				}

			}

		}

		return newTargetBits, nil
	}

	// nH := lastNode.height + 1

	// algo := fork.GetAlgoVer(algoname, nH)
	return fork.GetMinBits(algoname, nH), nil
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

// CalcWork calculates a work value from difficulty bits.  Bitcoin increases the difficulty for generating a block by decreasing the value which the generated hash must be less than.  This difficulty target is stored in each block header using a compact representation as described in the documentation for CompactToBig. The main chain is selected by choosing the chain that has the most proof of work (highest difficulty). Since a lower target difficulty value equates to higher actual difficulty, the work value which will be accumulated must be the inverse of the difficulty.  Also, in order to avoid potential division by zero and really small floating point numbers, the result adds 1 to the denominator and multiplies the numerator by 2^256.

func CalcWork(bits uint32, height int32, algover int32) *big.Int {

	// Return a work value of zero if the passed difficulty bits represent a negative number. Note this should not happen in practice with valid blocks, but an invalid block could trigger it.
	difficultyNum := CompactToBig(bits)

	// To make the difficulty values correlate to number of hash operations, multiply this difficulty base by the nanoseconds/hash figures in the fork algorithms list
	current := fork.GetCurrent(height)
	algoname := fork.List[current].AlgoVers[algover]
	difficultyNum = new(big.Int).Mul(difficultyNum, big.NewInt(fork.List[current].Algos[algoname].NSperOp))
	difficultyNum = new(big.Int).Quo(difficultyNum, big.NewInt(fork.List[current].WorkBase))

	if difficultyNum.Sign() <= 0 {

		return big.NewInt(0)
	}

	denominator := new(big.Int).Add(difficultyNum, bigOne)
	r := new(big.Int).Div(oneLsh256, denominator)
	return r
}

// CompactToBig converts a compact representation of a whole number N to an unsigned 32-bit number.  The representation is similar to IEEE754 floating point numbers.
/*
Like IEEE754 floating point, there are three basic components: the sign, the exponent, and the mantissa.  They are broken out as follows:

	* the most significant 8 bits represent the unsigned base 256 exponent
	* bit 23 (the 24th bit) represents the sign bit
	* the least significant 23 bits represent the mantissa

	-------------------------------------------------
	|   Exponent     |    Sign    |    Mantissa     |
	-------------------------------------------------
	| 8 bits [31-24] | 1 bit [23] | 23 bits [22-00] |
	-------------------------------------------------

The formula to calculate N is:

	N = (-1^sign) * mantissa * 256^(exponent-3)

This compact form is only used in bitcoin to encode unsigned 256-bit numbers which represent difficulty targets, thus there really is not a need for a sign bit, but it is implemented here to stay consistent with bitcoind.
*/

func CompactToBig(compact uint32) *big.Int {

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

// HashToBig converts a chainhash.Hash into a big.Int that can be used to perform math comparisons.

func HashToBig(hash *chainhash.Hash) *big.Int {

	// A Hash is in little-endian, but the big package wants the bytes in big-endian, so reverse them.
	buf := *hash
	blen := len(buf)

	for i := 0; i < blen/2; i++ {

		buf[i], buf[blen-1-i] = buf[blen-1-i], buf[i]
	}

	return new(big.Int).SetBytes(buf[:])
}
