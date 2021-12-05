package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCreds(t *testing.T) {
	t.Run("user:pass", func(t *testing.T) {
		creds, err := parseAuths([]string{"user:pass"})
		assert.NoError(t, err)
		assert.Equal(t, "pass", creds["user"])
	})

	t.Run("invalid format", func(t *testing.T) {
		_, err := parseAuths([]string{"userpass"})
		assert.Error(t, err)
		assert.Equal(t, "auth must be specified in the format username:password", err.Error())
	})

	t.Run("", func(t *testing.T) {
		_, err := parseAuths([]string{"us:er:pa:ss"})
		assert.Error(t, err)
		assert.Equal(t, "only one colon is allowed", err.Error())
	})
}
