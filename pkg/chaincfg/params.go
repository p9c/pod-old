package chaincfg

import (
	"encoding/hex"
	"errors"
	"math"
	"math/big"
	"strings"
	"time"

	"git.parallelcoin.io/pod/pkg/fork"

	"git.parallelcoin.io/pod/pkg/chaincfg/chainhash"
	"git.parallelcoin.io/pod/pkg/wire"
)

// These variables are the chain proof-of-work limit parameters for each default network.
var (
	// AllOnes is 32 bytes of 0xff, the maximum target
	AllOnes = func() big.Int {
		b := big.NewInt(1)
		t := make([]byte, 32)
		for i := range t {
			t[i] = ^byte(0)
		}
		b.SetBytes(t)
		return *b
	}()
	// mainPowLimit is the highest proof of work value a Parallelcoin block can have for the main network.
	mainPowLimit = func() big.Int {
		mplb, _ := hex.DecodeString("00000fffff000000000000000000000000000000000000000000000000000000")
		return *big.NewInt(0).SetBytes(mplb) //AllOnes.Rsh(&AllOnes, 0)
	}()
	MainPowLimit     = mainPowLimit
	MainPowLimitBits = BigToCompact(&MainPowLimit)
	scryptPowLimit   = func() big.Int {
		mplb, _ := hex.DecodeString("0000fffff0000000000000000000000000000000000000000000000000000000")
		return *big.NewInt(0).SetBytes(mplb) //AllOnes.Rsh(&AllOnes, 0)
	}()
	ScryptPowLimit     = scryptPowLimit
	ScryptPowLimitBits = BigToCompact(&scryptPowLimit)
	// regressionPowLbimit is the highest proof of work value a Bitcoin block can have for the regression test network.  It is the value 2^255 - 1, all ones, 256 bits.
	regressionPowLimit = &AllOnes
	testnetBits        = ScryptPowLimitBits
	testNet3PowLimit   = ScryptPowLimit
	// simNetPowLimit is the highest proof of work value a Bitcoin block can have for the simulation test network.  It is the value 2^255 - 1, all ones, 256 bits.
	simNetPowLimit                         = &AllOnes
	Interval                       int64   = 100
	MaxAdjustDown                  int64   = 10
	MaxAdjustUp                    int64   = 20
	TargetTimePerBlock             int64   = 300
	AveragingInterval              int64   = 10
	AveragingTargetTimespan                = TargetTimePerBlock * AveragingInterval
	TargetTimespan                         = Interval * TargetTimePerBlock
	TestnetCoefficient             float64 = 0.12
	TestnetInterval                int64   = 100
	TestnetMaxAdjustDown           int64   = 10
	TestnetMaxAdjustUp             int64   = 20
	TestnetTargetTimePerBlock      int64   = 9
	TestnetAveragingInterval       int64   = 1600
	TestnetAveragingTargetTimespan         = TestnetTargetTimePerBlock * TestnetAveragingInterval
	TestnetTargetTimespan                  = TestnetInterval * TestnetTargetTimePerBlock
)

// Checkpoint identifies a known good point in the block chain.  Using checkpoints allows a few optimizations for old blocks during initial download and also prevents forks from old blocks.
// Each checkpoint is selected based upon several factors.  See the documentation for blockchain.IsCheckpointCandidate for details on the selection criteria.
type Checkpoint struct {
	Height int32
	Hash   *chainhash.Hash
}

// DNSSeed identifies a DNS seed.
type DNSSeed struct {
	// Host defines the hostname of the seed.
	Host string
	// HasFiltering defines whether the seed supports filtering by service flags (wire.ServiceFlag).
	HasFiltering bool
}

// ConsensusDeployment defines details related to a specific consensus rule change that is voted in.  This is part of BIP0009.
type ConsensusDeployment struct {
	// BitNumber defines the specific bit number within the block version this particular soft-fork deployment refers to.
	BitNumber uint8
	// StartTime is the median block time after which voting on the deployment starts.
	StartTime uint64
	// ExpireTime is the median block time after which the attempted deployment expires.
	ExpireTime uint64
}

