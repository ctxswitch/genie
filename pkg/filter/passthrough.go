package filter

func Passthrough(s any) string {
	return s.(string)
}
