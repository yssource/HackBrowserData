package hackbrowserdata

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"path/filepath"

	jsongo "github.com/json-iterator/go"
	"github.com/jszwec/csvutil"
)

type outputType int

const (
	OutputJson outputType = iota + 1
	OutputCSV
)

type OutPutter struct {
	OutputType outputType
}

func NewOutPutter(outputType outputType) *OutPutter {
	return &OutPutter{OutputType: outputType}
}

func (o *OutPutter) Write(data BrowsingData, writer io.Writer) error {
	switch o.OutputType {
	case OutputCSV:
		encoder := csvutil.NewEncoder(csv.NewWriter(writer))
		return encoder.Encode(data)
	case OutputJson:
		encoder := jsongo.NewEncoder(writer)
		encoder.SetIndent("  ", "  ")
		encoder.SetEscapeHTML(false)
		return encoder.Encode(data)
	}
	return nil
}

func (o *OutPutter) CreateFile(filename string, appendtoFile bool) (*os.File, error) {
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
