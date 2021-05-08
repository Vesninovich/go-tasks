package commonerrors

// NotFound represents error that something was not found and that's an error
type NotFound struct {
	What string
}

func (e NotFound) Error() string {
	return e.What + " not found"
}