// Constants that define the deployment offset in the deployments field of the
// parameters for each deployment.  This is useful to be able to get the details
// of a specific deployment by name.
const (
	// DeploymentTestDummy defines the rule change deployment ID for testing
	// purposes.
	DeploymentTestDummy = iota
	// DeploymentCSV defines the rule change deployment ID for the CSV
	// soft-fork package. The CSV package includes the deployment of BIPS
	// 68, 112, and 113.
	DeploymentCSV
	// DeploymentSegwit defines the rule change deployment ID for the
	// Segregated Witness (segwit) soft-fork package. The segwit package
	// includes the deployment of BIPS 141, 142, 144, 145, 147 and 173.
	DeploymentSegwit
	// NOTE: DefinedDeployments must always come last since it is used to
	// determine how many defined deployments there currently are.
	// DefinedDeployments is the number of currently defined deployments.
	DefinedDeployments
)

// Params defines a Bitcoin network by its parameters.  These parameters may be
// used by Bitcoin applications to differentiate networks as well as addresses
// and keys for one network from those intended for use on another network.
type Params struct {
	// Name defines a human-readable identifier for the network.
	Name string
	// Net defines the magic bytes used to identify the network.
	Net wire.BitcoinNet
	// DefaultPort defines the default peer-to-peer port for the network.
	DefaultPort string
	// DNSSeeds defines a list of DNS seeds for the network that are used
	// as one method to discover peers.
	DNSSeeds []DNSSeed
	// GenesisBlock defines the first block of the chain.
	GenesisBlock *wire.MsgBlock
	// GenesisHash is the starting block hash.
	GenesisHash *chainhash.Hash
	// PowLimit defines the highest allowed proof of work value for a block
	// as a uint256.
	PowLimit *big.Int
	// PowLimitBits defines the highest allowed proof of work value for a
	// block in compact form.
	PowLimitBits uint32
	// These fields define the block heights at which the specified softfork
	// BIP became active.
	BIP0034Height int32
	BIP0065Height int32
	BIP0066Height int32
	// CoinbaseMaturity is the number of blocks required before newly mined
	// coins (coinbase transactions) can be spent.
	CoinbaseMaturity uint16
	// SubsidyReductionInterval is the interval of blocks before the subsidy
	// is reduced.
	SubsidyReductionInterval int32
	// TargetTimespan is the desired amount of time that should elapse
	// before the block difficulty requirement is examined to determine how
	// it should be changed in order to maintain the desired block
	// generation rate.
	TargetTimespan int64
	// TargetTimePerBlock is the desired amount of time to generate each
	// block. Same as TargetSpacing in legacy client.
	TargetTimePerBlock int64
	// RetargetAdjustmentFactor is the adjustment factor used to limit
	// the minimum and maximum amount of adjustment that can occur between
	// difficulty retargets.
	RetargetAdjustmentFactor int64
	// ReduceMinDifficulty defines whether the network should reduce the
	// minimum required difficulty after a long enough period of time has
	// passed without finding a block.  This is really only useful for test
	// networks and should not be set on a main network.
	ReduceMinDifficulty bool
	// MinDiffReductionTime is the amount of time after which the minimum
	// required difficulty should be reduced when a block hasn't been found.
	//
	// NOTE: This only applies if ReduceMinDifficulty is true.
	MinDiffReductionTime time.Duration
	// GenerateSupported specifies whether or not CPU mining is allowed.
	GenerateSupported bool
	// Checkpoints ordered from oldest to newest.
	Checkpoints []Checkpoint
	// These fields are related to voting on consensus rule changes as
	// defined by BIP0009.
	//
	// RuleChangeActivationThreshold is the number of blocks in a threshold
	// state retarget window for which a positive vote for a rule change
	// must be cast in order to lock in a rule change. It should typically
	// be 95% for the main network and 75% for test networks.
	//
	// MinerConfirmationWindow is the number of blocks in each threshold
	// state retarget window.
	//
	// Deployments define the specific consensus rule changes to be voted
	// on.
	RuleChangeActivationThreshold uint32
	MinerConfirmationWindow       uint32
	Deployments                   [DefinedDeployments]ConsensusDeployment
	// Mempool parameters
	RelayNonStdTxs bool
	// Human-readable part for Bech32 encoded segwit addresses, as defined
	// in BIP 173.
	Bech32HRPSegwit string
	// Address encoding magics
	PubKeyHashAddrID        byte // First byte of a P2PKH address
	ScriptHashAddrID        byte // First byte of a P2SH address
	PrivateKeyID            byte // First byte of a WIF private key
	WitnessPubKeyHashAddrID byte // First byte of a P2WPKH address
	WitnessScriptHashAddrID byte // First byte of a P2WSH address
	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID [4]byte
	HDPublicKeyID  [4]byte
	// BIP44 coin type used in the hierarchical deterministic path for
	// address generation.
	HDCoinType uint32
	// Parallelcoin specific difficulty adjustment parameters
	Coefficient             float64
	Interval                int64
	AveragingInterval       int64
	AveragingTargetTimespan int64
	MaxAdjustDown           int64
	MaxAdjustUp             int64
	TargetTimespanAdjDown   int64
	MinActualTimespan       int64
	MaxActualTimespan       int64
	// PowLimit defines the highest allowed proof of work value for a scrypt block as a uint256.
	ScryptPowLimit     *big.Int
	ScryptPowLimitBits uint32
}

