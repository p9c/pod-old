package main

import (
	"github.com/parallelcointeam/pod/btcutil"
	"github.com/parallelcointeam/pod/rpcclient"
	"github.com/parallelcointeam/pod/wire"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"
)

func main() {
	// Only override the handlers for notifications you care about. Also note most of these handlers will only be called if you register for notifications.  See the documentation of the rpcclient NotificationHandlers type for more details about each handler.
	ntfnHandlers := rpcclient.NotificationHandlers{
		OnFilteredBlockConnected: func(height int32, header *wire.BlockHeader, txns []*btcutil.Tx) {
			log.Printf("Block connected: %v (%d) %v",
				header.BlockHash(), height, header.Timestamp)
		},
		OnFilteredBlockDisconnected: func(height int32, header *wire.BlockHeader) {
			log.Printf("Block disconnected: %v (%d) %v",
				header.BlockHash(), height, header.Timestamp)
		},
	}
	// Connect to local pod RPC server using websockets.
	podHomeDir := btcutil.AppDataDir("pod", false)
	certs, err := ioutil.ReadFile(filepath.Join(podHomeDir, "rpc.cert"))
	if err != nil {
		log.Fatal(err)
	}
	connCfg := &rpcclient.ConnConfig{
		Host:         "localhost:11048",
		Endpoint:     "ws",
		User:         "yourrpcuser",
		Pass:         "yourrpcpass",
		Certificates: certs,
	}
	client, err := rpcclient.New(connCfg, &ntfnHandlers)
	if err != nil {
		log.Fatal(err)
	}
	// Register for block connect and disconnect notifications.
	if err := client.NotifyBlocks(); err != nil {
		log.Fatal(err)
	}
	log.Println("NotifyBlocks: Registration Complete")
	// Get the current block count.
	blockCount, err := client.GetBlockCount()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Block count: %d", blockCount)
	// For this example gracefully shutdown the client after 10 seconds. Ordinarily when to shutdown the client is highly application specific.
	log.Println("Client shutdown in 10 seconds...")
	time.AfterFunc(time.Second*10, func() {
		log.Println("Client shutting down...")
		client.Shutdown()
		log.Println("Client shutdown complete.")
	})
	// Wait until the client either shuts down gracefully (or the user terminates the process with Ctrl+C).
	client.WaitForShutdown()
}
