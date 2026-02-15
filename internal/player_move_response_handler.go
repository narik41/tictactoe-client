package internal

import (
	"encoding/json"

	"github.com/narik41/tictactoe-client/internal/decoder"
	"github.com/narik41/tictactoe-client/internal/templates"
	"github.com/narik41/tictactoe-message/core"
)

type PlayerMoveResponseHandler struct {
	cmdUI templates.UI
}

func NewPlayerMoveResponseHandler(cmdUI templates.UI) PlayerMoveResponseHandler {
	return PlayerMoveResponseHandler{
		cmdUI: cmdUI,
	}
}

func (a PlayerMoveResponseHandler) Handle(msg *decoder.DecodedMessage, client *Client) (*HandlerResponse, error) {
	jsonBytes, err := json.Marshal(msg.Payload)
	if err != nil {
		return nil, err
	}

	var playerMoveResponse core.Version1PositionMovedResponsePayload
	if err := json.Unmarshal(jsonBytes, &playerMoveResponse); err != nil {
		return nil, err
	}

	client.myTurn = playerMoveResponse.TurnSymbol == client.mySymbol
	client.board[playerMoveResponse.MovedToPosition] = playerMoveResponse.MovedByUser

	a.cmdUI.DisplayBoard(client.mySymbol, client.board)
	if client.myTurn {
		move, err := a.cmdUI.PromptForMove()
		if err != nil {
			return nil, err
		}
		client.board[move] = client.mySymbol
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
