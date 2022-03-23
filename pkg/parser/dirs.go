package parser

import (
	"fmt"
	"github.com/forando/refactory/pkg/schema"
)

const dirRecordItems = 1

func ParseDirs(fileName string) (*[]string, error) {
	records, err := readData(fileName)

	if err != nil {
		return nil, err
	}

	dirs := make([]string, 0)

	for _, record := range records {
		if len(record) != dirRecordItems {
			return nil, schema.ParsingError{Message: fmt.Sprintf("CSV Dir Record != %v items: %s", dirRecordItems, record)}
		}
		dirs = append(dirs, record[0])
	}

	return &dirs, nil
}
