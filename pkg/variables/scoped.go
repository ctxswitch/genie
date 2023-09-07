package variables

import "sync"

type ScopedVariables struct {
	vars []*Variables
	pos  int
	sync.Mutex
}

func NewScopedVariables(v *Variables) *ScopedVariables {
	vars := make([]*Variables, 0)
	vars = append(vars, v)

	return &ScopedVariables{
		vars: vars,
		pos:  0,
	}
}

func (v *ScopedVariables) NewScope() {
	v.Lock()
	defer v.Unlock()

	v.vars = append(v.vars, v.vars[v.pos].DeepCopy())
	v.pos++
}

func (v *ScopedVariables) ExitScope() {
	v.vars[v.pos] = nil
	v.pos--
}

// TODO: recursively go back to get a variable from a parent scope if it doesn't
// exist in the current scope.
func (v *ScopedVariables) Get(name string) (string, bool) {
	v.Lock()
	defer v.Unlock()

	return v.vars[v.pos].Get(name)
}

// TODO: recursively go back to set a variable in a parent scope if it doesn't
// exist in the current scope.  If we get to the end of the scopes and it's
// not found, create it in the current scope.
func (v *ScopedVariables) Set(name, value string) error {
	v.Lock()
	defer v.Unlock()

	return v.vars[v.pos].Set(name, value)
}

func (v *ScopedVariables) Len() int {
	v.Lock()
	defer v.Unlock()

	return v.pos + 1
}
