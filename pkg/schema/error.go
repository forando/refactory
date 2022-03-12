package schema

type ParsingError struct {
	Message string
}

func (e ParsingError) Error() string {
	return e.Message
}
