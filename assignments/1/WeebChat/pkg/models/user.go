package models

import (
	"WeebChat/pkg/websocket"
	"time"
)

type User struct {
	Id           string                     `json:"id"`   //like public id -> hash of secret
	Name         string                     `json:"name"` //name
	Secret       string                     `json:"-"`    //like password
	Ws           *websocket.ServerWebSocket `json:"-"`
	LastActiveAt time.Time                  `json:"lastActiveAt"`
}
