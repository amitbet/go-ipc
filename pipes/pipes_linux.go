package pipes

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"syscall"
	"time"
)

func GetPipePath(pipeName string) string{
	tmpDir, _ := ioutil.TempDir("", "named-pipes")
	pipePath := path.Join(tmpDir, pipeName)
	return pipePath
}

type NamedPipe struct {
	PipePath   string
	Incoming   chan string
	readingEnd *os.File
	writingEnd *os.File
}

func NewNamedPipe(pipePath string) *NamedPipe {
	np := &NamedPipe{
		PipePath: pipePath,
		Incoming : make(chan string),
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

func (np *NamedPipe) handleConnection(processName string) {
	//	var err error
	reader := bufio.NewReader(np.readingEnd)

	for {
		msg, err := reader.ReadBytes('\n')
		if err != nil {
			//if err.Error() != "EOF" {
			fmt.Println(processName+": handleConnection, Error: ", err)
			//}
			time.Sleep(500 * time.Millisecond)
			continue
		}
		fmt.Println(processName+": handleConnection, got message: ", string(msg))
		np.Incoming <- string(msg)
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
	fmt.Println("server: open read end: ", np.PipePath+"1")
	np.readingEnd, err = os.OpenFile(np.PipePath+"1", os.O_RDONLY|syscall.O_NONBLOCK, 0600)
	if err != nil {
		fmt.Printf("ListenAndServe, Error: %v\n", err)
	}
	go np.handleConnection("server")
	fmt.Println("server: open write end: ", np.PipePath+"2")
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
	np.Incoming = make(chan string, 10)

	//syscall.Mkfifo(np.p2Path, 0600)
	fmt.Println("client: open write end: ", np.PipePath+"1")
	np.writingEnd, err = os.OpenFile(np.PipePath+"1", os.O_WRONLY, 0600)
	if err != nil {
		fmt.Printf("Warn: %v\n", err)
	}

	fmt.Println("client: open read end: ", np.PipePath+"2")
	np.readingEnd, err = os.OpenFile(np.PipePath+"2", os.O_RDONLY|syscall.O_NONBLOCK, 0600)
	if err != nil {
		fmt.Printf("Connect, Error: %v\n", err)
		return err
	}

	go np.handleConnection("client")
	return nil
}

func (np *NamedPipe) WriteMessage(message string) {
	np.writingEnd.WriteString(message)
}

func (np *NamedPipe) ReadMessage() string {
	//var buff bytes.Buffer

	//_, err := io.Copy(&buff, np.readingEnd)
	msg, err := bufio.NewReader(np.readingEnd).ReadBytes('\n')

	// var b [5]byte
	// _, err := np.readingEnd.Read(b[:])
	if err != nil {
		fmt.Printf("ReadMessage, Error: %v\n", err)
	}
	fmt.Printf("ReadMessage, got: %s\n", string(msg))
	return string(msg)
}
