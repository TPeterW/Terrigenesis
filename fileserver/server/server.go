package server

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"terrigenesis/fileserver/utils"
)

var sessions []utils.Session

/*
StartServer Entry point for fileserver
*/
func StartServer() {
	handleInterrupt()

	// initialize session list
	// sessions := make([]utils.Session, 0)

	// port number
	portNum := 3000

	fmt.Println("Listening on port " + strconv.Itoa(portNum))
	http.HandleFunc("/", handler)
	http.HandleFunc("/monkey", monkeyHandler)
	http.ListenAndServe(":"+strconv.Itoa(portNum), nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Println("Get request")
	} else if r.Method == "POST" {
		fmt.Println("Post request")
	}
	splited := strings.Split(r.URL.Path[1:], "/")
	request := splited[0]

	switch request := splited[0]; request {
	// Establish Connection
	case "estabcon":
	// Close Connection
	case "closecon":
	// Print Working Directory
	case "pwd":
	// Change Directory
	case "chdir":
	// Make Directory
	case "mkdir":
	// Remove Directory
	case "rmdir":
	// Upload File
	case "upfile":
	// Download File
	case "downfile":
	// Remove File
	case "rmfile":
	// Move File (does not support rename)
	case "mvfile":
	default:
		// TODO: return error message
	}

	filename := splited[1]

	fmt.Println(request)
	fmt.Println(filename)
}

func monkeyHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is a monkey")
}

func handleInterrupt() {
	// handle keyboard interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			if sig != nil {
				fmt.Println("\rShutting down server...")
				os.Exit(0)
			}
		}
	}()
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
