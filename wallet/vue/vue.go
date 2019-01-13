package vue

import (
	"fmt"
	"log"

	"git.parallelcoin.io/pod/json"
	"git.parallelcoin.io/pod/rpcclient"
)

// BlockChainData is the response from getinfo
type BlockChainData struct {
	GetInfo *json.InfoWalletResult `json:"getinfo"`
}

// GetBlockChainData requests a getinfo command from the RPC
func (k *BlockChainData) GetBlockChainData() {
	connCfg := &rpcclient.ConnConfig{
		Host:     "localhost:11046",
		Endpoint: "ws",
		User:     "user",
		Pass:     "pa55word",
		TLS:      true,
	}
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Get the list of unspent transaction outputs (utxos) that the
	// connected wallet has at least one private key for.
	info, err := client.GetInfo()
	if err != nil {
		log.Fatal(err)
	}

	k.GetInfo = info
	// result := []byte(info.Result())
	// json.Unmarshal(bytes(info), &k.GetInfo)
	fmt.Println("fdffffff", k.GetInfo)

	// For this example gracefully shutdown the client after 10 seconds.
	// Ordinarily when to shutdown the client is highly application
	// specific.

}
