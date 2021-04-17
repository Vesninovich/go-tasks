package common

// InvalidInputError represents error on invalid input from user
type InvalidInputError struct {
	Reason string
}

func (e *InvalidInputError) Error() string {
	return "Invalid input: " + e.Reason
}
