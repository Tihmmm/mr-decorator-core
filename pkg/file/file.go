package file

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

func ReadCsv(dir, subpath string) ([][]string, error) {
	root, err := os.OpenRoot(dir)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error opening directory root: %s\n", err))
	}
	defer func(root *os.Root) {
		if err := root.Close(); err != nil {
			log.Printf("Error closing root: %s\n", err)
			return
		}
	}(root)

	file, err := root.Open(subpath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error opening csv file: %s\n", err))
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("Error closing file: %s\n", err)
			return
		}
	}(file)

	r := csv.NewReader(file)
	records := make([][]string, 0)
	for {
		record, err := r.Read()
		if errors.Is(err, io.EOF) || errors.Is(err, csv.ErrBareQuote) || errors.Is(err, csv.ErrFieldCount) || errors.Is(err, csv.ErrQuote) {
			break
		}
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Error parsing csv: %s\n", err))
		}
		records = append(records, record)
	}

	return records, nil
}

func ParseJsonFile(dir, subpath string, dest any) error {
	root, err := os.OpenRoot(dir)
	if err != nil {
		return errors.New(fmt.Sprintf("Error opening directory root: %s\n", err))
	}
	defer func(root *os.Root) {
		if err := root.Close(); err != nil {
			log.Printf("Error closing root: %s\n", err)
			return
		}
	}(root)

	jsonFile, err := root.Open(subpath)
	if err != nil {
		return errors.New(fmt.Sprintf("Error opening jsonFile file: %s, err: %s\n", subpath, err))
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("Error closing file: %s\n", err)
			return
		}
	}(jsonFile)

	jsonParser := json.NewDecoder(jsonFile)
	if err = jsonParser.Decode(dest); err != nil {
		return errors.New(fmt.Sprintf("Error decoding json file: %s, err: %s\n", subpath, err))
	}

	return nil
}

func DeleteDirectory(dir string) error {
	err := os.RemoveAll(dir)
	if err != nil {
		return errors.New(fmt.Sprintf("Error deleting directory %s: %s\n", dir, err))
	}

	return nil
}
