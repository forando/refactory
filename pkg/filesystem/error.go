package filesystem

type FsError struct {
	Message string
}

func (e FsError) Error() string {
	return e.Message
}
