package internal

import (
	"fmt"

	"github.com/narik41/tictactoe-message/core"
)

type MessageRouter struct {
	handlers map[core.Version1MessageType]MessageHandler
}
type HandlerResponse struct {
	MessageType core.Version1MessageType
	Payload     interface{}
	Relay       bool
}

type MessageHandler interface {
	Handle(msg *DecodedMessage, session *Client) (*HandlerResponse, error)
}

func NewMessageRouter() *MessageRouter {
	router := &MessageRouter{
		handlers: make(map[core.Version1MessageType]MessageHandler),
	}

	return router
}

func (r *MessageRouter) RegisterHandler(msgType core.Version1MessageType, handler MessageHandler) {
	r.handlers[msgType] = handler
}

func (r *MessageRouter) Route(msg *DecodedMessage, client *Client) (*HandlerResponse, error) {

	handler, exists := r.handlers[msg.MessageType]
	if !exists {
		return nil, fmt.Errorf("unknown message type: %s", msg.MessageType)
	}

	//log.Printf("Routing %s for client %s", msg.MessageType, client.gameID)
	response, err := handler.Handle(msg, client)
	if err != nil {
		return nil, fmt.Errorf("handler failed: %w", err)
	}

	return response, nil
}
