package ctl

import "encoding/json"

func runCtl() {
	j, _ := json.MarshalIndent(Config, "", "  ")
	Log.Tracef.Print("running with configuration:\n%s", string(j))
}
