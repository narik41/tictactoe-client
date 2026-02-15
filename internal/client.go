package internal

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"

	decoder2 "github.com/narik41/tictactoe-client/internal/decoder"
)

const (
	ServerAddr   = "127.0.0.1:9000"
	retryTimeout = 60 * time.Second
	retryDelay   = 2 * time.Second
)

type Client struct {
	conn     net.Conn
	myTurn   bool
	board    [9]string
	mySymbol string
	router   *MessageRouter
}

func NewClient(router *MessageRouter) *Client {
	return &Client{
		board:  [9]string{"", "", "", "", "", "", "", "", ""},
		router: router,
	}
}

func (c *Client) Connect(addr string) error {
	fmt.Println("Connecting to server. Server Address ", addr)
	deadline := time.Now().Add(retryTimeout)

	for {
		conn, err := net.Dial("tcp", addr)
		if err == nil {
			c.conn = conn
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

	responseSender := NewResponseSender()
	rw := bufio.NewReadWriter(bufio.NewReader(c.conn), bufio.NewWriter(c.conn))
	decoder := decoder2.NewMessageDecoder(rw)
	for {
		decodedMessage, err := decoder.Decode()
		if err != nil {
			log.Printf("Error %v", err)
			return
		}

		response, err := c.router.Route(decodedMessage, c)
		if err != nil {
			//sender.SendError(s, "HANDLER_ERROR", err2.Error())
			continue
		}

		if !response.Relay {
			continue
		}

		err = responseSender.Send(c, response)
		if err != nil {
			continue
		}
	}
}

func (c *Client) Disconnect() {
	err := c.conn.Close()
	if err != nil {
		log.Printf("Failed to close connection. %v", err)
	}
}
