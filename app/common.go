package app

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/urfave/cli.v1"
	"gopkg.in/yaml.v1"
)

func podHandleSave() {
	podconfig := filepath.Join(appConfigCommon.Datadir, podConfigFilename)
	fmt.Println("saving to", podconfig)
	if yp, e := yaml.Marshal(appConfigCommon); e == nil {
		EnsureDir(podconfig)
		if e := ioutil.WriteFile(podconfig, yp, 0600); e != nil {
			panic(e)
		}
	} else {
		panic(e)
	}
}

func confHandle(c *cli.Context) error {
	appendNum := false
	number := 1
	if c.IsSet("number") {
		appendNum = true
		number = c.Int("number")
		if number > 10 {
			return errors.New("cannot make more than 10 (0-9) test profiles")
		}
	}
	base := c.String("base")
	var working string
	fmt.Println("base:", base)
	for i := 0; i < number; i++ {
		working = "" + base
		if appendNum {
			working += fmt.Sprint(i)
		}
		apps := []string{"c", "n", "w", "s", "g"}

		for _, x := range apps {
			e := App.Run([]string{"pod", "-i", "-D", working, x})
			if e != nil {
				panic(e)
			}
		}
	}
	return nil
}
