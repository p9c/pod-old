package controller

import (
	"errors"
	"github.com/parallelcointeam/pod/chaincfg/chainhash"
)

// Address is the parameter and reply type for subscriptions
type Address struct {
	Address string
}

// Kopach is the protocol root for the kopach work protocol, a list of addresses of subscribed miner workers
type Kopach []Address

// Work is the data required to construct a valid block to solve given the option of version number and corresponding hash algorithm
type Work struct {
	PrevBlockHash *chainhash.Hash
	MerkleRoot    *chainhash.Hash
	TimeStamp     uint32
	Difficulties  [9]uint32
}

// Subscribe adds an address to the list of subscribers to push work to
func (k *Kopach) Subscribe(args *Address, reply *Address) (err error) {
	log.Info("Subscribe called with", *args)
	err = errors.New("already subscribed")
	for i := range *k {
		if (*k)[i].Address == (*args).Address {
			return
		}
	}
	err = nil
	*k = append(*k, *args)
	reply = args
	return
}

// Unsubscribe removes an address from the list of subscribers to push work to
func (k *Kopach) Unsubscribe(args *Address, reply *Address) (err error) {
	log.Info("Unsubscribe called with", *args)
	err = errors.New("not subscribed")
	for i := range *k {
		if (*k)[i].Address == (*args).Address {
			err = nil
			if len(*k)-1 > i {
				*k = append((*k)[:i], (*k)[i+1:]...)
			} else {
				*k = (*k)[:i]
			}
		}
	}
	err = nil
	log.Info("sending reply", *args)
	*reply = *args
	return
}
