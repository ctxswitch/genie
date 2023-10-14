package uuid

import (
	"math/rand"

	"github.com/google/uuid"
)

type Func func() string

type UUID struct {
	uniques int
	fn      Func
	cache   []string
}

func New(cfg Config) *UUID {
	u := &UUID{
		uniques: cfg.Uniques,
	}

	u.setFn(cfg.Type)
	return u
}

func (u *UUID) setFn(t string) {
	switch t {
	case "uuid1":
		u.fn = u.uuid1
	default:
		u.fn = u.uuid4
	}
}

func (u *UUID) Cache() []string {
	cache := make([]string, u.uniques)

	for i := 0; i < u.uniques; i++ {
		cache = append(cache, u.fn())
	}

	u.fn = u.cached
	return cache
}

func (u *UUID) Get() string {
	if u.cache == nil && u.uniques > 0 {
		u.cache = u.Cache()
	}

	return u.fn()
}

func (u *UUID) uuid1() string {
	id, e := uuid.NewUUID()
	if e != nil {
		// TODO: rethink this
		return "deadbeef-0000-0000-0000-000000000000"
	}

	return id.String()
}

func (u *UUID) uuid4() string {
	return uuid.NewString()
}

func (u *UUID) cached() string {
	return u.cache[rand.Intn(len(u.cache))] //nolint:gosec
}
