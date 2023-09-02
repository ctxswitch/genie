package random_string

import (
	"math/rand"

	"ctx.sh/genie/pkg/config"
)

type RandomString struct {
	size    uint32
	chars   []rune
	uniques uint32
	cache   []string
}

func New(settings config.RandomStringBlock) *RandomString {
	return &RandomString{
		size:    settings.Size,
		chars:   settings.Chars,
		uniques: settings.Uniques,
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
