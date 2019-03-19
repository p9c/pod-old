package app

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"gopkg.in/urfave/cli.v1"
)

func ctlHandleList(c *cli.Context) error {

	fmt.Println("running ctl listcommands")
	_ = ctlHandle(c)
	spew.Dump(ctlConfig)
	return nil
}
