package internal

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/narik41/tictactoe-message/core"
)

type LoginResponseHandler struct {
	cmdUI ClientUI
}

func NewLoginResponseHandler(cmdUI ClientUI) LoginResponseHandler {
	return LoginResponseHandler{
		cmdUI: cmdUI,
	}
}

func (a LoginResponseHandler) Handle(msg *DecodedMessage, client *Client) (*HandlerResponse, error) {

	jsonBytes, err := json.Marshal(msg.Payload)
	if err != nil {
		return nil, err
	}

	var loginResponsePayload core.Version1MessageLoginResponse
	if err := json.Unmarshal(jsonBytes, &loginResponsePayload); err != nil {
		return nil, err
	}

	if loginResponsePayload.IsAuthenticated {
		fmt.Println("Login Successful!!! Waiting for opponent !!!!!")
		fmt.Println("=================================================")
		return &HandlerResponse{
			Relay: false,
		}, nil
	}
	return nil, errors.New("Login Failed")
}
