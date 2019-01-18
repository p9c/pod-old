package shell

import "encoding/json"

func runNode() {
	j, _ := json.MarshalIndent(CombinedCfg, "", "  ")
	log.Tracef.Print("running with configuration:\n%s", string(j))
}
