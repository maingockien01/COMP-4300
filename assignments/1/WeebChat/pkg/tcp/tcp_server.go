package tcp

import {
	"net"
	"os"
}

const (
	BUFFER_SIZE = 1024
	CONN_TYPE = "tcp"
)

func handleIncomingRequest (conn net.Conn) {
	buffer := make([]byte, BUFFER_SIZE)
	err := conn.Read(buffer)

	if err != nil {
		log.Fatal(err)
	}

}

type TcpServer struct {
	Host			string
	Port			string
	HandleAcceptError 	func(error)
	HandleRequest 		func(net.Conn)
}

func NewTcpServer (host string, port string, handleAcceptError func(error), handleRequest func(net.Conn)) *TcpServer {
	tcpServer := TcpServer {
		Host: host,
		Port: port,
		HandleAcceptError: handleAcceptError,
		HandleRequest: handleRequest
	}

	return &tcpServer

}

func (s *TcpServer) Start () {
	listener, err := net.Listen(CONN_TYPE, s.Host+":"+s.Port)
	if err != nil {
		fmt.Println("Error listening: ", err.Error())
		os.Exit()
	}

	defer listener.Close()

	fmt.Println("Opened listener on " + host + ":" + port)

	for {
		conn, err:= listener.Accept()

		if err != nil {
			fmt.Println("Error on accepting: ", err.Error())
			s.HandleAcceptError(err)
		}

		go s.HandeRequest(conn)
	}
}

func ExtractData (conn net.Conn) ([]byte, error) {
	buffer := make([]byte, BUFFER_SIZE)
	_, err := conn.Read(buffer)
	return buffer, err
}