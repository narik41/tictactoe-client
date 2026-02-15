package internal

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type ClientUI interface {
	DisplayBoard(userSymbol string, board [9]string)
	PromptForMove() (int, error)
	PromptCredentials() (string, string, error)
}

type CMDClientUI struct {
	reader *bufio.Reader
}

func NewCMDClient() ClientUI {
	return &CMDClientUI{
		reader: bufio.NewReader(os.Stdin),
	}
}

func (c CMDClientUI) DisplayBoard(userSymbol string, board [9]string) {
	fmt.Println("╔═══════════════╗")
	fmt.Printf("║ You are: %s    ║\n", userSymbol)
	fmt.Println("╚═══════════════╝")
	fmt.Println("")

	for i := 0; i < 9; i++ {
		cell := board[i]
		if cell == "" {
			fmt.Printf(" %d ", i)
		} else {
			fmt.Printf(" %s ", cell)
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

func (c CMDClientUI) PromptForMove() (int, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n╔════════════════════════╗")
	fmt.Println("║   YOUR TURN!           ║")
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

func (m CMDClientUI) PromptCredentials() (string, string, error) {
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
