package client

import (
	"fmt"
	"os/signal"
	"syscall"

	"os"

	"golang.org/x/crypto/ssh/terminal"
)

/*
StartClient Entry point for fileclient
*/
func StartClient(args []string) {
	handleInterrupt()

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

func handleInterrupt() {
	// handle keyboard interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			if sig != nil {
				fmt.Println("\rEnding session...")
				// TODO: close conneciton first

				fmt.Println("Shutting down client...")
				os.Exit(0)
			}
		}
	}()
}
