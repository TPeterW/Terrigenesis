package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"terrigenesis/fileserver/utils"
)

/*
EstablishConnection Handles connection establishment
*/
func EstablishConnection(w http.ResponseWriter, r *http.Request) (string, bool) {
	var m utils.Message
	w.Header().Set("Content-Type", "application/json")

	var username string
	var ok bool
	if username, ok = utils.BasicAuth(r); ok {
		fmt.Println(">>> Authentication Passed")
		// compose header
		w.WriteHeader(200)
		m = utils.Message{Status: 200, Message: "ok"}
		j, err := json.Marshal(m)
		if err != nil {
			fmt.Println("ERR ", err)
		} else {
			w.Write(j)
			return username, true
		}
	} else {
		fmt.Println(">>> Cannot authenticate user")
		w.WriteHeader(401)
		m = utils.Message{Status: 401, Message: "Authorization Error"}
		j, err := json.Marshal(m)
		if err != nil {
			fmt.Println("ERR ", err)
		} else {
			w.Write(j)
		}
	}

	return username, false
}