// MainNetParams defines the network parameters for the main Bitcoin network.
var MainNetParams = Params{
	Name:        "mainnet",
	Net:         wire.MainNet,
	DefaultPort: "11047",
	DNSSeeds: []DNSSeed{
		{"seed1.parallelcoin.io", true},
		{"seed2.parallelcoin.io", true},
		{"seed3.parallelcoin.io", true},
		{"seed4.parallelcoin.io", true},
		{"seed5.parallelcoin.io", true},
	},
	// Chain parameters
	GenesisBlock:             &genesisBlock,
	GenesisHash:              &genesisHash,
	PowLimit:                 &mainPowLimit,
	PowLimitBits:             MainPowLimitBits, //0x1e0fffff,
	BIP0034Height:            1000000,          // Reserved for future change
	BIP0065Height:            1000000,
	BIP0066Height:            1000000,
	CoinbaseMaturity:         100,
	SubsidyReductionInterval: 250000,
	TargetTimespan:           TargetTimespan,
	TargetTimePerBlock:       TargetTimePerBlock,
	RetargetAdjustmentFactor: 2, // 50% less, 200% more (not used in parallelcoin)
	ReduceMinDifficulty:      false,
	MinDiffReductionTime:     0,
	GenerateSupported:        true,
	// Checkpoints ordered from oldest to newest.
	Checkpoints: []Checkpoint{
		// {11111, newHashFromStr("0000000069e244f73d78e8fd29ba2fd2ed618bd6fa2ee92559f542fdb26e7c1d")},
	},
	// Consensus rule change deployments.
	//
	// The miner confirmation window is defined as:
	//   target proof of work timespan / target proof of work spacing
	RuleChangeActivationThreshold: 1916, // 95% of MinerConfirmationWindow
	MinerConfirmationWindow:       2016, //
	Deployments: [DefinedDeployments]ConsensusDeployment{
		DeploymentTestDummy: {
			BitNumber:  28,
			StartTime:  1199145601, // January 1, 2008 UTC
			ExpireTime: 1230767999, // December 31, 2008 UTC
		},
		DeploymentCSV: {
			BitNumber:  0,
			StartTime:  1462060800, // May 1st, 2016
			ExpireTime: 1493596800, // May 1st, 2017
		},
		DeploymentSegwit: {
			BitNumber:  1,
			StartTime:  1479168000, // November 15, 2016 UTC
			ExpireTime: 1510704000, // November 15, 2017 UTC.
		},
	},
	// Mempool parameters
	RelayNonStdTxs: false,
	// Human-readable part for Bech32 encoded segwit addresses, as defined in
	// BIP 173.
	Bech32HRPSegwit: "bc", // always bc for main net
	// Address encoding magics
	PubKeyHashAddrID:        83,  // 0x00, // starts with 1
	ScriptHashAddrID:        9,   // 0x05, // starts with 3
	PrivateKeyID:            178, // 0x80, // starts with 5 (uncompressed) or K (compressed)
	WitnessPubKeyHashAddrID: 84,  // 0x06, // starts with p2
	WitnessScriptHashAddrID: 19,  // 0x0A, // starts with 7Xh
	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID: [4]byte{0x04, 0x88, 0xad, 0xe4}, // starts with xprv
	HDPublicKeyID:  [4]byte{0x04, 0x88, 0xb2, 0x1e}, // starts with xpub
	// BIP44 coin type used in the hierarchical deterministic path for
	// address generation.
	HDCoinType: 0,
	// Parallelcoin specific difficulty adjustment parameters
	Interval:                Interval,
	AveragingInterval:       AveragingInterval, // Extend to target timespan to adjust better to hashpower (30000/300=100) post hardfork
	AveragingTargetTimespan: AveragingTargetTimespan,
	MaxAdjustDown:           MaxAdjustDown,
	MaxAdjustUp:             MaxAdjustUp,
	TargetTimespanAdjDown:   AveragingTargetTimespan * (Interval + MaxAdjustDown) / Interval,
	MinActualTimespan:       AveragingTargetTimespan * (Interval - MaxAdjustUp) / Interval,
	MaxActualTimespan:       AveragingTargetTimespan * (Interval + MaxAdjustDown) / Interval,
	ScryptPowLimit:          &scryptPowLimit,
	ScryptPowLimitBits:      ScryptPowLimitBits,
}

