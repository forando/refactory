package factory

import (
	"github.com/forando/refactory/pkg/schema"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"log"
	"os"
)

func BootstrapPermissionSetTerragrunt(fileName string, module *schema.PermissionSetModule) {
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
	inputsBody.AppendUnstructuredTokens(module.SsoAdminInstanceArnAttr.BuildTokens(nil))

	if module.ManagedPolicyArnsAttr != nil {
		inputsBody.AppendUnstructuredTokens(module.ManagedPolicyArnsAttr.BuildTokens(nil))
	}

	if module.InlinePolicyDocumentsAttr != nil {
		inputsBody.AppendUnstructuredTokens(module.InlinePolicyDocumentsAttr.BuildTokens(nil))
	}

	if module.TagsAttr != nil {
		inputsBody.AppendUnstructuredTokens(module.TagsAttr.BuildTokens(nil))
	}

	rootBody.AppendNewline()

	rootBody.AppendBlock(module.PolicyDocument)

	if _, err := newFile.WriteTo(fw); err != nil {
		log.Fatal("Cannot write to the new file ", err)
	}
}
