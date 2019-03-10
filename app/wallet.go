package app

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"gopkg.in/urfave/cli.v1"
)

func walletHandle(c *cli.Context) error {
	fmt.Println("running wallet")
	if !c.Parent().IsSet("useproxy") {
		*walletConfig.Proxy = ""
	}
	spew.Dump(walletConfig)
	return nil
}
