package main

import (
	"github.com/parallelcointeam/pod/chaincfg/chainhash"
	"github.com/parallelcointeam/pod/database"
	"time"
)

// headersCmd defines the configuration options for the loadheaders command.
type headersCmd struct {
	Bulk bool `long:"bulk" description:"Use bulk loading of headers instead of one at a time"`
}

var (
	// headersCfg defines the configuration options for the command.
	headersCfg = headersCmd{
		Bulk: false,
	}
)

// Execute is the main entry point for the command.  It's invoked by the parser.
func (cmd *headersCmd) Execute(args []string) error {
	// Setup the global config options and ensure they are valid.
	if err := setupGlobalConfig(); err != nil {
		return err
	}
	// Load the block database.
	db, err := loadBlockDB()
	if err != nil {
		return err
	}
	defer db.Close()
	// NOTE: This code will only work for ffldb.  Ideally the package using the database would keep a metadata index of its own.
	blockIdxName := []byte("ffldb-blockidx")
	if !headersCfg.Bulk {
		err = db.View(func(tx database.Tx) error {
			totalHdrs := 0
			blockIdxBucket := tx.Metadata().Bucket(blockIdxName)
			blockIdxBucket.ForEach(func(k, v []byte) error {
				totalHdrs++
				return nil
			})
			log.Infof("Loading headers for %d blocks...", totalHdrs)
			numLoaded := 0
			startTime := time.Now()
			blockIdxBucket.ForEach(func(k, v []byte) error {
				var hash chainhash.Hash
				copy(hash[:], k)
				_, err := tx.FetchBlockHeader(&hash)
				if err != nil {
					return err
				}
				numLoaded++
				return nil
			})
			log.Infof("Loaded %d headers in %v", numLoaded,
				time.Since(startTime))
			return nil
		})
		return err
	}
	// Bulk load headers.
	err = db.View(func(tx database.Tx) error {
		blockIdxBucket := tx.Metadata().Bucket(blockIdxName)
		hashes := make([]chainhash.Hash, 0, 500000)
		blockIdxBucket.ForEach(func(k, v []byte) error {
			var hash chainhash.Hash
			copy(hash[:], k)
			hashes = append(hashes, hash)
			return nil
		})
		log.Infof("Loading headers for %d blocks...", len(hashes))
		startTime := time.Now()
		hdrs, err := tx.FetchBlockHeaders(hashes)
		if err != nil {
			return err
		}
		log.Infof("Loaded %d headers in %v", len(hdrs),
			time.Since(startTime))
		return nil
	})
	return err
}
