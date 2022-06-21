package shellexec

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"os/exec"
)

type CmdRunner interface {
	Init() error
	StatePull() (*[]byte, error)
	StateList() (*[]string, error)
	StateMove(src string, dest string, dryRun bool) error
	StateImport(address string, id string) error
	StateRemove(address string, dryRun bool) error
}

type iacTool struct {
	Name string
	Dir  string
}

type OutPutType int

const (
	StdOut   OutPutType = 0
	StdError OutPutType = 1
)

type OutPut struct {
	Type  OutPutType
	Bytes []byte
}

func GetCmdRunner(name string, path string) CmdRunner {
	if name == "terragrunt" {
		return NewTerragrunt(path)
	}
	return NewTerraform(path)
}

func (t *iacTool) Run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = t.Dir
	return cmd.Run()
}

func (t *iacTool) RunWitStdOutput(name string, args ...string) (*[]byte, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = t.Dir

	output, err := cmd.Output()
	if err != nil {
		if output != nil {
			return nil, errors.Errorf("Dir: %s, error: %s", cmd.Dir, string(output))
		}
		return nil, err
	}
	return &output, nil
}

func (t *iacTool) RunWithCombinedOutput(name string, args ...string) (*[]byte, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = t.Dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		if output != nil {
			return nil, errors.Errorf("Dir: %s, error: %s", cmd.Dir, string(output))
		}
		return nil, err
	}
	return &output, nil
}

func (t *iacTool) RunWithOutputChannel(ch chan<- *OutPut, name string, args ...string) {
	go func() {
		cmd := exec.Command(name, args...)
		cmd.Dir = t.Dir
		var err error
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			t.sendError(ch, err.Error())
			close(ch)
			return
		}
		reader := bufio.NewReader(stdout)

		var stderrBuf bytes.Buffer
		cmd.Stderr = io.MultiWriter(&stderrBuf)

		go func(reader io.Reader, ch chan<- *OutPut) {
			scanner := bufio.NewScanner(reader)
			for scanner.Scan() {
				output := scanner.Bytes()
				ch <- &OutPut{
					Type:  StdOut,
					Bytes: output,
				}
			}
		}(reader, ch)

		if err := cmd.Start(); err != nil {
			t.sendError(ch, err.Error())
			close(ch)
			return
		}

		if err := cmd.Wait(); err != nil {
			errMsg := err.Error()
			if len(stderrBuf.Bytes()) > 0 {
				errMsg = string(stderrBuf.Bytes())
			}
			t.sendError(ch, errMsg)
		}
		close(ch)
	}()
}

func (t *iacTool) sendError(ch chan<- *OutPut, errMsg string) {
	out := OutPut{
		Type:  StdError,
		Bytes: []byte(fmt.Sprintf("Dir: %s, error: %s", t.Dir, errMsg)),
	}
	ch <- &out
}
