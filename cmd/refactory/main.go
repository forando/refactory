package main

import (
	"flag"
	"github.com/cheggaaa/pb/v3"
	"github.com/forando/refactory/pkg/factory"
	"github.com/forando/refactory/pkg/filesystem"
	"github.com/forando/refactory/pkg/parser"
	"github.com/forando/refactory/pkg/schema"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"log"
	"strings"
)

func main() {
	inputDirFlag := flag.String("dir", "", "path to a dir to scan")
	outputDirFlag := flag.String("out", "", "path to a dir for the output")

	flag.Parse()

	if len(*inputDirFlag) == 0 {
		log.Fatal("requires path to a dir. please use the -dir flag to set it")
	}

	if len(*outputDirFlag) == 0 {
		log.Fatal("requires path to an output. please use the -out flag to set it")
	}

	fList, err := filesystem.GetTerraformFileNameList(*inputDirFlag)
	if err != nil {
		log.Fatal(err)
	}
	bar := pb.StartNew(len(fList))
	for _, file := range fList {
		body := parser.ParseFile(file)

		policyDocuments := parsePolicyDocuments(body, file)

		permissionSetModules, permissionSetNames := parsePermissionSets(body, file, policyDocuments)

		accountModules := parseAccounts(body, file, permissionSetNames)
		factory.Bootstrap(accountModules, permissionSetModules, *outputDirFlag)
		bar.Increment()
	}
	bar.Finish()
	log.Println("All Done!!!")
}

func parsePolicyDocuments(body *hclwrite.Body, file string) *map[string]*hclwrite.Block {
	documents := make(map[string]*hclwrite.Block)
	for _, block := range body.Blocks() {
		blockMetaData, err := parser.ParseBlockType(block)
		if err != nil {
			log.Fatalf("File %s: %s", file, err)
		}
		if blockMetaData.BlockType == schema.IamPolicyDocumentType {
			documents[blockMetaData.BlockName] = block
		}
	}
	return &documents
}

func parsePermissionSets(body *hclwrite.Body, file string, policyDocuments *map[string]*hclwrite.Block) (*schema.PermissionSetModules, *map[string]string) {
	fileTokens := strings.Split(file, "/")
	fileTokens = strings.Split(fileTokens[len(fileTokens)-1], ".")
	productTicket := fileTokens[0]
	permissionSetNames := make(map[string]string)
	var permissionSetModules schema.PermissionSetModules
	policyDocumentsCrossCheckList := make(map[string]bool)
	for _, block := range body.Blocks() {
		blockMetaData, err := parser.ParseBlockType(block)
		if err != nil {
			log.Fatalf("File %s: %s", file, err)
		}
		if blockMetaData.BlockType == schema.PermissionSetModuleType {
			module, err := parser.ParsePermissionSetModule(block.Body(), policyDocuments)
			if err != nil {
				log.Fatalf("File %s: %s", file, err)
			}
			module.ProductTicket = productTicket
			permissionSetModules = append(permissionSetModules, module)
			permissionSetNames[blockMetaData.BlockName] = module.PermissionSetName

			if len(module.PolicyDocumentName) > 0 {
				if _, found := policyDocumentsCrossCheckList[module.PolicyDocumentName]; found {
					log.Fatalf("File %s: multiple reference of %s policyDocument", file, module.PolicyDocumentName)
				}
				policyDocumentsCrossCheckList[module.PolicyDocumentName] = true
			}
		}
	}

	for pDoc := range *policyDocuments {
		if _, found := policyDocumentsCrossCheckList[pDoc]; !found {
			log.Fatalf("File %s: no reference for %s policyDocument", file, pDoc)
		}
	}

	return &permissionSetModules, &permissionSetNames
}

func parseAccounts(body *hclwrite.Body, file string, permissionSetNames *map[string]string) *schema.AccountModules {
	var modules schema.AccountModules
	for _, block := range body.Blocks() {
		blockMetaData, err := parser.ParseBlockType(block)
		if err != nil {
			log.Fatalf("File %s: %s", file, err)
		}
		if blockMetaData.BlockType == schema.AccountModuleType {
			module, err := parser.ParseAccountModule(block.Body(), permissionSetNames)
			if err != nil {
				log.Fatalf("File %s: %s", file, err)
			}
			modules = append(modules, module)
		}
	}
	return &modules
}
