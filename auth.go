package main

import (
	"crypto/subtle"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const credsPartCount = 2

func ParseCreds(creds string) (string, string, error) {
	if !strings.Contains(creds, ":") {
		return "", "", errors.New("credentials must be specified in the format username:password")
	}

	parts := strings.Split(creds, ":")
	if len(parts) != credsPartCount {
		return "", "", errors.New("cannot parse credentials")
	}

	return parts[0], parts[1], nil
}

// BasicAuth adds a basic authentication middleware
func BasicAuth(next http.Handler, realm string, creds map[string]string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok {
			basicAuthFailed(w, realm)
			return
		}
		credPass, credUserOk := creds[user]
		if !credUserOk || subtle.ConstantTimeCompare([]byte(pass), []byte(credPass)) != 1 {
			basicAuthFailed(w, realm)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func basicAuthFailed(w http.ResponseWriter, realm string) {
	w.Header().Add("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, realm))
	w.WriteHeader(http.StatusUnauthorized)
}
