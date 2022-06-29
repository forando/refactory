package factory

import (
	"fmt"
	"github.com/forando/refactory/pkg/schema"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

type AivenTerraform struct {
	Name string
	Dir  string
}

func NewAivenTerraform(dir string) AivenFactory {
	return &AivenTerraform{Name: "terraform", Dir: dir}
}

func (t *AivenTerraform) BootstrapNewModule(consumers *[]schema.AivenConsumerModule) error {

	const terragruntFileNme = "aiven_peering_verbose.tf"
	filePath := filepath.Join(t.Dir, terragruntFileNme)

	fw, osErr := os.Create(filePath)

	if osErr != nil {
		return osErr
	}

	newFile := hclwrite.NewEmptyFile()
	rootBody := newFile.Body()

	t.bootstrap(rootBody, consumers)

	if _, err := newFile.WriteTo(fw); err != nil {
		return errors.WithMessage(err, "Cannot write to the new file ")
	}
	fmt.Printf("A new module was created at %s\n", filePath)
	fmt.Println("Please have a look and adjust it according to your needs.")

	return nil
}

func (t *AivenTerraform) bootstrap(rootBody *hclwrite.Body, consumers *[]schema.AivenConsumerModule) {
	localsBlock := rootBody.AppendNewBlock("locals", nil)
	localsBody := localsBlock.Body()

	localsBody.SetAttributeValue("vpc_peering_connections", *createPeeringConnectionsVal(consumers))

	rootBody.AppendNewline()

	bootstrapPeeringConnectionsModule(rootBody, "local.vpc_peering_connections")
}
