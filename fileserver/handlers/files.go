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
	// TODO:

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
