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
	m := utils.Response{Status: 401, Message: "Authorization error"}
	json.NewEncoder(w).Encode(m)
}

/*
SessionNotFoundError Handles all unknown session errors
*/
func SessionNotFoundError(w http.ResponseWriter) {
	fmt.Println(">>> Cannot find session")
	w.WriteHeader(404)
	m := utils.Response{Status: 404, Message: "Session not found"}
	json.NewEncoder(w).Encode(m)
}

/*
IllegalArgumentsError Handles all bad requests
*/
func IllegalArgumentsError(w http.ResponseWriter) {
	fmt.Println(">>> Illegal arguments")
	w.WriteHeader(400)
	m := utils.Response{Status: 400, Message: "Bad request, possibily missing arguments"}
	json.NewEncoder(w).Encode(m)
}

/*
UnknownCommandError Handles all unknown commands
*/
func UnknownCommandError(w http.ResponseWriter) {
	fmt.Println(">>> Unknown commands")
	w.WriteHeader(406)
	m := utils.Response{Status: 406, Message: "Unknown command"}
	json.NewEncoder(w).Encode(m)
}

/*
FileNotFoundError Hanldes all file not found error
*/
func FileNotFoundError(w http.ResponseWriter) {
	fmt.Println(">>> File not found")
	w.WriteHeader(500)
	m := utils.Response{Status: 500, Message: "File not found"}
	json.NewEncoder(w).Encode(m)
}

/*
FileTypeError Handles all file type error
*/
func FileTypeError(w http.ResponseWriter) {
	fmt.Println(">>> File type error")
	w.WriteHeader(500)
	m := utils.Response{Status: 500, Message: "File type error, mistaking file for dir or V/V"}
	json.NewEncoder(w).Encode(m)
}

/*
FolderPermissionError Handles cases where user tries to access outside "./db"
*/
func FolderPermissionError(w http.ResponseWriter) {
	fmt.Println(">>> Folder permission error")
	w.WriteHeader(401)
	m := utils.Response{Status: 401, Message: "Cannot access outside database folder"}
	json.NewEncoder(w).Encode(m)
}

/*
GeneralError Handles general type error
*/
func GeneralError(w http.ResponseWriter, statusCode int, message string) {
	fmt.Printf("General error, %s\n", message)
	w.WriteHeader(statusCode)
	m := utils.Response{Status: statusCode, Message: message}
	json.NewEncoder(w).Encode(m)
}
