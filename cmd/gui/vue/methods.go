package vue

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"

	chaincfg "git.parallelcoin.io/dev/pod/pkg/chain/config"
	chainhash "git.parallelcoin.io/dev/pod/pkg/chain/hash"
	wtxmgr "git.parallelcoin.io/dev/pod/pkg/chain/tx/mgr"
	txrules "git.parallelcoin.io/dev/pod/pkg/chain/tx/rules"
	txscript "git.parallelcoin.io/dev/pod/pkg/chain/tx/script"
	"git.parallelcoin.io/dev/pod/pkg/chain/wire"
	rpcclient "git.parallelcoin.io/dev/pod/pkg/rpc/client"
	"git.parallelcoin.io/dev/pod/pkg/rpc/json"
	"git.parallelcoin.io/dev/pod/pkg/util"
	ec "git.parallelcoin.io/dev/pod/pkg/util/elliptic"
	"git.parallelcoin.io/dev/pod/pkg/wallet"
	waddrmgr "git.parallelcoin.io/dev/pod/pkg/wallet/addrmgr"
	chain "git.parallelcoin.io/dev/pod/pkg/wallet/chain"
	"github.com/google/martian/log"
)

// Error types to simplify the reporting of specific categories of
// errors, and their *json.RPCError creation.

type (

	// DeserializationError describes a failed deserializaion due to bad
	// user input.  It corresponds to json.ErrRPCDeserialization.

	DeserializationError struct {
		error
	}

	// InvalidParameterError describes an invalid parameter passed by
	// the user.  It corresponds to json.ErrRPCInvalidParameter.

	InvalidParameterError struct {
		error
	}

	// ParseError describes a failed parse due to bad user input.  It
	// corresponds to json.ErrRPCParse.

	ParseError struct {
		error
	}
)

type RPCInterface struct {
	MSG        interface{} `json="msg"`
	ERR        interface{} `json="err"`
	BlockCount int32       `json="blockcount"`

	Address              string `json="address"`
	createrawtransaction *string
	debuglevel           *string
	decoderawtransaction *json.TxRawDecodeResult
	decodescript         *json.DecodeScriptResult
	estimatefee          *float64
	generate             *[]string
	Getaddednodeinfo     *[]json.GetAddedNodeInfoResult
	Getbestblock         *json.GetBestBlockResult
	Getbestblockhash     *string
	Getblock             *json.GetBlockVerboseResult
	// Getblockcount         int32
	Getblockhash          *string
	Getblockheader        *json.GetBlockHeaderVerboseResult
	Getblocktemplate      *json.GetBlockTemplateResult
	Getblockchaininfo     *json.GetBlockChainInfoResult
	Getcfilter            *string
	Getcfilterheader      *string
	Getconnectioncount    *int32
	Getcurrentnet         *uint32
	Getdifficulty         *float64
	Getgenerate           *bool
	Gethashespersec       *float64
	Getheaders            *[]string
	Getinfo               *json.InfoWalletResult
	Getmempoolinfo        *json.GetMempoolInfoResult
	Getmininginfo         *json.GetMiningInfoResult
	Getnettotals          *json.GetNetTotalsResult
	Getnetworkhashps      *int64
	Getpeerinfo           *[]json.GetPeerInfoResult
	Getrawmempool         *json.GetRawMempoolVerboseResult
	Getrawtransaction     *json.TxRawResult
	Gettxout              *json.GetTxOutResult
	Searchrawtransactions *[]json.SearchRawTransactionsResult
	Sendrawtransaction    *string
	Stop                  *string
	Submitblock           *string
	Uptime                *int64
	Validateaddress       *json.ValidateAddressChainResult
	Verifychain           *bool
	Verifymessage         *bool
	Version               *map[string]json.VersionResult

	// Websocket commands.
	Session      *json.SessionResult
	Rescanblocks *[]json.RescannedBlock
}

// confirmed checks whether a transaction at height txHeight has met minconf
// confirmations for a blockchain at height curHeight.

func confirmed(minconf, txHeight, curHeight int32) bool {

	return confirms(txHeight, curHeight) >= minconf
}

// confirms returns the number of confirmations for a transaction in a block at
// height txHeight (or -1 for an unconfirmed tx) given the chain height
// curHeight.

func confirms(txHeight, curHeight int32) int32 {

	switch {

	case txHeight == -1, txHeight > curHeight:
		return 0
	default:
		return curHeight - txHeight + 1
	}

}

// requestHandler is a handler function to handle an unmarshaled and parsed
// request into a marshalable response.  If the error is a *json.RPCError
// or any of the above special error classes, the server will respond with
// the JSON-RPC appropiate error code.  All other errors use the wallet
// catch-all error code, json.ErrRPCWallet.

type requestHandler func(interface{}, *wallet.Wallet) (interface{}, error)

// requestHandlerChain is a requestHandler that also takes a parameter for

type requestHandlerChainRequired func(interface{}, *wallet.Wallet, *chain.RPCClient) (interface{}, error)

var rpcHandlers = map[string]struct {
	handler          requestHandler
	handlerWithChain requestHandlerChainRequired
	Handler          requestHandler
	HandlerWithChain requestHandlerChainRequired

	// Function variables cannot be compared against anything but nil, so
	// use a boolean to record whether help generation is necessary.  This
	// is used by the tests to ensure that help can be generated for every
	// implemented method.
	//
	// A single map and this bool is here is used rather than several maps
	// for the unimplemented handlers so every method has exactly one
	// handler function.
	noHelp bool
}{

	// Reference implementation wallet methods (implemented)

}

var RPCHandlers = &rpcHandlers

// unimplemented handles an unimplemented RPC request with the
// appropiate error.

func unimplemented(interface{}, *wallet.Wallet) (interface{}, error) {

	return nil, &json.RPCError{

		Code:    json.ErrRPCUnimplemented,
		Message: "Method unimplemented",
	}

}

// unsupported handles a standard bitcoind RPC request which is
// unsupported by btcwallet due to design differences.

func unsupported(interface{}, *wallet.Wallet) (interface{}, error) {

	return nil, &json.RPCError{

		Code:    -1,
		Message: "Request unsupported by mod",
	}

}

// lazyHandler is a closure over a requestHandler or passthrough request with
// the RPC server's wallet and chain server variables as part of the closure
// context.

type lazyHandler func() (interface{}, *json.RPCError)

// lazyApplyHandler looks up the best request handler func for the method,
// returning a closure that will execute it with the (required) wallet and
// (optional) consensus RPC server.  If no handlers are found and the
// chainClient is not nil, the returned handler performs RPC passthrough.

func lazyApplyHandler(request *json.Request, chainClient chain.Interface) lazyHandler {

	handlerData, ok := rpcHandlers[request.Method]

	if ok && handlerData.handlerWithChain != nil && WLT != nil && chainClient != nil {

		return func() (interface{}, *json.RPCError) {

			cmd, err := json.UnmarshalCmd(request)

			if err != nil {

				return nil, json.ErrRPCInvalidRequest
			}

			switch client := chainClient.(type) {

			case *chain.RPCClient:
				resp, err := handlerData.handlerWithChain(cmd,
					WLT, client)

				if err != nil {

					// return nil, jsonError(err)
				}

				return resp, nil
			default:

				return nil, &json.RPCError{

					Code:    -1,
					Message: "Chain RPC is inactive",
				}

			}

		}

	}

	if ok && handlerData.handler != nil && WLT != nil {

		return func() (interface{}, *json.RPCError) {

			cmd, err := json.UnmarshalCmd(request)

			if err != nil {

				return nil, json.ErrRPCInvalidRequest
			}

			resp, err := handlerData.handler(cmd, WLT)

			if err != nil {

				// return nil, jsonError(err)
			}

			return resp, nil
		}

	}

	// Fallback to RPC passthrough

	return func() (interface{}, *json.RPCError) {

		if chainClient == nil {

			return nil, &json.RPCError{

				Code:    -1,
				Message: "Chain RPC is inactive",
			}

		}

		switch client := chainClient.(type) {

		case *chain.RPCClient:
			resp, err := client.RawRequest(request.Method,
				request.Params)

			if err != nil {

				// return nil, jsonError(err)
			}

			return &resp, nil
		default:

			return nil, &json.RPCError{

				Code:    -1,
				Message: "Chain RPC is inactive",
			}

		}

	}

}

