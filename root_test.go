package fileserve

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRoot(t *testing.T) {
	t.Run("dot", func(t *testing.T) {
		r, err := newRoot(".")
		assert.NoError(t, err)
		assert.Equal(t, http.Dir("."), r)
	})

	t.Run("not a directory", func(t *testing.T) {
		_, err := newRoot("./LICENSE")
		assert.Error(t, err)
		assert.Equal(t, "cannot serve ./LICENSE: not a directory", err.Error())
	})

	t.Run("no such file or directory", func(t *testing.T) {
		_, err := newRoot("./XYZ")
		assert.Error(t, err)
		assert.Equal(t, "cannot serve ./XYZ: stat ./XYZ: no such file or directory", err.Error())
	})
}
