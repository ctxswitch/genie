package dict

// Map is a map of string to string.  This is currently not represented in the configs
// or templates.
type Map map[string]string

// Get implements the Resource interface and returns an empty string.
func (m *Map) Get() string { return "" }
