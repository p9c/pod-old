package app

import (
	"fmt"

	walletmain "git.parallelcoin.io/dev/pod/cmd/walletmain"

	"gopkg.in/urfave/cli.v1"
)

func walletHandle(c *cli.Context) error {
	fmt.Println("starting wallet")
	Configure()
	return walletmain.Main(&podConfig, activeNetParams)
}
