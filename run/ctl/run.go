package ctl

import "encoding/json"

func runCtl() {
	j, _ := json.MarshalIndent(Config, "", "  ")
	log.Tracef.Print("running with configuration:\n%s", string(j))
}
