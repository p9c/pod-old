package main

import (
	"fmt"

	"git.parallelcoin.io/pod/pkg/util/interrupt"
)

func main() {

	interrupt.AddHandler(func() {

		fmt.Println("IT'S THE END OF THE WORLD!")
	})
	<-interrupt.HandlersDone
}
