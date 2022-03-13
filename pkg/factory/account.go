package factory

import (
	"fmt"
	"github.com/forando/refactory/pkg/schema"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	"log"
	"os"
)

func BootstrapAccountTerragrunt(fileName string, module *schema.AccountModule) {
	fw, osErr := os.Create(fileName)

	if osErr != nil {
		log.Fatal("Cannot create new file ", osErr)
	}

	newFile := hclwrite.NewEmptyFile()
	rootBody := newFile.Body()

	includeBlock := rootBody.AppendNewBlock("include", []string{"root"})
	includeBody := includeBlock.Body()

	tokens := hclwrite.Tokens{
		{Type: hclsyntax.TokenIdent, Bytes: []byte(`find_in_parent_folders()`), SpacesBefore: 1},
	}

	includeBody.SetAttributeRaw("path", tokens)

	rootBody.AppendNewline()

	inputsBody := rootBody.AppendNewBlock("inputs =", nil).Body()

	inputsBody.SetAttributeValue(schema.AccName, cty.StringVal(module.AccountName))
	inputsBody.SetAttributeValue(schema.AccOrganizationalUnit, cty.StringVal(module.OrganizationalUnit))
	inputsBody.SetAttributeValue(schema.AccCostCenter, cty.NumberIntVal(int64(module.CostCenter)))
	inputsBody.SetAttributeValue(schema.AccKomuebProductTicket, cty.StringVal(module.ProductTicket))
	inputsBody.AppendNewline()
	inputsBody.SetAttributeValue(schema.AccOwnerEmail, cty.StringVal(module.OwnerEmail))
	inputsBody.SetAttributeValue(schema.AccOwnerJiraUsername, cty.StringVal(module.OwnerJiraUsername))

	if len(module.GroupPermissions) > 0 {
		inputsBody.AppendNewline()
		groupPermissionsBody := inputsBody.AppendNewBlock(fmt.Sprintf("%s =", schema.AccGroupPermissions), nil).Body()
		for key, vals := range module.GroupPermissions {
			ctyVals := make([]cty.Value, 0)
			for _, val := range vals {
				ctyVals = append(ctyVals, cty.StringVal(val))
			}
			key = fmt.Sprintf("%q", key)
			groupPermissionsBody.SetAttributeValue(key, cty.ListVal(ctyVals))
		}
	}

	if len(module.UserPermissions) > 0 {
		inputsBody.AppendNewline()
		userPermissionsBody := inputsBody.AppendNewBlock(fmt.Sprintf("%s =", schema.AccUserPermissions), nil).Body()
		for key, vals := range module.UserPermissions {
			ctyVals := make([]cty.Value, 0)
			for _, val := range vals {
				ctyVals = append(ctyVals, cty.StringVal(val))
			}
			key = fmt.Sprintf("%q", key)
			userPermissionsBody.SetAttributeValue(key, cty.ListVal(ctyVals))
		}
	}

	if module.PersonalDataProcessed {
		inputsBody.SetAttributeValue(schema.AccPersonalDataProcessed, cty.BoolVal(true))
	}

	if module.PersonalDataStored {
		inputsBody.SetAttributeValue(schema.AccPersonalDataStored, cty.BoolVal(true))
	}

	if len(module.RootEmail) > 0 {
		inputsBody.SetAttributeValue(schema.AccRootEmail, cty.StringVal(module.RootEmail))
	}

	if module.AccountBudget > 0 {
		inputsBody.SetAttributeValue(schema.AccAccountBudget, cty.NumberIntVal(int64(module.AccountBudget)))
	}

	if len(module.AccountBudgetEmail) > 0 {
		inputsBody.SetAttributeValue(schema.AccAccountBudgetEmail, cty.StringVal(module.AccountBudgetEmail))
	}

	if _, err := newFile.WriteTo(fw); err != nil {
		log.Fatal("Cannot write to the new file ", err)
	}
}
