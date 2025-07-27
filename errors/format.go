package errors

type FormatError struct {
}

func (e *FormatError) Error() string {
	return "parser not registered"
}
