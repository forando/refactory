package parser

import (
	"github.com/forando/refactory/pkg/schema"
	"github.com/hashicorp/hcl/v2/hclsimple"
)

func ParsePolicyDocumentBlock(bytes []byte) (*schema.PolicyDocument, error) {
	var document schema.PolicyDocument
	err := hclsimple.Decode(
		"dummy.hcl", bytes,
		nil, &document,
	)
	if err != nil {
		return nil, err
	}
	return &document, nil
}
