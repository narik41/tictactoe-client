package v1

import (
	"errors"

	"github.com/narik41/tictactoe-message/core"
)

type Version1MsgHandler struct{}

func NewVersion1MsgHandler() Version1MsgHandler {
	return Version1MsgHandler{}
}

func (m *Version1MsgHandler) Version1MessageHandler(v1MsgPayload interface{}) error {
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

//
//func (m *MessageHandler) doAuth() error {
//	log.Println("Server requires authentication.")
//	stdin := bufio.NewReader(os.Stdin)
//
//	log.Println("Enter username and password to login in the server.")
//	log.Println("Username: ")
//	fmt.Print(" ")
//	username, err := stdin.ReadString('\n')
//	if err != nil {
//		return fmt.Errorf("failed to read username: %w", err)
//	}
//
//	log.Println("Password: ")
//	fmt.Print(" ")
//	password, err := stdin.ReadString('\n')
//	if err != nil {
//		return fmt.Errorf("failed to read password: %w", err)
//	}
//
//	//return c.send(core.Message{
//	//	MessageType:     core.MsgAuth,
//	//	IsAuthenticated: false,
//	//	Payload: core.AuthMessagePayload{
//	//		Payload: core.AuthPayload{
//	//			Username: strings.TrimSpace(username),
//	//			Password: strings.TrimSpace(password),
//	//		},
//	//		Type: core.MsgAuthRequest,
//	//	},
//	//})
//}
