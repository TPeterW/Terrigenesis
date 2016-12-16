package utils

import (
	"encoding/base64"
	"net/http"
	"strings"
	"terrigenesis/secrets"
)

/*
BasicAuth Middleware for handling http basic auth
*/
func BasicAuth(r *http.Request) bool {
	s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(s) != 2 {
		return false
	}

	b, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil {
		return false
	}

	pair := strings.SplitN(string(b), ":", 2)
	if len(pair) != 2 {
		return false
	}

	return pair[0] == secrets.Username() && pair[1] == secrets.Password()
}
