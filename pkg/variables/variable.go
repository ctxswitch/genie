package variables

type Variables struct {
	// how to handle this?  I'll just use a map for now, but I think I want something more
	// with type refereces and such.
	vars map[string]string
}

func Parse(block []Config) (*Variables, error) {
	vars := make(map[string]string)
	for _, v := range block {
		vars[v.Name] = v.Value
	}
	return &Variables{
		vars: vars,
	}, nil
}

func (v *Variables) Get(name string) (string, bool) {
	val, ok := v.vars[name]
	return val, ok
}

func (v *Variables) Set(name, value string) error {
	v.vars[name] = value
	return nil
}

func (v *Variables) DeepCopy() *Variables {
	out := make(map[string]string)
	for k, v := range v.vars {
		out[k] = v
	}
	return &Variables{
		vars: out,
	}
}
