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
	conn     net.Conn
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
	for {
		str := np.ReadMessage()
		fmt.Println("got message: ", str)
		np.Incoming <- str
	}
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
		np.conn = conn
		go np.handleConnection(conn)
	}
}

func (np *NamedPipe) Connect() {
	conn, err := npipe.Dial(`\\.\pipe\` + np.PipeName)
	np.conn = conn
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	go np.handleConnection()
}

func (np *NamedPipe) WriteMessage(message string) {
	// w := bufio.NewWriter(np.Conn)
	// w.WriteString(message)
	strings.Trim(message, "\n")
	_, err := fmt.Fprintf(np.conn, message+"\n")
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
}

func (np *NamedPipe) ReadMessage() string {
	msg, err := bufio.NewReader(np.conn).ReadString('\n')
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	return msg
}
