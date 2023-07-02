package list

type List []string

func (l List) Get() string {
	return l[0]
}
