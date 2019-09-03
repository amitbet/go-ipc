package main

import (
	"fmt"
	"ipc/pipes"
	"os"
	"os/exec"
	"time"
)

var execerPath string

func main() {

	var execerPath = "./client/client"

	// Create named pipe
	pipeName := "stam"
	fmt.Println("Opening ipc as server")
	ppath:=pipes.GetPipePath(pipeName)
	np := pipes.NewNamedPipe(ppath)
	go np.ListenAndServe()
	fmt.Println("-running client")
	go func() {
		cmd := exec.Command(execerPath, np.PipePath)
		// Just to forward the stdout
		cmd.Stdout = os.Stdout
		//fmt.Println("running command: " + execerPath + " " + namedPipe)
		cmd.Run()
	}()
	fmt.Println("server: reading!")
	str := <-np.Incoming
	fmt.Println("server, got message1: ", str)
	str = <-np.Incoming
	fmt.Println("server, got message2: ", str)
	np.WriteMessage("yo client!\n")
	np.WriteMessage("yo client2!\n")
	time.Sleep(5 * time.Minute)
}
