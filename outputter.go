package hackbrowserdata

import (
	"encoding/csv"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/gocarina/gocsv"
	jsoniter "github.com/json-iterator/go"
)

type outputType int

const (
	OutputJson outputType = iota + 1
	OutputCSV
)

func (o outputType) String() string {
	switch o {
	case OutputJson:
		return ".json"
	case OutputCSV:
		return ".csv"
	default:
		return "unknown type"
	}
}

type OutPutter struct {
	OutputType outputType
}

func NewOutPutter(outputType outputType) *OutPutter {
	return &OutPutter{OutputType: outputType}
}

func (o *OutPutter) Write(data BrowsingData, writer *os.File) error {
	switch o.OutputType {
	case OutputCSV:
		gocsv.SetCSVWriter(func(w io.Writer) *gocsv.SafeCSVWriter {
			writer := csv.NewWriter(w)
			writer.Comma = ','
			return gocsv.NewSafeCSVWriter(writer)
		})
		return gocsv.MarshalFile(data, writer)
	case OutputJson:
		encoder := jsoniter.NewEncoder(writer)
		encoder.SetIndent("  ", "  ")
		encoder.SetEscapeHTML(false)
		return encoder.Encode(data)
	}
	return nil
}

func (o *OutPutter) CreateFile(filename string) (*os.File, error) {
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
	file, err = os.OpenFile(filename, os.O_TRUNC|os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return file, nil
}
