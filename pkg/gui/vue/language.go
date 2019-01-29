package vue

import (
	"fmt"

	"git.parallelcoin.io/pod/pkg/gui/jdb"
)

// "github.com/parallelcointeam/pod/rpcclient"

type Language struct {
	TrueStory               string `json:"truestory"`
	Wallet                  string `json:"wallet"`
	Interface               string `json:"interface"`
	Network                 string `json:"network"`
	Security                string `json:"security"`
	Mining                  string `json:"mining"`
	Overview                string `json:"overview"`
	Send                    string `json:"send"`
	Receive                 string `json:"receive"`
	AddressBook             string `json:"addressbook"`
	History                 string `json:"history"`
	Balances                string `json:"balances"`
	CurrentBalance          string `json:"currentbalance"`
	Total                   string `json:"total"`
	Unconfirmed             string `json:"unconfirmed"`
	Available               string `json:"available"`
	LatestTransactions      string `json:"latesttransactions"`
	Type                    string `json:"type"`
	Address                 string `json:"address"`
	Date                    string `json:"date"`
	Amount                  string `json:"amount"`
	PayTo                   string `json:"payto"`
	Label                   string `json:"label"`
	PerPage                 string `json:"perpage"`
	Paginated               string `json:"paginated"`
	SentTransactionsHistory string `json:"senttransactionshistory"`
	RequestPaymentsHistory  string `json:"requestpaymentshistory"`
	Message                 string `json:"message"`
}

func (l *Language) LanguageData(lang string) {
	if err := jdb.JDB.Read("lang", lang, &l); err != nil {
		fmt.Println("Error", err)
	}
}
