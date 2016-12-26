package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"terrigenesis/fileserver/handlers"
	"terrigenesis/fileserver/utils"
	"terrigenesis/secrets"
	"time"
)

var sessions []utils.Session
var quit chan struct{}

/*
StartServer Entry point for fileserver
*/
func StartServer() {
	handleInterrupt()

	quit = make(chan struct{})
	startTicker(quit)

	checkAndCreateDBDirectory()

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
	username, password, ok := r.BasicAuth()
	if ok && strings.Compare(username, secrets.Username()) == 0 && strings.Compare(password, secrets.Password()) == 0 {
		fmt.Println(">>> Authentication Passed")
	} else {
		handlers.AuthenticationError(w)
		fmt.Println()
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
	var index int
	var exists bool
	if session, index, exists = utils.SessionExist(sessions, strings.Join(r.URL.Query()["Token"], "")); !exists {
		handlers.SessionNotFoundError(w)
		return sessions
	}

	// now session is available for use
	switch request {
	// Print Working Directory
	case "pwd":
		fmt.Println(">>> Print Working Directory")
		session = handlers.PrintWorkingDirectory(w, r, session)

	// List all files under current directory
	case "dir":
		fmt.Println(">>> List Files")
		session = handlers.ListFiles(w, r, session)

	// Download File
	case "downfile":
		fmt.Println(">>> Download File")
		session = handlers.DownloadFile(w, r, session)

	default:
		handlers.UnknownCommandError(w)
		session.LastUsed = time.Now()
	}

	// replace the session with updated one
	sessions = utils.RemoveFromSlice(sessions, index)
	sessions = append(sessions, session)

	fmt.Printf("\n")
	return sessions
}

func handlePost(w http.ResponseWriter, r *http.Request, request string, sessions []utils.Session) []utils.Session {
	defer r.Body.Close()

	err := r.ParseMultipartForm(32 >> 20)
	if err != nil {
		handlers.IllegalArgumentsError(w)
		return sessions
	}
	token := strings.Join(r.Form["token"], "")

	var session utils.Session
	var index int
	var exists bool
	if session, index, exists = utils.SessionExist(sessions, token); !exists {
		handlers.SessionNotFoundError(w)
		return sessions
	}

	switch request {
	// Change Directory
	case "chdir":
		fmt.Println(">>> Change Directory")
		session = handlers.ChangeDir(w, r, r.Form, session)

	// Make Directory
	case "mkdir":
		fmt.Println(">>> Create Directory")
		session = handlers.MakeDir(w, r, r.Form, session)

	// Remove Directory
	case "rmdir":
		fmt.Println(">>> Remove Directory")
		session = handlers.RemoveDir(w, r, r.Form, session)

	// Upload File
	case "upfile":
		fmt.Println(">>> Upload File")
		session = handlers.UploadFile(w, r, session)

	// Remove File
	case "rmfile":
		fmt.Println(">>> Remove File")
		session = handlers.RemoveFile(w, r, r.Form, session)

	// Move File (does not support rename)
	case "mvfiledir":
		fmt.Println(">>> Move File Or Dir")
		session = handlers.MoveFileDir(w, r, r.Form, session)

	default:
		handlers.UnknownCommandError(w)
		session.LastUsed = time.Now()
	}

	sessions = utils.RemoveFromSlice(sessions, index)
	sessions = append(sessions, session)

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
				close(quit)
				os.Exit(0)
			}
		}
	}()
}

func startTicker(quit chan struct{}) {
	ticker := time.NewTicker(1 * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				for index, session := range sessions {
					// without action for 5 minutes
					if time.Since(session.LastUsed).Minutes() > 5 {
						sessions = utils.RemoveFromSlice(sessions, index)
						fmt.Printf("Session %v timed out\n\n", session)
					}
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func checkAndCreateDBDirectory() {
	info, err := os.Stat("./db")
	if err != nil {
		os.Mkdir("./db", 0744)
	} else {
		if !info.IsDir() {
			os.Remove("./db")
			os.Mkdir("./db", 0744)
		}
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
