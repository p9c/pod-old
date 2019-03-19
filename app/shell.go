package app

import "gopkg.in/urfave/cli.v1"

func shellHandleSave() {

	appConfigCommon.Save = false
}
func shellHandle(c *cli.Context) error {

	_ = podHandle(c)

	return nil
}
