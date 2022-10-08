package server

import (
	"WeebChat/pkg/services/protocols"
	"WeebChat/pkg/websocket"
	"encoding/json"
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
				return err
			}
		}

		return nil
	}

	return nil
}
