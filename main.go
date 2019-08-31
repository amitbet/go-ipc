package main

import (
	"fmt"
	"io/ioutil"
	"ipc/pipes"
	"os"
	"os/exec"
	"path"
)

var execerPath string

func main() {

	var execerPath = "./client/client"

	// Create named pipe
	pipeName := "stam"
	fmt.Println("Opening ipc as server")

	tmpDir, _ := ioutil.TempDir("", "named-pipes")
	pipePath := path.Join(tmpDir, pipeName)

	pipe := pipes.NewNamedPipe(pipePath)
	go pipe.ListenAndServe()

	go func() {
		cmd := exec.Command(execerPath, pipePath)
		// Just to forward the stdout
		cmd.Stdout = os.Stdout
		//fmt.Println("running command: " + execerPath + " " + namedPipe)
		cmd.Run()
	}()

	str := <-pipe.Incoming
	fmt.Println("got message: ", str)

}
