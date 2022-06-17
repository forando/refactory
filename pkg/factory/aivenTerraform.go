package factory

import (
	"fmt"
	"github.com/forando/refactory/pkg/schema"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	"log"
	"os"
	"path/filepath"
)

func BootstrapAivenModule(consumers *map[string]schema.AivenConsumerModule, dir string) {

	//length := len(*consumers)
	const terragruntFileNme = "aiven_peering_verbose.tf"
	filePath := filepath.Join(dir, terragruntFileNme)

	fw, osErr := os.Create(filePath)

	if osErr != nil {
		log.Fatal("Cannot create new file ", osErr)
	}

	newFile := hclwrite.NewEmptyFile()
	rootBody := newFile.Body()

	localsBlock := rootBody.AppendNewBlock("locals", nil)
	localsBody := localsBlock.Body()

	peeringConnections := make([]cty.Value, 0)
	//var vpcId string
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
		//vpcId = consumer.ConnectionAccepter.VpcId
	}
	if len(peeringConnections) > 0 {
		localsBody.SetAttributeValue("vpc_peering_connections", cty.ListVal(peeringConnections))
	} else {
		localsBody.SetAttributeValue("vpc_peering_connections", cty.ListValEmpty(cty.String))
	}
	//localsBody.SetAttributeValue("vpc_id", cty.StringVal(vpcId))

	rootBody.AppendNewline()

	moduleBlock := rootBody.AppendNewBlock("module", []string{"peering_connection"})
	moduleBody := moduleBlock.Body()

	moduleBody.SetAttributeValue("source", cty.StringVal("git@github.com:idealo/terraform-aiven-vpc-peering//modules/aiven-aws-peering-connections-acceptor?ref=v2.0.0"))

	//for_each = { for connection in local.vpc_peering_connections : "${local.vpc_id}/${connection.id}" => connection }

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

	if _, err := newFile.WriteTo(fw); err != nil {
		log.Fatal("Cannot write to the new file ", err)
	}
	fmt.Printf("A new module was created at %s\n", filePath)
	fmt.Println("Please have a look and adjust it according to your needs.")
}
