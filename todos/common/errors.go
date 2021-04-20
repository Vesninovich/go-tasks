package common

// InvalidInputError represents error on invalid input from user
type InvalidInputError struct {
	Reason string
}

func (e *InvalidInputError) Error() string {
	return "Invalid input: " + e.Reason
}

// NotFoundError represents error that something was not found and that's an error
type NotFoundError struct {
	What string
}

func (e *NotFoundError) Error() string {
	return e.What + " not found"
}
