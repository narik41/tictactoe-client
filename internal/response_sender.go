package internal

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"

	"github.com/narik41/tictactoe-message/core"
)

type ResponseSender struct {
}

func NewResponseSender() *ResponseSender {
	return &ResponseSender{}
}

func (rs *ResponseSender) Send(client *Client, response *HandlerResponse) error {

	msgBytes, err := rs.encodeMessage(response.Payload)
	if err != nil {
		return fmt.Errorf("failed to encode: %w", err)
	}

	if err := rs.writeToConn(client.conn, msgBytes); err != nil {
		return fmt.Errorf("failed to send: %w", err)
	}

	return nil
}

func (rs *ResponseSender) SendError(session *Client, errorCode, errorMessage string) error {
	errorPayload := &core.Version1MessagePayload{
		MessageType: core.ERROR,
		Payload: map[string]interface{}{
			"code":    errorCode,
			"message": errorMessage,
		},
	}

	msgBytes, err := rs.encodeMessage(errorPayload)
	if err != nil {
		return fmt.Errorf("failed to encode error: %w", err)
	}

	if err := rs.writeToConn(session.conn, msgBytes); err != nil {
		return fmt.Errorf("failed to send error: %w", err)
	}

	return nil
}

func (rs *ResponseSender) encodeMessage(payload interface{}) ([]byte, error) {

	msg := core.TicTacToeMessage{
		MessageId: core.UUID("msg"),
		Version:   "v1",
		Timestamp: core.GetNPTToUtcInMillisecond(),
		Payload:   payload,
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	msgBytes = append(msgBytes, '\n')

	return msgBytes, nil
}

func (rs *ResponseSender) writeToConn(conn net.Conn, data []byte) error {
	writer := bufio.NewWriter(conn)

	n, err := writer.Write(data)
	if err != nil {
		return err
	}

	if n != len(data) {
		return fmt.Errorf("incomplete write: wrote %d of %d bytes", n, len(data))
	}

	if err := writer.Flush(); err != nil {
		return err
	}

	return nil
}
