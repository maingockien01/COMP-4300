package models

import (
	"errors"
	"time"
)

type ChatService struct {
	Host     string
	Port     string
	Name     string
	Rooms    []*Room
	Users    []*User
	LastPing time.Time
}

func NewChatService(name string) *ChatService {
	return &ChatService{
		Name:  name,
		Rooms: make([]*Room, 0),
		Users: make([]*User, 0),
	}
}

func (c *ChatService) CreateRoom(name string) error {
	//Create a new room
	//Add room to rooms list
	room := c.GetRoom(name)

	if room != nil {
		return errors.New("Duplicate room name")
	}

	newRoom := NewRoom(name)
	c.Rooms = append(c.Rooms, newRoom)

	return nil
}

func (c *ChatService) DeleteRoom(name string) {
	//Delete room from rooms list
	for i, r := range c.Rooms {
		if r.Name == name {
			c.Rooms = append(c.Rooms[:i], c.Rooms[i+1:]...)
		}
	}

}

func (c *ChatService) GetRoom(name string) *Room {
	//Return room from rooms list
	for _, r := range c.Rooms {
		if r.Name == name {
			return r
		}
	}

	return nil
}

func (c *ChatService) GetRooms() []*Room {
	//Return rooms list
	return c.Rooms
}

func (c *ChatService) AddUser(user *User) {
	for _, u := range c.Users {
		if u.Id == user.Id {

			if u.Ws == nil {
				u.Ws = user.Ws
			}

			return
		}
	}

	user.LastActiveAt = time.Now()

	c.Users = append(c.Users, user)
}

func (c *ChatService) GetUser(userId string) *User {
	for _, u := range c.Users {
		if u.Id == userId {
			return u
		}
	}

	return nil
}

func (c *ChatService) JoinUser(userId string, roomName string) error {
	user := c.GetUser(userId)

	if user == nil {
		return errors.New("Found no user")
	}

	room := c.GetRoom(roomName)

	if room == nil {
		return errors.New("Found no room")
	}

	room.AddUser(user)

	return nil

}
