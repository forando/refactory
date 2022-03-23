package parser

import (
	"fmt"
	"github.com/forando/refactory/pkg/schema"
	"strings"
)

const importRecordItems = 3

func ParseImports(fileName string) (*[]schema.Import, error) {
	records, err := readData(fileName)

	if err != nil {
		return nil, err
	}

	imports := make([]schema.Import, 0)

	for _, record := range records {
		if len(record) != importRecordItems {
			return nil, schema.ParsingError{Message: fmt.Sprintf("CSV Import Record != %v items: %s", importRecordItems, record)}
		}
		imp := schema.Import{Module: record[0], Address: strings.ReplaceAll(record[1], "'", "\""), Id: record[2]}
		imports = append(imports, imp)
	}

	return &imports, nil
}