// // makeResponse makes the JSON-RPC response struct for the result and error
// // returned by a requestHandler.  The returned response is not ready for
// // marshaling and sending off to a client, but must be

// func makeResponse(id, result interface{}, err error) json.Response {

// 	idPtr := idPointer(id)

// 	if err != nil {

// 		return json.Response{

// 			ID:    idPtr,
// 			Error: jsonError(err),
// 		}
// 	}
// 	resultBytes, err := json.Marshal(result)

// 	if err != nil {

// 		return json.Response{

// 			ID: idPtr,

// 			Error: &json.RPCError{

// 				Code:    json.ErrRPCInternal.Code,
// 				Message: "Unexpected error marshalling result",
// 			},
// 		}
// 	}

// 	return json.Response{

// 		ID:     idPtr,
// 		Result: json.RawMessage(resultBytes),
// 	}
// }

// jsonError creates a JSON-RPC error from the Go error.

func jsonError(err error) *json.RPCError {

	if err == nil {

		return nil
	}

	code := json.ErrRPCWallet

	switch e := err.(type) {

	case json.RPCError:
		return &e
	case *json.RPCError:
		return e
	case DeserializationError:
		code = json.ErrRPCDeserialization
	case InvalidParameterError:
		code = json.ErrRPCInvalidParameter
	case ParseError:
		code = json.ErrRPCParse.Code
	case waddrmgr.ManagerError:

		switch e.ErrorCode {

		case waddrmgr.ErrWrongPassphrase:
			code = json.ErrRPCWalletPassphraseIncorrect
		}

	}

	return &json.RPCError{

		Code:    code,
		Message: err.Error(),
	}

}

// makeMultiSigScript is a helper function to combine common logic for
// AddMultiSig and CreateMultiSig.

func makeMultiSigScript(w *wallet.Wallet, keys []string, nRequired int) ([]byte, error) {

	keysesPrecious := make([]*util.AddressPubKey, len(keys))

	// The address list will made up either of addreseses (pubkey hash), for
	// which we need to look up the keys in wallet, straight pubkeys, or a
	// mixture of the two.

	for i, a := range keys {

		// try to parse as pubkey address
		a, err := decodeAddress(a, WLT.ChainParams())

		if err != nil {

			return nil, err
		}

		switch addr := a.(type) {

		case *util.AddressPubKey:
			keysesPrecious[i] = addr
		default:
			pubKey, err := WLT.PubKeyForAddress(addr)

			if err != nil {

				return nil, err
			}

			pubKeyAddr, err := util.NewAddressPubKey(
				pubKey.SerializeCompressed(), WLT.ChainParams())

			if err != nil {

				return nil, err
			}

			keysesPrecious[i] = pubKeyAddr
		}

	}

	return txscript.MultiSigScript(keysesPrecious, nRequired)
}

// addMultiSigAddress handles an addmultisigaddress request by adding a
// multisig address to the given wallet.

func addMultiSigAddress(icmd interface{}) (interface{}, error) {

	cmd := icmd.(*json.AddMultisigAddressCmd)

	// If an account is specified, ensure that is the imported account.

	if cmd.Account != nil && *cmd.Account != waddrmgr.ImportedAddrAccountName {

		// return nil, &ErrNotImportedAccount
	}

	secp256k1Addrs := make([]util.Address, len(cmd.Keys))

	for i, k := range cmd.Keys {

		addr, err := decodeAddress(k, WLT.ChainParams())

		if err != nil {

			// return nil, ParseError{err}
		}

		secp256k1Addrs[i] = addr
	}

	script, err := WLT.MakeMultiSigScript(secp256k1Addrs, cmd.NRequired)

	if err != nil {

		return nil, err
	}

	p2shAddr, err := WLT.ImportP2SHRedeemScript(script)

	if err != nil {

		return nil, err
	}

	return p2shAddr.EncodeAddress(), nil
}

// createMultiSig handles an createmultisig request by returning a
// multisig address for the given inputs.

func createMultiSig(icmd interface{}) (interface{}, error) {

	cmd := icmd.(*json.CreateMultisigCmd)

	script, err := makeMultiSigScript(WLT, cmd.Keys, cmd.NRequired)

	if err != nil {

		// return nil, ParseError{err}
	}

	address, err := util.NewAddressScriptHash(script, WLT.ChainParams())

	if err != nil {

		// above is a valid script, shouldn't happen.
		return nil, err
	}

	return json.CreateMultiSigResult{

			Address:      address.EncodeAddress(),
			RedeemScript: hex.EncodeToString(script),
		},
		nil
}

// dumpPrivKey handles a dumpprivkey request with the private key
// for a single address, or an appropiate error if the wallet
// is locked.

func dumpPrivKey(icmd interface{}) (interface{}, error) {

	cmd := icmd.(*json.DumpPrivKeyCmd)

	addr, err := decodeAddress(cmd.Address, WLT.ChainParams())

	if err != nil {

		return nil, err
	}

	key, err := WLT.DumpWIFPrivateKey(addr)

	if waddrmgr.IsError(err, waddrmgr.ErrLocked) {

		// Address was found, but the private key isn't
		// accessible.
		// return nil, &ErrWalletUnlockNeeded
	}

	return key, err
}

// dumpWallet handles a dumpwallet request by returning  all private
// keys in a wallet, or an appropiate error if the wallet is locked.
// TODO: finish this to match bitcoind by writing the dump to a file.

func dumpWallet(icmd interface{}) (interface{}, error) {

	keys, err := WLT.DumpPrivKeys()

	if waddrmgr.IsError(err, waddrmgr.ErrLocked) {

		// return nil, &ErrWalletUnlockNeeded
	}

	return keys, err
}

// getAddressesByAccount handles a getaddressesbyaccount request by returning
// all addresses for an account, or an error if the requested account does
// not exist.

func getAddressesByAccount(icmd interface{}) (interface{}, error) {

	cmd := icmd.(*json.GetAddressesByAccountCmd)

	account, err := WLT.AccountNumber(waddrmgr.KeyScopeBIP0044, cmd.Account)

	if err != nil {

		return nil, err
	}

	addrs, err := WLT.AccountAddresses(account)

	if err != nil {

		return nil, err
	}

	addrStrs := make([]string, len(addrs))

	for i, a := range addrs {

		addrStrs[i] = a.EncodeAddress()
	}

	return addrStrs, nil
}

// getBalance handles a getbalance request by returning the balance for an
// account (wallet), or an error if the requested account does not
// exist.

func getBalance(icmd interface{}) (interface{}, error) {

	cmd := icmd.(*json.GetBalanceCmd)

	var balance util.Amount
	var err error
	accountName := "*"

	if cmd.Account != nil {

		accountName = *cmd.Account
	}

	if accountName == "*" {

		balance, err = WLT.CalculateBalance(int32(*cmd.MinConf))

		if err != nil {

			return nil, err
		}

	} else {

		var account uint32
		account, err = WLT.AccountNumber(waddrmgr.KeyScopeBIP0044, accountName)

		if err != nil {

			return nil, err
		}

		bals, err := WLT.CalculateAccountBalances(account, int32(*cmd.MinConf))

		if err != nil {

			return nil, err
		}

		balance = bals.Spendable
	}

	return balance.ToDUO(), nil
}

// getBestBlock handles a getbestblock request by returning a JSON object
// with the height and hash of the most recently processed block.

func (rpc *RPCInterface) GetBestBlock() {

	// blk := WLT.Manager.SyncedTo()

	// result := &json.GetBestBlockResult{

	// 	Hash:   blk.Hash.String(),
	// 	Height: blk.Height,
	// }
	// rpc.GetBestBlock = result
}

