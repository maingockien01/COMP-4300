package models

type Room struct {
	Messages []Message
	Name     string
	secret   string
}
