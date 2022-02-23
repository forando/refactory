package factory

import (
	"fmt"
	"github.com/forando/refactory/pkg/filesystem"
	"github.com/forando/refactory/pkg/schema"
	_ "github.com/forando/refactory/pkg/schema"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	"log"
	"os"
	"strings"
)

// SaveToNewFile Saves Body Content to a new file
// and prints out the parsed content to the console
func SaveToNewFile(fileName string, body *hclwrite.Body) {
	fw, osErr := os.Create(fileName)

	if osErr != nil {
		log.Fatal("Cannot create new file ", osErr)
	}

	newFile := hclwrite.NewEmptyFile()
	rootBody := newFile.Body()

	for _, block := range body.Blocks() {
		log.Println(block.Type())
		for _, label := range block.Labels() {
			log.Println(label)
		}
		for key, attr := range block.Body().Attributes() {
			tokens := attr.Expr().BuildTokens(nil)

			log.Printf("attr[%s]=%s", key, removeQuotes(string(tokens.Bytes())))
		}
		rootBody.AppendBlock(block)
	}

	_, writeErr := newFile.WriteTo(fw)

	if writeErr != nil {
		log.Fatal("Cannot write to the new file ", writeErr)
	}
}

func removeQuotes(s string) string {
	out := strings.TrimSpace(s)
	if len(out) < 3 {
		return s
	}
	if out[len(out)-1] == '"' {
		out = out[:len(out)-1]
	}
	if out[0] == '"' {
		out = out[1:]
	}
	return out
}

/*block := hclwrite.NewBlock("foo", []string{})

blockBody := block.Body()

blockBody.SetAttributeValue("hello", cty.StringVal("world"))

var toks hclwrite.Tokens

rootBody.SetAttributeRaw("test", toks)

rootBody.AppendBlock(block)

block.BuildTokens(toks)

log.Println(string(toks.Bytes()))*/

/*rootBody.SetAttributeTraversal("test", hcl.Traversal{
	hcl.TraverseRoot{Name: "var"},
	hcl.TraverseAttr{Name: "name"},
})*/

func BootstrapAccountTerragrunt(fileName string) {
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

	inputsBody.SetAttributeValue("name", cty.StringVal("Test_Account_v5"))
	inputsBody.SetAttributeValue("organizational_unit", cty.StringVal("Sandbox"))
	inputsBody.SetAttributeValue("cost_center", cty.NumberIntVal(1187))
	inputsBody.SetAttributeValue("komueb_product_ticket", cty.StringVal("KOMUEB-3188"))
	inputsBody.AppendNewline()
	inputsBody.SetAttributeValue("owner_email", cty.StringVal("stefan.neben@idealo.de"))
	inputsBody.SetAttributeValue("owner_jira_username", cty.StringVal("m.brettschneider"))

	inputsBody.AppendNewline()
	groupPermissions := make(map[string]cty.Value)
	groupPermissions["Cloud Shuttle"] = cty.ListVal([]cty.Value{
		cty.StringVal("AWSAdministratorAccess"), cty.StringVal("AWSReadOnlyAccess"),
	})
	inputsBody.SetAttributeValue("group_permissions", cty.ObjectVal(groupPermissions))

	inputsBody.AppendNewline()
	userPermissions := make(map[string]cty.Value)
	userPermissions["stefan.neben@idealo.de"] = cty.ListVal([]cty.Value{
		cty.StringVal("AWSPowerUserAccess"),
	})
	inputsBody.SetAttributeValue("user_permissions", cty.ObjectVal(userPermissions))

	_, writeErr := newFile.WriteTo(fw)

	if writeErr != nil {
		log.Fatal("Cannot write to the new file ", writeErr)
	}
}

func Bootstrap(modules *schema.AccountModules) {
	fs := filesystem.NewOsFs()
	for _, module := range *modules {
		fs.MakeDirs(fmt.Sprintf("%s/%s", module.ProductTicket, module.AccountName))
	}
}
