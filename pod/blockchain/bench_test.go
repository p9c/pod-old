package blockchain

import (
	"github.com/parallelcointeam/pod/btcutil"
	"testing"
)

// BenchmarkIsCoinBase performs a simple benchmark against the IsCoinBase function.
func BenchmarkIsCoinBase(b *testing.B) {
	tx, _ := btcutil.NewBlock(&Block100000).Tx(1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IsCoinBase(tx)
	}
}

// BenchmarkIsCoinBaseTx performs a simple benchmark against the IsCoinBaseTx function.
func BenchmarkIsCoinBaseTx(b *testing.B) {
	tx := Block100000.Transactions[1]
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IsCoinBaseTx(tx)
	}
}
