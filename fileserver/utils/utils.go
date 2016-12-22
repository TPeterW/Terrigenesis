package utils

import "time"

/*
Session An open session between client and server
*/
type Session struct {
	Token    string    // token
	CWD      string    // current working directory
	LastUsed time.Time // last action time
}

/*
Response Response format from server to client
*/
type Response struct {
	Status   int
	Message  string
	Token    string
	CWD      string
	DirFiles []string
}

/*
PostBody Format for body of post request
*/
type PostBody struct {
	Token string // session token
}

/*
SessionExist Check if a token is for one of the existinng sessions
*/
func SessionExist(sessions []Session, token string) (Session, bool) {
	for _, session := range sessions {
		if session.Token == token {
			return session, true
		}
	}
	return Session{}, false // return empty session not to be used
}
