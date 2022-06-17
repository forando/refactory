package schema

type CallerIdentity struct {
	UserId  string `hcl:"UserId,attr"`
	Account string `hcl:"Account,attr"`
	Arn     string `hcl:"Arn,attr"`
}
