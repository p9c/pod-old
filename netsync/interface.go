package netsync

import (
	"git.parallelcoin.io/pod/blockchain"
	"git.parallelcoin.io/pod/chaincfg"
	"git.parallelcoin.io/pod/chaincfg/chainhash"
	"git.parallelcoin.io/pod/node/mempool"
	"git.parallelcoin.io/pod/peer"
	"git.parallelcoin.io/pod/util"
	"git.parallelcoin.io/pod/wire"
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
