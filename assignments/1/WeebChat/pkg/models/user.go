package models

import (
	"WeebChat/pkg/websocket"
	"time"
)

type User struct {
	Id           string //like public id -> hash of secret
	Name         string //name
	Secret       string //like password
	Ws           *websocket.ServerWebSocket
	LastActiveAt time.Time
}
