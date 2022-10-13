package server

import (
	"WeebChat/pkg/models"
	"WeebChat/pkg/services/protocols"
	"WeebChat/pkg/tcp"
	"WeebChat/pkg/websocket"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
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
	router.Handle("/echo", http.HandlerFunc(s.HandlerEcho))
	router.Handle("/chat", http.HandlerFunc(HandlerWrapper(s.HandlerFrame)))
	router.Handle("/rooms", http.HandlerFunc(s.HandlerGetRooms))
	router.Handle("/users/room/", http.HandlerFunc(s.HandleGetUserIn))
	router.Handle("/user", http.HandlerFunc(s.HandleAddUser))
	router.Handle("/room", http.HandlerFunc(s.HandlerRoom))
	router.Handle("/messages/", http.HandlerFunc(s.HandleMessageRestful))

	server.Handler = router

	s.Server = server

	s.ChatService.Rooms = append(s.ChatService.Rooms, &models.Room{
		Messages: make([]*models.Message, 0),
		Name:     "HelloWorld",
		Users:    make([]*models.User, 0),
		Lock:     sync.Mutex{},
		Capacity: 0,
		Limit:    10,
	})

	s.ChatService.Rooms[0].AppendMessage(&models.User{Id: "00", Name: "default", Ws: nil, LastActiveAt: time.Now()}, models.NewMessage("00", "Hello World"))

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
	fmt.Println("Handling frame...")
	payload := frame.ParseText()
	var protocolMetada protocols.Protocol
	err := json.Unmarshal([]byte(payload), &protocolMetada)

	fmt.Println(protocolMetada)

	if err != nil {
		return err
	}

	fmt.Println("Handling ", protocolMetada.Metadata.Type, " ...")
	switch protocolMetada.Metadata.Type {
	case protocols.TYPE_MESSAGE:
		return s.HandleMessage(payload, frame, ws)
	case protocols.TYPE_ROOM:
		return s.HandleRoom(payload, frame, ws)
	case protocols.TYPE_USER:
		return s.HandleUser(payload, frame, ws)
	default:
		fmt.Println("Found no match type handler")
		return errors.New("Errors on unexpected type handler")
	}

}

func (s *ChatServiceServer) PushMessage(ws *websocket.ServerWebSocket, messages ...*models.Message) error {
	if ws == nil {
		return errors.New("There is no socket")
	}
	protocolMessagePush := protocols.ProtocolMessage{
		Metadata: protocols.ProtocolMetadata{
			Version:   protocols.V1,
			From:      s.ChatService.Name,
			Direction: protocols.DIRECTION_PUSH,
			Type:      protocols.TYPE_MESSAGE,
		},
		Data: messages,
	}

	jsonPayload, err := json.Marshal(protocolMessagePush)
	if err != nil {
		fmt.Println("Error on jsonifying message for pushing message ", err)
		return err
	}

	frame := websocket.NewFrameMessage(jsonPayload)

	fmt.Println("---")
	fmt.Println(*frame)

	err = ws.Send(*frame)

	if err != nil {
		fmt.Println("Error on pushing message into socket ", err)
		return err
	}

	return nil
}
