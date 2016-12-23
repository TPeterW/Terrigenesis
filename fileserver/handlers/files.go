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
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"terrigenesis/fileserver/utils"
)

/*
DownloadFile Handles requests to download a file
*/
func DownloadFile(w http.ResponseWriter, r *http.Request, session utils.Session) {
	var fileName []string
	if fileName = r.URL.Query()["filename"]; fileName == nil {
		IllegalArgumentsError(w)
	}

	pathToFile := session.CWD + "/" + strings.Join(fileName, "")
	var entry os.FileInfo
	var err error
	if entry, err = os.Stat(pathToFile); err == nil {
		// entry exists
		if !entry.IsDir() {
			// is not directory
			http.ServeFile(w, r, pathToFile)
		}
		FileTypeError(w)
	}

	FileNotFoundError(w)
}

/*
UploadFile Handles requests to upload file
*/
func UploadFile(w http.ResponseWriter, r *http.Request, session utils.Session) {
	file, handler, err := r.FormFile("file")
	defer file.Close()
	if err != nil {
		GeneralError(w, 500, "Error parsing file")
	} else {
		f, err := os.OpenFile(session.CWD+"/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		defer f.Close()
		if err != nil {
			GeneralError(w, 500, "Error writing file")
		} else {
			io.Copy(f, file)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			m := utils.Response{Status: 200, Message: "Successfully uploaded file: " + handler.Filename}
			json.NewEncoder(w).Encode(m)
		}
	}
}

/*
RemoveFile Handles requests to remove file
*/
func RemoveFile(w http.ResponseWriter, r *http.Request, form url.Values, session utils.Session) {
	if form["filename"] == nil {
		IllegalArgumentsError(w)
		return
	}

	pathToFile := session.CWD + "/" + strings.Join(form["filename"], "")
	var entry os.FileInfo
	var err error
	if entry, err = os.Stat(pathToFile); err == nil {
		// entry exists
		if !entry.IsDir() {
			// is not directory
			if removeErr := os.Remove(pathToFile); removeErr != nil {
				GeneralError(w, 500, "Error removing file")
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(200)
				m := utils.Response{Status: 200, Message: "Successfully removed file"}
				json.NewEncoder(w).Encode(m)
			}
		} else {
			// is directory
			GeneralError(w, 500, "File is a directory")
		}
	} else {
		FileNotFoundError(w) // other options not very possible
	}
}

/*
MoveFileDir Moves a specific file to a specific location
*/
func MoveFileDir(w http.ResponseWriter, r *http.Request, form url.Values, session utils.Session) {
	if form["origin"] == nil || form["destination"] == nil {
		IllegalArgumentsError(w)
		return
	}

	pointerSplited := strings.Split(session.CWD, "/")
	origin := session.CWD + "/" + strings.Join(form["origin"], "")
	destSplited := strings.Split(strings.Join(form["destination"], ""), "/")

	for _, dest := range destSplited {
		if strings.Compare(dest, ".") == 0 {
			continue
		} else if strings.Compare(dest, "..") == 0 {
			// up one level
			curDir := strings.Join(pointerSplited, "")
			if strings.Compare(curDir, "./db") != 0 && strings.Compare(curDir, "./db/") != 0 {
				// when not top level
				pointerSplited = pointerSplited[:len(pointerSplited)-1]
			} else {
				FolderPermissionError(w)
				return
			}
		} else {
			// down one level
			path := strings.Join(pointerSplited, "/") + "/" + dest
			if _, err := os.Stat(path); err == nil {
				// entry exists
				pointerSplited = append(pointerSplited, dest)
			} else {
				// entry doesn't exist
				FileNotFoundError(w)
				return
			}
		}
	}
	destination := strings.Join(pointerSplited, "/") + "/" + strings.Join(form["origin"], "")
	if err := os.Rename(origin, destination); err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		m := utils.Response{Status: 200, Message: "Successfully moved file/directory"}
		json.NewEncoder(w).Encode(m)
	} else {
		GeneralError(w, 500, "Unable to move file/directory")
	}
}
