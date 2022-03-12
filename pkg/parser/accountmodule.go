package parser

import (
	"fmt"
	"github.com/forando/refactory/pkg/schema"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"strconv"
	"strings"
)

var organizationalUnits = map[string]string{
	"organizational_unit_workloads_prod":      "Workloads_Prod",
	"organizational_unit_workloads_sdlc":      "Workloads_Sdlc",
	"organizational_unit_infrastructure_prod": "Infrastructure_Prod",
	"organizational_unit_infrastructure_sdlc": "Infrastructure_Sdlc",
	"organizational_unit_sandbox":             "Sandbox",
	"organizational_unit_security":            "Security",
	"organizational_unit_transitional":        "Transitional",
	"organizational_unit_suspended":           "Suspended",
}

func ParseAccountModule(body *hclwrite.Body, permissionSetNames map[string]string) (*schema.AccountModule, error) {
	var module schema.AccountModule

	attr := body.Attributes()

	if err := module.CheckAllAttributes(&attr); err != nil {
		return nil, err
	}

	var key string

	key = "name"
	accountNameAttr := attr[key]
	if accountNameAttr == nil {
		return makeError(key)
	}
	module.AccountName = getExpressionAsString(accountNameAttr)

	key = "organizational_unit"
	organizationalUnitAttr := attr[key]
	if organizationalUnitAttr == nil {
		return makeError(key)
	}
	ouKey := string(organizationalUnitAttr.Expr().BuildTokens(nil)[2].Bytes)
	organizationalUnit, ok := organizationalUnits[ouKey]
	if !ok {
		return nil, schema.ParsingError{fmt.Sprintf("cannot find [%s] key in organizationalUnits map", ouKey)}
	}
	module.OrganizationalUnit = organizationalUnit

	key = "cost_center"
	costCenterAttr := attr[key]
	if costCenterAttr == nil {
		return makeError(key)
	}
	costCenterStr := getExpressionAsString(costCenterAttr)
	intVar, intErr := strconv.Atoi(costCenterStr)
	if intErr != nil {
		return nil, schema.ParsingError{Message: fmt.Sprintf("cannot parse %s vlaue of [%s] into int", key, costCenterStr)}
	}
	module.CostCenter = intVar

	key = "komueb_product_ticket"
	productTicketAttr := attr[key]
	if productTicketAttr == nil {
		return makeError(key)
	}
	module.ProductTicket = getExpressionAsString(productTicketAttr)

	key = "owner_email"
	emailAttr := attr[key]
	if emailAttr == nil {
		return makeError(key)
	}
	module.OwnerEmail = getExpressionAsString(emailAttr)

	key = "owner_jira_username"
	jiraUserNameAttr := attr[key]
	if jiraUserNameAttr == nil {
		return makeError(key)
	}
	module.OwnerJiraUsername = getExpressionAsString(jiraUserNameAttr)

	key = "group_permissions"
	groupPermissionsAttr := attr[key]
	if groupPermissionsAttr == nil {
		return makeError(key)
	}
	groupPermissions, err := parsePermissions(groupPermissionsAttr, key, permissionSetNames)
	if err != nil {
		return nil, err
	}
	module.GroupPermissions = groupPermissions

	key = "user_permissions"
	userPermissionsAttr := attr[key]
	if userPermissionsAttr != nil {
		userPermissions, err := parsePermissions(userPermissionsAttr, key, permissionSetNames)
		if err != nil {
			return nil, err
		}
		module.UserPermissions = userPermissions
	}

	key = "personal_data_processed"
	personalDataProcessedAttr := attr[key]
	if personalDataProcessedAttr != nil {
		personalDataProcessedStr := getExpressionAsString(personalDataProcessedAttr)
		boolVar, boolErr := strconv.ParseBool(personalDataProcessedStr)
		if boolErr != nil {
			return nil, schema.ParsingError{Message: fmt.Sprintf("cannot parse %s vlaue of [%s] into bool", key, personalDataProcessedStr)}
		}
		module.PersonalDataProcessed = boolVar
	}

	key = "personal_data_stored"
	personalDataStoredAttr := attr[key]
	if personalDataStoredAttr != nil {
		personalDataStoredStr := getExpressionAsString(personalDataStoredAttr)
		boolVar, boolErr := strconv.ParseBool(personalDataStoredStr)
		if boolErr != nil {
			return nil, schema.ParsingError{Message: fmt.Sprintf("cannot parse %s vlaue of [%s] into bool", key, personalDataStoredStr)}
		}
		module.PersonalDataStored = boolVar
	}

	key = "root_email"
	rootEmailAttr := attr[key]
	if rootEmailAttr != nil {
		module.RootEmail = getExpressionAsString(rootEmailAttr)
	}

	key = "account_budget"
	accountBudgetAttr := attr[key]
	if accountBudgetAttr != nil {
		accountBudgetStr := getExpressionAsString(accountBudgetAttr)
		intVar, intErr := strconv.Atoi(accountBudgetStr)
		if intErr != nil {
			return nil, schema.ParsingError{Message: fmt.Sprintf("cannot parse %s vlaue of [%s] into int", key, accountBudgetStr)}
		}
		module.AccountBudget = intVar
	}

	key = "account_budget_email"
	accountBudgetEmailAttr := attr[key]
	if accountBudgetEmailAttr != nil {
		module.AccountBudgetEmail = getExpressionAsString(accountBudgetEmailAttr)
	}

	return &module, nil
}

