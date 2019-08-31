package pipes

import (
	"fmt"
	"net"
	"gopkg.in/natefinch/npipe.v2"
	// "bytes"
	// "io"
	"bufio"
)

type NamedPipe struct {
	PipePath string
	conn     net.Conn
	Incoming chan string

}

func NewNamedPipe(pipeName string) *NamedPipe {

	np := &NamedPipe{
		PipePath: `\\.\pipe\` + pipeName,
	}

	return np
}

func (np *NamedPipe) handleConnection() {
	for {
		str := np.ReadMessage()
		fmt.Println("got message: ", str)
		np.Incoming <- str
	}
}

func (np *NamedPipe) ListenAndServe() error {
	ln, err := npipe.Listen(np.PipePath)
	if err != nil {
		fmt.Printf("ListenAndServe, Error: %v\n", err)
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("ListenAndServe2, Error: %v\n", err)
			continue
		}
		np.conn = conn
		go np.handleConnection()
		return nil
	}
}

func (np *NamedPipe) Connect() error {
	fmt.Println("-Connect:",np.PipePath)
	conn, err := npipe.Dial(np.PipePath)
	np.conn = conn
	fmt.Println("-after dial-")
	if err != nil {
		fmt.Printf("Connect, Error: %v\n", err)
		return err
	}
	fmt.Println("-launching read handling routine-")
	go np.handleConnection()
	return nil
}

// func (np *NamedPipe) WriteMessage(message string) {
// 	// w := bufio.NewWriter(np.Conn)
// 	// w.WriteString(message)
// 	strings.Trim(message, "\n")
// 	_, err := fmt.Fprintf(np.conn, message+"\n")
// 	if err != nil {
// 		fmt.Printf("Error: %v\n", err)
// 	}
// }

func (np *NamedPipe) ReadMessage() string {
	msg, err := bufio.NewReader(np.conn).ReadString('\n')
	if err != nil {
		fmt.Printf("ReadMessage, Error: %v\n", err)
	}
	return msg
}


func (np *NamedPipe) WriteMessage(message string) {
	np.conn.Write([]byte(message))
}

// func (np *NamedPipe) ReadMessage() string {
// 	var buff bytes.Buffer

// 	_, err := io.Copy(&buff, np.conn)
// 	//msg, err := bufio.NewReader(np.readingEnd).ReadString('\n')
// 	if err != nil {
// 		fmt.Printf("ReadMessage, Error: %v\n", err)
// 	}
// 	return string(buff.Bytes())
// }