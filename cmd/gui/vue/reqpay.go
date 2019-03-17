package vue

import (
	"encoding/json"
	"fmt"

	"git.parallelcoin.io/dev/pod/cmd/gui/jdb"
)

// var VAB []AddressBook

type RequestedPayment struct {
	Time    string `json:"time"`
	Label   string `json:"label"`
	Address string `json:"address"`
	Amount  string `json:"amount"`
	Desc    string `json:"desc"`
}

type RequestedPaymentHistory struct {
	RequestedPayment []RequestedPayment `json:"label"`
}

func (ab *RequestedPaymentHistory) RequestedPaymentHistory() {

	reqpays, err := jdb.JDB.ReadAll("reqpay")

	if err != nil {

		fmt.Println("Error", err)
	}

	for _, f := range reqpays {

		var reqPay RequestedPayment

		if err := json.Unmarshal([]byte(f), &reqPay); err != nil {

			fmt.Println("Error", err)
		}

		ab.RequestedPayment = append(ab.RequestedPayment, reqPay)
	}

	fmt.Println("Ersssssssssssssssssssssssssssror", ab.RequestedPayment)

}

func (ab *RequestedPayment) RequestedPaymentWrite(time, label, address, amount, desc string) {

	ab.Time = time
	ab.Label = label
	ab.Address = address
	ab.Amount = amount
	ab.Desc = desc
	jdb.JDB.Write("addressbook", ab.Label, ab)
	fmt.Println("Ersssssssssssssssssssssssssssror", ab)

}
