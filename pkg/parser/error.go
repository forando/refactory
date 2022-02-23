package parser

type ParsingError struct {
	message string
}

func (e ParsingError) Error() string {
	return e.message
}