// getBestBlockHash handles a getbestblockhash request by returning the hash
// of the most recently processed block.

func (rpc *RPCInterface) getBestBlockHash(icmd interface{}) (interface{}, error) {

	blk := WLT.Manager.SyncedTo()
	return blk.Hash.String(), nil
}

// getBlockCount handles a getblockcount request by returning the chain height
// of the most recently processed block.

// func (rpc *RPCInterface) GetBlockCount() {

func (rpc *RPCInterface) GetBlockCount() {

	blk := WLT.Manager.SyncedTo()
	rpc.BlockCount = blk.Height
	// rpc.BlockCount = "New block: " + fmt.Sprint(blk.Height)
}

// getInfo handles a getinfo request by returning the a structure containing
// information about the current state of btcwallet.
// exist.

func (rpc *RPCInterface) GetInfo() {

	// Call down to pod for all of the information in this command known
	// by them.
	info, err := WLT.ChainClient().(*chain.RPCClient).GetInfo()
	// info, err := chainClient.GetInfo()

	if err != nil {

		// return nil, err
	}

	bal, err := WLT.CalculateBalance(1)

	if err != nil {

		// return nil, err
	}

	// TODO(davec): This should probably have a database version as opposed
	// to using the manager version.
	info.WalletVersion = int32(waddrmgr.LatestMgrVersion)
	info.Balance = bal.ToDUO()
	info.PaytxFee = float64(txrules.DefaultRelayFeePerKb)
	// We don't set the following since they don't make much sense in the
	// wallet architecture:
	//  - unlocked_until
	//  - errors
	rpc.Getinfo = info
}

func decodeAddress(s string, params *chaincfg.Params) (util.Address, error) {

	addr, err := util.DecodeAddress(s, params)

	if err != nil {

		msg := fmt.Sprintf("Invalid address %q: decode failed with %#q", s, err)

		return nil, &json.RPCError{

			Code:    json.ErrRPCInvalidAddressOrKey,
			Message: msg,
		}

	}

	if !addr.IsForNet(params) {

		msg := fmt.Sprintf("Invalid address %q: not intended for use on %s",
			addr, params.Name)

		return nil, &json.RPCError{

			Code:    json.ErrRPCInvalidAddressOrKey,
			Message: msg,
		}

	}

	return addr, nil
}

// getAccount handles a getaccount request by returning the account name
// associated with a single address.

func getAccount(icmd interface{}) (interface{}, error) {

	cmd := icmd.(*json.GetAccountCmd)

	addr, err := decodeAddress(cmd.Address, WLT.ChainParams())

	if err != nil {

		return nil, err
	}

	// Fetch the associated account
	account, err := WLT.AccountOfAddress(addr)

	if err != nil {

		// return nil, &ErrAddressNotInWallet
	}

	acctName, err := WLT.AccountName(waddrmgr.KeyScopeBIP0044, account)

	if err != nil {

		// return nil, &ErrAccountNameNotFound
	}

	return acctName, nil
}

// getAccountAddress handles a getaccountaddress by returning the most
// recently-created chained address that has not yet been used (does not yet
// appear in the blockchain, or any tx that has arrived in the pod mempool).
// If the most recently-requested address has been used, a new address (the
// next chained address in the keypool) is used.  This can fail if the keypool
// runs out (and will return json.ErrRPCWalletKeypoolRanOut if that happens).

func getAccountAddress(icmd interface{}) (interface{}, error) {

	cmd := icmd.(*json.GetAccountAddressCmd)

	account, err := WLT.AccountNumber(waddrmgr.KeyScopeBIP0044, cmd.Account)

	if err != nil {

		return nil, err
	}

	addr, err := WLT.CurrentAddress(account, waddrmgr.KeyScopeBIP0044)

	if err != nil {

		return nil, err
	}

	return addr.EncodeAddress(), err
}

// getUnconfirmedBalance handles a getunconfirmedbalance extension request
// by returning the current unconfirmed balance of an account.

func (rpc *RPCInterface) getUnconfirmedBalance(icmd interface{}) (interface{}, error) {

	cmd := icmd.(*json.GetUnconfirmedBalanceCmd)

	acctName := "default"

	if cmd.Account != nil {

		acctName = *cmd.Account
	}

	account, err := WLT.AccountNumber(waddrmgr.KeyScopeBIP0044, acctName)

	if err != nil {

		return nil, err
	}

	bals, err := WLT.CalculateAccountBalances(account, 1)

	if err != nil {

		return nil, err
	}

	return (bals.Total - bals.Spendable).ToDUO(), nil
}

// importPrivKey handles an importprivkey request by parsing
// a WIF-encoded private key and adding it to an account.

func importPrivKey(icmd interface{}) (interface{}, error) {

	cmd := icmd.(*json.ImportPrivKeyCmd)

	// Ensure that private keys are only imported to the correct account.
	//
	// Yes, Label is the account name.

	if cmd.Label != nil && *cmd.Label != waddrmgr.ImportedAddrAccountName {

		// return nil, &ErrNotImportedAccount
	}

	wif, err := util.DecodeWIF(cmd.PrivKey)

	if err != nil {

		return nil, &json.RPCError{

			Code:    json.ErrRPCInvalidAddressOrKey,
			Message: "WIF decode failed: " + err.Error(),
		}

	}

	if !wif.IsForNet(WLT.ChainParams()) {

		return nil, &json.RPCError{

			Code:    json.ErrRPCInvalidAddressOrKey,
			Message: "Key is not intended for " + WLT.ChainParams().Name,
		}

	}

	// Import the private key, handling any errors.
	_, err = WLT.ImportPrivateKey(waddrmgr.KeyScopeBIP0044, wif, nil, *cmd.Rescan)

	switch {

	case waddrmgr.IsError(err, waddrmgr.ErrDuplicateAddress):
		// Do not return duplicate key errors to the client.
		return nil, nil
	case waddrmgr.IsError(err, waddrmgr.ErrLocked):
		// return nil, &ErrWalletUnlockNeeded
	}

	return nil, err
}

// keypoolRefill handles the keypoolrefill command. Since we handle the keypool
// automatically this does nothing since refilling is never manually required.

func keypoolRefill(icmd interface{}) (interface{}, error) {

	return nil, nil
}

// createNewAccount handles a createnewaccount request by creating and
// returning a new account. If the last account has no transaction history
// as per BIP 0044 a new account cannot be created so an error will be returned.

func createNewAccount(icmd interface{}) (interface{}, error) {

	cmd := icmd.(*json.CreateNewAccountCmd)

	// The wildcard * is reserved by the rpc server with the special meaning
	// of "all accounts", so disallow naming accounts to this string.

	if cmd.Account == "*" {

		// return nil, &ErrReservedAccountName
	}

	_, err := WLT.NextAccount(waddrmgr.KeyScopeBIP0044, cmd.Account)

	if waddrmgr.IsError(err, waddrmgr.ErrLocked) {

		return nil, &json.RPCError{

			Code: json.ErrRPCWalletUnlockNeeded,
			Message: "Creating an account requires the wallet to be unlocked. " +
				"Enter the wallet passphrase with walletpassphrase to unlock",
		}

	}

	return nil, err
}

// renameAccount handles a renameaccount request by renaming an account.
// If the account does not exist an appropiate error will be returned.

func renameAccount(icmd interface{}) (interface{}, error) {

	cmd := icmd.(*json.RenameAccountCmd)

	// The wildcard * is reserved by the rpc server with the special meaning
	// of "all accounts", so disallow naming accounts to this string.

	if cmd.NewAccount == "*" {

		// return nil, &ErrReservedAccountName
	}

	// Check that given account exists
	account, err := WLT.AccountNumber(waddrmgr.KeyScopeBIP0044, cmd.OldAccount)

	if err != nil {

		return nil, err
	}

	return nil, WLT.RenameAccount(waddrmgr.KeyScopeBIP0044, account, cmd.NewAccount)
}

