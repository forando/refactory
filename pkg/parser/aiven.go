package parser

import (
	"fmt"
	"github.com/forando/refactory/pkg/schema"
	"github.com/pkg/errors"
	"strings"
)

type ModuleType int

const (
	Unknown       ModuleType = 0
	AivenProducer ModuleType = 1
	AivenConsumer ModuleType = 2
)

type aivenCandidate struct {
	moduleType                      ModuleType
	aivenVpcPeeringConnection       *schema.AivenVpcPeeringConnection
	awsRoutResources                *map[string]schema.AwsRoutResource
	awsNetworkAclRules              *map[string]schema.AwsNetworkAclRule
	awsVpcPeeringConnectionAccepter *schema.AwsVpcPeeringConnectionAccepter
}

const (
	aivenVpcPeeringConnection       string = "aiven_vpc_peering_connection"
	awsNetworkAclRule               string = "aws_network_acl_rule"
	awsRoute                        string = "aws_route"
	awsVpcPeeringConnectionAccepter string = "aws_vpc_peering_connection_accepter"
)

var allowedResourceTypes = map[string]interface{}{
	aivenVpcPeeringConnection:       nil,
	awsNetworkAclRule:               nil,
	awsRoute:                        nil,
	awsVpcPeeringConnectionAccepter: nil,
}

func ParseAivenStateBytes(bytes *[]byte) (*map[string]schema.AivenProducerModule, *map[string]schema.AivenConsumerModule, error) {
	state, err := parseTfStateBytes(bytes)
	if err != nil {
		return nil, nil, err
	}
	return ParseAivenState(state)
}

func ParseAivenStateFile(file string) (*map[string]schema.AivenProducerModule, *map[string]schema.AivenConsumerModule, error) {
	state, err := parseTfStateFile(file)
	if err != nil {
		return nil, nil, err
	}
	return ParseAivenState(state)
}

func ParseAivenState(state *schema.TerraformState) (*map[string]schema.AivenProducerModule, *map[string]schema.AivenConsumerModule, error) {

	candidates := make(map[string]*aivenCandidate)
	for _, resource := range state.Resources {
		if len(resource.Module) == 0 || resource.Type == "null_resource" {
			continue
		}
		moduleName := resource.Module

		switch resource.Mode {
		case "data":
			continue
		case "managed":
			break
		default:
			return nil, nil, errors.Errorf("module: %s has usupported mode: [%s]", resource.Module, resource.Mode)
		}

		if _, found := allowedResourceTypes[resource.Type]; !found {
			continue
		}

		var found bool
		candidate, found := candidates[moduleName]
		if !found {
			candidate = &aivenCandidate{moduleType: Unknown}
			candidates[moduleName] = candidate
		}

		if resource.Type == aivenVpcPeeringConnection {
			candidate.moduleType = AivenProducer
		}
		if err := populateCandidate(candidate, &resource); err != nil {
			return nil, nil, err
		}
	}
	if producers, consumers, err := getAivenModules(&candidates); err != nil {
		return nil, nil, err
	} else {
		return producers, consumers, nil
	}
}

