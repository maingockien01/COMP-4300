/*
 * TCP server will be in charge of transportation layer working with tcp connection
 */

package tcp

import (
	e "WeebChat/pkg/err"
	"net"
	"time"
)

const DEADLINE_DURATION_READ = 2 * time.Second
const DEADLINE_DURATION_WRITE = 20 * time.Second

type TCPServer struct {
	listener      *net.TCPListener
	Host          string // Ex: "127.0.0.1"
	Port          string // Ex: "8080"
	ClientHandler func(TCPRequest) TCPResponse
}

type TCPRequest struct {
	Request []byte
	Size    int
}

type TCPResponse struct {
	Response []byte
}

func (server *TCPServer) Setup() {
	address := server.Host + ":" + server.Port

	tcpAddr, err := net.ResolveTCPAddr("tcp4", address)
	e.HandleFatalError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	e.HandleFatalError(err)

	server.listener = listener
}

func (server *TCPServer) Stop() {

	server.listener.Close()
}

func (server *TCPServer) Start() {
	for {
		conn, err := server.listener.Accept()

		e.LogError(err)
		if err != nil {
			continue
		}

		go server.handleClient(conn)
	}

}

func (server *TCPServer) handleClient(conn net.Conn) {
	//conn.SetReadDeadline(time.Now().Add(DEADLINE_DURATION_READ))
	//conn.SetWriteDeadline(time.Now().Add(DEADLINE_DURATION_WRITE))
	//TODO check if conn is still alive
	tcpRequest, err := HandlerReadBuffer(conn)
	e.LogError(err)
	if err != nil {
		e.HandleFatalError(err)
	}

	tcpResponse := server.ClientHandler(tcpRequest)

	err = HandlerWriteBack(conn, tcpResponse)
	e.LogError(err)
	e.HandleFatalError(err)

	HandlerCloseClient(conn)
}
