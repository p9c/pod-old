package libs

import (
	"io/ioutil"
)

type cSS map[string][]byte

var buefycss, _ = ioutil.ReadFile("./gui/assets/vue/css/buefy.css")
var vuelayerscss, _ = ioutil.ReadFile("./gui/assets/vue/css/vuelayers.css")
var scrollercss, _ = ioutil.ReadFile("./gui/assets/vue/js/scroller.css")

var fontawesome, _ = ioutil.ReadFile("./gui/assets/vue/css/fontawesome.css")
var materialdesignicons, _ = ioutil.ReadFile("./gui/assets/vue/css/materialdesignicons.css")

var CSS cSS = cSS{

	"fontawesome":         fontawesome,
	"materialdesignicons": materialdesignicons,
	"buefycss":            buefycss,
	"vuelayerscss":        vuelayerscss,
	"scrollercss":         scrollercss,
}
