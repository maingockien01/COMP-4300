package models

import (
	"net"
	"time"
)

type User struct {
	Id     string //like public id -> hash of secret
	Name   string //name
	Secret string //like password
	Conn   net.Conn
	LastActiveAt time.Time
}
