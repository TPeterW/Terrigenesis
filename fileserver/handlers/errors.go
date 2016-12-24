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
	w.Header().Set("Content-Type", "application/json")
	m := utils.Response{Status: 401, Message: "Authorization error"}
	json.NewEncoder(w).Encode(m)
}

/*
SessionNotFoundError Handles all unknown session errors
*/
func SessionNotFoundError(w http.ResponseWriter) {
	fmt.Println(">>> Cannot find session")
	w.WriteHeader(404)
	w.Header().Set("Content-Type", "application/json")
	m := utils.Response{Status: 404, Message: "Session not found"}
	json.NewEncoder(w).Encode(m)
}

/*
IllegalArgumentsError Handles all bad requests
*/
func IllegalArgumentsError(w http.ResponseWriter) {
	fmt.Println(">>> Illegal arguments")
	w.WriteHeader(400)
	w.Header().Set("Content-Type", "application/json")
	m := utils.Response{Status: 400, Message: "Bad request, possibily missing arguments"}
	json.NewEncoder(w).Encode(m)
}

/*
UnknownCommandError Handles all unknown commands
*/
func UnknownCommandError(w http.ResponseWriter) {
	fmt.Println(">>> Unknown commands")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(406)
	m := utils.Response{Status: 406, Message: "Unknown command"}
	json.NewEncoder(w).Encode(m)
}

/*
FileNotFoundError Hanldes all file not found error
*/
func FileNotFoundError(w http.ResponseWriter) {
	fmt.Println(">>> File not found")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)
	m := utils.Response{Status: 500, Message: "File not found"}
	json.NewEncoder(w).Encode(m)
}

/*
FileTypeError Handles all file type error
*/
func FileTypeError(w http.ResponseWriter) {
	fmt.Println(">>> File type error")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)
	m := utils.Response{Status: 500, Message: "File type error, mistaking file for dir or V/V"}
	json.NewEncoder(w).Encode(m)
}

/*
FolderPermissionError Handles cases where user tries to access outside "./db"
*/
func FolderPermissionError(w http.ResponseWriter) {
	fmt.Println(">>> Folder permission error")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(401)
	m := utils.Response{Status: 401, Message: "Cannot access outside database folder"}
	json.NewEncoder(w).Encode(m)
}

/*
FileExistError Hanles cases where file or directory already exists
*/
func FileExistError(w http.ResponseWriter) {
	fmt.Println(">>> Already exists")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(409)
	m := utils.Response{Status: 401, Message: "File or directory already exists"}
	json.NewEncoder(w).Encode(m)
}

/*
GeneralError Handles general type error
*/
func GeneralError(w http.ResponseWriter, statusCode int, err error) {
	fmt.Printf("General error, %s\n", err.Error())
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	m := utils.Response{Status: statusCode, Message: err.Error()}
	json.NewEncoder(w).Encode(m)
}
