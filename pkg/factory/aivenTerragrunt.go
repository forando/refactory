package factory

import (
	"github.com/forando/refactory/pkg/filesystem"
	"github.com/forando/refactory/pkg/schema"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/pkg/errors"
	"github.com/zclconf/go-cty/cty"
	"os"
	"path/filepath"
)

type AivenTerragrunt struct {
	Name string
	Dir  string
}

func NewAivenTerragrunt(dir string) AivenFactory {
	return &AivenTerragrunt{Name: "terragrunt", Dir: dir}
}

func (t *AivenTerragrunt) BootstrapNewModule(consumers *[]schema.AivenConsumerModule) error {
	fs := filesystem.NewOsFs()
	newConfigDir := filepath.Join(t.Dir, "new_config")
	fs.MakeDirs(newConfigDir)
	aivenModuleDir := filepath.Join(newConfigDir, "aiven-aws-vpc-peering-acceptor")
	fs.MakeDirs(aivenModuleDir)
	if err := t.bootstrapAivenModule(aivenModuleDir); err != nil {
		return err
	}
	if err := t.bootstrapTerragruntConfig(newConfigDir, consumers); err != nil {
		return err
	}
	return nil
}

func (t *AivenTerragrunt) bootstrapAivenModule(dir string) error {
	if err := t.bootstrapAivenModuleMain(dir); err != nil {
		return err
	}
	if err := t.bootstrapAivenModuleVariables(dir); err != nil {
		return err
	}
	if err := t.bootstrapAivenModuleOutputs(dir); err != nil {
		return err
	}
	return nil
}

func (t *AivenTerragrunt) bootstrapAivenModuleMain(dir string) error {
	filePath := filepath.Join(dir, "main.tf")

	fw, osErr := os.Create(filePath)

	if osErr != nil {
		return osErr
	}

	newFile := hclwrite.NewEmptyFile()
	rootBody := newFile.Body()

	bootstrapPeeringConnectionsModule(rootBody, "var.vpc_peering_connections")

	if _, err := newFile.WriteTo(fw); err != nil {
		return errors.WithMessage(err, "Cannot write to the new file ")
	}

	return nil
}

func (t *AivenTerragrunt) bootstrapAivenModuleVariables(dir string) error {

	filePath := filepath.Join(dir, "variables.tf")

	fw, osErr := os.Create(filePath)

	if osErr != nil {
		return osErr
	}

	newFile := hclwrite.NewEmptyFile()
	rootBody := newFile.Body()

	vpcPeeringConnectionsBlock := rootBody.AppendNewBlock("variable", []string{"vpc_peering_connections"})
	vpcPeeringConnectionsBody := vpcPeeringConnectionsBlock.Body()

	vpcPeeringConnectionsBody.SetAttributeRaw("type", hclwrite.Tokens{
		{Type: hclsyntax.TokenIdent, Bytes: []byte("list(map(string))"), SpacesBefore: 1},
	})
	vpcPeeringConnectionsBody.SetAttributeValue("description", cty.StringVal("Your list of peering connection objects"))

	if _, err := newFile.WriteTo(fw); err != nil {
		return errors.WithMessage(err, "Cannot write to the new file ")
	}

	return nil
}

func (t *AivenTerragrunt) bootstrapAivenModuleOutputs(dir string) error {
	filePath := filepath.Join(dir, "outputs.tf")

	_, osErr := os.Create(filePath)

	if osErr != nil {
		return osErr
	}
	return nil
}

func (t *AivenTerragrunt) bootstrapTerragruntConfig(dir string, consumers *[]schema.AivenConsumerModule) error {
	filePath := filepath.Join(dir, "terragrunt.hcl")

	fw, osErr := os.Create(filePath)

	if osErr != nil {
		return osErr
	}

	newFile := hclwrite.NewEmptyFile()
	rootBody := newFile.Body()

	terraformBlock := rootBody.AppendNewBlock("terraform", nil)
	terraformBody := terraformBlock.Body()

	terraformBody.SetAttributeRaw("source", hclwrite.Tokens{
		{Type: hclsyntax.TokenOQuote, Bytes: []byte{'"'}},
		{Type: hclsyntax.TokenStringLit, Bytes: []byte("${path_relative_from_include()}/modules//aiven-aws-vpc-peering-acceptor")},
		{Type: hclsyntax.TokenCQuote, Bytes: []byte{'"'}},
	})

	rootBody.AppendNewline()

	includeBlock := rootBody.AppendNewBlock("include", nil)
	includeBody := includeBlock.Body()

	includeBody.SetAttributeRaw("path", hclwrite.Tokens{
		{Type: hclsyntax.TokenIdent, Bytes: []byte("find_in_parent_folders()")},
	})

	rootBody.AppendNewline()

	inputsBlock := rootBody.AppendNewBlock("inputs", nil)
	inputsBody := inputsBlock.Body()

	inputsBody.SetAttributeValue("vpc_peering_connections", *createPeeringConnectionsVal(consumers))

	if _, err := newFile.WriteTo(fw); err != nil {
		return errors.WithMessage(err, "Cannot write to the new file ")
	}

	return nil
}
