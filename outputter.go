package hackbrowserdata

import (
	"errors"
	"os"
	"path/filepath"
)

type OutPutter struct {
	Json bool
}

func OutPutToJson(data BrowsingData) (err error) {
	switch data {
	case &WebkitPassword{}:

	}
	return err
}

func (o *OutPutter) createFile(filename string, appendtoFile bool) (*os.File, error) {
	if filename == "" {
		return nil, errors.New("empty filename")
	}

	dir := filepath.Dir(filename)

	if dir != "" {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err := os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				return nil, err
			}
		}
	}

	var file *os.File
	var err error
	if appendtoFile {
		file, err = os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	} else {
		file, err = os.Create(filename)
	}
	if err != nil {
		return nil, err
	}

	return file, nil
}
