package internal

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

const (
	ServerAddr   = "127.0.0.1:8080"
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

		err = c.msgHandler.ProcessMessage(msg)
		if err != nil {
			log.Printf("Error handling message: %v", err)
			return
		}
	}
}

func (c *Client) Disconnect() {
	err := c.conn.Close()
	if err != nil {
		log.Printf("Failed to close connection. %v", err)
	}
}
