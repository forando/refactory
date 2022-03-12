package parser

import (
	"fmt"
	"github.com/forando/refactory/pkg/schema"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func ParsePermissionSetModule(body *hclwrite.Body) (*schema.PermissionSetModule, error) {
	var module schema.PermissionSetModule

	attr := body.Attributes()

	if err := module.CheckAllAttributes(&attr); err != nil {
		return nil, err
	}

	var key string

	key = "name"
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

	key = "ssoadmin_instance_arn"
	ssoAdminInstanceArnAttr := attr[key]
	if ssoAdminInstanceArnAttr == nil {
		return makePermissionSetError(key)
	}
	module.SsoAdminInstanceArnAttr = ssoAdminInstanceArnAttr

	key = "inline_policy_documents"
	inlinePolicyDocumentsAttr := attr[key]
	if inlinePolicyDocumentsAttr != nil {
		module.InlinePolicyDocumentsAttr = inlinePolicyDocumentsAttr
	}

	key = "managed_policy_arns"
	managedPolicyArnsAttr := attr[key]
	if managedPolicyArnsAttr != nil {
		module.ManagedPolicyArnsAttr = managedPolicyArnsAttr
	}

	key = "tags"
	tagsAttr := attr[key]
	if tagsAttr != nil {
		module.TagsAttr = tagsAttr
	}

	return &module, nil
}

func makePermissionSetError(key string) (*schema.PermissionSetModule, error) {
	return nil, schema.ParsingError{Message: fmt.Sprintf("[%s] property not found in aws-ssoadmin-permission-set module", key)}
}
