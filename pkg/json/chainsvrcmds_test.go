package json_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"git.parallelcoin.io/pod/pkg/json"
	"git.parallelcoin.io/pod/pkg/wire"
)

// TestChainSvrCmds tests all of the chain server commands marshal and unmarshal into valid results include handling of optional fields being omitted in the marshalled command, while optional fields with defaults have the default assigned on unmarshalled commands.
func TestChainSvrCmds(
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
			name: "addnode",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("addnode", "127.0.0.1", json.ANRemove)
			},
			staticCmd: func() interface{} {
				return json.NewAddNodeCmd("127.0.0.1", json.ANRemove)
			},
			marshalled:   `{"jsonrpc":"1.0","method":"addnode","params":["127.0.0.1","remove"],"id":1}`,
			unmarshalled: &json.AddNodeCmd{Addr: "127.0.0.1", SubCmd: json.ANRemove},
		},
		{
			name: "createrawtransaction",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("createrawtransaction", `[{"txid":"123","vout":1}]`,
					`{"456":0.0123}`)
			},
			staticCmd: func() interface{} {
				txInputs := []json.TransactionInput{
					{Txid: "123", Vout: 1},
				}
				amounts := map[string]float64{"456": .0123}
				return json.NewCreateRawTransactionCmd(txInputs, amounts, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"createrawtransaction","params":[[{"txid":"123","vout":1}],{"456":0.0123}],"id":1}`,
			unmarshalled: &json.CreateRawTransactionCmd{
				Inputs:  []json.TransactionInput{{Txid: "123", Vout: 1}},
				Amounts: map[string]float64{"456": .0123},
			},
		},
		{
			name: "createrawtransaction optional",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("createrawtransaction", `[{"txid":"123","vout":1}]`,
					`{"456":0.0123}`, int64(12312333333))
			},
			staticCmd: func() interface{} {
				txInputs := []json.TransactionInput{
					{Txid: "123", Vout: 1},
				}
				amounts := map[string]float64{"456": .0123}
				return json.NewCreateRawTransactionCmd(txInputs, amounts, json.Int64(12312333333))
			},
			marshalled: `{"jsonrpc":"1.0","method":"createrawtransaction","params":[[{"txid":"123","vout":1}],{"456":0.0123},12312333333],"id":1}`,
			unmarshalled: &json.CreateRawTransactionCmd{
				Inputs:   []json.TransactionInput{{Txid: "123", Vout: 1}},
				Amounts:  map[string]float64{"456": .0123},
				LockTime: json.Int64(12312333333),
			},
		},
		{
			name: "decoderawtransaction",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("decoderawtransaction", "123")
			},
			staticCmd: func() interface{} {
				return json.NewDecodeRawTransactionCmd("123")
			},
			marshalled:   `{"jsonrpc":"1.0","method":"decoderawtransaction","params":["123"],"id":1}`,
			unmarshalled: &json.DecodeRawTransactionCmd{HexTx: "123"},
		},
		{
			name: "decodescript",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("decodescript", "00")
			},
			staticCmd: func() interface{} {
				return json.NewDecodeScriptCmd("00")
			},
			marshalled:   `{"jsonrpc":"1.0","method":"decodescript","params":["00"],"id":1}`,
			unmarshalled: &json.DecodeScriptCmd{HexScript: "00"},
		},
		{
			name: "getaddednodeinfo",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("getaddednodeinfo", true)
			},
			staticCmd: func() interface{} {
				return json.NewGetAddedNodeInfoCmd(true, nil)
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getaddednodeinfo","params":[true],"id":1}`,
			unmarshalled: &json.GetAddedNodeInfoCmd{DNS: true, Node: nil},
		},
		{
			name: "getaddednodeinfo optional",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("getaddednodeinfo", true, "127.0.0.1")
			},
			staticCmd: func() interface{} {
				return json.NewGetAddedNodeInfoCmd(true, json.String("127.0.0.1"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getaddednodeinfo","params":[true,"127.0.0.1"],"id":1}`,
			unmarshalled: &json.GetAddedNodeInfoCmd{
				DNS:  true,
				Node: json.String("127.0.0.1"),
			},
		},
		{
			name: "getbestblockhash",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("getbestblockhash")
			},
			staticCmd: func() interface{} {
				return json.NewGetBestBlockHashCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getbestblockhash","params":[],"id":1}`,
			unmarshalled: &json.GetBestBlockHashCmd{},
		},
		{
			name: "getblock",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("getblock", "123")
			},
			staticCmd: func() interface{} {
				return json.NewGetBlockCmd("123", nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblock","params":["123"],"id":1}`,
			unmarshalled: &json.GetBlockCmd{
				Hash:      "123",
				Verbose:   json.Bool(true),
				VerboseTx: json.Bool(false),
			},
		},
		{
			name: "getblock required optional1",
			newCmd: func() (interface{}, error) {

				// Intentionally use a source param that is more pointers than the destination to exercise that path.
				verbosePtr := json.Bool(true)
				return json.NewCmd("getblock", "123", &verbosePtr)
			},
			staticCmd: func() interface{} {
				return json.NewGetBlockCmd("123", json.Bool(true), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblock","params":["123",true],"id":1}`,
			unmarshalled: &json.GetBlockCmd{
				Hash:      "123",
				Verbose:   json.Bool(true),
				VerboseTx: json.Bool(false),
			},
		},
		{
			name: "getblock required optional2",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("getblock", "123", true, true)
			},
			staticCmd: func() interface{} {
				return json.NewGetBlockCmd("123", json.Bool(true), json.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblock","params":["123",true,true],"id":1}`,
			unmarshalled: &json.GetBlockCmd{
				Hash:      "123",
				Verbose:   json.Bool(true),
				VerboseTx: json.Bool(true),
			},
		},
		{
			name: "getblockchaininfo",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("getblockchaininfo")
			},
			staticCmd: func() interface{} {
				return json.NewGetBlockChainInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getblockchaininfo","params":[],"id":1}`,
			unmarshalled: &json.GetBlockChainInfoCmd{},
		},
		{
			name: "getblockcount",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("getblockcount")
			},
			staticCmd: func() interface{} {
				return json.NewGetBlockCountCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getblockcount","params":[],"id":1}`,
			unmarshalled: &json.GetBlockCountCmd{},
		},
		{
			name: "getblockhash",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("getblockhash", 123)
			},
			staticCmd: func() interface{} {
				return json.NewGetBlockHashCmd(123)
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getblockhash","params":[123],"id":1}`,
			unmarshalled: &json.GetBlockHashCmd{Index: 123},
		},
		{
			name: "getblockheader",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("getblockheader", "123")
			},
			staticCmd: func() interface{} {
				return json.NewGetBlockHeaderCmd("123", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblockheader","params":["123"],"id":1}`,
			unmarshalled: &json.GetBlockHeaderCmd{
				Hash:    "123",
				Verbose: json.Bool(true),
			},
		},
		{
			name: "getblocktemplate",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("getblocktemplate")
			},
			staticCmd: func() interface{} {
				return json.NewGetBlockTemplateCmd(nil)
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getblocktemplate","params":[],"id":1}`,
			unmarshalled: &json.GetBlockTemplateCmd{Request: nil},
		},
		{
			name: "getblocktemplate optional - template request",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("getblocktemplate", `{"mode":"template","capabilities":["longpoll","coinbasetxn"]}`)
			},
			staticCmd: func() interface{} {
				template := json.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
				}
				return json.NewGetBlockTemplateCmd(&template)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblocktemplate","params":[{"mode":"template","capabilities":["longpoll","coinbasetxn"]}],"id":1}`,
			unmarshalled: &json.GetBlockTemplateCmd{
				Request: &json.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
				},
			},
		},
		{
			name: "getblocktemplate optional - template request with tweaks",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("getblocktemplate", `{"mode":"template","capabilities":["longpoll","coinbasetxn"],"sigoplimit":500,"sizelimit":100000000,"maxversion":2}`)
			},
			staticCmd: func() interface{} {
				template := json.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
					SigOpLimit:   500,
					SizeLimit:    100000000,
					MaxVersion:   2,
				}
				return json.NewGetBlockTemplateCmd(&template)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblocktemplate","params":[{"mode":"template","capabilities":["longpoll","coinbasetxn"],"sigoplimit":500,"sizelimit":100000000,"maxversion":2}],"id":1}`,
			unmarshalled: &json.GetBlockTemplateCmd{
				Request: &json.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
					SigOpLimit:   int64(500),
					SizeLimit:    int64(100000000),
					MaxVersion:   2,
				},
			},
		},
		{
			name: "getblocktemplate optional - template request with tweaks 2",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("getblocktemplate", `{"mode":"template","capabilities":["longpoll","coinbasetxn"],"sigoplimit":true,"sizelimit":100000000,"maxversion":2}`)
			},
			staticCmd: func() interface{} {
				template := json.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
					SigOpLimit:   true,
					SizeLimit:    100000000,
					MaxVersion:   2,
				}
				return json.NewGetBlockTemplateCmd(&template)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblocktemplate","params":[{"mode":"template","capabilities":["longpoll","coinbasetxn"],"sigoplimit":true,"sizelimit":100000000,"maxversion":2}],"id":1}`,
			unmarshalled: &json.GetBlockTemplateCmd{
				Request: &json.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
					SigOpLimit:   true,
					SizeLimit:    int64(100000000),
					MaxVersion:   2,
				},
			},
		},
		{
			name: "getcfilter",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("getcfilter", "123",
					wire.GCSFilterRegular)
			},
			staticCmd: func() interface{} {
				return json.NewGetCFilterCmd("123",
					wire.GCSFilterRegular)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getcfilter","params":["123",0],"id":1}`,
			unmarshalled: &json.GetCFilterCmd{
				Hash:       "123",
				FilterType: wire.GCSFilterRegular,
			},
		},
		{
			name: "getcfilterheader",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("getcfilterheader", "123",
					wire.GCSFilterRegular)
			},
			staticCmd: func() interface{} {
				return json.NewGetCFilterHeaderCmd("123",
					wire.GCSFilterRegular)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getcfilterheader","params":["123",0],"id":1}`,
			unmarshalled: &json.GetCFilterHeaderCmd{
				Hash:       "123",
				FilterType: wire.GCSFilterRegular,
			},
		},
		{
			name: "getchaintips",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("getchaintips")
			},
			staticCmd: func() interface{} {
				return json.NewGetChainTipsCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getchaintips","params":[],"id":1}`,
			unmarshalled: &json.GetChainTipsCmd{},
		},
		{
			name: "getconnectioncount",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("getconnectioncount")
			},
			staticCmd: func() interface{} {
				return json.NewGetConnectionCountCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getconnectioncount","params":[],"id":1}`,
			unmarshalled: &json.GetConnectionCountCmd{},
		},
		{
			name: "getdifficulty",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("getdifficulty", "123")
			},
			staticCmd: func() interface{} {
				return json.NewGetDifficultyCmd("123")
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getdifficulty","params":["123"],"id":1}`,
			unmarshalled: &json.GetDifficultyCmd{Algo: "123"},
		},
		{
			name: "getgenerate",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("getgenerate")
			},
			staticCmd: func() interface{} {
				return json.NewGetGenerateCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getgenerate","params":[],"id":1}`,
			unmarshalled: &json.GetGenerateCmd{},
		},
		{
			name: "gethashespersec",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("gethashespersec")
			},
			staticCmd: func() interface{} {
				return json.NewGetHashesPerSecCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"gethashespersec","params":[],"id":1}`,
			unmarshalled: &json.GetHashesPerSecCmd{},
		},
		{
			name: "getinfo",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("getinfo")
			},
			staticCmd: func() interface{} {
				return json.NewGetInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getinfo","params":[],"id":1}`,
			unmarshalled: &json.GetInfoCmd{},
		},
		{
			name: "getmempoolentry",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("getmempoolentry", "txhash")
			},
			staticCmd: func() interface{} {
				return json.NewGetMempoolEntryCmd("txhash")
			},
			marshalled: `{"jsonrpc":"1.0","method":"getmempoolentry","params":["txhash"],"id":1}`,
			unmarshalled: &json.GetMempoolEntryCmd{
				TxID: "txhash",
			},
		},
		{
			name: "getmempoolinfo",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("getmempoolinfo")
			},
			staticCmd: func() interface{} {
				return json.NewGetMempoolInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getmempoolinfo","params":[],"id":1}`,
			unmarshalled: &json.GetMempoolInfoCmd{},
		},
		{
			name: "getmininginfo",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("getmininginfo")
			},
			staticCmd: func() interface{} {
				return json.NewGetMiningInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getmininginfo","params":[],"id":1}`,
			unmarshalled: &json.GetMiningInfoCmd{},
		},
		{
			name: "getnetworkinfo",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("getnetworkinfo")
			},
			staticCmd: func() interface{} {
				return json.NewGetNetworkInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getnetworkinfo","params":[],"id":1}`,
			unmarshalled: &json.GetNetworkInfoCmd{},
		},
		{
			name: "getnettotals",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("getnettotals")
			},
			staticCmd: func() interface{} {
				return json.NewGetNetTotalsCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getnettotals","params":[],"id":1}`,
			unmarshalled: &json.GetNetTotalsCmd{},
		},
		{
			name: "getnetworkhashps",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("getnetworkhashps")
			},
			staticCmd: func() interface{} {
				return json.NewGetNetworkHashPSCmd(nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getnetworkhashps","params":[],"id":1}`,
			unmarshalled: &json.GetNetworkHashPSCmd{
				Blocks: json.Int(120),
				Height: json.Int(-1),
			},
		},
		{
			name: "getnetworkhashps optional1",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("getnetworkhashps", 200)
			},
			staticCmd: func() interface{} {
				return json.NewGetNetworkHashPSCmd(json.Int(200), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getnetworkhashps","params":[200],"id":1}`,
			unmarshalled: &json.GetNetworkHashPSCmd{
				Blocks: json.Int(200),
				Height: json.Int(-1),
			},
		},
		{
			name: "getnetworkhashps optional2",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("getnetworkhashps", 200, 123)
			},
			staticCmd: func() interface{} {
				return json.NewGetNetworkHashPSCmd(json.Int(200), json.Int(123))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getnetworkhashps","params":[200,123],"id":1}`,
			unmarshalled: &json.GetNetworkHashPSCmd{
				Blocks: json.Int(200),
				Height: json.Int(123),
			},
		},
		{
			name: "getpeerinfo",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("getpeerinfo")
			},
			staticCmd: func() interface{} {
				return json.NewGetPeerInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getpeerinfo","params":[],"id":1}`,
			unmarshalled: &json.GetPeerInfoCmd{},
		},
		{
			name: "getrawmempool",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("getrawmempool")
			},
			staticCmd: func() interface{} {
				return json.NewGetRawMempoolCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawmempool","params":[],"id":1}`,
			unmarshalled: &json.GetRawMempoolCmd{
				Verbose: json.Bool(false),
			},
		},
		{
			name: "getrawmempool optional",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("getrawmempool", false)
			},
			staticCmd: func() interface{} {
				return json.NewGetRawMempoolCmd(json.Bool(false))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawmempool","params":[false],"id":1}`,
			unmarshalled: &json.GetRawMempoolCmd{
				Verbose: json.Bool(false),
			},
		},
		{
			name: "getrawtransaction",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("getrawtransaction", "123")
			},
			staticCmd: func() interface{} {
				return json.NewGetRawTransactionCmd("123", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawtransaction","params":["123"],"id":1}`,
			unmarshalled: &json.GetRawTransactionCmd{
				Txid:    "123",
				Verbose: json.Int(0),
			},
		},
		{
			name: "getrawtransaction optional",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("getrawtransaction", "123", 1)
			},
			staticCmd: func() interface{} {
				return json.NewGetRawTransactionCmd("123", json.Int(1))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawtransaction","params":["123",1],"id":1}`,
			unmarshalled: &json.GetRawTransactionCmd{
				Txid:    "123",
				Verbose: json.Int(1),
			},
		},
		{
			name: "gettxout",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("gettxout", "123", 1)
			},
			staticCmd: func() interface{} {
				return json.NewGetTxOutCmd("123", 1, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"gettxout","params":["123",1],"id":1}`,
			unmarshalled: &json.GetTxOutCmd{
				Txid:           "123",
				Vout:           1,
				IncludeMempool: json.Bool(true),
			},
		},
		{
			name: "gettxout optional",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("gettxout", "123", 1, true)
			},
			staticCmd: func() interface{} {
				return json.NewGetTxOutCmd("123", 1, json.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"gettxout","params":["123",1,true],"id":1}`,
			unmarshalled: &json.GetTxOutCmd{
				Txid:           "123",
				Vout:           1,
				IncludeMempool: json.Bool(true),
			},
		},
		{
			name: "gettxoutproof",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("gettxoutproof", []string{"123", "456"})
			},
			staticCmd: func() interface{} {
				return json.NewGetTxOutProofCmd([]string{"123", "456"}, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"gettxoutproof","params":[["123","456"]],"id":1}`,
			unmarshalled: &json.GetTxOutProofCmd{
				TxIDs: []string{"123", "456"},
			},
		},
		{
			name: "gettxoutproof optional",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("gettxoutproof", []string{"123", "456"},
					json.String("000000000000034a7dedef4a161fa058a2d67a173a90155f3a2fe6fc132e0ebf"))
			},
			staticCmd: func() interface{} {
				return json.NewGetTxOutProofCmd([]string{"123", "456"},
					json.String("000000000000034a7dedef4a161fa058a2d67a173a90155f3a2fe6fc132e0ebf"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"gettxoutproof","params":[["123","456"],` +
				`"000000000000034a7dedef4a161fa058a2d67a173a90155f3a2fe6fc132e0ebf"],"id":1}`,
			unmarshalled: &json.GetTxOutProofCmd{
				TxIDs:     []string{"123", "456"},
				BlockHash: json.String("000000000000034a7dedef4a161fa058a2d67a173a90155f3a2fe6fc132e0ebf"),
			},
		},
		{
			name: "gettxoutsetinfo",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("gettxoutsetinfo")
			},
			staticCmd: func() interface{} {
				return json.NewGetTxOutSetInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"gettxoutsetinfo","params":[],"id":1}`,
			unmarshalled: &json.GetTxOutSetInfoCmd{},
		},
		{
			name: "getwork",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("getwork")
			},
			staticCmd: func() interface{} {
				return json.NewGetWorkCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getwork","params":[],"id":1}`,
			unmarshalled: &json.GetWorkCmd{
				Data: nil,
			},
		},
		{
			name: "getwork optional",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("getwork", "00112233")
			},
			staticCmd: func() interface{} {
				return json.NewGetWorkCmd(json.String("00112233"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getwork","params":["00112233"],"id":1}`,
			unmarshalled: &json.GetWorkCmd{
				Data: json.String("00112233"),
			},
		},
		{
			name: "help",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("help")
			},
			staticCmd: func() interface{} {
				return json.NewHelpCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"help","params":[],"id":1}`,
			unmarshalled: &json.HelpCmd{
				Command: nil,
			},
		},
		{
			name: "help optional",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("help", "getblock")
			},
			staticCmd: func() interface{} {
				return json.NewHelpCmd(json.String("getblock"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"help","params":["getblock"],"id":1}`,
			unmarshalled: &json.HelpCmd{
				Command: json.String("getblock"),
			},
		},
		{
			name: "invalidateblock",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("invalidateblock", "123")
			},
			staticCmd: func() interface{} {
				return json.NewInvalidateBlockCmd("123")
			},
			marshalled: `{"jsonrpc":"1.0","method":"invalidateblock","params":["123"],"id":1}`,
			unmarshalled: &json.InvalidateBlockCmd{
				BlockHash: "123",
			},
		},
		{
			name: "ping",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("ping")
			},
			staticCmd: func() interface{} {
				return json.NewPingCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"ping","params":[],"id":1}`,
			unmarshalled: &json.PingCmd{},
		},
		{
			name: "preciousblock",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("preciousblock", "0123")
			},
			staticCmd: func() interface{} {
				return json.NewPreciousBlockCmd("0123")
			},
			marshalled: `{"jsonrpc":"1.0","method":"preciousblock","params":["0123"],"id":1}`,
			unmarshalled: &json.PreciousBlockCmd{
				BlockHash: "0123",
			},
		},
		{
			name: "reconsiderblock",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("reconsiderblock", "123")
			},
			staticCmd: func() interface{} {
				return json.NewReconsiderBlockCmd("123")
			},
			marshalled: `{"jsonrpc":"1.0","method":"reconsiderblock","params":["123"],"id":1}`,
			unmarshalled: &json.ReconsiderBlockCmd{
				BlockHash: "123",
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("searchrawtransactions", "1Address")
			},
			staticCmd: func() interface{} {
				return json.NewSearchRawTransactionsCmd("1Address", nil, nil, nil, nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address"],"id":1}`,
			unmarshalled: &json.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     json.Int(1),
				Skip:        json.Int(0),
				Count:       json.Int(100),
				VinExtra:    json.Int(0),
				Reverse:     json.Bool(false),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("searchrawtransactions", "1Address", 0)
			},
			staticCmd: func() interface{} {
				return json.NewSearchRawTransactionsCmd("1Address",
					json.Int(0), nil, nil, nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0],"id":1}`,
			unmarshalled: &json.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     json.Int(0),
				Skip:        json.Int(0),
				Count:       json.Int(100),
				VinExtra:    json.Int(0),
				Reverse:     json.Bool(false),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("searchrawtransactions", "1Address", 0, 5)
			},
			staticCmd: func() interface{} {
				return json.NewSearchRawTransactionsCmd("1Address",
					json.Int(0), json.Int(5), nil, nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0,5],"id":1}`,
			unmarshalled: &json.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     json.Int(0),
				Skip:        json.Int(5),
				Count:       json.Int(100),
				VinExtra:    json.Int(0),
				Reverse:     json.Bool(false),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("searchrawtransactions", "1Address", 0, 5, 10)
			},
			staticCmd: func() interface{} {
				return json.NewSearchRawTransactionsCmd("1Address",
					json.Int(0), json.Int(5), json.Int(10), nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0,5,10],"id":1}`,
			unmarshalled: &json.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     json.Int(0),
				Skip:        json.Int(5),
				Count:       json.Int(10),
				VinExtra:    json.Int(0),
				Reverse:     json.Bool(false),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("searchrawtransactions", "1Address", 0, 5, 10, 1)
			},
			staticCmd: func() interface{} {
				return json.NewSearchRawTransactionsCmd("1Address",
					json.Int(0), json.Int(5), json.Int(10), json.Int(1), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0,5,10,1],"id":1}`,
			unmarshalled: &json.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     json.Int(0),
				Skip:        json.Int(5),
				Count:       json.Int(10),
				VinExtra:    json.Int(1),
				Reverse:     json.Bool(false),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("searchrawtransactions", "1Address", 0, 5, 10, 1, true)
			},
			staticCmd: func() interface{} {
				return json.NewSearchRawTransactionsCmd("1Address",
					json.Int(0), json.Int(5), json.Int(10), json.Int(1), json.Bool(true), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0,5,10,1,true],"id":1}`,
			unmarshalled: &json.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     json.Int(0),
				Skip:        json.Int(5),
				Count:       json.Int(10),
				VinExtra:    json.Int(1),
				Reverse:     json.Bool(true),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("searchrawtransactions", "1Address", 0, 5, 10, 1, true, []string{"1Address"})
			},
			staticCmd: func() interface{} {
				return json.NewSearchRawTransactionsCmd("1Address",
					json.Int(0), json.Int(5), json.Int(10), json.Int(1), json.Bool(true), &[]string{"1Address"})
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0,5,10,1,true,["1Address"]],"id":1}`,
			unmarshalled: &json.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     json.Int(0),
				Skip:        json.Int(5),
				Count:       json.Int(10),
				VinExtra:    json.Int(1),
				Reverse:     json.Bool(true),
				FilterAddrs: &[]string{"1Address"},
			},
		},
		{
			name: "sendrawtransaction",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("sendrawtransaction", "1122")
			},
			staticCmd: func() interface{} {
				return json.NewSendRawTransactionCmd("1122", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendrawtransaction","params":["1122"],"id":1}`,
			unmarshalled: &json.SendRawTransactionCmd{
				HexTx:         "1122",
				AllowHighFees: json.Bool(false),
			},
		},
		{
			name: "sendrawtransaction optional",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("sendrawtransaction", "1122", false)
			},
			staticCmd: func() interface{} {
				return json.NewSendRawTransactionCmd("1122", json.Bool(false))
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendrawtransaction","params":["1122",false],"id":1}`,
			unmarshalled: &json.SendRawTransactionCmd{
				HexTx:         "1122",
				AllowHighFees: json.Bool(false),
			},
		},
		{
			name: "setgenerate",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("setgenerate", true)
			},
			staticCmd: func() interface{} {
				return json.NewSetGenerateCmd(true, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"setgenerate","params":[true],"id":1}`,
			unmarshalled: &json.SetGenerateCmd{
				Generate:     true,
				GenProcLimit: json.Int(-1),
			},
		},
		{
			name: "setgenerate optional",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("setgenerate", true, 6)
			},
			staticCmd: func() interface{} {
				return json.NewSetGenerateCmd(true, json.Int(6))
			},
			marshalled: `{"jsonrpc":"1.0","method":"setgenerate","params":[true,6],"id":1}`,
			unmarshalled: &json.SetGenerateCmd{
				Generate:     true,
				GenProcLimit: json.Int(6),
			},
		},
		{
			name: "stop",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("stop")
			},
			staticCmd: func() interface{} {
				return json.NewStopCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"stop","params":[],"id":1}`,
			unmarshalled: &json.StopCmd{},
		},
		{
			name: "submitblock",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("submitblock", "112233")
			},
			staticCmd: func() interface{} {
				return json.NewSubmitBlockCmd("112233", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"submitblock","params":["112233"],"id":1}`,
			unmarshalled: &json.SubmitBlockCmd{
				HexBlock: "112233",
				Options:  nil,
			},
		},
		{
			name: "submitblock optional",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("submitblock", "112233", `{"workid":"12345"}`)
			},
			staticCmd: func() interface{} {
				options := json.SubmitBlockOptions{
					WorkID: "12345",
				}
				return json.NewSubmitBlockCmd("112233", &options)
			},
			marshalled: `{"jsonrpc":"1.0","method":"submitblock","params":["112233",{"workid":"12345"}],"id":1}`,
			unmarshalled: &json.SubmitBlockCmd{
				HexBlock: "112233",
				Options: &json.SubmitBlockOptions{
					WorkID: "12345",
				},
			},
		},
		{
			name: "uptime",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("uptime")
			},
			staticCmd: func() interface{} {
				return json.NewUptimeCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"uptime","params":[],"id":1}`,
			unmarshalled: &json.UptimeCmd{},
		},
		{
			name: "validateaddress",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("validateaddress", "1Address")
			},
			staticCmd: func() interface{} {
				return json.NewValidateAddressCmd("1Address")
			},
			marshalled: `{"jsonrpc":"1.0","method":"validateaddress","params":["1Address"],"id":1}`,
			unmarshalled: &json.ValidateAddressCmd{
				Address: "1Address",
			},
		},
		{
			name: "verifychain",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("verifychain")
			},
			staticCmd: func() interface{} {
				return json.NewVerifyChainCmd(nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"verifychain","params":[],"id":1}`,
			unmarshalled: &json.VerifyChainCmd{
				CheckLevel: json.Int32(3),
				CheckDepth: json.Int32(288),
			},
		},
		{
			name: "verifychain optional1",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("verifychain", 2)
			},
			staticCmd: func() interface{} {
				return json.NewVerifyChainCmd(json.Int32(2), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"verifychain","params":[2],"id":1}`,
			unmarshalled: &json.VerifyChainCmd{
				CheckLevel: json.Int32(2),
				CheckDepth: json.Int32(288),
			},
		},
		{
			name: "verifychain optional2",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("verifychain", 2, 500)
			},
			staticCmd: func() interface{} {
				return json.NewVerifyChainCmd(json.Int32(2), json.Int32(500))
			},
			marshalled: `{"jsonrpc":"1.0","method":"verifychain","params":[2,500],"id":1}`,
			unmarshalled: &json.VerifyChainCmd{
				CheckLevel: json.Int32(2),
				CheckDepth: json.Int32(500),
			},
		},
		{
			name: "verifymessage",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("verifymessage", "1Address", "301234", "test")
			},
			staticCmd: func() interface{} {
				return json.NewVerifyMessageCmd("1Address", "301234", "test")
			},
			marshalled: `{"jsonrpc":"1.0","method":"verifymessage","params":["1Address","301234","test"],"id":1}`,
			unmarshalled: &json.VerifyMessageCmd{
				Address:   "1Address",
				Signature: "301234",
				Message:   "test",
			},
		},
		{
			name: "verifytxoutproof",
			newCmd: func() (interface{}, error) {

				return json.NewCmd("verifytxoutproof", "test")
			},
			staticCmd: func() interface{} {
				return json.NewVerifyTxOutProofCmd("test")
			},
			marshalled: `{"jsonrpc":"1.0","method":"verifytxoutproof","params":["test"],"id":1}`,
			unmarshalled: &json.VerifyTxOutProofCmd{
				Proof: "test",
			},
		},
	}
	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Marshal the command as created by the new static command
		// creation function.
		marshalled, err := json.MarshalCmd(testID, test.staticCmd())
		if err != nil {
			t.Errorf("MarshalCmd #%d (%s) unexpected error: %v", i,
				test.name, err)
			continue
		}
		if !bytes.Equal(marshalled, []byte(test.marshalled)) {

			t.Errorf("Test #%d (%s) unexpected marshalled data - "+
				"got %s, want %s", i, test.name, marshalled,
				test.marshalled)
			t.Errorf("\n%s\n%s", marshalled, test.marshalled)
			continue
		}
		// Ensure the command is created without error via the generic
		// new command creation function.
		cmd, err := test.newCmd()
		if err != nil {
			t.Errorf("Test #%d (%s) unexpected NewCmd error: %v ",
				i, test.name, err)
		}
		// Marshal the command as created by the generic new command
		// creation function.
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

