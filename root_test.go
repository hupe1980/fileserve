package fileserve

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestNewDirRoot(t *testing.T) {
	t.Run("dot", func(t *testing.T) {
		r, err := NewDirRoot(".")
		assert.NoError(t, err)
		assert.Equal(t, rootFilesystem{
			FileSystem: http.Dir("."),
			fileOptions: fileOptions{
				hideDirListing: false,
				hideDotFiles:   false,
			},
		}, r)
	})

	t.Run("not a directory", func(t *testing.T) {
		_, err := NewDirRoot("./LICENSE")
		assert.Error(t, err)
		assert.Equal(t, "cannot serve ./LICENSE: not a directory", err.Error())
	})

	t.Run("no such file or directory", func(t *testing.T) {
		_, err := NewDirRoot("./XYZ")
		assert.Error(t, err)
		assert.Equal(t, "cannot serve ./XYZ: stat ./XYZ: no such file or directory", err.Error())
	})
}

func TestRootFilesystem(t *testing.T) {
	appFS := afero.NewMemMapFs()

	// create test files and directories
	err := afero.WriteFile(appFS, "index.html", []byte("indexPage"), 0644)
	assert.NoError(t, err)

	err = afero.WriteFile(appFS, ".dot", []byte("dot"), 0644)
	assert.NoError(t, err)

	err = appFS.Mkdir("noindex", 0755)
	assert.NoError(t, err)
	err = afero.WriteFile(appFS, "noindex/a.txt", []byte("a"), 0644)
	assert.NoError(t, err)
	err = afero.WriteFile(appFS, "noindex/b.txt", []byte("b"), 0644)
	assert.NoError(t, err)

	err = appFS.Mkdir("empty", 0755)
	assert.NoError(t, err)

	t.Run("read indexPage", func(t *testing.T) {
		rootFs := rootFilesystem{
			FileSystem: afero.NewHttpFs(appFS),
			fileOptions: fileOptions{
				hideDirListing: false,
				hideDotFiles:   false,
			},
		}

		indexPage, err := rootFs.Open("index.html")
		assert.NoError(t, err)
		defer indexPage.Close()

		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(indexPage)
		assert.NoError(t, err)
		assert.Equal(t, []byte("indexPage"), buf.Bytes())
	})

	t.Run("read dir", func(t *testing.T) {
		rootFs := rootFilesystem{
			FileSystem: afero.NewHttpFs(appFS),
			fileOptions: fileOptions{
				hideDirListing: false,
				hideDotFiles:   false,
			},
		}

		root, err := rootFs.Open(".")
		assert.NoError(t, err)
		defer root.Close()

		fi, err := root.Readdir(-1)
		assert.NoError(t, err)
		assert.Equal(t, 4, len(fi))
	})

	t.Run("read dir - hide dot", func(t *testing.T) {
		rootFs := rootFilesystem{
			FileSystem: afero.NewHttpFs(appFS),
			fileOptions: fileOptions{
				hideDirListing: false,
				hideDotFiles:   true,
			},
		}

		root, err := rootFs.Open(".")
		assert.NoError(t, err)
		defer root.Close()

		fi, err := root.Readdir(-1)
		assert.NoError(t, err)
		assert.Equal(t, 3, len(fi))
	})

	t.Run("read noindex dir - hide dir listing", func(t *testing.T) {
		rootFs := rootFilesystem{
			FileSystem: afero.NewHttpFs(appFS),
			fileOptions: fileOptions{
				hideDirListing: true,
				hideDotFiles:   false,
			},
		}

		root, err := rootFs.Open("noindex")
		assert.NoError(t, err)
		defer root.Close()

		_, err = root.Stat()
		assert.Error(t, err)
		assert.Equal(t, "file does not exist", err.Error())
	})
}
