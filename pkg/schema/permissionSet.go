package schema

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/pkg/errors"
)

type PermissionSet struct {
	Policy PolicyType
	Val    string
}

const (
	PsName                 string = "name"
	PsSsoAdminInstanceArn  string = "ssoadmin_instance_arn"
	PsManagedPolicyArns    string = "managed_policy_arns"
	PsInlinePolicyDocument string = "inline_policy_document"
	PsTags                 string = "tags"
)

type PermissionSetModule struct {
	ModuleName            string
	PermissionSetName     string
	NameAttr              *hclwrite.Attribute
	PolicyDocumentName    string
	PolicyDocument        *PolicyDocument
	ManagedPolicyArnsAttr *hclwrite.Attribute
	TagsAttr              *hclwrite.Attribute
	ProductTicket         string
}

type PermissionSetModules []*PermissionSetModule

var permissionSetKnownAttrs = map[string]bool{
	"source":                 true,
	"name":                   true,
	"ssoadmin_instance_arn":  true,
	"inline_policy_document": true,
	"managed_policy_arns":    true,
	"tags":                   true,
}

func (module *PermissionSetModule) CheckAllAttributes(attrs *map[string]*hclwrite.Attribute) error {

	for key := range *attrs {
		if _, ok := permissionSetKnownAttrs[key]; !ok {
			return errors.Errorf("unknown attribute [%s] in aws-ssoadmin-permission-set module", key)
		}
	}

	return nil
}
