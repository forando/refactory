package schema

type AccountModule struct {
	ProductTicket      string
	AccountName        string
	OrganizationalUnit string
	CostCenter         int
	OwnerEmail         string
	OwnerJiraUsername  string
	GroupPermissions   map[string][]string
	UserPermissions    map[string][]string
}

type AccountModules []*AccountModule
