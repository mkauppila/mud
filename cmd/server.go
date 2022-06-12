package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"

	_ "net/http/pprof"

	"github.com/mkauppila/mud/internal/game"
	"github.com/mkauppila/mud/internal/server"
)


func setupPprof(host string, port int16) {
	fmt.Printf("Start pprof at %s:%d\n", host, port)

	// pprof is by default added to DefaultServerMux
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), nil)
	// err := http.ListenAndServe("localhost:5000", nil)
	if err != nil {
		panic(err)
	}
}

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

	go setupPprof("localhost", 8000)

	<-exitC

	// server.Disconnect()
	// for _, client := range clients {
	// 	client.Disconnect()
	// }
}
