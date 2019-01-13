package btcutil_test
import (
	"fmt"
	"math"
	"github.com/parallelcointeam/pod/btcutil"
)
func ExampleAmount() {
	a := btcutil.Amount(0)
	fmt.Println("Zero Satoshi:", a)
	a = btcutil.Amount(1e8)
	fmt.Println("100,000,000 Satoshis:", a)
	a = btcutil.Amount(1e5)
	fmt.Println("100,000 Satoshis:", a)
	// Output:
	// Zero Satoshi: 0 DUO
	// 100,000,000 Satoshis: 1 DUO
	// 100,000 Satoshis: 0.001 DUO
}
func ExampleNewAmount() {
	amountOne, err := btcutil.NewAmount(1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountOne) //Output 1
	amountFraction, err := btcutil.NewAmount(0.01234567)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountFraction) //Output 2
	amountZero, err := btcutil.NewAmount(0)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountZero) //Output 3
	amountNaN, err := btcutil.NewAmount(math.NaN())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountNaN) //Output 4
	// Output: 1 DUO
	// 0.01234567 DUO
	// 0 DUO
	// invalid bitcoin amount
}
func ExampleAmount_unitConversions() {
	amount := btcutil.Amount(44433322211100)
	fmt.Println("Satoshi to kDUO:", amount.Format(btcutil.AmountKiloDUO))
	fmt.Println("Satoshi to DUO:", amount)
	fmt.Println("Satoshi to MilliDUO:", amount.Format(btcutil.AmountMilliDUO))
	fmt.Println("Satoshi to MicroDUO:", amount.Format(btcutil.AmountMicroDUO))
	fmt.Println("Satoshi to Satoshi:", amount.Format(btcutil.AmountSatoshi))
	// Output:
	// Satoshi to kDUO: 444.333222111 kDUO
	// Satoshi to DUO: 444333.222111 DUO
	// Satoshi to MilliDUO: 444333222.111 mDUO
	// Satoshi to MicroDUO: 444333222111 μDUO
	// Satoshi to Satoshi: 44433322211100 Satoshi
}
