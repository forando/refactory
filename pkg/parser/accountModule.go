package parser

import (
	"github.com/forando/refactory/pkg/schema"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/pkg/errors"
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

func ParseAccountModule(body *hclwrite.Body, permissionSetNames *map[string]string) (*schema.AccountModule, error) {
	var module schema.AccountModule

	attrs := Attributes{Map: body.Attributes(), ModuleName: "aws-account"}

	if err := module.CheckAllAttributes(&attrs.Map); err != nil {
		return nil, err
	}

	if val, err := attrs.getStr(schema.AccName); err == nil {
		module.AccountName = val
	} else {
		return nil, err
	}

	if attr, err := attrs.getAttr(schema.AccOrganizationalUnit); err == nil {
		ouKey := string(attr.Expr().BuildTokens(nil)[2].Bytes)
		organizationalUnit, ok := organizationalUnits[ouKey]
		if !ok {
			return nil, errors.Errorf("cannot find [%s] key in organizationalUnits map", ouKey)
		}
		module.OrganizationalUnit = organizationalUnit
	} else {
		return nil, err
	}

	if val, keyNotFound, errInt := attrs.getInt(schema.AccCostCenter); keyNotFound == nil && errInt == nil {
		module.CostCenter = val
	} else if keyNotFound != nil {
		return nil, keyNotFound
	} else {
		return nil, errInt
	}

	if val, err := attrs.getStr(schema.AccKomuebProductTicket); err == nil {
		module.ProductTicket = val
	} else {
		return nil, err
	}

	if val, err := attrs.getStr(schema.AccOwnerEmail); err == nil {
		module.OwnerEmail = val
	} else {
		return nil, err
	}

	if val, err := attrs.getStr(schema.AccOwnerJiraUsername); err == nil {
		module.OwnerJiraUsername = val
	} else {
		return nil, err
	}

	if attr, err := attrs.getAttr(schema.AccGroupPermissions); err == nil {
		groupPermissions, err := parsePermissions(attr, schema.AccGroupPermissions, permissionSetNames)
		if err != nil {
			return nil, err
		}
		module.GroupPermissions = groupPermissions
	} else {
		return nil, err
	}

	if attr, err := attrs.getAttr(schema.AccUserPermissions); err == nil {
		userPermissions, err := parsePermissions(attr, schema.AccUserPermissions, permissionSetNames)
		if err != nil {
			return nil, err
		}
		module.UserPermissions = userPermissions
	}

	if val, keyNotFound, errBool := attrs.getBool(schema.AccPersonalDataProcessed); keyNotFound == nil && errBool == nil {
		module.PersonalDataProcessed = val
	} else if errBool != nil {
		return nil, errBool
	}

	if val, keyNotFound, errBool := attrs.getBool(schema.AccPersonalDataStored); keyNotFound == nil && errBool == nil {
		module.PersonalDataStored = val
	} else if errBool != nil {
		return nil, errBool
	}

	if val, err := attrs.getStr(schema.AccRootEmail); err == nil {
		module.RootEmail = val
	}

	if val, keyNotFound, errInt := attrs.getInt(schema.AccAccountBudget); keyNotFound == nil && errInt == nil {
		module.AccountBudget = val
	} else if errInt != nil {
		return nil, errInt
	}

	if val, err := attrs.getStr(schema.AccAccountBudgetEmail); err == nil {
		module.AccountBudgetEmail = val
	}

	return &module, nil
}

func parsePermissions(permissionAttr *hclwrite.Attribute, permissionsType string, permissionSetNames *map[string]string) (map[string][]*schema.PermissionSet, error) {

	expr := permissionAttr.Expr()
	exprTokens := expr.BuildTokens(nil)

	oBrace := exprTokens[0]
	cBrace := exprTokens[len(exprTokens)-1]
	if oBrace.Type != hclsyntax.TokenOBrace || cBrace.Type != hclsyntax.TokenCBrace {
		return nil, errors.Errorf("%s expression is not enclosed in braces", permissionsType)
	}
	inside := exprTokens[1 : len(exprTokens)-1]

	res := make(map[string][]*schema.PermissionSet)

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

				if _, found := res[key]; found {
					return nil, errors.Errorf("%s expression. Key [%s] already exists", permissionsType, key)
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
					return nil, errors.Errorf("%s expression. Key [%s]. Close bracket found before open one at token %v", permissionsType, key, index)
				}

				openBracketFound = false
				traversingKeyTokens = true

				value, err := getValueTokens(&inside, start, end, permissionSetNames)
				if err != nil {
					return nil, errors.Errorf("%s expression. Key [%s]. %s", permissionsType, key, err.Error())
				}
				res[key] = value
			}
		}
	}

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

func getValueTokens(valueTokens *hclwrite.Tokens, start int, end int, permissionSetNames *map[string]string) ([]*schema.PermissionSet, error) {
	values := make([]*schema.PermissionSet, 0)

	if end < start+1 {
		return values, nil
	}

	innerTokens := (*valueTokens)[start+1 : end]

	var tokens hclwrite.Tokens
	var value *schema.PermissionSet
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

func toValue(tokens hclwrite.Tokens, permissionSetNames *map[string]string) (*schema.PermissionSet, error) {

	for index, token := range tokens {
		if string(token.Bytes) == string('"') {
			return &schema.PermissionSet{Policy: schema.ManagedPolicy, Val: quotedTokensToString(tokens)}, nil
		}
		if string(token.Bytes) == "module" {
			if len(tokens) < index+3 {
				return nil, errors.Errorf("module reference does not have enough tokens, expected %v but actual %v", index+3, len(tokens))
			}
			moduleNameToken := tokens[index+2]
			moduleName := string(moduleNameToken.Bytes)

			permissionSetName, found := (*permissionSetNames)[moduleName]
			if !found {
				return nil, errors.Errorf("Cannot find [%s] permission set", moduleName)
			}
			return &schema.PermissionSet{Policy: schema.InlinePolicy, Val: permissionSetName}, nil
		}
	}

	return nil, errors.New("tokens don't contain neither quoted strings nor module references")
}
