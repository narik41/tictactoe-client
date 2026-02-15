package main

import (
	"fmt"
	"log"

	"github.com/narik41/tictactoe-client/internal"
	"github.com/narik41/tictactoe-message/core"
)

func main() {
	fmt.Println("!!! Starting the tic tac toe !!!")
	cmdUI := internal.NewCMDClient()

	router := internal.NewMessageRouter()
	router.RegisterHandler(core.MSG_LOGIN_RESPONSE, internal.NewLoginResponseHandler(cmdUI))
	router.RegisterHandler(core.GAME_START, internal.NewGameStartHandler(cmdUI))
	router.RegisterHandler(core.GAME_END, internal.NewGameEndHandler(cmdUI))
	router.RegisterHandler(core.PLAYER_MOVE_RESPONSE, internal.NewPlayerMoveResponseHandler(cmdUI))
	router.RegisterHandler(core.MSG_LOGIN_REQUEST, internal.NewLoginRequestHandler(cmdUI))

	client := internal.NewClient(router)
	err := client.Connect(internal.ServerAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect()
	client.Start()
}