// TestChainSvrCmdErrors ensures any errors that occur in the command during custom mashal and unmarshal are as expected.
func TestChainSvrCmdErrors(
	t *testing.T) {

	t.Parallel()
	tests := []struct {
		name       string
		result     interface{}
		marshalled string
		err        error
	}{
		{
			name:       "template request with invalid type",
			result:     &json.TemplateRequest{},
			marshalled: `{"mode":1}`,
			err:        &json.UnmarshalTypeError{},
		},
		{
			name:       "invalid template request sigoplimit field",
			result:     &json.TemplateRequest{},
			marshalled: `{"sigoplimit":"invalid"}`,
			err:        json.Error{ErrorCode: json.ErrInvalidType},
		},
		{
			name:       "invalid template request sizelimit field",
			result:     &json.TemplateRequest{},
			marshalled: `{"sizelimit":"invalid"}`,
			err:        json.Error{ErrorCode: json.ErrInvalidType},
		},
	}
	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		err := json.Unmarshal([]byte(test.marshalled), &test.result)
		if reflect.TypeOf(err) != reflect.TypeOf(test.err) {

			t.Errorf("Test #%d (%s) wrong error - got %T (%v), "+
				"want %T", i, test.name, err, err, test.err)
			continue
		}
		if terr, ok := test.err.(json.Error); ok {
			gotErrorCode := err.(json.Error).ErrorCode
			if gotErrorCode != terr.ErrorCode {
				t.Errorf("Test #%d (%s) mismatched error code "+
					"- got %v (%v), want %v", i, test.name,
					gotErrorCode, terr, terr.ErrorCode)
				continue
			}
		}
	}
}
