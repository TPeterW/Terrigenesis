package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"terrigenesis/fileserver/utils"
)

/*
AuthenticationError Handles all authentication errors
*/
func AuthenticationError(w http.ResponseWriter) {
	fmt.Println(">>> Cannot authenticate user")
	w.WriteHeader(401)
	m := utils.Message{Status: 401, Message: "Authorization Error"}
	j, err := json.Marshal(m)
	if err != nil {
		fmt.Println("ERR ", err)
	} else {
		w.Write(j)
	}
}
