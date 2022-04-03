package parser

import (
	"github.com/forando/refactory/pkg/schema"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/pkg/errors"
)

func ParsePermissionSetModule(body *hclwrite.Body, policyDocuments *map[string]*schema.PolicyDocument) (*schema.PermissionSetModule, error) {
	var module schema.PermissionSetModule

	attrs := Attributes{Map: body.Attributes(), ModuleName: "aws-ssoadmin-permission-set"}

	if err := module.CheckAllAttributes(&attrs.Map); err != nil {
		return nil, err
	}

	if attr, err := attrs.getAttr(schema.PsName); err == nil {
		module.PermissionSetName = getExpressionAsString(attr)
		module.NameAttr = attr
	} else {
		return nil, err
	}

	if attr, err := attrs.getAttr(schema.PsInlinePolicyDocument); err == nil {
		pDocName := string(attr.Expr().BuildTokens(nil)[4].Bytes)
		pDockBlock, found := (*policyDocuments)[pDocName]
		if !found {
			return nil, errors.Errorf("cannot find [%s] policyDocument", pDocName)
		}
		module.PolicyDocument = pDockBlock
		module.PolicyDocumentName = pDocName
	}

	if attr, err := attrs.getAttr(schema.PsManagedPolicyArns); err == nil {
		module.ManagedPolicyArnsAttr = attr
	}

	if attr, err := attrs.getAttr(schema.PsTags); err == nil {
		module.TagsAttr = attr
	}

	return &module, nil
}
