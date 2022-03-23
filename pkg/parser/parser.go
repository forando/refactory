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

	dataBlockType := labels[0]
	if dataBlockType != schema.IamPolicyDocumentType {
		return nil, schema.ParsingError{Message: fmt.Sprintf("Unexpected Data Block Type: [%s]", dataBlockType)}
	}

	if len(labels) < 2 {
		return nil, schema.ParsingError{Message: fmt.Sprintf("Data Block: [%s] does not have name", dataBlockType)}
	}

	return &schema.BlockMetaData{BlockType: schema.IamPolicyDocumentType, BlockName: labels[1]}, nil
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

func ParseTfState(file string) (*map[string][]schema.TfImport, error) {
	state, err := parseTfStateFile(file)
	if err != nil {
		return nil, err
	}
	importArgs := make(map[string][]schema.TfImport)
	for _, resource := range state.Resources {
		if len(resource.Module) == 0 || resource.Type == "null_resource" {
			continue
		}
		var moduleName string
		moduleName, err = parseModuleName(resource.Module)
		if err != nil {
			return nil, err
		}

		addressTokens := make([]string, 0)
		switch resource.Mode {
		case "data":
			continue
		case "managed":
			addressTokens = append(addressTokens, resource.Type, resource.Name)
		default:
			return nil, schema.ParsingError{Message: fmt.Sprintf("module: %s has usupported mode: [%s]", resource.Module, resource.Mode)}
		}

		var tfImports []schema.TfImport
		found := false
		if tfImports, found = importArgs[moduleName]; !found {
			tfImports = make([]schema.TfImport, 0)
		}

		for _, instance := range resource.Instances {
			var id string
			address := strings.Join(addressTokens, ".")
			if resource.Type == "aws_s3_bucket_object" {
				if len(instance.Attrs.Bucket) == 0 {
					return nil, schema.ParsingError{Message: fmt.Sprintf("module: %s of type: %s does not have bucket property", resource.Module, resource.Type)}
				}
				id = fmt.Sprintf("%s/%s", instance.Attrs.Bucket, instance.Attrs.Id)
			} else {
				id = instance.Attrs.Id
			}
			if resource.Type == "aws_ssoadmin_account_assignment" {
				if len(instance.IndexKey) == 0 {
					return nil, schema.ParsingError{Message: fmt.Sprintf("module: %s of type: %s does not have index_key property", resource.Module, resource.Type)}
				}
				address = fmt.Sprintf("%s['%s']", address, instance.IndexKey)
			}
			tfImport := schema.TfImport{Address: address, Id: id}
			tfImports = append(tfImports, tfImport)
		}
		importArgs[moduleName] = tfImports
	}
	return &importArgs, nil
}
