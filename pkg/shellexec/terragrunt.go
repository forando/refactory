package shellexec

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/forando/refactory/pkg/schema"
	"github.com/pkg/errors"
)

type Terragrunt struct {
	iacTool
}

func NewTerragrunt(dir string) *Terragrunt {
	return &Terragrunt{iacTool{Name: "terragrunt", Dir: dir}}
}

func (t *Terragrunt) Init() error {
	data := make(chan *OutPut)
	fmt.Printf("Initialazing %s...\n", t.Name)
	t.RunWithOutputChannel(data, t.Name, "init", "--terragrunt-non-interactive", "-input=false")
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
	return t.RunWitStdOutput(t.Name, "state", "pull")
}

func (t *Terragrunt) StateList() (*[]string, error) {
	data := make(chan *OutPut)
	t.RunWithOutputChannel(data, t.Name, "state", "list")
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
		color.Cyan("%s state rm --dry-run %q\n", t.Name, address)
		t.RunWithOutputChannel(data, t.Name, "state", "rm", "-dry-run", address)
	} else {
		color.Cyan("%s state rm %q\n", t.Name, address)
		t.RunWithOutputChannel(data, t.Name, "state", "rm", address)
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
		color.Cyan("%s state mv --dry-run %q %q\n", t.Name, src, dest)
		t.RunWithOutputChannel(data, t.Name, "state", "mv", "-dry-run", src, dest)
	} else {
		color.Cyan("%s state mv %q %q\n", t.Name, src, dest)
		t.RunWithOutputChannel(data, t.Name, "state", "mv", src, dest)
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

func (t *Terragrunt) StateImport(address string, id string) error {
	data := make(chan *OutPut)
	color.Cyan("%s import %q %q\n", t.Name, address, id)
	t.RunWithOutputChannel(data, t.Name, "import", address, id)
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
	outputBytes, err := t.RunWithCombinedOutput(t.Name, "import", address, id)
	output := string(*outputBytes)
	return &output, err
}

func (t *Terragrunt) RollBackImports(imports *[]schema.Import) {
	for _, imp := range *imports {
		t.Run(t.Name, "state", "rm", imp.Address)
	}
}
