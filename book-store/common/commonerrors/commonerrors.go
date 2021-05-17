package commonerrors

// NotFound represents error that something was not found and that's an error
type NotFound struct {
	What string
}

func (e NotFound) Error() string {
	return e.What + " not found"
}

// InvalidInput represents error that someone provided invalid input
type InvalidInput struct {
	Reason string
}

func (e InvalidInput) Error() string {
	return "Invalid input: " + e.Reason
}
