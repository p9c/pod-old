package app

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/urfave/cli"
)

func walletHandle(c *cli.Context) error {
	fmt.Println("running wallet")
	if !c.Parent().IsSet("useproxy") {
		*walletConfig.Proxy = ""
	}
	spew.Dump(walletConfig)
	return nil
}
