package factory

import (
	"github.com/forando/refactory/pkg/filesystem"
	"github.com/forando/refactory/pkg/schema"
	_ "github.com/forando/refactory/pkg/schema"
	"log"
	"path/filepath"
)

func Bootstrap(accountModules *schema.AccountModules, pSetModules *schema.PermissionSetModules, importResources *map[string][]schema.TfImport, output string, org schema.Org) {
	const terragruntFileNme = "terragrunt.hcl"
	fs := filesystem.NewOsFs()
	filePath := filepath.Join(output, terragruntFileNme)

	BootstrapRootTerragrunt(filePath, org)
	bootstrapAccounts(accountModules, output, &fs, importResources, filePath, terragruntFileNme)
	bootstrapPermissionStets(pSetModules, &fs, output, importResources, terragruntFileNme, org)
}

func bootstrapAccounts(accountModules *schema.AccountModules, output string, fs *filesystem.FS, importResources *map[string][]schema.TfImport, filePath string, terragruntFileNme string) {
	for _, module := range *accountModules {

		dirPath := filepath.Join(output, module.ProductTicket)
		(*fs).MakeDirs(dirPath)

		dirPath = filepath.Join(dirPath, module.AccountName)
		(*fs).MakeDirs(dirPath)

		BootstrapImports(dirPath, module.ModuleName, importResources)

		filePath = filepath.Join(dirPath, terragruntFileNme)
		BootstrapAccountTerragrunt(filePath, module)
	}
}

func bootstrapPermissionStets(pSetModules *schema.PermissionSetModules, fs *filesystem.FS, output string, importResources *map[string][]schema.TfImport, terragruntFileNme string, org schema.Org) {
	for _, module := range *pSetModules {
		path := filepath.Join(output, module.ProductTicket)
		checkDirExists(fs, path)

		dirPath := filepath.Join(output, module.ProductTicket, "PermissionSets", module.PermissionSetName)
		(*fs).MakeDirs(dirPath)

		BootstrapImports(dirPath, module.ModuleName, importResources)

		filePath := filepath.Join(dirPath, terragruntFileNme)
		BootstrapPermissionSetTerragrunt(filePath, module, org)
	}
}

func checkDirExists(fs *filesystem.FS, path string) {
	if exists, err := (*fs).Exists(path); !exists {
		if err != nil {
			log.Fatal(err)
		}
		log.Fatalf("Dir %s dos not exist", path)
	}
}
