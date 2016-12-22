package handlers

// type FileInfo interface {
//         Name() string       // base name of the file
//         Size() int64        // length in bytes for regular files; system-dependent for others
//         Mode() FileMode     // file mode bits
//         ModTime() time.Time // modification time
//         IsDir() bool        // abbreviation for Mode().IsDir()
//         Sys() interface{}   // underlying data source (can return nil)
// }

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"terrigenesis/fileserver/utils"
)

/*
PrintWorkingDirectory Returns current working directory
*/
func PrintWorkingDirectory(w http.ResponseWriter, r *http.Request, session utils.Session) {
	// compose header
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	// generate response body
	m := utils.Response{Status: 200, CWD: session.CWD}

	// convert response to json
	json.NewEncoder(w).Encode(m)
}

/*
ListFiles Returns a list containing all files in current directory of the session
*/
func ListFiles(w http.ResponseWriter, r *http.Request, session utils.Session) {
	curDir := session.CWD

	// compose header
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	m := utils.Response{Status: 200, CWD: session.CWD}

	files, err := ioutil.ReadDir(curDir)
	if err != nil {
		log.Fatal(err)
		m.Status = 500
		m.Message = "Error reading directory"
	} else {
		for _, file := range files {
			if file.IsDir() {
				m.DirFiles = append(m.DirFiles, file.Name()+"/")
			} else {
				m.DirFiles = append(m.DirFiles, file.Name())
			}
		}
	}

	json.NewEncoder(w).Encode(m)
}
