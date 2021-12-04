package fileserve

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicAuth(t *testing.T) {
	withBasicAuth := BasicAuth("restricted", map[string]string{"user": "pass"})

	handler := withBasicAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "helloworld")
	}))

	ts := httptest.NewServer(handler)

	t.Run("unauthorized", func(t *testing.T) {
		r, err := http.Get(ts.URL)
		assert.NoError(t, err)
		defer r.Body.Close()
		assert.Equal(t, http.StatusUnauthorized, r.StatusCode)
	})

	t.Run("Authorized", func(t *testing.T) {
		req, err := http.NewRequest("Get", ts.URL, nil)
		assert.NoError(t, err)

		req.SetBasicAuth("user", "pass")

		r, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer r.Body.Close()
		assert.Equal(t, http.StatusOK, r.StatusCode)
	})
}
