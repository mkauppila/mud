package server

import (
	"net"
	"sync"
	"time"

	"github.com/mkauppila/mud/internal/game"
)

type Server struct {
	clientsMutex sync.RWMutex
	clients      map[ClientId]*Client
	world        *game.World
	timeStep     time.Duration
	idGenerator  IdGenerator
}

func NewServer(idGenerator IdGenerator, world *game.World) Server {
	return Server{
		clientsMutex: sync.RWMutex{},
		clients:      make(map[ClientId]*Client),
		world:        world,
		timeStep:     time.Second,
		idGenerator:  idGenerator,
	}
}

func (s *Server) AddNewClient(conn net.Conn) error {
	clientId, err := s.idGenerator()
	if err != nil {
		return err
	}

	client := NewClient(conn, clientId, s.world)
	s.clientsMutex.Lock()
	s.clients[clientId] = client
	s.clientsMutex.Unlock()

	go client.Listen()
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
	go s.world.RunGameLooop()
}
