package shellexec

import (
	"fmt"
	"github.com/pkg/errors"
)

type Terraform struct {
	iacTool
}

func NewTerraform(dir string) *Terraform {
	return &Terraform{iacTool{Name: "terraform", Dir: dir}}
}

func (t *Terraform) Init() error {
	data := make(chan *OutPut)
	fmt.Printf("Initialazing %s...\n", t.Name)
	t.RunWithOutputChannel(data, t.Name, "init", "-input=false")
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
		fmt.Printf("%s state rm --dry-run %q\n", t.Name, address)
		t.RunWithOutputChannel(data, t.Name, "state", "rm", "-dry-run", address)
	} else {
		fmt.Printf("%s state rm %q\n", t.Name, address)
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
		fmt.Printf("%s state mv --dry-run %q %q\n", t.Name, src, dest)
		t.RunWithOutputChannel(data, t.Name, "state", "mv", "-dry-run", src, dest)
	} else {
		fmt.Printf("%s state mv %q %q\n", t.Name, src, dest)
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
	fmt.Printf("%s import %q %q\n", t.Name, address, id)
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
