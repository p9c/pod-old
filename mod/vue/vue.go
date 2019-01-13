package vue

import (
	"fmt"
	"log"

	"github.com/parallelcointeam/pod/btcjson"
	"github.com/parallelcointeam/pod/rpcclient"
)

type BlockChainData struct {
	GetInfo *btcjson.InfoWalletResult `json:"getinfo"`
}

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

	// fmt.Println("Daaaaaaaa", info.Blocks)
	k.GetInfo = info
	// result := []byte(info.Result())
	// json.Unmarshal(bytes(info), &k.GetInfo)
	fmt.Println("fdffffff", k.GetInfo)

	// For this example gracefully shutdown the client after 10 seconds.
	// Ordinarily when to shutdown the client is highly application
	// specific.

}
