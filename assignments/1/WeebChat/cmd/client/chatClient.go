package main

import (
	"WeebChat/pkg/client"
	"WeebChat/pkg/models"
	"os"
)

func main() {
	chatUrl := os.Getenv("CHAT_URL")
	userId := os.Getenv("USER_ID")
	userName := os.Getenv("USER_NAME")

	user := &models.User{
		Id:   userId,
		Name: userName,
	}

	chatClient := client.NewChatClient(user, chatUrl)

	commander := client.NewCommander()

	commander.ChatClient = chatClient

	commander.Run()
}
