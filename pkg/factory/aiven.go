package factory

import (
	"github.com/forando/refactory/pkg/schema"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

type AivenFactory interface {
	BootstrapNewModule(consumers *[]schema.AivenConsumerModule) error
}

func GetAivenFactory(toolName string, newModulePath string) AivenFactory {
	if toolName == "terragrunt" {
		return NewAivenTerragrunt(newModulePath)
	}
	return NewAivenTerraform(newModulePath)
}

func bootstrapPeeringConnectionsModule(rootBody *hclwrite.Body, vpcPeeringConnectionsVar string) {
	moduleBlock := rootBody.AppendNewBlock("module", []string{schema.NewModuleDefaultName})
	moduleBody := moduleBlock.Body()

	moduleBody.SetAttributeValue("source", cty.StringVal("git@github.com:idealo/terraform-aiven-vpc-peering//modules/aiven-aws-peering-connections-acceptor?ref=v2.0.1"))

	moduleBody.SetAttributeRaw("for_each", hclwrite.Tokens{
		{Type: hclsyntax.TokenIdent, Bytes: []byte("{ for connection in ")},
		{Type: hclsyntax.TokenStringLit, Bytes: []byte(vpcPeeringConnectionsVar)},
		{Type: hclsyntax.TokenStringLit, Bytes: []byte(" : ")},
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

func createPeeringConnectionsVal(consumers *[]schema.AivenConsumerModule) *cty.Value {
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
	var res cty.Value
	if len(peeringConnections) > 0 {
		res = cty.ListVal(peeringConnections)
	} else {
		res = cty.ListValEmpty(cty.String)
	}
	return &res
}
