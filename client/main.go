package main

import (
	"flag"
	"fmt"
	"ipc/pipes"
)

func main() {
	flag.Parse()
	pipePath := flag.Args()[0]

	var pipeClient *pipes.NamedPipe
	//for i := 0; i < 3; i++ {
	fmt.Println("Opening ipc as client")
	// set the correct pipe path: (constructor function will concat it with stuff)
	pipeClient = &pipes.NamedPipe{PipePath:pipePath}

	err := pipeClient.Connect()
	if err != nil {
		fmt.Println("Client Main, Error: ", err)
	}
	// time.Sleep(500 * time.Millisecond)
	//}
	fmt.Println("client: writing!")
	pipeClient.WriteMessage("hello\n")
	pipeClient.WriteMessage("hello2\n")
	fmt.Println("client: reading!")
	data:=<-pipeClient.Incoming
	fmt.Println("client got data: ",data)
	 
}
