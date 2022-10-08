package server

import (
	"WeebChat/pkg/websocket"
	"fmt"
	"log"
	"net/http"
)

func (s *ChatServiceServer) HandlerEcho(w http.ResponseWriter, r *http.Request) {
	websocket.WebSocketHandler(w, r, Echo)
}

func Echo(frame websocket.Frame, ws *websocket.ServerWebSocket) error {
	fmt.Println("Received request")
	switch frame.Opcode {
	case 8: // Close
		return nil
	case 9: // Ping
		frame.Opcode = 10
		fallthrough
	case 0: // Continuation
		fallthrough
	case 1: // Text
		fallthrough
	case 2: // Binary
		if err := ws.Send(frame); err != nil {
			log.Println("Error sending", err)
			return err
		}
	}

	return nil
}
