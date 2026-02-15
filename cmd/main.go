package main

import (
	"fmt"
	"log"

	"github.com/narik41/tictactoe-client/internal"
)

func main() {
	fmt.Println("!!! Starting the tic tac toe client !!!")

	cmdUI := internal.NewCMDClient()
	client := internal.NewClient("tictactoe", cmdUI)
	err := client.Connect(internal.ServerAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect()
	client.Start()
}
