package schema

const (
	IngressTcp = "ingress_tcp"
	IngressUdp = "ingress_udp"
)

type AivenProducerModule struct {
	Name              string
	PeeringConnection AivenVpcPeeringConnection
	Consumer          *AivenConsumerModule
}

type AivenVpcPeeringConnection struct {
	Key
	AivenProjectVpcId string
	VpcId             string
	AccountId         string
}

type AivenConsumerModule struct {
	Name               string
	AwsRoutResources   map[string]AwsRoutResource
	AwsNetworkAclRules map[string]AwsNetworkAclRule
	ConnectionAccepter AwsVpcPeeringConnectionAccepter
}

type AwsNetworkAclRule struct {
	Key
	IngressRuleNumber int64
	IngressDenyToPort int64
}

type AwsVpcPeeringConnectionAccepter struct {
	Key
	VpcId               string
	PeeringConnectionId string
}

type AwsRoutResource struct {
	Key
}

type Key struct {
	Address string
	Id      string
}
