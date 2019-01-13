package jdb

import (
	"io/ioutil"
)

type VPages map[string]string

var app, _ = ioutil.ReadFile("./gui/assets/pages/app.html")
var home, _ = ioutil.ReadFile("./gui/assets/pages/home.html")
var send, _ = ioutil.ReadFile("./gui/assets/pages/send.html")
var receive, _ = ioutil.ReadFile("./gui/assets/pages/receive.html")
var history, _ = ioutil.ReadFile("./gui/assets/pages/history.html")
var addressbook, _ = ioutil.ReadFile("./gui/assets/pages/addressbook.html")

// var settings, _ = ioutil.ReadFile("./gui/assets/pages/settings.html")
var peers, _ = ioutil.ReadFile("./gui/assets/pages/peers.html")
var blocks, _ = ioutil.ReadFile("./gui/assets/pages/blocks.html")
var about, _ = ioutil.ReadFile("./gui/assets/pages/about.html")
var help, _ = ioutil.ReadFile("./gui/assets/pages/help.html")

var settings, _ = ioutil.ReadFile("./gui/assets/pages/settings/settings.html")
var ifc, _ = ioutil.ReadFile("./gui/assets/pages/settings/interface.html")
var network, _ = ioutil.ReadFile("./gui/assets/pages/settings/network.html")
var security, _ = ioutil.ReadFile("./gui/assets/pages/settings/security.html")
var mining, _ = ioutil.ReadFile("./gui/assets/pages/settings/mining.html")

var VPG VPages = VPages{
	"app":         string(app),
	"home":        string(home),
	"send":        string(send),
	"receive":     string(receive),
	"history":     string(history),
	"addressbook": string(addressbook),
	"peers":       string(peers),
	"blocks":      string(blocks),
	"about":       string(about),
	"help":        string(help),

	"settings": string(settings),
	"ifc":      string(ifc),
	"network":  string(network),
	"security": string(security),
	"mining":   string(mining),
}

// func init() {
// 	fmt.Println("daaaaaa", string(home))
// }
