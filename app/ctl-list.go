package app

import (
	"fmt"
	"time"

	"git.parallelcoin.io/dev/pod/cmd/ctl"
	"gopkg.in/urfave/cli.v1"
)

func ctlHandleList(c *cli.Context) error {
	fmt.Println("Here are the available commands. Pausing a moment as it is a long list...")
	time.Sleep(2 * time.Second)
	ctl.ListCommands()
	return nil
}