// getNewAddress handles a getnewaddress request by returning a new
// address for an account.  If the account does not exist an appropiate
// error is returned.
// TODO: Follow BIP 0044 and warn if number of unused addresses exceeds
// the gap limit.

func (rpc *RPCInterface) GetNewAddress() {

	acctName := "default"
	account, err := WLT.AccountNumber(waddrmgr.KeyScopeBIP0044, acctName)

	if err != nil {

		rpc.ERR = err
	}

	addr, err := WLT.NewAddress(account, waddrmgr.KeyScopeBIP0044)

	if err != nil {

		rpc.ERR = err
	}

	// Return the new payment address string.
	rpc.MSG = "Created new address:" + addr.EncodeAddress()
	rpc.Address = addr.EncodeAddress()
	fmt.Println("dsdsdsdsdsdsdsd", rpc.Address)
	fmt.Println("MSGMSGMSGMSGMSGMSGMSGMSG", rpc.MSG)
}

// getRawChangeAddress handles a getrawchangeaddress request by creating
// and returning a new change address for an account.
//
// Note: bitcoind allows specifying the account as an optional parameter,
// but ignores the parameter.

func getRawChangeAddress(icmd interface{}) (interface{}, error) {

	cmd := icmd.(*json.GetRawChangeAddressCmd)

	acctName := "default"

	if cmd.Account != nil {

		acctName = *cmd.Account
	}

	account, err := WLT.AccountNumber(waddrmgr.KeyScopeBIP0044, acctName)

	if err != nil {

		return nil, err
	}

	addr, err := WLT.NewChangeAddress(account, waddrmgr.KeyScopeBIP0044)

	if err != nil {

		return nil, err
	}

	// Return the new payment address string.
	return addr.EncodeAddress(), nil
}

// getReceivedByAccount handles a getreceivedbyaccount request by returning
// the total amount received by addresses of an account.

func getReceivedByAccount(icmd interface{}) (interface{}, error) {

	cmd := icmd.(*json.GetReceivedByAccountCmd)

	account, err := WLT.AccountNumber(waddrmgr.KeyScopeBIP0044, cmd.Account)

	if err != nil {

		return nil, err
	}

	// TODO: This is more inefficient that it could be, but the entire
	// algorithm is already dominated by reading every transaction in the
	// wallet's history.
	results, err := WLT.TotalReceivedForAccounts(
		waddrmgr.KeyScopeBIP0044, int32(*cmd.MinConf),
	)

	if err != nil {

		return nil, err
	}

	acctIndex := int(account)

	if account == waddrmgr.ImportedAddrAccount {

		acctIndex = len(results) - 1
	}

	return results[acctIndex].TotalReceived.ToDUO(), nil
}

// getReceivedByAddress handles a getreceivedbyaddress request by returning
// the total amount received by a single address.

func getReceivedByAddress(icmd interface{}) (interface{}, error) {

	cmd := icmd.(*json.GetReceivedByAddressCmd)

	addr, err := decodeAddress(cmd.Address, WLT.ChainParams())

	if err != nil {

		return nil, err
	}

	total, err := WLT.TotalReceivedForAddr(addr, int32(*cmd.MinConf))

	if err != nil {

		return nil, err
	}

	return total.ToDUO(), nil
}

// getTransaction handles a gettransaction request by returning details about
// a single transaction saved by wallet.

func getTransaction(icmd interface{}) (interface{}, error) {

	cmd := icmd.(*json.GetTransactionCmd)

	txHash, err := chainhash.NewHashFromStr(cmd.Txid)

	if err != nil {

		return nil, &json.RPCError{

			Code:    json.ErrRPCDecodeHexString,
			Message: "Transaction hash string decode failed: " + err.Error(),
		}

	}

	details, err := wallet.UnstableAPI(WLT).TxDetails(txHash)

	if err != nil {

		return nil, err
	}

	if details == nil {

		// return nil, &ErrNoTransactionInfo
	}

	syncBlock := WLT.Manager.SyncedTo()

	// TODO: The serialized transaction is already in the DB, so
	// reserializing can be avoided here.
	var txBuf bytes.Buffer
	txBuf.Grow(details.MsgTx.SerializeSize())
	err = details.MsgTx.Serialize(&txBuf)

	if err != nil {

		return nil, err
	}

	// TODO: Add a "generated" field to this result type.  "generated":true
	// is only added if the transaction is a coinbase.

	ret := json.GetTransactionResult{

		TxID:            cmd.Txid,
		Hex:             hex.EncodeToString(txBuf.Bytes()),
		Time:            details.Received.Unix(),
		TimeReceived:    details.Received.Unix(),
		WalletConflicts: []string{}, // Not saved
		//Generated:     blockchain.IsCoinBaseTx(&details.MsgTx),
	}

	if details.Block.Height != -1 {

		ret.BlockHash = details.Block.Hash.String()
		ret.BlockTime = details.Block.Time.Unix()
		ret.Confirmations = int64(confirms(details.Block.Height, syncBlock.Height))
	}

	var (
		debitTotal  util.Amount
		creditTotal util.Amount // Excludes change
		fee         util.Amount
		feeF64      float64
	)

	for _, deb := range details.Debits {

		debitTotal += deb.Amount
	}

	for _, cred := range details.Credits {

		if !cred.Change {

			creditTotal += cred.Amount
		}

	}

	// Fee can only be determined if every input is a debit.

	if len(details.Debits) == len(details.MsgTx.TxIn) {

		var outputTotal util.Amount

		for _, output := range details.MsgTx.TxOut {

			outputTotal += util.Amount(output.Value)
		}

		fee = debitTotal - outputTotal
		feeF64 = fee.ToDUO()
	}

	if len(details.Debits) == 0 {

		// Credits must be set later, but since we know the full length
		// of the details slice, allocate it with the correct cap.
		ret.Details = make([]json.GetTransactionDetailsResult, 0, len(details.Credits))

	} else {

		ret.Details = make([]json.GetTransactionDetailsResult, 1, len(details.Credits)+1)

		ret.Details[0] = json.GetTransactionDetailsResult{

			// Fields left zeroed:
			//   InvolvesWatchOnly
			//   Account
			//   Address
			//   Vout
			//
			// TODO(jrick): Address and Vout should always be set,
			// but we're doing the wrong thing here by not matching
			// core.  Instead, gettransaction should only be adding
			// details for transaction outputs, just like
			// listtransactions (but using the short result format).
			Category: "send",
			Amount:   (-debitTotal).ToDUO(), // negative since it is a send
			Fee:      &feeF64,
		}

		ret.Fee = feeF64
	}

	credCat := wallet.RecvCategory(details, syncBlock.Height, WLT.ChainParams()).String()

	for _, cred := range details.Credits {

		// Change is ignored.

		if cred.Change {

			continue
		}

		var address string
		var accountName string
		_, addrs, _, err := txscript.ExtractPkScriptAddrs(
			details.MsgTx.TxOut[cred.Index].PkScript, WLT.ChainParams())

		if err == nil && len(addrs) == 1 {

			addr := addrs[0]
			address = addr.EncodeAddress()
			account, err := WLT.AccountOfAddress(addr)

			if err == nil {

				name, err := WLT.AccountName(waddrmgr.KeyScopeBIP0044, account)

				if err == nil {

					accountName = name
				}

			}

		}

		ret.Details = append(ret.Details, json.GetTransactionDetailsResult{

			// Fields left zeroed:
			//   InvolvesWatchOnly
			//   Fee
			Account:  accountName,
			Address:  address,
			Category: credCat,
			Amount:   cred.Amount.ToDUO(),
			Vout:     cred.Index,
		})

	}

	ret.Amount = creditTotal.ToDUO()
	return ret, nil
}

