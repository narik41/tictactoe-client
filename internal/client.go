package internal

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/narik41/tictactoe-message/core"
)

const (
	ServerAddr   = "127.0.0.1:9000"
	retryTimeout = 60 * time.Second
	retryDelay   = 2 * time.Second
)

type Client struct {
	conn         net.Conn
	reader       *bufio.Reader
	writer       *bufio.Writer
	mu           sync.Mutex
	symbol       string
	opponent     string
	gameID       string
	myTurn       bool
	gameActive   bool
	board        [9]string
	name         string
	currentBoard [9]string
	mySymbol     string
	ui           ClientUI
}

func NewClient(name string, ui ClientUI) *Client {
	return &Client{
		name:  name,
		board: [9]string{"", "", "", "", "", "", "", "", ""},
		ui:    ui,
	}
}

func (c *Client) Connect(addr string) error {
	fmt.Println("Connecting to server. Server Address ", addr)
	deadline := time.Now().Add(retryTimeout)

	for {
		conn, err := net.Dial("tcp", addr)
		if err == nil {
			c.conn = conn
			c.reader = bufio.NewReader(conn)
			c.writer = bufio.NewWriter(conn)
			fmt.Println("Connected to ", conn.RemoteAddr())
			break
		}

		fmt.Println("Failed to connect to server. Retrying in ", retryDelay)
		if time.Now().After(deadline) {
			panic("Failed to connect within 60 seconds")
		}

		fmt.Println("Connection failed. Retrying in 2 seconds...")
		time.Sleep(retryDelay)
	}
	return nil
}

func (c *Client) Start() {
	defer c.Disconnect()

	newMsgPayload, err := c.doAuth()
	if err != nil {
		log.Printf("Error handling message: %v", err)
		return
	}

	// send login request
	err = c.sendMessage(newMsgPayload)
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return
	}
	rw := bufio.NewReadWriter(bufio.NewReader(c.conn), bufio.NewWriter(c.conn))
	for {

		line, err := rw.ReadString('\n')
		if err != nil {
			log.Printf("Error %v", err)
			return
		}

		message, err := core.DecodeMessage([]byte(line))
		if err != nil {
			log.Printf("Invalid message from: %v", err)
			continue
		}

		jsonBytes, err := json.Marshal(message.Payload)
		if err != nil {
			log.Println(err)
			return
		}

		var v1Msg *core.Version1MessagePayload
		if err := json.Unmarshal(jsonBytes, &v1Msg); err != nil {
			log.Println(err)
		}

		if v1Msg.MessageType == core.MSG_LOGIN_RESPONSE {
			fmt.Println("Waiting for other player to join the game....")
			continue
		}

		if v1Msg.MessageType == core.PLAYER_MOVE {
			jsonBytes, err = json.Marshal(message.Payload)
			if err != nil {
				log.Println(err)
				return
			}

			var playerMovePayload *core.Version1PositionMoveRequestPayload
			if err := json.Unmarshal(jsonBytes, &playerMovePayload); err != nil {
				log.Println(err)
			}
			c.currentBoard[playerMovePayload.Position] = playerMovePayload.Symbol
		}

		if v1Msg.MessageType == core.GAME_START {
			jsonBytes, err = json.Marshal(v1Msg.Payload)
			if err != nil {
				log.Println(err)
				return
			}

			var gameStartPayload *core.Version1GameStartPayload
			if err := json.Unmarshal(jsonBytes, &gameStartPayload); err != nil {
				log.Println(err)
			}
			c.mySymbol = gameStartPayload.YourSymbol
			c.myTurn = gameStartPayload.YourTurn
		} else if v1Msg.MessageType == core.PLAYER_MOVE_RESPONSE {
			jsonBytes, err = json.Marshal(v1Msg.Payload)
			if err != nil {
				log.Println(err)
				return
			}

			var gameStartPayload *core.Version1PositionMovedResponsePayload
			if err := json.Unmarshal(jsonBytes, &gameStartPayload); err != nil {
				log.Println(err)
			}
			c.myTurn = gameStartPayload.TurnSymbol == c.mySymbol
			c.currentBoard[gameStartPayload.MovedToPosition] = gameStartPayload.MovedByUser
			c.ui.DisplayBoard(c.mySymbol, c.board)
		}

		if !c.myTurn {
			continue
		}

		// display board for the user to enter
		c.ui.DisplayBoard(c.mySymbol, c.board)

		move, err := c.ui.PromptForMove()
		if err != nil {
			fmt.Printf("Error prompting for move: %v", err)
			continue
		}
		c.currentBoard[move] = c.mySymbol
		err = c.sendMove(move)
		if err != nil {
			log.Printf("Error sending message: %v", err)
			return
		}
	}
}

func (m *Client) doAuth() (*core.TicTacToeMessage, error) {
	username, password, err := m.ui.PromptCredentials()
	if err != nil {
		return nil, err
	}

	fmt.Printf("\n Authenticating as '%s'...\n", username)

	loginData := &core.TicTacToeMessage{
		MessageId: UUID("msg"),
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

func (c *Client) sendMove(position int) error {
	msg := &core.TicTacToeMessage{
		MessageId: UUID("msg"),
		Version:   "v1",
		Payload: &core.Version1MessagePayload{
			MessageType: core.PLAYER_MOVE,
			Payload: &core.Version1PositionMoveRequestPayload{
				Position: position,
				Symbol:   c.mySymbol,
			},
		},
	}
	return c.sendMessage(msg)
}

func (c *Client) sendMessage(msg *core.TicTacToeMessage) error {

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	msgBytes = append(msgBytes, []byte("\n")...)
	writer := bufio.NewWriter(c.conn)
	if err := json.NewEncoder(writer).Encode(&msgBytes); err != nil {
		log.Printf("Encoding error: %v", err)
		return err
	}

	err = writer.Flush()
	if err != nil {
		log.Printf("Flush error: %v", err)
		return err
	}
	return nil
}

func (c *Client) Disconnect() {
	err := c.conn.Close()
	if err != nil {
		log.Printf("Failed to close connection. %v", err)
	}
}

func UUID(prefix string) string {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		log.Printf("UUID generation error: %v", err)
		return ""
	}
	return fmt.Sprintf("%s-%s", prefix, newUUID.String())
}
