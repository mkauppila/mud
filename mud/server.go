package mud

import (
	"net"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Server struct {
	actions      chan ServerAction
	clientsMutex sync.Mutex
	clients      map[uuid.UUID]*Client
	registry     *CommandRegistry
	world        *World
	timeStep     time.Duration
}

func NewServer() Server {
	return Server{
		actions:      make(chan ServerAction),
		clientsMutex: sync.Mutex{},
		clients:      make(map[uuid.UUID]*Client),
		registry:     NewLoginCommandRegistry(LoginParseCommand),
		world:        NewWorld(),
		timeStep:     time.Second,
	}
}

func (s *Server) AddNewClient(conn net.Conn) {
	s.clientsMutex.Lock()
	defer s.clientsMutex.Unlock()

	clientId, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}

	client := NewClient(conn, clientId, s.registry)
	s.clients[clientId] = client

	go client.Listen(s.actions)
	go client.Broadcast()
}

func (s *Server) removeClientAtIndex(clientId uuid.UUID) {
	s.clientsMutex.Lock()
	defer s.clientsMutex.Unlock()

	delete(s.clients, clientId)
}

func (s *Server) Run() {
	ticker := time.NewTicker(s.timeStep)
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

	s.world.UpdateCharacterStates(s.timeStep)
}
