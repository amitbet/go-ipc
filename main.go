package main

import (
	"fmt"
	"ipc/pipes"
)

var execerPath string

func main() {

	go func() {
		fmt.Println("Opening ipc as client")
		pipeClient := pipes.NewNamedPipe("stam")
		pipeClient.Connect()
		pipeClient.WriteMessage("hello/n")
	}()

	fmt.Println("Opening ipc as server")
	pipe := pipes.NewNamedPipe("stam")
	go pipe.ListenAndServe()
	str := <-pipe.Incoming
	fmt.Println("got message: ", str)

}
