package parser

import (
	"fmt"
	"github.com/forando/refactory/pkg/schema"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"strconv"
)

func ParseAccountModule(body *hclwrite.Body) (*schema.AccountModule, error) {
	var module schema.AccountModule

	attr := body.Attributes()

	var key string

	key = "name"
	accountNameAttr := attr[key]
	if accountNameAttr == nil {
		return makeError(key)
	}

	key = "organizational_unit"
	organizationalUnitAttr := attr[key]
	if organizationalUnitAttr == nil {
		return makeError(key)
	}

	key = "cost_center"
	costCenterAttr := attr[key]
	if costCenterAttr == nil {
		return makeError(key)
	}
	costCenterStr := getExpressionAsString(costCenterAttr)
	intVar, intErr := strconv.Atoi(costCenterStr)
	if intErr != nil {
		return nil, ParsingError{fmt.Sprintf("cannot parse [%s] costCenter value into int", costCenterStr)}
	}

	key = "komueb_product_ticket"
	productTicketAttr := attr[key]
	if productTicketAttr == nil {
		return makeError(key)
	}

	key = "owner_email"
	emailAttr := attr[key]
	if emailAttr == nil {
		return makeError(key)
	}

	key = "owner_jira_username"
	jiraUserNameAttr := attr[key]
	if jiraUserNameAttr == nil {
		return makeError(key)
	}

	module.AccountName = getExpressionAsString(accountNameAttr)
	module.OrganizationalUnit = getExpressionAsString(organizationalUnitAttr)
	module.CostCenter = intVar
	module.ProductTicket = getExpressionAsString(productTicketAttr)
	module.OwnerEmail = getExpressionAsString(emailAttr)
	module.OwnerJiraUsername = getExpressionAsString(jiraUserNameAttr)

	return &module, nil
}

func makeError(key string) (*schema.AccountModule, error) {
	return nil, ParsingError{fmt.Sprintf("[%s] property not found in aws-account module", key)}
}
