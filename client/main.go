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
	pipeClient = pipes.NewNamedPipe(pipePath)
	err := pipeClient.Connect()
	if err != nil {
		fmt.Println("Error: ", err)
	}
	// time.Sleep(500 * time.Millisecond)
	//}
	pipeClient.WriteMessage("hello\n")
}
