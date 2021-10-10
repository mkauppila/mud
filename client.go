package main

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
}

func NewClient(conn net.Conn, id int) *Client {
	client := Client{
		id:        id,
		conn:      conn,
		broadcast: make(chan string),
		reply:     make(chan string),
	}
	return &client
}

// work: send parsed and formed message from the
// the clent to the server to be processed? Needs to which
// client send it and when (server tick will be handled separately)

func (c *Client) Listen(work chan<- Command) {
	reader := bufio.NewReader(c.conn)

	work <- Command{command: "connect", contents: "", clientId: c.id}
	connectReply := <-c.reply
	c.directReply(connectReply)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			// disconnect command tells the server to clean up this client
			// the break here breaks the Listen loop
			work <- Command{command: "disconnect", contents: line, clientId: c.id}
			break
		}
		work <- Command{command: "say", contents: line, clientId: c.id}

		// wait for reply to player's command
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

func (c *Client) Write(b []byte) (n int, err error) {
	return c.conn.Write(b)
}

func (c *Client) Read(b []byte) (n int, err error) {
	return c.conn.Read(b)
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
