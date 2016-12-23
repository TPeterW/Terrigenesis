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
	"net/url"
	"os"
	"strings"
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

/*
ChangeDir Changes to a certain directory (direct parent or child)
*/
func ChangeDir(w http.ResponseWriter, r *http.Request, form url.Values, session utils.Session) utils.Session {
	if form["dirname"] == nil {
		IllegalArgumentsError(w)
		return session
	}

	curSplited := strings.Split(session.CWD, "/")
	destSplited := strings.Split(strings.Join(form["dirname"], ""), "/")

	for _, dest := range destSplited {
		if strings.Compare(dest, ".") == 0 {
			continue
		} else if strings.Compare(dest, "..") == 0 {
			// up one level
			curDir := strings.Join(curSplited, "/")
			if strings.Compare(curDir, "./db") != 0 && strings.Compare(curDir, "./db/") != 0 {
				// when not top level
				curSplited = curSplited[:len(curSplited)-1]
			} else {
				FolderPermissionError(w)
				return session
			}
		} else {
			// down one level
			pathToNewDir := strings.Join(curSplited, "/") + "/" + dest
			if entry, err := os.Stat(pathToNewDir); err == nil {
				// dir exists
				if !entry.IsDir() {
					// not a directory
					FileTypeError(w)
					return session
				}
				// is actually a directory
				curSplited = append(curSplited, dest)
			} else {
				// dir doesn't exist
				FileNotFoundError(w)
				return session
			}
		}
	}

	session.CWD = strings.TrimRight(strings.Join(curSplited, "/"), "/")
	w.WriteHeader(200)
	m := utils.Response{Status: 200, CWD: session.CWD}
	json.NewEncoder(w).Encode(m)

	return session
}

/*
MakeDir Creates a specific directory
*/
func MakeDir(w http.ResponseWriter, r *http.Request, form url.Values, session utils.Session) {
	if form["dirname"] == nil {
		IllegalArgumentsError(w)
		return
	}

	pathToDir := session.CWD + "/" + strings.Join(form["dirname"], "")
	if entry, err := os.Stat(pathToDir); err == nil {
		if entry.IsDir() {
			FileExistError(w)
			return
		}
	}
	os.Mkdir(pathToDir, os.ModeDir)
	m := utils.Response{Status: 200, Message: "Successfully created directory: " + strings.Join(form["dirname"], "")}
	json.NewEncoder(w).Encode(m)
}

/*
RemoveDir Removes a specific directory
*/
func RemoveDir(w http.ResponseWriter, r *http.Request, form url.Values, session utils.Session) {
	if form["dirname"] == nil {
		IllegalArgumentsError(w)
		return
	}

	pathToDir := session.CWD + "/" + strings.Join(form["dirname"], "")
	var entry os.FileInfo
	var err error
	if entry, err = os.Stat(pathToDir); err == nil {
		// entry exists
		if entry.IsDir() {
			// is actually a directory
			if removeErr := os.Remove(pathToDir); removeErr != nil {
				GeneralError(w, 500, "Error removing directory")
			} else {
				w.WriteHeader(200)
				m := utils.Response{Status: 200, Message: "Successfully removed directory"}
				json.NewEncoder(w).Encode(m)
			}
		} else {
			// is not directory
			GeneralError(w, 500, "Directory is a file")
		}
	} else {
		FileNotFoundError(w) // other options not very possible
	}
}
