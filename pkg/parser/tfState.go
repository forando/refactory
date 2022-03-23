package parser

import (
	"fmt"
	"github.com/forando/refactory/pkg/schema"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"strings"
)

func parseTfStateFile(file string) (*schema.TerraformState, error) {
	var state schema.TerraformState
	err := hclsimple.DecodeFile(file, nil, &state)
	if err != nil {
		return nil, err
	}
	return &state, nil
}

func parseModuleName(module string) (string, error) {
	split := strings.Split(module, ".")
	if len(split) != 2 || len(split[1]) == 0 {
		return "", schema.ParsingError{Message: fmt.Sprintf("cannot parse module name: %s", module)}
	}
	if split[0] != "module" {
		return "", schema.ParsingError{Message: fmt.Sprintf("bad module name: %s, does not prefixed with 'module' keyword", module)}
	}
	return split[1], nil
}
