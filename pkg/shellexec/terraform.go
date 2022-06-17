package shellexec

import (
	"fmt"
	"github.com/pkg/errors"
)

type Terraform struct {
	CmdRunner
}

func NewTerraform(dir string) *Terraform {
	return &Terraform{CmdRunner{Dir: dir}}
}

func (t *Terraform) Name() string {
	return "terraform"
}

func (t *Terraform) Init() error {
	data := make(chan *OutPut)
	t.RunWithOutputChannel(data, "terraform", "init", "-input=false")
	for output := range data {
		if output.Type == StdError {
			return errors.New(string(output.Bytes))
		}
		fmt.Println(string(output.Bytes))
	}
	return nil
}

func (t *Terraform) StatePull() (*[]byte, error) {
	return t.RunWitStdOutput("terraform", "state", "pull")
}

func (t *Terraform) StateList() (*[]string, error) {
	data := make(chan *OutPut)
	t.RunWithOutputChannel(data, "terraform", "state", "list")
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
		fmt.Printf("terraform state rm --dry-run %q\n", address)
		t.RunWithOutputChannel(data, "terraform", "state", "rm", "-dry-run", address)
	} else {
		fmt.Printf("terraform state rm %q\n", address)
		t.RunWithOutputChannel(data, "terraform", "state", "rm", address)
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
		fmt.Printf("terraform state mv --dry-run %q %q\n", src, dest)
		t.RunWithOutputChannel(data, "terraform", "state", "mv", "-dry-run", src, dest)
	} else {
		fmt.Printf("terraform state mv %q %q\n", src, dest)
		t.RunWithOutputChannel(data, "terraform", "state", "mv", src, dest)
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
