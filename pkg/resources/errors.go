package resources

type Error string

func (r Error) Error() string { return string(r) }

const (
	NotFoundError            Error = "not found"
	InvalidResourceTypeError Error = "invalid resource type requested"
)