// RegressionNetParams defines the network parameters for the regression test
// Bitcoin network.  Not to be confused with the test Bitcoin network (version
// 3), this network is sometimes simply called "testnet".
var RegressionNetParams = Params{
	Name:        "regtest",
	Net:         wire.TestNet,
	DefaultPort: "31047",
	DNSSeeds:    []DNSSeed{},
	// Chain parameters
	GenesisBlock:             &regTestGenesisBlock,
	GenesisHash:              &regTestGenesisHash,
	PowLimit:                 regressionPowLimit,
	PowLimitBits:             0x207fffff,
	CoinbaseMaturity:         100,
	BIP0034Height:            100000000, // Not active - Permit ver 1 blocks
	BIP0065Height:            100000000, // Used by regression tests
	BIP0066Height:            100000000, // Used by regression tests
	SubsidyReductionInterval: 150,
	TargetTimespan:           30000, // 14 days
	TargetTimePerBlock:       300,   // 5 minutes
	RetargetAdjustmentFactor: 2,     // 50% less, 200% more
	ReduceMinDifficulty:      true,
	MinDiffReductionTime:     300 * 2,
	GenerateSupported:        true,
	// Checkpoints ordered from oldest to newest.
	Checkpoints: nil,
	// Consensus rule change deployments.
	//
	// The miner confirmation window is defined as:
	//   target proof of work timespan / target proof of work spacing
	RuleChangeActivationThreshold: 108, // 75%  of MinerConfirmationWindow
	MinerConfirmationWindow:       144,
	Deployments: [DefinedDeployments]ConsensusDeployment{
		DeploymentTestDummy: {
			BitNumber:  28,
			StartTime:  0,             // Always available for vote
			ExpireTime: math.MaxInt64, // Never expires
		},
		DeploymentCSV: {
			BitNumber:  0,
			StartTime:  0,             // Always available for vote
			ExpireTime: math.MaxInt64, // Never expires
		},
		DeploymentSegwit: {
			BitNumber:  1,
			StartTime:  0,             // Always available for vote
			ExpireTime: math.MaxInt64, // Never expires.
		},
	},
	// Mempool parameters
	RelayNonStdTxs: true,
	// Human-readable part for Bech32 encoded segwit addresses, as defined in
	// BIP 173.
	Bech32HRPSegwit: "bcrt", // always bcrt for reg test net
	// Address encoding magics
	PubKeyHashAddrID: 0x00,
	ScriptHashAddrID: 0x05,
	PrivateKeyID:     0x80,
	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID: [4]byte{0x04, 0x35, 0x83, 0x94}, // starts with tprv
	HDPublicKeyID:  [4]byte{0x04, 0x35, 0x87, 0xcf}, // starts with tpub
	// BIP44 coin type used in the hierarchical deterministic path for
	// address generation.
	HDCoinType: 1,
	// Parallelcoin specific difficulty adjustment parameters
	Interval:                Interval,
	AveragingInterval:       10, // Extend to target timespan to adjust better to hashpower (30000/300=100) post hardfork
	AveragingTargetTimespan: AveragingTargetTimespan,
	MaxAdjustDown:           MaxAdjustDown,
	MaxAdjustUp:             MaxAdjustUp,
	TargetTimespanAdjDown:   AveragingTargetTimespan * (Interval + MaxAdjustDown) / Interval,
	MinActualTimespan:       AveragingTargetTimespan * (Interval - MaxAdjustUp) / Interval,
	MaxActualTimespan:       AveragingTargetTimespan * (Interval + MaxAdjustDown) / Interval,
	ScryptPowLimit:          &scryptPowLimit,
	ScryptPowLimitBits:      ScryptPowLimitBits,
}

