package shellexec

import (
	"github.com/forando/refactory/pkg/schema"
	"github.com/pkg/errors"
	"os/exec"
)

type Terragrunt struct {
	Dir string
}

func (t *Terragrunt) Init() error {
	cmd := exec.Command("terragrunt", "init", "--terragrunt-non-interactive", "-input=false")
	cmd.Dir = t.Dir

	output, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Errorf("Dir: %s, error: %s", t.Dir, string(output))
	}
	return nil
}

func (t *Terragrunt) Import(address string, id string) (*string, error) {
	cmd := exec.Command("terragrunt", "import", address, id)
	cmd.Dir = t.Dir

	outputBytes, err := cmd.CombinedOutput()
	output := string(outputBytes)
	if err != nil {
		return nil, errors.Errorf("Dir: %s, error: %s", t.Dir, output)
	}
	return &output, nil
}

func (t *Terragrunt) RollBackImports(imports *[]schema.Import) {
	for _, imp := range *imports {
		cmd := exec.Command("terragrunt", "state", "rm", imp.Address)
		cmd.Dir = t.Dir
		cmd.Run()
	}
}
