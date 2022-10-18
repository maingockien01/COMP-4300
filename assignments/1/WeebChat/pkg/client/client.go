package client

import (
	"WeebChat/pkg/models"
	"WeebChat/pkg/services/protocols"
	"WeebChat/pkg/websocket"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
)

type Client struct {
	ClientUser *models.User
	Rooms      []*models.Room
	Ws         *websocket.ClientWebSocket
	ChatUrl    string
}

func (c *Client) RegisterUser(user *models.User) error {
	registerUrl := c.ChatUrl + "/user"
	jsonBody, err := json.Marshal(c.ClientUser)

	if err != nil {
		return err
	}

	_, err = PostJson(registerUrl, []byte(jsonBody))

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetRoomList() []*models.Room {
	var rooms []*models.Room

	roomListUrl := c.ChatUrl + "/rooms"

	GetJson(roomListUrl, &rooms)

	return rooms
}

func (c *Client) CreateRoom(roomName string, limit int) error {
	url := c.ChatUrl + "/room"
	room := models.Room{
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

func (c *Client) JoinRoom(roomName string) error {
	url := c.ChatUrl + "/users/room/" + roomName
	jsonBody, err := json.Marshal(c.ClientUser)

	if err != nil {
		return err
	}

	_, err = PostJson(url, []byte(jsonBody))

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) LeaveRoom(roomName string) error {
	url := c.ChatUrl + "/users/room/" + roomName
	jsonBody, err := json.Marshal(c.ClientUser)

	if err != nil {
		return err
	}

	_, err = DeleteJson(url, []byte(jsonBody))

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) SendMessageRest(roomName string, content string) error {
	url := c.ChatUrl + "/messages/" + roomName
	message := models.Message{
		Content: content,
		Sender:  c.ClientUser.Name,
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

func (c *Client) PullMessageRest(roomName string, from int, to int) []*models.Message {
	var messages []*models.Message

	url := c.ChatUrl + "/messages/" + roomName

	GetJson(url, &messages)

	return messages
}

func (c *Client) AppendMessage(message *models.Message) error {
	room := c.GetRoom(message.Room)

	if room == nil {
		return nil
	}

	if room.IsMessageIn(message) {
		return errors.New("message already in room")
	}

	room.AppendMessage(c.ClientUser, message)

	return nil
}

func (c *Client) GetRoom(roomName string) *models.Room {
	for _, room := range c.Rooms {
		if room.Name == roomName {
			return room
		}
	}

	return nil
}

func (c *Client) OpenChatSocket() error {

	url, err := url.Parse(c.ChatUrl + "/chat")

	if err != nil {
		return err
	}

	ws, err := websocket.NewClientWebSocket(url, nil)

	if err != nil {
		return err
	}

	c.Ws = ws

	greetings := protocols.ProtocolUser{
		Metadata: protocols.ProtocolMetadata{
			From:      c.ClientUser.Name,
			Direction: protocols.DIRECTION_GREETING,
			Type:      protocols.TYPE_USER,
			Version:   protocols.V1,
		},
		Data: *c.ClientUser,
	}

	jsonGreetingBody, err := json.Marshal(greetings)

	if err != nil {
		return err
	}

	greetingFrame := websocket.NewFrameMessage(jsonGreetingBody)

	return ws.Send(*greetingFrame)
}

func (c *Client) SendMessageSocket(roomName string, content string) error {
	messages := make([]*models.Message, 0)
	messages = append(messages, &models.Message{
		Content: content,
		Sender:  c.ClientUser.Name,
		Room:    roomName,
	})

	messageProtocol := protocols.ProtocolMessage{
		Metadata: protocols.ProtocolMetadata{
			From:      c.ClientUser.Name,
			Direction: protocols.DIRECTION_UPDATE,
			Type:      protocols.TYPE_MESSAGE,
			Version:   protocols.V1,
		},
		Data: messages,
	}

	jsonMessageBody, err := json.Marshal(messageProtocol)

	if err != nil {
		return err
	}

	messageFrame := websocket.NewFrameMessage(jsonMessageBody)

	return c.Ws.Send(*messageFrame)
}

func (c *Client) ReceiveMessageSocket(roomName string, handleMessage func(*models.Message)) error {
	frame, err := c.Ws.Receive()

	if err != nil {
		return err
	}

	frameBody := frame.ParseText()

	var messageProtocol protocols.ProtocolMessage

	err = json.Unmarshal([]byte(frameBody), &messageProtocol)

	if err != nil {
		return err
	}

	for _, message := range messageProtocol.Data {
		handleMessage(message)
	}

	return nil
}

func (c *Client) CloseChatSocket() error {
	return c.Ws.Close()
}

func (c *Client) handleNewMessage(message *models.Message) {
	err := c.AppendMessage(message)

	if err != nil {
		fmt.Printf("[%s] %s: %s\n", message.Room, message.Sender, message.Content)
		return
	}
}
