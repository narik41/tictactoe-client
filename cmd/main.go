package main

import (
	"fmt"
	"log"

	"github.com/narik41/tictactoe-client/internal"
	"github.com/narik41/tictactoe-client/internal/templates"
	"github.com/narik41/tictactoe-message/core"
)

func main() {
	fmt.Println("!!! Starting the tic tac toe !!!")
	consoleUI := templates.NewConsoleUI()

	router := internal.NewMessageRouter()
	router.RegisterHandler(core.MSG_LOGIN_RESPONSE, internal.NewLoginResponseHandler(consoleUI))
	router.RegisterHandler(core.GAME_START, internal.NewGameStartHandler(consoleUI))
	router.RegisterHandler(core.GAME_END, internal.NewGameEndHandler(consoleUI))
	router.RegisterHandler(core.PLAYER_MOVE_RESPONSE, internal.NewPlayerMoveResponseHandler(consoleUI))
	router.RegisterHandler(core.MSG_LOGIN_REQUEST, internal.NewLoginRequestHandler(consoleUI))

	client := internal.NewClient(router)
	err := client.Connect(internal.ServerAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect()
	client.Start()
}
