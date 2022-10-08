package models

import (
	"time"
)

type Message struct {
	Sender   string
	Content  string
	Room     string
	TimeAt   time.Time
	Position int
}

func NewMessage(sender string, content string) *Message {
	message := Message{
		Sender:  sender,
		Content: content,
		TimeAt:  time.Now(),
	}

	return &message
}
