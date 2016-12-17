package utils

import (
	"time"
)

/*
Session representing an open session between client and server
*/
type Session struct {
	Token    string    // token
	CWD      string    // current working directory
	LastUsed time.Time // last action time
}

/*
Message Response format from server to client
*/
type Message struct {
	Status  int
	Message string
	Token   string
}
