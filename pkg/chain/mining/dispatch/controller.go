package controller

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"git.parallelcoin.io/pod/pkg/util/clog"

	"git.parallelcoin.io/pod/pkg/chain"
	"git.parallelcoin.io/pod/pkg/chain/config"
	"git.parallelcoin.io/pod/pkg/chain/fork"
	"git.parallelcoin.io/pod/pkg/chain/mining"
	"git.parallelcoin.io/pod/pkg/util"
	"git.parallelcoin.io/pod/pkg/chain/wire"
	"github.com/xtaci/kcp-go"
)

const (

	// maxNonce is the maximum value a nonce can be in a block header.
	maxNonce = ^uint32(0) // 2^32 - 1

	// maxExtraNonce is the maximum value an extra nonce used in a coinbase transaction can be.
	maxExtraNonce = ^uint64(0) // 2^64 - 1
)

// Config is a descriptor containing the controller configuration.
type Config struct {

	// Blockchain gives access for the miner to information about the chain
	Blockchain *blockchain.BlockChain

	// ChainParams identifies which chain parameters the cpu miner is associated with.
	ChainParams *chaincfg.Params

	// BlockTemplateGenerator identifies the instance to use in order to generate block templates that the miner will attempt to solve.
	BlockTemplateGenerator *mining.BlkTmplGenerator

	// MiningAddrs is a list of payment addresses to use for the generated blocks.  Each generated block will randomly choose one of thec.
	MiningAddrs []util.Address

	// ProcessBlock defines the function to call with any solved blocks. It typically must run the provided block through the same set of rules and handling as any other block coming from the network.
	ProcessBlock func(*util.Block, blockchain.BehaviorFlags) (bool, error)

	// MinerListener is the listener that will accept miner subscriptions and such
	MinerListener string

	// MinerKey is generated from the password specified in the main configuration for miner port using Stribog hash to derive the nonce, Argon2i to expand the password, and a final pass of Keccak
	MinerKey []byte

	// ConnectedCount defines the function to use to obtain how many other peers the server is connected to.  This is used by the automatic persistent mining routine to determine whether or it should attempt mining.  This is useful because there is no point in mining when not connected to any peers since there would no be anyone to send any found blocks to.
	ConnectedCount func() int32

	// IsCurrent defines the function to use to obtain whether or not the block chain is current.  This is used by the automatic persistent mining routine to determine whether or it should attempt mining. This is useful because there is no point in mining if the chain is not current since any solved blocks would be on a side chain and and up orphaned anyways.
	IsCurrent func() bool
}

// Controller delivers new work to miner clients
type Controller struct {
	sync.Mutex
	b                *blockchain.BlockChain
	g                *mining.BlkTmplGenerator
	cfg              Config
	started          bool
	submitBlockLock  sync.Mutex
	wg               sync.WaitGroup
	workerWg         sync.WaitGroup
	updateNumWorkers chan struct{}
	quit             chan struct{}
	listener         *kcp.Listener
	workers          []*kcp.UDPSession
}

