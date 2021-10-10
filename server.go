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
	go client.Broadcast()
	s.clients = append(s.clients, client)
}

func (s *Server) Run() {
	ticker := time.NewTicker(time.Duration(time.Millisecond * 1000))
	defer ticker.Stop()

	var commands []Command
	for {
		select {
		case command, ok := <-s.work:
			if !ok {
				panic("work queue closed")
			}
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
		fmt.Println("process command: ", msg)
		switch msg.command {
		case "say":
			for _, c := range s.clients {
				if c.id == msg.clientId {
					c.reply <- fmt.Sprintf("You said %s\n", msg.contents)
				} else {
					c.broadcast <- fmt.Sprintf("They said %s\n", msg.contents)
				}
			}
		case "connect":
			for _, c := range s.clients {
				if c.id == msg.clientId {
					c.reply <- "You joined\n"
				} else {
					c.broadcast <- fmt.Sprintf("Client %d joined!\n", msg.clientId)
				}
			}
		case "disconnect":
			for i, c := range s.clients {
				if c.id == msg.clientId {
					c.reply <- "You are disconnecting"
					c.Disconnect()
					// update the clients list
					s.clients[i] = s.clients[len(s.clients)-1]
					s.clients = s.clients[:len(s.clients)-1]
				} else {
					c.broadcast <- fmt.Sprintf("Client %d disconnecting...\n", msg.clientId)
				}
			}
		}
	}

	fmt.Println("All commands processed for this tick")
}
