package server

import (
	"WeebChat/pkg/services/protocols"
	"WeebChat/pkg/websocket"
	"encoding/json"
	"errors"
	"fmt"
	"math"
)

func (s *ChatServiceServer) HandleMessage(payload string, frame websocket.Frame, ws *websocket.ServerWebSocket) error {
	var protocolMessage protocols.ProtocolMessage
	err := json.Unmarshal([]byte(payload), &protocolMessage)

	if err != nil {
		return err
	}

	switch protocolMessage.Metadata.Direction {
	case protocols.DIRECTION_PULL:
		messages := protocolMessage.Data

		switch len(messages) {
		case 1:
			message := messages[0]
			room := s.ChatService.GetRoom(message.Room)

			if room == nil {
				fmt.Println("Found no room - ", message.Room)
				return errors.New("Found no room")
			}

			user := s.ChatService.GetUser(message.Sender)

			if user == nil {
				fmt.Println("Found no user - ", message.Sender)
				return errors.New("Found no user")
			}

			returnMessages := room.GetMessages(message.Position, math.MaxInt)

			return s.PushMessage(user.Ws, returnMessages...)

		case 2:
			message1 := messages[0]
			message2 := messages[0]

			room := s.ChatService.GetRoom(message1.Room)

			if room == nil {
				fmt.Println("Found no room - ", message1.Room)
				return errors.New("Found no room")
			}

			user := s.ChatService.GetUser(message1.Sender)

			if user == nil {
				fmt.Println("Found no user - ", message1.Sender)
				return errors.New("Found no user")
			}

			returnMessages := room.GetMessages(message1.Position, message2.Position)

			return s.PushMessage(user.Ws, returnMessages...)
		default:
			return errors.New("Invalid number of messages")
		}
	case protocols.DIRECTION_UPDATE:
		messages := protocolMessage.Data
		for _, message := range messages {
			room := s.ChatService.GetRoom(message.Room)
			if room == nil {
				fmt.Println("Found no room - ", message.Room)
				continue
			}

			user := s.ChatService.GetUser(message.Sender)

			if user == nil {
				fmt.Println("Found no user - ", message.Sender)
				continue
			}

			message = room.AppendMessage(user, message)

			usersInRoom := room.GetUsers()

			for i := 0; i < len(usersInRoom); i++ {
				userInRoom := usersInRoom[i]

				s.PushMessage(userInRoom.Ws, message)
			}
		}
	}

	return nil
}
