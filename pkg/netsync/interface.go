package netsync

import (
	"git.parallelcoin.io/pod/pkg/blockchain"
	"git.parallelcoin.io/pod/pkg/chaincfg"
	"git.parallelcoin.io/pod/pkg/chaincfg/chainhash"
	"git.parallelcoin.io/pod/pkg/peer"
	"git.parallelcoin.io/pod/pkg/util"
	"git.parallelcoin.io/pod/pkg/wire"
	"git.parallelcoin.io/pod/cmd/node/mempool"
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
