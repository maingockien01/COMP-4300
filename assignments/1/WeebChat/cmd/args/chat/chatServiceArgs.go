package main

import (
	chat "WeebChat/pkg/services/chat/server"
	"fmt"
	"os"
)

func main() {
	argsWithoutProgs := os.Args[1:]
	host := argsWithoutProgs[0]
	port := argsWithoutProgs[1]
	serverName := argsWithoutProgs[2]

	fmt.Printf("Setting up chat service on %s:%s\n", host, port)
	chatServer := chat.NewChatServiceServer(host, port, serverName)
	chatServer.Setup()
	chatServer.Start()
}
