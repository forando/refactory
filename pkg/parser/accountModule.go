package parser

import (
	"github.com/forando/refactory/pkg/schema"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/pkg/errors"
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

	inlinePolicies := make(map[string]interface{})
	for _, val := range *permissionSetNames {
		inlinePolicies[val] = nil
	}

	ctx := buildContext(permissionSetNames)

	exprParser := ExpressionParser{Context: ctx}

	var err error
	var rawPermissions *map[string][]byte
	bytes := permissionAttr.Expr().BuildTokens(nil).Bytes()
	if rawPermissions, err = exprParser.ParseObjectExpr(permissionsType, bytes); err != nil {
		return nil, err
	}

	res := make(map[string][]*schema.PermissionSet)

	for key, bytes := range *rawPermissions {
		var permissions []string
		permissions, err = exprParser.ParseArrayExpr(key, bytes)
		if err != nil {
			return nil, err
		}
		permissionSets := make([]*schema.PermissionSet, 0)
		for _, permission := range permissions {
			policyType := schema.ManagedPolicy
			if _, found := inlinePolicies[permission]; found {
				policyType = schema.InlinePolicy
			}
			permissionSets = append(permissionSets, &schema.PermissionSet{Policy: policyType, Val: permission})
		}
		res[key] = permissionSets
	}

	return res, nil
}
