package client

import "strings"

type Commander struct {
	Commands   map[string]func(string)
	ChatClient *Client
}

func NewCommander() *Commander {
	commander := Commander{
		Commands: make(map[string]func(string)),
	}

	return &commander
}

func (c *Commander) ParseCommand(input string) func(string) {
	command := strings.Split(input, " ")[0]
	return c.Commands[command]
}
