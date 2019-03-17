package libs

import (
	"io/ioutil"
)

type VImgs map[string][]byte

var mainBg, _ = ioutil.ReadFile("./gui/assets/imgs/mainbg.jpg")
var mainBgS, _ = ioutil.ReadFile("./gui/assets/imgs/mainbgs.jpg")

var VIM VImgs = VImgs{

	"mainbg":  mainBg,
	"mainbgs": mainBgS,
}
