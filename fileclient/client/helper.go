package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"terrigenesis/fileserver/utils"
	"terrigenesis/secrets"
	"time"
)

/*
makeGetRequest Creates and sends GET requests, returns response
*/
func makeGetRequest(timeout time.Duration, url string, query url.Values, del delegate) (utils.Response, bool) {
	client := http.Client{Timeout: timeout}
	if !strings.HasPrefix(url, "/") {
		url = "/" + url
	}
	req, err := http.NewRequest("GET", secrets.URL()+url, nil)
	if err != nil {
		fmt.Println("Cannot form request")
		return utils.Response{}, false
	}

	req.SetBasicAuth(del.username, del.password)
	req.URL.RawQuery = query.Encode()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request")
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

/*
makePostRequest Creates and sends POST requests, returns response
*/
func makePostRequest(timeout time.Duration, url string, formFields map[string]string, filename string, del delegate) (utils.Response, bool) {
	if !strings.HasPrefix(url, "/") {
		url = "/" + url
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for key, val := range formFields {
		writer.WriteField(key, val)
	}

	req, err := http.NewRequest("POST", secrets.URL()+url, body)
	if err != nil {
		fmt.Println("Cannot form request")
		return utils.Response{}, false
	}

	// TODO: upload file
	// if strings.Compare(filename, "") != 0 {

	// }

	writer.Close()

	req.SetBasicAuth(del.username, del.password)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.ContentLength = int64(body.Len())

	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error sending request")
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
