package v1

import (
	"errors"

	"github.com/narik41/tictactoe-message/core"
)

type AuthMessageHandler struct{}

func (m *AuthMessageHandler) Version1MessageHandler(v1MsgPayload interface{}) error {
	v1Msg, ok := v1MsgPayload.(core.Version1MessagePayload)
	if !ok {
		return errors.New("v1 payload is not a TicTacToeMessage")
	}

	switch v1Msg.MessageType {
	case core.MSG_LOGIN_REQUEST:
		//return m.doAuth()
	default:
		return errors.New("v1 payload is not a TicTacToeMessage")
	}
	return nil
}
