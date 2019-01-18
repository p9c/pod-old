package shell

import "encoding/json"

func runShell() {
	j, _ := json.MarshalIndent(CombinedCfg, "", "  ")
	Log.Tracef.Print("running with configuration:\n%s", string(j))
}
