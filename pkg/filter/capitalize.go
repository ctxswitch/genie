package filter

import "strings"

func Capitalize(s any) string {
	return strings.ToUpper(s.(string))
}
