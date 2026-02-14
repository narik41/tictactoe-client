package main

import (
	"fmt"
	"log"

	"github.com/narik41/tictactoe-client/internal"
)

func main() {
	fmt.Println("!!! Starting the tic tac toe client !!!")

	v1MsgHandler := internal.NewVersion1MsgHandler()
	msgHandler := internal.NewMessageHandler(v1MsgHandler)
	msgReceiver := internal.NewMsgReceiver()
	client := internal.NewClient("tictactoe", msgReceiver, msgHandler)
	err := client.Connect(internal.ServerAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect()
	client.Start()
}
