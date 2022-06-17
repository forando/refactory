package shellexec

import (
	"github.com/forando/refactory/pkg/parser"
	"github.com/forando/refactory/pkg/schema"
	"github.com/pkg/errors"
	"os/exec"
)

func AwsGetCallerIdentity() (*schema.CallerIdentity, error) {
	cmd := exec.Command("aws", "sts", "get-caller-identity")
	output, err := cmd.CombinedOutput()
	if err != nil {
		if output != nil {
			return nil, errors.Errorf("Cannot login to AWS account: %s", string(output))
		}
		return nil, err
	}
	return parser.ParseCallerIdentity(output)
}
