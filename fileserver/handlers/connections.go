package handlers

import (
	"encoding/json"
	"fmt"
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
	if ok := utils.BasicAuth(r); ok {
		fmt.Println(">>> Authentication Passed")
		// compose header
		w.WriteHeader(200)

		// generate session token
		token := uuid.NewV4()
		m := utils.Message{Status: 200, Token: token.String()}

		// generate response in json
		j, err := json.Marshal(m)
		if err != nil {
			fmt.Println("ERR ", err)
		} else {
			w.Write(j)
			return token.String(), true
		}
	} else {
		AuthenticationError(w)
	}

	return token.String(), false
}

/*
CloseConnection Handles connection terminalization
*/
func CloseConnection(w http.ResponseWriter, r *http.Request, sessions []utils.Session) {
	w.Header().Set("Content-Type", "application/json")

	if ok := utils.BasicAuth(r); ok {
		fmt.Println(">>> Authentication Passed")

		m := utils.Message{Status: 200, Message: "Session not found"}
		w.WriteHeader(404)

		for i := 0; i < len(sessions); i++ {
			if sessions[i].Token == strings.Join(r.URL.Query()["token"], "") {
				removeFromSlice(sessions, i)
				m = utils.Message{Status: 200, Message: "Successfully closed session"}
				w.WriteHeader(200)
			}
		}

		j, err := json.Marshal(m)
		if err != nil {
			fmt.Println("ERR ", err)
		} else {
			w.Write(j)
		}
	} else {
		AuthenticationError(w)
	}
}

func removeFromSlice(s []utils.Session, i int) []utils.Session {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}
