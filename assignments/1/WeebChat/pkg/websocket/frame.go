package websocket

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"unicode/utf8"
)

var CLOSE_CODES map[int]string = map[int]string{
	1000: "NormalError",
	1001: "GoingAwayError",
	1002: "ProtocolError",
	1003: "UnknownType",
	1007: "TypeError",
	1008: "PolicyError",
	1009: "MessageTooLargeError",
	1010: "ExtensionError",
	1011: "UnexpectedError",
}

const (
	// TextMessage denotes a text data message. The text message payload is
	// interpreted as UTF-8 encoded text data.
	MESSAGE_TEXT = 1

	// BinaryMessage denotes a binary data message.
	MESSAGE_BINARY = 2

	// CloseMessage denotes a close control message. The optional message
	// payload contains a numeric code and text. Use the FormatCloseMessage
	// function to format a close message payload.
	MESSAGE_CLOSE = 8

	// PingMessage denotes a ping control message. The optional message payload
	// is UTF-8 encoded text.
	MESSAGE_PING = 9

	// PongMessage denotes a pong control message. The optional message payload
	// is UTF-8 encoded text.
	MESSAGE_PONG = 10
)

type Frame struct {
	IsFragment bool
	Opcode     byte
	Reserved   byte
	IsMasked   bool
	Length     uint64
	Payload    []byte
}

func NewFramePong() (f *Frame) {
	f.Opcode = MESSAGE_PONG
	return f
}

func NewFramePing() (f *Frame) {
	f.Opcode = MESSAGE_PING
	return f
}

func NewFrameMessage(message []byte) (f *Frame) {
	f.Opcode = MESSAGE_TEXT
	f.Payload = message

	return f
}

func (f *Frame) ParseText() string {
	return string(f.Payload)
}

func (f *Frame) IsControl() bool {
	return f.Opcode&0x08 == 0x08
}

func (f *Frame) HasReservedOpcode() bool {
	return f.Opcode > 10 || (f.Opcode >= 3 && f.Opcode <= 7)
}

func (f *Frame) ClodeCode() uint16 {
	var code uint16
	binary.Read(bytes.NewReader(f.Payload), binary.BigEndian, &code)
	return code
}

func (f *Frame) Validate(ws *WebSocket) error {
	if !f.IsMasked {
		ws.Status = 1002
		return errors.New("protocol error: unmasked client frame")
	}
	if f.IsControl() && (f.Length > 125 || f.IsFragment) {
		ws.Status = 1002
		return errors.New("protocol error: all control frames MUST have a payload length of 125 bytes or less and MUST NOT be fragmented")
	}
	if f.HasReservedOpcode() {
		ws.Status = 1002
		return errors.New("protocol error: opcode " + fmt.Sprintf("%x", f.Opcode) + " is reserved")
	}
	if f.Reserved > 0 {
		ws.Status = 1002
		return errors.New("protocol error: RSV " + fmt.Sprintf("%x", f.Reserved) + " is reserved")
	}
	if f.Opcode == 1 && !f.IsFragment && !utf8.Valid(f.Payload) {
		ws.Status = 1007
		return errors.New("wrong code: invalid UTF-8 text message ")
	}
	if f.Opcode == 8 {
		if f.Length >= 2 {
			code := binary.BigEndian.Uint16(f.Payload[:2])
			reason := utf8.Valid(f.Payload[2:])
			if code >= 5000 || (code < 3000 && CLOSE_CODES[int(code)] == "") {
				ws.Status = 1002
				return errors.New(CLOSE_CODES[1002] + " Wrong Code")
			}
			if f.Length > 2 && !reason {
				ws.Status = 1007
				return errors.New(CLOSE_CODES[1007] + " invalid UTF-8 reason message")
			}
		} else if f.Length != 0 {
			ws.Status = 1002
			return errors.New(CLOSE_CODES[1002] + " Wrong Code")
		}
	}
	return nil
}
