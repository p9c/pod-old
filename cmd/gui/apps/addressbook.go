package apps

import (
	"encoding/json"
	"fmt"

	"git.parallelcoin.io/pod/cmd/gui/jdb"
	"git.parallelcoin.io/pod/cmd/gui/vue"
)

// var VAB []AddressBook

type AddressBookLabel struct {
	Label   string `json:"label"`
	Address string `json:"address"`
}
type AddressBook struct {
	AddressBookLabel []AddressBookLabel `json:"labels"`
}

func init() {

	vue.MODS["addressbook"] = AddressBook{}
	vue.MODS["addressbooklabel"] = AddressBookLabel{}
}
func (ab *AddressBook) AddressBookData() {

	ab.AddressBookLabel = nil
	addressbooks, err := jdb.JDB.ReadAll("addressbook")
	if err != nil {
		fmt.Println("Error", err)
	}

	for _, f := range addressbooks {
		var addressbook AddressBookLabel
		if err := json.Unmarshal([]byte(f), &addressbook); err != nil {
			fmt.Println("Error", err)
		}
		ab.AddressBookLabel = append(ab.AddressBookLabel, addressbook)
	}
	// fmt.Println("Ersssssssssssssssssssssssssssror", ab.AddressBookLabel)
}
func (ab *AddressBookLabel) AddressBookLabelWrite(label, address string) {
	ab.Label = label
	ab.Address = address
	jdb.JDB.Write("addressbook", ab.Label, ab)
	fmt.Println("Ersssssssssssssssssssssssssssror", ab)

}
func (ab *AddressBookLabel) AddressBookLabelDelete(label string) {
	if err := jdb.JDB.Delete("addressbook", label); err != nil {
		fmt.Println("Error", err)
	}
	fmt.Println("Ersssssssssssssssssssssssssssror", ab)
}
