package main

import (
	"fmt"
	"net"
	"time"
)

type Server struct {
	work    chan Command
	clients []*Client
}

func NewServer() Server {
	return Server{
		work:    make(chan Command),
		clients: []*Client{},
	}
}

func (s *Server) AddNewConnection(conn net.Conn) {
	client := NewClient(conn, len(s.clients))
	go client.Listen(s.work)
	go client.Serve()
	s.clients = append(s.clients, client)
}

func (s *Server) Run() {
	// basically, read messages until a tick completes,
	// then run them all and start from the beginning
	ticker := time.NewTicker(time.Duration(time.Millisecond * 1000))
	defer ticker.Stop()

	var commands []Command
	for {
		select {
		case command, ok := <-s.work:
			if !ok {
				panic("work queue closed")
			}
			// but handle disconnects immediatley!
			// or as the first thing during the message handling!
			commands = append(commands, command)
		case _, ok := <-ticker.C:
			if !ok {
				panic("ticker closed")
			}

			s.processCommands(commands)
			commands = make([]Command, 0)
		}
	}
}

func (s *Server) processCommands(commands []Command) {
	fmt.Println("process commands: ", len(commands))
	for _, msg := range commands {
		switch msg.command {
		case "say":
			for _, c := range s.clients {
				c.outgoing <- msg.contents
			}
		case "connect":
			for _, c := range s.clients {
				if c.id != msg.clientId {
					c.outgoing <- fmt.Sprintf("Client %d joined!\n", msg.clientId)
				}
			}
		case "disconnect":
			for i, c := range s.clients {
				if c.id == msg.clientId {
					c.Disconnect()
					// update the clients list
					s.clients[i] = s.clients[len(s.clients)-1]
					s.clients = s.clients[:len(s.clients)-1]
				} else {
					c.outgoing <- fmt.Sprintf("Client %d disconnecting...\n", msg.clientId)
				}
			}
		}
	}
}
