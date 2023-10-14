package filter

import "errors"

type Func func(any) string

var FilterMap = map[string]Func{
	"capitalize":  Capitalize,
	"passthrough": Passthrough,
}

func Lookup(name string) (Func, error) {
	fn, ok := FilterMap[name]
	if !ok {
		return nil, errors.New("Unknown filter")
	}
	return fn, nil
}
