package pipes

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"syscall"
)

type NamedPipe struct {
	PipePath   string
	Incoming   chan string
	readingEnd *os.File
	writingEnd *os.File
}

func NewNamedPipe(pipeName string) *NamedPipe {

	tmpDir, _ := ioutil.TempDir("", "named-pipes")
	pipePath := path.Join(tmpDir, pipeName)

	np := &NamedPipe{
		PipePath: pipePath,
	}
	return np
}

// func NewNamedPipe(pipename string) *NamedPipe {
// 	np := &NamedPipe{
// 		PipeName: pipename,
// 	}
// 	// conn, err := npipe.Dial(`\\.\pipe\` + np.PipeName)
// 	// if err != nil {
// 	// 	fmt.Fprintf(conn, "Error: %v", err)
// 	// }

// 	return np
// }

func (np *NamedPipe) handleConnection() {
	//	var err error

	for {
		str := np.ReadMessage()
		fmt.Println("handleConnection, got message: ", str)
		np.Incoming <- str
	}
}

// func (np *NamedPipe) waitForOtherSide() {
// 	var err error
// 	for i := 0; i < 500; i++ {
// 		fmt.Println("server: open write end")
// 		np.writingEnd, err = os.OpenFile(np.p2Path, os.O_WRONLY, 0600)
// 		if err != nil {
// 			fmt.Printf("Warn: %v\n", err)
// 		}
// 	}
// }

func (np *NamedPipe) ListenAndServe() error {
	var err error

	fmt.Println("Running IPC server")
	// Create named pipe
	syscall.Mkfifo(np.PipePath+"1", 0600)
	syscall.Mkfifo(np.PipePath+"2", 0600)
	fmt.Println("server: open read end")
	np.readingEnd, err = os.OpenFile(np.PipePath+"1", os.O_RDONLY, 0600)
	if err != nil {
		fmt.Printf("ListenAndServe, Error: %v\n", err)
	}
	go np.handleConnection()
	fmt.Println("server: open write end")
	np.writingEnd, err = os.OpenFile(np.PipePath+"2", os.O_WRONLY, 0600)
	if err != nil {
		fmt.Printf("Warn: %v\n", err)
	}
	// var buff bytes.Buffer
	// fmt.Println("Parent Waiting for initial connection & hello message")
	// io.Copy(&buff, np.readingEnd)

	// if string(buff.Bytes()) != "hello" {
	// 	fmt.Printf("Error: %v\n", err)
	// 	return errors.New("initiation sequence error: invalid handshake")
	// }

	return nil
}

func (np *NamedPipe) Connect() error {
	var err error
	//syscall.Mkfifo(np.p2Path, 0600)
	fmt.Println("client: open write end")
	np.writingEnd, err = os.OpenFile(np.PipePath+"1", os.O_WRONLY, 0600)
	if err != nil {
		fmt.Printf("Warn: %v\n", err)
	}

	fmt.Println("client: open read end")
	np.readingEnd, err = os.OpenFile(np.PipePath+"2", os.O_RDONLY, 0600)
	if err != nil {
		fmt.Printf("Connect, Error: %v\n", err)
		return err
	}

	go np.handleConnection()
	return nil
}

func (np *NamedPipe) WriteMessage(message string) {
	np.writingEnd.Write([]byte(message))
}

func (np *NamedPipe) ReadMessage() string {
	var buff bytes.Buffer

	_, err := io.Copy(&buff, np.readingEnd)
	//msg, err := bufio.NewReader(np.readingEnd).ReadString('\n')
	if err != nil {
		fmt.Printf("ReadMessage, Error: %v\n", err)
	}
	return string(buff.Bytes())
}
