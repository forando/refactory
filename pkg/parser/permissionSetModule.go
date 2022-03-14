package parser

import (
	"fmt"
	"github.com/forando/refactory/pkg/schema"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func ParsePermissionSetModule(body *hclwrite.Body, policyDocuments *map[string]*schema.PolicyDocument) (*schema.PermissionSetModule, error) {
	var module schema.PermissionSetModule

	attr := body.Attributes()

	if err := module.CheckAllAttributes(&attr); err != nil {
		return nil, err
	}

	var key string

	key = schema.PsName
	nameAttr := attr[key]
	if nameAttr == nil {
		return makePermissionSetError(key)
	}
	module.PermissionSetName = getExpressionAsString(nameAttr)
	module.NameAttr = nameAttr

	key = "source"
	sourceAttr := attr[key]
	if sourceAttr == nil {
		return makePermissionSetError(key)
	}
	module.SourceAttr = sourceAttr

	key = schema.PsSsoAdminInstanceArn
	ssoAdminInstanceArnAttr := attr[key]
	if ssoAdminInstanceArnAttr == nil {
		return makePermissionSetError(key)
	}
	module.SsoAdminInstanceArnAttr = ssoAdminInstanceArnAttr

	key = schema.PsInlinePolicyDocument
	inlinePolicyDocumentsAttr := attr[key]
	if inlinePolicyDocumentsAttr != nil {

		pDocName := string(inlinePolicyDocumentsAttr.Expr().BuildTokens(nil)[4].Bytes)
		pDockBlock, found := (*policyDocuments)[pDocName]
		if !found {
			return nil, schema.ParsingError{Message: fmt.Sprintf("cannot find [%s] policyDocument", pDocName)}
		}
		module.InlinePolicyDocumentAttr = inlinePolicyDocumentsAttr
		module.PolicyDocument = pDockBlock
		module.PolicyDocumentName = pDocName
	}

	key = schema.PsManagedPolicyArns
	managedPolicyArnsAttr := attr[key]
	if managedPolicyArnsAttr != nil {
		module.ManagedPolicyArnsAttr = managedPolicyArnsAttr
	}

	key = schema.PsTags
	tagsAttr := attr[key]
	if tagsAttr != nil {
		module.TagsAttr = tagsAttr
	}

	return &module, nil
}

func makePermissionSetError(key string) (*schema.PermissionSetModule, error) {
	return nil, schema.ParsingError{Message: fmt.Sprintf("[%s] property not found in aws-ssoadmin-permission-set module", key)}
}
