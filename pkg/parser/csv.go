package parser

import (
	"encoding/csv"
	"os"
)

func readData(fileName string) ([][]string, error) {

	f, err := os.Open(fileName)

	if err != nil {
		return [][]string{}, err
	}

	//TODO: potential cause for `too many open files` error
	defer f.Close()

	r := csv.NewReader(f)
	r.Comma = ';'

	records, err := r.ReadAll()

	if err != nil {
		return [][]string{}, err
	}

	return records, nil
}
