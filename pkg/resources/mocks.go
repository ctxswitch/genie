package resources

import (
	"ctx.sh/genie/pkg/config"
	"ctx.sh/genie/pkg/resources/list"
)

func MockResources() *Resources {
	lists := map[string]Resource{
		"name":      list.New(config.ListBlock{"Jim Halpert"}),
		"greeting":  list.New(config.ListBlock{"Hello"}),
		"statement": list.New(config.ListBlock{"I'm sorry Mr. Buttlicker"}),
	}

	return &Resources{
		Lists: lists,
	}
}
