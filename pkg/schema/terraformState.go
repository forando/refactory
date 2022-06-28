package schema

import "github.com/hashicorp/hcl/v2"

type TerraformState struct {
	Version          string     `hcl:"version,attr"`
	TerraformVersion string     `hcl:"terraform_version,attr"`
	Serial           int        `hcl:"serial,attr"`
	Lineage          string     `hcl:"lineage,attr"`
	Outputs          hcl.Body   `hcl:"lineage,body"`
	Resources        []Resource `hcl:"resources,block"`
	Rest             hcl.Body   `hcl:"rest,remain"`
}

type Resource struct {
	Module    string     `hcl:"module,optional"`
	Mode      string     `hcl:"mode,attr"`
	Type      string     `hcl:"type,attr"`
	Name      string     `hcl:"name,attr"`
	Provider  string     `hcl:"provider,attr"`
	Instances []Instance `hcl:"instances,block"`
}

type Instance struct {
	Attrs Attributes `hcl:"attributes,block"`
	Rest  hcl.Body   `hcl:"rest,remain"`
}

type IndexedInstance struct {
	IndexKey string   `hcl:"index_key,attr"`
	Rest     hcl.Body `hcl:"rest,remain"`
}

type Attributes struct {
	Id     string   `hcl:"id,attr"`
	Bucket string   `hcl:"bucket,optional"`
	Rest   hcl.Body `hcl:"rest,remain"`
}

type IngressAttributes struct {
	IngressProtocol     string   `hcl:"protocol,optional"`
	IngressNetworkAclId string   `hcl:"network_acl_id,optional"`
	IngressRuleNumber   int64    `hcl:"rule_number,optional"`
	IngressToPort       int64    `hcl:"to_port,optional"`
	IngressEgress       bool     `hcl:"egress,optional"`
	Rest                hcl.Body `hcl:"rest,remain"`
}

type RouteAttributes struct {
	RouteDestinationCidrBlock string   `hcl:"destination_cidr_block,optional"`
	RouteTableId              string   `hcl:"route_table_id,optional"`
	Rest                      hcl.Body `hcl:"rest,remain"`
}

type PeeringAccepterAttributes struct {
	VpcId               string   `hcl:"vpc_id,optional"`
	PeeringConnectionId string   `hcl:"vpc_peering_connection_id,optional"`
	Rest                hcl.Body `hcl:"rest,remain"`
}

type PeeringConnectionAttributes struct {
	VpcId            string   `hcl:"vpc_id,optional"`
	PeerVpcId        string   `hcl:"peer_vpc,optional"`
	PeerCloudAccount string   `hcl:"peer_cloud_account,optional"`
	Rest             hcl.Body `hcl:"rest,remain"`
}

type TfImport struct {
	Address string
	Id      string
}
