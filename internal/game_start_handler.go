package internal

import (
	"encoding/json"

	"github.com/narik41/tictactoe-message/core"
)

type GameStartHandler struct {
	cmdUI ClientUI
}

func NewGameStartHandler(cmdUI ClientUI) GameStartHandler {
	return GameStartHandler{
		cmdUI: cmdUI,
	}
}

func (a GameStartHandler) Handle(msg *DecodedMessage, client *Client) (*HandlerResponse, error) {
	//log.Println("GameStartHandler: Handling the game start msg type")
	jsonBytes, err := json.Marshal(msg.Payload)
	if err != nil {
		return nil, err
	}

	var gameStartPayload core.Version1GameStartPayload
	if err := json.Unmarshal(jsonBytes, &gameStartPayload); err != nil {
		return nil, err
	}

	client.mySymbol = gameStartPayload.YourSymbol
	client.myTurn = gameStartPayload.YourTurn
	a.cmdUI.DisplayBoard(client.mySymbol, client.board)
	if client.myTurn {
		move, err := a.cmdUI.PromptForMove()
		if err != nil {
			return nil, err
		}

		v1MsgPayload := &core.Version1MessagePayload{
			MessageType: core.PLAYER_MOVE,
			Payload: &core.Version1PositionMoveRequestPayload{
				Position: move,
				Symbol:   client.mySymbol,
			},
		}

		return &HandlerResponse{
			Relay:       true,
			Payload:     v1MsgPayload,
			MessageType: core.PLAYER_MOVE,
		}, nil

	}

	return &HandlerResponse{
		Relay: false,
	}, nil
}
