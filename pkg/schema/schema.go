package schema

import (
	"fmt"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

const (
	AccName                  string = "name"
	AccOrganizationalUnit    string = "organizational_unit"
	AccCostCenter            string = "cost_center"
	AccKomuebProductTicket   string = "komueb_product_ticket"
	AccOwnerEmail            string = "owner_email"
	AccOwnerJiraUsername     string = "owner_jira_username"
	AccPersonalDataProcessed string = "personal_data_processed"
	AccPersonalDataStored    string = "personal_data_stored"
	AccGroupPermissions      string = "group_permissions"
	AccUserPermissions       string = "user_permissions"
	AccRootEmail             string = "root_email"
	AccAccountBudget         string = "account_budget"
	AccAccountBudgetEmail    string = "account_budget_email"
)

type AccountModule struct {
	ProductTicket         string
	AccountName           string
	OrganizationalUnit    string
	CostCenter            int
	OwnerEmail            string
	OwnerJiraUsername     string
	GroupPermissions      map[string][]string
	UserPermissions       map[string][]string
	PersonalDataProcessed bool
	PersonalDataStored    bool
	RootEmail             string
	AccountBudget         int
	AccountBudgetEmail    string
}

type AccountModules []*AccountModule

const (
	PsName                 string = "name"
	PsSsoAdminInstanceArn  string = "ssoadmin_instance_arn"
	PsManagedPolicyArns    string = "managed_policy_arns"
	PsInlinePolicyDocument string = "inline_policy_document"
	PsTags                 string = "tags"
)

type PermissionSetModule struct {
	SourceAttr                *hclwrite.Attribute
	SsoAdminInstanceArnAttr   *hclwrite.Attribute
	PermissionSetName         string
	NameAttr                  *hclwrite.Attribute
	InlinePolicyDocumentsAttr *hclwrite.Attribute
	PolicyDocumentName        string
	PolicyDocument            *hclwrite.Block
	ManagedPolicyArnsAttr     *hclwrite.Attribute
	TagsAttr                  *hclwrite.Attribute
	ProductTicket             string
}

type PermissionSetModules []*PermissionSetModule

type BlockMetaData struct {
	BlockType string
	BlockName string
}

const (
	AccountModuleType       string = "aws-account"
	PermissionSetModuleType string = "aws-ssoadmin-permission-set"
	IamPolicyDocumentType   string = "aws_iam_policy_document"
)

const (
	ModuleBlock string = "module"
	DataBlock   string = "data"
)

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
			return ParsingError{Message: fmt.Sprintf("unknown attribute [%s] in aws-ssoadmin-permission-set module", key)}
		}
	}

	return nil
}

var accountModuleKnownAttrs = map[string]bool{
	"source":                  true,
	"name":                    true,
	"organizational_unit":     true,
	"cost_center":             true,
	"komueb_product_ticket":   true,
	"owner_email":             true,
	"owner_jira_username":     true,
	"group_permissions":       true,
	"user_permissions":        true,
	"personal_data_processed": true,
	"personal_data_stored":    true,
	"root_email":              true,
	"account_budget":          true,
	"account_budget_email":    true,
	"depends_on":              true,
}

func (module *AccountModule) CheckAllAttributes(attrs *map[string]*hclwrite.Attribute) error {

	for key := range *attrs {
		if _, ok := accountModuleKnownAttrs[key]; !ok {
			return ParsingError{Message: fmt.Sprintf("unknown attribute [%s] in aws-account module", key)}
		}
	}

	return nil
}
