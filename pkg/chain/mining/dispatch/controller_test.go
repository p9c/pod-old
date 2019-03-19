package controller

import (
	"fmt"
	"testing"
	"time"

	"github.com/dogeerf/rpcx"
	"github.com/xtaci/kcp-go"
	"golang.org/x/crypto/argon2"
)

const cryptKey = "rpcx-key"
const cryptSalt = "rpcx-salt"

var minerAddr = "127.0.0.1:11012"
var nodeAddr = "127.0.0.1:11011"

var ready = make(chan bool)
var K Kopach

func TestKCPClientServer(
	t *testing.T) {

	fmt.Println("Testing KCP client/server connection")
	go runServer(t)
	time.Sleep(time.Second)
	runClient(t)
}
func runServer(
	t *testing.T) {

	server := rpcx.NewServer()
	server.RegisterName("Kopach", &K)
	pass := argon2.IDKey([]byte(cryptKey), []byte(cryptSalt), 1, 4096, 32, 32)
	bc, err := kcp.NewAESBlockCrypt(pass)
	if err != nil {

		fmt.Println(err)
	}
	ln, err := kcp.ListenWithOptions(nodeAddr, bc, 10, 3)
	if err != nil {

		fmt.Println(err)
	}
	fmt.Println("Running server")
	server.ServeListener(ln)
	ready <- true
	fmt.Println("Finished running server")
}
func runClient(
	t *testing.T) {

	pass := argon2.Key([]byte(cryptKey), []byte(cryptSalt), 1, 4096, 32, 32)
	bc, _ := kcp.NewAESBlockCrypt(pass)

	s := &rpcx.DirectClientSelector{Network: "kcp", Address: nodeAddr}
	client := rpcx.NewClient(s)
	client.Block = bc
	args := &minerAddr
	var reply *string
	fmt.Println("server", nodeAddr, "Calling subscribe", *args, client.PluginContainer)
	err := client.Call("Kopach.Subscribe", &args, &reply)
	fmt.Println("Got reply", reply)
	if err != nil {

		fmt.Printf("error for Kopach: %s, %v\n", *args, err)
	} else {

		fmt.Printf("Subscribe: %s : %s\n", *args, *reply)
	}
	fmt.Println("closing connection")
	client.Close()
}
