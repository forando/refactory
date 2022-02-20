package parser

import (
	"github.com/forando/refactory/pkg/filesystem"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"log"
)

func ParseFile(fileName string) *hclwrite.Body {
	fs := filesystem.NewOsFs()
	b, err := fs.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Cannot read file %s: %v", fileName, err)
	}

	file, parsedDiags := hclwrite.ParseConfig(b, fileName, hcl.InitialPos)

	if parsedDiags.HasErrors() {
		log.Fatal("Cannot parse file ", parsedDiags.Error())
	}

	return file.Body()
}