// TestNet3Params defines the network parameters for the test Bitcoin network
// (version 3).  Not to be confused with the regression test network, this
// network is sometimes simply called "testnet".
var TestNet3Params = Params{
	Name:        "testnet",
	Net:         wire.TestNet3,
	DefaultPort: "21047",
	DNSSeeds:    []DNSSeed{
		// {"testnet-seed.bitcoin.jonasschnelli.ch", true},
	},
	// Chain parameters
	GenesisBlock:             &testNet3GenesisBlock,
	GenesisHash:              &testNet3GenesisHash,
	PowLimit:                 &fork.FirstPowLimit,    //fork&testNet3PowLimit,
	PowLimitBits:             fork.FirstPowLimitBits, //testnetBits,
	BIP0034Height:            1000000,                // 0000000023b3a96d3484e5abb3755c413e7d41500f8e2a5c3f0dd01299cd8ef8
	BIP0065Height:            1000000,                // 00000000007f6655f22f98e72ed80d8b06dc761d5da09df0fa1dc4be4f861eb6
	BIP0066Height:            1000000,                // 000000002104c8c45e99a8853285a3b592602a3ccde2b832481da85e9e4ba182
	CoinbaseMaturity:         100,
	SubsidyReductionInterval: 250000,
	TargetTimespan:           TestnetTargetTimespan,
	TargetTimePerBlock:       TestnetTargetTimePerBlock,
	RetargetAdjustmentFactor: 2,
	ReduceMinDifficulty:      false,
	MinDiffReductionTime:     0, // time.Minute * 10, // TargetTimePerBlock * 2
	GenerateSupported:        true,
	// Checkpoints ordered from oldest to newest.
	Checkpoints: []Checkpoint{
		// {546, newHashFromStr("000000002a936ca763904c3c35fce2f3556c559c0214345d31b1bcebf76acb70")},
	},
	// Consensus rule change deployments.
	//
	// The miner confirmation window is defined as:
	//   target proof of work timespan / target proof of work spacing
	RuleChangeActivationThreshold: 1512, // 75% of MinerConfirmationWindow
	MinerConfirmationWindow:       2016,
	Deployments: [DefinedDeployments]ConsensusDeployment{
		DeploymentTestDummy: {
			BitNumber:  28,
			StartTime:  1199145601, // January 1, 2008 UTC
			ExpireTime: 1230767999, // December 31, 2008 UTC
		},
		DeploymentCSV: {
			BitNumber:  0,
			StartTime:  1456790400, // March 1st, 2016
			ExpireTime: 1493596800, // May 1st, 2017
		},
		DeploymentSegwit: {
			BitNumber:  1,
			StartTime:  1462060800, // May 1, 2016 UTC
			ExpireTime: 1493596800, // May 1, 2017 UTC.
		},
	},
	// Mempool parameters
	RelayNonStdTxs: true,
	// Human-readable part for Bech32 encoded segwit addresses, as defined in
	// BIP 173.
	Bech32HRPSegwit: "tb", // always tb for test net
	// Address encoding magics
	PubKeyHashAddrID:        18,   // starts with m or n
	ScriptHashAddrID:        188,  // starts with 2
	WitnessPubKeyHashAddrID: 0x03, // starts with QW
	WitnessScriptHashAddrID: 0x28, // starts with T7n
	PrivateKeyID:            239,  // starts with 9 (uncompressed) or c (compressed)
	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID: [4]byte{0x04, 0x35, 0x83, 0x94}, // starts with tprv
	HDPublicKeyID:  [4]byte{0x04, 0x35, 0x87, 0xcf}, // starts with tpub
	// BIP44 coin type used in the hierarchical deterministic path for
	// address generation.
	HDCoinType: 1,
	// Parallelcoin specific difficulty adjustment parameters
	Coefficient:             TestnetCoefficient,
	Interval:                TestnetInterval,
	AveragingInterval:       TestnetAveragingInterval, // Extend to target timespan to adjust better to hashpower (30000/300=100) post hardforkTestnet
	AveragingTargetTimespan: TestnetAveragingTargetTimespan,
	MaxAdjustDown:           TestnetMaxAdjustDown,
	MaxAdjustUp:             TestnetMaxAdjustUp,
	TargetTimespanAdjDown:   TestnetAveragingTargetTimespan * (TestnetInterval + TestnetMaxAdjustDown) / TestnetInterval,
	MinActualTimespan:       TestnetAveragingTargetTimespan * (TestnetInterval - TestnetMaxAdjustUp) / TestnetInterval,
	MaxActualTimespan:       TestnetAveragingTargetTimespan * (TestnetInterval + TestnetMaxAdjustDown) / TestnetInterval,
	ScryptPowLimit:          &scryptPowLimit,
	ScryptPowLimitBits:      ScryptPowLimitBits,
}

