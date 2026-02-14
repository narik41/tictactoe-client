package internal

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
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
	conn        net.Conn
	reader      *bufio.Reader
	writer      *bufio.Writer
	mu          sync.Mutex
	symbol      string
	opponent    string
	gameID      string
	myTurn      bool
	gameActive  bool
	board       [9]string
	name        string
	msgReceiver MsgReceiver
	msgHandler  MessageHandler
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
	for {
		msg, err := c.msgReceiver.Receive(c)
		if err != nil {
			log.Println("Error reading from server:", err)
			return
		}

		newMsgPayload, err := c.msgHandler.ProcessMessage(msg)
		if err != nil {
			log.Printf("Error handling message: %v", err)
			return
		}
		if newMsgPayload == nil {
			return
		}

		err = c.sendMessage(newMsgPayload)
		if err != nil {
			return
		}
	}
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
