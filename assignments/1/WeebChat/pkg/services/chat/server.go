package chat

import (
	"WeebChat/pkg/models"
	"WeebChat/pkg/tcp"
)

type ChatServiceServer struct {
	Host        string
	Port        string
	ChatService models.ChatService
	tcpServer   *tcp.TCPServer
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

	s.tcpServer = &tcp.TCPServer{
		Host:          s.Host,
		Port:          s.Port,
		ClientHandler: s.handlers,
	}

	s.tcpServer.Setup()

	return nil
}

func (s *ChatServiceServer) PingDiscoveryService(discoveryService models.DiscoveryService) {
	//Open connection to discovery service
	//Register service with discovery service protocols
}

func (s *ChatServiceServer) Start() {
	//Start listening for connections
	//Handle connections
}

func (s *ChatServiceServer) Stop() {
	//Stop listening for connections
	//Close connections
}

func (s *ChatServiceServer) handlers(tcpReq tcp.TCPRequest) (tcpRes tcp.TCPResponse, isClose bool) {
	//Parse request
	//Handle request
	//Return response
}
