package fileserve

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strings"
)

const DefaultIndexPage = "index.html"

type DirRootOptions struct {
	ShowDirListing bool
	ShowDotFiles   bool
}

func NewDirRoot(dir string, optFns ...func(o *DirRootOptions)) (http.FileSystem, error) {
	fileInfo, err := os.Stat(dir)
	if err != nil {
		return nil, fmt.Errorf("cannot serve %v: %v", dir, err)
	} else if !fileInfo.IsDir() {
		return nil, fmt.Errorf("cannot serve %v: not a directory", dir)
	} else {
		file, err := os.Open(dir)
		if err != nil {
			return nil, fmt.Errorf("cannot serve %v: %v", dir, err)
		}
		file.Close()
	}

	options := DirRootOptions{
		ShowDirListing: false,
		ShowDotFiles:   false,
	}

	for _, fn := range optFns {
		fn(&options)
	}

	return dirFilesystem{
		FileSystem:     http.Dir(dir),
		showDirListing: options.ShowDirListing,
		showDotFiles:   options.ShowDotFiles,
	}, nil
}

type dirFilesystem struct {
	http.FileSystem
	showDirListing bool
	showDotFiles   bool
}

func (fs dirFilesystem) Open(name string) (http.File, error) {
	f, err := fs.FileSystem.Open(name)
	if err != nil {
		return nil, err
	}

	if fs.showDirListing && fs.showDotFiles {
		return f, nil
	}

	return wrappedFile{
		File:           f,
		showDirListing: fs.showDirListing,
		showDotFiles:   fs.showDotFiles,
	}, nil
}

type wrappedFile struct {
	http.File
	showDirListing bool
	showDotFiles   bool
}

func (f wrappedFile) Stat() (os.FileInfo, error) {
	s, err := f.File.Stat()
	if err != nil {
		return nil, err
	}

	if s.IsDir() && !f.showDirListing {
		return f.findIndexPage(s, "index.html")
	}

	return s, err
}

func (f wrappedFile) Readdir(count int) ([]fs.FileInfo, error) {
	fi, err := f.File.Readdir(count)
	if err != nil {
		return nil, err
	}

	if f.showDotFiles {
		return fi, nil
	}

	filtered := []fs.FileInfo{}

	for _, f := range fi {
		if !strings.HasPrefix(f.Name(), ".") {
			filtered = append(filtered, f)
		}
	}

	return filtered, err
}

func (f wrappedFile) findIndexPage(s fs.FileInfo, indexPageName string) (os.FileInfo, error) {
LOOP:
	for {
		fl, err := f.File.Readdir(1)
		switch err {
		case io.EOF:
			break LOOP
		case nil:
			for _, f := range fl {
				if f.Name() == indexPageName {
					return s, err
				}
			}
		default:
			return nil, err
		}
	}

	return nil, os.ErrNotExist
}
