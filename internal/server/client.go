package server

import (
	"bufio"
	"fmt"
	"net"

	"github.com/google/uuid"
	"github.com/mkauppila/mud/internal/game"
)

type ClientId string
type IdGenerator func() (ClientId, error)

func UuidGenerator() (ClientId, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return ClientId(uuid.Nil.String()), err
	}
	return ClientId(id.String()), nil
}

type Client struct {
	id        ClientId
	conn      net.Conn
	broadcast chan string
	reply     chan string

	world game.Worlder
}

func NewClient(conn net.Conn, id ClientId, world game.Worlder) *Client {
	client := &Client{
		id:        id,
		conn:      conn,
		broadcast: make(chan string),
		reply:     make(chan string),
		world:     world,
	}

	return client
}

// func DisconnectCommandAction(command Command, clientId ClientId) WorldAction {
// 	return func(s *Server) error {
// 		client := s.getClient(clientId)
// 		if client == nil {
// 			return ErrUnknownClientId{id: clientId}
// 		}
// 		if ch := s.world.GetCharacter(game.ClientId(clientId)); ch != nil {
// 			s.world.RemoveCharacterOnDisconnect(ch)
// 			s.world.BroadcastToOtherCharactersInRoom(
// 				ch,
// 				fmt.Sprintf("%v disconnecting...\n", ch.Name),
// 			)
// 		} else {
// 			return ErrUnknownCharacter{id: clientId, action: command.command}
// 		}
// 		client.Disconnect()
// 		s.removeClient(clientId)
// 		return nil
// 	}
// }

func (c *Client) Listen() {
	reader := bufio.NewReader(c.conn)

	c.world.ClientJoined(
		game.ClientId(c.id),
		func(message string) {
			c.reply <- message
		},
		func(message string) {
			c.broadcast <- message
		},
	)
	c.directReply("Connected to the server...")

	// actions <- c.registry.ConnectAction(c.id)
	// connectReply := <-c.reply
	// c.directReply(connectReply)

	// on connect create the client
	// and put the character in the right mode withing hte
	// these should be only concerned by creating/destroying
	// the tcp client, all the other input syould go try to world

	// Worlder needs another set of fucntions
	// that call server side functionatilies
	// ServerInterface Reply/Broadcast(clientid, message)

	//

	// this should ideally only handle the tcp connection
	// pass the parsed message forward that would handle the rest
	//  that would handle command registry and events

	// let's lift the connect/disconnect from the game all together,
	// it's not really a business for teh game anyway

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			// disconnect command tells the server to clean up this client
			// the break here breaks the Listen loop
			// actions <- c.registry.DisconnectAction(c.id)
			<-c.reply
			break
		}

		c.world.PassMessageToClient(line, game.ClientId(c.id))

		// Wait for client's reply
		commandReply, ok := <-c.reply
		if !ok {
			break
		}
		c.directReply(commandReply)
	}

	fmt.Printf("Client %s disconnected (listen)\n", c.id)
}

func (c *Client) directReply(message string) {
	writer := bufio.NewWriter(c.conn)

	_, err := writer.WriteString(message)
	if err != nil {
		fmt.Println("Failed to write")
	}

	err = writer.Flush()
	if err != nil {
		fmt.Println("Failed to flush")
	}
}

func (c *Client) Broadcast() {
	writer := bufio.NewWriter(c.conn)
	for {
		message, ok := <-c.broadcast
		if !ok {
			break
		}

		_, err := writer.WriteString(message)
		if err != nil {
			fmt.Println("Failed to write")
		}

		err = writer.Flush()
		if err != nil {
			fmt.Println("Failed to flush")
		}
	}

	fmt.Printf("Client %s disconnected (server)\n", c.id)
}

func (c *Client) Disconnect() {
	fmt.Printf("Disconnecting %s\n", c.id)

	close(c.broadcast)
	close(c.reply)

	err := c.conn.Close()
	if err != nil {
		panic(err)
	}
}
