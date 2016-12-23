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
	Token    string // session token
	Filename string // filename to move or delete
	Dirname  string // dirname to move or delete
	Dest     string // destination of moving
	File     string // file to upload
}

/*
SessionExist Check if a token is for one of the existinng sessions
*/
func SessionExist(sessions []Session, token string) (Session, int, bool) {
	for index, session := range sessions {
		if session.Token == token {
			return session, index, true
		}
	}
	return Session{}, -1, false // return empty session not to be used
}

/*
RemoveFromSlice Remove a session from a slice
*/
func RemoveFromSlice(s []Session, i int) []Session {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}
