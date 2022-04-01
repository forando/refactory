package schema

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/pkg/errors"
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
	ModuleName            string
	ProductTicket         string
	AccountName           string
	OrganizationalUnit    string
	CostCenter            int
	OwnerEmail            string
	OwnerJiraUsername     string
	GroupPermissions      map[string][]*PermissionSet
	UserPermissions       map[string][]*PermissionSet
	PersonalDataProcessed bool
	PersonalDataStored    bool
	RootEmail             string
	AccountBudget         int
	AccountBudgetEmail    string
}

type AccountModules []*AccountModule

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
			return errors.Errorf("unknown attribute [%s] in aws-account module", key)
		}
	}
	return nil
}
