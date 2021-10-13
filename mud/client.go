package mud

import (
	"bufio"
	"fmt"
	"net"
)

type Client struct {
	id        int
	conn      net.Conn
	broadcast chan string
	reply     chan string
	commander *CommandHandler
}

// TODO: client id should be an uuid
func NewClient(conn net.Conn, id int, commander *CommandHandler) Client {
	client := Client{
		id:        id,
		conn:      conn,
		broadcast: make(chan string),
		reply:     make(chan string),
		commander: commander,
	}
	return client
}

func (c *Client) Listen(work chan<- ServerAction) {
	reader := bufio.NewReader(c.conn)

	work <- c.commander.ConnectAction(c.id)
	connectReply := <-c.reply
	c.directReply(connectReply)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			// disconnect command tells the server to clean up this client
			// the break here breaks the Listen loop
			work <- c.commander.DisconnectAction(c.id)
			<-c.reply
			break
		}

		work <- c.commander.InputToAction(line, c.id)

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
