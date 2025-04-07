package model

import "time"

// Basic message
type Message struct {
	Sender  string    `json:"sender"`
	Time    time.Time `json:"time"`
	Content string    `json:"content"`
}
