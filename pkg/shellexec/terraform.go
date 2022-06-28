package shellexec

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/pkg/errors"
)

type Terraform struct {
	iacTool
}

func NewTerraform(dir string) *Terraform {
	return &Terraform{iacTool{Name: "terraform", Dir: dir}}
}

func (t *Terraform) Init(backendConfig string) error {
	data := make(chan *OutPut)
	fmt.Printf("Initialazing %s...\n", t.Name)
	if len(backendConfig) > 0 {
		backendConfigFlag := fmt.Sprintf("-backend-config=%s", backendConfig)
		t.RunWithOutputChannel(data, t.Name, "init", "-input=false", backendConfigFlag)
	} else {
		t.RunWithOutputChannel(data, t.Name, "init", "-input=false")
	}
	for output := range data {
		if output.Type == StdError {
			return errors.New(string(output.Bytes))
		}
		fmt.Println(string(output.Bytes))
	}
	return nil
}

func (t *Terraform) StatePull() (*[]byte, error) {
	return t.RunWitStdOutput(t.Name, "state", "pull")
}

func (t *Terraform) StateList() (*[]string, error) {
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

func (t *Terraform) StateRemove(address string, dryRun bool) error {
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

func (t *Terraform) StateMove(src string, dest string, dryRun bool) error {
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

func (t *Terraform) StateImport(address string, id string) error {
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
