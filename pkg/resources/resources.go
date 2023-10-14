package resources

import (
	"stvz.io/genie/pkg/resources/integer_range"
	"stvz.io/genie/pkg/resources/ipaddr"
	"stvz.io/genie/pkg/resources/list"
	"stvz.io/genie/pkg/resources/random_string"
	"stvz.io/genie/pkg/resources/timestamp"
	"stvz.io/genie/pkg/resources/uuid"
)

type Resource interface {
	Get() string
}

type Resources struct {
	Lists         map[string]Resource
	IntegerRanges map[string]Resource
	RandomStrings map[string]Resource
	UUIDs         map[string]Resource
	Timestamps    map[string]Resource
	Maps          map[string]Resource
	IPAddrs       map[string]Resource
}

func New(block Config) *Resources {
	return &Resources{
		IntegerRanges: parseIntegerRanges(block),
		RandomStrings: parseRandomStrings(block),
		Lists:         parseLists(block),
		Timestamps:    parseTimestamps(block),
		UUIDs:         parseUUIDs(block),
		IPAddrs:       parseIPAddrs(block),
	}
}

func parseIntegerRanges(res Config) map[string]Resource {
	out := make(map[string]Resource)

	for k, v := range res.IntegerRanges {
		out[k] = integer_range.New(v)
	}

	return out
}

func parseRandomStrings(res Config) map[string]Resource {
	out := make(map[string]Resource)

	for k, v := range res.RandomStrings {
		out[k] = random_string.New(v)
	}

	return out
}

func parseLists(res Config) map[string]Resource {
	out := make(map[string]Resource)

	for k, v := range res.Lists {
		out[k] = list.New(v)
	}

	return out
}

func parseTimestamps(res Config) map[string]Resource {
	out := make(map[string]Resource)

	for k, v := range res.Timestamps {
		out[k] = timestamp.New(v)
	}

	return out
}

func parseUUIDs(res Config) map[string]Resource {
	out := make(map[string]Resource)

	for k, v := range res.UUIDs {
		out[k] = uuid.New(v)
	}

	return out
}

func parseIPAddrs(res Config) map[string]Resource {
	out := make(map[string]Resource)

	for k, v := range res.IPAddrs {
		out[k] = ipaddr.New(v)
	}

	return out
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
		return r.GetUUID(name)
	case "timestamp":
		return r.GetTimestamp(name)
	case "map":
		return r.GetMap(name)
	case "ipaddr":
		return r.GetIPAddr(name)
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
	if resource, ok := r.Lists[name]; ok {
		return resource, nil
	}

	return nil, NotFoundError
}

func (r *Resources) GetIntegerRange(name string) (Resource, error) {
	if resource, ok := r.IntegerRanges[name]; ok {
		return resource, nil
	}

	return nil, NotFoundError
}

func (r *Resources) GetRandomString(name string) (Resource, error) {
	if resource, ok := r.RandomStrings[name]; ok {
		return resource, nil
	}

	return nil, NotFoundError
}

func (r *Resources) GetUUID(name string) (Resource, error) {
	if resource, ok := r.UUIDs[name]; ok {
		return resource, nil
	}

	return nil, NotFoundError
}

func (r *Resources) GetTimestamp(name string) (Resource, error) {
	if resource, ok := r.Timestamps[name]; ok {
		return resource, nil
	}

	return nil, NotFoundError
}

func (r *Resources) GetMap(name string) (Resource, error) {
	if resource, ok := r.Maps[name]; ok {
		return resource, nil
	}

	return nil, NotFoundError
}

func (r *Resources) GetIPAddr(name string) (Resource, error) {
	if resource, ok := r.IPAddrs[name]; ok {
		return resource, nil
	}

	return nil, NotFoundError
}
