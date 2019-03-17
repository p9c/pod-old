package libs

import (
	"io/ioutil"
)

type VIcons map[string]string

var icoLogo, _ = ioutil.ReadFile("./gui/assets/icons/IcoLogo.svg")
var icoReceive, _ = ioutil.ReadFile("./gui/assets/icons/IcoReceive.svg")
var icoReceived, _ = ioutil.ReadFile("./gui/assets/icons/IcoReceived.svg")
var icoOverview, _ = ioutil.ReadFile("./gui/assets/icons/IcoOverview.svg")
var icoSettings, _ = ioutil.ReadFile("./gui/assets/icons/IcoSettings.svg")
var icoSend, _ = ioutil.ReadFile("./gui/assets/icons/IcoSend.svg")
var icoSent, _ = ioutil.ReadFile("./gui/assets/icons/IcoSent.svg")
var icoHistory, _ = ioutil.ReadFile("./gui/assets/icons/IcoHistory.svg")
var icoHelp, _ = ioutil.ReadFile("./gui/assets/icons/IcoHelp.svg")
var icoInfo, _ = ioutil.ReadFile("./gui/assets/icons/IcoInfo.svg")
var icoPeers, _ = ioutil.ReadFile("./gui/assets/icons/IcoPeers.svg")
var icoBlocks, _ = ioutil.ReadFile("./gui/assets/icons/IcoBlocks.svg")
var icoBalance, _ = ioutil.ReadFile("./gui/assets/icons/IcoBalance.svg")
var icoLink, _ = ioutil.ReadFile("./gui/assets/icons/IcoLink.svg")
var icoLinkOff, _ = ioutil.ReadFile("./gui/assets/icons/IcoLinkOff.svg")
var icoLoading, _ = ioutil.ReadFile("./gui/assets/icons/IcoLoading.svg")
var icoAddressBook, _ = ioutil.ReadFile("./gui/assets/icons/IcoAddressBook.svg")
var icoUnconfirmed, _ = ioutil.ReadFile("./gui/assets/icons/IcoUnconfirmed.svg")
var icoTxNumber, _ = ioutil.ReadFile("./gui/assets/icons/IcoTxNumber.svg")

var VIC VIcons = VIcons{

	"logo":        string(icoLogo),
	"overview":    string(icoOverview),
	"send":        string(icoSend),
	"sent":        string(icoSent),
	"receive":     string(icoReceive),
	"received":    string(icoReceived),
	"history":     string(icoHistory),
	"settings":    string(icoSettings),
	"help":        string(icoHelp),
	"info":        string(icoInfo),
	"peers":       string(icoPeers),
	"balance":     string(icoBalance),
	"blocks":      string(icoBlocks),
	"link":        string(icoLink),
	"linkOff":     string(icoLinkOff),
	"loading":     string(icoLoading),
	"addressbook": string(icoAddressBook),
	"unconfirmed": string(icoUnconfirmed),
	"txnumber":    string(icoTxNumber),
}
