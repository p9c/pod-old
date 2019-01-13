package gui

import (
	"net/url"

	"github.com/parallelcointeam/mod/gui/jdb"
	"github.com/parallelcointeam/mod/vue"
	"github.com/zserge/webview"
)

type VDATA struct {
	Pages map[string]string `json:"pages"`
	Icons map[string]string `json:"icons"`
	Imgs  map[string][]byte `json:"imgs"`
}

func GUI() {
	// libs := jdb.VueLibs
	// pages := jdb.VuePages
	w := webview.New(webview.Settings{
		Title:     "ParallelCoin - DUO - True Story",
		Width:     1400,
		Height:    800,
		URL:       `data:text/html,` + url.PathEscape(string(jdb.VLB["apphtml"])),
		Debug:     true,
		Resizable: false,
	})
	defer w.Exit()
	w.Dispatch(func() {
		// w.Bind("blockchaindata", []interface{}{(*btcjson.InfoWalletResult)(nil)})
		w.Bind("blockchaindata", &vue.BlockChainData{})

		//w.Bind("icons", &icons)
		w.Bind("vuedata", &VDATA{
			Pages: jdb.VPG,
			Icons: jdb.VIC,
			Imgs:  jdb.VIM,
		})

		w.InjectCSS(string(jdb.VLB["buefycss"]))
		w.InjectCSS(string(jdb.VLB["appcss"]))

		w.Eval(string(jdb.VLB["vue"]))
		w.Eval(string(jdb.VLB["easybar"]))
		w.Eval(string(jdb.VLB["buefyjs"]))

		w.Eval(string(jdb.VLB["settings"]))
		w.Eval(string(jdb.VLB["comp"]))

		// w.Eval(string(jdb.VPG["home"]))

		w.Eval(string(jdb.VLB["appjs"]))
		// fmt.Println("daaaaaa", imgs)
	})

	w.Run()
}
