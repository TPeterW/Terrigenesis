package client

import (
	"fmt"
	"net/url"
	"strings"
	"terrigenesis/fileserver/utils"
	"time"
)

type delegate struct {
	token    string
	username string
	password string
}

/*
openConnection Open connection
*/
func openConnection(del delegate) (delegate, bool) {
	var target utils.Response
	var ok bool
	target, ok = makePostRequest(10*time.Second, "/estabcon", nil, "", del)
	if ok {
		fmt.Printf("Succesfully connected to server\n>>> ")
		del.token = target.Token
		return del, true
	}

	return del, false
}

/*
closeConnection Close connection
*/
func closeConnection(del delegate) bool {
	var ok bool
	query := url.Values{}
	query.Add("Token", del.token)

	_, ok = makeGetRequest(10*time.Second, "/closecon", query, del)
	return ok
}

/*
printWorkingDirectory Print working directory
*/
func printWorkingDirectory(del delegate) {
	var target utils.Response
	var ok bool
	query := url.Values{}
	query.Add("Token", del.token)

	target, ok = makeGetRequest(10*time.Second, "pwd", query, del)
	if ok {
		fmt.Printf(target.CWD + "\n>>> ")
	}
}

/*
listFiles List files
*/
func listFiles(del delegate) {
	var target utils.Response
	var ok bool
	query := url.Values{}
	query.Add("Token", del.token)

	target, ok = makeGetRequest(10*time.Second, "dir", query, del)
	if ok {
		if len(target.DirFiles) > 0 {
			fmt.Printf(strings.Join(target.DirFiles, "\t") + "\n>>> ")
		} else {
			fmt.Printf(">>> ")
		}
	}
}

/*
changeDir Change directory
*/
func changeDir(del delegate, dirname string) {
	// bodyBuf := &bytes.Buffer{}
	// bodyWriter := multipart.NewWriter(bodyBuf)
	// bodyWriter.WriteField("token", del.token)
	// bodyWriter.WriteField("dirname", dirname)

	form := make(map[string]string)
	form["token"] = del.token
	form["dirname"] = dirname

	makePostRequest(10*time.Second, "chdir", form, "", del)
	fmt.Printf(">>> ")
}