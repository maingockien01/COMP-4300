package models

import (
	"errors"
	"math"
	"sync"
	"time"
)

type Room struct {
	Messages []*Message `json:"-"`
	Name     string     `json:"name"`
	secret   *string    `json:"-"`
	Users    []*User    `json:"-"`
	Lock     sync.Mutex `json:"-"`
	Capacity int        `json:"capacity"`
	Limit    int        `json:"limit"`
}

type WrongSecretError struct{}

const NUM_LATEST_MESSAGE = 10

const DEFAULT_ROOM_LIMIT = 5

func (e WrongSecretError) Error() string {
	return "Wrong secret"
}

func (r *Room) AddUser(user *User) error {

	if r.Capacity >= r.Limit {
		return errors.New("Capacity limited")
	}

	if r.IsUserIn(user.Id) {
		return errors.New("User already in the room")
	}

	user.LastActiveAt = time.Now()
	r.Users = append(r.Users, user)

	r.Capacity = len(r.Users)

	return nil
}

func (r *Room) RemoveUser(user *User) error {
	for i, u := range r.Users {
		if u.Id == user.Id {
			r.Users = append(r.Users[:i], r.Users[i+1:]...)
			r.Capacity = len(r.Users)

			return nil
		}
	}

	return errors.New("No users found")
}

func (r *Room) IsUserIn(userId string) bool {
	for _, user := range r.Users {
		if user.Id == userId {
			return true
		}
	}
	return false
}

func (r *Room) AppendMessage(sender *User, message *Message) *Message {

	if !r.IsUserIn(sender.Id) {
		return nil
	}

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

func (r *Room) GetUser(userId string) *User {
	for _, user := range r.Users {
		if user.Id == userId {
			return user
		}
	}
	return nil
}

func NewRoom(name string) *Room {
	room := Room{
		Name:     name,
		Capacity: 0,
		Limit:    DEFAULT_ROOM_LIMIT,
		Lock:     sync.Mutex{},
	}

	return &room
}

func (r *Room) GetName() string {
	return r.Name
}
