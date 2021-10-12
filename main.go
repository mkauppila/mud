package main

import (
	"fmt"
	"net"

	mud "github.com/mkauppila/mud/mud"
)

func main() {
	address := "localhost:6000"
	fmt.Printf("starting at %s\n", address)

	ln, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	server := mud.NewServer()
	go server.Run()

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		server.AddNewConnection(conn)
	}

	// server.Disconnect()
	// for _, client := range clients {
	// 	client.Disconnect()
	// }
}
