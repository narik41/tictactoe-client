package internal

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

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
	msgReceiver  MsgReceiver
	msgHandler   MessageHandler
	currentBoard [9]string
	mySymbol     string
}

func NewClient(name string, msgReceiver MsgReceiver, msgHandler MessageHandler) *Client {
	return &Client{
		name:        name,
		board:       [9]string{"", "", "", "", "", "", "", "", ""},
		msgReceiver: msgReceiver,
		msgHandler:  msgHandler,
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

	// Ask username and password
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
	for {

		// display board for the user to enter
		c.displayBoard()

		move, err := c.promptForMove()
		if err != nil {
			fmt.Printf("Error prompting for move: %v", err)
			continue
		}
		c.currentBoard[move] = "X"
		err = c.sendMove(move)
		if err != nil {
			log.Printf("Error sending message: %v", err)
			return
		}
	}
}

func (m *Client) doAuth() (*core.TicTacToeMessage, error) {
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

func (m *Client) promptCredentials() (string, string, error) {
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("TICTACTOE GAME - LOGIN")
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

func (c *Client) sendMove(position int) error {
	msg := &core.TicTacToeMessage{
		MessageId: "",
		Version:   "v1",
		Payload: &core.Version1MessagePayload{
			MessageType: core.PLAYER_MOVE,
			Payload: &core.Version1MessageLoginRequestPayload{
				Position: position,
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

	rw := bufio.NewReadWriter(bufio.NewReader(c.conn), bufio.NewWriter(c.conn))
	if err := json.NewEncoder(rw.Writer).Encode(&msgBytes); err != nil {
		log.Printf("Encoding error: %v", err)
		return err
	}

	err = rw.Flush()
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

func (c *Client) displayBoard() {
	fmt.Println("╔═══════════════╗")
	fmt.Printf("║ You are: %s    ║\n", c.mySymbol)
	fmt.Println("╚═══════════════╝")
	fmt.Println("")

	for i := 0; i < 9; i++ {
		cell := c.currentBoard[i]
		if cell == "" {
			fmt.Printf(" %d ", i) // Show position number
		} else {
			fmt.Printf(" %s ", cell) // Show X or O
		}

		if i%3 == 2 {
			fmt.Println()
			if i < 6 {
				fmt.Println("---|---|---")
			}
		} else {
			fmt.Print("|")
		}
	}
	fmt.Println()
}
