package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os/user"
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
	query := make(url.Values)
	query.Add("Token", del.token)

	_, _, ok = makeGetRequest(10*time.Second, "/closecon", query, del)
	return ok
}

/*
printWorkingDirectory Print working directory
*/
func printWorkingDirectory(del delegate) {
	var target utils.Response
	var ok bool
	query := make(url.Values)
	query.Add("Token", del.token)

	_, target, ok = makeGetRequest(10*time.Second, "pwd", query, del)
	if ok {
		fmt.Println(target.CWD)
	}
}

/*
listFiles List files
*/
func listFiles(del delegate) {
	var target utils.Response
	var ok bool
	query := make(url.Values)
	query.Add("Token", del.token)

	_, target, ok = makeGetRequest(10*time.Second, "dir", query, del)
	if ok {
		if len(target.DirFiles) > 0 {
			fmt.Printf(strings.Join(target.DirFiles, "\t") + "\n")
		}
	}
}

/*
downloadFile Download file
*/
func downloadFile(del delegate, filename string) {
	// substitute "~" with actual home directory
	usr, err := user.Current()
	if err != nil {
		fmt.Println("SysErr: " + err.Error())
		return
	}

	var resp *http.Response
	var target utils.Response
	var ok bool
	query := make(url.Values)
	query.Add("Token", del.token)
	query.Add("filename", filename)
	resp, target, ok = makeGetRequest(30*time.Second, "downfile", query, del)

	if ok {
		fileBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("SysErr: " + err.Error())
		} else {
			ioutil.WriteFile(usr.HomeDir+"/Downloads/"+strings.Split(filename, "/")[len(strings.Split(filename, "/"))-1], fileBytes, 0666)
		}

		// out, err := os.OpenFile(usr.HomeDir+"/Downloads/"+strings.Split(filename, "/")[len(strings.Split(filename, "/"))-1], os.O_RDWR|os.O_CREATE, 0666)
		// defer out.Close()

		// if err != nil {
		// 	fmt.Println("SysErr: " + err.Error())
		// 	os.Remove(usr.HomeDir + "/Downloads/" + strings.Split(filename, "/")[len(strings.Split(filename, "/"))-1])
		// } else {
		// 	_, err = io.Copy(out, body)
		// 	if err != nil {
		// 		fmt.Println("SysErr: " + err.Error())
		// 		os.Remove(usr.HomeDir + "/Downloads/" + strings.Split(filename, "/")[len(strings.Split(filename, "/"))-1])
		// 	}
		// }
	} else {
		fmt.Print(target.Message)
	}
}

/*
uploadFile Upload file
*/
func uploadFile(del delegate, filename string) {
	form := make(map[string]string)
	form["token"] = del.token

	response, ok := makePostRequest(60*time.Second, "upfile", form, filename, del)
	if !ok {
		fmt.Println(response.Message)
	}
}

/*
changeDir Change directory
*/
func changeDir(del delegate, dirname string) {
	form := make(map[string]string)
	form["token"] = del.token
	form["dirname"] = dirname

	makePostRequest(10*time.Second, "chdir", form, "", del)
}

/*
makeDir Make directory
*/
func makeDir(del delegate, dirname string) {
	form := make(map[string]string)
	form["token"] = del.token
	form["dirname"] = dirname

	makePostRequest(10*time.Second, "mkdir", form, "", del)
}

/*
removeDir Remove directory
*/
func removeDir(del delegate, dirname string) {
	form := make(map[string]string)
	form["token"] = del.token
	form["dirname"] = dirname

	response, ok := makePostRequest(10*time.Second, "rmdir", form, "", del)
	if !ok {
		fmt.Println(response.Message)
	}
}

func removeFile(del delegate, args []string) {
	form := make(map[string]string)
	form["token"] = del.token

	for _, filename := range strings.Split(strings.Join(args, " "), ",") {
		filename = strings.TrimPrefix(filename, " ")
		filename = strings.TrimSuffix(filename, " ")
		form["filename"] = filename
		makePostRequest(10*time.Second, "rmfile", form, "", del)
	}
}