// SimNetParams defines the network parameters for the simulation test Bitcoin
// network.  This network is similar to the normal test network except it is
// intended for private use within a group of individuals doing simulation
// testing.  The functionality is intended to differ in that the only nodes
// which are specifically specified are used to create the network rather than
// following normal discovery rules.  This is important as otherwise it would
// just turn into another public testnet.
var SimNetParams = Params{
	Name:        "simnet",
	Net:         wire.SimNet,
	DefaultPort: "41047",
	DNSSeeds:    []DNSSeed{}, // NOTE: There must NOT be any seeds.
	// Chain parameters
	GenesisBlock:             &simNetGenesisBlock,
	GenesisHash:              &simNetGenesisHash,
	PowLimit:                 simNetPowLimit,
	PowLimitBits:             0x207fffff,
	BIP0034Height:            0, // Always active on simnet
	BIP0065Height:            0, // Always active on simnet
	BIP0066Height:            0, // Always active on simnet
	CoinbaseMaturity:         100,
	SubsidyReductionInterval: 210000,
	TargetTimespan:           30000, // 14 days
	TargetTimePerBlock:       300,   // 10 minutes
	RetargetAdjustmentFactor: 2,     // 25% less, 400% more
	ReduceMinDifficulty:      true,
	MinDiffReductionTime:     time.Minute * 10, // TargetTimePerBlock * 2
	GenerateSupported:        true,
	// Checkpoints ordered from oldest to newest.
	Checkpoints: nil,
	// Consensus rule change deployments.
	//
	// The miner confirmation window is defined as:
	//   target proof of work timespan / target proof of work spacing
	RuleChangeActivationThreshold: 75, // 75% of MinerConfirmationWindow
	MinerConfirmationWindow:       100,
	Deployments: [DefinedDeployments]ConsensusDeployment{
		DeploymentTestDummy: {
			BitNumber:  28,
			StartTime:  0,             // Always available for vote
			ExpireTime: math.MaxInt64, // Never expires
		},
		DeploymentCSV: {
			BitNumber:  0,
			StartTime:  0,             // Always available for vote
			ExpireTime: math.MaxInt64, // Never expires
		},
		DeploymentSegwit: {
			BitNumber:  1,
			StartTime:  0,             // Always available for vote
			ExpireTime: math.MaxInt64, // Never expires.
		},
	},
	// Mempool parameters
	RelayNonStdTxs: true,
	// Human-readable part for Bech32 encoded segwit addresses, as defined in
	// BIP 173.
	Bech32HRPSegwit: "sb", // always sb for sim net
	// Address encoding magics
	PubKeyHashAddrID:        0x3f, // starts with S
	ScriptHashAddrID:        0x7b, // starts with s
	PrivateKeyID:            0x64, // starts with 4 (uncompressed) or F (compressed)
	WitnessPubKeyHashAddrID: 0x19, // starts with Gg
	WitnessScriptHashAddrID: 0x28, // starts with ?
	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID: [4]byte{0x04, 0x20, 0xb9, 0x00}, // starts with sprv
	HDPublicKeyID:  [4]byte{0x04, 0x20, 0xbd, 0x3a}, // starts with spub
	// BIP44 coin type used in the hierarchical deterministic path for
	// address generation.
	HDCoinType: 115, // ASCII for s
	// Parallelcoin specific difficulty adjustment parameters
	Interval:                100,
	AveragingInterval:       10, // Extend to target timespan to adjust better to hashpower (30000/300=100) post hardfork
	AveragingTargetTimespan: 10 * 300,
	MaxAdjustDown:           10,
	MaxAdjustUp:             20,
	TargetTimespanAdjDown:   300 * (100 + 10) / 100,
	MinActualTimespan:       10 * 300 * (100 - 20) / 100,
	MaxActualTimespan:       10 * 300 * (100 + 10) / 100,
	ScryptPowLimit:          &scryptPowLimit,
	ScryptPowLimitBits:      ScryptPowLimitBits,
}
var (
	// ErrDuplicateNet describes an error where the parameters for a Bitcoin
	// network could not be set due to the network already being a standard
	// network or previously-registered into this package.
	ErrDuplicateNet = errors.New("duplicate Bitcoin network")
	// ErrUnknownHDKeyID describes an error where the provided id which
	// is intended to identify the network for a hierarchical deterministic
	// private extended key is not registered.
	ErrUnknownHDKeyID = errors.New("unknown hd private extended key bytes")
)
var (
	registeredNets       = make(map[wire.BitcoinNet]struct{})
	pubKeyHashAddrIDs    = make(map[byte]struct{})
	scriptHashAddrIDs    = make(map[byte]struct{})
	bech32SegwitPrefixes = make(map[string]struct{})
	hdPrivToPubKeyIDs    = make(map[[4]byte][]byte)
)

