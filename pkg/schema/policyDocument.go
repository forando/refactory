package schema

type PolicyDocument struct {
	Statements []Statement `hcl:"statement,block"`
}

type Statement struct {
	Sid       string   `hcl:"sid,optional"`
	Effect    string   `hcl:"effect,optional"`
	Actions   []string `hcl:"actions"`
	Resources []string `hcl:"resources"`
}

const (
	EffectAllow string = "Allow"
	EffectDeny  string = "Deny"
)

type PolicyType int

const (
	ManagedPolicy PolicyType = 0
	InlinePolicy  PolicyType = 1
)
