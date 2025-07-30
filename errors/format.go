package errors

type FormatError struct {
	Err string
}

func (e *FormatError) Error() string {
	return "parser not registered: " + e.Err
}
