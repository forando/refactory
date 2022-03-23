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
	Attrs    Attributes `hcl:"attributes,block"`
	IndexKey string     `hcl:"index_key,optional"`
	Rest     hcl.Body   `hcl:"rest,remain"`
}

type Attributes struct {
	Id     string   `hcl:"id,attr"`
	Bucket string   `hcl:"bucket,optional"`
	Rest   hcl.Body `hcl:"rest,remain"`
}

type TfImport struct {
	Address string
	Id      string
}