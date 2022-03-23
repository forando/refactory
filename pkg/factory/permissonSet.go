package factory

import (
	"github.com/forando/refactory/pkg/schema"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	"log"
	"os"
)

const ssoAdminInstanceArnProd = "arn:aws:sso:::instance/ssoins-69873834cda94459"
const ssoAdminInstanceArnTest = "arn:aws:sso:::instance/ssoins-69877a64b4162b02"

func BootstrapPermissionSetTerragrunt(fileName string, module *schema.PermissionSetModule, org schema.Org) {
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

	inputsBody.AppendUnstructuredTokens(module.NameAttr.BuildTokens(nil))

	if org == schema.ProdOrg {
		inputsBody.SetAttributeValue(schema.PsSsoAdminInstanceArn, cty.StringVal(ssoAdminInstanceArnProd))
	} else {
		inputsBody.SetAttributeValue(schema.PsSsoAdminInstanceArn, cty.StringVal(ssoAdminInstanceArnTest))
	}

	if module.ManagedPolicyArnsAttr != nil {
		inputsBody.AppendUnstructuredTokens(module.ManagedPolicyArnsAttr.BuildTokens(nil))
	}

	if module.InlinePolicyDocumentAttr != nil {
		inputsBody.AppendUnstructuredTokens(hclwrite.Tokens{
			{Type: hclsyntax.TokenStringLit, Bytes: []byte("# language=JSON")},
			{Type: hclsyntax.TokenNewline, Bytes: []byte("\n")},
		})
		inputsBody.SetAttributeRaw(schema.PsInlinePolicyDocument, *buildInlinePolicyTokens(module.PolicyDocument))
	}

	if module.TagsAttr != nil {
		inputsBody.AppendUnstructuredTokens(module.TagsAttr.BuildTokens(nil))
	}

	if _, err := newFile.WriteTo(fw); err != nil {
		log.Fatal("Cannot write to the new file ", err)
	}
}
