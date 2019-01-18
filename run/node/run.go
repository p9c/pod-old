package node

import "encoding/json"

func runNode() {
	j, _ := json.MarshalIndent(CombinedCfg, "", "  ")
	Log.Tracef.Print("running with configuration:\n%s", string(j))
}
