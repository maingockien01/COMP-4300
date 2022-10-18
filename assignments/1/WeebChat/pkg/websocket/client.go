package websocket

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
)

type ClientWebSocket struct {
	WebSocket
	url *url.URL
	//Cookies??
	//Context??
}

func NewClientWebSocket(url *url.URL, header http.Header) (*ClientWebSocket, error) {
	address := "localhost:8010"

	fmt.Println("Open websocket to ", address)

	dialer := &net.Dialer{}

	conn, err := dialer.Dial("tcp", address)

	if err != nil {
		return nil, err
	}

	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	ws := &ClientWebSocket{
		url: url,
		WebSocket: WebSocket{
			Conn:   conn,
			Rw:     rw,
			Header: header,
			Status: 1000,
		},
	}

	err = ws.Handshake()

	if err != nil {
		return nil, err
	}

	return ws, nil
}

func (ws *ClientWebSocket) Handshake() error {
	url := ws.url

	fmt.Println("Handshake with ", url)

	req := &http.Request{
		ProtoMajor: 1,
		ProtoMinor: 1,
		Method:     http.MethodGet,
		URL:        url,
		Proto:      "HTTP/1.1",
		Header:     make(http.Header),
		Host:       url.Host,
	}
	req.Header.Add("Connection", "Upgrade")
	req.Header.Add("Upgrade", "websocket")
	req.Header.Add("Host", url.Host)

	challengeKey, err := getChallengeKey()

	if err != nil {
		return err
	}

	req.Header.Add("Sec-WebSocket-Key", challengeKey)

	err = req.Write(ws.Rw)
	if err != nil {
		return err
	}

	res, err := http.ReadResponse(ws.Rw.Reader, req)

	if err != nil {
		return err
	}

	if res.StatusCode != 101 ||
		!headersContainsValue(res.Header, "Upgrade", "websocket") ||
		!headersContainsValue(res.Header, "Connection", "upgrade") ||
		res.Header.Get("Sec-Websocket-Accept") != getSecWebsocketAcceptHash(challengeKey) {
		return errors.New("bad handshake")
	}

	return nil
}

func getChallengeKey() (string, error) {
	p := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, p); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(p), nil
}

// TODO: test this
func getAddress(url *url.URL) string {

	return url.String()
}
