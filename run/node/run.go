package node

import "encoding/json"

func runNode() {
	j, _ := json.MarshalIndent(Config, "", "  ")
	log.Tracef.Print("running with configuration:\n%s", string(j))
}
