package gui

import (
	"fmt"
	"net/url"

	"github.com/parallelcointeam/mod/gui/apps"
	"github.com/parallelcointeam/mod/gui/conf"
	"github.com/parallelcointeam/mod/gui/libs"
	"github.com/parallelcointeam/mod/gui/vue"
	"github.com/parallelcointeam/mod/wallet"
	"github.com/zserge/webview"
)

type VDATA struct {
	Config conf.Conf         `json:"config"`
	Pages  map[string]string `json:"pages"`
	Icons  map[string]string `json:"icons"`
	Imgs   map[string][]byte `json:"imgs"`
}

func GUI(wlt *wallet.Wallet) {
	apps.InitApps()
	vue.WLT = wlt
	w := webview.New(webview.Settings{
		Title:     "ParallelCoin - DUO - True Story",
		Width:     1800,
		Height:    960,
		URL:       `data:text/html,` + url.PathEscape(string(libs.APP["apphtml"])),
		Debug:     true,
		Resizable: false,
	})
	defer w.Exit()
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
