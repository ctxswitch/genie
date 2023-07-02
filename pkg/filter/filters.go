package filter

import "errors"

type FilterFunc func(any) string

var FilterMap = map[string]FilterFunc{
	"capitalize":  Capitalize,
	"passthrough": Passthrough,
}

func Lookup(name string) (FilterFunc, error) {
	fn, ok := FilterMap[name]
	if !ok {
		return nil, errors.New("Unknown filter")
	}
	return fn, nil
}
