package internal

import (
	"encoding/json"
	"log"

	"github.com/narik41/tictactoe-message/core"
)

type GameEndHandler struct {
	cmdUI ClientUI
}

func NewGameEndHandler(cmdUI ClientUI) GameEndHandler {
	return GameEndHandler{
		cmdUI: cmdUI,
	}
}

func (a GameEndHandler) Handle(msg *DecodedMessage, client *Client) (*HandlerResponse, error) {
	log.Println("GameEndHandler: Handling the end start msg type")
	jsonBytes, err := json.Marshal(msg.Payload)
	if err != nil {
		return nil, err
	}

	var gameStartPayload core.Version1GameEndPayload
	if err := json.Unmarshal(jsonBytes, &gameStartPayload); err != nil {
		return nil, err
	}

	client.mySymbol = gameStartPayload.Winner
	a.cmdUI.DisplayBoard(client.mySymbol, client.board)
	if client.myTurn {

	}

	return &HandlerResponse{
		Relay: false,
	}, nil
}
