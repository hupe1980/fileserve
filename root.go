package fileserve

import (
	"fmt"
	"net/http"
	"os"
)

func NewDirRoot(dir string) (http.FileSystem, error) {
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

	return http.Dir(dir), nil
}
