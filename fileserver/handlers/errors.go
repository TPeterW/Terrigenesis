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
	m := utils.Message{Status: 401, Message: "Authorization error"}
	json.NewEncoder(w).Encode(m)
}

/*
SessionNotFoundError Handles all unknown session errors
*/
func SessionNotFoundError(w http.ResponseWriter) {
	fmt.Println(">>> Cannot find session")
	w.WriteHeader(404)
	m := utils.Message{Status: 404, Message: "Session not found"}
	json.NewEncoder(w).Encode(m)
}

/*
IllegalArgumentsError Handles all bad requests
*/
func IllegalArgumentsError(w http.ResponseWriter) {
	fmt.Println(">>> Illegal Arguments")
	w.WriteHeader(400)
	m := utils.Message{Status: 400, Message: "Bad request, possibily missing arguments"}
	json.NewEncoder(w).Encode(m)
}
