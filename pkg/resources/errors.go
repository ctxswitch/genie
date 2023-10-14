package resources

// Error is a custom error type for resources.
type Error string

// Error implements the error interface.
func (r Error) Error() string { return string(r) }

const (
	// NotFoundError is returned when a resource is not found.
	NotFoundError Error = "not found"
	// InvalidResourceTypeError is returned when a resource type is not valid.
	InvalidResourceTypeError Error = "invalid resource type requested"
)
