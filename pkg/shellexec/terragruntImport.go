package shellexec

import (
	"fmt"
	"os/exec"
)

func ExecTerragruntImport(dir string, address string, id string) (*string, error) {
	cmd := exec.Command("terragrunt", "import", address, id)
	cmd.Dir = dir

	outputBytes, err := cmd.CombinedOutput()
	output := string(outputBytes)
	if err != nil {
		return nil, shellError{Message: fmt.Sprintf("Dir: %s, error: %s", dir, output)}
	}
	return &output, nil
}
