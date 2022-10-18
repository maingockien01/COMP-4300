package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const (
	COMMAND_ROOM_LIST   = "\\room_list"
	COMMAND_ROOM_JOIN   = "\\room_join"
	COMMAND_ROOM_CREATE = "\\room_create"
	COMMAND_ROOM_LEAVE  = "\\room_leave"

	COMMAND_MESSAGE_SEND = "\\message"

	COMMAND_HELP = "\\help"
)

type Commander struct {
	Commands   map[string]func(*Commander, string)
	ChatClient *ChatClient
}

func NewCommander() *Commander {
	commander := Commander{
		Commands:   make(map[string]func(*Commander, string)),
		ChatClient: nil,
	}

	commander.Commands[COMMAND_ROOM_LIST] = CommandRoomList
	commander.Commands[COMMAND_ROOM_JOIN] = CommandRoomJoin
	commander.Commands[COMMAND_ROOM_CREATE] = CommandRoomCreate
	commander.Commands[COMMAND_ROOM_LEAVE] = CommandRoomLeave

	commander.Commands[COMMAND_MESSAGE_SEND] = CommandMessageSend

	commander.Commands[COMMAND_HELP] = Help

	return &commander
}

func (c *Commander) ParseCommand(input string) func(*Commander, string) {
	command := strings.Split(input, " ")[0]

	return c.Commands[command]
}

func (c *Commander) ParseData(input string) string {
	command := strings.Split(input, " ")[0]

	data := strings.TrimSpace(strings.TrimPrefix(input, command))

	return data
}

func CommandRoomList(commander *Commander, data string) {
	rooms := commander.ChatClient.GetRoomList()
	commander.ChatClient.Rooms = rooms

	jsonByte, err := json.MarshalIndent(rooms, "", "\t")

	if err != nil {
		fmt.Println("Err ", err)
		return
	}

	fmt.Println(string(jsonByte))
}

func CommandRoomCreate(commander *Commander, data string) {
	err := commander.ChatClient.CreateRoom(data, 10)
	if err == nil {
		fmt.Println("Room is created")
	} else {
		fmt.Println("Error on creating room: ", err)
	}
}

func CommandRoomLeave(commander *Commander, data string) {
	err := commander.ChatClient.LeaveRoom(data)
	if err == nil {
		fmt.Println("User left room")
	} else {
		fmt.Println("Error on leaving room: ", err)
	}
}

func CommandMessageSend(commander *Commander, data string) {
	roomName := commander.ChatClient.CurrentRoomName
	if roomName == nil {
		fmt.Println("Join a room before sending a message")
	} else {
		commander.ChatClient.SendMessageSocket(*roomName, data)
	}
}

func CommandRoomJoin(commander *Commander, data string) {
	err := commander.ChatClient.JoinRoom(data)

	if err == nil {
		fmt.Println("Joined!")
	} else {
		fmt.Println("Error on joining room ", err)
	}
}

func Help(c *Commander, data string) {
	fmt.Println("Command list:")
	fmt.Println("Use `\\room_list` to print out all chat rooms in server")
	fmt.Println("Use `\\room_join [roomName]` to join/switch to room")
	fmt.Println("Use `\\room_create [roomName]` to create rooms in server")
	fmt.Println("Use `\\room_leave [roomName]` to leave chat rooms in server")
	fmt.Println("Use `\\message [message]` to send message to current chat room")
	fmt.Println("Use `\\help ` to print all commands")

}

func (c *Commander) Run() {
	if c.ChatClient == nil {
		panic("chat client is nil")
	}

	err := c.ChatClient.OpenChatSocket()

	if err != nil {
		panic("chat client unable to open websocket with server" + err.Error())
	}

	go func() {
		for {
			rooms := c.ChatClient.GetRoomList()
			c.ChatClient.Rooms = rooms
			err := c.ChatClient.ReceiveMessageSocket(c.ChatClient.handleNewMessage)
			if err != nil {
				fmt.Println("There is error with receiving message: ", err.Error())
			}
		}
	}()

	reader := bufio.NewReader(os.Stdin)

	for {
		if c.ChatClient.CurrentRoomName == nil {
			fmt.Printf("[No room]\t")
		} else {
			fmt.Printf("[%s]\t", *c.ChatClient.CurrentRoomName)
		}

		command, _ := reader.ReadString('\n')
		command = strings.Replace(command, "\n", "", -1)

		commandFunction := c.ParseCommand(command)
		data := c.ParseData(command)

		if commandFunction == nil {
			fmt.Println("Invalid command!")
		} else {
			commandFunction(c, data)
		}

	}
}
