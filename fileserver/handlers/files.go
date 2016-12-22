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
	"fmt"
	"net/http"
	"strings"
	"terrigenesis/fileserver/utils"
)

/*
DownloadFile Handles requests to download a file
*/
func DownloadFile(w http.ResponseWriter, r *http.Request, session utils.Session) bool {
	fmt.Println(">>> Download file")
	var fileName []string
	if fileName = r.URL.Query()["filename"]; fileName == nil {
		IllegalArgumentsError(w)
		return false
	}

	// TODO:

	http.ServeFile(w, r, session.CWD+"/"+strings.Join(fileName, ""))
	return true
}

/*
UploadFile Handles requests to upload file
*/
func UploadFile(w http.ResponseWriter, r *http.Request, session utils.Session) bool {
	// TODO:

	return true
}
