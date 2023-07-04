package resources

import "ctx.sh/genie/pkg/resources/list"

func MockResources() *Resources {
	lists := map[string]Resource{
		"name":      list.List{"Jim Halpert"},
		"greeting":  list.List{"Hello"},
		"statement": list.List{"I'm sorry Mr. Buttlicker"},
	}

	return &Resources{
		lists: lists,
	}
}
