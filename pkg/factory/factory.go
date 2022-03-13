package factory

import (
	"github.com/forando/refactory/pkg/filesystem"
	"github.com/forando/refactory/pkg/schema"
	_ "github.com/forando/refactory/pkg/schema"
	"log"
	"path/filepath"
)

func Bootstrap(accountModules *schema.AccountModules, pSetModules *schema.PermissionSetModules, output string) {
	const terragruntFileNme = "terragrunt.hcl"
	fs := filesystem.NewOsFs()
	for _, module := range *accountModules {
		dirPath := filepath.Join(output, module.ProductTicket, module.AccountName)
		fs.MakeDirs(dirPath)

		filePath := filepath.Join(dirPath, terragruntFileNme)
		BootstrapAccountTerragrunt(filePath, module)
	}
	for _, module := range *pSetModules {
		checkDirExists(&fs, filepath.Join(output, module.ProductTicket))

		dirPath := filepath.Join(output, module.ProductTicket, "PermissionSets", module.PermissionSetName)
		fs.MakeDirs(dirPath)

		filePath := filepath.Join(dirPath, terragruntFileNme)
		BootstrapPermissionSetTerragrunt(filePath, module)
	}
}

func checkDirExists(fs *filesystem.FS, path string) {
	if exists, err := (*fs).Exists(path); err != nil || !exists {
		if !exists {
			log.Fatalf("dir: [%s] expected but not found", path)
		} else {
			log.Fatal(err)
		}
	}
}
