package main

import (
	"fmt"
	"net"

	"github.com/mkauppila/mud/internal/game"
	"github.com/mkauppila/mud/internal/server"
)

func main() {
	address := "localhost:6000"
	fmt.Printf("starting at %s\n", address)

	ln, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	world := game.NewWorld()
	server := server.NewServer(server.UuidGenerator, world)
	go server.Run()

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		go func() {
			err := server.AddNewClient(conn)
			if err != nil {
				panic(err)
			}
		}()
	}

	// server.Disconnect()
	// for _, client := range clients {
	// 	client.Disconnect()
	// }
}
