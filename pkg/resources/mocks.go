package resources

import "stvz.io/genie/pkg/resources/list"

func MockResources() *Resources {
	lists := map[string]Resource{
		"name":      list.New(list.Config{"Jim Halpert"}),
		"greeting":  list.New(list.Config{"Hello"}),
		"statement": list.New(list.Config{"I'm sorry Mr. Buttlicker"}),
	}

	return &Resources{
		Lists: lists,
	}
}
