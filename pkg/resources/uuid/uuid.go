package uuid

import (
	"math/rand"

	"ctx.sh/genie/pkg/config"
	"github.com/google/uuid"
)

type UuidFunc func() string

type Uuid struct {
	uuidType string
	uniques  int
	fn       UuidFunc
	cache    []string
}

func FromConfig(options config.Uuid) *Uuid {
	u := &Uuid{
		uniques: options.Uniques,
	}

	switch options.Type {
	case "uuid1":
		u.fn = u.uuid1
	default:
		u.fn = u.uuid4
	}

	return u
}

func (u *Uuid) Cache() []string {
	cache := make([]string, u.uniques)

	for i := 0; i < u.uniques; i++ {
		cache = append(cache, u.fn())
	}

	u.fn = u.cached
	return cache
}

func (u *Uuid) Get() string {
	if u.cache == nil && u.uniques > 0 {
		u.cache = u.Cache()
	}

	return u.fn()
}

func (u *Uuid) uuid1() string {
	id, e := uuid.NewUUID()
	if e != nil {
		return "deadbeef-0000-0000-0000-000000000000"
	}

	return id.String()
}

func (u *Uuid) uuid4() string {
	return uuid.NewString()
}

func (u *Uuid) cached() string {
	return u.cache[rand.Intn(len(u.cache))]
}
