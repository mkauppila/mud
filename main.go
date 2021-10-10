package main

import (
	"bufio"
	"fmt"
	"net"
)

// Client

type Client struct {
	id       int
	conn     net.Conn
	outgoing chan string
	incoming chan string
}

func NewClient(conn net.Conn, id int) *Client {
	client := Client{
		id:       id,
		conn:     conn,
		outgoing: make(chan string),
		incoming: make(chan string),
	}
	return &client
}

func (c *Client) Listen() {
	reader := bufio.NewReader(c.conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("End of Reader goroutine for ", c.id)
			break
		}
		fmt.Println("->: ", line)
		// TODO: parse the command
		// TODO: send the message to the server to be handled?
		c.outgoing <- line
	}
}

func (c *Client) Serve() {
	writer := bufio.NewWriter(c.conn)
	for {
		message := <-c.outgoing
		writer.WriteString(message)
		writer.Flush()
	}
}

func (c *Client) Write(b []byte) (n int, err error) {
	return c.conn.Write(b)
}

func (c *Client) Read(b []byte) (n int, err error) {
	return c.conn.Read(b)
}

func (c *Client) Disconnect() {
	err := c.conn.Close()
	if err != nil {
		panic(err)
	}
}

// Server

type Server struct {
	clients []*Client
}

func NewServer() Server {
	return Server{}
}

func (s *Server) AddNewConnection(conn net.Conn) {
	client := NewClient(conn, len(s.clients))
	go client.Listen()
	go client.Serve()
	s.clients = append(s.clients, client)
}

func (s *Server) HandleChat() {
	// for {
	// 	var b = make([]byte, 1024)
	// 	n, err := client.Read(b)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	message := string(b[:n])

	// 	if len(message) > 0 {
	// 		output := fmt.Sprintf("%d -> %s", client.id, message)

	// 		for _, c := range clients {
	// 			fmt.Println("hello ", clients)
	// 			_, err := c.Write([]byte(output))
	// 			if err != nil {
	// 				panic(err)
	// 			}
	// 		}
	// 	}
	// }
}

func main() {
	address := "localhost:6000"
	fmt.Printf("starting at %s\n", address)

	ln, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	server := NewServer()
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		server.AddNewConnection(conn)
	}

	// server.Disconnect()
	// for _, client := range clients {
	// 	client.Disconnect()
	// }
}
