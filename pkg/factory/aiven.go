package factory

import (
	"fmt"
	"github.com/forando/refactory/pkg/schema"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/pkg/errors"
	"github.com/zclconf/go-cty/cty"
	"os"
	"path/filepath"
)

type AivenTerraform struct {
	Name string
	Dir  string
}

func NewAivenTerraform(dir string) *AivenTerraform {
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

	peeringConnections := make([]cty.Value, 0)
	for _, consumer := range *consumers {
		tcpRule := consumer.AwsNetworkAclRules[schema.IngressTcp]
		udpRule := consumer.AwsNetworkAclRules[schema.IngressUdp]

		entries := map[string]cty.Value{
			"id":                            cty.StringVal(consumer.ConnectionAccepter.PeeringConnectionId),
			"vpc_id":                        cty.StringVal(consumer.ConnectionAccepter.VpcId),
			"ingress_tcp_rule_number":       cty.NumberIntVal(tcpRule.IngressRuleNumber),
			"ingress_udp_rule_number":       cty.NumberIntVal(udpRule.IngressRuleNumber),
			"aws_nacl_ingress_deny_to_port": cty.NumberIntVal(udpRule.IngressDenyToPort),
		}

		peeringConnections = append(peeringConnections, cty.ObjectVal(entries))
	}
	if len(peeringConnections) > 0 {
		localsBody.SetAttributeValue("vpc_peering_connections", cty.ListVal(peeringConnections))
	} else {
		localsBody.SetAttributeValue("vpc_peering_connections", cty.ListValEmpty(cty.String))
	}

	rootBody.AppendNewline()

	moduleBlock := rootBody.AppendNewBlock("module", []string{schema.NewModuleDefaultName})
	moduleBody := moduleBlock.Body()

	moduleBody.SetAttributeValue("source", cty.StringVal("git@github.com:idealo/terraform-aiven-vpc-peering//modules/aiven-aws-peering-connections-acceptor?ref=v2.0.1"))

	moduleBody.SetAttributeRaw("for_each", hclwrite.Tokens{
		{Type: hclsyntax.TokenStringLit, Bytes: []byte(" { for connection in local.vpc_peering_connections : ")},
		{Type: hclsyntax.TokenOQuote, Bytes: []byte{'"'}},
		{Type: hclsyntax.TokenStringLit, Bytes: []byte("${connection.vpc_id}/${connection.id}")},
		{Type: hclsyntax.TokenCQuote, Bytes: []byte{'"'}},
		{Type: hclsyntax.TokenStringLit, Bytes: []byte(" => connection }")},
	})
	moduleBody.AppendNewline()

	moduleBody.SetAttributeRaw("aws_vpc_peering_connection_id", hclwrite.Tokens{
		{Type: hclsyntax.TokenStringLit, Bytes: []byte(" each.value.id")},
	})

	moduleBody.SetAttributeRaw("aws_vpc_id", hclwrite.Tokens{
		{Type: hclsyntax.TokenStringLit, Bytes: []byte(" each.value.vpc_id")},
	})

	moduleBody.SetAttributeRaw("ingress_tcp_rule_number", hclwrite.Tokens{
		{Type: hclsyntax.TokenStringLit, Bytes: []byte(" each.value.ingress_tcp_rule_number")},
	})
	moduleBody.SetAttributeRaw("ingress_udp_rule_number", hclwrite.Tokens{
		{Type: hclsyntax.TokenStringLit, Bytes: []byte(" each.value.ingress_udp_rule_number")},
	})
	moduleBody.SetAttributeRaw("aws_nacl_ingress_deny_to_port", hclwrite.Tokens{
		{Type: hclsyntax.TokenStringLit, Bytes: []byte(" each.value.aws_nacl_ingress_deny_to_port")},
	})
}
