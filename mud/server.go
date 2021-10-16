package mud

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type Server struct {
	actions      chan ServerAction
	clientsMutex sync.Mutex
	clients      []Client
	commander    *CommandHandler
}

func NewServer() Server {
	return Server{
		actions:      make(chan ServerAction),
		clientsMutex: sync.Mutex{},
		clients:      []Client{},
		commander:    NewCommandHandler(),
	}
}

func (s *Server) AddNewClient(conn net.Conn) {
	s.clientsMutex.Lock()
	defer s.clientsMutex.Unlock()

	client := NewClient(conn, len(s.clients), s.commander)
	s.clients = append(s.clients, client)

	go client.Listen(s.actions)
	go client.Broadcast()
}

func (s *Server) removeClientAtIndex(index int) {
	s.clientsMutex.Lock()
	defer s.clientsMutex.Unlock()

	s.clients[index] = s.clients[len(s.clients)-1]
	s.clients = s.clients[:len(s.clients)-1]
}

func (s *Server) Run() {
	ticker := time.NewTicker(time.Duration(time.Millisecond * 1000))
	defer ticker.Stop()

	var actions []ServerAction
	for {
		select {
		case command, ok := <-s.actions:
			if !ok {
				panic("actions channel closed")
			}
			actions = append(actions, command)
		case _, ok := <-ticker.C:
			if !ok {
				panic("ticker closed")
			}

			s.processServerActions(actions)
			actions = make([]ServerAction, 0)
		}
	}
}

func (s *Server) processServerActions(actions []ServerAction) {
	for _, action := range actions {
		err := action(s)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("All actions processed for this tick")
}
