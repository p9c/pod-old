package main

import (
	"fmt"
	"github.com/obsilp/rmnp"
	"net"
	"time"
)

var (
	serverAddr = "127.0.0.1:11011"
	myListener = "127.0.0.1:11012"
	server     = rmnp.NewServer(myListener)
	client     *rmnp.Client
	connected  = false
	subscribed = false
	heartbeat  = time.Second * 3 / 2
)

func main() {
	server.ClientConnect = clientConnect
	server.ClientDisconnect = clientDisconnect
	server.ClientTimeout = clientTimeout
	server.ClientValidation = validateClient
	server.PacketHandler = handleServerPacket
	server.Start()
	defer server.Stop()
	for {
		if !connected {
			subscribed = false
			connected = true
			client = rmnp.NewClient(serverAddr)
			client.ServerConnect = serverConnect
			client.ServerDisconnect = serverDisconnect
			client.ServerTimeout = serverTimeout
			client.PacketHandler = handleClientPacket
			client.ConnectWithData([]byte("kopach"))
		} else {
			time.Sleep(heartbeat)
		}
	}
}

// Client callbacks

func serverConnect(conn *rmnp.Connection, data []byte) {
	// fmt.Println("serverConnect")
	if !subscribed && connected {
		fmt.Println("subscribe", serverAddr)
		conn.SendReliableOrdered([]byte("subscribe " + myListener))
		time.Sleep(heartbeat)
	}
	for conn.Addr != nil && connected {
		fmt.Println("ping", conn.Addr)
		conn.SendReliableOrdered([]byte("ping " + myListener))
		time.Sleep(heartbeat)
	}
}

func serverDisconnect(conn *rmnp.Connection, data []byte) {
	// fmt.Println("server disconnect")
	subscribed = false
	connected = false
	conn.Disconnect([]byte("disconn"))
}

func serverTimeout(conn *rmnp.Connection, data []byte) {
	// fmt.Println("server timeout")
	subscribed = false
	connected = false
	conn.Disconnect([]byte("timeout"))
}

func handleClientPacket(conn *rmnp.Connection, data []byte, channel rmnp.Channel) {
	fmt.Println("->" + string(data))
	if string(data)[:10] == "subscribed" {
		subscribed = true
	}
}

// Client callbacks

func clientConnect(conn *rmnp.Connection, data []byte) {
	// fmt.Println("clientConnection")
	if string(data) != "nachalnik" {
		conn.Disconnect([]byte("wrong handshake"))
	}
	subscribed = true
	connected = true
}

func clientDisconnect(conn *rmnp.Connection, data []byte) {
	// fmt.Println("client disconnect")
	subscribed = false
	connected = false
}

func clientTimeout(conn *rmnp.Connection, data []byte) {
	// fmt.Println("client timeout")
	subscribed = false
	connected = false
}

func validateClient(addr *net.UDPAddr, data []byte) (valid bool) {
	// fmt.Println("validateClient")
	valid = string(data) == "nachalnik"
	if !valid {
		fmt.Println("Wrong handshake from", addr.IP, addr.Port, addr.Network(), addr.Zone)
	}
	return
}

func handleServerPacket(conn *rmnp.Connection, data []byte, channel rmnp.Channel) {
	// fmt.Println("handleServerPacket", string(data))
	str := string(data)
	switch {
	case !connected:
		conn.Disconnect([]byte("already connected"))
	case str[:10] == "subscribed":
		subscribed = true
	case str == "already subscribed" && connected:
		subscribed = true
	case str[:4] == "ping" && subscribed && connected:
		conn.SendReliableOrdered([]byte("pong " + str[5:]))
	}
}
