package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRoot(t *testing.T) {
	t.Run("dot", func(t *testing.T) {
		r, err := Root(".")
		assert.NoError(t, err)
		assert.Equal(t, http.Dir("."), r)
	})

	t.Run("not a directory", func(t *testing.T) {
		_, err := Root("./LICENSE")
		assert.Error(t, err)
		assert.Equal(t, "cannot serve ./LICENSE: not a directory", err.Error())
	})

	t.Run("no such file or directory", func(t *testing.T) {
		_, err := Root("./XYZ")
		assert.Error(t, err)
		assert.Equal(t, "cannot serve ./XYZ: stat ./XYZ: no such file or directory", err.Error())
	})
}
