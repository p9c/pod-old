package app

import (
	"encoding/json"
	"time"

	cl "git.parallelcoin.io/pod/pkg/clog"
)

func runShell() {
	j, _ := json.MarshalIndent(ShellConfig, "", "  ")
	log <- cl.Tracef{"running with configuration:\n%s", string(j)}
	time.Sleep(time.Second)
}