// submitBlock submits the passed block to network after ensuring it passes all of the consensus validation rules.
func (c *Controller) submitBlock(block *util.Block) bool {
	c.submitBlockLock.Lock()
	defer c.submitBlockLock.Unlock()

	// Ensure the block is not stale since a new block could have shown up while the solution was being found.  Typically that condition is detected and all work on the stale block is halted to start work on a new block, but the check only happens periodically, so it is possible a block was found and submitted in between.
	msgBlock := block.MsgBlock()
	if !msgBlock.Header.PrevBlock.IsEqual(&c.g.BestSnapshot().Hash) {

		log <- cl.Debugf{
			"Block submitted via miner with previous block %s is stale",
			msgBlock.Header.PrevBlock,
		}
		return false
	}

	// Process this block using the same rules as blocks coming from other nodes.  This will in turn relay it to the network like normal.
	isOrphan, err := c.cfg.ProcessBlock(block, blockchain.BFNone)
	if err != nil {
		// Anything other than a rule violation is an unexpected error, so log that error as an internal error.
		if _, ok := err.(blockchain.RuleError); !ok {
			log <- cl.Error{
				"Unexpected error while processing block submitted via miner worker:",
				err,
			}
			return false
		}
		log <- cl.Debug{"Block submitted via miner rejected:", err}
		return false
	}
	if isOrphan {
		log <- cl.Dbg("Block submitted via miner is an orphan")
		return false
	}

	// The block was accepted.
	coinbaseTx := block.MsgBlock().Transactions[0].TxOut[0]
	prevHeight := block.Height() - 1
	prevBlock, _ := c.b.BlockByHeight(prevHeight)
	prevTime := prevBlock.MsgBlock().Header.Timestamp.Unix()
	since := block.MsgBlock().Header.Timestamp.Unix() - prevTime

	Log.Infc(func() string {
		return fmt.Sprintf(
			"new block height %d %s %10d %08x %v %s %ds since prev",
			block.Height(),
			block.MsgBlock().BlockHashWithAlgos(block.Height()),
			block.MsgBlock().Header.Timestamp.Unix(),
			block.MsgBlock().Header.Bits,
			util.Amount(coinbaseTx.Value),
			fork.GetAlgoName(block.MsgBlock().Header.Version,
				block.Height()),
			since,
		)
	})
	return true
}

// solveBlock attempts to find some combination of a nonce, extra nonce, and current timestamp which makes the passed block hash to a value less than the target difficulty.  The timestamp is updated periodically and the passed block is modified with all tweaks during this process.  This means that when the function returns true, the block is ready for submission. This function will return early with false when conditions that trigger a stale block such as a new block showing up or periodically when there are new transactions and enough time has elapsed without finding a solution.
func (c *Controller) solveBlock(msgBlock *wire.MsgBlock, blockHeight int32, testnet bool, ticker *time.Ticker, submissionReceived chan *wire.MsgBlock, quit chan struct{}) bool {

	// Choose a random extra nonce offset for this block template and worker.
	enOffset, err := wire.RandomUint64()
	if err != nil {
		log <- cl.Error{"Unexpected error while generating random extra nonce offset:", err}
		enOffset = 0
	}
	header := &msgBlock.Header
	targetDifficulty := blockchain.CompactToBig(header.Bits)
	lastGenerated := time.Now()
	lastTxUpdate := c.g.TxSource().LastUpdated()

	// Note that the entire extra nonce range is iterated and the offset is added relying on the fact that overflow will wrap around 0 as provided by the Go spec.
	for extraNonce := uint64(0); extraNonce < maxExtraNonce; extraNonce++ {
		// Update the extra nonce in the block template with the new value by regenerating the coinbase script and setting the merkle root to the new value.
		c.g.UpdateExtraNonce(msgBlock, blockHeight, extraNonce+enOffset)
		// Search through the entire nonce range for a solution while periodically checking for early quit and stale block conditions along with updates to the speed monitor.
		for i := uint32(0); i <= maxNonce; i++ {
			select {
			case <-quit:
				// fmt.Println("chan:<-quit")
				return false
			case <-ticker.C:
				// fmt.Println("chan:<-ticker.C")
				// The current block is stale if the best block has changed.
				best := c.g.BestSnapshot()
				if !header.PrevBlock.IsEqual(&best.Hash) {

					return false
				}
				// The current block is stale if the memory pool has been updated since the block template was generated and it has been at least one minute.
				if lastTxUpdate != c.g.TxSource().LastUpdated() &&
					time.Now().After(lastGenerated.Add(time.Minute)) {

					return false
				}
				c.g.UpdateBlockTime(msgBlock)
			case <-submissionReceived:
				// fmt.Println("chan:<-submissionReceived")
				// Here we will send out the updated block to subscribed client workers
			default:
			}
			// Instead of directly attempting to solve blocks, here we send out the new block data to subscribed miner clients and process submissions they return when they find a solution

			header.Nonce = i
			hash := header.BlockHashWithAlgos(int32(fork.GetCurrent(blockHeight)))
			// The block is solved when the new block hash is less than the target difficulty.  Yay!
			if blockchain.HashToBig(&hash).Cmp(targetDifficulty) <= 0 {
				return true
			}
		}
	}
	return false
}

