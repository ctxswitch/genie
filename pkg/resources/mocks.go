package resources

func MockResources() *Resources {
	lists := map[string]Resource{
		"name": nil,
	}

	return &Resources{
		lists: lists,
	}
}
