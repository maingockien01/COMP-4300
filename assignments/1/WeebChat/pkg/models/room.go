package models

import (
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

func (e WrongSecretError) Error() string {
	return "Wrong secret"
}

func (r *Room) AddUser(user *User, secret *string) error {
	if r.secret != nil && *(r.secret) != *secret {
		return WrongSecretError{}
	}
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

func (r *Room) AppendMessage(sender *User, message *Message) {
	r.Lock.Lock()
	defer r.Lock.Unlock()
	message.From = sender.Id
	r.Messages = append(r.Messages, message)
	sender.LastActiveAt = time.Now()
}

func (r *Room) GetMessages() []*Message {
	return r.Messages
}

func (r *Room) GetUsers() []*User {
	return r.Users
}

func NewRoom(name string, secret *string) *Room {
	room := Room{
		Name:   name,
		secret: secret,
	}

	return &room
}

func (r *Room) GetName() string {
	return r.Name
}
