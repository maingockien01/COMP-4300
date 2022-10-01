package discovery

import (
	. "WeebChat/pkg/models"
	"net/http"
	"time"
)

type DiscoveryServiceServer struct {
	Host             string
	Port             string
	DiscoveryService *DiscoveryService
	Server           *http.Server
}

func NewDiscoveryServiceServer(host string, port string) *DiscoveryServiceServer {
	server := DiscoveryServiceServer{
		Host:             host,
		Port:             port,
		DiscoveryService: NewDiscoveryService(),
	}

	return &server
}

func (s *DiscoveryServiceServer) Setup() error {
	//Setup server
	server := &http.Server{
		Addr:         s.Host + ":" + s.Port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	//Setup handlers
	router := http.NewServeMux()
	router.Handle("/services/chat", http.HandlerFunc(s.HandlerServicesChat))
	server.Handler = router

	s.Server = server
	return nil
}

func (s *DiscoveryServiceServer) Start() error {
	//Start listening for connections
	return s.Server.ListenAndServe()
}
