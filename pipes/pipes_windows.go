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

func (np *NamedPipe) ListenAndServe()error {
	ln, err := npipe.Listen(`\\.\pipe\` + np.PipeName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		np.conn = conn
		go np.handleConnection(conn)
		return nil
	}
}

func (np *NamedPipe) Connect() error {
	conn, err := npipe.Dial(`\\.\pipe\` + np.PipeName)
	np.conn = conn
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}
	go np.handleConnection()
	return nil
}

func (np *NamedPipe) WriteMessage(message string) {
	// w := bufio.NewWriter(np.Conn)
	// w.WriteString(message)
	strings.Trim(message, "\n")
	_, err := fmt.Fprintf(np.conn, message+"\n")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func (np *NamedPipe) ReadMessage() string {
	msg, err := bufio.NewReader(np.conn).ReadString('\n')
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	return msg
}