// These generators create the following global variables in this package:
//
//   var localeHelpDescs map[string]func() map[string]string
//   var requestUsages string
//
// localeHelpDescs maps from locale strings (e.g. "en_US") to a function that
// builds a map of help texts for each RPC server method.  This prevents help
// text maps for every locale map from being rooted and created during init.
// Instead, the appropiate function is looked up when help text is first needed
// using the current locale and saved to the global below for futher reuse.
//
// requestUsages contains single line usages for every supported request,
// separated by newlines.  It is set during init.  These usages are used for all
// locales.
//
//go:generate go run ../../internal/rpchelp/genrpcserverhelp.go legacyrpc
//go:generate gofmt -w rpcserverhelp.go

// listAccounts handles a listaccounts request by returning a map of account
// names to their balances.

func listAccounts(icmd interface{}) (interface{}, error) {

	cmd := icmd.(*json.ListAccountsCmd)

	accountBalances := map[string]float64{}
	results, err := WLT.AccountBalances(waddrmgr.KeyScopeBIP0044, int32(*cmd.MinConf))

	if err != nil {

		return nil, err
	}

	for _, result := range results {

		accountBalances[result.AccountName] = result.AccountBalance.ToDUO()
	}

	// Return the map.  This will be marshaled into a JSON object.
	return accountBalances, nil
}

// listLockUnspent handles a listlockunspent request by returning an slice of
// all locked outpoints.

func listLockUnspent(icmd interface{}) (interface{}, error) {

	return WLT.LockedOutpoints(), nil
}

// listReceivedByAccount handles a listreceivedbyaccount request by returning
// a slice of objects, each one containing:
//  "account": the receiving account;
//  "amount": total amount received by the account;
//  "confirmations": number of confirmations of the most recent transaction.
// It takes two parameters:
//  "minconf": minimum number of confirmations to consider a transaction -
//             default: one;
//  "includeempty": whether or not to include addresses that have no transactions -
//                  default: false.

func listReceivedByAccount(icmd interface{}) (interface{}, error) {

	cmd := icmd.(*json.ListReceivedByAccountCmd)

	results, err := WLT.TotalReceivedForAccounts(
		waddrmgr.KeyScopeBIP0044, int32(*cmd.MinConf),
	)

	if err != nil {

		return nil, err
	}

	jsonResults := make([]json.ListReceivedByAccountResult, 0, len(results))

	for _, result := range results {

		jsonResults = append(jsonResults, json.ListReceivedByAccountResult{

			Account:       result.AccountName,
			Amount:        result.TotalReceived.ToDUO(),
			Confirmations: uint64(result.LastConfirmation),
		})

	}

	return jsonResults, nil
}

// listReceivedByAddress handles a listreceivedbyaddress request by returning
// a slice of objects, each one containing:
//  "account": the account of the receiving address;
//  "address": the receiving address;
//  "amount": total amount received by the address;
//  "confirmations": number of confirmations of the most recent transaction.
// It takes two parameters:
//  "minconf": minimum number of confirmations to consider a transaction -
//             default: one;
//  "includeempty": whether or not to include addresses that have no transactions -
//                  default: false.

func listReceivedByAddress(icmd interface{}) (interface{}, error) {

	cmd := icmd.(*json.ListReceivedByAddressCmd)

	// Intermediate data for each address.

	type AddrData struct {

		// Total amount received.
		amount util.Amount
		// Number of confirmations of the last transaction.
		confirmations int32
		// Hashes of transactions which include an output paying to the address
		tx []string
		// Account which the address belongs to
		account string
	}

	syncBlock := WLT.Manager.SyncedTo()

	// Intermediate data for all addresses.
	allAddrData := make(map[string]AddrData)
	// Create an AddrData entry for each active address in the account.
	// Otherwise we'll just get addresses from transactions later.
	sortedAddrs, err := WLT.SortedActivePaymentAddresses()

	if err != nil {

		return nil, err
	}

	for _, address := range sortedAddrs {

		// There might be duplicates, just overwrite them.
		allAddrData[address] = AddrData{}
	}

	minConf := *cmd.MinConf
	var endHeight int32

	if minConf == 0 {

		endHeight = -1

	} else {

		endHeight = syncBlock.Height - int32(minConf) + 1
	}

	err = wallet.UnstableAPI(WLT).RangeTransactions(0, endHeight, func(details []wtxmgr.TxDetails) (bool, error) {

		confirmations := confirms(details[0].Block.Height, syncBlock.Height)

		for _, tx := range details {

			for _, cred := range tx.Credits {

				pkScript := tx.MsgTx.TxOut[cred.Index].PkScript
				_, addrs, _, err := txscript.ExtractPkScriptAddrs(
					pkScript, WLT.ChainParams())

				if err != nil {

					// Non standard script, skip.
					continue
				}

				for _, addr := range addrs {

					addrStr := addr.EncodeAddress()
					addrData, ok := allAddrData[addrStr]

					if ok {

						addrData.amount += cred.Amount
						// Always overwrite confirmations with newer ones.
						addrData.confirmations = confirmations

					} else {

						addrData = AddrData{

							amount:        cred.Amount,
							confirmations: confirmations,
						}

					}

					addrData.tx = append(addrData.tx, tx.Hash.String())
					allAddrData[addrStr] = addrData
				}

			}

		}

		return false, nil
	})

	if err != nil {

		return nil, err
	}

	// Massage address data into output format.
	numAddresses := len(allAddrData)
	ret := make([]json.ListReceivedByAddressResult, numAddresses, numAddresses)
	idx := 0

	for address, addrData := range allAddrData {

		ret[idx] = json.ListReceivedByAddressResult{

			Address:       address,
			Amount:        addrData.amount.ToDUO(),
			Confirmations: uint64(addrData.confirmations),
			TxIDs:         addrData.tx,
		}

		idx++
	}

	return ret, nil
}

// listSinceBlock handles a listsinceblock request by returning an array of maps
// with details of sent and received wallet transactions since the given block.

func listSinceBlock(icmd interface{}, chainClient *chain.RPCClient) (interface{}, error) {

	cmd := icmd.(*json.ListSinceBlockCmd)

	syncBlock := WLT.Manager.SyncedTo()
	targetConf := int64(*cmd.TargetConfirmations)

	// For the result we need the block hash for the last block counted
	// in the blockchain due to confirmations. We send this off now so that
	// it can arrive asynchronously while we figure out the rest.
	gbh := chainClient.GetBlockHashAsync(int64(syncBlock.Height) + 1 - targetConf)

	var start int32

	if cmd.BlockHash != nil {

		hash, err := chainhash.NewHashFromStr(*cmd.BlockHash)

		if err != nil {

			// return nil, DeserializationError{err}
		}

		block, err := chainClient.GetBlockVerboseTx(hash)

		if err != nil {

			return nil, err
		}

		start = int32(block.Height) + 1
	}

	txInfoList, err := WLT.ListSinceBlock(start, -1, syncBlock.Height)

	if err != nil {

		return nil, err
	}

	// Done with work, get the response.
	blockHash, err := gbh.Receive()

	if err != nil {

		return nil, err
	}

	res := json.ListSinceBlockResult{

		Transactions: txInfoList,
		LastBlock:    blockHash.String(),
	}

	return res, nil
}

// listTransactions handles a listtransactions request by returning an
// array of maps with details of sent and recevied wallet transactions.

func listTransactions(icmd interface{}) (interface{}, error) {

	cmd := icmd.(*json.ListTransactionsCmd)

	// TODO: ListTransactions does not currently understand the difference
	// between transactions pertaining to one account from another.  This
	// will be resolved when wtxmgr is combined with the waddrmgr namespace.

	if cmd.Account != nil && *cmd.Account != "*" {

		// For now, don't bother trying to continue if the user
		// specified an account, since this can't be (easily or
		// efficiently) calculated.

		return nil, &json.RPCError{

			Code:    json.ErrRPCWallet,
			Message: "Transactions are not yet grouped by account",
		}

	}

	return WLT.ListTransactions(*cmd.From, *cmd.Count)
}

// listAddressTransactions handles a listaddresstransactions request by
// returning an array of maps with details of spent and received wallet
// transactions.  The form of the reply is identical to listtransactions,
// but the array elements are limited to transaction details which are
// about the addresess included in the request.

