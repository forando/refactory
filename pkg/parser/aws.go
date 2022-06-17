package parser

import (
	"github.com/forando/refactory/pkg/schema"
	"github.com/hashicorp/hcl/v2/hclsimple"
)

func ParseCallerIdentity(bytes []byte) (*schema.CallerIdentity, error) {
	var parsed schema.CallerIdentity
	err := hclsimple.Decode("dummy.json", bytes, nil, &parsed)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}
