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
	http.HandleFunc("/", mainHandler)
	http.ListenAndServe(":"+strconv.Itoa(portNum), nil)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("\n", r)

	// root page, render snakey
	if len(r.URL.Path) <= 1 {
		renderSnake(w)
		return
	}

	// basic authentication
	if ok := utils.BasicAuth(r); ok {
		fmt.Println(">>> Authentication Passed")
	} else {
		handlers.AuthenticationError(w)
		return
	}

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
	defer r.Body.Close()

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

	// List all files under current directory
	case "dir":
		fmt.Println(">>> List Files")
		handlers.ListFiles(w, r, session)

	// Download File
	case "downfile":
		fmt.Println(">>> Download File")
		handlers.DownloadFile(w, r, session)

	default:
		handlers.UnknownCommandError(w)
	}

	// none of the actions will modify the session

	fmt.Printf("\n")
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
		handlers.RemoveDir(w, r, body, session)

	// Upload File
	case "upfile":
		fmt.Println(">>> Upload File")
		handlers.UploadFile(w, r, session)

	// Remove File
	case "rmfile":
		fmt.Println(">>> Remove File")
		handlers.RemoveFile(w, r, body, session)

	// Move File (does not support rename)
	case "mvfile":
		fmt.Println(">>> Move File")

	default:
		handlers.UnknownCommandError(w)
	}

	// TODO: replace original session with current one

	fmt.Printf("\n")
	return sessions
}

// render a snake game
func renderSnake(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<!DOCTYPE html><html lang=\"en\"><head><meta charset=\"UTF-8\"><title>Snakey!!!</title><canvas id=\"canvas\" width=\"400\" height=\"400\"></canvas></head><body><script>window.onkeydown=((ctx,snake,food,direction,move,draw)=>((loop,newFood,timer)=>Array.from({length:400}).forEach((_e,i)=>draw(ctx,i,\"black\"))||(timer=setInterval(()=>loop(newFood)||clearInterval(timer)||console.log(timer)||alert('Game Over'),200))&&(e=>direction=snake[1]-snake[0]==(move=[-1,-20,1,20][(e||event).keyCode-37]||direction)?direction:move))((newFood)=>snake.unshift(move=snake[0]+direction)&&snake.indexOf(move,1)>0||move<0||move>399||direction==1&&move%20==0||direction==-1&&move%20==19?false:(draw(ctx,move,\"green\")||move==food?newFood()&draw(ctx,food,\"red\"):draw(ctx,snake.pop(),\"Black\"))!==[],()=>Array.from({length:8000}).some(e=>snake.indexOf(food=~~(Math.random()*400))===-1)))(document.getElementById('canvas').getContext('2d'),[42,41],43,1,null,(ctx,node,color)=>(ctx.fillStyle=color)&ctx.fillRect(node%20*20+1,~~(node/20)*20+1,18,18));</script></body></html>")
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
