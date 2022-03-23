package shellexec

import (
	"github.com/forando/refactory/pkg/schema"
	"os/exec"
)

func ExecTerragruntRollBackImports(dir string, imports *[]schema.Import) {
	for _, imp := range *imports {
		cmd := exec.Command("terragrunt", "state", "rm", imp.Address)
		cmd.Dir = dir
		cmd.Run()
	}
}
