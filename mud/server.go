package mud

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type Server struct {
	work         chan ServerAction
	clientsMutex sync.Mutex
	clients      []Client
	commander    *CommandHandler
}

func NewServer() Server {
	return Server{
		work:         make(chan ServerAction),
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

	go client.Listen(s.work)
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