func listAddressTransactions(icmd interface{}) (interface{}, error) {

	cmd := icmd.(*json.ListAddressTransactionsCmd)

	if cmd.Account != nil && *cmd.Account != "*" {

		return nil, &json.RPCError{

			Code:    json.ErrRPCInvalidParameter,
			Message: "Listing transactions for addresses may only be done for all accounts",
		}

	}

	// Decode addresses.
	hash160Map := make(map[string]struct{})

	for _, addrStr := range cmd.Addresses {

		addr, err := decodeAddress(addrStr, WLT.ChainParams())

		if err != nil {

			return nil, err
		}

		hash160Map[string(addr.ScriptAddress())] = struct{}{}
	}

	return WLT.ListAddressTransactions(hash160Map)
}

// listAllTransactions handles a listalltransactions request by returning
// a map with details of sent and recevied wallet transactions.  This is
// similar to ListTransactions, except it takes only a single optional
// argument for the account name and replies with all transactions.

func listAllTransactions(icmd interface{}) (interface{}, error) {

	cmd := icmd.(*json.ListAllTransactionsCmd)

	if cmd.Account != nil && *cmd.Account != "*" {

		return nil, &json.RPCError{

			Code:    json.ErrRPCInvalidParameter,
			Message: "Listing all transactions may only be done for all accounts",
		}

	}

	return WLT.ListAllTransactions()
}

// listUnspent handles the listunspent command.

func listUnspent(icmd interface{}) (interface{}, error) {

	cmd := icmd.(*json.ListUnspentCmd)

	var addresses map[string]struct{}

	if cmd.Addresses != nil {

		addresses = make(map[string]struct{})
		// confirm that all of them are good:

		for _, as := range *cmd.Addresses {

			a, err := decodeAddress(as, WLT.ChainParams())

			if err != nil {

				return nil, err
			}

			addresses[a.EncodeAddress()] = struct{}{}
		}

	}

	return WLT.ListUnspent(int32(*cmd.MinConf), int32(*cmd.MaxConf), addresses)
}

// lockUnspent handles the lockunspent command.

func lockUnspent(icmd interface{}) (interface{}, error) {

	cmd := icmd.(*json.LockUnspentCmd)

	switch {

	case cmd.Unlock && len(cmd.Transactions) == 0:
		WLT.ResetLockedOutpoints()
	default:

		for _, input := range cmd.Transactions {

			txHash, err := chainhash.NewHashFromStr(input.Txid)

			if err != nil {

				// return nil, ParseError{err}
			}

			op := wire.OutPoint{Hash: *txHash, Index: input.Vout}

			if cmd.Unlock {

				WLT.UnlockOutpoint(op)

			} else {

				WLT.LockOutpoint(op)
			}

		}

	}

	return true, nil
}

// makeOutputs creates a slice of transaction outputs from a pair of address
// strings to amounts.  This is used to create the outputs to include in newly
// created transactions from a JSON object describing the output destinations
// and amounts.

func makeOutputs(pairs map[string]util.Amount, chainParams *chaincfg.Params) ([]*wire.TxOut, error) {

	outputs := make([]*wire.TxOut, 0, len(pairs))

	for addrStr, amt := range pairs {

		addr, err := util.DecodeAddress(addrStr, chainParams)

		if err != nil {

			return nil, fmt.Errorf("cannot decode address: %s", err)
		}

		pkScript, err := txscript.PayToAddrScript(addr)

		if err != nil {

			return nil, fmt.Errorf("cannot create txout script: %s", err)
		}

		outputs = append(outputs, wire.NewTxOut(int64(amt), pkScript))
	}

	return outputs, nil
}

// sendPairs creates and sends payment transactions.
// It returns the transaction hash in string format upon success
// All errors are returned in json.RPCError format
func sendPairs(w *wallet.Wallet, amounts map[string]util.Amount,

	account uint32, minconf int32, feeSatPerKb util.Amount) (string, error) {

	outputs, err := makeOutputs(amounts, WLT.ChainParams())

	if err != nil {

		return "", err
	}

	txHash, err := WLT.SendOutputs(outputs, account, minconf, feeSatPerKb)

	if err != nil {

		if err == txrules.ErrAmountNegative {

			// return "", ErrNeedPositiveAmount
		}

		if waddrmgr.IsError(err, waddrmgr.ErrLocked) {

			// return "", &ErrWalletUnlockNeeded
		}

		switch err.(type) {

		case json.RPCError:
			return "", err
		}

		return "", &json.RPCError{

			Code:    json.ErrRPCInternal.Code,
			Message: err.Error(),
		}

	}

	txHashStr := txHash.String()
	log.Infof("Successfully sent transaction %v", txHashStr)
	return txHashStr, nil
}

func isNilOrEmpty(s *string) bool {

	return s == nil || *s == ""
}

// sendFrom handles a sendfrom RPC request by creating a new transaction
// spending unspent transaction outputs for a wallet to another payment
// address.  Leftover inputs not sent to the payment address or a fee for
// the miner are sent back to a new address in the wallet.  Upon success,
// the TxID for the created transaction is returned.

func sendFrom(icmd interface{}, chainClient *chain.RPCClient) (interface{}, error) {

	cmd := icmd.(*json.SendFromCmd)

	// Transaction comments are not yet supported.  Error instead of
	// pretending to save them.

	if !isNilOrEmpty(cmd.Comment) || !isNilOrEmpty(cmd.CommentTo) {

		return nil, &json.RPCError{

			Code:    json.ErrRPCUnimplemented,
			Message: "Transaction comments are not yet supported",
		}

	}

	account, err := WLT.AccountNumber(
		waddrmgr.KeyScopeBIP0044, cmd.FromAccount,
	)

	if err != nil {

		return nil, err
	}

	// Check that signed integer parameters are positive.

	if cmd.Amount < 0 {

		// return nil, ErrNeedPositiveAmount
	}

	minConf := int32(*cmd.MinConf)

	if minConf < 0 {

		// return nil, ErrNeedPositiveMinconf
	}

	// Create map of address and amount pairs.
	amt, err := util.NewAmount(cmd.Amount)

	if err != nil {

		return nil, err
	}

	pairs := map[string]util.Amount{

		cmd.ToAddress: amt,
	}

	return sendPairs(WLT, pairs, account, minConf,
		txrules.DefaultRelayFeePerKb)
}

// sendMany handles a sendmany RPC request by creating a new transaction
// spending unspent transaction outputs for a wallet to any number of
// payment addresses.  Leftover inputs not sent to the payment address
// or a fee for the miner are sent back to a new address in the wallet.
// Upon success, the TxID for the created transaction is returned.

func sendMany(icmd interface{}) (interface{}, error) {

	cmd := icmd.(*json.SendManyCmd)

	// Transaction comments are not yet supported.  Error instead of
	// pretending to save them.

	if !isNilOrEmpty(cmd.Comment) {

		return nil, &json.RPCError{

			Code:    json.ErrRPCUnimplemented,
			Message: "Transaction comments are not yet supported",
		}

	}

	account, err := WLT.AccountNumber(waddrmgr.KeyScopeBIP0044, cmd.FromAccount)

	if err != nil {

		return nil, err
	}

	// Check that minconf is positive.
	minConf := int32(*cmd.MinConf)

	if minConf < 0 {

		// return nil, ErrNeedPositiveMinconf
	}

	// Recreate address/amount pairs, using dcrutil.Amount.
	pairs := make(map[string]util.Amount, len(cmd.Amounts))

	for k, v := range cmd.Amounts {

		amt, err := util.NewAmount(v)

		if err != nil {

			return nil, err
		}

		pairs[k] = amt
	}

	return sendPairs(WLT, pairs, account, minConf, txrules.DefaultRelayFeePerKb)
}

// sendToAddress handles a sendtoaddress RPC request by creating a new
// transaction spending unspent transaction outputs for a wallet to another
// payment address.  Leftover inputs not sent to the payment address or a fee
// for the miner are sent back to a new address in the wallet.  Upon success,
// the TxID for the created transaction is returned.

