package shellexec

type shellError struct {
	Message string
}

func (e shellError) Error() string {
	return e.Message
}
