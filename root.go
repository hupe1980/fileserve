package main

import (
	"fmt"
	"net/http"
	"os"
)

func Root(dir string) (http.Dir, error) {
	fileInfo, err := os.Stat(dir)
	if err != nil {
		return "", fmt.Errorf("cannot serve %v: %v", dir, err)
	} else if !fileInfo.IsDir() {
		return "", fmt.Errorf("cannot serve %v: not a directory", dir)
	} else {
		file, err := os.Open(dir)
		if err != nil {
			return "", fmt.Errorf("cannot serve %v: %v", dir, err)
		}
		file.Close()
	}

	return http.Dir(dir), nil
}