func (rpc *RPCInterface) SendToAddress(vaddress string, vlabel string, vamount interface{}) {

	amount, err := strconv.ParseFloat(vamount.(string), 64)

	amt, err := util.NewAmount(amount)

	if err != nil {

		rpc.ERR = jsonError(err)
	}

	// Check that signed integer parameters are positive.

	if amt < 0 {

		// return nil, ErrNeedPositiveAmount
	}

	// Mock up map of address and amount pairs.

	pairs := map[string]util.Amount{

		vaddress: amt,
	}

	// sendtoaddress always spends from the default account, this matches bitcoind
	txid, _ := sendPairs(WLT, pairs, waddrmgr.DefaultAccountNum, 1, txrules.DefaultRelayFeePerKb)

	if txid != "" {

		rpc.MSG = "Inserting unconfirmed transaction " + txid
	}

	// fmt.Println("TTTTwwwwwwwwwwwwwwwwwTTTTTTTTTTTTTTTTTTTTTTTTTTTTT", rpc.MSG)

}

// setTxFee sets the transaction fee per kilobyte added to transactions.

func setTxFee(icmd interface{}) (interface{}, error) {

	cmd := icmd.(*json.SetTxFeeCmd)

	// Check that amount is not negative.

	if cmd.Amount < 0 {

		// return nil, ErrNeedPositiveAmount
	}

	// A boolean true result is returned upon success.
	return true, nil
}

// signMessage signs the given message with the private key for the given
// address

func signMessage(icmd interface{}) (interface{}, error) {

	cmd := icmd.(*json.SignMessageCmd)

	addr, err := decodeAddress(cmd.Address, WLT.ChainParams())

	if err != nil {

		return nil, err
	}

	privKey, err := WLT.PrivKeyForAddress(addr)

	if err != nil {

		return nil, err
	}

	var buf bytes.Buffer
	wire.WriteVarString(&buf, 0, "Bitcoin Signed Message:\n")
	wire.WriteVarString(&buf, 0, cmd.Message)
	messageHash := chainhash.DoubleHashB(buf.Bytes())
	sigbytes, err := ec.SignCompact(ec.S256(), privKey,
		messageHash, true)

	if err != nil {

		return nil, err
	}

	return base64.StdEncoding.EncodeToString(sigbytes), nil
}

// signRawTransaction handles the signrawtransaction command.

func signRawTransaction(icmd interface{}, chainClient *chain.RPCClient) (interface{}, error) {

	cmd := icmd.(*json.SignRawTransactionCmd)

	serializedTx, err := decodeHexStr(cmd.RawTx)

	if err != nil {

		return nil, err
	}

	var tx wire.MsgTx
	err = tx.Deserialize(bytes.NewBuffer(serializedTx))

	if err != nil {

		// e := errors.New("TX decode failed")
		// return nil, DeserializationError{e}
	}

	var hashType txscript.SigHashType

	switch *cmd.Flags {

	case "ALL":
		hashType = txscript.SigHashAll
	case "NONE":
		hashType = txscript.SigHashNone
	case "SINGLE":
		hashType = txscript.SigHashSingle
	case "ALL|ANYONECANPAY":
		hashType = txscript.SigHashAll | txscript.SigHashAnyOneCanPay
	case "NONE|ANYONECANPAY":
		hashType = txscript.SigHashNone | txscript.SigHashAnyOneCanPay
	case "SINGLE|ANYONECANPAY":
		hashType = txscript.SigHashSingle | txscript.SigHashAnyOneCanPay
	default:
		// e := errors.New("Invalid sighash parameter")
		// return nil, InvalidParameterError{e}
	}

	// TODO: really we probably should look these up with pod anyway to
	// make sure that they match the blockchain if present.
	inputs := make(map[wire.OutPoint][]byte)
	scripts := make(map[string][]byte)
	var cmdInputs []json.RawTxInput

	if cmd.Inputs != nil {

		cmdInputs = *cmd.Inputs
	}

	for _, rti := range cmdInputs {

		inputHash, err := chainhash.NewHashFromStr(rti.Txid)

		if err != nil {

			// return nil, DeserializationError{err}
		}

		script, err := decodeHexStr(rti.ScriptPubKey)

		if err != nil {

			return nil, err
		}

		// redeemScript is only actually used iff the user provided
		// private keys. In which case, it is used to get the scripts
		// for signing. If the user did not provide keys then we always
		// get scripts from the wallet.
		// Empty strings are ok for this one and hex.DecodeString will
		// DTRT.

		if cmd.PrivKeys != nil && len(*cmd.PrivKeys) != 0 {

			redeemScript, err := decodeHexStr(rti.RedeemScript)

			if err != nil {

				return nil, err
			}

			addr, err := util.NewAddressScriptHash(redeemScript,
				WLT.ChainParams())

			if err != nil {

				// return nil, DeserializationError{err}
			}

			scripts[addr.String()] = redeemScript
		}

		inputs[wire.OutPoint{

			Hash:  *inputHash,
			Index: rti.Vout,
		}] = script
	}

	// Now we go and look for any inputs that we were not provided by
	// querying pod with getrawtransaction. We queue up a bunch of async
	// requests and will wait for replies after we have checked the rest of
	// the arguments.
	requested := make(map[wire.OutPoint]rpcclient.FutureGetTxOutResult)

	for _, txIn := range tx.TxIn {

		// Did we get this outpoint from the arguments?

		if _, ok := inputs[txIn.PreviousOutPoint]; ok {

			continue
		}

		// Asynchronously request the output script.
		requested[txIn.PreviousOutPoint] = chainClient.GetTxOutAsync(
			&txIn.PreviousOutPoint.Hash, txIn.PreviousOutPoint.Index,
			true)
	}

	// Parse list of private keys, if present. If there are any keys here
	// they are the keys that we may use for signing. If empty we will
	// use any keys known to us already.
	var keys map[string]*util.WIF

	if cmd.PrivKeys != nil {

		keys = make(map[string]*util.WIF)

		for _, key := range *cmd.PrivKeys {

			wif, err := util.DecodeWIF(key)

			if err != nil {

				// return nil, DeserializationError{err}
			}

			if !wif.IsForNet(WLT.ChainParams()) {

				// s := "key network doesn't match wallet's"
				// return nil, DeserializationError{errors.New(s)}
			}

			addr, err := util.NewAddressPubKey(wif.SerializePubKey(),
				WLT.ChainParams())

			if err != nil {

				// return nil, DeserializationError{err}
			}

			keys[addr.EncodeAddress()] = wif
		}

	}

	// We have checked the rest of the args. now we can collect the async
	// txs. TODO: If we don't mind the possibility of wasting work we could
	// move waiting to the following loop and be slightly more asynchronous.

	for outPoint, resp := range requested {

		result, err := resp.Receive()

		if err != nil {

			return nil, err
		}

		script, err := hex.DecodeString(result.ScriptPubKey.Hex)

		if err != nil {

			return nil, err
		}

		inputs[outPoint] = script
	}

	// All args collected. Now we can sign all the inputs that we can.
	// `complete' denotes that we successfully signed all outputs and that
	// all scripts will run to completion. This is returned as part of the
	// reply.
	signErrs, err := WLT.SignTransaction(&tx, hashType, inputs, keys, scripts)

	if err != nil {

		return nil, err
	}

	var buf bytes.Buffer
	buf.Grow(tx.SerializeSize())

	// All returned errors (not OOM, which panics) encounted during
	// bytes.Buffer writes are unexpected.

	if err = tx.Serialize(&buf); err != nil {

		panic(err)
	}

	signErrors := make([]json.SignRawTransactionError, 0, len(signErrs))

	for _, e := range signErrs {

		input := tx.TxIn[e.InputIndex]

		signErrors = append(signErrors, json.SignRawTransactionError{

			TxID:      input.PreviousOutPoint.Hash.String(),
			Vout:      input.PreviousOutPoint.Index,
			ScriptSig: hex.EncodeToString(input.SignatureScript),
			Sequence:  input.Sequence,
			Error:     e.Error.Error(),
		})

	}

	return json.SignRawTransactionResult{

			Hex:      hex.EncodeToString(buf.Bytes()),
			Complete: len(signErrors) == 0,
			Errors:   signErrors,
		},
		nil
}

