package main

import (
	chat "WeebChat/pkg/services/chat/server"
	"fmt"
	"os"
)

func main() {
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	serverName := os.Getenv("SERVER_NAME")

	fmt.Printf("Setting up chat service on %s:%s\n", host, port)
	chatServer := chat.NewChatServiceServer(host, port, serverName)
	chatServer.Setup()
	chatServer.Start()
}
