package v1

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/narik41/tictactoe-message/core"
)

type Version1MsgHandler struct{}

func NewVersion1MsgHandler() Version1MsgHandler {
	return Version1MsgHandler{}
}

func (m *Version1MsgHandler) Version1MessageHandler(v1MsgPayload interface{}) (*core.TicTacToeMessage, error) {
	jsonBytes, err := json.Marshal(v1MsgPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal: %w", err)
	}

	// Unmarshal to struct
	var v1Msg *core.Version1MessagePayload
	if err := json.Unmarshal(jsonBytes, &v1Msg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}

	switch v1Msg.MessageType {
	case core.MSG_LOGIN_REQUEST:
		return m.doAuth(v1Msg)
	case core.MSG_LOGIN_RESPONSE:

	default:
		return nil, errors.New("v1 payload is not a TicTacToeMessage")
	}
	return nil, nil
}

func (m *Version1MsgHandler) loginRsp(v1Msg *core.Version1MessagePayload) (*core.TicTacToeMessage, error) {
	fmt.Printf("Successfully Login!!! Searching for the players")

	return nil, nil
}

func (m *Version1MsgHandler) doAuth(v1Msg *core.Version1MessagePayload) (*core.TicTacToeMessage, error) {
	username, password, err := m.promptCredentials()
	if err != nil {
		return nil, err
	}

	fmt.Printf("\n Authenticating as '%s'...\n", username)

	loginData := &core.TicTacToeMessage{
		MessageId: "",
		Version:   "v1",
		Payload: &core.Version1MessagePayload{
			MessageType: core.MSG_LOGIN_PAYLOAD,
			Payload: &core.Version1MessageLoginPayload{
				Username: username,
				Password: password,
			},
		},
	}

	return loginData, nil
}

func (m *Version1MsgHandler) promptCredentials() (string, string, error) {
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("           TICTACTOE GAME - LOGIN")
	fmt.Println(strings.Repeat("=", 50))
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nUsername: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		return "", "", fmt.Errorf("failed to read username: %w", err)
	}
	username = strings.TrimSpace(username)

	if username == "" {
		return "", "", fmt.Errorf("username cannot be empty")
	}

	// Ask for password
	fmt.Print("Password: ")
	password, err := reader.ReadString('\n')
	if err != nil {
		return "", "", fmt.Errorf("failed to read password: %w", err)
	}
	password = strings.TrimSpace(password)

	if password == "" {
		return "", "", fmt.Errorf("password cannot be empty")
	}

	fmt.Println(strings.Repeat("=", 50))

	return username, password, nil
}
