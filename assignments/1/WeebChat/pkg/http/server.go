package http

import (
	"WeebChat/pkg/tcp"
)

type HttpServer struct {
	tcpServer tcp.TCPServer
	Host      string
	Port      string
	Handler   func(HttpRequest, HttpResponse) (HttpResponse, error)
}

func (server *HttpServer) Setup() error {

	if server.Host == "" || server.Port == "" {
		return tcp.MissingRequiredField{}
	}

	server.tcpServer = tcp.TCPServer{
		Host:          server.Host,
		Port:          server.Port,
		ClientHandler: server.handlers,
	}

	server.tcpServer.Setup()

	return nil
}

func (server *HttpServer) Start() error {
	server.tcpServer.Start()
	return nil
}

func (server *HttpServer) Stop() {
	server.tcpServer.Stop()
}

func (server *HttpServer) handlers(tcpReq tcp.TCPRequest) (tcpRes tcp.TCPResponse, isClose bool) {
	//TODO: to be orgainized
	isClose = true

	req, err := parseHttpRequest(tcpReq.Request)

	if err != nil {
		res := handleParsingHttpError(req, err)
		tcpRes = parseTcpResponse(res)

		return
	}

	res := initHttpResponse(req)

	if err != nil {
		_, res := handleInternalError(req, res)
		tcpRes = parseTcpResponse(res)
		return
	}

	//Abstract: webserver/application layer will in charge of this
	if server.Handler == nil {
		_, res := handleInternalError(req, res)
		tcpRes = parseTcpResponse(res)
		return
	}

	res, err = server.Handler(req, res)

	if err != nil {
		_, res := handleInternalError(req, res)
		tcpRes = parseTcpResponse(res)
		return
	}

	tcpRes = parseTcpResponse(res)
	return tcpRes, true //always close connection after handling request (handle keep alive case later, if needed)
}