func parsePermissions(permissionAttr *hclwrite.Attribute, permissionsType string, permissionSetNames map[string]string) (map[string][]string, error) {

	expr := permissionAttr.Expr()
	/*tokens := hclwrite.Tokens{
		{Type: hclsyntax.TokenIdent, Bytes: []byte(`dummy`)},
	}*/
	exprTokens := expr.BuildTokens(nil)

	oBrace := exprTokens[0]
	cBrace := exprTokens[len(exprTokens)-1]
	if oBrace.Type != hclsyntax.TokenOBrace || cBrace.Type != hclsyntax.TokenCBrace {
		return nil, schema.ParsingError{Message: fmt.Sprintf("%s expression is not enclosed in braces", permissionsType)}
	}
	inside := exprTokens[1 : len(exprTokens)-1]

	res := make(map[string][]string)

	var tokens hclwrite.Tokens
	var key string
	traversingKeyTokens := true
	openBracketFound := false
	start := 0
	end := 0
	for index, token := range inside {
		if traversingKeyTokens {
			keyProcessed := processKeyToken(token, &tokens)
			if keyProcessed {
				key = quotedTokensToString(tokens)

				if _, ok := res[key]; ok {
					return nil, schema.ParsingError{Message: fmt.Sprintf("%s expression. Key [%s] already exists", permissionsType, key)}
				}

				tokens = hclwrite.Tokens{}
				traversingKeyTokens = false
			}
		} else {

			if token.Type == hclsyntax.TokenOBrack {
				start = index
				openBracketFound = true
				continue
			}

			if token.Type == hclsyntax.TokenCBrack {
				end = index

				if !openBracketFound {
					return nil, schema.ParsingError{Message: fmt.Sprintf("%s expression. Key [%s]. Close bracket found before open one at token %v", permissionsType, key, index)}
				}

				openBracketFound = false
				traversingKeyTokens = true

				value, err := getValueTokens(&inside, start, end, permissionSetNames)
				if err != nil {
					return nil, schema.ParsingError{Message: fmt.Sprintf("%s expression. Key [%s]. %s", permissionsType, key, err.Error())}
				}
				res[key] = value
			}
		}
	}

	/*file, parsedDiags := hclwrite.ParseConfig(tokens.Bytes(), "", hcl.InitialPos)

	if parsedDiags.HasErrors() {
		log.Fatal("Cannot parse bytes from groupPermission block ", parsedDiags.Error())
	}

	body := file.Body()

	for _, val := range body.Blocks() {
		log.Printf("Block %s", string(val.Body().BuildTokens(nil).Bytes()))
	}*/

	return res, nil
}

