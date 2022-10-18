package websocket

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const HANDSHAKE_STATUS = 101

type ServerWebSocket struct {
	WebSocket
}

func (ws *ServerWebSocket) Handshake() error {
	hash := getSecWebsocketAcceptHash(ws.Header.Get("Sec-Websocket-Key"))

	headerLines := []string{
		"HTTP/1.1 101 Web Socket Protocol Handshake",
		"Server: go/websocket",
		"Upgrade: WebSocket",
		"Connection: Upgrade",
		"Sec-WebSocket-Accept: " + hash,
		"", // required for extra CRLF
		"", // required for extra CRLF
	}

	return ws.Write([]byte(strings.Join(headerLines, "\r\n")))
}

func getSecWebsocketAcceptHash(key string) string {
	sha := sha1.New()
	sha.Write(([]byte(key)))
	sha.Write([]byte(HASH_KEY_APPEND))

	return base64.StdEncoding.EncodeToString(sha.Sum(nil))
}

func NewServerWebSocket(w http.ResponseWriter, req *http.Request) (*ServerWebSocket, error) {
	hj, ok := w.(http.Hijacker)
	if !ok {
		return nil, errors.New("webserver doesn't support http hijacking")

	}

	conn, rw, err := hj.Hijack()

	if err != nil {
		return nil, err
	}

	return &ServerWebSocket{
		WebSocket{
			Conn:   conn,
			Rw:     rw,
			Header: req.Header,
			Status: 1000,
		},
	}, nil
}

func WebSocketHandler(w http.ResponseWriter, req *http.Request, frameHandler func(Frame, *ServerWebSocket) error) {
	var ws *ServerWebSocket
	ws, err := NewServerWebSocket(w, req)

	if err != nil {
		log.Println("Error on opening websocket: ", err)
		return
	}

	err = ws.Handshake()

	if err != nil {
		log.Println("Error on init handshake: ", err)
		return
	}

	for {
		frame, err := ws.Receive()

		if err != nil {
			fmt.Println("Error on receiving frame: ", err)
		}

		frameHandler(frame, ws)

	}
}
