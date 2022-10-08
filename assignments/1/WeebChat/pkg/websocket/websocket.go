package websocket

import (
	"bufio"
	"encoding/binary"
	"io"
	"math"
	"net"
	"net/http"
)

const HASH_KEY_APPEND = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
const BUFFER_SIZE = 1024

type WebSocket struct {
	WebSocketInterface
	Conn   net.Conn
	Rw     *bufio.ReadWriter
	Header http.Header
	Status uint16
}

type WebSocketInterface interface {
	Handshake()
	Write([]byte) error
	Receive() (Frame, error)
	Send(Frame) error
	Close() error
}

func (ws *WebSocket) Write(data []byte) error {
	_, err := ws.Rw.Write(data)
	if err != nil {
		return err
	}

	return ws.Rw.Flush()
}

func (ws *WebSocket) Receive() (frame Frame, err error) {
	head, err := ws.read(2)

	if err != nil {
		return
	}

	frame.IsFragment = (head[0] & 0x80) == 0x00
	frame.Opcode = head[0] & 0x0F
	frame.Reserved = head[0] & 0x70

	frame.IsMasked = (head[1] & 0x80) == 0x80

	var length uint64
	length = uint64(head[1] & 0x7F)

	if length == 126 {
		data, err := ws.read(2)

		if err != nil {
			return frame, err
		}
		length = uint64(binary.BigEndian.Uint16(data))
	} else if length == 127 {
		data, err := ws.read(8)
		if err != nil {
			return frame, err
		}
		length = uint64(binary.BigEndian.Uint64(data))
	}

	mask, err := ws.read(4)

	if err != nil {
		return
	}

	frame.Length = length

	payload, err := ws.read(int(length)) // possible data loss

	if err != nil {
		return frame, err
	}

	for i := uint64(0); i < length; i++ {
		payload[i] ^= mask[i%4]
	}

	frame.Payload = payload

	err = frame.Validate(ws)

	return

}

func (ws *WebSocket) read(byteToRead int) ([]byte, error) {
	data := make([]byte, 0)

	for {
		if len(data) == byteToRead {
			break
		}
		// Temporary slice to read chunk
		bufferSize := BUFFER_SIZE
		remaining := byteToRead - len(data)
		if bufferSize > remaining {
			bufferSize = remaining
		}
		temp := make([]byte, bufferSize)

		n, err := ws.Rw.Read(temp)
		if err != nil && err != io.EOF {
			return data, err
		}

		data = append(data, temp[:n]...)
	}
	return data, nil

}

func (ws *WebSocket) Send(frame Frame) error {
	data := make([]byte, 2)
	data[0] = 0x80 | frame.Opcode
	if frame.IsFragment {
		data[0] &= 0x7F
	}

	if frame.Length <= 125 {
		data[1] = byte(frame.Length)
		data = append(data, frame.Payload...)
	} else if frame.Length > 125 && float64(frame.Length) < math.Pow(2, 16) {
		data[1] = byte(126)
		size := make([]byte, 2)
		binary.BigEndian.PutUint16(size, uint16(frame.Length))
		data = append(data, size...)
		data = append(data, frame.Payload...)
	} else if float64(frame.Length) >= math.Pow(2, 16) {
		data[1] = byte(127)
		size := make([]byte, 8)
		binary.BigEndian.PutUint64(size, frame.Length)
		data = append(data, size...)
		data = append(data, frame.Payload...)
	}
	return ws.Write(data)
}

func (ws *WebSocket) Close() error {
	f := Frame{}
	f.Opcode = 8
	f.Length = 2
	f.Payload = make([]byte, 2)
	binary.BigEndian.PutUint16(f.Payload, ws.Status)
	if err := ws.Send(f); err != nil {
		return err
	}
	return ws.Conn.Close()
}
