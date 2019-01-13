package jdb

import (
	"fmt"

	scribble "github.com/nanobox-io/golang-scribble"
)

var VueLibs map[string][]byte

// var VuePages map[string]string
var VueIcons map[string]string
var VueImgs map[string]string

func init() {
	dir := "./gui/jdb/"
	db, err := scribble.New(dir, nil)
	if err != nil {
		fmt.Println("Error", err)
	}
	ICF := InfConf{
		Lang:  "en",
		Deno:  "min",
		Fiat:  "min",
		Theme: "min",
		CCSS:  "min",
		Start: "min",
		Tray:  true,
	}
	NCF := NetConf{
		TLS:     true,
		Network: "network",
		RPC:     "rpc",
		SRPC:    "srpc",
		TLSpub:  "tlspub",
		TLSpri:  "tlspri",
		Proxy:   "rpc",
	}
	SCF := SecConf{
		Network: "network",
	}
	MCF := MiningConf{
		Algo:  "network",
		CPU:   "network",
		Cores: 6,
	}
	db.Write("conf", "interface", ICF)
	db.Write("conf", "network", NCF)
	db.Write("conf", "security", SCF)
	db.Write("conf", "mining", MCF)

	// db.Write("data", "vlibs", VLB)
	// db.Write("data", "vpages", VPG)
	// db.Write("data", "vicons", VIC)
	// db.Write("data", "vimgs", VIM)

	if err := db.Read("data", "vlibs", &VueLibs); err != nil {
		fmt.Println("Error", err)
	}
	// if err := db.Read("data", "vpages", &VuePages); err != nil {
	// 	fmt.Println("Error", err)
	// }
	if err := db.Read("data", "vicons", &VueIcons); err != nil {
		fmt.Println("Error", err)
	}
	if err := db.Read("data", "vimgs", &VueImgs); err != nil {
		fmt.Println("Error", err)
	}
	// if err := db.Read("data", "vicons", &VueIcons); err != nil {
	// 	fmt.Println("Error", err)
	// }

	// if err := db.Read("vicons", "vicons", &VIcons); err != nil {
	// 	fmt.Println("Error", err)
	// }
	// if err := db.Read("data", "settings", &Settings); err != nil {
	// 	fmt.Println("Error", err)
	// }

}
