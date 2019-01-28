package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	gomod, err := ioutil.ReadFile("go.mod")
	if err != nil {
		log.Fatal(err)
	}
	s := string(gomod)
	S := strings.Split(s, "\n")
	S = strings.Split(S[0], " ")
	args := `-X `
	args += S[1]
	args += `/app.Stamp=`
	args += time.Now().UTC().Format("v06.01.02.15")

	var dir string
	dir, err = os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	dir += "/bin/pod"
	cmd := exec.Command("go", "build", "-o", dir, "-ldflags", args, "-v")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		fmt.Println("ERR", err)
		os.Exit(1)
	}
}
