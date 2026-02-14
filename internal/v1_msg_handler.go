package internal

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
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
		 
	case core.MSG_LOGIN_RESPONSE:
		// sessionId will be provided from server
		// and openont player too

	default:
		return nil, errors.New("v1 payload is not a TicTacToeMessage")
	}
	return nil, nil
}

func (m *Version1MsgHandler) loginRsp(v1Msg *core.Version1MessagePayload) (*core.TicTacToeMessage, error) {
	fmt.Printf("Successfully Login!!! Searching for the players")

	return nil, nil
}

func (c *Client) promptForMove() (int, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("╔════════════════════════╗")
	fmt.Println("║      YOUR TURN!        ║")
	fmt.Println("╚════════════════════════╝")
	fmt.Print("Enter position (0-8): ")

	input, err := reader.ReadString('\n')
	if err != nil {
		return -1, err
	}

	input = strings.TrimSpace(input)
	position, err := strconv.Atoi(input)
	if err != nil {
		return -1, fmt.Errorf("invalid input, enter a number 0-8")
	}

	if position < 0 || position > 8 {
		return -1, fmt.Errorf("position must be between 0 and 8")
	}

	return position, nil
}
