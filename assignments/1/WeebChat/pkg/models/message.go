package models

import (
	"time"
)

type Message struct {
	Sender   string    `json:"sender"`
	Content  string    `json:"content"`
	Room     string    `json:"room"`
	TimeAt   time.Time `json:"timeAt"`
	Position int       `json:"position"`
}

func NewMessage(sender string, content string) *Message {
	message := Message{
		Sender:  sender,
		Content: content,
		TimeAt:  time.Now(),
	}

	return &message
}
