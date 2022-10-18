package server

import (
	"WeebChat/pkg/models"
	"WeebChat/pkg/services/protocols"
	"WeebChat/pkg/websocket"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func (s *ChatServiceServer) HandleRoom(payload string, frame websocket.Frame, ws *websocket.ServerWebSocket) error {
	var protocolRoom protocols.ProtocolRoom
	err := json.Unmarshal([]byte(payload), &protocolRoom)

	if err != nil {
		return err
	}

	switch protocolRoom.Metadata.Direction {
	case protocols.DIRECTION_JOIN:
		for i := 0; i < len(protocolRoom.Data); i++ {
			r := &protocolRoom.Data[i]
			err := s.ChatService.JoinUser(protocolRoom.Metadata.From, r.Name)

			if err != nil {
				fmt.Println("Error on joining rooms ... ", err)
				return err
			}

			//Push latest messages to clients
			roomInDatabase := s.ChatService.GetRoom(r.Name)

			if roomInDatabase == nil {
				return errors.New("there is no room")
			}

			messages := roomInDatabase.GetLatestMessages()

			err = s.PushMessage(ws, messages...)

			if err != nil {
				fmt.Println("Error on pushing messages ... ", err)
				return err
			}
		}

		return nil
	}

	return nil
}

func (s *ChatServiceServer) HandlerGetRooms(w http.ResponseWriter, r *http.Request) {
	rooms := s.ChatService.GetRooms()

	ReturnRestResponse(w, rooms, http.StatusOK)
}

func (s *ChatServiceServer) HandlerRoom(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:

		var room models.Room

		err := json.NewDecoder(r.Body).Decode(&room)

		if err != nil {
			fmt.Println("Error on handle room: Decode request body ", err)
			ReturnRestResponse(w, INTERNAL_ERROR, http.StatusInternalServerError)
			return
		}

		if room.Name == "" {
			ReturnRestResponse(w, models.ErrorResponse{
				Message: "Invalid body",
				Code:    http.StatusBadRequest,
			}, http.StatusBadRequest)
			return
		}

		err = s.ChatService.CreateRoom(room.Name)

		if err != nil {
			ReturnRestResponse(w, models.ErrorResponse{
				Message: "Duplicate room name",
				Code:    http.StatusBadRequest,
			}, http.StatusBadRequest)
			return
		}

		newRoom := s.ChatService.GetRoom(room.Name)

		ReturnRestResponse(w, newRoom, http.StatusOK)

		return
	default:
		res := models.ErrorResponse{
			Message: "Not found",
			Code:    http.StatusNotFound,
		}

		ReturnRestResponse(w, res, http.StatusNotFound)
		return
	}
}
