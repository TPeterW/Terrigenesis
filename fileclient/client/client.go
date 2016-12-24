package client

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"terrigenesis/secrets"

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
		fmt.Printf("Error reading your password: %v", err)
		os.Exit(1)
	}

	fmt.Println("\nConnnecting to " + secrets.URL())
	var del delegate
	var ok bool
	del.username = username
	del.password = string(bytePasswd)
	del, ok = openConnection(del)
	if !ok {
		os.Exit(1)
	} else {
		handleInterrupt(del)
	}

	doExit := false
	for !doExit {
		reader := bufio.NewReader(os.Stdin)
		for text, _ := reader.ReadString('\n'); strings.Compare(strings.TrimSuffix(text, "\n"), "closecon") != 0; text, _ = reader.ReadString('\n') {
			text = strings.TrimSuffix(text, "\n")
			middleware(text, del)
		}
		fmt.Println("Ending session...")
		if ok := closeConnection(del); ok {
			fmt.Println("Shutting down client...")
			doExit = true
		} else {
			fmt.Printf("Failed to close session\n>>> ")
			doExit = false
		}
	}
}

/*
middleware Does different things depending on the request
*/
func middleware(cmd string, del delegate) {
	// TODO:
	switch cmd {
	case "upfile":
		// TODO:
	// TODO:
	default:
		fmt.Printf("Unrecognized command\n>>> ")
	}
}

func handleInterrupt(del delegate) {
	// handle keyboard interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			if sig != nil {
				fmt.Println("\rEnding session...")
				if ok := closeConnection(del); ok {
					fmt.Println("Shutting down client...")
					os.Exit(0)
				} else {
					fmt.Printf("Failed to close session\n>>> ")
				}
			}
		}
	}()
}
