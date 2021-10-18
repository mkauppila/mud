package mud

import (
	"bufio"
	"fmt"
	"net"

	"github.com/google/uuid"
)

type Client struct {
	id        uuid.UUID
	conn      net.Conn
	broadcast chan string
	reply     chan string
	registry  *CommandRegistry

	Character *Character
}

func NewClient(conn net.Conn, id uuid.UUID, registry *CommandRegistry, character *Character) Client {
	client := Client{
		id:        id,
		conn:      conn,
		broadcast: make(chan string),
		reply:     make(chan string),
		registry:  registry,
		Character: character,
	}

	client.Character.Reply = func(message string) {
		client.reply <- message
	}
	client.Character.Broadcast = func(message string) {
		client.broadcast <- message
	}

	return client
}

func (c *Client) Listen(actions chan<- ServerAction) {
	reader := bufio.NewReader(c.conn)

	actions <- c.registry.ConnectAction(c.id)
	connectReply := <-c.reply
	c.directReply(connectReply)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			// disconnect command tells the server to clean up this client
			// the break here breaks the Listen loop
			actions <- c.registry.DisconnectAction(c.id)
			<-c.reply
			break
		}

		actions <- c.registry.InputToAction(line, c.id)

		// Wait for player's reply
		commandReply, ok := <-c.reply
		if !ok {
			break
		}
		c.directReply(commandReply)
	}

	fmt.Printf("Client %d disconnected (listen)\n", c.id)
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

	fmt.Printf("Client %d disconnected (server)\n", c.id)
}

func (c *Client) Disconnect() {
	fmt.Printf("Disconnecting %d\n", c.id)

	close(c.broadcast)
	close(c.reply)

	err := c.conn.Close()
	if err != nil {
		panic(err)
	}
}
