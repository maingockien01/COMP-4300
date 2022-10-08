package websocket

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
)

type ClientWebSocket struct {
	WebSocket
	url *url.URL
	//Cookies??
	//Context??
}

func NewClientWebSocket(url *url.URL, header http.Header) (*ClientWebSocket, error) {
	address := getAddress(url)

	conn, err := net.Dial("TCP", address)

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

	return ws, nil
}

func (ws *ClientWebSocket) Handshake() error {
	url := ws.url
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
		return errors.New("Bad handshake")
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

//TODO: test this
func getAddress(url *url.URL) string {
	addr := url.Host
	if i := strings.LastIndex(url.Host, ":"); i > strings.LastIndex(url.Host, "]") {

	} else {
		switch url.Scheme {
		case "wss":
			addr += ":443"
		case "https":
			addr += ":443"
		default:
			addr += ":80"
		}
	}
	return addr
}