// String returns the hostname of the DNS seed in human-readable form.
func (d DNSSeed) String() string {
	return d.Host
}

// Register registers the network parameters for a Bitcoin network.  This may
// error with ErrDuplicateNet if the network is already registered (either
// due to a previous Register call, or the network being one of the default
// networks).
// Network parameters should be registered into this package by a main package
// as early as possible.  Then, library packages may lookup networks or network
// parameters based on inputs and work regardless of the network being standard
// or not.
func Register(params *Params) error {
	if _, ok := registeredNets[params.Net]; ok {
		return ErrDuplicateNet
	}
	registeredNets[params.Net] = struct{}{}
	pubKeyHashAddrIDs[params.PubKeyHashAddrID] = struct{}{}
	scriptHashAddrIDs[params.ScriptHashAddrID] = struct{}{}
	hdPrivToPubKeyIDs[params.HDPrivateKeyID] = params.HDPublicKeyID[:]
	// A valid Bech32 encoded segwit address always has as prefix the
	// human-readable part for the given net followed by '1'.
	bech32SegwitPrefixes[params.Bech32HRPSegwit+"1"] = struct{}{}
	return nil
}

// mustRegister performs the same function as Register except it panics if there
// is an error.  This should only be called from package init functions.
func mustRegister(params *Params) {
	if err := Register(params); err != nil {
		panic("failed to register network: " + err.Error())
	}
}

// IsPubKeyHashAddrID returns whether the id is an identifier known to prefix a
// pay-to-pubkey-hash address on any default or registered network.  This is
// used when decoding an address string into a specific address type.  It is up
// to the caller to check both this and IsScriptHashAddrID and decide whether an
// address is a pubkey hash address, script hash address, neither, or
// undeterminable (if both return true).
func IsPubKeyHashAddrID(id byte) bool {
	_, ok := pubKeyHashAddrIDs[id]
	return ok
}

