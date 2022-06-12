package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/mkauppila/mud/internal/game"
	"github.com/mkauppila/mud/internal/server"
)

func main() {
	exitC := make(chan struct{})

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)
	go func() {
		for sig := range signals {
			fmt.Println("Going to shut down due to: ", sig)
			exitC <- struct{}{}
		}
	}()

	world := game.NewWorld()
	server := server.NewServer(server.UuidGenerator, world)
	go server.StartAcceptingConnections()
	go world.RunGameLoop()

	<-exitC

	// server.Disconnect()
	// for _, client := range clients {
	// 	client.Disconnect()
	// }
}
