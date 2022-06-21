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
	AwsRoutResources   map[string]AwsRouteResource
	AwsNetworkAclRules map[string]AwsNetworkAclRule
	ConnectionAccepter AwsVpcPeeringConnectionAccepter
}

type AwsNetworkAclRule struct {
	Key
	IngressProtocol     string
	IngressNetworkAclId string
	IngressRuleNumber   int64
	IngressDenyToPort   int64
	IngressEgress       bool
}

type AwsVpcPeeringConnectionAccepter struct {
	Key
	VpcId               string
	PeeringConnectionId string
}

type AwsRouteResource struct {
	Key
	RouteTableId              string
	RouteDestinationCidrBlock string
}

type Key struct {
	Address string
	Id      string
}
