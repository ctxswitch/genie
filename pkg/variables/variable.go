package variables

// Variables is a collection of variable names and values.
type Variables struct {
	// how to handle this?  I'll just use a map for now, but I think I want something more
	// with type references and such.
	vars map[string]string
}

// Parse parses a collection of variable configs into a Variables object.
func Parse(block []Config) (*Variables, error) {
	vars := make(map[string]string)
	for _, v := range block {
		vars[v.Name] = v.Value
	}
	return &Variables{
		vars: vars,
	}, nil
}

// Get returns the value of a variable in the collection.
func (v *Variables) Get(name string) (string, bool) {
	val, ok := v.vars[name]
	return val, ok
}

// Set sets the value of a variable in the collection.
func (v *Variables) Set(name, value string) error {
	v.vars[name] = value
	return nil
}

// DeepCopy returns a deep copy of the variables.
func (v *Variables) DeepCopy() *Variables {
	out := make(map[string]string)
	for k, v := range v.vars {
		out[k] = v
	}
	return &Variables{
		vars: out,
	}
}
