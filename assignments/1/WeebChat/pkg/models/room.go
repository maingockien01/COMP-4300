package models

import (
	"math"
	"sync"
	"time"
)

type Room struct {
	Messages []*Message
	Name     string
	secret   *string
	Users    []*User
	Lock     sync.Mutex
}

type WrongSecretError struct{}

const NUM_LATEST_MESSAGE = 10

func (e WrongSecretError) Error() string {
	return "Wrong secret"
}

func (r *Room) AddUser(user *User) error {

	user.LastActiveAt = time.Now()
	r.Users = append(r.Users, user)

	return nil
}

func (r *Room) RemoveUser(user *User) {
	for i, u := range r.Users {
		if u.Id == user.Id {
			r.Users = append(r.Users[:i], r.Users[i+1:]...)
		}
	}
}

func (r *Room) AppendMessage(sender *User, message *Message) *Message {
	r.Lock.Lock()
	defer r.Lock.Unlock()
	message.Sender = sender.Id
	message.Position = len(r.Messages)
	r.Messages = append(r.Messages, message)
	message.Room = r.Name
	sender.LastActiveAt = time.Now()

	return message
}

func (r *Room) GetAllMessages() []*Message {
	return r.Messages
}

func (r *Room) GetLatestMessages() []*Message {
	if len(r.Messages) > NUM_LATEST_MESSAGE {
		messages := r.Messages[len(r.Messages)-NUM_LATEST_MESSAGE:]
		return messages
	} else {
		return r.Messages
	}
}

func (r *Room) GetMessages(from int, to int) []*Message {
	from = int(math.Max(0, float64(from)))
	to = int(math.Min(float64(len(r.Messages)-1), float64(to)))

	return r.Messages[from:to]
}

func (r *Room) GetUsers() []*User {
	return r.Users
}

func NewRoom(name string) *Room {
	room := Room{
		Name: name,
	}

	return &room
}

func (r *Room) GetName() string {
	return r.Name
}
