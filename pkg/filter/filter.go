package filter

import "errors"

// Func is a function that takes an argument of any type and returns a string. It
// is used to implement expression compatible filters.
type Func func(any) string

// FilterMap is a map of filter names to filter functions.
var FilterMap = map[string]Func{ //nolint:gochecknoglobals
	"capitalize":  Capitalize,
	"passthrough": Passthrough,
}

// Lookup returns the filter function for the given name. If the name is not
// found then an error is returned.
func Lookup(name string) (Func, error) {
	fn, ok := FilterMap[name]
	if !ok {
		return nil, errors.New("unknown filter")
	}
	return fn, nil
}
