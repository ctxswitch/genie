package variables

func MockVariables() *Variables {
	return &Variables{
		vars: map[string]string{
			"name": "Dwight Schrute",
		},
	}
}
