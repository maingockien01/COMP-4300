package client

import (
	"WeebChat/pkg/models"
	"WeebChat/pkg/websocket"
)

type Client struct {
	ClientUser *models.User
	Rooms      []*models.Room
	Ws         *websocket.ClientWebSocket
	ChatUrl    string
}

func (c *Client) RegisterUser(user *models.User) error {
	res, err := PostJson()

}

func (c *Client) GetRoomList() []*models.Room {
	var rooms []*models.Room

	roomListUrl := c.ChatUrl + "/rooms"

	GetJson(roomListUrl, &rooms)

	return rooms
}

func (c *Client) JoinRoom(roomName string) {
}

func (c *Client) SendMessage(roomName string, content string) {

}

func (c *Client) PullMessage(roomName string, from int, to int) []*models.Message {

}

func (c *Client) AppendMessage(message *models.Message) error {

}
