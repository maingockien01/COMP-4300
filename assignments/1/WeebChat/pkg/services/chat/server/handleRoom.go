package server

import (
	"WeebChat/pkg/services/protocols"
	"WeebChat/pkg/websocket"
	"encoding/json"
	"errors"
	"fmt"
)

func (s *ChatServiceServer) HandleRoom(payload string, frame websocket.Frame, ws *websocket.ServerWebSocket) error {
	var protocolRoom protocols.ProtocolRoom
	err := json.Unmarshal([]byte(payload), &protocolRoom)

	if err != nil {
		return err
	}

	switch protocolRoom.Metadata.Direction {
	case protocols.DIRECTION_JOIN:
		fmt.Println("User joining rooms...")
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
				return errors.New("There is no room")
			}

			messages := roomInDatabase.GetLatestMessages()
			fmt.Println("Pushing messages ... ", messages)
			fmt.Println(frame)
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