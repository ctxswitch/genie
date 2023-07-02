package resources

import "ctx.sh/genie/pkg/resources/list"

func MockResources() *Resources {
	lists := map[string]Resource{
		"name":     list.List{"Jim Halpert"},
		"greeting": list.List{"Hello"},
	}

	return &Resources{
		lists: lists,
	}
}
