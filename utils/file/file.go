package file

import (
	"encoding/csv"
	"os"
	"path/filepath"
)

func createFile(path string) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return
	}
	var _, err = os.Stat(path)
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if err != nil {
			return
		}
		defer file.Close()
	}
}

func ReadCsvFile(path string) ([][]string, error) {
	createFile(path)
	f, err := os.Open(path)
	if err != nil {
		return [][]string{}, err
	}
	defer f.Close()
	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return [][]string{}, err
	}
	return records, nil
}

func WriteCsvFile(path string, records [][]string) error {
	createFile(path)
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		return err
	}
	w := csv.NewWriter(f)
	err = w.WriteAll(records)
	return err
}
