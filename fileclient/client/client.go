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
		for text, _ := reader.ReadString('\n'); strings.Compare(strings.TrimSuffix(text, "\n"), "closecon") != 0 && strings.Compare(strings.TrimSuffix(text, "\n"), "exit") != 0; text, _ = reader.ReadString('\n') {
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
	commands := strings.Split(cmd, " ")
	switch commands[0] {
	case "":
		break
	case "help":
		fmt.Println("pwd\t\t: []\t\t\t- print working directory")
		fmt.Println("ch, chdir\t: [dirname]\t\t- change into directory")
		fmt.Println("ls, dir\t\t: []\t\t\t- list files under current directory")
		fmt.Println("mkdir\t\t: [dirname]\t\t- create new directory")
		fmt.Println("rmdir\t\t: [dirname]\t\t- remove a directory")
		fmt.Println("rm, rmfile\t: [filename]\t\t- remove a file")
		fmt.Println("mv, mvfiledir\t: [filename] [destination] - moves a file or directory (no renaming)")
		fmt.Println("download\t: [filename]\t\t- download a file to local Downloads folder")
		fmt.Println("upload\t\t: [pathToFile]\t\t- upload a local file to server")
		// TODO:
	case "pwd":
		printWorkingDirectory(del)

	case "ls":
		fallthrough
	case "dir":
		if len(commands) > 1 {
			fmt.Println("Too many arguments for \"" + commands[1] + "\"")
		} else {
			listFiles(del)
		}

	case "download":
		fallthrough
	case "downfile":
		if len(commands) < 2 {
			fmt.Println("Missing filename")
		} else {
			downloadFile(del, strings.Join(commands[1:], " "))
		}

	case "upload":
		fallthrough
	case "upfile":
		if len(commands) < 2 {
			fmt.Println("Missing filename")
		} else {
			uploadFile(del, strings.Join(commands[1:], " "))
		}

	case "cd":
		fallthrough
	case "chdir":
		if len(commands) < 2 {
			fmt.Println("Too few arguments for \"" + commands[0] + "\"")
		} else if len(commands) > 2 {
			fmt.Println("Too many arguments")
		} else {
			changeDir(del, commands[1])
		}

	case "mkdir":
		if len(commands) < 2 {
			fmt.Println("Too few arguments for \"mkdir\"")
		} else if len(commands) > 2 {
			fmt.Println("Too many arguments for \"mkdir\"")
		} else {
			makeDir(del, commands[1])
		}

	case "rmdir":
		if len(commands) < 2 {
			fmt.Println("Too few arguments for \"rmdir\"")
		} else if len(commands) > 2 {
			fmt.Println("Too many arguments for \"rmdir\"")
		} else {
			removeDir(del, commands[1])
		}

	case "rm":
		fallthrough
	case "rmfile":
		if len(commands) < 2 {
			fmt.Println("Too few arguments for \"" + commands[0] + "\"")
		} else {
			removeFile(del, commands[1:])
		}

	case "mv":
		fallthrough
	case "mvfiledir":
		if len(commands) < 3 {
			fmt.Println("Too few arguments for \"" + commands[0] + "\"")
		} else {
			moveFileOrDir(del, commands[1:], commands[0])
		}

	default:
		fmt.Println("Unrecognized command")
	}
	fmt.Printf(">>> ")
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
					fmt.Printf("Failed to close session\n")
					os.Exit(1)
				}
			}
		}
	}()
}