func populateCandidate(candidate *aivenCandidate, resource *schema.Resource) error {
	switch resource.Type {
	case aivenVpcPeeringConnection:
		length := len(resource.Instances)
		if length != 1 {
			return wrongInstanceNumber(resource.Module, resource.Type, length)
		}
		if candidate.aivenVpcPeeringConnection != nil {
			return resourceAlreadySet(resource.Module, resource.Type)
		}
		key := schema.Key{
			Address: fmt.Sprintf("%s.%s", resource.Type, resource.Name),
			Id:      resource.Instances[0].Attrs.Id,
		}
		candidate.aivenVpcPeeringConnection = &schema.AivenVpcPeeringConnection{
			Key:               key,
			AccountId:         resource.Instances[0].Attrs.PeerCloudAccount,
			AivenProjectVpcId: *resource.Instances[0].Attrs.VpcId,
			VpcId:             *resource.Instances[0].Attrs.PeerVpcId,
		}
		return nil
	case awsVpcPeeringConnectionAccepter:
		length := len(resource.Instances)
		if length != 1 {
			return wrongInstanceNumber(resource.Module, resource.Type, length)
		}
		if candidate.awsVpcPeeringConnectionAccepter != nil {
			return resourceAlreadySet(resource.Module, resource.Type)
		}
		key := schema.Key{
			Address: fmt.Sprintf("%s.%s", resource.Type, resource.Name),
			Id:      resource.Instances[0].Attrs.Id,
		}
		candidate.awsVpcPeeringConnectionAccepter = &schema.AwsVpcPeeringConnectionAccepter{
			Key:                 key,
			VpcId:               *resource.Instances[0].Attrs.VpcId,
			PeeringConnectionId: *resource.Instances[0].Attrs.PeeringConnectionId,
		}
		return nil
	case awsNetworkAclRule:
		length := len(resource.Instances)
		if length != 1 {
			return wrongInstanceNumber(resource.Module, resource.Type, length)
		}
		var rules map[string]schema.AwsNetworkAclRule
		if candidate.awsNetworkAclRules != nil {
			rules = *candidate.awsNetworkAclRules
		} else {
			rules = make(map[string]schema.AwsNetworkAclRule)
			candidate.awsNetworkAclRules = &rules
		}
		if _, found := rules[resource.Name]; found {
			return errors.Errorf(
				"module: %s of type: %s already hase object with name = %s",
				resource.Module,
				resource.Type,
				resource.Name,
			)
		}
		key := schema.Key{
			Address: fmt.Sprintf("%s.%s", resource.Type, resource.Name),
			Id:      resource.Instances[0].Attrs.Id,
		}
		rules[resource.Name] = schema.AwsNetworkAclRule{
			Key:               key,
			IngressRuleNumber: resource.Instances[0].Attrs.IngressRuleNumber,
			IngressDenyToPort: resource.Instances[0].Attrs.IngressToPort,
		}
		return nil
	case awsRoute:
		routes := make(map[string]schema.AwsRoutResource)
		if candidate.awsRoutResources != nil {
			routes = *candidate.awsRoutResources
		}
		for _, instance := range resource.Instances {
			if len(instance.IndexKey) == 0 {
				//return errors.Errorf("module: %s of type: %s does not have index_key property", resource.Module, resource.Type)
				return nil
			}
			if _, found := routes[instance.IndexKey]; found {
				return errors.Errorf(
					"module: %s of type: %s already hase instance object with index_key = %s",
					resource.Module,
					resource.Type,
					instance.IndexKey,
				)
			}
			address := fmt.Sprintf("%s.%s[\"%s\"]", resource.Type, resource.Name, instance.IndexKey)
			key := schema.Key{
				Address: address,
				Id:      resource.Instances[0].Attrs.Id,
			}
			route := schema.AwsRoutResource{Key: key}
			routes[instance.IndexKey] = route
		}
		if candidate.awsRoutResources == nil {
			candidate.awsRoutResources = &routes
		}
		return nil
	default:
		return errors.Errorf(
			"Unexpected module: %s of type: %s",
			resource.Module,
			resource.Type,
		)
	}
}

func wrongInstanceNumber(module string, tp string, length int) error {
	return errors.Errorf(
		"module: %s of type: %s expected to have exactly 1 instance, found %d",
		module,
		tp,
		length,
	)
}

func resourceAlreadySet(module string, tp string) error {
	return errors.Errorf(
		"module: %s has more than 1 %s resource",
		module,
		tp,
	)
}

func getAivenModules(candidates *map[string]*aivenCandidate) (*map[string]schema.AivenProducerModule, *map[string]schema.AivenConsumerModule, error) {
	producers := make(map[string]schema.AivenProducerModule)
	consumers := make(map[string]schema.AivenConsumerModule)

	for key, candidate := range *candidates {
		if candidate.awsRoutResources == nil {
			continue
		}
		if candidate.awsNetworkAclRules == nil {
			continue
		}
		if candidate.awsVpcPeeringConnectionAccepter == nil {
			continue
		}
		consumers[key] = schema.AivenConsumerModule{
			Name:               key,
			AwsNetworkAclRules: *candidate.awsNetworkAclRules,
			AwsRoutResources:   *candidate.awsRoutResources,
			ConnectionAccepter: *candidate.awsVpcPeeringConnectionAccepter,
		}
	}
	for key, candidate := range *candidates {
		if candidate.moduleType == AivenProducer {
			if candidate.aivenVpcPeeringConnection == nil {
				continue
			}
			found := false
			for subKey, consumer := range consumers {
				if strings.Contains(subKey, key) {
					producers[key] = schema.AivenProducerModule{
						Name:              key,
						PeeringConnection: *candidate.aivenVpcPeeringConnection,
						Consumer:          &consumer,
					}
					delete(consumers, subKey)
					found = true
					break
				}
			}
			if !found {
				producers[key] = schema.AivenProducerModule{
					Name:              key,
					PeeringConnection: *candidate.aivenVpcPeeringConnection,
				}
			}
		}
	}

	return &producers, &consumers, nil
}
