package common

type InvalidInputError struct {
	Reason string
}

func (e *InvalidInputError) Error() string {
	return "Invalid input: " + e.Reason
}
