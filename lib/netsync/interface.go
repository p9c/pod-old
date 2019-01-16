package netsync

import (
	"git.parallelcoin.io/pod/lib/blockchain"
	"git.parallelcoin.io/pod/lib/chaincfg"
	"git.parallelcoin.io/pod/lib/chaincfg/chainhash"
	"git.parallelcoin.io/pod/lib/peer"
	"git.parallelcoin.io/pod/lib/util"
	"git.parallelcoin.io/pod/lib/wire"
	"git.parallelcoin.io/pod/module/node/mempool"
)

// PeerNotifier exposes methods to notify peers of status changes to transactions, blocks, etc. Currently server (in the main package) implements this interface.
type PeerNotifier interface {
	AnnounceNewTransactions(newTxs []*mempool.TxDesc)
	UpdatePeerHeights(latestBlkHash *chainhash.Hash, latestHeight int32, updateSource *peer.Peer)
	RelayInventory(invVect *wire.InvVect, data interface{})
	TransactionConfirmed(tx *util.Tx)
}

// Config is a configuration struct used to initialize a new SyncManager.
type Config struct {
	PeerNotifier       PeerNotifier
	Chain              *blockchain.BlockChain
	TxMemPool          *mempool.TxPool
	ChainParams        *chaincfg.Params
	DisableCheckpoints bool
	MaxPeers           int
	FeeEstimator       *mempool.FeeEstimator
}
