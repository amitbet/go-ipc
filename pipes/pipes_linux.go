package pipes

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"syscall"
)

type NamedPipe struct {
	PipeName   string
	p1Path     string
	p2Path     string
	Incoming   chan string
	readingEnd *os.File
	writingEnd *os.File
}

func NewNamedPipe(pipename string) *NamedPipe {
	tmpDir, _ := ioutil.TempDir("", "named-pipes")
	p1path := path.Join(tmpDir, pipename+"-P1")
	p2path := path.Join(tmpDir, pipename+"-P2")

	np := &NamedPipe{
		PipeName: pipename,
		p1Path:   p1path,
		p2Path:   p2path,
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
	for {
		str := np.ReadMessage()
		fmt.Println("got message: ", str)
		np.Incoming <- str
	}
}

func (np *NamedPipe) ListenAndServe() {
	var err error

	fmt.Println("Running IPC server")
	// Create named pipe
	syscall.Mkfifo(np.p1Path, 0600)
	syscall.Mkfifo(np.p2Path, 0600)
	np.readingEnd, err = os.OpenFile(np.p1Path, os.O_RDONLY, 0600)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	np.writingEnd, err = os.OpenFile(np.p2Path, os.O_WRONLY, 0600)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	go np.handleConnection()
}

func (np *NamedPipe) Connect() {
	var err error
	fmt.Println("Opening named pipe for reading")
	np.readingEnd, err = os.OpenFile(np.p2Path, os.O_RDONLY, 0600)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	np.writingEnd, err = os.OpenFile(np.p1Path, os.O_WRONLY, 0600)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	go np.handleConnection()
}

func (np *NamedPipe) WriteMessage(message string) {
	w := bufio.NewWriter(np.readingEnd)
	w.WriteString(message)

	// strings.Trim(message, "\n")
	// _, err := fmt.Fprintf(np.conn, message+"\n")
	// if err != nil {
	// 	fmt.Printf("Error: %v", err)
	// }
}

func (np *NamedPipe) ReadMessage() string {
	// reader := bufio.NewReader(np.readingEnd)

	msg, err := bufio.NewReader(np.readingEnd).ReadString('\n')
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	return msg
}
