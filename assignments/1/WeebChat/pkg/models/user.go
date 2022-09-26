package models

type User struct {
	Id     string //like public id -> hash of secret
	Name   string //name
	secret string //like password
}
