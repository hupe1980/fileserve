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

type Root interface {
	http.FileSystem
}

type RootOptions struct {
	HideDirListing bool
	HideDotFiles   bool
}

func NewDirRoot(dir string, optFns ...func(o *RootOptions)) (Root, error) {
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

	options := RootOptions{
		HideDirListing: false,
		HideDotFiles:   false,
	}

	for _, fn := range optFns {
		fn(&options)
	}

	return rootFilesystem{
		FileSystem: http.Dir(dir),
		fileOptions: fileOptions{
			hideDirListing: options.HideDirListing,
			hideDotFiles:   options.HideDotFiles,
		},
	}, nil
}

func NewFSRoot(distFS fs.FS, optFns ...func(o *RootOptions)) (Root, error) {
	options := RootOptions{
		HideDirListing: false,
		HideDotFiles:   false,
	}

	for _, fn := range optFns {
		fn(&options)
	}

	return rootFilesystem{
		FileSystem: http.FS(distFS),
		fileOptions: fileOptions{
			hideDirListing: options.HideDirListing,
			hideDotFiles:   options.HideDotFiles,
		},
	}, nil
}

type rootFilesystem struct {
	http.FileSystem
	fileOptions
}

func (fs rootFilesystem) Open(name string) (http.File, error) {
	f, err := fs.FileSystem.Open(name)
	if err != nil {
		return nil, err
	}

	if !fs.hideDirListing && !fs.hideDotFiles {
		return f, nil
	}

	return wrappedFile{
		File:        f,
		fileOptions: fs.fileOptions,
	}, nil
}

type fileOptions struct {
	hideDirListing bool
	hideDotFiles   bool
}

type wrappedFile struct {
	http.File
	fileOptions
}

func (f wrappedFile) Stat() (os.FileInfo, error) {
	s, err := f.File.Stat()
	if err != nil {
		return nil, err
	}

	if s.IsDir() && f.hideDirListing {
		return f.findIndexPage(s, "index.html")
	}

	return s, err
}

func (f wrappedFile) Readdir(count int) ([]fs.FileInfo, error) {
	fi, err := f.File.Readdir(count)
	if err != nil {
		return nil, err
	}

	if f.hideDotFiles {
		filtered := []fs.FileInfo{}

		for _, f := range fi {
			if !strings.HasPrefix(f.Name(), ".") {
				filtered = append(filtered, f)
			}
		}

		return filtered, err
	}

	return fi, nil
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
