package main

import (
	"WeebChat/pkg/client"
	"WeebChat/pkg/models"
	"os"
)

func main() {
	argsWithoutProgs := os.Args[1:]

	chatUrl := argsWithoutProgs[0]
	userId := argsWithoutProgs[1]
	userName := argsWithoutProgs[2]

	user := &models.User{
		Id:   userId,
		Name: userName,
	}

	chatClient := client.NewChatClient(user, chatUrl)

	commander := client.NewCommander()

	commander.ChatClient = chatClient

	commander.Run()
}
