package client

import (
	"fmt"
	"syscall"

	"os"

	"golang.org/x/crypto/ssh/terminal"
)

/*
StartClient Entry point for fileclient
*/
func StartClient(args []string) {
	username := args[0]

	fmt.Print("Please input password: ")
	bytePasswd, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Printf("Error reading you rpassword: %v", err)
		os.Exit(1)
	}

	// TODO: establish connection
	fmt.Println()
	fmt.Println(username)
	fmt.Println(string(bytePasswd))
}
