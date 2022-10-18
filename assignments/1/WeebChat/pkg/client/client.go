package client

import (
	"WeebChat/pkg/models"
	"WeebChat/pkg/services/protocols"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

const PING_PERIOD = 50 * time.Second //secs

const WRITE_WAIT = 10 * time.Second

type ChatClient struct {
	ChatClientUser  *models.User
	Rooms           []*models.Room
	Ws              *websocket.Conn
	ChatUrl         string
	CurrentRoomName *string
}

func NewChatClient(user *models.User, chat_server_url string) *ChatClient {
	ChatClient := &ChatClient{
		ChatClientUser: user,
		Rooms:          make([]*models.Room, 0),
		ChatUrl:        chat_server_url,
		Ws:             nil,
	}

	return ChatClient
}

func (c *ChatClient) RegisterUser(user *models.User) error {
	registerUrl := c.ChatUrl + "/user"
	jsonBody, err := json.Marshal(c.ChatClientUser)

	if err != nil {
		return err
	}

	_, err = PostJson(registerUrl, []byte(jsonBody))

	if err != nil {
		return err
	}

	return nil
}

func (c *ChatClient) GetRoomList() []*models.Room {
	var rooms []*models.Room

	roomListUrl := "http://" + c.ChatUrl + "/rooms"

	err := GetJson(roomListUrl, &rooms)

	if err != nil {
		fmt.Println("Err on getting ", roomListUrl, ": ", err)
		return nil
	}

	return rooms
}

func (c *ChatClient) CreateRoom(roomName string, limit int) error {
	url := "http://" + c.ChatUrl + "/room"
	room := &models.Room{
		Name: roomName,
	}

	jsonBody, err := json.Marshal(room)

	if err != nil {
		return err
	}

	_, err = PostJson(url, []byte(jsonBody))

	if err != nil {
		return err
	}

	return nil
}

func (c *ChatClient) JoinRoom(roomName string) error {
	users, err := c.GetUserList(roomName)

	if err != nil {
		return err
	}

	for _, user := range users {
		if user.Id == c.ChatClientUser.Id {
			c.CurrentRoomName = &roomName
			return nil
		}
	}

	url := "http://" + c.ChatUrl + "/users/room/" + roomName
	jsonBody, err := json.Marshal(*c.ChatClientUser)
	if err != nil {
		return err
	}

	_, err = PostJson(url, []byte(jsonBody))

	if err != nil {
		return err
	}

	c.CurrentRoomName = &roomName

	return nil
}

func (c *ChatClient) GetUserList(roomName string) ([]*models.User, error) {
	var users []*models.User

	url := "http://" + c.ChatUrl + "/users/room/" + roomName

	err := GetJson(url, &users)

	return users, err
}

func (c *ChatClient) LeaveRoom(roomName string) error {
	url := "http://" + c.ChatUrl + "/users/room/" + roomName
	jsonBody, err := json.Marshal(*c.ChatClientUser)

	if err != nil {
		return err
	}

	_, err = DeleteJson(url, []byte(jsonBody))

	if err != nil {
		return err
	}

	if *c.CurrentRoomName == roomName {
		c.CurrentRoomName = nil
	}

	return nil
}

func (c *ChatClient) SendMessageRest(roomName string, content string) error {
	url := "http://" + c.ChatUrl + "/messages/" + roomName
	message := models.Message{
		Content: content,
		Sender:  c.ChatClientUser.Id,
		Room:    roomName,
	}

	jsonBody, err := json.Marshal(message)

	if err != nil {
		return err
	}

	resBody, err := PostJson(url, []byte(jsonBody))

	if err != nil {
		return err
	}

	var returnMessage models.Message

	err = json.Unmarshal(resBody, &returnMessage)

	if err != nil {
		return err
	}

	return c.AppendMessage(&returnMessage)
}

func (c *ChatClient) PullMessageRest(roomName string, from int, to int) []*models.Message {
	var messages []*models.Message

	url := "http://" + c.ChatUrl + "/messages/" + roomName

	GetJson(url, &messages)

	return messages
}

func (c *ChatClient) AppendMessage(message *models.Message) error {
	room := c.GetRoom(message.Room)

	if room == nil {
		return nil
	}

	if room.IsMessageIn(message) {
		return errors.New("message already in room")
	}

	room.AppendMessage(c.ChatClientUser, message)

	return nil
}

func (c *ChatClient) GetRoom(roomName string) *models.Room {
	for _, room := range c.Rooms {
		if room.Name == roomName {
			return room
		}
	}

	return nil
}

func (c *ChatClient) OpenChatSocket() error {

	url, err := url.Parse("ws://" + c.ChatUrl + "/chat")

	if err != nil {
		return err
	}

	ws, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		return err
	}

	c.Ws = ws

	go c.setPing()

	greetings := protocols.ProtocolUser{
		Metadata: protocols.ProtocolMetadata{
			From:      c.ChatClientUser.Name,
			Direction: protocols.DIRECTION_GREETING,
			Type:      protocols.TYPE_USER,
			Version:   protocols.V1,
		},
		Data: *c.ChatClientUser,
	}

	return ws.WriteJSON(greetings)
}

func (c *ChatClient) SendMessageSocket(roomName string, content string) error {
	messages := make([]*models.Message, 0)
	messages = append(messages, &models.Message{
		Content: content,
		Sender:  c.ChatClientUser.Id,
		Room:    roomName,
	})

	messageProtocol := protocols.ProtocolMessage{
		Metadata: protocols.ProtocolMetadata{
			From:      c.ChatClientUser.Name,
			Direction: protocols.DIRECTION_UPDATE,
			Type:      protocols.TYPE_MESSAGE,
			Version:   protocols.V1,
		},
		Data: messages,
	}

	return c.Ws.WriteJSON(messageProtocol)
}

func (c *ChatClient) ReceiveMessageSocket(handleMessage func(*models.Message)) error {
	var messageProtocol protocols.ProtocolMessage

	err := c.Ws.ReadJSON(&messageProtocol)

	if err != nil {
		return err
	}

	for _, message := range messageProtocol.Data {
		handleMessage(message)
	}

	return nil
}

func (c *ChatClient) CloseChatSocket() error {
	return c.Ws.Close()
}

func (c *ChatClient) handleNewMessage(message *models.Message) {
	err := c.AppendMessage(message)

	if err == nil {
		fmt.Printf("\n[%s]\t%s: %s\n", message.Room, message.Sender, message.Content)
		return
	} else {
		fmt.Println(err)
	}
}

func (c *ChatClient) setPing() {
	ticker := time.NewTicker(PING_PERIOD)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := c.Ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(WRITE_WAIT)); err != nil {
				log.Println("ping:", err)
			}
		}
	}
}
