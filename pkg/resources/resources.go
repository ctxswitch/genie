package resources

import (
	"stvz.io/genie/pkg/resources/integer_range"
	"stvz.io/genie/pkg/resources/ipaddr"
	"stvz.io/genie/pkg/resources/list"
	"stvz.io/genie/pkg/resources/random_string"
	"stvz.io/genie/pkg/resources/timestamp"
	"stvz.io/genie/pkg/resources/uuid"
)

// Resource defines the interface all resources must adhear to.
type Resource interface {
	Get() string
}

// Resources contains a map of all configured resources.
type Resources struct {
	Lists         map[string]Resource
	IntegerRanges map[string]Resource
	RandomStrings map[string]Resource
	UUIDs         map[string]Resource
	Timestamps    map[string]Resource
	Maps          map[string]Resource
	IPAddrs       map[string]Resource
}

// New returns a new collection of resources.
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

// parseIntegerRanges parses a map of integer ranges into a map of resources.
func parseIntegerRanges(res Config) map[string]Resource {
	out := make(map[string]Resource)

	for k, v := range res.IntegerRanges {
		out[k] = integer_range.New(v)
	}

	return out
}

// parseRandomStrings parses a map of random strings into a map of resources.
func parseRandomStrings(res Config) map[string]Resource {
	out := make(map[string]Resource)

	for k, v := range res.RandomStrings {
		out[k] = random_string.New(v)
	}

	return out
}

// parseLists parses a map of lists into a map of resources.
func parseLists(res Config) map[string]Resource {
	out := make(map[string]Resource)

	for k, v := range res.Lists {
		out[k] = list.New(v)
	}

	return out
}

// parseTimestamps parses a map of timestamps into a map of resources.
func parseTimestamps(res Config) map[string]Resource {
	out := make(map[string]Resource)

	for k, v := range res.Timestamps {
		out[k] = timestamp.New(v)
	}

	return out
}

// parseUUIDs parses a map of UUIDs into a map of resources.
func parseUUIDs(res Config) map[string]Resource {
	out := make(map[string]Resource)

	for k, v := range res.UUIDs {
		out[k] = uuid.New(v)
	}

	return out
}

// parseIPAddrs parses a map of IP addresses into a map of resources.
func parseIPAddrs(res Config) map[string]Resource {
	out := make(map[string]Resource)

	for k, v := range res.IPAddrs {
		out[k] = ipaddr.New(v)
	}

	return out
}

// Get returns a resource by type and name.
// TODO: we could probably just pass in the full resource name
// complete with the delimited name and type.  Something similar
// was done for the sinks.
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

// GetList returns a list resource by name.
func (r *Resources) GetList(name string) (Resource, error) {
	if resource, ok := r.Lists[name]; ok {
		return resource, nil
	}

	return nil, NotFoundError
}

// GetIntegerRange returns an integer range resource by name.
func (r *Resources) GetIntegerRange(name string) (Resource, error) {
	if resource, ok := r.IntegerRanges[name]; ok {
		return resource, nil
	}

	return nil, NotFoundError
}

// GetRandomString returns a random string resource by name.
func (r *Resources) GetRandomString(name string) (Resource, error) {
	if resource, ok := r.RandomStrings[name]; ok {
		return resource, nil
	}

	return nil, NotFoundError
}

// GetUUID returns a UUID resource by name.
func (r *Resources) GetUUID(name string) (Resource, error) {
	if resource, ok := r.UUIDs[name]; ok {
		return resource, nil
	}

	return nil, NotFoundError
}

// GetTimestamp returns a timestamp resource by name.
func (r *Resources) GetTimestamp(name string) (Resource, error) {
	if resource, ok := r.Timestamps[name]; ok {
		return resource, nil
	}

	return nil, NotFoundError
}

// GetMap returns a map resource by name.
func (r *Resources) GetMap(name string) (Resource, error) {
	if resource, ok := r.Maps[name]; ok {
		return resource, nil
	}

	return nil, NotFoundError
}

// GetIPAddr returns an IP address resource by name.
func (r *Resources) GetIPAddr(name string) (Resource, error) {
	if resource, ok := r.IPAddrs[name]; ok {
		return resource, nil
	}

	return nil, NotFoundError
}
