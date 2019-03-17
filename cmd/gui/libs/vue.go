package libs

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type VPages map[string]string

type aPP map[string][]byte
type jSL map[string][]byte

type vJS map[string][]byte

var APP map[string][]byte = map[string][]byte{

	"apphtml": apphtml,
	"appjs":   appjs,
	"appcss":  appcss,
}

var JSL jSL = jSL{}

var PGS VPages = VPages{}

var VJS vJS = vJS{}
var appcss, _ = ioutil.ReadFile("./gui/assets/vue/app.css")

// var commands, _ = ioutil.ReadFile("./gui/assets/vue/js/commands.js")
// var tasks, _ = ioutil.ReadFile("./gui/assets/vue/js/tasks.js")

// var set, _ = ioutil.ReadFile("./gui/assets/vue/components/settings.js")

// var homePage, _ = ioutil.ReadFile("./gui/assets/vue/components/pages/home.js")
// var sendPage, _ = ioutil.ReadFile("./gui/assets/vue/components/pages/send.js")
// var receivePage, _ = ioutil.ReadFile("./gui/assets/vue/components/pages/receive.js")
// var addressbookPage, _ = ioutil.ReadFile("./gui/assets/vue/components/pages/addressbook.js")
// var historyPage, _ = ioutil.ReadFile("./gui/assets/vue/components/pages/history.js")
// var peersPage, _ = ioutil.ReadFile("./gui/assets/vue/components/pages/peers.js")
// var blocksPage, _ = ioutil.ReadFile("./gui/assets/vue/components/pages/blocks.js")
// var helpPage, _ = ioutil.ReadFile("./gui/assets/vue/components/pages/help.js")
// var consolePage, _ = ioutil.ReadFile("./gui/assets/vue/components/pages/console.js")

// var vrfsendjs, _ = ioutil.ReadFile("./gui/assets/vue/components/modals/vrfsend.js")
var apphtml, _ = ioutil.ReadFile("./gui/assets/vue/app.html")
var appjs, _ = ioutil.ReadFile("./gui/assets/vue/app.js")

func init() {

	// files, err := filepath.Glob("./gui/assets/vue/components/pages/*.js")

	// if err != nil {

	// 	fmt.Print(err)
	// 	os.Exit(1)
	// }

	// for _, f := range files {

	// 	_, fn := filepath.Split(f)
	// 	fl, _ := ioutil.ReadFile(f)
	// 	VJS[strings.TrimSuffix(fn, filepath.Ext(fn))] = fl
	// }

	// appLibs, err := filepath.Glob("./gui/assets/vue/*")

	// if err != nil {

	// 	fmt.Print(err)
	// 	os.Exit(1)
	// }

	// for _, appLib := range appLibs {

	// 	_, fn := filepath.Split(appLib)
	// 	fl, _ := ioutil.ReadFile(appLib)
	// 	PGS[strings.TrimSuffix(fn, filepath.Ext(fn))] = string(fl)
	// }

	vueLibs, err := filepath.Glob("./gui/assets/vue/js/*.js")

	if err != nil {

		fmt.Print(err)
		os.Exit(1)
	}

	for _, vueLib := range vueLibs {

		_, fn := filepath.Split(vueLib)
		fl, _ := ioutil.ReadFile(vueLib)
		JSL[strings.TrimSuffix(fn, filepath.Ext(fn))] = fl
	}

}
