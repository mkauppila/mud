package server

import (
	"net"
	"sync"
	"time"

	"github.com/mkauppila/mud/internal/game"
)

type Server struct {
	actions      chan ServerAction
	clientsMutex sync.RWMutex
	clients      map[ClientId]*Client
	world        *game.World
	timeStep     time.Duration
	idGenerator  IdGenerator
}

func NewServer(idGenerator IdGenerator) Server {
	return Server{
		actions:      make(chan ServerAction),
		clientsMutex: sync.RWMutex{},
		clients:      make(map[ClientId]*Client),
		world:        game.NewWorld(),
		timeStep:     time.Second,
		idGenerator:  idGenerator,
	}
}

func (s *Server) AddNewClient(conn net.Conn) error {
	clientId, err := s.idGenerator()
	if err != nil {
		return err
	}

	client := NewClient(conn, clientId, NewLoginCommandRegistry())
	s.clientsMutex.Lock()
	s.clients[clientId] = client
	s.clientsMutex.Unlock()

	go client.Listen(s.actions)
	go client.Broadcast()

	return nil
}

func (s *Server) removeClient(id ClientId) {
	s.clientsMutex.Lock()
	delete(s.clients, id)
	s.clientsMutex.Unlock()
}

func (s *Server) getClient(id ClientId) *Client {
	s.clientsMutex.RLock()
	defer s.clientsMutex.RUnlock()
	client, ok := s.clients[id]
	if ok {
		return client
	} else {
		return nil
	}
}

func (s *Server) Run() {
	ticker := time.NewTicker(s.timeStep)
	defer ticker.Stop()

	// lets move to this world
	// it can handle all the events
	// a couple of them connect/disconnect need to be augmented by
	// the server though

	// world will be contact between the game and the server
	// logic

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

	s.world.UpdateCharacterStates(s.timeStep)
}
