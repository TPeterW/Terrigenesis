package handlers

import "net/http"
import "terrigenesis/fileserver/utils"
import "encoding/json"

/*
PrintWorkingDirectory Returns current working directory
*/
func PrintWorkingDirectory(w http.ResponseWriter, r *http.Request, session utils.Session) {
	w.Header().Set("Content-Type", "application/json")
	// compose header
	w.WriteHeader(200)

	// generate response body
	m := utils.Response{Status: 200, CWD: session.CWD}

	// convert response to json
	json.NewEncoder(w).Encode(m)
}
