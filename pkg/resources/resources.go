package resources

import (
	"ctx.sh/genie/pkg/resources/integer_range"
	"ctx.sh/genie/pkg/resources/list"
	"ctx.sh/genie/pkg/resources/random_string"
	"ctx.sh/genie/pkg/resources/timestamp"
	"ctx.sh/genie/pkg/resources/uuid"

	"ctx.sh/genie/pkg/config"
)

type Resource interface {
	Get() string
}

type Resources struct {
	lists         map[string]Resource
	integerRanges map[string]Resource
	randomStrings map[string]Resource
	uuids         map[string]Resource
	timestamps    map[string]Resource
	maps          map[string]Resource
}

func NewResources() *Resources {
	return &Resources{
		lists:         make(map[string]Resource),
		integerRanges: make(map[string]Resource),
		randomStrings: make(map[string]Resource),
		uuids:         make(map[string]Resource),
		timestamps:    make(map[string]Resource),
		maps:          make(map[string]Resource),
	}
}

func FromConfig(cfg *config.Configs) (*Resources, error) {
	r := NewResources()

	for k := range cfg.Resources.Lists {
		r.lists[k] = list.FromConfig(cfg.Resources.Lists[k])
	}

	for k := range cfg.Resources.IntegerRanges {
		r.integerRanges[k] = integer_range.FromConfig(cfg.Resources.IntegerRanges[k])
	}

	for k := range cfg.Resources.RandomStrings {
		r.randomStrings[k] = random_string.FromConfig(cfg.Resources.RandomStrings[k])
	}

	for k := range cfg.Resources.Timestamps {
		r.timestamps[k] = timestamp.FromConfig(cfg.Resources.Timestamps[k])
	}

	for k := range cfg.Resources.Uuids {
		r.uuids[k] = uuid.FromConfig(cfg.Resources.Uuids[k])
	}

	return r, nil
}

func (r *Resources) Get(rtype string, name string) (Resource, error) {
	switch rtype {
	case "list":
		return r.GetList(name)
	case "integer_range":
		return r.GetIntegerRange(name)
	case "random_string":
		return r.GetRandomString(name)
	case "uuid":
		return r.GetUuid(name)
	case "timestamp":
		return r.GetTimestamp(name)
	case "map":
		return r.GetMap(name)
	default:
		return nil, InvalidResourceTypeError
	}
}

func (r *Resources) MustGet(rtype string, name string) Resource {
	res, err := r.Get(rtype, name)
	if err != nil {
		panic(err)
	}

	return res
}

func (r *Resources) GetList(name string) (Resource, error) {
	if resource, ok := r.lists[name]; ok {
		return resource, nil
	}

	return nil, NotFoundError
}

func (r *Resources) GetIntegerRange(name string) (Resource, error) {
	if resource, ok := r.integerRanges[name]; ok {
		return resource, nil
	}

	return nil, NotFoundError
}

func (r *Resources) GetRandomString(name string) (Resource, error) {
	if resource, ok := r.randomStrings[name]; ok {
		return resource, nil
	}

	return nil, NotFoundError
}

func (r *Resources) GetUuid(name string) (Resource, error) {
	if resource, ok := r.uuids[name]; ok {
		return resource, nil
	}

	return nil, NotFoundError
}

func (r *Resources) GetTimestamp(name string) (Resource, error) {
	if resource, ok := r.timestamps[name]; ok {
		return resource, nil
	}

	return nil, NotFoundError
}

func (r *Resources) GetMap(name string) (Resource, error) {
	if resource, ok := r.maps[name]; ok {
		return resource, nil
	}

	return nil, NotFoundError
}
