package list

type Null struct{}

func (Null) Get() string {
	return "null"
}
