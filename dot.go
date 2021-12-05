package fileserve

import (
	"net/http"
	"strings"
)

func NoDot(err string, code int) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			pathParts := strings.Split(r.URL.Path, "/")
			for _, part := range pathParts {
				if strings.HasPrefix(part, ".") {
					http.Error(w, err, code)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
