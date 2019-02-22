package gui

import (
	"fmt"
	"net/url"
	"path/filepath"

	"git.parallelcoin.io/pod/cmd/gui/jdb"

	"git.parallelcoin.io/pod/cmd/gui/apps"
	"git.parallelcoin.io/pod/cmd/gui/conf"
	"git.parallelcoin.io/pod/cmd/gui/libs"
	"git.parallelcoin.io/pod/cmd/gui/vue"
	"git.parallelcoin.io/pod/cmd/shell"
	scribble "github.com/nanobox-io/golang-scribble"
	"github.com/zserge/webview"
)

type VDATA struct {
	Config conf.Conf         `json:"config"`
	Pages  map[string]string `json:"pages"`
	Icons  map[string]string `json:"icons"`
	Imgs   map[string][]byte `json:"imgs"`
}

// GUI is the main entry point for the GUI interface
func GUI(
	sh *shell.Config,
) {
	// wlt =
	var err error
	jdb.JDB, err = scribble.New(filepath.Join(sh.DataDir, "gui"), nil)
	if err != nil {
		panic(err)
	}
	w := webview.New(webview.Settings{
		Title:     "ParallelCoin - DUO - True Story",
		Width:     1800,
		Height:    960,
		URL:       `data:text/html,` + url.PathEscape(string(libs.APP["apphtml"])),
		Debug:     true,
		Resizable: false,
	})
	defer w.Exit()
	// Here we need to check for and create wallet :
	// Next start up shell

	//  Now start as normal with something in `wlt`
	apps.InitApps()
	// vue.WLT = _
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
