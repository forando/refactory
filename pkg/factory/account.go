package factory

import (
	"fmt"
	"github.com/forando/refactory/pkg/schema"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	"log"
	"os"
	"path"
)

func BootstrapAccountTerragrunt(fileName string, module *schema.AccountModule) {
	fw, osErr := os.Create(fileName)

	if osErr != nil {
		log.Fatal("Cannot create new file ", osErr)
	}

	newFile := hclwrite.NewEmptyFile()
	rootBody := newFile.Body()

	buildIncludeRoot(rootBody)

	buildDependencies(rootBody, module)

	buildInputs(rootBody, module)

	if _, err := newFile.WriteTo(fw); err != nil {
		log.Fatal("Cannot write to the new file ", err)
	}
}

func buildIncludeRoot(rootBody *hclwrite.Body) {
	includeBlock := rootBody.AppendNewBlock("include", []string{"root"})
	includeBody := includeBlock.Body()

	tokens := hclwrite.Tokens{
		{Type: hclsyntax.TokenIdent, Bytes: []byte(`find_in_parent_folders()`), SpacesBefore: 1},
	}

	includeBody.SetAttributeRaw("path", tokens)

	rootBody.AppendNewline()
}

func buildDependencies(rootBody *hclwrite.Body, module *schema.AccountModule) {
	permissionSets := make(map[string]bool)
	for _, val := range module.GroupPermissions {
		for _, perm := range val {
			if perm.Policy == schema.InlinePolicy {
				permissionSets[perm.Val] = true
			}
		}
	}
	for _, val := range module.UserPermissions {
		for _, perm := range val {
			if perm.Policy == schema.InlinePolicy {
				permissionSets[perm.Val] = true
			}
		}
	}
	if len(permissionSets) > 0 {
		dependenciesBody := rootBody.AppendNewBlock("dependencies", nil).Body()
		perms := make([]cty.Value, 0)
		for perm, _ := range permissionSets {
			pathToModule := path.Join("..", "PermissionSets", perm)
			perms = append(perms, cty.StringVal(pathToModule))
		}
		dependenciesBody.SetAttributeValue("paths", cty.ListVal(perms))

		rootBody.AppendNewline()
	}
}

func buildInputs(rootBody *hclwrite.Body, module *schema.AccountModule) {

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
			for _, perm := range vals {
				ctyVals = append(ctyVals, cty.StringVal(perm.Val))
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
			for _, perm := range vals {
				ctyVals = append(ctyVals, cty.StringVal(perm.Val))
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
}
