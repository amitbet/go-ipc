package pipes

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"gopkg.in/natefinch/npipe.v2"
)

type NamedPipe struct {
	PipeName string
	Conn     net.Conn
	Incoming chan string
}

func NewNamedPipe(pipename string) *NamedPipe {
	np := &NamedPipe{
		PipeName: pipename,
	}
	// conn, err := npipe.Dial(`\\.\pipe\` + np.PipeName)
	// if err != nil {
	// 	fmt.Fprintf(conn, "Error: %v", err)
	// }

	return np
}

func (np *NamedPipe) handleConnection(conn net.Conn) {
	str := np.ReadMessage()
	fmt.Println("got message: ", str)
	np.Incoming <- str
}

func (np *NamedPipe) ListenAndServe() {
	ln, err := npipe.Listen(`\\.\pipe\` + np.PipeName)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Error: %v", err)
			continue
		}
		np.Conn = conn
		go np.handleConnection(conn)
	}
}

func (np *NamedPipe) Connect() {
	conn, err := npipe.Dial(`\\.\pipe\` + np.PipeName)
	np.Conn = conn
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
}

func (np *NamedPipe) WriteMessage(message string) {
	// w := bufio.NewWriter(np.Conn)
	// w.WriteString(message)
	strings.Trim(message, "\n")
	_, err := fmt.Fprintf(np.Conn, message+"\n")
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
}

func (np *NamedPipe) ReadMessage() string {
	//reader := bufio.NewReader(np.Conn)

	msg, err := bufio.NewReader(np.Conn).ReadString('\n')
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	return msg
}
