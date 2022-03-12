package parser

import (
	"fmt"
	"github.com/forando/refactory/pkg/filesystem"
	"github.com/forando/refactory/pkg/schema"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"log"
	"strings"
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

func ParseBlockType(block *hclwrite.Block) (*schema.BlockMetaData, error) {

	switch block.Type() {
	case schema.DataBlock:
		return parseDataBlock(block)
	case schema.ModuleBlock:
		return parseModuleBlock(block)
	default:
		return nil, schema.ParsingError{Message: fmt.Sprintf("block of type [%s] found but not expected", block.Type())}
	}
}

func parseDataBlock(block *hclwrite.Block) (*schema.BlockMetaData, error) {
	labels := block.Labels()

	if len(labels) == 0 {
		return nil, schema.ParsingError{Message: "Data Block does not have labels"}
	}

	if labels[0] != schema.IamPolicyDocumentType {
		return nil, schema.ParsingError{Message: fmt.Sprintf("Unexpected Data Block Type: [%s]", labels[0])}
	}

	return &schema.BlockMetaData{BlockType: schema.IamPolicyDocumentType, BlockName: labels[0]}, nil
}

func parseModuleBlock(block *hclwrite.Block) (*schema.BlockMetaData, error) {

	labels := block.Labels()

	if len(labels) == 0 {
		return nil, schema.ParsingError{Message: "Module Block does not have labels"}
	}

	attr := block.Body().Attributes()

	key := "source"
	sourceAttr := attr[key]
	if sourceAttr == nil {
		return nil, schema.ParsingError{Message: fmt.Sprintf("Module block [%s] does not have [%s] attribute", labels[0], key)}
	}

	source := getExpressionAsString(sourceAttr)
	tokens := strings.Split(source, "/")

	moduleType := tokens[len(tokens)-1]

	switch moduleType {
	case schema.AccountModuleType:
		return &schema.BlockMetaData{BlockType: schema.AccountModuleType, BlockName: labels[0]}, nil
	case schema.PermissionSetModuleType:
		return &schema.BlockMetaData{BlockType: schema.PermissionSetModuleType, BlockName: labels[0]}, nil
	default:
		return nil, schema.ParsingError{Message: fmt.Sprintf("Unexpected Module Block Type: [%s]", moduleType)}
	}
}
