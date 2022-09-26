package tcp

import (
	e "WeebChat/pkg/err"
	"bufio"
	"fmt"
	"net"
)

const BUFFER_SIZE = 4096

func HandlerCloseClient(conn net.Conn) {
	conn.Close()
}

func HandlerWriteBack(conn net.Conn, tcpResponse TCPResponse) error {
	//TODO: check if conn is open before write ??? do I need to check
	_, err := conn.Write(tcpResponse.Response)

	e.LogError(err)

	return err
}

func HandlerReadBuffer(conn net.Conn) (TCPRequest, error) {
	//TODO: handle read large message bigger than buffer?
	var returnErr error

	request := make([]byte, 0)
	requestSize := 0

	buffer := make([]byte, BUFFER_SIZE)
	reader := bufio.NewReader(conn)

	// if reader.Size() == 0 {
	// 	break
	// }

	n, err := reader.Read(buffer)

	if err != nil {
		returnErr = err
		// break
	}

	request = append(request, buffer[:n]...)

	requestSize += n

	// if n < BUFFER_SIZE {
	// 	break
	// }

	tcpRequest := TCPRequest{
		Request: request,
		Size:    requestSize,
	}

	fmt.Println(requestSize)

	e.LogError(returnErr)

	return tcpRequest, returnErr

}