func processKeyToken(currentToken *hclwrite.Token, keyTokens *hclwrite.Tokens) bool {
	if currentToken.Type == hclsyntax.TokenEqual {
		return true
	} else {
		*keyTokens = hclwrite.Tokens{
			currentToken,
		}.BuildTokens(*keyTokens)
		return false
	}
}

func getValueTokens(valueTokens *hclwrite.Tokens, start int, end int, permissionSetNames map[string]string) ([]string, error) {
	values := make([]string, 0)

	if end < start+1 {
		return values, nil
	}

	innerTokens := (*valueTokens)[start+1 : end]

	var tokens hclwrite.Tokens
	var value string
	var err error

	valueProcessed := false

	for _, token := range innerTokens {
		valueProcessed = processValueToken(token, &tokens)
		if valueProcessed {
			value, err = toValue(tokens, permissionSetNames)
			if err != nil {
				return nil, err
			}
			values = append(values, value)
			tokens = hclwrite.Tokens{}
		}
	}

	//if only one value in the array
	if !valueProcessed {
		value, err = toValue(tokens, permissionSetNames)
		if err != nil {
			//if there is a trailing comma and a new line char
			if len(tokens) <= 1 && tokens[0].Type == hclsyntax.TokenNewline {
				return values, nil
			}
			return nil, err
		}
		values = append(values, value)
	}

	return values, nil
}

func processValueToken(currentToken *hclwrite.Token, keyTokens *hclwrite.Tokens) bool {
	if currentToken.Type == hclsyntax.TokenComma {
		return true
	} else {
		*keyTokens = hclwrite.Tokens{
			currentToken,
		}.BuildTokens(*keyTokens)
		return false
	}
}

func quotedTokensToString(tokens hclwrite.Tokens) string {
	var ts string
	if len(tokens) < 3 {
		return string(tokens.Bytes())
	}
	startQuoteFound := false
	start := 0
	end := 0
	for index, token := range tokens {
		if string(token.Bytes) == string('"') {
			if startQuoteFound {
				end = index
				break
			} else {
				start = index
				startQuoteFound = true
			}
		}
	}
	if end < (start + 1) {
		return string(tokens.Bytes())
	}
	ts = string(tokens[start+1 : end].Bytes())
	return strings.TrimSpace(ts)
}

func makeError(key string) (*schema.AccountModule, error) {
	return nil, schema.ParsingError{Message: fmt.Sprintf("[%s] property not found in aws-account module", key)}
}

/*type PermissionValueTokensType int

const (
	Literal         PermissionValueTokensType = 0
	ModuleReference PermissionValueTokensType = 1
)

type PermissionValueTokens struct {
	ValueType PermissionValueTokensType
	Value     string
}*/

func toValue(tokens hclwrite.Tokens, permissionSetNames map[string]string) (string, error) {

	for index, token := range tokens {
		if string(token.Bytes) == string('"') {
			return quotedTokensToString(tokens), nil
		}
		if string(token.Bytes) == "module" {
			if len(tokens) < index+3 {
				return "", schema.ParsingError{Message: fmt.Sprintf("module reference does not have enough tokens, expected %v but actual %v", index+3, len(tokens))}
			}
			moduleNameToken := tokens[index+2]
			moduleName := string(moduleNameToken.Bytes)

			permissionSetName, ok := permissionSetNames[moduleName]
			if !ok {
				return "", schema.ParsingError{Message: fmt.Sprintf("Cannot find [%s] permission set", moduleName)}
			}
			return permissionSetName, nil
		}
	}

	return "", schema.ParsingError{Message: fmt.Sprintf("tokens don't contain neither quoted strings nor module references")}
}
