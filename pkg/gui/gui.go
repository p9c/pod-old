package gui

import (
	"fmt"
	"net/url"

	"git.parallelcoin.io/pod/pkg/gui/apps"
	"git.parallelcoin.io/pod/pkg/gui/conf"
	"git.parallelcoin.io/pod/pkg/gui/libs"
	"git.parallelcoin.io/pod/pkg/gui/vue"
	"git.parallelcoin.io/pod/pkg/wallet"
	"github.com/zserge/webview"
)

type VDATA struct {
	Config conf.Conf         `json:"config"`
	Pages  map[string]string `json:"pages"`
	Icons  map[string]string `json:"icons"`
	Imgs   map[string][]byte `json:"imgs"`
}

var wlt *wallet.Wallet

// GUI is the main entry point for the GUI interface
func GUI() {
	w := webview.New(webview.Settings{
		Title:     "ParallelCoin - DUO - True Story",
		Width:     1800,
		Height:    960,
		URL:       `data:text/html,` + url.PathEscape(string(libs.APP["apphtml"])),
		Debug:     true,
		Resizable: false,
	})
	defer w.Exit()
	// Here we need to check for and create wallet

	// Next start up shell

	//  Now start as normal with something in `wlt`
	apps.InitApps()
	vue.WLT = wlt
	w.Dispatch(func() {

		w.Bind("blockchaindata", &vue.BlockChain{})
		// w.Bind("sendtoaddress", &vue.SendToAddress{})
		w.Bind("language", &vue.Language{})
		w.Bind("addressbook", &apps.AddressBook{})
		w.Bind("addressbooklabel", &apps.AddressBookLabel{})
		for mn, md := range vue.MODS {
			w.Bind(mn, &md)
			fmt.Println("aaaaaaaaaaaaaaaaaaaaaaaaaaaaawwwwwwww", mn)
		}
		fmt.Println("vue.MODSvue.MODSvue.MODSvue.MODSvue.MODSvue.MODSvue.MODSvue.MODSvue.MODSvue.MODSvue.MODSvue.MODSvue.MODSvue.MODS", vue.MODS)

		w.Bind("reqpays", &vue.RequestedPaymentHistory{})
		w.Bind("reqpay", &vue.RequestedPayment{})

		w.Bind("rpcinterface", &vue.RPCInterface{})

		w.Bind("conf", &conf.Conf{})

		w.Bind("vuedata", &VDATA{
			Config: conf.VCF,
			Pages:  libs.PGS,
			Icons:  libs.VIC,
			Imgs:   libs.VIM,
		})

		for _, c := range libs.CSS {
			w.InjectCSS(string(c))
		}

		for _, j := range libs.JSL {
			w.Eval(string(j))
		}
		for _, v := range libs.VJS {
			w.Eval(string(v))
		}

		w.Eval(string(libs.APP["appjs"]))
		w.InjectCSS(string(libs.APP["appcss"]))

	})

	w.Run()
}
