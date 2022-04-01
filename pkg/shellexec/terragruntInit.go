package shellexec

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log"
	"os/exec"
	"regexp"
)

var createS3Regexp = regexp.MustCompile(`Remote state S3 bucket.+does not exist.+Would you like Terragrunt to create it`)

func ExecTerragruntInitWithStdIn(dir string) error {
	var err error
	cmd := exec.Command("terragrunt", "init")
	cmd.Dir = dir

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return errors.Errorf("Dir: %s, error: %s", dir, err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return errors.Errorf("Dir: %s, error: %s", dir, err)
	}
	reader := bufio.NewReader(stdout)

	var stderrBuf bytes.Buffer
	cmd.Stderr = io.MultiWriter(&stderrBuf)

	done := make(chan string)

	go func(reader io.Reader) {
		scanner := bufio.NewScanner(reader)
		var output string
		for scanner.Scan() {
			currentOutput := scanner.Text()
			output = fmt.Sprintf("%s%s", output, currentOutput)
			if createS3Regexp.MatchString(currentOutput) {
				if _, err := stdin.Write([]byte("y\n")); err != nil {
					output = fmt.Sprintf("%s error:\n%s", output, err)
					break
				}
			}
		}
		done <- output
	}(reader)

	if err := cmd.Start(); err != nil {
		return errors.Errorf("Dir: %s, error: %s", dir, err)
	}
	output := <-done
	log.Printf("Dir: %s, output: %s", dir, output)

	if err := cmd.Wait(); err != nil {
		if len(stderrBuf.Bytes()) > 0 {
			output = fmt.Sprintf("%s\n%s", output, string(stderrBuf.Bytes()))
		}
		return errors.Errorf("Dir: %s, error: %s", dir, err)
	}
	return nil
}

func ExecTerragruntInit(dir string) error {
	cmd := exec.Command("terragrunt", "init")
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Errorf("Dir: %s, error: %s", dir, string(output))
	}
	return nil
}
