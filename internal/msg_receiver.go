package internal

import (
	"fmt"
	"log"

	"github.com/narik41/tictactoe-message/core"
)

type MsgReceiver struct {
}

func NewMsgReceiver() MsgReceiver {
	return MsgReceiver{}
}

func (m *MsgReceiver) Receive(c *Client) (*core.TicTacToeMessage, error) {
	//log.Printf("Message received from server: %s", c.conn.RemoteAddr())
	line, err := c.reader.ReadBytes('\n')
	if err != nil {
		log.Printf("Error reading from server: %v", err)
		return &core.TicTacToeMessage{}, fmt.Errorf("receive error: %w", err)
	}

	message, err := core.DecodeMessage(line)
	if err != nil {
		log.Printf("Error unmarshalling message: %v", err)
		return &core.TicTacToeMessage{}, fmt.Errorf("decode error: %w", err)
	}

	return message, err
}
