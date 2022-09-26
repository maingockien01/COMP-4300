package http

import (
	"WeebChat/pkg/tcp"
)

//Throws HttpRequestInvalidError
func parseHttpRequest(reqByte []byte) (req HttpRequest, err error) {
	req, err = ParseRequest(string(reqByte))

	return req, err
}

func handleParsingHttpError(req HttpRequest, err error) (res HttpResponse) {
	//TODO: get parsing error
	return res
}

func initHttpResponse(req HttpRequest) (res HttpResponse) {
	//TODO: to be implemented
	res.HttpVersion = req.HttpVersion
	res.Request = req

	return res
}

func handleInternalError(req HttpRequest, res HttpResponse) (HttpRequest, HttpResponse) {
	return HttpRequest{}, HttpResponse{}
}

func parseTcpResponse(res HttpResponse) (tcpRes tcp.TCPResponse) {
	tcpRes.Response = []byte(res.String())

	return tcpRes
}
