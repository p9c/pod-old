package blockchain

import (
	"encoding/hex"
	"fmt"
	"github.com/parallelcointeam/pod/chaincfg/chainhash"
	"github.com/parallelcointeam/pod/fork"
	"math"
	"math/big"
	"math/rand"
	"time"
)

var (
	scryptPowLimit = func() big.Int {
		mplb, _ := hex.DecodeString("000000039fcaa04ac30b6384471f337748ef5c87c7aeffce5e51770ce6283137,")
		return *big.NewInt(0).SetBytes(mplb) //AllOnes.Rsh(&AllOnes, 0)
	}()
	ScryptPowLimit     = scryptPowLimit
	ScryptPowLimitBits = BigToCompact(&scryptPowLimit)
	// bigOne is 1 represented as a big.Int.  It is defined here to avoid the overhead of creating it multiple times.
	bigOne = big.NewInt(1)
	// oneLsh256 is 1 shifted left 256 bits.  It is defined here to avoid the overhead of creating it multiple times.
	oneLsh256 = new(big.Int).Lsh(bigOne, 256)
)

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

// CompactToBig converts a compact representation of a whole number N to an unsigned 32-bit number.  The representation is similar to IEEE754 floating
// point numbers.
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

// calcEasiestDifficulty calculates the easiest possible difficulty that a block can have given starting difficulty bits and a duration.  It is mainly used to verify that claimed proof of work by a block is sane as compared to a known good checkpoint.
func (b *BlockChain) calcEasiestDifficulty(bits uint32, duration time.Duration) uint32 {
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

// findPrevTestNetDifficulty returns the difficulty of the previous block which did not have the special testnet minimum difficulty rule applied. This function MUST be called with the chain state lock held (for writes).
func (b *BlockChain) findPrevTestNetDifficulty(startNode *blockNode) uint32 {
	// Search backwards through the chain for the last block without the special rule applied.
	iterNode := startNode
	for iterNode != nil && iterNode.height%b.blocksPerRetarget != 0 &&
		iterNode.bits == b.chainParams.PowLimitBits {
		iterNode = iterNode.parent
	}
	// Return the found difficulty or the minimum difficulty if no appropriate block was found.
	lastBits := b.chainParams.PowLimitBits
	if iterNode != nil {
		lastBits = iterNode.bits
	}
	return lastBits
}

// calcNextRequiredDifficulty calculates the required difficulty for the block after the passed previous block node based on the difficulty retarget rules. This function differs from the exported  CalcNextRequiredDifficulty in that the exported version uses the current best chain as the previous block node while this function accepts any block node.
func (b *BlockChain) calcNextRequiredDifficulty(lastNode *blockNode, newBlockTime time.Time, algoname string, l bool) (newTargetBits uint32, err error) {
	switch fork.GetCurrent(lastNode.height + 1) {
	case 0:
		nH := lastNode.height + 1
		algo := fork.GetAlgoVer(algoname, nH)
		newTargetBits = fork.GetMinBits(algoname, nH)
		if lastNode == nil {
			return newTargetBits, nil
		}
		prevNode := lastNode
		if prevNode.version != algo {
			prevNode = prevNode.GetPrevWithAlgo(algo)
		}
		firstNode := prevNode.GetPrevWithAlgo(algo)
		for i := int64(1); firstNode != nil && i < b.chainParams.AveragingInterval; i++ {
			firstNode = firstNode.RelativeAncestor(1).GetPrevWithAlgo(algo)
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
		newTarget := new(big.Int).Mul(oldTarget, big.NewInt(adjustedTimespan))
		newTarget = newTarget.Div(newTarget, big.NewInt(b.chainParams.AveragingTargetTimespan))
		if newTarget.Cmp(CompactToBig(newTargetBits)) > 0 {
			newTarget.Set(CompactToBig(newTargetBits))
		}
		newTargetBits = BigToCompact(newTarget)
		log.Debugf("Difficulty retarget at block height %d, old %08x new %08x", lastNode.height+1, prevNode.bits, newTargetBits)
		log.Tracef("Old %08x New %08x", prevNode.bits, oldTarget, newTargetBits, CompactToBig(newTargetBits))
		log.Tracef("Actual timespan %v, adjusted timespan %v, target timespan %v",
			actualTimespan,
			adjustedTimespan,
			b.chainParams.AveragingTargetTimespan)
		return newTargetBits, nil

	case 1: // Plan 9 from Crypto Space
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
			l = l.GetPrevWithAlgo(algo)
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
		for ; counter < int(b.chainParams.AveragingInterval) && pb.height > 1; counter++ {
			p := pb.RelativeAncestor(1)
			if p != nil {
				if p.height == 0 {
					return fork.SecondPowLimitBits, nil
				}
				pb = p.GetPrevWithAlgo(algo)
			} else {
				break
			}
			if pb != nil && pb.height > 0 {
				// only add the timestamp if is not the same as the previous
				if float64(pb.timestamp) != timestamps[len(timestamps)-1] {
					timestamps = append(timestamps, float64(pb.timestamp))
				}
			} else {
				break
			}
		}
		allTimeAverage, trailTimeAverage := float64(b.chainParams.TargetTimePerBlock), float64(b.chainParams.TargetTimePerBlock)
		startHeight := fork.List[1].ActivationHeight
		if b.chainParams.Name == "testnet" {
			startHeight = 1
		}
		trailHeight := int32(int64(lastNode.height) - b.chainParams.AveragingInterval*int64(len(fork.List[1].Algos)))
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
			target := b.chainParams.TargetTimePerBlock * numalgos
			adjustment = 1.0
			counter = 0
			for i := 0; i < len(timestamps)-1; i++ {
				factor := 0.9
				if i == 0 {
					f := factor
					for j := 0; j < i; j++ {
						f = f * factor
					}
					factor = f
				} else {
					factor = 1.0
				}
				adjustment = timestamps[i] - timestamps[i+1]
				adjustment = adjustment * factor
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
		ttpb := float64(b.chainParams.TargetTimePerBlock)
		allTimeDivergence := allTimeAverage / ttpb
		trailTimeDivergence := trailTimeAverage / ttpb
		weighted := adjusted / targetAdjusted
		adjustment = (weighted*weighted*weighted + allTimeDivergence + trailTimeDivergence) / 3.0
		if adjustment < 0 {
			fmt.Println("negative weight adjustment")
			adjustment = allTimeDivergence
		}
		// d := adjustment - 1.0
		// adjustment = 1.0 + (d*d*d+d+d*d)
		if math.IsNaN(adjustment) {
			return lastNode.bits, nil
		}
		bigadjustment := big.NewFloat(adjustment)
		bigoldtarget := big.NewFloat(1.0).SetInt(CompactToBig(last.bits))
		bigfnewtarget := big.NewFloat(1.0).Mul(bigadjustment, bigoldtarget)
		newtarget, _ := bigfnewtarget.Int(nil)
		if newtarget == nil {
			return newTargetBits, nil
		}
		mintarget := CompactToBig(newTargetBits)
		var delay uint16
		if newtarget.Cmp(mintarget) < 0 {
			newTargetBits = BigToCompact(newtarget)
			if b.chainParams.Name == "testnet" {
				rand.Seed(time.Now().UnixNano())
				delay = uint16(rand.Int()) >> 6
				// fmt.Printf("%s testnet delay %dms algo %s\n", time.Now().Format("2006-01-02 15:04:05.000000"), delay, algoname)
				time.Sleep(time.Millisecond * time.Duration(delay))
			}
			if l {
				log.Debugf("mining %d, old %08x new %08x average %3.2f trail %3.2f weighted %3.2f blocks in window: %d adjustment %0.1f%% algo %s delayed %dms",
					lastNode.height+1, last.bits, newTargetBits, allTimeAverage, trailTimeAverage, weighted*ttpb, counter, (1-adjustment)*100, fork.List[1].AlgoVers[algo], delay)
				if b.chainParams.Name == "testnet" && int64(lastNode.height) < b.chainParams.TargetTimePerBlock+1 && lastNode.height > 0 {
					time.Sleep(time.Second * time.Duration(b.chainParams.TargetTimePerBlock))
				}
			}
		}
		return newTargetBits, nil
	}
	nH := lastNode.height + 1
	// algo := fork.GetAlgoVer(algoname, nH)
	return fork.GetMinBits(algoname, nH), nil
}

// CalcNextRequiredDifficulty calculates the required difficulty for the block after the end of the current best chain based on the difficulty retarget rules. This function is safe for concurrent access.
func (b *BlockChain) CalcNextRequiredDifficulty(timestamp time.Time, algo string) (difficulty uint32, err error) {
	b.chainLock.Lock()
	difficulty, err = b.calcNextRequiredDifficulty(b.bestChain.Tip(), timestamp, algo, true)
	b.chainLock.Unlock()
	return
}
