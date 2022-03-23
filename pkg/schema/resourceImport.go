package schema

type ResourceCount struct {
	Dir     string
	Total   int
	Message string
}

type ImportResource struct {
	Address string
	Id      string
}

type Status int

const (
	Ok  Status = 0
	Err Status = 1
)

type Result interface {
	GetStatus() Status
	GetResourceCount() int
}

type ErrResult struct {
	Status        Status
	Dir           string
	ResourceCount int
	Message       string
}

func (e *ErrResult) GetStatus() Status {
	return e.Status
}

func (e *ErrResult) GetResourceCount() int {
	return e.ResourceCount
}

type OkResult struct {
	Status        Status
	Dir           string
	ResourceCount int
	Message       string
}

func (ok *OkResult) GetStatus() Status {
	return ok.Status
}

func (ok *OkResult) GetResourceCount() int {
	return ok.ResourceCount
}

type Done struct {
	Status         Status
	Dir            string
	FailedResource *ImportResource
	ResourceCount  int
	Message        string
	CleanErrors    string
}

func (done *Done) GetStatus() Status {
	return done.Status
}

func (done *Done) GetResourceCount() int {
	return done.ResourceCount
}
