package filter

// Passthrough returns the string s.  It will be removed later as it was only
// added as a test.
func Passthrough(s any) string {
	return s.(string)
}
