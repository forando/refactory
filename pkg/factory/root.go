package factory

import (
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	"log"
	"os"
)

func BootstrapRootTerragrunt(filePath string) {
	fw, osErr := os.Create(filePath)

	if osErr != nil {
		log.Fatal("Cannot create new file ", osErr)
	}

	newFile := hclwrite.NewEmptyFile()
	rootBody := newFile.Body()

	buildRemoteState(rootBody)
	buildTerraformBlock(rootBody)
	buildGenerateBlock(rootBody)

	if _, err := newFile.WriteTo(fw); err != nil {
		log.Fatal("Cannot write to the new file ", err)
	}
}

func buildRemoteState(body *hclwrite.Body) {
	remoteStateBlock := body.AppendNewBlock("remote_state", nil)
	remoteStateBody := remoteStateBlock.Body()

	remoteStateBody.SetAttributeValue("backend", cty.StringVal("s3"))
	remoteStateBody.SetAttributeValue("backend", cty.StringVal("s3"))
	generateBody := remoteStateBody.AppendNewBlock("generate =", nil).Body()
	generateBody.SetAttributeValue("path", cty.StringVal("backend.tf"))
	generateBody.SetAttributeValue("if_exists", cty.StringVal("overwrite"))
	configBody := remoteStateBody.AppendNewBlock("config =", nil).Body()
	configBody.SetAttributeValue("bucket", cty.StringVal("idealo-test-org-tg-state"))
	configBody.SetAttributeRaw("key", hclwrite.Tokens{
		{Type: hclsyntax.TokenOQuote, Bytes: []byte{'"'}},
		{Type: hclsyntax.TokenStringLit, Bytes: []byte("${path_relative_to_include()}/terraform.tfstate")},
		{Type: hclsyntax.TokenCQuote, Bytes: []byte{'"'}},
	})
	//configBody.SetAttributeValue("key", cty.StringVal("${path_relative_to_include()}/terraform.tfstate"))
	configBody.SetAttributeValue("region", cty.StringVal("eu-central-1"))
	configBody.SetAttributeValue("encrypt", cty.BoolVal(true))
	configBody.SetAttributeValue("dynamodb_table", cty.StringVal("idealo-test-org-tg-lock"))
}

func buildTerraformBlock(body *hclwrite.Body) {
	body.AppendNewline()

	terraformBlock := body.AppendNewBlock("terraform", nil)
	terraformBody := terraformBlock.Body()
	terraformBody.SetAttributeRaw("source", hclwrite.Tokens{
		{Type: hclsyntax.TokenStringLit, Bytes: []byte(" regexall(")},
		{Type: hclsyntax.TokenOQuote, Bytes: []byte{'"'}},
		{Type: hclsyntax.TokenStringLit, Bytes: []byte(".*/PermissionSets/.*")},
		{Type: hclsyntax.TokenCQuote, Bytes: []byte{'"'}},
		{Type: hclsyntax.TokenComma, Bytes: []byte{','}},
		{Type: hclsyntax.TokenStringLit, Bytes: []byte("get_original_terragrunt_dir()) ? ")},
		{Type: hclsyntax.TokenOQuote, Bytes: []byte{'"'}},
		{Type: hclsyntax.TokenStringLit, Bytes: []byte("${get_path_to_repo_root()}//modules/aws-ssoadmin-permission-set")},
		{Type: hclsyntax.TokenCQuote, Bytes: []byte{'"'}},
		{Type: hclsyntax.TokenColon, Bytes: []byte{':'}},
		{Type: hclsyntax.TokenOQuote, Bytes: []byte{'"'}},
		{Type: hclsyntax.TokenStringLit, Bytes: []byte("${get_path_to_repo_root()}//modules/aws-account")},
		{Type: hclsyntax.TokenCQuote, Bytes: []byte{'"'}},
	})
}

func buildGenerateBlock(body *hclwrite.Body) {
	body.AppendNewline()

	generateBlock := body.AppendNewBlock("generate", []string{"provider"})
	generateBody := generateBlock.Body()
	generateBody.SetAttributeValue("path", cty.StringVal("provider.tf"))
	generateBody.SetAttributeValue("if_exists", cty.StringVal("overwrite_terragrunt"))
	generateBody.AppendNewline()

	providerAwsBlock := hclwrite.NewBlock("provider", []string{"aws"})
	providerAwsBody := providerAwsBlock.Body()
	providerAwsBody.SetAttributeValue("region", cty.StringVal("eu-central-1"))
	providerAwsBody.SetAttributeValue("allowed_account_ids", cty.ListVal([]cty.Value{cty.StringVal("573275350257")}))

	providerControlTowerBlock := hclwrite.NewBlock("provider", []string{"controltower"})
	providerControlTowerBody := providerControlTowerBlock.Body()
	providerControlTowerBody.SetAttributeValue("region", cty.StringVal("eu-central-1"))

	providerJiraBlock := hclwrite.NewBlock("provider", []string{"jira"})
	providerJiraBody := providerJiraBlock.Body()
	providerJiraBody.SetAttributeValue("url", cty.StringVal("https://jira.eu.idealo.com/issues"))

	providerAwsTokens := providerAwsBlock.BuildTokens(nil)
	providerControlTowerTokens := providerControlTowerBlock.BuildTokens(nil)
	providerJiraTokens := providerJiraBlock.BuildTokens(nil)

	tokens := hclwrite.Tokens{
		{Type: hclsyntax.TokenEOF, Bytes: []byte("<<EOF")},
		{Type: hclsyntax.TokenNewline, Bytes: []byte("\n")},
	}

	tokens = append(tokens, providerAwsTokens...)
	tokens = append(tokens, hclwrite.Tokens{
		{Type: hclsyntax.TokenNewline, Bytes: []byte("\n")},
	}...)
	tokens = append(tokens, providerControlTowerTokens...)
	tokens = append(tokens, hclwrite.Tokens{
		{Type: hclsyntax.TokenNewline, Bytes: []byte("\n")},
	}...)
	tokens = append(tokens, providerJiraTokens...)
	tokens = append(tokens, hclwrite.Tokens{
		{Type: hclsyntax.TokenEOF, Bytes: []byte("EOF")},
		{Type: hclsyntax.TokenNewline, Bytes: []byte("\n")},
	}...)

	generateBody.AppendUnstructuredTokens(hclwrite.Tokens{
		{Type: hclsyntax.TokenStringLit, Bytes: []byte("# language=HCL")},
		{Type: hclsyntax.TokenNewline, Bytes: []byte("\n")},
	})
	//generateBody.AppendUnstructuredTokens(tokens)
	generateBody.SetAttributeRaw("contents", tokens)
}
