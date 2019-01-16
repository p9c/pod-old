package main

import (
	"errors"
	"fmt"
)

var factory factoryReset

func (n *factoryReset) Execute(args []string) (err error) {
	fmt.Println("pod", version())
	if n.OnlyPurge && n.Purge {
		return errors.New("conflicting purge and onlypurge options")
	}
	fmt.Print("resetting configurations to factory ")
	if !n.Really {
		fmt.Println("... just kidding, must add `--really` to commandline options to confirm")
		return errors.New("not confirmed")
	}
	if n.Purge {
		fmt.Print("... also deleting data")
	}
	fmt.Println()
	flip := false
	if n.Ctl {
		if n.OnlyPurge {
			fmt.Println("purging data for ctl")
		} else {
			fmt.Println("resetting settings for ctl")
			if n.Purge {
				fmt.Println("purging data for ctl")
			}
		}
		flip = true
	}
	if n.Node {
		if n.OnlyPurge {
			fmt.Println("purging data for node")
		} else {
			fmt.Println("resetting settings for node")
			if n.Purge {
				fmt.Println("purging data for node")
			}
		}
		flip = true
	}
	if n.Wallet {
		if n.OnlyPurge {
			fmt.Println("purging data for wallet")
		} else {
			fmt.Println("resetting settings for wallet")
			if n.Purge {
				fmt.Println("purging data for wallet")
			}
		}
		flip = true
	}
	if n.WalletGUI {
		if n.OnlyPurge {
			fmt.Println("purging data for walletgui")
		} else {
			fmt.Println("resetting settings for walletgui")
			if n.Purge {
				fmt.Println("purging data for walletgui")
			}
		}
		flip = true
	}
	if n.WalletNode {
		if n.OnlyPurge {
			fmt.Println("purging data for walletnode")
		} else {
			fmt.Println("resetting settings for walletnode")
			if n.Purge {
				fmt.Println("purging data for walletnode")
			}
		}
		flip = true
	}
	if !flip {
		if n.OnlyPurge {
			fmt.Println("purging data for all subcommands")
		} else {
			fmt.Println("resetting all settings for all subcommands")
			if n.Purge {
				fmt.Println("purging data for all subcommands")
			}
		}
	}
	return
}
