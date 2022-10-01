package models

import (
	"time"
)

type ChatService struct {
	Host     string
	Port     string
	Name     string
	Rooms    []*Room
	LastPing time.Time
}

func NewChatService(name string) *ChatService {
	return &ChatService{
		Name: name,
	}
}

func (c *ChatService) CreateRoom(name string) {
	//Create a new room
	//Add room to rooms list
}

func (c *ChatService) DeleteRoom(name string) {
	//Delete room from rooms list
}

func (c *ChatService) GetRoom(name string) {
	//Return room from rooms list
}

func (c *ChatService) GetRooms() {
	//Return rooms list
}
