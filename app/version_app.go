package app

import (
	"fmt"

	"github.com/tucnak/climax"
)

// VersionCmd is a command to send RPC queries to bitcoin RPC protocol server for node and wallet queries
var VersionCmd = climax.Command{
	Name:  "version",
	Brief: "prints the version of pod",
	Help:  "",
	Handle: func(ctx climax.Context) int {
		fmt.Println(Version())
		return 0
	},
}
