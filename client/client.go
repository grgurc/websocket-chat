package client

import (
	"github.com/gorilla/websocket"
	"github.com/grgurc/websocket-chat/model"
)

// Each connected client will have 2 goroutines responsible for it
// 1. Runs in loop listening for messages
// 2. Sends messages to server
type Client struct {
	conn *websocket.Conn
	in   chan model.Message // Receives messages from server
	out  chan model.Message // equals server.incoming, send msgs from client here
}

func NewClient(conn *websocket.Conn, in, out chan model.Message) *Client {
	return &Client{
		conn: conn,
		in:   in,
		out:  out,
	}
}

func (c *Client) Run() {
	go c.Listen()

	// Write msgs from server to client
	for {
		msg := <-c.in
		if err := c.conn.WriteJSON(msg); err != nil {
			panic(err)
		}
	}
}

func (c *Client) Listen() {
	for {
		var msg model.Message
		if err := c.conn.ReadJSON(&msg); err != nil {
			break
		}
		c.out <- msg
	}
}
