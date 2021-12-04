package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCreds(t *testing.T) {
	t.Run("user:pass", func(t *testing.T) {
		user, pass, err := ParseCreds("user:pass")
		assert.NoError(t, err)
		assert.Equal(t, "user", user)
		assert.Equal(t, "pass", pass)
	})

	t.Run("invalid format", func(t *testing.T) {
		_, _, err := ParseCreds("userpass")
		assert.Error(t, err)
		assert.Equal(t, "credentials must be specified in the format username:password", err.Error())
	})

	t.Run("too many parts", func(t *testing.T) {
		_, _, err := ParseCreds("us:er:pa:ss")
		assert.Error(t, err)
		assert.Equal(t, "cannot parse credentials", err.Error())
	})
}
