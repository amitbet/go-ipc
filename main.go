package main

import (
	"fmt"
	"ipc/pipes"
	"os"
	"os/exec"
)

var execerPath string

func main() {

	var execerPath = "./client/client.exe"

	// Create named pipe
	pipeName := "stam"
	fmt.Println("Opening ipc as server")

	np := pipes.NewNamedPipe(pipeName)
	go np.ListenAndServe()

	go func() {
		cmd := exec.Command(execerPath, np.PipePath)
		// Just to forward the stdout
		cmd.Stdout = os.Stdout
		//fmt.Println("running command: " + execerPath + " " + namedPipe)
		cmd.Run()
	}()

	str := <-np.Incoming
	fmt.Println("got message: ", str)

}
