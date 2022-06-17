package parser

import (
	"github.com/forando/refactory/pkg/schema"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/pkg/errors"
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

func parseTfStateBytes(bytes *[]byte) (*schema.TerraformState, error) {
	var state schema.TerraformState
	err := hclsimple.Decode("dummy.json", *bytes, nil, &state)
	if err != nil {
		return nil, err
	}
	return &state, nil
}

func parseModuleName(module string) (string, error) {
	split := strings.Split(module, ".")
	if len(split) != 2 || len(split[1]) == 0 {
		return "", errors.Errorf("cannot parse module name: %s", module)
	}
	if split[0] != "module" {
		return "", errors.Errorf("bad module name: %s, does not prefixed with 'module' keyword", module)
	}
	return split[1], nil
}
