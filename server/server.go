package server

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/grgurc/websocket-chat/client"
	"github.com/grgurc/websocket-chat/model"
)

type Server struct {
	ConnChan chan *websocket.Conn

	clientChans []chan model.Message // send msgs to clients
	messages    []model.Message      // ordered list of all messages
	incoming    chan model.Message   // forward to clients
}

// Creates a new server that handles saving messages
// and sending received messages to clients.
// To add a new client, simply push a ConnectionRequest to
// the ConnectionRequests channel
func NewServer() *Server {
	return &Server{
		clientChans: make([]chan model.Message, 0),
		messages:    make([]model.Message, 0),
		ConnChan:    make(chan *websocket.Conn, 5),
		incoming:    make(chan model.Message, 1000),
	}
}

func (s *Server) Run() {
	for {
		select {
		case conn := <-s.ConnChan:
			s.ConnectClient(conn)
		case msg := <-s.incoming:
			s.ReceiveMessage(msg)
		}
	}
}

func (s *Server) ConnectClient(conn *websocket.Conn) error {
	clientIncoming := make(chan model.Message, 10)

	client := client.NewClient(conn, clientIncoming, s.incoming)

	go client.Run()
	s.clientChans = append(s.clientChans, clientIncoming)

	// Send all messages so far to client
	// !!! Prolly slow af, work on this...
	for _, m := range s.messages {
		clientIncoming <- m
	}

	return nil
}

// Save msg and forward to all clients
func (s *Server) ReceiveMessage(msg model.Message) {
	s.messages = append(s.messages, msg)
	log.Println(s.messages)

	for _, ch := range s.clientChans {
		ch <- msg
	}
}
