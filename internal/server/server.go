package server

import (
	"fmt"
	"net"
	"sync"

	"github.com/mkauppila/mud/internal/game"
)

type Server struct {
	clientsMutex sync.RWMutex
	clients      map[ClientId]*Client
	world        game.Worlder
	idGenerator  IdGenerator
}

func NewServer(idGenerator IdGenerator, world game.Worlder) Server {
	return Server{
		clientsMutex: sync.RWMutex{},
		clients:      make(map[ClientId]*Client),
		world:        world,
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

func (s *Server) StartAcceptingConnections() {
	address := "localhost:6000"
	fmt.Printf("starting at %s\n", address)

	ln, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		go func() {
			err := s.AddNewClient(conn)
			if err != nil {
				panic(err)
			}
		}()
	}
}
