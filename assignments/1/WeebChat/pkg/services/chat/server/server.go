package server

import (
	"WeebChat/pkg/models"
	"WeebChat/pkg/services/protocols"
	"WeebChat/pkg/tcp"
	"WeebChat/pkg/websocket"
	"encoding/json"
	"net/http"
	"time"
)

type ChatServiceServer struct {
	Host        string
	Port        string
	ChatService models.ChatService
	Ws          websocket.ServerWebSocket
	Server      *http.Server
}

func NewChatServiceServer(host string, port string, service_name string) *ChatServiceServer {
	server := ChatServiceServer{
		Host: host,
		Port: port,
		ChatService: models.ChatService{
			Name: service_name,
		},
	}

	return &server
}

func (s *ChatServiceServer) Setup() error {

	if s.Host == "" || s.Port == "" {
		return tcp.MissingRequiredField{}
	}
	server := &http.Server{
		Addr:         s.Host + ":" + s.Port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	//Setup handlers
	router := http.NewServeMux()

	//Add handlers
	router.Handle("/echo/{id}", http.HandlerFunc(s.HandlerEcho))
	router.Handle("/chat", http.HandlerFunc(HandlerWrapper(s.HandlerFrame)))
	server.Handler = router

	s.Server = server
	return nil
}

func (s *ChatServiceServer) Start() error {
	//Start listening for connections
	return s.Server.ListenAndServe()
}

func (s *ChatServiceServer) PingDiscoveryService(discoveryService models.DiscoveryService) {
	//Open connection to discovery service
	//Register service with discovery service protocols
}

func (s *ChatServiceServer) Stop() {
	//Stop listening for connections
	//Close connections
	s.Server.Close()
}

func HandlerWrapper(handlerFrame func(websocket.Frame, *websocket.ServerWebSocket) error) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		websocket.WebSocketHandler(w, r, handlerFrame)
	}
}

func (s *ChatServiceServer) HandlerFrame(frame websocket.Frame, ws *websocket.ServerWebSocket) error {
	payload := frame.ParseText()
	var protocolMetada protocols.ProtocolMetadata
	err := json.Unmarshal([]byte(payload), &protocolMetada)

	if err != nil {
		return err
	}

	switch protocolMetada.Type {
	case protocols.TYPE_MESSAGE:
		var protocolMessage protocols.ProtocolMessage
		err := json.Unmarshal([]byte(payload), &protocolMessage)

		if err != nil {
			return err
		}
		return nil
	case protocols.TYPE_ROOM:
		return s.HandleRoom(payload, frame, ws)
	case protocols.TYPE_USER:
		return s.HandleUser(payload, frame, ws)
	}

	return nil
}
