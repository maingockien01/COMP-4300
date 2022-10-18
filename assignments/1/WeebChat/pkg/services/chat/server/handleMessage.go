package server

import (
	"WeebChat/pkg/models"
	"WeebChat/pkg/services/protocols"
	"WeebChat/pkg/websocket"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strings"
)

func (s *ChatServiceServer) HandleMessageRestful(w http.ResponseWriter, r *http.Request) {
	roomName := strings.TrimPrefix(r.URL.Path, "/messages/")
	room := s.ChatService.GetRoom(roomName)

	if room == nil {
		res := models.ErrorResponse{
			Message: "Not found room",
			Code:    http.StatusNotFound,
		}

		ReturnRestResponse(w, res, http.StatusNotFound)
	}

	switch r.Method {
	case http.MethodGet:

		messages := room.Messages

		ReturnRestResponse(w, messages, http.StatusOK)
		return
	case http.MethodPost:
		var message models.Message

		err := json.NewDecoder(r.Body).Decode(&message)

		if err != nil {
			ReturnRestResponse(w, INTERNAL_ERROR, http.StatusInternalServerError)
			return
		}

		err = s.AddMessage(&message, roomName)

		if err != nil {
			res := models.ErrorResponse{
				Message: err.Error(),
				Code:    http.StatusBadRequest,
			}

			ReturnRestResponse(w, res, http.StatusBadRequest)
		}

		return
	}
}

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
				return errors.New("found no room")
			}

			user := s.ChatService.GetUser(message.Sender)

			if user == nil {
				fmt.Println("Found no user - ", message.Sender)
				return errors.New("found no user")
			}

			returnMessages := room.GetMessages(message.Position, math.MaxInt)

			return s.PushMessage(user.Ws, returnMessages...)

		case 2:
			message1 := messages[0]
			message2 := messages[0]

			room := s.ChatService.GetRoom(message1.Room)

			if room == nil {
				fmt.Println("Found no room - ", message1.Room)
				return errors.New("found no room")
			}

			user := s.ChatService.GetUser(message1.Sender)

			if user == nil {
				fmt.Println("Found no user - ", message1.Sender)
				return errors.New("found no user")
			}

			returnMessages := room.GetMessages(message1.Position, message2.Position)

			return s.PushMessage(user.Ws, returnMessages...)
		default:
			return errors.New("invalid number of messages")
		}
	case protocols.DIRECTION_UPDATE:
		messages := protocolMessage.Data
		for _, message := range messages {
			s.AddMessage(message, message.Room)
		}
	}

	return nil
}

func (s *ChatServiceServer) AddMessage(message *models.Message, roomName string) error {
	room := s.ChatService.GetRoom(message.Room)
	if room == nil {
		fmt.Println("Found no room - ", message.Room)
		return errors.New("found no room")
	}

	user := room.GetUser(message.Sender)

	if user == nil {
		fmt.Println("Found no user - ", message.Sender)
		return errors.New("found no user")
	}

	message = room.AppendMessage(user, message)

	if message != nil {
		return s.BroadcastMessage(message, room.Name)
	}

	return errors.New("message is not appended")
}

func (s *ChatServiceServer) BroadcastMessage(message *models.Message, roomName string) error {
	room := s.ChatService.GetRoom(roomName)

	if room == nil {
		fmt.Println("Found no room - ", roomName)
		return errors.New("found no room")
	}

	if message != nil {
		usersInRoom := room.GetUsers()

		for i := 0; i < len(usersInRoom); i++ {
			userInRoom := usersInRoom[i]

			s.PushMessage(userInRoom.Ws, message)
		}
	}

	return nil
}
