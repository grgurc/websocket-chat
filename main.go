package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/grgurc/websocket-chat/server"
)

func main() {
	s := server.NewServer()
	go s.Run()

	u := websocket.Upgrader{}
	// listen for websocket connections, just forward to server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, err := u.Upgrade(w, r, http.Header{})
		if err != nil {
			http.Error(w, "error creating websocket connection", http.StatusInternalServerError)
		}

		s.ConnChan <- conn
	})

	log.Println("Listening on port 8000")
	http.ListenAndServe(":8000", nil)
}
