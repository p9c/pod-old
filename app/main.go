package app

import (
	"github.com/urfave/cli"
	"os"
)



var App = cli.App{
	Name:        "pod",
	Version:     "v0.9.9",
	Description: "Parallelcoin Pod Suite -- All-in-one everything for Parallelcoin!",
	Copyright:   "Legacy portions derived from btcsuite/btcd under ISC licence. The remainder is already in your possession. Use it wisely.",
}

func Main() int {
	e := App.Run(os.Args)
	if e != nil {
		return 1
	}
	return 0
}
