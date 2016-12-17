package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"terrigenesis/fileserver/handlers"
	"terrigenesis/fileserver/utils"
	"time"
)

var sessions []utils.Session

/*
StartServer Entry point for fileserver
*/
func StartServer() {
	handleInterrupt()

	// initialize session list
	sessions := make([]utils.Session, 0)

	// port number
	portNum := 3000

	fmt.Println("Listening on port " + strconv.Itoa(portNum))
	fmt.Printf("Current sessions %v\n", sessions)
	http.HandleFunc("/", handler)
	http.ListenAndServe(":"+strconv.Itoa(portNum), nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("\n", r)

	splited := strings.Split(r.URL.Path[1:], "/")

	if r.Method == "GET" {
		switch request := splited[0]; request {
		// Close Connection
		case "closecon":
			fmt.Println(">>> Closing Connection")
			sessions = handlers.CloseConnection(w, r, sessions)
			fmt.Printf("Current sessions: %v\n\n", sessions)

		default:
			sessions = handleGet(w, r, request, sessions)
		}
	} else if r.Method == "POST" {
		switch request := splited[0]; request {
		// Establish Connection
		case "estabcon":
			fmt.Println(">>> Establish Connection")
			if token, ok := handlers.EstablishConnection(w, r); ok {
				sessions = append(sessions, utils.Session{Token: token, CWD: "./db", LastUsed: time.Now()})
			}
			fmt.Printf("Current sessions: %v\n\n", sessions)

		default:
			sessions = handlePost(w, r, request, sessions)
		}
	}
}

func handleGet(w http.ResponseWriter, r *http.Request, request string, sessions []utils.Session) []utils.Session {
	if r.URL.Query()["Token"] == nil {
		handlers.IllegalArgumentsError(w)
		return sessions
	}
	var session utils.Session
	var exists bool
	if session, exists = utils.SessionExist(sessions, strings.Join(r.URL.Query()["Token"], "")); !exists {
		handlers.SessionNotFoundError(w)
		return sessions
	}

	// now session is available for use
	switch request {
	// Print Working Directory
	case "pwd":
		fmt.Println(">>> Print Working Directory")
		handlers.PrintWorkingDirectory(w, r, session)

	// Download File
	case "downfile":
		fmt.Println(">>> Download File")

	default:
		// TODO: render a snake game
	}

	// TODO: replace original session with current one

	return sessions
}

func handlePost(w http.ResponseWriter, r *http.Request, request string, sessions []utils.Session) []utils.Session {
	defer r.Body.Close()

	var body utils.PostBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		handlers.IllegalArgumentsError(w)
		return sessions
	}
	var session utils.Session
	var exists bool
	if session, exists = utils.SessionExist(sessions, body.Token); !exists {
		handlers.SessionNotFoundError(w)
		return sessions
	}

	switch request {
	// Change Directory
	case "chdir":
		fmt.Println(">>> Change Directory")

	// Make Directory
	case "mkdir":
		fmt.Println(">>> Create Directory")

	// Remove Directory
	case "rmdir":
		fmt.Println(">>> Remove Directory")

	// Upload File
	case "upfile":
		fmt.Println(">>> Upload File")
		handlers.UploadFile(w, r, session)

	// Remove File
	case "rmfile":
		fmt.Println(">>> Remove File")

	// Move File (does not support rename)
	case "mvfile":
		fmt.Println(">>> Move File")

	default:
		// TODO: render a snake game
	}

	// TODO: replace original session with current one

	return sessions
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
