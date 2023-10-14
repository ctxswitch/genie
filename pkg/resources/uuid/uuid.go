package uuid

import (
	"math/rand"

	"github.com/google/uuid"
)

// Func is a function that returns a string.
type Func func() string

// UUID is a resource that generates UUIDs.
type UUID struct {
	uniques int
	fn      Func
	cache   []string
}

// New returns a new UUID resource.
func New(cfg Config) *UUID {
	u := &UUID{
		uniques: cfg.Uniques,
	}

	u.setFn(cfg.Type)
	return u
}

// setFn sets the UUID generation function based on the string
// representation.
func (u *UUID) setFn(t string) {
	switch t {
	case "uuid1":
		u.fn = u.uuid1
	default:
		u.fn = u.uuid4
	}
}

// Cache generates a cache of unique UUIDs.
func (u *UUID) Cache() []string {
	cache := make([]string, u.uniques)

	for i := 0; i < u.uniques; i++ {
		cache = append(cache, u.fn())
	}

	u.fn = u.cached
	return cache
}

// Get implements the Resource interface and returns a random uuid.
func (u *UUID) Get() string {
	if u.cache == nil && u.uniques > 0 {
		u.cache = u.Cache()
	}
	return u.fn()
}

// uuid1 generates a UUIDv1.
func (u *UUID) uuid1() string {
	id, e := uuid.NewUUID()
	if e != nil {
		// TODO: rethink this
		return "deadbeef-0000-0000-0000-000000000000"
	}

	return id.String()
}

// uuid4 generates a UUIDv4.
func (u *UUID) uuid4() string {
	return uuid.NewString()
}

// cached returns a random UUID from the cache.
func (u *UUID) cached() string {
	return u.cache[rand.Intn(len(u.cache))] //nolint:gosec
}
