package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"strings"
	"terrigenesis/fileserver/utils"
	"terrigenesis/secrets"
	"time"
)

/*
makeGetRequest Creates and sends GET requests, returns response
*/
func makeGetRequest(timeout time.Duration, url string, query url.Values, del delegate) (*http.Response, utils.Response, bool) {
	client := http.Client{Timeout: timeout}
	if !strings.HasPrefix(url, "/") {
		url = "/" + url
	}
	req, err := http.NewRequest("GET", secrets.URL()+url, nil)
	if err != nil {
		fmt.Println("Cannot form request")
		return nil, utils.Response{}, false
	}

	req.SetBasicAuth(del.username, del.password)
	req.URL.RawQuery = query.Encode()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request")
		return resp, utils.Response{}, false
	}

	if strings.Compare(url, "/downfile") == 0 {
		return resp, utils.Response{}, true
	}

	var target utils.Response
	json.NewDecoder(resp.Body).Decode(&target)
	if resp.StatusCode != http.StatusOK {
		fmt.Println(target.Message)
		return resp, utils.Response{}, false
	}
	return resp, target, true
}

/*
makePostRequest Creates and sends POST requests, returns response
*/
func makePostRequest(timeout time.Duration, url string, formFields map[string]string, filename string, del delegate) (utils.Response, bool) {
	if !strings.HasPrefix(url, "/") {
		url = "/" + url
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// form field data
	for key, val := range formFields {
		writer.WriteField(key, val)
	}

	// form file data
	if strings.Compare(filename, "") != 0 {
		// substitute "~" with actual home directory
		usr, err := user.Current()
		if err != nil {
			fmt.Println("SysErr: " + err.Error())
		} else {
			if strings.HasPrefix(filename, "~") {
				filename = strings.Replace(filename, "~", usr.HomeDir, -1)
			}
		}

		fileWriter, err := writer.CreateFormFile("file", strings.Split(filename, "/")[len(strings.Split(filename, "/"))-1])
		if err != nil {
			fmt.Print("Error writing to buffer")
			return utils.Response{}, false
		}

		f, err := os.Open(filename)
		if err != nil {
			fmt.Print(err.Error())
			return utils.Response{}, false
		}
		_, err = io.Copy(fileWriter, f)
		if err != nil {
			fmt.Print("SysErr: " + err.Error())
			return utils.Response{}, false
		}
	}

	req, err := http.NewRequest("POST", secrets.URL()+url, body)
	if err != nil {
		fmt.Println("SysErr: " + err.Error())
		return utils.Response{}, false
	}

	writer.Close()

	req.SetBasicAuth(del.username, del.password)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.ContentLength = int64(body.Len())

	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("SysErr: " + err.Error())
		return utils.Response{}, false
	}

	var target utils.Response
	json.NewDecoder(resp.Body).Decode(&target)
	if resp.StatusCode != http.StatusOK {
		fmt.Println(target.Message)
		return utils.Response{}, false
	}
	return target, true
}
