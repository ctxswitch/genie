package random_string // nolint:revive

import (
	"math/rand"
)

type RandomString struct {
	size    uint32
	chars   []rune
	uniques uint32
	cache   []string
}

func New(cfg Config) *RandomString {
	return &RandomString{
		size:    cfg.Size,
		chars:   cfg.Chars,
		uniques: cfg.Uniques,
	}
}

func (r *RandomString) Cache() []string {
	c := make([]string, r.uniques)

	for i := range c {
		c[i] = r.randomizer()
	}

	return c
}

func (r *RandomString) Get() string {
	if r.uniques > 0 {
		if r.cache == nil {
			r.cache = r.Cache()
		}

		return r.cache[rand.Intn(len(r.cache))] //nolint:gosec
	}

	return r.randomizer()
}

func (r *RandomString) randomizer() string {
	c := make([]rune, r.size)
	for i := range c {
		c[i] = r.chars[rand.Intn(len(r.chars))] //nolint:gosec
	}

	return string(c)
}
