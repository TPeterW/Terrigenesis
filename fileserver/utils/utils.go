package utils

import "time"

/*
Session Representing an open session between client and server
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

/*
SessionExist Check if a token is for one of the existinng sessions
*/
func SessionExist(sessions []Session, token string) bool {
	for _, session := range sessions {
		if session.Token == token {
			return true
		}
	}
	return false
}
