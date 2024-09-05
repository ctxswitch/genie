package variables

import "sync"

// ScopedVariables is a stack of variables that represent the specific
// scope that a variable is in.  It's currently not used past the global
// scope, but as functions like "for" are added, it will be necessary move
// in and out of their own scopes.
// TODO: This is super naive right now as it doesn't handle more interesting
// cases like variable access from a parent for vars that don't exist in the
// current scope.
type ScopedVariables struct {
	vars []*Variables
	pos  int
	sync.Mutex
}

// NewScopedVariables creates a new scoped variables object.
func NewScopedVariables(v *Variables) *ScopedVariables {
	vars := make([]*Variables, 0)
	vars = append(vars, v)

	return &ScopedVariables{
		vars: vars,
		pos:  0,
	}
}

// NewScope creates a new scope and pushes it onto the stack.
func (v *ScopedVariables) NewScope() {
	v.Lock()
	defer v.Unlock()

	v.vars = append(v.vars, v.vars[v.pos].DeepCopy())
	v.pos++
}

// ExitScope pops the current scope off the stack and restores the previous
// scope.
func (v *ScopedVariables) ExitScope() {
	v.vars[v.pos] = nil
	v.pos--
}

// Get returns the value of a variable if it exists in the current scope.
// TODO: recursively go back to get a variable from a parent scope if it doesn't
// exist in the current scope.
func (v *ScopedVariables) Get(name string) (string, bool) {
	v.Lock()
	defer v.Unlock()

	return v.vars[v.pos].Get(name)
}

// Set sets the value of a variable in the current scope.
// TODO: recursively go back to set a variable in a parent scope if it doesn't
// exist in the current scope.  If we get to the end of the scopes and it's
// not found, create it in the current scope.
func (v *ScopedVariables) Set(name, value string) error {
	v.Lock()
	defer v.Unlock()

	return v.vars[v.pos].Set(name, value)
}

// Len returns the number of scopes in the stack.
func (v *ScopedVariables) Len() int {
	v.Lock()
	defer v.Unlock()

	return v.pos + 1
}
