package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"terrigenesis/fileserver/utils"

	uuid "github.com/satori/go.uuid"
)

/*
EstablishConnection Handles connection establishment
*/
func EstablishConnection(w http.ResponseWriter, r *http.Request) (string, bool) {
	w.Header().Set("Content-Type", "application/json")

	var token uuid.UUID
	// compose header
	w.WriteHeader(200)

	// generate session token
	token = uuid.NewV4()
	m := utils.Response{Status: 200, Token: token.String()}

	// generate response in json
	err := json.NewEncoder(w).Encode(m)
	return token.String(), err == nil
}

/*
CloseConnection Handles connection terminalization
*/
func CloseConnection(w http.ResponseWriter, r *http.Request, sessions []utils.Session) []utils.Session {
	w.Header().Set("Content-Type", "application/json")

	m := utils.Response{Status: 200, Message: "Session not found"}
	w.WriteHeader(404)

	for i := 0; i < len(sessions); i++ {
		if sessions[i].Token == strings.Join(r.URL.Query()["Token"], "") {
			sessions = utils.RemoveFromSlice(sessions, i)
			m = utils.Response{Status: 200, Message: "Successfully closed session"}
			w.WriteHeader(200)
		}
	}

	json.NewEncoder(w).Encode(m)

	return sessions
}
