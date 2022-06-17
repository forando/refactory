package shellexec

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"os/exec"
)

type OutPutType int

const (
	StdOut   OutPutType = 0
	StdError OutPutType = 1
)

type OutPut struct {
	Type  OutPutType
	Bytes []byte
}

type CmdRunner struct {
	Dir string
}

func (r *CmdRunner) Run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = r.Dir
	return cmd.Run()
}

func (r *CmdRunner) RunWitStdOutput(name string, args ...string) (*[]byte, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = r.Dir

	output, err := cmd.Output()
	if err != nil {
		if output != nil {
			return nil, errors.Errorf("Dir: %s, error: %s", cmd.Dir, string(output))
		}
		return nil, err
	}
	return &output, nil
}

func (r *CmdRunner) RunWithCombinedOutput(name string, args ...string) (*[]byte, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = r.Dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		if output != nil {
			return nil, errors.Errorf("Dir: %s, error: %s", cmd.Dir, string(output))
		}
		return nil, err
	}
	return &output, nil
}

func (r *CmdRunner) RunWithOutputChannel(ch chan<- *OutPut, name string, args ...string) {
	go func() {
		cmd := exec.Command(name, args...)
		cmd.Dir = r.Dir
		var err error
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			r.sendError(ch, err.Error())
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
			r.sendError(ch, err.Error())
			close(ch)
			return
		}

		if err := cmd.Wait(); err != nil {
			errMsg := err.Error()
			if len(stderrBuf.Bytes()) > 0 {
				errMsg = string(stderrBuf.Bytes())
			}
			r.sendError(ch, errMsg)
		}
		close(ch)
	}()
}

func (r *CmdRunner) sendError(ch chan<- *OutPut, errMsg string) {
	out := OutPut{
		Type:  StdError,
		Bytes: []byte(fmt.Sprintf("Dir: %s, error: %s", r.Dir, errMsg)),
	}
	ch <- &out
}
