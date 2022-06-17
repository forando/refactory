package shellexec

type IaC interface {
	Name() string
	Init() error
	StatePull() (*[]byte, error)
	StateList() (*[]string, error)
	StateMove(src string, dest string, dryRun bool) error
	StateRemove(address string, dryRun bool) error
}
