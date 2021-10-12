package main

import (
	"fmt"
	"net"
	"time"
)

type Server struct {
	work      chan ServerAction
	clients   []*Client // might not need to be a pointer, slices are pointer-like?
	commander *CommandHandler
}

func NewServer() Server {
	return Server{
		work:      make(chan ServerAction),
		clients:   []*Client{},
		commander: NewCommandHandler(),
	}
}

func (s *Server) AddNewConnection(conn net.Conn) {
	client := NewClient(conn, len(s.clients), s.commander)
	go client.Listen(s.work)
	go client.Broadcast()
	s.clients = append(s.clients, client)
}

func (s *Server) Run() {
	ticker := time.NewTicker(time.Duration(time.Millisecond * 1000))
	defer ticker.Stop()

	var actions []ServerAction
	for {
		select {
		case command, ok := <-s.work:
			if !ok {
				panic("work queue closed")
			}
			actions = append(actions, command)
		case _, ok := <-ticker.C:
			if !ok {
				panic("ticker closed")
			}

			s.processCommands(actions)
			actions = make([]ServerAction, 0)
		}
	}
}

func (s *Server) processCommands(actions []ServerAction) {
	for _, action := range actions {
		err := action(s)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("All actions processed for this tick")
}
