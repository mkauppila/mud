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

// TODO: client id should be an uuid
// TODO: pass in the command parser
func NewClient(conn net.Conn, id int) *Client {
	client := Client{
		id:        id,
		conn:      conn,
		broadcast: make(chan string),
		reply:     make(chan string),
	}
	return &client
}

type ServerAction func(server *Server) error
type CommandAction func(command Command) ServerAction

func ConnectCommandAction(command Command) ServerAction {
	return func(s *Server) error {
		for _, c := range s.clients {
			if c.id == command.clientId {
				c.reply <- "You joined\n"
			} else {
				c.broadcast <- fmt.Sprintf("Client %d joined!\n", command.clientId)
			}
		}
		return nil
	}
}

func DisconnectCommandAction(command Command) ServerAction {
	return func(s *Server) error {
		for i, c := range s.clients {
			if c.id == command.clientId {
				c.reply <- "You are disconnecting"
				c.Disconnect()
				// update the clients list
				s.clients[i] = s.clients[len(s.clients)-1]
				s.clients = s.clients[:len(s.clients)-1]
			} else {
				c.broadcast <- fmt.Sprintf("Client %d disconnecting...\n", command.clientId)
			}
		}
		return nil
	}
}

// Needs client id
func SayCommandAction(command Command) ServerAction {
	return func(s *Server) error {
		for _, c := range s.clients {
			if c.id == command.clientId {
				c.reply <- fmt.Sprintf("You said %s\n", command.contents)
			} else {
				c.broadcast <- fmt.Sprintf("They said %s\n", command.contents)
			}
		}
		return nil
	}
}

func (c *Client) Listen(work chan<- ServerAction) {
	reader := bufio.NewReader(c.conn)

	work <- ConnectCommandAction(Command{command: "connect", contents: "", clientId: c.id})
	connectReply := <-c.reply
	c.directReply(connectReply)

	// m := map[string]CommandAction{}
	// m["say"] = SayCommandAction
	// m["connect"] = ConnectCommandAction
	// m["disconnect"] = DisconnectCommandAction
	// fmt.Println(m)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			// disconnect command tells the server to clean up this client
			// the break here breaks the Listen loop
			work <- DisconnectCommandAction(Command{command: "disconnect", contents: line, clientId: c.id})
			break
		}

		command := ParseCommand(line)
		command.clientId = c.id

		work <- SayCommandAction(command)

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

/*
func (c *Client) Write(b []byte) (n int, err error) {
	return c.conn.Write(b)
}

func (c *Client) Read(b []byte) (n int, err error) {
	return c.conn.Read(b)
}
*/

func (c *Client) Disconnect() {
	fmt.Printf("Disconnecting %d\n", c.id)

	close(c.broadcast)
	close(c.reply)

	err := c.conn.Close()
	if err != nil {
		panic(err)
	}
}
