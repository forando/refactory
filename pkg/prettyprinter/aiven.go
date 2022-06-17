package prettyprinter

import (
	"fmt"
	"github.com/forando/refactory/pkg/schema"
)

func PrintAivenResources(producers *map[string]schema.AivenProducerModule, consumers *map[string]schema.AivenConsumerModule) {

	fmt.Printf("producers: %d, consumers: %d\n", len(*producers), len(*consumers))
	if len(*producers) > 0 {
		fmt.Printf("Producders: %d\n", len(*producers))
		for _, producer := range *producers {
			fmt.Printf("  Module: %s\n", producer.Name)
			fmt.Printf("  |__PeeringConnection\n")
			if producer.Consumer != nil {
				printProducer(&producer, "  |  ")
				fmt.Printf("  |__ConnectionAccepter Resources:\n")
				printConsumer(producer.Consumer, "     ")
			} else {
				printProducer(&producer, "     ")
			}
			println()
		}
	}
	if len(*consumers) > 0 {
		fmt.Printf("Consumers: %d\n", len(*consumers))
		for _, consumer := range *consumers {
			printConsumer(&consumer, "  ")
			println()
		}
	}
}

func printProducer(producer *schema.AivenProducerModule, prefix string) {
	fmt.Printf("%s|__Id: %s\n", prefix, producer.PeeringConnection.Id)
	fmt.Printf("%s|__Address: %s\n", prefix, producer.PeeringConnection.Address)
	fmt.Printf("%s|__AivenProjectVpcId: %s\n", prefix, producer.PeeringConnection.AivenProjectVpcId)
	fmt.Printf("%s|__VpcId: %s\n", prefix, producer.PeeringConnection.VpcId)
	fmt.Printf("%s|__AccountId: %s\n", prefix, producer.PeeringConnection.AccountId)
}

func printConsumer(consumer *schema.AivenConsumerModule, indent string) {
	fmt.Printf("%sModule: %s\n", indent, consumer.Name)
	fmt.Printf("%s|__ConnectionAccepter\n", indent)
	fmt.Printf("%s|  |__Id: %s\n", indent, consumer.ConnectionAccepter.Id)
	fmt.Printf("%s|  |__Address: %s\n", indent, consumer.ConnectionAccepter.Address)
	fmt.Printf("%s|  |__VpcId: %s\n", indent, consumer.ConnectionAccepter.VpcId)
	fmt.Printf("%s|  |__PeeringConnectionId: %s\n", indent, consumer.ConnectionAccepter.PeeringConnectionId)
	fmt.Printf("%s|__AwsNetworkAclRules:\n", indent)
	index := 0
	for key, val := range consumer.AwsNetworkAclRules {
		fmt.Printf("%s|   |_____%s:\n", indent, key)
		if index == 0 {
			fmt.Printf("%s|   |      |__Id: %s\n", indent, val.Id)
			fmt.Printf("%s|   |      |__Address: %s\n", indent, val.Address)
			fmt.Printf("%s|   |      |__IngressRuleNumber: %d\n", indent, val.IngressRuleNumber)
			fmt.Printf("%s|   |      |__IngressDenyToPort: %d\n", indent, val.IngressDenyToPort)
		} else {
			fmt.Printf("%s|          |__Id: %s\n", indent, val.Id)
			fmt.Printf("%s|          |__Address: %s\n", indent, val.Address)
			fmt.Printf("%s|          |__IngressRuleNumber: %d\n", indent, val.IngressRuleNumber)
			fmt.Printf("%s|          |__IngressDenyToPort: %d\n", indent, val.IngressDenyToPort)
		}
		index++
	}
	fmt.Printf("%s|__AwsRoutResources:\n", indent)
	for key, val := range consumer.AwsRoutResources {
		fmt.Printf("%s    |_____%s {Id: %s, Address: %s}\n", indent, key, val.Id, val.Address)
	}
}
