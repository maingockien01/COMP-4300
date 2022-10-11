package server

import (
	"WeebChat/pkg/services/protocols"
	"WeebChat/pkg/websocket"
	"encoding/json"
	"fmt"
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
		s.ChatService.AddUser(&protocolUser.Data)
		return nil
	}

	return nil
}