// IsScriptHashAddrID returns whether the id is an identifier known to prefix a
// pay-to-script-hash address on any default or registered network.  This is
// used when decoding an address string into a specific address type.  It is up
// to the caller to check both this and IsPubKeyHashAddrID and decide whether an
// address is a pubkey hash address, script hash address, neither, or
// undeterminable (if both return true).
func IsScriptHashAddrID(id byte) bool {
	_, ok := scriptHashAddrIDs[id]
	return ok
}

// IsBech32SegwitPrefix returns whether the prefix is a known prefix for segwit
// addresses on any default or registered network.  This is used when decoding
// an address string into a specific address type.
func IsBech32SegwitPrefix(prefix string) bool {
	prefix = strings.ToLower(prefix)
	_, ok := bech32SegwitPrefixes[prefix]
	return ok
}

// HDPrivateKeyToPublicKeyID accepts a private hierarchical deterministic
// extended key id and returns the associated public key id.  When the provided
// id is not registered, the ErrUnknownHDKeyID error will be returned.
func HDPrivateKeyToPublicKeyID(id []byte) ([]byte, error) {
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

// newHashFromStr converts the passed big-endian hex string into a
// chainhash.Hash.  It only differs from the one available in chainhash in that
// it panics on an error since it will only (and must only) be called with
// hard-coded, and therefore known good, hashes.
func newHashFromStr(hexStr string) *chainhash.Hash {
	hash, err := chainhash.NewHashFromStr(hexStr)
	if err != nil {
		// Ordinarily I don't like panics in library code since it
		// can take applications down without them having a chance to
		// recover which is extremely annoying, however an exception is
		// being made in this case because the only way this can panic
		// is if there is an error in the hard-coded hashes.  Thus it
		// will only ever potentially panic on init and therefore is
		// 100% predictable.
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

// CompactToBig converts a compact representation of a whole number N to an
// unsigned 32-bit number.  The representation is similar to IEEE754 floating
// point numbers.
// Like IEEE754 floating point, there are three basic components: the sign,
// the exponent, and the mantissa.  They are broken out as follows:
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
// This compact form is only used in bitcoin to encode unsigned 256-bit numbers
// which represent difficulty targets, thus there really is not a need for a
// sign bit, but it is implemented here to stay consistent with bitcoind.
func CompactToBig(compact uint32) *big.Int {
	// Extract the mantissa, sign bit, and exponent.
	mantissa := compact & 0x007fffff
	isNegative := compact&0x00800000 != 0
	exponent := uint(compact >> 24)
	// Since the base for the exponent is 256, the exponent can be treated
	// as the number of bytes to represent the full 256-bit number.  So,
	// treat the exponent as the number of bytes and shift the mantissa
	// right or left accordingly.  This is equivalent to:
	// N = mantissa * 256^(exponent-3)
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

// BigToCompact converts a whole number N to a compact representation using
// an unsigned 32-bit number.  The compact representation only provides 23 bits
// of precision, so values larger than (2^23 - 1) only encode the most
// significant digits of the number.  See CompactToBig for details.
func BigToCompact(n *big.Int) uint32 {
	// No need to do any work if it's zero.
	if n.Sign() == 0 {
		return 0
	}
	// Since the base for the exponent is 256, the exponent can be treated
	// as the number of bytes.  So, shift the number right or left
	// accordingly.  This is equivalent to:
	// mantissa = mantissa / 256^(exponent-3)
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
	// When the mantissa already has the sign bit set, the number is too
	// large to fit into the available 23-bits, so divide the number by 256
	// and increment the exponent accordingly.
	if mantissa&0x00800000 != 0 {
		mantissa >>= 8
		exponent++
	}
	// Pack the exponent, sign bit, and mantissa into an unsigned 32-bit
	// int and return it.
	compact := uint32(exponent<<24) | mantissa
	if n.Sign() < 0 {
		compact |= 0x00800000
	}
	return compact
}
