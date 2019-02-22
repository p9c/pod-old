package main

import (
	"fmt"
	"github.com/obsilp/rmnp"
	"net"
	"time"
)

var (
	myListener = "127.0.0.1:11011"
	server     = rmnp.NewServer(myListener)
	clients    = make(map[string]*rmnp.Client)
	clientsMap = make(map[string]string)
	heartbeat  = time.Second * 3 / 2
	mapWrite   = false
)

func main(
	) {
	server.ClientConnect = clientConnect
	server.ClientDisconnect = clientDisconnect
	server.ClientTimeout = clientTimeout
	server.ClientValidation = validateClient
	server.PacketHandler = handleServerPacket
	server.Start()
	defer server.Stop()
	select {}
}

// Client callbacks

func serverConnect(
	conn *rmnp.Connection, data []byte) {
	// fmt.Println("serverConnect")
	originAddr := conn.Addr.String()
	for conn.Addr != nil {
		for i := range clients {
			returnAddr := clientsMap[originAddr]
			time.Sleep(heartbeat)
			if i == returnAddr {
				fmt.Println("ping", conn.Addr)
				conn.SendReliableOrdered([]byte("ping " + myListener))
				break
			}
		}
	}
}

func serverDisconnect(
	conn *rmnp.Connection, data []byte) {
	fmt.Println("server disconnect")
	addr := clientsMap[conn.Addr.String()]
	for mapWrite {
	}
	mapWrite = true
	delete(clients, clientsMap[addr])
	delete(clientsMap, addr)
	mapWrite = false
	conn.Disconnect([]byte("disconn"))
}

func serverTimeout(
	conn *rmnp.Connection, data []byte) {
	// fmt.Println("server timeout")
	addr := clientsMap[conn.Addr.String()]
	for mapWrite {
	}
	mapWrite = true
	delete(clients, clientsMap[addr])
	delete(clientsMap, addr)
	mapWrite = false
	conn.Disconnect([]byte("timeout"))
}

func handleClientPacket(
	conn *rmnp.Connection, data []byte, channel rmnp.Channel) {
	fmt.Println("->", string(data))
}

// Server callbacks

func clientConnect(
	conn *rmnp.Connection, data []byte) {
	// fmt.Println("clientConnection")
	if string(data) != "kopach" {
		conn.Disconnect([]byte("wrong handshake"))
	}
}

func clientDisconnect(
	conn *rmnp.Connection, data []byte) {
	// fmt.Println("client disconnect")
	addr := clientsMap[conn.Addr.String()]
	for mapWrite {
	}
	mapWrite = true
	delete(clients, clientsMap[addr])
	delete(clientsMap, addr)
	mapWrite = false
	conn.Disconnect([]byte("timeout"))
}

func clientTimeout(
	conn *rmnp.Connection, data []byte) {
	// fmt.Println("client timeout")
	addr := clientsMap[conn.Addr.String()]
	for mapWrite {
	}
	mapWrite = true
	delete(clients, clientsMap[addr])
	delete(clientsMap, addr)
	mapWrite = false
	conn.Disconnect([]byte("timeout"))
}

func validateClient(
	addr *net.UDPAddr, data []byte) (valid bool) {
	// fmt.Println("validateClient")
	valid = string(data) == "kopach"
	if !valid {
		fmt.Println("Wrong handshake from", addr.IP, addr.Port, addr.Network(), addr.Zone)
	}
	return
}

func handleServerPacket(
	conn *rmnp.Connection, data []byte, channel rmnp.Channel) {
	// fmt.Println("handleServerPacket", string(data))
	str := string(data)
	switch {
	case str[:9] == "subscribe":
		addr := string(str[10:])
		if _, ok := clients[addr]; ok {
			// conn.Disconnect([]byte("already subscribed"))
			break
		}
		newClient := rmnp.NewClient(addr)
		newClient.ServerConnect = serverConnect
		newClient.ServerDisconnect = serverDisconnect
		newClient.ServerTimeout = serverTimeout
		newClient.PacketHandler = handleClientPacket
		newClient.ConnectWithData([]byte("nachalnik"))
		for mapWrite {
		}
		mapWrite = true
		clients[addr] = newClient
		clientsMap[newClient.Server.Addr.String()] = addr
		mapWrite = false
		fmt.Println("sub", addr)
		conn.SendReliableOrdered([]byte("subscribed " + addr))
	case string(str[:4]) == "ping":
		conn.SendReliableOrdered([]byte("pong " + myListener))
	}
}
