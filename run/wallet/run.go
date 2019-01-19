package wallet

import "encoding/json"

func runNode() {
	j, _ := json.MarshalIndent(Config, "", "  ")
	Log.Tracef.Print("running with configuration:\n%s", string(j))
}
