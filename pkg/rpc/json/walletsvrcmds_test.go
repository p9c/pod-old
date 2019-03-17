package json_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"git.parallelcoin.io/dev/pod/pkg/rpc/json"
)

// TestWalletSvrCmds tests all of the wallet server commands marshal and unmarshal into valid results include handling of optional fields being omitted in the marshalled command, while optional fields with defaults have the default assigned on unmarshalled commands.
func TestWalletSvrCmds(
	t *testing.T) {

	t.Parallel()
	testID := int(1)
	tests := []struct {
		name         string
		newCmd       func() (interface{}, error)
		staticCmd    func() interface{}
		marshalled   string
		unmarshalled interface{}
	}{
		{
			name: "addmultisigaddress",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("addmultisigaddress", 2, []string{"031234", "035678"})
			},
			staticCmd: func() interface{} {

				keys := []string{"031234", "035678"}
				return json.NewAddMultisigAddressCmd(2, keys, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"addmultisigaddress","params":[2,["031234","035678"]],"id":1}`,
			unmarshalled: &json.AddMultisigAddressCmd{
				NRequired: 2,
				Keys:      []string{"031234", "035678"},
				Account:   nil,
			},
		},
		{
			name: "addmultisigaddress optional",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("addmultisigaddress", 2, []string{"031234", "035678"}, "test")
			},
			staticCmd: func() interface{} {

				keys := []string{"031234", "035678"}
				return json.NewAddMultisigAddressCmd(2, keys, json.String("test"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"addmultisigaddress","params":[2,["031234","035678"],"test"],"id":1}`,
			unmarshalled: &json.AddMultisigAddressCmd{
				NRequired: 2,
				Keys:      []string{"031234", "035678"},
				Account:   json.String("test"),
			},
		},
		{
			name: "addwitnessaddress",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("addwitnessaddress", "1address")
			},
			staticCmd: func() interface{} {

				return json.NewAddWitnessAddressCmd("1address")
			},
			marshalled: `{"jsonrpc":"1.0","method":"addwitnessaddress","params":["1address"],"id":1}`,
			unmarshalled: &json.AddWitnessAddressCmd{
				Address: "1address",
			},
		},
		{
			name: "createmultisig",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("createmultisig", 2, []string{"031234", "035678"})
			},
			staticCmd: func() interface{} {

				keys := []string{"031234", "035678"}
				return json.NewCreateMultisigCmd(2, keys)
			},
			marshalled: `{"jsonrpc":"1.0","method":"createmultisig","params":[2,["031234","035678"]],"id":1}`,
			unmarshalled: &json.CreateMultisigCmd{
				NRequired: 2,
				Keys:      []string{"031234", "035678"},
			},
		},
		{
			name: "dumpprivkey",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("dumpprivkey", "1Address")
			},
			staticCmd: func() interface{} {

				return json.NewDumpPrivKeyCmd("1Address")
			},
			marshalled: `{"jsonrpc":"1.0","method":"dumpprivkey","params":["1Address"],"id":1}`,
			unmarshalled: &json.DumpPrivKeyCmd{
				Address: "1Address",
			},
		},
		{
			name: "encryptwallet",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("encryptwallet", "pass")
			},
			staticCmd: func() interface{} {

				return json.NewEncryptWalletCmd("pass")
			},
			marshalled: `{"jsonrpc":"1.0","method":"encryptwallet","params":["pass"],"id":1}`,
			unmarshalled: &json.EncryptWalletCmd{
				Passphrase: "pass",
			},
		},
		{
			name: "estimatefee",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("estimatefee", 6)
			},
			staticCmd: func() interface{} {

				return json.NewEstimateFeeCmd(6)
			},
			marshalled: `{"jsonrpc":"1.0","method":"estimatefee","params":[6],"id":1}`,
			unmarshalled: &json.EstimateFeeCmd{
				NumBlocks: 6,
			},
		},
		{
			name: "estimatepriority",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("estimatepriority", 6)
			},
			staticCmd: func() interface{} {

				return json.NewEstimatePriorityCmd(6)
			},
			marshalled: `{"jsonrpc":"1.0","method":"estimatepriority","params":[6],"id":1}`,
			unmarshalled: &json.EstimatePriorityCmd{
				NumBlocks: 6,
			},
		},
		{
			name: "getaccount",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("getaccount", "1Address")
			},
			staticCmd: func() interface{} {

				return json.NewGetAccountCmd("1Address")
			},
			marshalled: `{"jsonrpc":"1.0","method":"getaccount","params":["1Address"],"id":1}`,
			unmarshalled: &json.GetAccountCmd{
				Address: "1Address",
			},
		},
		{
			name: "getaccountaddress",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("getaccountaddress", "acct")
			},
			staticCmd: func() interface{} {

				return json.NewGetAccountAddressCmd("acct")
			},
			marshalled: `{"jsonrpc":"1.0","method":"getaccountaddress","params":["acct"],"id":1}`,
			unmarshalled: &json.GetAccountAddressCmd{
				Account: "acct",
			},
		},
		{
			name: "getaddressesbyaccount",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("getaddressesbyaccount", "acct")
			},
			staticCmd: func() interface{} {

				return json.NewGetAddressesByAccountCmd("acct")
			},
			marshalled: `{"jsonrpc":"1.0","method":"getaddressesbyaccount","params":["acct"],"id":1}`,
			unmarshalled: &json.GetAddressesByAccountCmd{
				Account: "acct",
			},
		},
		{
			name: "getbalance",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("getbalance")
			},
			staticCmd: func() interface{} {

				return json.NewGetBalanceCmd(nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getbalance","params":[],"id":1}`,
			unmarshalled: &json.GetBalanceCmd{
				Account: nil,
				MinConf: json.Int(1),
			},
		},
		{
			name: "getbalance optional1",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("getbalance", "acct")
			},
			staticCmd: func() interface{} {

				return json.NewGetBalanceCmd(json.String("acct"), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getbalance","params":["acct"],"id":1}`,
			unmarshalled: &json.GetBalanceCmd{
				Account: json.String("acct"),
				MinConf: json.Int(1),
			},
		},
		{
			name: "getbalance optional2",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("getbalance", "acct", 6)
			},
			staticCmd: func() interface{} {

				return json.NewGetBalanceCmd(json.String("acct"), json.Int(6))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getbalance","params":["acct",6],"id":1}`,
			unmarshalled: &json.GetBalanceCmd{
				Account: json.String("acct"),
				MinConf: json.Int(6),
			},
		},
		{
			name: "getnewaddress",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("getnewaddress")
			},
			staticCmd: func() interface{} {

				return json.NewGetNewAddressCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getnewaddress","params":[],"id":1}`,
			unmarshalled: &json.GetNewAddressCmd{
				Account: nil,
			},
		},
		{
			name: "getnewaddress optional",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("getnewaddress", "acct")
			},
			staticCmd: func() interface{} {

				return json.NewGetNewAddressCmd(json.String("acct"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getnewaddress","params":["acct"],"id":1}`,
			unmarshalled: &json.GetNewAddressCmd{
				Account: json.String("acct"),
			},
		},
		{
			name: "getrawchangeaddress",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("getrawchangeaddress")
			},
			staticCmd: func() interface{} {

				return json.NewGetRawChangeAddressCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawchangeaddress","params":[],"id":1}`,
			unmarshalled: &json.GetRawChangeAddressCmd{
				Account: nil,
			},
		},
		{
			name: "getrawchangeaddress optional",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("getrawchangeaddress", "acct")
			},
			staticCmd: func() interface{} {

				return json.NewGetRawChangeAddressCmd(json.String("acct"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawchangeaddress","params":["acct"],"id":1}`,
			unmarshalled: &json.GetRawChangeAddressCmd{
				Account: json.String("acct"),
			},
		},
		{
			name: "getreceivedbyaccount",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("getreceivedbyaccount", "acct")
			},
			staticCmd: func() interface{} {

				return json.NewGetReceivedByAccountCmd("acct", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getreceivedbyaccount","params":["acct"],"id":1}`,
			unmarshalled: &json.GetReceivedByAccountCmd{
				Account: "acct",
				MinConf: json.Int(1),
			},
		},
		{
			name: "getreceivedbyaccount optional",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("getreceivedbyaccount", "acct", 6)
			},
			staticCmd: func() interface{} {

				return json.NewGetReceivedByAccountCmd("acct", json.Int(6))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getreceivedbyaccount","params":["acct",6],"id":1}`,
			unmarshalled: &json.GetReceivedByAccountCmd{
				Account: "acct",
				MinConf: json.Int(6),
			},
		},
		{
			name: "getreceivedbyaddress",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("getreceivedbyaddress", "1Address")
			},
			staticCmd: func() interface{} {

				return json.NewGetReceivedByAddressCmd("1Address", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getreceivedbyaddress","params":["1Address"],"id":1}`,
			unmarshalled: &json.GetReceivedByAddressCmd{
				Address: "1Address",
				MinConf: json.Int(1),
			},
		},
		{
			name: "getreceivedbyaddress optional",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("getreceivedbyaddress", "1Address", 6)
			},
			staticCmd: func() interface{} {

				return json.NewGetReceivedByAddressCmd("1Address", json.Int(6))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getreceivedbyaddress","params":["1Address",6],"id":1}`,
			unmarshalled: &json.GetReceivedByAddressCmd{
				Address: "1Address",
				MinConf: json.Int(6),
			},
		},
		{
			name: "gettransaction",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("gettransaction", "123")
			},
			staticCmd: func() interface{} {

				return json.NewGetTransactionCmd("123", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"gettransaction","params":["123"],"id":1}`,
			unmarshalled: &json.GetTransactionCmd{
				Txid:             "123",
				IncludeWatchOnly: json.Bool(false),
			},
		},
		{
			name: "gettransaction optional",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("gettransaction", "123", true)
			},
			staticCmd: func() interface{} {

				return json.NewGetTransactionCmd("123", json.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"gettransaction","params":["123",true],"id":1}`,
			unmarshalled: &json.GetTransactionCmd{
				Txid:             "123",
				IncludeWatchOnly: json.Bool(true),
			},
		},
		{
			name: "getwalletinfo",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("getwalletinfo")
			},
			staticCmd: func() interface{} {

				return json.NewGetWalletInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getwalletinfo","params":[],"id":1}`,
			unmarshalled: &json.GetWalletInfoCmd{},
		},
		{
			name: "importprivkey",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("importprivkey", "abc")
			},
			staticCmd: func() interface{} {

				return json.NewImportPrivKeyCmd("abc", nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"importprivkey","params":["abc"],"id":1}`,
			unmarshalled: &json.ImportPrivKeyCmd{
				PrivKey: "abc",
				Label:   nil,
				Rescan:  json.Bool(true),
			},
		},
		{
			name: "importprivkey optional1",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("importprivkey", "abc", "label")
			},
			staticCmd: func() interface{} {

				return json.NewImportPrivKeyCmd("abc", json.String("label"), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"importprivkey","params":["abc","label"],"id":1}`,
			unmarshalled: &json.ImportPrivKeyCmd{
				PrivKey: "abc",
				Label:   json.String("label"),
				Rescan:  json.Bool(true),
			},
		},
		{
			name: "importprivkey optional2",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("importprivkey", "abc", "label", false)
			},
			staticCmd: func() interface{} {

				return json.NewImportPrivKeyCmd("abc", json.String("label"), json.Bool(false))
			},
			marshalled: `{"jsonrpc":"1.0","method":"importprivkey","params":["abc","label",false],"id":1}`,
			unmarshalled: &json.ImportPrivKeyCmd{
				PrivKey: "abc",
				Label:   json.String("label"),
				Rescan:  json.Bool(false),
			},
		},
		{
			name: "keypoolrefill",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("keypoolrefill")
			},
			staticCmd: func() interface{} {

				return json.NewKeyPoolRefillCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"keypoolrefill","params":[],"id":1}`,
			unmarshalled: &json.KeyPoolRefillCmd{
				NewSize: json.Uint(100),
			},
		},
		{
			name: "keypoolrefill optional",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("keypoolrefill", 200)
			},
			staticCmd: func() interface{} {

				return json.NewKeyPoolRefillCmd(json.Uint(200))
			},
			marshalled: `{"jsonrpc":"1.0","method":"keypoolrefill","params":[200],"id":1}`,
			unmarshalled: &json.KeyPoolRefillCmd{
				NewSize: json.Uint(200),
			},
		},
		{
			name: "listaccounts",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("listaccounts")
			},
			staticCmd: func() interface{} {

				return json.NewListAccountsCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listaccounts","params":[],"id":1}`,
			unmarshalled: &json.ListAccountsCmd{
				MinConf: json.Int(1),
			},
		},
		{
			name: "listaccounts optional",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("listaccounts", 6)
			},
			staticCmd: func() interface{} {

				return json.NewListAccountsCmd(json.Int(6))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listaccounts","params":[6],"id":1}`,
			unmarshalled: &json.ListAccountsCmd{
				MinConf: json.Int(6),
			},
		},
		{
			name: "listaddressgroupings",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("listaddressgroupings")
			},
			staticCmd: func() interface{} {

				return json.NewListAddressGroupingsCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"listaddressgroupings","params":[],"id":1}`,
			unmarshalled: &json.ListAddressGroupingsCmd{},
		},
		{
			name: "listlockunspent",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("listlockunspent")
			},
			staticCmd: func() interface{} {

				return json.NewListLockUnspentCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"listlockunspent","params":[],"id":1}`,
			unmarshalled: &json.ListLockUnspentCmd{},
		},
		{
			name: "listreceivedbyaccount",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("listreceivedbyaccount")
			},
			staticCmd: func() interface{} {

				return json.NewListReceivedByAccountCmd(nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaccount","params":[],"id":1}`,
			unmarshalled: &json.ListReceivedByAccountCmd{
				MinConf:          json.Int(1),
				IncludeEmpty:     json.Bool(false),
				IncludeWatchOnly: json.Bool(false),
			},
		},
		{
			name: "listreceivedbyaccount optional1",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("listreceivedbyaccount", 6)
			},
			staticCmd: func() interface{} {

				return json.NewListReceivedByAccountCmd(json.Int(6), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaccount","params":[6],"id":1}`,
			unmarshalled: &json.ListReceivedByAccountCmd{
				MinConf:          json.Int(6),
				IncludeEmpty:     json.Bool(false),
				IncludeWatchOnly: json.Bool(false),
			},
		},
		{
			name: "listreceivedbyaccount optional2",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("listreceivedbyaccount", 6, true)
			},
			staticCmd: func() interface{} {

				return json.NewListReceivedByAccountCmd(json.Int(6), json.Bool(true), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaccount","params":[6,true],"id":1}`,
			unmarshalled: &json.ListReceivedByAccountCmd{
				MinConf:          json.Int(6),
				IncludeEmpty:     json.Bool(true),
				IncludeWatchOnly: json.Bool(false),
			},
		},
		{
			name: "listreceivedbyaccount optional3",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("listreceivedbyaccount", 6, true, false)
			},
			staticCmd: func() interface{} {

				return json.NewListReceivedByAccountCmd(json.Int(6), json.Bool(true), json.Bool(false))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaccount","params":[6,true,false],"id":1}`,
			unmarshalled: &json.ListReceivedByAccountCmd{
				MinConf:          json.Int(6),
				IncludeEmpty:     json.Bool(true),
				IncludeWatchOnly: json.Bool(false),
			},
		},
		{
			name: "listreceivedbyaddress",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("listreceivedbyaddress")
			},
			staticCmd: func() interface{} {

				return json.NewListReceivedByAddressCmd(nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaddress","params":[],"id":1}`,
			unmarshalled: &json.ListReceivedByAddressCmd{
				MinConf:          json.Int(1),
				IncludeEmpty:     json.Bool(false),
				IncludeWatchOnly: json.Bool(false),
			},
		},
		{
			name: "listreceivedbyaddress optional1",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("listreceivedbyaddress", 6)
			},
			staticCmd: func() interface{} {

				return json.NewListReceivedByAddressCmd(json.Int(6), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaddress","params":[6],"id":1}`,
			unmarshalled: &json.ListReceivedByAddressCmd{
				MinConf:          json.Int(6),
				IncludeEmpty:     json.Bool(false),
				IncludeWatchOnly: json.Bool(false),
			},
		},
		{
			name: "listreceivedbyaddress optional2",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("listreceivedbyaddress", 6, true)
			},
			staticCmd: func() interface{} {

				return json.NewListReceivedByAddressCmd(json.Int(6), json.Bool(true), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaddress","params":[6,true],"id":1}`,
			unmarshalled: &json.ListReceivedByAddressCmd{
				MinConf:          json.Int(6),
				IncludeEmpty:     json.Bool(true),
				IncludeWatchOnly: json.Bool(false),
			},
		},
		{
			name: "listreceivedbyaddress optional3",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("listreceivedbyaddress", 6, true, false)
			},
			staticCmd: func() interface{} {

				return json.NewListReceivedByAddressCmd(json.Int(6), json.Bool(true), json.Bool(false))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaddress","params":[6,true,false],"id":1}`,
			unmarshalled: &json.ListReceivedByAddressCmd{
				MinConf:          json.Int(6),
				IncludeEmpty:     json.Bool(true),
				IncludeWatchOnly: json.Bool(false),
			},
		},
		{
			name: "listsinceblock",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("listsinceblock")
			},
			staticCmd: func() interface{} {

				return json.NewListSinceBlockCmd(nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listsinceblock","params":[],"id":1}`,
			unmarshalled: &json.ListSinceBlockCmd{
				BlockHash:           nil,
				TargetConfirmations: json.Int(1),
				IncludeWatchOnly:    json.Bool(false),
			},
		},
		{
			name: "listsinceblock optional1",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("listsinceblock", "123")
			},
			staticCmd: func() interface{} {

				return json.NewListSinceBlockCmd(json.String("123"), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listsinceblock","params":["123"],"id":1}`,
			unmarshalled: &json.ListSinceBlockCmd{
				BlockHash:           json.String("123"),
				TargetConfirmations: json.Int(1),
				IncludeWatchOnly:    json.Bool(false),
			},
		},
		{
			name: "listsinceblock optional2",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("listsinceblock", "123", 6)
			},
			staticCmd: func() interface{} {

				return json.NewListSinceBlockCmd(json.String("123"), json.Int(6), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listsinceblock","params":["123",6],"id":1}`,
			unmarshalled: &json.ListSinceBlockCmd{
				BlockHash:           json.String("123"),
				TargetConfirmations: json.Int(6),
				IncludeWatchOnly:    json.Bool(false),
			},
		},
		{
			name: "listsinceblock optional3",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("listsinceblock", "123", 6, true)
			},
			staticCmd: func() interface{} {

				return json.NewListSinceBlockCmd(json.String("123"), json.Int(6), json.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listsinceblock","params":["123",6,true],"id":1}`,
			unmarshalled: &json.ListSinceBlockCmd{
				BlockHash:           json.String("123"),
				TargetConfirmations: json.Int(6),
				IncludeWatchOnly:    json.Bool(true),
			},
		},
		{
			name: "listtransactions",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("listtransactions")
			},
			staticCmd: func() interface{} {

				return json.NewListTransactionsCmd(nil, nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listtransactions","params":[],"id":1}`,
			unmarshalled: &json.ListTransactionsCmd{
				Account:          nil,
				Count:            json.Int(10),
				From:             json.Int(0),
				IncludeWatchOnly: json.Bool(false),
			},
		},
		{
			name: "listtransactions optional1",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("listtransactions", "acct")
			},
			staticCmd: func() interface{} {

				return json.NewListTransactionsCmd(json.String("acct"), nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listtransactions","params":["acct"],"id":1}`,
			unmarshalled: &json.ListTransactionsCmd{
				Account:          json.String("acct"),
				Count:            json.Int(10),
				From:             json.Int(0),
				IncludeWatchOnly: json.Bool(false),
			},
		},
		{
			name: "listtransactions optional2",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("listtransactions", "acct", 20)
			},
			staticCmd: func() interface{} {

				return json.NewListTransactionsCmd(json.String("acct"), json.Int(20), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listtransactions","params":["acct",20],"id":1}`,
			unmarshalled: &json.ListTransactionsCmd{
				Account:          json.String("acct"),
				Count:            json.Int(20),
				From:             json.Int(0),
				IncludeWatchOnly: json.Bool(false),
			},
		},
		{
			name: "listtransactions optional3",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("listtransactions", "acct", 20, 1)
			},
			staticCmd: func() interface{} {

				return json.NewListTransactionsCmd(json.String("acct"), json.Int(20),
					json.Int(1), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listtransactions","params":["acct",20,1],"id":1}`,
			unmarshalled: &json.ListTransactionsCmd{
				Account:          json.String("acct"),
				Count:            json.Int(20),
				From:             json.Int(1),
				IncludeWatchOnly: json.Bool(false),
			},
		},
		{
			name: "listtransactions optional4",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("listtransactions", "acct", 20, 1, true)
			},
			staticCmd: func() interface{} {

				return json.NewListTransactionsCmd(json.String("acct"), json.Int(20),
					json.Int(1), json.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listtransactions","params":["acct",20,1,true],"id":1}`,
			unmarshalled: &json.ListTransactionsCmd{
				Account:          json.String("acct"),
				Count:            json.Int(20),
				From:             json.Int(1),
				IncludeWatchOnly: json.Bool(true),
			},
		},
		{
			name: "listunspent",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("listunspent")
			},
			staticCmd: func() interface{} {

				return json.NewListUnspentCmd(nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listunspent","params":[],"id":1}`,
			unmarshalled: &json.ListUnspentCmd{
				MinConf:   json.Int(1),
				MaxConf:   json.Int(9999999),
				Addresses: nil,
			},
		},
		{
			name: "listunspent optional1",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("listunspent", 6)
			},
			staticCmd: func() interface{} {

				return json.NewListUnspentCmd(json.Int(6), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listunspent","params":[6],"id":1}`,
			unmarshalled: &json.ListUnspentCmd{
				MinConf:   json.Int(6),
				MaxConf:   json.Int(9999999),
				Addresses: nil,
			},
		},
		{
			name: "listunspent optional2",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("listunspent", 6, 100)
			},
			staticCmd: func() interface{} {

				return json.NewListUnspentCmd(json.Int(6), json.Int(100), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listunspent","params":[6,100],"id":1}`,
			unmarshalled: &json.ListUnspentCmd{
				MinConf:   json.Int(6),
				MaxConf:   json.Int(100),
				Addresses: nil,
			},
		},
		{
			name: "listunspent optional3",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("listunspent", 6, 100, []string{"1Address", "1Address2"})
			},
			staticCmd: func() interface{} {

				return json.NewListUnspentCmd(json.Int(6), json.Int(100),
					&[]string{"1Address", "1Address2"})
			},
			marshalled: `{"jsonrpc":"1.0","method":"listunspent","params":[6,100,["1Address","1Address2"]],"id":1}`,
			unmarshalled: &json.ListUnspentCmd{
				MinConf:   json.Int(6),
				MaxConf:   json.Int(100),
				Addresses: &[]string{"1Address", "1Address2"},
			},
		},
		{
			name: "lockunspent",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("lockunspent", true, `[{"txid":"123","vout":1}]`)
			},
			staticCmd: func() interface{} {

				txInputs := []json.TransactionInput{
					{Txid: "123", Vout: 1},
				}
				return json.NewLockUnspentCmd(true, txInputs)
			},
			marshalled: `{"jsonrpc":"1.0","method":"lockunspent","params":[true,[{"txid":"123","vout":1}]],"id":1}`,
			unmarshalled: &json.LockUnspentCmd{
				Unlock: true,
				Transactions: []json.TransactionInput{
					{Txid: "123", Vout: 1},
				},
			},
		},
		{
			name: "move",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("move", "from", "to", 0.5)
			},
			staticCmd: func() interface{} {

				return json.NewMoveCmd("from", "to", 0.5, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"move","params":["from","to",0.5],"id":1}`,
			unmarshalled: &json.MoveCmd{
				FromAccount: "from",
				ToAccount:   "to",
				Amount:      0.5,
				MinConf:     json.Int(1),
				Comment:     nil,
			},
		},
		{
			name: "move optional1",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("move", "from", "to", 0.5, 6)
			},
			staticCmd: func() interface{} {

				return json.NewMoveCmd("from", "to", 0.5, json.Int(6), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"move","params":["from","to",0.5,6],"id":1}`,
			unmarshalled: &json.MoveCmd{
				FromAccount: "from",
				ToAccount:   "to",
				Amount:      0.5,
				MinConf:     json.Int(6),
				Comment:     nil,
			},
		},
		{
			name: "move optional2",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("move", "from", "to", 0.5, 6, "comment")
			},
			staticCmd: func() interface{} {

				return json.NewMoveCmd("from", "to", 0.5, json.Int(6), json.String("comment"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"move","params":["from","to",0.5,6,"comment"],"id":1}`,
			unmarshalled: &json.MoveCmd{
				FromAccount: "from",
				ToAccount:   "to",
				Amount:      0.5,
				MinConf:     json.Int(6),
				Comment:     json.String("comment"),
			},
		},
		{
			name: "sendfrom",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("sendfrom", "from", "1Address", 0.5)
			},
			staticCmd: func() interface{} {

				return json.NewSendFromCmd("from", "1Address", 0.5, nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendfrom","params":["from","1Address",0.5],"id":1}`,
			unmarshalled: &json.SendFromCmd{
				FromAccount: "from",
				ToAddress:   "1Address",
				Amount:      0.5,
				MinConf:     json.Int(1),
				Comment:     nil,
				CommentTo:   nil,
			},
		},
		{
			name: "sendfrom optional1",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("sendfrom", "from", "1Address", 0.5, 6)
			},
			staticCmd: func() interface{} {

				return json.NewSendFromCmd("from", "1Address", 0.5, json.Int(6), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendfrom","params":["from","1Address",0.5,6],"id":1}`,
			unmarshalled: &json.SendFromCmd{
				FromAccount: "from",
				ToAddress:   "1Address",
				Amount:      0.5,
				MinConf:     json.Int(6),
				Comment:     nil,
				CommentTo:   nil,
			},
		},
		{
			name: "sendfrom optional2",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("sendfrom", "from", "1Address", 0.5, 6, "comment")
			},
			staticCmd: func() interface{} {

				return json.NewSendFromCmd("from", "1Address", 0.5, json.Int(6),
					json.String("comment"), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendfrom","params":["from","1Address",0.5,6,"comment"],"id":1}`,
			unmarshalled: &json.SendFromCmd{
				FromAccount: "from",
				ToAddress:   "1Address",
				Amount:      0.5,
				MinConf:     json.Int(6),
				Comment:     json.String("comment"),
				CommentTo:   nil,
			},
		},
		{
			name: "sendfrom optional3",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("sendfrom", "from", "1Address", 0.5, 6, "comment", "commentto")
			},
			staticCmd: func() interface{} {

				return json.NewSendFromCmd("from", "1Address", 0.5, json.Int(6),
					json.String("comment"), json.String("commentto"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendfrom","params":["from","1Address",0.5,6,"comment","commentto"],"id":1}`,
			unmarshalled: &json.SendFromCmd{
				FromAccount: "from",
				ToAddress:   "1Address",
				Amount:      0.5,
				MinConf:     json.Int(6),
				Comment:     json.String("comment"),
				CommentTo:   json.String("commentto"),
			},
		},
		{
			name: "sendmany",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("sendmany", "from", `{"1Address":0.5}`)
			},
			staticCmd: func() interface{} {

				amounts := map[string]float64{"1Address": 0.5}
				return json.NewSendManyCmd("from", amounts, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendmany","params":["from",{"1Address":0.5}],"id":1}`,
			unmarshalled: &json.SendManyCmd{
				FromAccount: "from",
				Amounts:     map[string]float64{"1Address": 0.5},
				MinConf:     json.Int(1),
				Comment:     nil,
			},
		},
		{
			name: "sendmany optional1",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("sendmany", "from", `{"1Address":0.5}`, 6)
			},
			staticCmd: func() interface{} {

				amounts := map[string]float64{"1Address": 0.5}
				return json.NewSendManyCmd("from", amounts, json.Int(6), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendmany","params":["from",{"1Address":0.5},6],"id":1}`,
			unmarshalled: &json.SendManyCmd{
				FromAccount: "from",
				Amounts:     map[string]float64{"1Address": 0.5},
				MinConf:     json.Int(6),
				Comment:     nil,
			},
		},
		{
			name: "sendmany optional2",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("sendmany", "from", `{"1Address":0.5}`, 6, "comment")
			},
			staticCmd: func() interface{} {

				amounts := map[string]float64{"1Address": 0.5}
				return json.NewSendManyCmd("from", amounts, json.Int(6), json.String("comment"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendmany","params":["from",{"1Address":0.5},6,"comment"],"id":1}`,
			unmarshalled: &json.SendManyCmd{
				FromAccount: "from",
				Amounts:     map[string]float64{"1Address": 0.5},
				MinConf:     json.Int(6),
				Comment:     json.String("comment"),
			},
		},
		{
			name: "sendtoaddress",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("sendtoaddress", "1Address", 0.5)
			},
			staticCmd: func() interface{} {

				return json.NewSendToAddressCmd("1Address", 0.5, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendtoaddress","params":["1Address",0.5],"id":1}`,
			unmarshalled: &json.SendToAddressCmd{
				Address:   "1Address",
				Amount:    0.5,
				Comment:   nil,
				CommentTo: nil,
			},
		},
		{
			name: "sendtoaddress optional1",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("sendtoaddress", "1Address", 0.5, "comment", "commentto")
			},
			staticCmd: func() interface{} {

				return json.NewSendToAddressCmd("1Address", 0.5, json.String("comment"),
					json.String("commentto"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendtoaddress","params":["1Address",0.5,"comment","commentto"],"id":1}`,
			unmarshalled: &json.SendToAddressCmd{
				Address:   "1Address",
				Amount:    0.5,
				Comment:   json.String("comment"),
				CommentTo: json.String("commentto"),
			},
		},
		{
			name: "setaccount",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("setaccount", "1Address", "acct")
			},
			staticCmd: func() interface{} {

				return json.NewSetAccountCmd("1Address", "acct")
			},
			marshalled: `{"jsonrpc":"1.0","method":"setaccount","params":["1Address","acct"],"id":1}`,
			unmarshalled: &json.SetAccountCmd{
				Address: "1Address",
				Account: "acct",
			},
		},
		{
			name: "settxfee",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("settxfee", 0.0001)
			},
			staticCmd: func() interface{} {

				return json.NewSetTxFeeCmd(0.0001)
			},
			marshalled: `{"jsonrpc":"1.0","method":"settxfee","params":[0.0001],"id":1}`,
			unmarshalled: &json.SetTxFeeCmd{
				Amount: 0.0001,
			},
		},
		{
			name: "signmessage",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("signmessage", "1Address", "message")
			},
			staticCmd: func() interface{} {

				return json.NewSignMessageCmd("1Address", "message")
			},
			marshalled: `{"jsonrpc":"1.0","method":"signmessage","params":["1Address","message"],"id":1}`,
			unmarshalled: &json.SignMessageCmd{
				Address: "1Address",
				Message: "message",
			},
		},
		{
			name: "signrawtransaction",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("signrawtransaction", "001122")
			},
			staticCmd: func() interface{} {

				return json.NewSignRawTransactionCmd("001122", nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"signrawtransaction","params":["001122"],"id":1}`,
			unmarshalled: &json.SignRawTransactionCmd{
				RawTx:    "001122",
				Inputs:   nil,
				PrivKeys: nil,
				Flags:    json.String("ALL"),
			},
		},
		{
			name: "signrawtransaction optional1",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("signrawtransaction", "001122", `[{"txid":"123","vout":1,"scriptPubKey":"00","redeemScript":"01"}]`)
			},
			staticCmd: func() interface{} {

				txInputs := []json.RawTxInput{
					{
						Txid:         "123",
						Vout:         1,
						ScriptPubKey: "00",
						RedeemScript: "01",
					},
				}
				return json.NewSignRawTransactionCmd("001122", &txInputs, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"signrawtransaction","params":["001122",[{"txid":"123","vout":1,"scriptPubKey":"00","redeemScript":"01"}]],"id":1}`,
			unmarshalled: &json.SignRawTransactionCmd{
				RawTx: "001122",
				Inputs: &[]json.RawTxInput{
					{
						Txid:         "123",
						Vout:         1,
						ScriptPubKey: "00",
						RedeemScript: "01",
					},
				},
				PrivKeys: nil,
				Flags:    json.String("ALL"),
			},
		},
		{
			name: "signrawtransaction optional2",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("signrawtransaction", "001122", `[]`, `["abc"]`)
			},
			staticCmd: func() interface{} {

				txInputs := []json.RawTxInput{}
				privKeys := []string{"abc"}
				return json.NewSignRawTransactionCmd("001122", &txInputs, &privKeys, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"signrawtransaction","params":["001122",[],["abc"]],"id":1}`,
			unmarshalled: &json.SignRawTransactionCmd{
				RawTx:    "001122",
				Inputs:   &[]json.RawTxInput{},
				PrivKeys: &[]string{"abc"},
				Flags:    json.String("ALL"),
			},
		},
		{
			name: "signrawtransaction optional3",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("signrawtransaction", "001122", `[]`, `[]`, "ALL")
			},
			staticCmd: func() interface{} {

				txInputs := []json.RawTxInput{}
				privKeys := []string{}
				return json.NewSignRawTransactionCmd("001122", &txInputs, &privKeys,
					json.String("ALL"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"signrawtransaction","params":["001122",[],[],"ALL"],"id":1}`,
			unmarshalled: &json.SignRawTransactionCmd{
				RawTx:    "001122",
				Inputs:   &[]json.RawTxInput{},
				PrivKeys: &[]string{},
				Flags:    json.String("ALL"),
			},
		},
		{
			name: "walletlock",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("walletlock")
			},
			staticCmd: func() interface{} {

				return json.NewWalletLockCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"walletlock","params":[],"id":1}`,
			unmarshalled: &json.WalletLockCmd{},
		},
		{
			name: "walletpassphrase",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("walletpassphrase", "pass", 60)
			},
			staticCmd: func() interface{} {

				return json.NewWalletPassphraseCmd("pass", 60)
			},
			marshalled: `{"jsonrpc":"1.0","method":"walletpassphrase","params":["pass",60],"id":1}`,
			unmarshalled: &json.WalletPassphraseCmd{
				Passphrase: "pass",
				Timeout:    60,
			},
		},
		{
			name: "walletpassphrasechange",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("walletpassphrasechange", "old", "new")
			},
			staticCmd: func() interface{} {

				return json.NewWalletPassphraseChangeCmd("old", "new")
			},
			marshalled: `{"jsonrpc":"1.0","method":"walletpassphrasechange","params":["old","new"],"id":1}`,
			unmarshalled: &json.WalletPassphraseChangeCmd{
				OldPassphrase: "old",
				NewPassphrase: "new",
			},
		},
	}
	t.Logf("Running %d tests", len(tests))

	for i, test := range tests {

		// Marshal the command as created by the new static command creation function.
		marshalled, err := json.MarshalCmd(testID, test.staticCmd())

		if err != nil {

			t.Errorf("MarshalCmd #%d (%s) unexpected error: %v", i, test.name, err)
			continue
		}

		if !bytes.Equal(marshalled, []byte(test.marshalled)) {

			t.Errorf("Test #%d (%s) unexpected marshalled data - "+
				"got %s, want %s", i, test.name, marshalled,
				test.marshalled)
			continue
		}
		// Ensure the command is created without error via the generic new command creation function.
		cmd, err := test.newCmd()

		if err != nil {

			t.Errorf("Test #%d (%s) unexpected NewCmd error: %v ",
				i, test.name, err)
		}
		// Marshal the command as created by the generic new command creation function.
		marshalled, err = json.MarshalCmd(testID, cmd)

		if err != nil {

			t.Errorf("MarshalCmd #%d (%s) unexpected error: %v", i,
				test.name, err)
			continue
		}

		if !bytes.Equal(marshalled, []byte(test.marshalled)) {

			t.Errorf("Test #%d (%s) unexpected marshalled data - "+
				"got %s, want %s", i, test.name, marshalled,
				test.marshalled)
			continue
		}
		var request json.Request

		if err := json.Unmarshal(marshalled, &request); err != nil {

			t.Errorf("Test #%d (%s) unexpected error while "+
				"unmarshalling JSON-RPC request: %v", i,
				test.name, err)
			continue
		}
		cmd, err = json.UnmarshalCmd(&request)

		if err != nil {

			t.Errorf("UnmarshalCmd #%d (%s) unexpected error: %v", i,
				test.name, err)
			continue
		}

		if !reflect.DeepEqual(cmd, test.unmarshalled) {

			t.Errorf("Test #%d (%s) unexpected unmarshalled command "+
				"- got %s, want %s", i, test.name,
				fmt.Sprintf("(%T) %+[1]v", cmd),
				fmt.Sprintf("(%T) %+[1]v\n", test.unmarshalled))
			continue
		}
	}
}
