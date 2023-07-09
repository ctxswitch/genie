package resources

import (
	"ctx.sh/genie/pkg/config"
	"ctx.sh/genie/pkg/resources/list"
)

func MockResources() *Resources {
	lists := map[string]Resource{
		"name":      list.FromConfig(config.List{"Jim Halpert"}),
		"greeting":  list.FromConfig(config.List{"Hello"}),
		"statement": list.FromConfig(config.List{"I'm sorry Mr. Buttlicker"}),
	}

	return &Resources{
		lists: lists,
	}
}
