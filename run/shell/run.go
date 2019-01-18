package shell

import "encoding/json"

func runShell() {
	j, _ := json.MarshalIndent(Config, "", "  ")
	Log.Tracef.Print("running with configuration:\n%s", string(j))
}
