package server

import (
	"WeebChat/pkg/models"
	"WeebChat/pkg/services/protocols"
	"WeebChat/pkg/websocket"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (s *ChatServiceServer) HandleUser(payload string, frame websocket.Frame, ws *websocket.ServerWebSocket) error {
	var protocolUser protocols.ProtocolUser

	err := json.Unmarshal([]byte(payload), &protocolUser)

	if err != nil {
		return nil
	}

	switch protocolUser.Metadata.Direction {
	case protocols.DIRECTION_GREETING:
		fmt.Println("Adding user...")
		protocolUser.Data.Ws = ws
		return s.ChatService.AddUser(&protocolUser.Data)
	}

	return nil
}

func (s *ChatServiceServer) HandleAddUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var user models.User

		err := json.NewDecoder(r.Body).Decode(&user)

		if err != nil {
			ReturnRestResponse(w, INTERNAL_ERROR, http.StatusInternalServerError)
			return
		}

		err = s.ChatService.AddUser(&user)

		if err != nil {
			res := models.ErrorResponse{
				Message: err.Error(),
				Code:    http.StatusBadRequest,
			}
			ReturnRestResponse(w, res, res.Code)
			return
		}

		w.WriteHeader(http.StatusOK)

		return

		// case http.MethodDelete:
		// 	var user models.User

		// 	err := json.NewDecoder(r.Body).Decode(&user)

		// 	if err != nil {
		// 		ReturnRestResponse(w, INTERNAL_ERROR, http.StatusInternalServerError)
		// 		return
		// 	}

		// 	w.WriteHeader(http.StatusOK)
		// 	return
	}
}

func (s *ChatServiceServer) HandleGetUserIn(w http.ResponseWriter, r *http.Request) {
	roomName := strings.TrimPrefix(r.URL.Path, "/users/room/")
	room := s.ChatService.GetRoom(roomName)

	if room == nil {
		res := models.ErrorResponse{
			Message: "Found no rooom",
			Code:    http.StatusNotFound,
		}
		ReturnRestResponse(w, res, res.Code)
		return
	}
	switch r.Method {
	case http.MethodGet:
		users := room.GetUsers()

		ReturnRestResponse(w, users, http.StatusOK)
		return
	case http.MethodPost:
		var user models.User

		err := json.NewDecoder(r.Body).Decode(&user)

		if err != nil {
			ReturnRestResponse(w, INTERNAL_ERROR, http.StatusInternalServerError)
			return
		}

		if s.ChatService.GetUser(user.Id) == nil {
			res := models.ErrorResponse{
				Message: "User has not been registered",
				Code:    http.StatusBadRequest,
			}
			ReturnRestResponse(w, res, res.Code)
			return
		}

		err = s.ChatService.JoinUser(user.Id, roomName)

		if err != nil {
			res := models.ErrorResponse{
				Message: err.Error(),
				Code:    http.StatusBadRequest,
			}
			ReturnRestResponse(w, res, res.Code)
			return
		}

		w.WriteHeader(http.StatusOK)

		return

	case http.MethodDelete:
		var user models.User

		err := json.NewDecoder(r.Body).Decode(&user)

		if err != nil {
			ReturnRestResponse(w, INTERNAL_ERROR, http.StatusInternalServerError)
			return
		}

		err = room.RemoveUser(&user)

		if err != nil {
			res := models.ErrorResponse{
				Message: err.Error(),
				Code:    http.StatusBadRequest,
			}
			ReturnRestResponse(w, res, res.Code)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}

}
