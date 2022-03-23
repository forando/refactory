package schema

type BlockMetaData struct {
	BlockType string
	BlockName string
}

const (
	AccountModuleType       string = "aws-account"
	PermissionSetModuleType string = "aws-ssoadmin-permission-set"
	IamPolicyDocumentType   string = "aws_iam_policy_document"
)

const (
	ModuleBlock string = "module"
	DataBlock   string = "data"
)

type Org int

const (
	TestOrg Org = 0
	ProdOrg Org = 1
)
