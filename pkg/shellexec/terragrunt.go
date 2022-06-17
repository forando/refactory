package shellexec

import (
	"fmt"
	"github.com/forando/refactory/pkg/schema"
	"github.com/pkg/errors"
)

type Terragrunt struct {
	CmdRunner
}

func NewTerragrunt(dir string) *Terragrunt {
	return &Terragrunt{CmdRunner{Dir: dir}}
}

func (t *Terragrunt) Name() string {
	return "terragrunt"
}

func (t *Terragrunt) Init() error {
	data := make(chan *OutPut)
	t.RunWithOutputChannel(data, "terragrunt", "init", "--terragrunt-non-interactive", "-input=false")
	for output := range data {
		msg := string(output.Bytes)
		if output.Type == StdError {
			return errors.New(msg)
		}
		fmt.Println(msg)
	}
	return nil
}

func (t *Terragrunt) StatePull() (*[]byte, error) {
	return t.RunWitStdOutput("terragrunt", "state", "pull")
}

func (t *Terragrunt) StateList() (*[]string, error) {
	data := make(chan *OutPut)
	t.RunWithOutputChannel(data, "terragrunt", "state", "list")
	resources := make([]string, 0)
	for output := range data {
		msg := string(output.Bytes)
		if output.Type == StdError {
			return nil, errors.New(msg)
		}
		fmt.Println(msg)
		resources = append(resources, msg)
	}
	return &resources, nil
}

func (t *Terragrunt) StateRemove(address string, dryRun bool) error {
	data := make(chan *OutPut)
	if dryRun {
		fmt.Printf("terragrunt state rm --dry-run %q\n", address)
		t.RunWithOutputChannel(data, "terragrunt", "state", "rm", "-dry-run", address)
	} else {
		fmt.Printf("terragrunt state rm %q\n", address)
		t.RunWithOutputChannel(data, "terragrunt", "state", "rm", address)
	}
	resources := make([]string, 0)
	for output := range data {
		msg := string(output.Bytes)
		if output.Type == StdError {
			return errors.New(msg)
		}
		fmt.Println(msg)
		resources = append(resources, msg)
	}
	return nil
}

func (t *Terragrunt) StateMove(src string, dest string, dryRun bool) error {
	data := make(chan *OutPut)
	if dryRun {
		fmt.Printf("terragrunt state mv --dry-run %q %q\n", src, dest)
		t.RunWithOutputChannel(data, "terragrunt", "state", "mv", "-dry-run", src, dest)
	} else {
		fmt.Printf("terragrunt state mv %q %q\n", src, dest)
		t.RunWithOutputChannel(data, "terragrunt", "state", "mv", src, dest)
	}
	resources := make([]string, 0)
	for output := range data {
		msg := string(output.Bytes)
		if output.Type == StdError {
			return errors.New(msg)
		}
		fmt.Println(msg)
		resources = append(resources, msg)
	}
	return nil
}

func (t *Terragrunt) Import(address string, id string) (*string, error) {
	outputBytes, err := t.RunWithCombinedOutput("terragrunt", "import", address, id)
	output := string(*outputBytes)
	return &output, err
}

func (t *Terragrunt) RollBackImports(imports *[]schema.Import) {
	for _, imp := range *imports {
		t.Run("terragrunt", "state", "rm", imp.Address)
	}
}
