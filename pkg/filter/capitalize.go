package filter

import "strings"

// Capitalize returns the string s with all Unicode letters mapped to their
// Unicode upper case.
func Capitalize(s any) string {
	return strings.ToUpper(s.(string))
}