// generateBlocks is a worker that is controlled by the miningWorkerController. It is self contained in that it creates block templates and attempts to solve them while detecting when it is performing stale work and reacting accordingly by generating a new block template.  When a block is solved, it is submitted. It must be run as a goroutine.
func (c *Controller) generateBlocks(quit chan struct{}) {


	// Start a ticker which is used to signal checks for stale work and updates to the speed monitor.
	ticker := time.NewTicker(time.Second / 2)
	defer ticker.Stop()

	// Create a channel to receive block submissions
	var submission chan *wire.MsgBlock
out:
	for {
		select {
		case <-quit: // Quit when the miner is stopped.
			// fmt.Println("chan:<-quit")
			break out
		default: // Non-blocking select to fall through
		}
		// Wait until there is a connection to at least one other peer since there is no way to relay a found block or receive transactions to work on when there are no connected peers.
		if c.cfg.ConnectedCount() == 0 {
			time.Sleep(time.Second)
			continue
		}
		// No point in searching for a solution before the chain is synced.  Also, grab the same lock as used for block submission, since the current block will be changing and this would otherwise end up building a new block template on a block that is in the process of becoming stale.
		c.submitBlockLock.Lock()
		curHeight := c.g.BestSnapshot().Height
		if curHeight != 0 && !c.cfg.IsCurrent() {

			c.submitBlockLock.Unlock()
			time.Sleep(time.Second)
			continue
		}
		// Choose a payment address at random
		rand.Seed(time.Now().UnixNano())
		payToAddr := c.cfg.MiningAddrs[rand.Intn(len(c.cfg.MiningAddrs))]
		// Create a new block template using the available transactions in the memory pool as a source of transactions to potentially include in the block.
		template, err := c.g.NewBlockTemplate(payToAddr, "")
		c.submitBlockLock.Unlock()
		if err != nil {
			log <- cl.Error{"Failed to create new block template: %v", err}
			continue
		}
		// Attempt to solve the block.  The function will exit early with false when conditions that trigger a stale block, so a new block template can be generated.  When the return is true a solution was found, so submit the solved block.
		if c.solveBlock(template.Block, curHeight+1, c.cfg.ChainParams.Name == "testnet", ticker, submission, quit) {

			block := util.NewBlock(template.Block)
			c.submitBlock(block)
		}
	}
	c.workerWg.Done()
}

func (c *Controller) minerController() {

	c.workerWg.Add(1)
	quit := make(chan struct{})
	go c.generateBlocks(quit)
out:
	for {
		select {
		case <-c.quit:
			// fmt.Println("chan:<-c.quit")
			close(quit)
			break out
		}
	}
}

// Start begins the miner controller process. Calling this function when the miner controller has already been started will have no effect.
func (c *Controller) Start() {

	c.Lock()
	defer c.Unlock()
	if c.started {
		return
	}

	c.quit = make(chan struct{})
	c.wg.Add(1)
	go c.minerController()
	c.started = true
	log <- cl.Inf("Miner controller started")
}

// Stop gracefully stops the mining process by signalling all workers, and the speed monitor to quit.  Calling this function when the miner controller has not already been started will have no effect.
func (c *Controller) Stop() {

	c.Lock()
	defer c.Unlock()
	if !c.started {
		return
	}
	close(c.quit)
	c.wg.Wait()
	c.started = false
	log <- cl.Inf("Miner controller stopped")
}

// IsMining returns whether or not the miner controller has been started and is therefore currenting mining. This function is safe for concurrent access.
func (c *Controller) IsMining() bool {
	c.Lock()
	defer c.Unlock()
	return c.started
}

// New returns a new instance of a CPU miner for the provided configuration. Use Start to begin the mining process.  See the documentation for Controller type for more details.
func New(
	cfg *Config) *Controller {
	return &Controller{
		b:   cfg.Blockchain,
		g:   cfg.BlockTemplateGenerator,
		cfg: *cfg,
	}
}
