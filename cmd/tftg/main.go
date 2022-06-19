package main

import (
	"flag"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/forando/refactory/pkg/factory"
	"github.com/forando/refactory/pkg/filesystem"
	"github.com/forando/refactory/pkg/parser"
	"github.com/forando/refactory/pkg/resourceimport"
	"github.com/forando/refactory/pkg/schema"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"log"
	"os"
	"strings"
)

func main() {

	inputDirFlag := flag.String("dir", "", "path to a dir to scan")
	inputDirListFileFlag := flag.String("dir-list", "", "path to csv file with dirList to be imported")
	stateFileFlag := flag.String("state", "", "path to a terraform.json state file")
	outputDirFlag := flag.String("out", "", "path to a dir for the output")
	orgFlag := flag.String("org", "", "AWS Organization to make refactoring for (test or prod)")

	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		fmt.Printf("%s [FLAGS...] [bootstrap|import]\n", os.Args[0])
		fmt.Println("COMMANDS:")
		fmt.Println("  bootstrap: bootstraps terragunt file structure from terraform files")
		fmt.Println("  import: runs `terragunt import ...` command for each bootstrapped terragrunt module")
		fmt.Println("FLAGS:")
		flag.PrintDefaults()
	}

	flag.Parse()

	command := flag.Arg(0)
	if len(command) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	switch command {
	case "bootstrap":
		if len(*inputDirFlag) == 0 {
			log.Fatal("requires path to a dir. please use the -dir flag to set it")
		}
		if len(*stateFileFlag) == 0 {
			log.Fatal("requires path to a terraform.json state file. please use the -state flag to set it")
		}
		if len(*outputDirFlag) == 0 {
			log.Fatal("requires path to an output. please use the -out flag to set it")
		}
		if len(*orgFlag) == 0 {
			log.Fatal("requires Organization. please use the -org flag to set it")
		}
		var org schema.Org
		switch *orgFlag {
		case "prod":
			org = schema.ProdOrg
		case "test":
			org = schema.TestOrg
		default:
			log.Fatal("-org flag can be set either to prod or to test")
		}

		bootstrap(*inputDirFlag, *stateFileFlag, *outputDirFlag, org)
	case "import":
		if len(*inputDirFlag) == 0 && len(*inputDirListFileFlag) == 0 {
			log.Fatal("requires either path to a dir or path to a file with dirList. please use either -dir or -dir-list flag to set it")
		}
		importState(*inputDirFlag, *inputDirListFileFlag)
	default:
		log.Printf("Unknown command [%s]", command)
		flag.Usage()
		os.Exit(1)
	}
}

func importState(dir string, file string) {
	var names *[]string
	var err error
	if len(dir) > 0 {
		names, err = filesystem.GetTerragruntModuleNameList(dir)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		names, err = parser.ParseDirs(file)
		if err != nil {
			log.Fatal(err)
		}
	}

	ch, emptyResources := resourceimport.Start(*names)
	all := len(*names)
	log.Printf("All directories: %v", all)
	for _, emptyResource := range *emptyResources {
		log.Printf("Dir: %s, error: %s", emptyResource.Dir, emptyResource.Message)
	}
	donesOk := make([]schema.Done, 0)
	donesErr := make([]schema.Done, 0)
	cleanErrors := make([]string, 0)
	for res := range ch {
		if res.Status == schema.Ok {
			donesOk = append(donesOk, res)
		} else {
			donesErr = append(donesErr, res)
		}
		cleanErrors = append(cleanErrors, res.CleanErrors)
	}
	log.Printf("Processed dirs: %v; ok: %v; err: %v", all, len(donesOk), len(donesErr))
	for _, doneErr := range donesErr {
		log.Printf("Dir [%s], Resource [%s], Id [%s], error: %s", doneErr.Dir, doneErr.FailedResource.Address, doneErr.FailedResource.Id, doneErr.Message)
	}
	log.Println("Errors during cleaning:")
	for _, cleanError := range cleanErrors {
		log.Println(cleanError)
	}
}

func bootstrap(inputDir string, stateFile string, outputDir string, org schema.Org) {

	fList, err := filesystem.GetTerraformFileNameList(inputDir)
	if err != nil {
		log.Fatal(err)
	}
	bar := pb.StartNew(len(fList))
	imports, err := parser.ParseTfState(stateFile)
	if err != nil {
		log.Fatalf("Cannot parse tfState; %s", err)
	}

	moduleNames := make(map[string]bool)
	for _, file := range fList {
		body := parser.ParseFile(file)

		policyDocuments := parsePolicyDocuments(body, file)

		permissionSetModules, permissionSetNames := parsePermissionSets(body, file, policyDocuments)

		accountModules := parseAccounts(body, file, permissionSetNames)

		for _, module := range *accountModules {
			moduleNames[module.ModuleName] = true
		}
		for _, module := range *permissionSetModules {
			moduleNames[module.ModuleName] = true
		}

		factory.Bootstrap(accountModules, permissionSetModules, imports, outputDir, org)
		bar.Increment()
	}

	checkConsistency(&moduleNames, imports)

	bar.Finish()
	log.Println("All Done!!!")
}

func parsePolicyDocuments(body *hclwrite.Body, file string) *map[string]*schema.PolicyDocument {
	documents := make(map[string]*schema.PolicyDocument)
	for _, block := range body.Blocks() {
		blockMetaData, err := parser.ParseBlockType(block)
		if err != nil {
			log.Fatalf("File %s: %s", file, err)
		}
		if blockMetaData.BlockType == schema.IamPolicyDocumentType {
			document, err := parser.ParsePolicyDocumentBlock(block.Body().BuildTokens(nil).Bytes())
			if err != nil {
				log.Fatalf("Cannot parse policyDocument %s: %s", blockMetaData.BlockName, err)
			}
			documents[blockMetaData.BlockName] = document
		}
	}
	return &documents
}

func parsePermissionSets(body *hclwrite.Body, file string, policyDocuments *map[string]*schema.PolicyDocument) (*schema.PermissionSetModules, *map[string]string) {
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
			module.ModuleName = blockMetaData.BlockName
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
			module.ModuleName = blockMetaData.BlockName
			modules = append(modules, module)
		}
	}
	return &modules
}

func checkConsistency(moduleNames *map[string]bool, imports *map[string][]schema.TfImport) {
	moduleCount := len(*moduleNames)
	importsCount := len(*imports)

	if moduleCount > importsCount {
		for moduleName := range *moduleNames {
			if _, found := (*imports)[moduleName]; !found {
				log.Printf("stateCount = %v not equal to moduleCount = %v", importsCount, moduleCount)
				log.Fatalf("missing %s module in state", moduleName)
			}
		}

	}
	if moduleCount < importsCount {
		for imp := range *imports {
			if _, found := (*moduleNames)[imp]; !found {
				log.Printf("stateCount = %v not equal to moduleCount = %v", importsCount, moduleCount)
				log.Fatalf("missing %s module in modules configuration", imp)
			}
		}

	}
}
