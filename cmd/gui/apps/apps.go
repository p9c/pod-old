package apps

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"git.parallelcoin.io/dev/pod/cmd/gui/libs"
)

func InitApps() {

	// JS
	jsLibs, err := filepath.Glob("./gui/apps/*/js/*.js")

	if err != nil {

		fmt.Print(err)
		os.Exit(1)
	}

	for _, jsLib := range jsLibs {

		_, fn := filepath.Split(jsLib)
		fl, _ := ioutil.ReadFile(jsLib)
		libs.VJS[strings.TrimSuffix(fn, filepath.Ext(fn))] = fl
	}

	// jsLibs, err := filepath.Glob("./apps/*/js/*.js")

	// if err != nil {

	// 	fmt.Print(err)
	// 	os.Exit(1)
	// }

	// for _, jsLib := range jsLibs {

	// 	_, fn := filepath.Split(jsLib)
	// 	fl, _ := ioutil.ReadFile(jsLib)
	// 	libs.VJS[strings.TrimSuffix(fn, filepath.Ext(fn))] = fl
	// }

	// HTML
	libHTMLs, err := filepath.Glob("./gui/apps/*/html/*.html")

	if err != nil {

		fmt.Print(err)
		os.Exit(1)
	}

	for _, libHTML := range libHTMLs {

		_, fn := filepath.Split(libHTML)
		fl, _ := ioutil.ReadFile(libHTML)
		libs.PGS[strings.TrimSuffix(fn, filepath.Ext(fn))] = string(fl)
	}

}