// validateAddress handles the validateaddress command.

func validateAddress(icmd interface{}) (interface{}, error) {

	cmd := icmd.(*json.ValidateAddressCmd)

	result := json.ValidateAddressWalletResult{}
	addr, err := decodeAddress(cmd.Address, WLT.ChainParams())

	if err != nil {

		// Use result zero value (IsValid=false).
		return result, nil
	}

	// We could put whether or not the address is a script here,
	// by checking the type of "addr", however, the reference
	// implementation only puts that information if the script is
	// "ismine", and we follow that behaviour.
	result.Address = addr.EncodeAddress()
	result.IsValid = true

	ainfo, err := WLT.AddressInfo(addr)

	if err != nil {

		if waddrmgr.IsError(err, waddrmgr.ErrAddressNotFound) {

			// No additional information available about the address.
			return result, nil
		}

		return nil, err
	}

	// The address lookup was successful which means there is further
	// information about it available and it is "mine".
	result.IsMine = true
	acctName, err := WLT.AccountName(waddrmgr.KeyScopeBIP0044, ainfo.Account())

	if err != nil {

		// return nil, &ErrAccountNameNotFound
	}

	result.Account = acctName

	switch ma := ainfo.(type) {

	case waddrmgr.ManagedPubKeyAddress:
		result.IsCompressed = ma.Compressed()
		result.PubKey = ma.ExportPubKey()

	case waddrmgr.ManagedScriptAddress:
		result.IsScript = true

		// The script is only available if the manager is unlocked, so
		// just break out now if there is an error.
		script, err := ma.Script()

		if err != nil {

			break
		}

		result.Hex = hex.EncodeToString(script)

		// This typically shouldn't fail unless an invalid script was
		// imported.  However, if it fails for any reason, there is no
		// further information available, so just set the script type
		// a non-standard and break out now.
		class, addrs, reqSigs, err := txscript.ExtractPkScriptAddrs(
			script, WLT.ChainParams())

		if err != nil {

			result.Script = txscript.NonStandardTy.String()
			break
		}

		addrStrings := make([]string, len(addrs))

		for i, a := range addrs {

			addrStrings[i] = a.EncodeAddress()
		}

		result.Addresses = addrStrings

		// Multi-signature scripts also provide the number of required
		// signatures.
		result.Script = class.String()

		if class == txscript.MultiSigTy {

			result.SigsRequired = int32(reqSigs)
		}

	}

	return result, nil
}

// verifyMessage handles the verifymessage command by verifying the provided
// compact signature for the given address and message.

func verifyMessage(icmd interface{}) (interface{}, error) {

	cmd := icmd.(*json.VerifyMessageCmd)

	addr, err := decodeAddress(cmd.Address, WLT.ChainParams())

	if err != nil {

		return nil, err
	}

	// decode base64 signature
	sig, err := base64.StdEncoding.DecodeString(cmd.Signature)

	if err != nil {

		return nil, err
	}

	// Validate the signature - this just shows that it was valid at all.
	// we will compare it with the key next.
	var buf bytes.Buffer
	wire.WriteVarString(&buf, 0, "Bitcoin Signed Message:\n")
	wire.WriteVarString(&buf, 0, cmd.Message)
	expectedMessageHash := chainhash.DoubleHashB(buf.Bytes())
	pk, wasCompressed, err := ec.RecoverCompact(ec.S256(), sig,
		expectedMessageHash)

	if err != nil {

		return nil, err
	}

	var serializedPubKey []byte

	if wasCompressed {

		serializedPubKey = pk.SerializeCompressed()

	} else {

		serializedPubKey = pk.SerializeUncompressed()
	}

	// Verify that the signed-by address matches the given address

	switch checkAddr := addr.(type) {

	case *util.AddressPubKeyHash: // ok
		return bytes.Equal(util.Hash160(serializedPubKey), checkAddr.Hash160()[:]), nil
	case *util.AddressPubKey: // ok
		return string(serializedPubKey) == checkAddr.String(), nil
	default:
		return nil, errors.New("address type not supported")
	}

}

// walletIsLocked handles the walletislocked extension request by
// returning the current lock state (false for unlocked, true for locked)
// of an account.

func WalletIsLocked(icmd interface{}) (interface{}, error) {

	return WLT.Locked(), nil
}

// walletLock handles a walletlock request by locking the all account
// wallets, returning an error if any wallet is not encrypted (for example,
// a watching-only wallet).

func walletLock(icmd interface{}) (interface{}, error) {

	WLT.Lock()
	return nil, nil
}

// walletPassphrase responds to the walletpassphrase request by unlocking
// the wallet.  The decryption key is saved in the wallet until timeout
// seconds expires, after which the wallet is locked.

func (rpc *RPCInterface) WalletPassphrase(wpp string, tmo int64) {

	timeout := time.Second * time.Duration(tmo)
	var unlockAfter <-chan time.Time

	if timeout != 0 {

		unlockAfter = time.After(timeout)
	}

	err := jsonError(WLT.Unlock([]byte(wpp), unlockAfter))

	if err != nil {

		rpc.ERR = err.Message
	}

	fmt.Println("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaawppwppwppwppwppwppa", wpp)
	fmt.Println("aaaaaaaaaaaaaaaaaaaaaaaaarrrrrraaaaaaaaaaaaaaaaaaaaaaaaaatmotmotmotmotmotmoaa", tmo)
	fmt.Println("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", unlockAfter)
	fmt.Println("ffffffffffffffffffffffffffffffffffffffffffffffffff", rpc.MSG)
	fmt.Println("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	fmt.Println("aaaaaaaaaaaaaaaaaffffffffaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	fmt.Println("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	fmt.Println("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", rpc.MSG)

}

// func walletPassphrase(icmd interface{}) (interface{}, error) {

// 	cmd := icmd.(*json.WalletPassphraseCmd)

// 	timeout := time.Second * time.Duration(cmd.Timeout)
// 	var unlockAfter <-chan time.Time

// 	if timeout != 0 {

// 		unlockAfter = time.After(timeout)
// 	}
// 	err := WLT.Unlock([]byte(cmd.Passphrase), unlockAfter)
// 	return nil, err
// }

// walletPassphraseChange responds to the walletpassphrasechange request
// by unlocking all accounts with the provided old passphrase, and
// re-encrypting each private key with an AES key derived from the new
// passphrase.
//
// If the old passphrase is correct and the passphrase is changed, all
// wallets will be immediately locked.

func walletPassphraseChange(icmd interface{}) (interface{}, error) {

	cmd := icmd.(*json.WalletPassphraseChangeCmd)

	err := WLT.ChangePrivatePassphrase([]byte(cmd.OldPassphrase),
		[]byte(cmd.NewPassphrase))

	if waddrmgr.IsError(err, waddrmgr.ErrWrongPassphrase) {

		return nil, &json.RPCError{

			Code:    json.ErrRPCWalletPassphraseIncorrect,
			Message: "Incorrect passphrase",
		}

	}

	return nil, err
}

// decodeHexStr decodes the hex encoding of a string, possibly prepending a
// leading '0' character if there is an odd number of bytes in the hex string.
// This is to prevent an error for an invalid hex string when using an odd
// number of bytes when calling hex.Decode.

func decodeHexStr(hexStr string) ([]byte, error) {

	if len(hexStr)%2 != 0 {

		hexStr = "0" + hexStr
	}

	decoded, err := hex.DecodeString(hexStr)

	if err != nil {

		return nil, &json.RPCError{

			Code:    json.ErrRPCDecodeHexString,
			Message: "Hex string decode failed: " + err.Error(),
		}

	}

	return decoded, nil
}
