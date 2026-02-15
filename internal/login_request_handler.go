package internal

import (
	"encoding/json"
	"fmt"

	"github.com/narik41/tictactoe-client/internal/decoder"
	"github.com/narik41/tictactoe-client/internal/templates"
	"github.com/narik41/tictactoe-message/core"
)

type LoginRequestHandler struct {
	clientUi templates.UI
}

func NewLoginRequestHandler(clientUi templates.UI) LoginRequestHandler {
	return LoginRequestHandler{
		clientUi: clientUi,
	}
}

func (a LoginRequestHandler) Handle(msg *decoder.DecodedMessage, client *Client) (*HandlerResponse, error) {
	jsonBytes, err := json.Marshal(msg.Payload)
	if err != nil {
		return nil, err
	}

	var gameStartPayload core.Version1GameStartPayload
	if err := json.Unmarshal(jsonBytes, &gameStartPayload); err != nil {
		return nil, err
	}

	return a.DisplayLoginForm()
}

func (a LoginRequestHandler) DisplayLoginForm() (*HandlerResponse, error) {

	username, password, err := a.clientUi.PromptCredentials()
	if err != nil {
		return nil, err
	}

	fmt.Printf("\n Authenticating as '%s'...\n", username)
	v1MsgPayload := &core.Version1MessagePayload{
		MessageType: core.MSG_LOGIN_PAYLOAD,
		Payload: &core.Version1MessageLoginPayload{
			Username: username,
			Password: password,
		},
	}

	return &HandlerResponse{
		Relay:       true,
		Payload:     v1MsgPayload,
		MessageType: core.MSG_LOGIN_PAYLOAD,
	}, nil
}
